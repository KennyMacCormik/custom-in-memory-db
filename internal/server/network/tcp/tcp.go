package tcp

import (
	"custom-in-memory-db/internal/server/network"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"
)

// From https://pkg.go.dev/net#Listen.
// The network must be "tcp", "tcp4", "tcp6", "unix" or "unixpacket".
const listenNetwork = "tcp4"

type Server struct {
	listener net.Listener
	deadline time.Duration
	maxConn  int
	lg       *slog.Logger
}

func (c *connMeter) new(maxConn int) {
	c.maxConn = maxConn
	c.cond = sync.NewCond(&c.mtx)
}

func (s *Server) New(host, port string, deadline time.Duration, maxConn int, lg *slog.Logger) error {
	const suf = "TcpServer.New()"
	var err error
	address := strings.Join([]string{host, port}, ":")
	s.listener, err = net.Listen(listenNetwork, address)
	if err != nil {
		return fmt.Errorf("%s failed: %w", suf, err)
	}
	s.deadline = deadline
	s.maxConn = maxConn
	s.lg = lg

	return nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Listen(f network.Handler) {
	var msg string
	var cm connMeter
	if s.maxConn > 0 {
		cm.new(s.maxConn)
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			msg = "cannot accept connection"
			s.lg.Error(msg, "error", err.Error())
		}

		cm.incConnCount()
		go s.handleClient(conn, &cm, f, s.lg)
	}
}

func (s *Server) handleClient(conn net.Conn, cm *connMeter, handler network.Handler, lg *slog.Logger) {
	const suf = "server.handleClient()"
	defer cm.decConnCount()
	defer conn.Close()

	ilg := lg.With("ID", uuid.New(), "remoteAddr", conn.RemoteAddr().String())

	err := conn.SetDeadline(time.Now().Add(s.deadline))
	// how to unit-test this????
	if err != nil {
		ilg.Error(fmt.Sprintf("%s.SetDeadline()", suf), "error", err.Error())
	}

	ilg.Debug(fmt.Sprintf("%s conn opened", suf))

	result, err := handler(conn, ilg)
	ilg.Debug(fmt.Sprintf("%s", suf), "handlerResult", result)
	if err != nil {
		_, err = conn.Write([]byte(err.Error()))
		if err != nil {
			ilg.Error(fmt.Sprintf("%s.conn.Write()", suf), "error", err.Error())
			return
		}
		return
	}

	_, err = conn.Write([]byte(result))
	if err != nil {
		ilg.Error(fmt.Sprintf("%s.conn.Write()", suf), "error", err.Error())
	}
	ilg.Debug(fmt.Sprintf("%s", suf), "respondedToClient", "done")
}

type connMeter struct {
	currConn int
	maxConn  int
	mtx      sync.Mutex
	cond     *sync.Cond
}

// incConnCount only allows further execution if maxConn is less than goMax
func (c *connMeter) incConnCount() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	for c.currConn == c.maxConn {
		c.cond.Wait()
	}

	c.currConn++
}

// decConnCount decrements maxConn and calls waiting Server
func (c *connMeter) decConnCount() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	c.currConn--
	c.cond.Signal()
}
