package compute

import (
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/internal/server/db/storage"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

const defaultOk = "OK\n"

type Compute interface {
	Exec(cmd parser.Command, lg *slog.Logger) (string, error)
	Close() error
}

// Comp is an instance of the Compute interface
type Comp struct {
	st storage.Storage
}

// New initializes Comp with storage interface it will be working with
func New(st storage.Storage) Compute {
	return &(Comp{st: st})
}

func (c *Comp) Close() error {
	closer, ok := c.st.(io.Closer)
	if ok {
		return closer.Close()
	}

	return nil
}

func (c *Comp) Exec(cmd parser.Command, lg *slog.Logger) (string, error) {
	switch cmd.Command {
	case "GET":
		r, err := c.st.Get(cmd.Arg1)
		if err != nil {
			return "", fmt.Errorf("error getting value: %v", err)
		}
		return r, nil
	case "SET":
		err := c.st.Set(cmd.Arg1, cmd.Arg2)
		if err != nil {
			return "", err
		}
		return defaultOk, nil
	case "DEL":
		err := c.st.Del(cmd.Arg1)
		if err != nil {
			return "", fmt.Errorf("error deleting value: %v", err)
		}
		return defaultOk, nil
	default:
		return "", errors.New("unknown command")
	}
}
