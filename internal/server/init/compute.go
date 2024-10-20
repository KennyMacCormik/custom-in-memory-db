package init

import (
	"custom-in-memory-db/internal/server/db/compute"
	"custom-in-memory-db/internal/server/db/storage"
	"log/slog"
)

func Compute(st storage.Storage, lg *slog.Logger) compute.Compute {
	lg.Info("compute init done")
	return compute.New(st)
}
