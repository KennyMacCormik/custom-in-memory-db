package compute

import (
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	"fmt"
)

type Compute interface {
	ControlLoop(p parser.Parser, st storage.Storage)
}

type Comp struct {
	validCommands []string
}

func (c *Comp) New() {
	c.validCommands = []string{"GET", "SET", "DEL", "QUIT", "EXIT"}
}

// how to test?
func (c *Comp) ControlLoop(p parser.Parser, st storage.Storage) {
	for {
		comm, err := p.Read(c.validCommands)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch comm.Command {
		case "GET":
			r, err := st.Get(comm)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(r)
		case "SET":
			err := st.Set(comm)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "DEL":
			err := st.Del(comm)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "EXIT", "QUIT":
			return
		}
	}
}
