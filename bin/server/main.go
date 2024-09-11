package main

import (
	"custom-in-memory-db/internal/server/compute"
	"custom-in-memory-db/internal/server/parser"
	_map "custom-in-memory-db/internal/server/storage/map"
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

	// Read cmd args

	// Run app
	// dunno why I need to use &bp but comp works fine
	comp.ControlLoop(&bp, &st)
}
