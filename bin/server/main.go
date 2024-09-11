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

	// Start CLI loop
	for {
		fmt.Println(bp.Read())
	}
}
