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
	st := _map.Storage{}
	st.New()
	lg.Info("map storage init done")
	return &st, nil
}

func initWalStorage(conf cmd.Config, lg *slog.Logger) (storage.Storage, error) {
	st := wal.Storage{}
	pr := Parser(conf, lg)
	err := st.New(conf)
	if err == nil {
		if conf.Wal.Recover {
			err = st.Recover(conf, pr, lg)
		}
		lg.Info("wal storage init done")
	}
	return &st, err
}

func Storage(conf cmd.Config, lg *slog.Logger) (storage.Storage, error) {
	const suf = "init.Storage()"
	switch conf.Engine.Type {
	case "map":
		return initMapStorage(lg)
	case "wal":
		return initWalStorage(conf, lg)
	default:
		return nil, fmt.Errorf("%s failed: unknown engine type %s", suf, conf.Engine.Type)
	}
}
