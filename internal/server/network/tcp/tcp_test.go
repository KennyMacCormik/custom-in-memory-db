package tcp

import (
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
	"time"
)

const ip = "0.0.0.0"
const port = "8080"
const timeout = 1
const goMax = 100

var nilLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestServer_NewAndClose(t *testing.T) {
	srv, err := New(ip, port, timeout*time.Second, goMax, nilLogger)
	assert.NoError(t, err)
	assert.NotNil(t, srv)

	err = srv.Close()
	assert.NoError(t, err)
}

func TestConnMeter_New(t *testing.T) {
	var cm connMeter
	cm.new(goMax)

	assert.Equal(t, goMax, cm.maxConn)
	assert.NotEqual(t, nil, cm.cond)
}
