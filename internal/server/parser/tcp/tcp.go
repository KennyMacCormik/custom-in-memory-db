package tcp

import (
	"bufio"
	"custom-in-memory-db/internal/server/parser"
	"fmt"
	"io"
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

func (t *TcpParser) Close() error {
	return t.listener.Close()
}

func (t *TcpParser) Write(response string, wc io.WriteCloser, lg *slog.Logger) error {
	_, err := wc.Write([]byte(response))
	if err != nil {
		msg := "write response error"
		lg.Error(msg, "error", err.Error())
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}

// Read reads TCP input and tries to compose it into valid Command struct
func (t *TcpParser) Read(vc []string, lg *slog.Logger) (parser.Command, io.WriteCloser, error) {
	var msg string
	conn, err := t.listener.Accept()
	if err != nil {
		msg = "cannot accept connection"
		lg.Error(msg, "error", err.Error())
		return parser.Command{}, nil, fmt.Errorf("%s: %w", msg, err)
	}

	err = conn.SetDeadline(time.Now().Add(t.deadline))
	// how to unit-test this????
	if err != nil {
		msg = "connection deadline cannot be set"
		lg.Error(msg, "error", err.Error())
		return parser.Command{}, nil, fmt.Errorf("%s: %w", msg, err)
	}

	lg.Debug("connection accepted", "ip", conn.RemoteAddr().String())
	lg.Debug("conn deadline set", "value", t.deadline.String())

	r := bufio.NewReader(conn)
	result, err := parser.BufferRead(r, vc, lg)
	if err != nil {
		msg = "cannot read connection"
		lg.Error(msg, "error", err.Error())
		return parser.Command{}, nil, err
	}

	return result, conn, nil
}
