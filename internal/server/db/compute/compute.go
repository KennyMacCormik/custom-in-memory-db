package compute

import (
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/internal/server/db/storage"
	"errors"
	"fmt"
	"log/slog"
)

const defaultOk = "OK"

type Compute interface {
	Exec(cmd parser.Command, lg *slog.Logger) (string, error)
}

// Comp is an instance of the Compute interface
type Comp struct {
	st storage.Storage
}

// New initializes Comp with storage interface it will be working with
func (c *Comp) New(st storage.Storage) {
	c.st = st
}

func (c *Comp) Exec(cmd parser.Command, lg *slog.Logger) (string, error) {
	switch cmd.Command {
	case "GET":
		r, err := c.st.Get(cmd.Args[0])
		if err != nil {
			return "", fmt.Errorf("error getting value: %v", err)
		}
		return r, nil
	case "SET":
		err := c.st.Set(cmd.Args[0], cmd.Args[1])
		if err != nil {
			return "", fmt.Errorf("error settings value: %v", err)
		}
		return defaultOk, nil
	case "DEL":
		err := c.st.Del(cmd.Args[0])
		if err != nil {
			return "", fmt.Errorf("error deleting value: %v", err)
		}
		return defaultOk, nil
	default:
		return "", errors.New("unknown command")
	}
}
