package db

import (
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/internal/server/network"
	"errors"
	"fmt"
	"io"
	"log/slog"
)

type Database struct {
	comp        compute.Compute
	pr          parser.Parser
	netEndpoint network.Endpoint
	lg          *slog.Logger
}

func New(comp compute.Compute, netEndpoint network.Endpoint, pr parser.Parser, lg *slog.Logger) Database {
	return Database{comp: comp, netEndpoint: netEndpoint, pr: pr, lg: lg}
}

func (d *Database) Close() error {
	var err1, err2 error
	closer, ok := d.netEndpoint.(io.Closer)
	if ok {
		err1 = closer.Close()
		if err1 != nil {
			d.lg.Error("Database.Close().netEndpoint.Close() failed", "error", errors.Unwrap(err1).Error())
		}
	}

	err2 = d.comp.Close()
	if err2 != nil {
		d.lg.Error("Database.Close().Compute.Close() failed", "error", errors.Unwrap(err2).Error())
	}

	if err1 == nil && err2 == nil {
		return nil
	}
	return fmt.Errorf("Database.Close() failed")
}

func (d *Database) ListenClient() {
	handler := func(r io.Reader, lg *slog.Logger) (string, error) {
		result, err := d.HandleRequest(r, lg)
		return result, err
	}
	d.lg.Info("listening")
	d.netEndpoint.Listen(handler)
}

func (d *Database) HandleRequest(r io.Reader, lg *slog.Logger) (string, error) {
	const suf = "database.HandleRequest()"
	cmd, err := d.pr.Read(r, lg)
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
