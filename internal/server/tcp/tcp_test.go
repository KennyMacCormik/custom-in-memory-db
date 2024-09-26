package tcp

import (
	"bytes"
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/mocks/compute"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net"
	"testing"
	"time"
)

const ip = "0.0.0.0"
const port = "8080"
const timeout = 1
const goMax = 100

func TestServer_NewAndClose(t *testing.T) {
	srv := Server{}
	err := srv.New(ip, port, timeout*time.Second, goMax, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)

	err = srv.Close()
	assert.NoError(t, err)
}

func TestConnMeter_New(t *testing.T) {
	var cm connMeter
	cm.New(goMax)

	assert.Equal(t, goMax, cm.maxConn)
	assert.NotEqual(t, nil, cm.cond)
}

// TODO fix pipes
func TestHandleClient(t *testing.T) {
	testCase := struct {
		input     string
		compInput parser.Command
		result    string
	}{
		input: "GET 1\n",
		compInput: parser.Command{
			Command: "GET",
			Args:    []string{"1"},
		},
		result: "2",
	}

	var cm connMeter
	cm.New(goMax)

	comp := compute.NewMockCompute(t)
	comp.EXPECT().Exec(testCase.compInput).Return(testCase.result, nil)

	var resp []byte
	logger := &bytes.Buffer{}
	lg := slog.New(slog.NewTextHandler(logger, nil))

	r, w := net.Pipe()
	go func() {
		w.Write([]byte(testCase.input))
		for {
			// why not working?
			w.Read(resp)
		}
	}()

	handleClient(r, timeout*time.Second, comp, &cm, lg)

	assert.Equal(t, "qwe", string(logger.Bytes()))
	assert.Equal(t, "qwe", string(resp))

	r.Close()
	w.Close()
}
