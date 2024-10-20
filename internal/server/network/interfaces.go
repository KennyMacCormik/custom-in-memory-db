package network

import (
	"io"
	"log/slog"
)

type Handler func(r io.Reader, lg *slog.Logger) (string, error)

type Endpoint interface {
	Listen(f Handler)
	Close() error
}
