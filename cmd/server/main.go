package main

import (
	"custom-in-memory-db/internal/server/cmd"
	myinit "custom-in-memory-db/internal/server/init"
	"errors"
	"os"
)

const errExit = 1

func main() {
	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		lg := myinit.Logger(conf)
		lg.Error("config init error", "error", errors.Unwrap(err).Error())
		os.Exit(errExit)
	}

	lg := myinit.Logger(conf)
	lg.Info("config init success")

	db := myinit.Database(conf, lg)
	defer db.Close()
	db.ListenClient()
}
