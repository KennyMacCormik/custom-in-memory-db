package shared

import (
	"custom-in-memory-db/cmd/client/cmd/conf"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const errExit = 1

func InvokeTcp(cfg conf.Config, c string) []byte {
	addr := strings.Join([]string{cfg.Server, strconv.Itoa(cfg.Port)}, ":")
	logVerbose(fmt.Sprintf("DEBUG connecting to: %s\n", addr), cfg.Verbose)
	conn, err := net.Dial("tcp4", addr)
	logError(fmt.Errorf("cannot dial tcp server: %w", err))

	err = conn.SetDeadline(time.Now().Add(cfg.Timeout))
	logError(fmt.Errorf("cannot set deadline: %w", err))
	logVerbose(fmt.Sprintf("DEBUG connection timeout: %s\n", cfg.Timeout), cfg.Verbose)

	logVerbose(fmt.Sprintf("DEBUG sending command: %q\n", c), cfg.Verbose)

	_, err = conn.Write([]byte(c))
	logError(fmt.Errorf("cannot send command: %w", err))
	logVerbose(fmt.Sprintf("DEBUG Done\n"), cfg.Verbose)

	resp := make([]byte, 4096)
	resp, err = waitForResponse(conn, resp)
	logError(fmt.Errorf("read response failed: %w", err))

	return resp
}

func logVerbose(s string, verbose bool) {
	if verbose {
		fmt.Print(s)
	}
}

func logError(err error) {
	if errors.Unwrap(err) != nil {
		fmt.Println(err)
		os.Exit(errExit)
	}
}

func waitForResponse(conn net.Conn, s []byte) ([]byte, error) {
	for {
		n, err := conn.Read(s)
		if err != nil && errors.Is(err, os.ErrDeadlineExceeded) {
			return nil, errors.New("timeout waiting for response")

		}
		if err != nil {
			return nil, fmt.Errorf("unexpected error: %w", err)
		}
		if n > 0 {
			break
		}
	}

	return s, nil
}
