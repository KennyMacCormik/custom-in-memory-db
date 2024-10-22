package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/network"
	"custom-in-memory-db/internal/server/network/tcp"
	"errors"
	"log/slog"
	"os"
	"strconv"
)

func TcpServer(conf cmd.Config, lg *slog.Logger) network.Endpoint {
	srv, err := tcp.New(conf.Network.Host, strconv.Itoa(conf.Network.Port), conf.Network.Timeout, conf.Network.MaxConn, lg)
	if err != nil {
		lg.Error("failed to init tcp server", "error", errors.Unwrap(err).Error())
		os.Exit(errExit)
	}
	lg.Info("tcp server init done")
	lg.Debug("tcp server params", "Host", conf.Network.Host,
		"Port", conf.Network.Port, "Timeout", conf.Network.Timeout, "MaxConn", conf.Network.MaxConn)

	return srv
}
