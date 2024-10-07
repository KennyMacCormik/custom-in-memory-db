package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	_map "custom-in-memory-db/internal/server/db/storage/map"
	"custom-in-memory-db/internal/server/db/wal"
	"custom-in-memory-db/internal/server/tcp"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func TcpServer(conf cmd.Config, lg *slog.Logger) (*tcp.Server, error) {
	const suf = "init.TcpServer() failed"
	srv := tcp.Server{}
	err := srv.New(conf.Net.NET_ADDR, strconv.Itoa(conf.Net.NET_PORT), conf.Net.NET_TIMEOUT, conf.Net.NET_MAX_CONN, lg)
	if err != nil {
		lg.Error(suf, "error", err.Error())
		return nil, fmt.Errorf("%s: %w", suf, err)
	}
	lg.Debug("tcp server init done")

	return &srv, nil
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

func Storage(conf cmd.Config, lg *slog.Logger) (storage.Storage, error) {
	const suf = "init.TcpServer() failed"
	switch conf.Engine.APP_STORAGE {
	case "mem":
		return initMapStorage(lg), nil
	case "wal":
		st := initMapStorage(lg)

		writer, err := initWriter(conf, lg)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", suf, err)
		}

		wl := initWal(conf, st, writer, lg)

		if conf.Wal.WAL_SEG_RECOVER {
			err = wl.Recover(conf, lg)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", suf, err)
			}
		}

		wl.Start()

		return wl, nil
	default:
		errMsg := "unexpected APP_STORAGE"
		lg.Error(suf, "error", errMsg)
		return nil, fmt.Errorf("%s: %w", suf, errors.New(errMsg))
	}
}

func initMapStorage(lg *slog.Logger) *_map.MapStorage {
	st := _map.MapStorage{}
	st.New()
	lg.Info("map storage init done")

	return &st
}

func initWriter(conf cmd.Config, lg *slog.Logger) (*wal.Writer, error) {
	const suf = "wal.Writer.New() error"
	writer := wal.Writer{}
	if err := writer.New(conf); err != nil {
		lg.Error(suf, "error", err.Error())
		return nil, fmt.Errorf("%s: %w", suf, err)
	}
	lg.Info("writer init done")

	return &writer, nil
}

func initWal(conf cmd.Config, st storage.Storage, writer wal.WriterInterface, lg *slog.Logger) *wal.Wal {
	wl := wal.Wal{}
	wl.New(conf, st, writer)
	lg.Info("wal init done")

	return &wl
}
