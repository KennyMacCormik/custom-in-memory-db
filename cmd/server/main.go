package main

import (
	"custom-in-memory-db/internal/server/cmd"
	db2 "custom-in-memory-db/internal/server/db"
	"custom-in-memory-db/internal/server/db/compute"
	"fmt"
	"io"
	"log/slog"
)

const errExit = 1

func main() {
	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		panic(fmt.Errorf("config init error: %w", err))
	}

	lg := initLogger(conf)
	lg.Info("config init success")

	st := initStorage(conf, lg)
	defer st.Close()

	// Init compute layer
	comp := compute.Comp{}
	comp.New(st)
	lg.Info("compute init done")

	// Init db layer
	db := db2.Database{}
	db.New(&comp)
	lg.Info("db init done")

	srv := intiTcpServer(conf, lg)
	defer srv.Close()
	lg.Info("tcp server init done")

	handler := func(r io.Reader, lg *slog.Logger) (string, error) {
		result, err := db.HandleRequest(r, lg)
		return result, err
	}
	lg.Info("listening")
	srv.Listen(handler)
}
