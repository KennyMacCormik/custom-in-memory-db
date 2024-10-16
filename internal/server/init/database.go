package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db"
	"errors"
	"log/slog"
	"os"
)

const errExit = 1

func Database(conf cmd.Config, lg *slog.Logger) *db.Database {

	st, err := Storage(conf, lg)
	if err != nil {
		lg.Error("storage init failed", "error", errors.Unwrap(err).Error())
		os.Exit(errExit)
	}

	comp := Compute(st, lg)

	pr := Parser(conf, lg)

	tcp := TcpServer(conf, lg)

	database := db.Database{}
	database.New(comp, tcp, pr, lg)
	lg.Info("db init done")

	return &database
}
