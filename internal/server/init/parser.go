package init

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/db/parser"
	"log/slog"
)

func Parser(conf cmd.Config, lg *slog.Logger) parser.Parser {
	pr := parser.Parse{}
	pr.New(conf)
	lg.Info("parser init done")
	return &pr
}
