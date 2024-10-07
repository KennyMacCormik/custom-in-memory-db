package init

import (
	"custom-in-memory-db/internal/server/cmd"
	_map "custom-in-memory-db/internal/server/db/storage/map"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

func TestTcpServer_Positive(t *testing.T) {
	var engine = cmd.Engine{
		APP_STORAGE: "mem",
	}
	var rep = cmd.Replication{
		REP_TYPE: "master",
		// Feel free to choose any port for this test
		REP_PORT: 8082,
	}
	var network = cmd.Network{
		NET_ADDR: "127.0.0.1",
		// Feel free to choose any port for this test
		NET_PORT: 8080,
	}
	var conf = cmd.Config{
		Engine: engine,
		Rep:    rep,
		Net:    network,
	}
	nilLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	srv, err := TcpServer(conf, nilLogger)
	assert.NoError(t, err)
	_ = srv.Close()
}

func TestTcpServer_Negative(t *testing.T) {
	var engine = cmd.Engine{
		APP_STORAGE: "mem",
	}
	var rep = cmd.Replication{
		REP_TYPE: "master",
		// Feel free to choose any port for this test
		REP_PORT: 8082,
	}
	var network = cmd.Network{
		NET_ADDR: "1271.0.0.1",
		// Feel free to choose any port for this test
		NET_PORT: 8080,
	}
	var conf = cmd.Config{
		Engine: engine,
		Rep:    rep,
		Net:    network,
	}
	var errStr = "init.TcpServer() failed: tcp listener init error: listen tcp4: lookup 1271.0.0.1: no such host"
	nilLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	srv, err := TcpServer(conf, nilLogger)
	assert.Nil(t, srv)
	assert.EqualError(t, err, errStr)
}

func TestLogger_Text(t *testing.T) {
	var engine = cmd.Engine{
		APP_STORAGE: "mem",
	}
	var rep = cmd.Replication{
		REP_TYPE: "master",
		// Feel free to choose any port for this test
		REP_PORT: 8082,
	}
	var network = cmd.Network{
		NET_ADDR: "127.0.0.1",
		// Feel free to choose any port for this test
		NET_PORT: 8080,
	}
	var logger = cmd.Logging{
		LOG_FORMAT: "text",
	}
	var conf = cmd.Config{
		Engine: engine,
		Rep:    rep,
		Net:    network,
		Log:    logger,
	}

	lg := Logger(conf)
	assert.NotNil(t, lg)
	assert.IsType(t, &slog.Logger{}, lg)

	logger = cmd.Logging{
		LOG_FORMAT: "json",
	}
	conf = cmd.Config{
		Engine: engine,
		Rep:    rep,
		Net:    network,
		Log:    logger,
	}

	lg = Logger(conf)
	assert.NotNil(t, lg)
	assert.IsType(t, &slog.Logger{}, lg)
}

func TestInitMapStorage(t *testing.T) {
	nilLogger := slog.New(slog.NewTextHandler(io.Discard, nil))
	st := initMapStorage(nilLogger)

	assert.NotNil(t, st)
	assert.IsType(t, &_map.MapStorage{}, st)
}
