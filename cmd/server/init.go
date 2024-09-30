package main

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	_map "custom-in-memory-db/internal/server/db/storage/map"
	"custom-in-memory-db/internal/server/db/wal"
	"custom-in-memory-db/internal/server/tcp"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func initTcpServer(conf cmd.Config, lg *slog.Logger) *tcp.Server {
	srv := tcp.Server{}
	err := srv.New(conf.Net.NET_ADDR, strconv.Itoa(conf.Net.NET_PORT), conf.Net.NET_TIMEOUT, conf.Net.NET_MAX_CONN, lg)
	if err != nil {
		lg.Error("failed init tcp server", "error", err.Error())
		os.Exit(errExit)
	}

	return &srv
}

func initLogger(conf cmd.Config) *slog.Logger {
	logLevelMap := map[string]slog.Level{
		"debug": -4,
		"info":  0,
		"warn":  4,
		"error": 8,
	}

	var logLevel = new(slog.LevelVar)
	logLevel.Set(logLevelMap[conf.Log.LOG_LEVEL])

	if conf.Log.LOG_FORMAT == "text" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
}

func initStorage(conf cmd.Config, lg *slog.Logger) storage.Storage {
	switch conf.Engine.APP_STORAGE {
	case "mem":
		st := _map.MapStorage{}
		st.New()
		lg.Info("storage init done")

		return &st
	case "wal":
		st := _map.MapStorage{}
		st.New()
		lg.Info("storage init done")

		wl := wal.Wal{}
		err := wl.New(conf, &st)
		if err != nil {
			str := fmt.Errorf("wal init error: %w", err).Error()
			lg.Error("wal init error", "error", str)
			panic(str)
		}
		lg.Info("wal init done")

		return &wl
	default:
		str := fmt.Sprintf("unknown storage: %s", conf.Engine.APP_STORAGE)
		lg.Error("storage init error", "error", str)
		panic(str)
	}
}
