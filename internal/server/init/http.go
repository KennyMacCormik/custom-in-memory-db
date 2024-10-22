package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/network"
	http2 "custom-in-memory-db/internal/server/network/http"
	"log/slog"
)

func HttpServer(conf cmd.Config, lg *slog.Logger) network.Endpoint {
	http := http2.Server{}
	http.New(conf, lg)
	return &http
}
