package init

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

const errExit = 1

func TcpServer(conf cmd.Config, lg *slog.Logger) *tcp.Server {
	srv := tcp.Server{}
	err := srv.New(conf.Net.NET_ADDR, strconv.Itoa(conf.Net.NET_PORT), conf.Net.NET_TIMEOUT, conf.Net.NET_MAX_CONN, lg)
	if err != nil {
		lg.Error("failed init tcp server", "error", err.Error())
		os.Exit(errExit)
	}

	return &srv
}

func Logger(conf cmd.Config) *slog.Logger {
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

func Storage(conf cmd.Config, lg *slog.Logger) storage.Storage {
	switch conf.Engine.APP_STORAGE {
	case "mem":
		return initMapStorage(lg)
	case "wal":
		st := initMapStorage(lg)
		writer := initWriter(conf, lg)
		wl := initWal(conf, st, writer, lg)

		return wl
	default:
		lg.Error("storage init error", "error",
			fmt.Sprintf("unknown storage: %s", conf.Engine.APP_STORAGE))
		os.Exit(errExit)
	}
	// why I need this?
	return nil
}

func initMapStorage(lg *slog.Logger) *_map.MapStorage {
	st := _map.MapStorage{}
	st.New()
	lg.Info("storage init done")

	return &st
}

func initWriter(conf cmd.Config, lg *slog.Logger) *wal.Writer {
	const suf = "wal.Writer.New() error"
	writer := wal.Writer{}
	if err := writer.New(conf); err != nil {
		lg.Error(suf, "error", err.Error())
		os.Exit(errExit)
	}
	lg.Info("writer init done")

	return &writer
}

func initWal(conf cmd.Config, st storage.Storage, writer wal.WriterInterface, lg *slog.Logger) *wal.Wal {
	const suf = "wal.Wal{}.New() error"
	wl := wal.Wal{}
	err := wl.New(conf, st, writer)
	if err != nil {
		lg.Error(suf, "error", err.Error())
		os.Exit(errExit)
	}
	lg.Info("wal init done")

	return &wl
}
