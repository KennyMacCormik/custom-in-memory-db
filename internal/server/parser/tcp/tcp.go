package tcp

import (
	"bufio"
	"custom-in-memory-db/internal/server/parser"
	"fmt"
	"log/slog"
	"net"
	"time"
)

type TcpParser struct {
	listener net.Listener
	deadline time.Duration
}

func (t *TcpParser) New(ip, port string, deadline time.Duration) error {
	var err error
	t.listener, err = net.Listen("tcp4", ip+":"+port)
	if err != nil {
		return fmt.Errorf("tcp listener init error: %w", err)
	}
	t.deadline = deadline
	return nil
}

func (t *TcpParser) Close() {
	t.listener.Close()
}

// Read reads TCP input and tries to compose it into valid Command struct
func (t *TcpParser) Read(vc []string, lg *slog.Logger) (parser.Command, error) {
	conn, err := t.listener.Accept()
	if err != nil {
		lg.Error("connection handling error", "error", err.Error())
		return parser.Command{}, err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(t.deadline))
	// how to unit-test this????
	if err != nil {
		lg.Error("connection deadline cannot be set", "error", err.Error())
		return parser.Command{}, err
	}

	lg.Debug("connection accepted", "ip", conn.RemoteAddr().String())
	lg.Debug("conn deadline set", "value", t.deadline.String())

	r := bufio.NewReader(conn)
	result, err := parser.BufferRead(r, vc, lg)
	if err != nil {
		return parser.Command{}, err
	}

	return result, nil
}
