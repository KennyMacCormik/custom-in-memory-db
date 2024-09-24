package compute

import (
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	"github.com/google/uuid"
	"log/slog"
)

type Compute interface {
	HandleRequest(p parser.Parser, st storage.Storage, lg *slog.Logger) error
}

type Comp struct {
	validCommands []string
}

func (c *Comp) New() {
	c.validCommands = []string{"GET", "SET", "DEL"}
}

func (c *Comp) HandleRequest(p parser.Parser, st storage.Storage, lg *slog.Logger) error {
	var r string
	var err error
	// unique id for a request
	lg = lg.With("ID", uuid.New())

	comm, wc, err := p.Read(c.validCommands, lg)
	if err != nil {
		lg.Error("failed reading command", "error", err.Error())
		return err
	}
	defer wc.Close()

	switch comm.Command {
	case "GET":
		r, err = st.Get(comm)
		if err != nil {
			lg.Error("failed to execute GET request", "error", err.Error())
			// Is it safe to ignore Write error?
			_ = p.Write(err.Error(), wc, lg)
			return err
		}
		// Is it safe to ignore Write error?
		_ = p.Write(r, wc, lg)
		return nil
	case "SET":
		err = st.Set(comm)
		if err != nil {
			lg.Error("failed to execute SET request", "error", err.Error())
			// Is it safe to ignore Write error?
			_ = p.Write(err.Error(), wc, lg)
			return err
		}
		// Is it safe to ignore Write error?
		_ = p.Write("OK", wc, lg)
		return nil
	case "DEL":
		err = st.Del(comm)
		if err != nil {
			lg.Error("failed to execute DEL request", "error", err.Error())
			// Is it safe to ignore Write error?
			_ = p.Write(err.Error(), wc, lg)
			return err
		}
		// Is it safe to ignore Write error?
		_ = p.Write("OK", wc, lg)
		return nil
	default:
		return nil
	}
}
