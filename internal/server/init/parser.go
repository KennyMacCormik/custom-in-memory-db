package init

import (
	"custom-in-memory-db/internal/server/db/parser"
	"log/slog"
)

func Parser(lg *slog.Logger) parser.Parser {
	pr := parser.Parse{}
	pr.New()
	lg.Info("parser init done")
	return &pr
}
