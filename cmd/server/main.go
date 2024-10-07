package main

import (
	"custom-in-memory-db/internal/server/cmd"
	db2 "custom-in-memory-db/internal/server/db"
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/init"
	"fmt"
	"io"
	"log/slog"
)

func main() {
	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		panic(fmt.Errorf("config init error: %w", err))
	}

	lg := init.Logger(conf)
	lg.Info("config init success")

	st := init.Storage(conf, lg)
	// how to mute "unhandled error" warning?
	defer st.Close()
	if conf.Wal.WAL_SEG_RECOVER {
		err = st.Recover(conf, lg)
		if err != nil {
			panic(err)
		}
	}

	// Init compute layer
	comp := compute.Comp{}
	comp.New(st)
	lg.Info("compute init done")

	// Init db layer
	db := db2.Database{}
	db.New(&comp)
	lg.Info("db init done")

	srv := init.TcpServer(conf, lg)
	// how to mute "unhandled error" warning?
	defer srv.Close()
	lg.Info("tcp server init done")

	handler := func(r io.Reader, lg *slog.Logger) (string, error) {
		result, err := db.HandleRequest(r, lg)
		return result, err
	}
	lg.Info("listening")
	srv.Listen(handler)
}
