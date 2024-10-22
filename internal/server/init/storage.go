package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/storage"
	_map "custom-in-memory-db/internal/server/db/storage/map"
	"custom-in-memory-db/internal/server/db/storage/wal"
	"fmt"
	"log/slog"
)

func initMapStorage(lg *slog.Logger) (storage.Storage, error) {
	lg.Info("map storage init done")
	return _map.New(), nil
}

func initWalStorage(conf cmd.Config, st storage.Storage, lg *slog.Logger) (storage.Storage, error) {
	wl, err := wal.New(conf, st, lg)
	if err != nil {
		return nil, err
	}
	lg.Info("wal storage init done")
	return wl, nil
}

func Storage(conf cmd.Config, lg *slog.Logger) (storage.Storage, error) {
	const suf = "init.Storage()"
	switch conf.Engine.Type {
	case "map":
		return initMapStorage(lg)
	case "wal":
		st, _ := initMapStorage(lg)
		return initWalStorage(conf, st, lg)
	default:
		return nil, fmt.Errorf("%s failed: unknown engine type %s", suf, conf.Engine.Type)
	}
}
