package compute

import (
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	"github.com/google/uuid"
	"io"
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
	var err error
	// unique id for a request
	lg = lg.With("ID", uuid.New())

	comm, wc, err := p.Read(c.validCommands, lg)
	if err != nil {
		lg.Error("failed reading command", "error", err.Error())
		// how to decouple wc from err var?
		if wc != nil {
			_, err1 := wc.Write([]byte(err.Error()))
			wc.Close()
			if err1 != nil {
				lg.Error("failed write to connection", "error", err1.Error())
				return err
			}
		}
		return err
	}

	defer wc.Close()

	err = callStorage(wc, comm, st, lg)
	if err != nil {
		lg.Error("failed calling storage", "error", err.Error())
		return err
	}

	return nil
}

func callStorage(wc io.Writer, comm parser.Command, st storage.Storage, lg *slog.Logger) error {
	switch comm.Command {
	case "GET":
		r, err := st.Get(comm)
		if err != nil {
			lg.Error("failed to execute GET request", "error", err.Error())
			_, err1 := wc.Write([]byte(err.Error()))
			if err1 != nil {
				lg.Error("failed write to connection", "error", err1.Error())
				return err
			}
			return err
		}

		_, err = wc.Write([]byte(r))
		if err != nil {
			lg.Error("failed write to connection", "error", err.Error())
			return err
		}
		return nil
	case "SET":
		err := st.Set(comm)
		if err != nil {
			lg.Error("failed to execute SET request", "error", err.Error())
			_, err1 := wc.Write([]byte(err.Error()))
			if err1 != nil {
				lg.Error("failed write to connection", "error", err1.Error())
				return err
			}
			return err
		}

		_, err = wc.Write([]byte("OK"))
		if err != nil {
			lg.Error("failed write to connection", "error", err.Error())
			return err
		}
		return nil
	case "DEL":
		err := st.Del(comm)
		if err != nil {
			lg.Error("failed to execute DEL request", "error", err.Error())
			_, err1 := wc.Write([]byte(err.Error()))
			if err1 != nil {
				lg.Error("failed write to connection", "error", err1.Error())
				return err
			}
			return err
		}

		_, err = wc.Write([]byte("OK"))
		if err != nil {
			lg.Error("failed write to connection", "error", err.Error())
			return err
		}
		return nil
	default:
		return nil
	}
}
