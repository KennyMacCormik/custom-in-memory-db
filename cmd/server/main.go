package main

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/internal/server/compute"
	"custom-in-memory-db/internal/server/parser"
	"custom-in-memory-db/internal/server/parser/stdin"
	"custom-in-memory-db/internal/server/parser/tcp"
	"custom-in-memory-db/internal/server/storage"
	_map "custom-in-memory-db/internal/server/storage/map"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func main() {
	// Init logger
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	// Init config
	conf := cmd.Config{}
	err := conf.New()
	if err != nil {
		lg.Error("config init error", "error", err.Error())
		panic(err)
	}

	// Init compute layer
	comp := compute.Comp{}
	comp.New()

	run(initParser(conf, lg),
		initStorage(conf, lg),
		&comp, lg)
}

func run(p parser.Parser, st storage.Storage, comp compute.Compute, lg *slog.Logger) {
	defer p.Close()
	for {
		_ = comp.HandleRequest(p, st, lg)
	}
}

func initStorage(conf cmd.Config, lg *slog.Logger) storage.Storage {
	switch conf.Eng.Storage {
	case "map":
		st := _map.MapStorage{}
		st.New()
		return &st
	default:
		str := fmt.Sprintf("unknown storage: %s", conf.Eng.Storage)
		lg.Error("storage init error", "error", str)
		panic(str)
	}
}

func initParser(conf cmd.Config, lg *slog.Logger) parser.Parser {
	switch conf.Eng.Input {
	case "tcp4":
		pars := tcp.TcpParser{}
		err := pars.New(conf.Net.Address, strconv.Itoa(conf.Net.Port), conf.Net.IdleTimeout)
		if err != nil {
			lg.Error("parser init error", "error", fmt.Sprintf("cannot init tcp parser: %s", err.Error()))
			panic(err)
		}
		return &pars
	case "stdin":
		pars := stdin.BuffParser{}
		pars.New(os.Stdin, os.Stdout)
		return &pars
	default:
		str := fmt.Sprintf("unknown parser: %s", conf.Eng.Input)
		lg.Error("parser init error", "error", str)
		panic(str)
	}
}
