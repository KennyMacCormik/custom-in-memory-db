package main

import (
	"custom-in-memory-db/internal/server/compute/parser"
	"fmt"
	"os"
)

func main() {

	// Init
	bp := parser.BuffParser{}
	bp.New(os.Stdin)

	// Read cmd args

	// Run app
	run(&bp)
}

func run(p parser.Parser) {
	for {
		fmt.Println(p.Read())
	}
}
