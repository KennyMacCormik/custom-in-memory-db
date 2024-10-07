package db

import (
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/db/parser"
	"fmt"
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
	const suf = "database.HandleRequest()"
	cmd, err := parser.Read(r, lg)
	if err != nil {
		lg.Error(fmt.Sprintf("%s.parser.Read()", suf), "error", err.Error())
		return "", err
	}

	lg.Debug(fmt.Sprintf("%s.parser", suf), "result", fmt.Sprintf("%+v", cmd))

	result, err := d.comp.Exec(cmd, lg)
	if err != nil {
		lg.Error(fmt.Sprintf("%s.compute.Exec()", suf), "error", err.Error())
		return "", err
	}

	lg.Debug(fmt.Sprintf("%s.compute", suf), "result", fmt.Sprintf("%+v", result))

	return result, nil
}
