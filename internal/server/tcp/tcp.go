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

const goMax = 100

type Server struct {
	listener net.Listener
	deadline time.Duration

	lg *slog.Logger

	currConn int
	maxConn  int
	mtx      sync.Mutex
	cond     *sync.Cond
}

func (s *Server) New(ip, port string, deadline time.Duration, lg *slog.Logger) error {
	var err error

	s.listener, err = net.Listen("tcp4", ip+":"+port)
	if err != nil {
		return fmt.Errorf("tcp listener init error: %w", err)
	}
	s.deadline = deadline

	s.lg = lg

	s.currConn = 0
	s.maxConn = goMax
	s.cond = sync.NewCond(&s.mtx)
	return nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) Listen(c compute.Compute) {
	var msg string
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			msg = "cannot accept connection"
			s.lg.Error(msg, "error", err.Error())
		}

		s.incConnCount()
		go s.handleClient(conn, c)
	}
}

func (s *Server) handleClient(conn net.Conn, c compute.Compute) {
	defer s.decConnCount()
	defer conn.Close()

	lg := s.lg.With("ID", uuid.New(), "remoteAddr", conn.RemoteAddr().String())

	err := conn.SetDeadline(time.Now().Add(s.deadline))
	// how to unit-test this????
	if err != nil {
		lg.Error("connection deadline cannot be set", "error", err.Error())
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
func (s *Server) incConnCount() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.currConn == s.maxConn {
		s.cond.Wait()
	}

	s.currConn++
}

// unlock decrements goNum and calls waiting main.go
func (s *Server) decConnCount() {
	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.currConn--
	s.cond.Signal()
}
