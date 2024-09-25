package compute

import (
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	"fmt"
)

const defaultOk = "OK"

type Compute interface {
	Exec(cmd parser.Command) (string, error)
}

// Instance of the Compute interface
type Comp struct {
	st storage.Storage
}

// New initializes Comp with storage interface it will be working with
func (c *Comp) New(st storage.Storage) {
	c.st = st
}

func (c *Comp) Exec(cmd parser.Command) (string, error) {
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
	}

	return "", nil
}
