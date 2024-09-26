package main

import (
	"fmt"
	"io"
	net2 "net"
	"time"
)

func invokeTCP(net string, addr string, port string, cmd string) ([]byte, error) {
	conn, err := net2.Dial(net, addr+":"+port)
	if err != nil {
		fmt.Println(fmt.Errorf("cannot dial tcp server: %w", err).Error())
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(10000 * time.Second))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot set deadline: %w", err))
		return nil, err
	}

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		fmt.Println(fmt.Errorf("cannot send command: %w", err).Error())
		return nil, err
	}

	var n int
	resp := make([]byte, 4096)
	for n == 0 {
		n, err = conn.Read(resp)
		if err != nil && err != io.EOF {
			fmt.Println(fmt.Errorf("cannot read response: %w", err).Error())
			return nil, err
		}
	}

	return resp, nil
}
