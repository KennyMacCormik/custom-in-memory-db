package db

import (
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/db/parser"
	"io"
	"log/slog"
)

type Database struct {
	comp compute.Compute
}

func (d *Database) New(comp compute.Compute) {
	d.comp = comp
}

func (d *Database) HandleRequest(r io.Reader, lg *slog.Logger) (string, error) {
	cmd, err := parser.Read(r, lg)
	if err != nil {
		lg.Error("parsing error", "error", err.Error())
		return "", err
	}

	result, err := d.comp.Exec(cmd)
	if err != nil {
		lg.Error("executing error", "error", err.Error())
		return "", err
	}

	return result, nil
}
