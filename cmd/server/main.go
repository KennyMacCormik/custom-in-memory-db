package main

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db"
	"custom-in-memory-db/internal/server/db/compute"
	myinit "custom-in-memory-db/internal/server/init"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const errExit = 1

func main() {
	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		panic(fmt.Errorf("config init error: %w", err))
	}

	lg := myinit.Logger(conf)
	lg.Info("config init success")

	st, err := myinit.Storage(conf, lg)
	if err != nil {
		os.Exit(errExit)
	}
	// how to mute "unhandled error" warning?
	defer st.Close()

	// Init compute layer
	comp := compute.Comp{}
	comp.New(st)
	lg.Info("compute init done")

	// Init db layer
	database := db.Database{}
	database.New(&comp)

	srv, err := myinit.TcpServer(conf, lg)
	if err != nil {
		os.Exit(errExit)
	}
	// how to mute "unhandled error" warning?
	defer srv.Close()
	lg.Info("tcp server init done")

	handler := func(r io.Reader, lg *slog.Logger) (string, error) {
		return database.HandleRequest(r, lg)
	}
	lg.Info("listening")
	srv.Listen(handler)
}
