package wal

import (
	"bufio"
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/internal/server/db/seg"
	"fmt"
	"log/slog"
	"os"
)

// Recover loads wal to the running storage.Storage
func Recover(seg *seg.Segments, setFunc func(k, v string) error, delFunc func(k string) error, lg *slog.Logger) error {
	for _, file := range seg.ExportSegNames() {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to load %q: %w", file, err)
		}

		loadFile(f, setFunc, delFunc, lg)
	}

	return nil
}

// loadFile reads commands from a provided file and commits them to the running storage.Storage
func loadFile(f *os.File, set func(k, v string) error, del func(k string) error, lg *slog.Logger) {
	pr := parser.New()
	reader := bufio.NewReader(f)
	var err error = nil
	for ; err == nil; _, err = reader.Peek(1) {
		c, e := pr.Read(reader, lg)
		if e != nil {
			// skip incorrect line in a file
			continue
		}
		if c.Command == "SET" {
			_ = set(c.Arg1, c.Arg2)
			continue
		}
		_ = del(c.Arg1)
	}
}
