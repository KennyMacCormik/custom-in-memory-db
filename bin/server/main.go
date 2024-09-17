package main

import (
	"custom-in-memory-db/internal/server/compute"
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/storage"
	_map "custom-in-memory-db/internal/server/storage/map"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// Init
	bp := parser.BuffParser{}
	bp.New(os.Stdin)

	comp := compute.Comp{}
	comp.New()

	st := _map.MapStorage{}
	st.New()

	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))

	// Run app
	run(&bp, &st, &comp, lg)
}

func run(p parser.Parser, st storage.Storage, comp compute.Compute, lg *slog.Logger) {
	for {
		str := comp.HandleRequest(p, st, lg)
		if str == compute.Stop {
			break
		}
		if str != compute.Ok &&
			str != compute.ErrParse &&
			str != compute.ErrGet &&
			str != compute.ErrSet &&
			str != compute.ErrDel &&
			str != compute.ErrUnexpectedCommand {
			fmt.Println(str)
		}
	}
}
