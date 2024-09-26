package main

import (
	"flag"
	"strings"
)

func ReadFlags() CmdInput {
	result := CmdInput{}

	bindFlags(&result)
	flag.Parse()
	result.Cmd = strings.Join(flag.Args(), " ")
	result.Cmd += "\n"

	return result
}

func bindFlags(c *CmdInput) {
	flag.StringVar(&c.Address, "address", "127.0.0.1", "Address to bind to")
	flag.StringVar(&c.Address, "a", "127.0.0.1", "Address to bind to")

	flag.IntVar(&c.Port, "port", 8080, "Port to bind to")
	flag.IntVar(&c.Port, "p", 8080, "Port to bind to")
}
