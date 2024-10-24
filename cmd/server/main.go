package main

import (
	"custom-in-memory-db/internal/server/cmd"
	myinit "custom-in-memory-db/internal/server/init"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

const errExit = 1

func main() {
	// Init config
	conf, err := cmd.New()
	if err != nil {
		lg := myinit.Logger(conf)
		lg.Error("config init error", "error", errors.Unwrap(err).Error())
		os.Exit(errExit)
	}

	lg := myinit.Logger(conf)
	lg.Info("config init success")

	db := myinit.Database(conf, lg)
	defer db.Close()
	go db.ListenClient()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lg.Info("Shutdown Server ...")
}
