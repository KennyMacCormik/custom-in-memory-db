package main

import (
	"custom-in-memory-db/internal/server/cmd"
	db2 "custom-in-memory-db/internal/server/db"
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/db/storage"
	"custom-in-memory-db/internal/server/db/storage/map"
	"custom-in-memory-db/internal/server/tcp"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

const errExit = 1

func main() {
	// Init logger
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		lg.Error("config init error", "error", err.Error())
		os.Exit(errExit)
	}

	// Init storage
	st := initStorage(conf, lg)

	// Init compute layer
	comp := compute.Comp{}
	comp.New(st)

	// Init db
	db := db2.Database{}
	db.New(&comp)

	// Init tcp server
	srv := tcp.Server{}
	err = srv.New(conf.Net.Address, strconv.Itoa(conf.Net.Port), conf.Net.IdleTimeout, conf.Net.MaxConn, lg)
	if err != nil {
		lg.Error("failed init tcp server", "error", err.Error())
		os.Exit(errExit)
	}
	defer srv.Close()

	srv.Listen(func(r io.Reader, lg *slog.Logger) (string, error) {
		result, err := db.HandleRequest(r, lg)
		return result, err
	})
}

func initStorage(conf cmd.Config, lg *slog.Logger) storage.Storage {
	switch conf.Eng.Storage {
	case "map":
		st := _map.MapStorage{}
		st.New()
		return &st
	default:
		str := fmt.Sprintf("unknown storage: %s", conf.Eng.Storage)
		lg.Error("storage init error", "error", str)
		panic(str)
	}
}
