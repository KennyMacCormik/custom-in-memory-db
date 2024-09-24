package stdin

import (
	"bufio"
	"custom-in-memory-db/internal/server/parser"
	"io"
	"log/slog"
)

type BuffWriteCloser struct {
	writer *bufio.Writer
}

func (b *BuffWriteCloser) New(w io.Writer) {
	b.writer = bufio.NewWriter(w)
}

func (b *BuffWriteCloser) Write(p []byte) (n int, err error) {
	return b.writer.Write(p)
}

func (b *BuffWriteCloser) Close() error {
	return b.writer.Flush()
}

type BuffParser struct {
	reader io.Reader
	writer io.Writer
}

// New creates new buffer to read input from
func (bp *BuffParser) New(in io.Reader, out io.Writer) {
	bp.reader = in
	bp.writer = out
}

func (bp *BuffParser) Write(response string, wc io.WriteCloser, lg *slog.Logger) error {
	return nil
}

func (bp *BuffParser) Close() error {
	return nil
}

// Read reads buffer input and tries to compose it into valid Command struct
func (bp *BuffParser) Read(vc []string, lg *slog.Logger) (parser.Command, io.WriteCloser, error) {
	r := bufio.NewReader(bp.reader)
	result, err := parser.BufferRead(r, vc, lg)
	if err != nil {
		return parser.Command{}, nil, err
	}

	var wc BuffWriteCloser
	wc.New(bp.writer)

	return result, &wc, nil
}
