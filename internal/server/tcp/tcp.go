package tcp

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"
)

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

type Server struct {
	listener net.Listener
	deadline time.Duration
	maxConn  int
	lg       *slog.Logger
}

type Handler func(r io.Reader, lg *slog.Logger) (string, error)

func (c *connMeter) New(maxConn int) {
	c.maxConn = maxConn
	c.cond = sync.NewCond(&c.mtx)
}

func (s *Server) New(ip, port string, deadline time.Duration, maxConn int, lg *slog.Logger) error {
	var err error

	s.listener, err = net.Listen("tcp4", ip+":"+port)
	if err != nil {
		return fmt.Errorf("tcp listener init error: %w", err)
	}
	s.deadline = deadline
	s.maxConn = maxConn
	s.lg = lg

	return nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Listen(f Handler) {
	var msg string
	var cm connMeter
	if s.maxConn > 0 {
		cm.New(s.maxConn)
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

func (s *Server) handleClient(conn net.Conn, cm *connMeter, f Handler, lg *slog.Logger) {
	defer cm.decConnCount()
	defer conn.Close()

	ilg := lg.With("ID", uuid.New(), "remoteAddr", conn.RemoteAddr().String())

	err := conn.SetDeadline(time.Now().Add(s.deadline))
	// how to unit-test this????
	if err != nil {
		ilg.Error("connection deadline cannot be set", "error", err.Error())
	}

	result, err := f(conn, ilg)
	if err != nil {
		ilg.Error("executing error", "error", err.Error())
		_, err = conn.Write([]byte(err.Error()))
		if err != nil {
			ilg.Error("connection writing error", "error", err.Error())
		}
	}

	_, err = conn.Write([]byte(result))
	if err != nil {
		ilg.Error("result writing error", "error", err.Error())
	}
}
