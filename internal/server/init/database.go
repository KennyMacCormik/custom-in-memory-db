package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db"
	"custom-in-memory-db/internal/server/network"
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

	net := initNetworkEndpoint(conf, lg)
	if net == nil {
		lg.Error("network init failed: unknown network type")
		os.Exit(errExit)
	}

	database := db.Database{}
	database.New(comp, net, pr, lg)
	lg.Info("db init done")

	return &database
}

func initNetworkEndpoint(conf cmd.Config, lg *slog.Logger) network.Endpoint {
	switch conf.Network.Endpoint {
	case "tcp":
		return TcpServer(conf, lg)
	case "http":
		return HttpServer(conf, lg)
	default:
		return nil
	}
}
