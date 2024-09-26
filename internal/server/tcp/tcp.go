package tcp

import (
	"custom-in-memory-db/internal/server/compute"
	"custom-in-memory-db/internal/server/parser"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"sync"
	"time"
)

type Server struct {
	listener net.Listener
	deadline time.Duration
	maxConn  int
	lg       *slog.Logger
}

type connMeter struct {
	currConn int
	maxConn  int
	mtx      sync.Mutex
	cond     *sync.Cond
}

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

func (s *Server) Listen(c compute.Compute) {
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
		go handleClient(conn, s.deadline, c, &cm, s.lg)
	}
}

func handleClient(conn net.Conn, deadline time.Duration, c compute.Compute, cm *connMeter, lg *slog.Logger) {
	defer cm.decConnCount()
	defer conn.Close()

	ilg := lg.With("ID", uuid.New(), "remoteAddr", conn.RemoteAddr().String())

	err := conn.SetDeadline(time.Now().Add(deadline))
	// how to unit-test this????
	if err != nil {
		ilg.Error("connection deadline cannot be set", "error", err.Error())
	}

	handleCommand(conn, c, lg)
}

func handleCommand(conn net.Conn, c compute.Compute, lg *slog.Logger) {
	cmd, err := parser.Read(conn, lg)
	if err != nil {
		lg.Error("parsing error", "error", err.Error())
		_, err = conn.Write([]byte(err.Error()))
		if err != nil {
			lg.Error("connection writing error", "error", err.Error())
		}
	}

	result, err := c.Exec(cmd)
	if err != nil {
		lg.Error("executing error", "error", err.Error())
		_, err = conn.Write([]byte(err.Error()))
		if err != nil {
			lg.Error("connection writing error", "error", err.Error())
		}
	}

	_, err = conn.Write([]byte(result))
	if err != nil {
		lg.Error("result writing error", "error", err.Error())
	}
}

// lock only allows further execution if goNum is less than goMax
func (c *connMeter) incConnCount() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	for c.currConn == c.maxConn {
		c.cond.Wait()
	}

	c.currConn++
}

// unlock decrements goNum and calls waiting main.go
func (c *connMeter) decConnCount() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	c.currConn--
	c.cond.Signal()
}
