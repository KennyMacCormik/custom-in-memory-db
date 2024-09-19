package stdin

import (
	"bufio"
	"custom-in-memory-db/internal/server/parser"
	"io"
	"log/slog"
)

type BuffParser struct {
	reader *bufio.Reader
}

// New creates new buffer to read input from
func (bp *BuffParser) New(in io.Reader) {
	bp.reader = bufio.NewReader(in)
}

// Read reads buffer input and tries to compose it into valid Command struct
func (bp *BuffParser) Read(vc []string, lg *slog.Logger) (parser.Command, error) {
	result, err := parser.BufferRead(bp.reader, vc, lg)
	if err != nil {
		return parser.Command{}, err
	}

	return result, nil
}
