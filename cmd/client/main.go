package main

import (
	"fmt"
	"os"
	"strconv"
)

const net = "tcp"
const errExit = 1

type CmdInput struct {
	Address string
	Port    int
	Cmd     string
}

func main() {
	// read cmdline input
	cmdInput := ReadFlags()

	result, err := invokeTCP(net, cmdInput.Address, strconv.Itoa(cmdInput.Port), cmdInput.Cmd)
	if err != nil {
		fmt.Println(fmt.Errorf("tcp error: %w", err))
		os.Exit(errExit)
	}

	fmt.Println(string(result))
}
