package main

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/compute"
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/parser/stdin"
	"custom-in-memory-db/internal/server/storage"
	_map "custom-in-memory-db/internal/server/storage/map"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	// Init logger
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))
	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		lg.Error("config init error", "error", err.Error())
		panic(err)
	}
	// Init compute layer
	bp := stdin.BuffParser{}
	bp.New(os.Stdin)
	// Init parser
	comp := compute.Comp{}
	comp.New()
	// Init storage
	switch conf.Eng.Storage {
	case "map":
		st := _map.MapStorage{}
		st.New()
		// Run app
		run(&bp, &st, &comp, lg)
	default:
		lg.Error("storage init error", "error", fmt.Sprintf("unknown storage: %s", conf.Eng.Storage))
		panic(err)
	}
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
