package compute

import (
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	"github.com/google/uuid"
	"log/slog"
)

// I need to return something to run loop. Why not this? I'm already logging everything.
const Ok = "ok"
const Stop = "Stop"
const ErrParse = "ErrParce"
const ErrGet = "ErrGet"
const ErrSet = "ErrSet"
const ErrDel = "ErrDel"
const ErrUnexpectedCommand = "UnexpectedCommand"

type Compute interface {
	HandleRequest(p parser.Parser, st storage.Storage, lg *slog.Logger) string
}

type Comp struct {
	validCommands []string
}

func (c *Comp) New() {
	c.validCommands = []string{"GET", "SET", "DEL", "QUIT", "EXIT"}
}

func (c *Comp) HandleRequest(p parser.Parser, st storage.Storage, lg *slog.Logger) string {
	// unique id for a request
	lg = lg.With("ID", uuid.New())

	comm, err := p.Read(c.validCommands, lg)
	if err != nil {
		lg.Error(ErrParse, "error", err.Error())
		return ErrParse
	}

	switch comm.Command {
	case "GET":
		r, err := st.Get(comm)
		if err != nil {
			lg.Error(ErrGet, "error", err.Error())
			return ErrGet
		}
		lg.Debug(Ok)
		return r
	case "SET":
		err := st.Set(comm)
		if err != nil {
			lg.Error(ErrSet, "error", err.Error())
			return ErrSet
		}
		lg.Debug(Ok)
		return Ok
	case "DEL":
		err := st.Del(comm)
		if err != nil {
			lg.Error(ErrDel, "error", err.Error())
			return ErrDel
		}
		lg.Debug(Ok)
		return Ok
	case "EXIT", "QUIT":
		return Stop
	}
	lg.Error(ErrUnexpectedCommand)
	return ErrUnexpectedCommand
}
