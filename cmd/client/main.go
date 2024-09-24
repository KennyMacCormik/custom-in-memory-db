package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

type CmdInput struct {
	Address string
	Port    int
	Cmd     string
}

func main() {
	// read flags
	cmnd := ReadFlags()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println(fmt.Errorf("cannot dial tcp server: %w", err).Error())
		os.Exit(1)
	}
	/*
		err = conn.SetDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println(fmt.Errorf("cannot set deadline: %w", err))
			os.Exit(1)
		}
	*/

	_, err = conn.Write([]byte(cmnd.Cmd))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot send command: %w", err).Error())
		os.Exit(1)
	}

	/*
		resp := make([]byte, 1024)
		_, err = conn.Read(resp)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot read response: %w", err).Error())
			os.Exit(1)
		}

	*/
}

func ReadFlags() CmdInput {
	result := CmdInput{}

	bindFlags(&result)
	flag.Parse()
	result.Cmd = strings.Join(flag.Args(), " ")

	return result
}

func bindFlags(c *CmdInput) {
	flag.StringVar(&c.Address, "address", "127.0.0.1", "Address to bind to")
	flag.StringVar(&c.Address, "a", "127.0.0.1", "Address to bind to")

	flag.IntVar(&c.Port, "port", 8080, "Port to bind to")
	flag.IntVar(&c.Port, "p", 8080, "Port to bind to")
}
