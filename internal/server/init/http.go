package init

import (
	"custom-in-memory-db/internal/server/network"
	http2 "custom-in-memory-db/internal/server/network/http"
	"log/slog"
)

func HttpServer(lg *slog.Logger) network.Endpoint {
	http := http2.Server{}
	http.New(lg)
	return &http
}
