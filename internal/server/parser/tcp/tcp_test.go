package tcp

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"testing"
	"time"
)

type testCase struct {
	Query    string
	Expected string
	Error    error
}

func toString(E interface{}) string {
	return fmt.Sprintf("%v", E)
}

func TestTcpParser_BogusIp(t *testing.T) {
	var tcp TcpParser
	err := tcp.New("1271.0.0.1", "8080", 10*time.Second)
	if err.Error() != "tcp listener init error: listen tcp4: lookup 1271.0.0.1: no such host" {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
}

func TestTcpParser_BogusPort(t *testing.T) {
	var tcp TcpParser
	err := tcp.New("127.0.0.1", "80808080", 10*time.Second)
	if err.Error() != "tcp listener init error: listen tcp4: address 80808080: invalid port" {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
}

func TestTcpParser_BusyPort(t *testing.T) {
	var tcp, tcpFail TcpParser
	err := tcp.New("127.0.0.1", "8080", 10*time.Second)
	if err != nil {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
	defer tcp.Close()
	err = tcpFail.New("127.0.0.1", "8080", 10*time.Second)
	if err.Error() != "tcp listener init error: listen tcp4 127.0.0.1:8080: bind: address already in use" {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
}

func TestTcpParser_ClosedListener(t *testing.T) {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))

	var tcp TcpParser
	err := tcp.New("127.0.0.1", "8080", 10*time.Second)
	if err != nil {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
	tcp.Close()

	_, err = tcp.Read([]string{"GET", "SET", "DEL", "QUIT", "EXIT"}, lg)
	if err.Error() != "accept tcp4 127.0.0.1:8080: use of closed network connection" {
		t.Fatalf("Test failed. Unexpected error from closed listener: %v", err)
	}
}

func TestTcpParser_GeneralCases(t *testing.T) {
	var testCases = []testCase{
		// basic commands
		{
			"GET 2",
			"{GET [2]}",
			nil,
		},
		{
			"SET 2 1",
			"{SET [2 1]}",
			nil,
		},
		{
			"DEL 2",
			"{DEL [2]}",
			nil,
		},
		{
			"QUIT",
			"{QUIT []}",
			nil,
		},
		{
			"EXIT",
			"{EXIT []}",
			nil,
		},
		// invalid amount of args
		{
			"GET",
			"",
			errors.New("parsing error: argument validation error: expected 1 argument, got 0"),
		},
		{
			"GET 2 2",
			"",
			errors.New("parsing error: argument validation error: expected 1 argument, got 2"),
		},
		{
			"SET 2",
			"",
			errors.New("parsing error: argument validation error: expected 2 arguments, got 1"),
		},
		{
			"SET 2 2 2",
			"",
			errors.New("parsing error: argument validation error: expected 2 arguments, got 3"),
		},
		{
			"DEL 2 2 2",
			"",
			errors.New("parsing error: argument validation error: expected 1 argument, got 3"),
		},
		{
			"QUIT 1",
			"",
			errors.New("parsing error: argument validation error: expected 0 arguments, got 1"),
		},
		{
			"EXIT 2 2",
			"",
			errors.New("parsing error: argument validation error: expected 0 arguments, got 2"),
		},
		// invalid commands
		{
			"DeL 2",
			"",
			errors.New("parsing error: argument validation error: invalid command: DeL"),
		},
		{
			"qwerty 2",
			"",
			errors.New("parsing error: argument validation error: invalid command: qwerty"),
		},
		// trimming and tabs
		{
			"  DEL     2     ",
			"{DEL [2]}",
			nil,
		},
		{
			"  SET\t\t2          1   ",
			"{SET [2 1]}",
			nil,
		},
		// syntax positive
		{
			"GET 21",
			"{GET [21]}",
			nil,
		},
		{
			"GET 2a1",
			"{GET [2a1]}",
			nil,
		},
		{
			"GET aBc",
			"{GET [aBc]}",
			nil,
		},
		{
			"GET a*Bc**",
			"{GET [a*Bc**]}",
			nil,
		},
		{
			"GET a__B_c",
			"{GET [a__B_c]}",
			nil,
		},
		{
			"GET a/B/c",
			"{GET [a/B/c]}",
			nil,
		},
		// syntax negative
		{
			"GET w$ord",
			"",
			errors.New("parsing error: argument validation error: invalid argument 1: expected printascii,containsany=*_/|alphanum|numeric|alpha"),
		},
		{
			"SET test w\\ord",
			"",
			errors.New("parsing error: argument validation error: invalid argument 2: expected printascii,containsany=*_/|alphanum|numeric|alpha"),
		},
		{
			"\t\tDEL \t w@rd   ",
			"",
			errors.New("parsing error: argument validation error: invalid argument 1: expected printascii,containsany=*_/|alphanum|numeric|alpha"),
		},
		{
			"\t\tDEL \t w-rd   ",
			"",
			errors.New("parsing error: argument validation error: invalid argument 1: expected printascii,containsany=*_/|alphanum|numeric|alpha"),
		},
	}
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))

	var tcp TcpParser
	err := tcp.New("127.0.0.1", "8080", 10*time.Second)
	if err != nil {
		t.Fatalf("Test failed. Cannot create tcp server: %v", err)
	}
	defer tcp.Close()

	for _, val := range testCases {
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			t.Fatalf("Test failed. Cannot dial tcp server: %v", err)
		}
		quit := make(chan bool)
		go func() {
			for {
				select {
				case <-quit:
					return
				default:
					conn.Write([]byte(val.Query + "\n"))
				}
			}
		}()

		res, err := tcp.Read([]string{"GET", "SET", "DEL", "QUIT", "EXIT"}, lg)
		quit <- true
		// error expected and present
		if val.Error != nil && err != nil {
			if err.Error() != val.Error.Error() {
				t.Errorf("case %v: expected error: %v, got: %v", toString(val), val.Error, err)
				continue
			}
		}
		// error expected and NOT present
		if val.Error != nil && err == nil {
			t.Errorf("case %v: expected error: %v, got no error", toString(val), val.Error)
			continue
		}
		// error NOT expected and present
		if val.Error == nil && err != nil {
			t.Errorf("case %v: expected value: %v, got error: %v", toString(val), val.Expected, err)
			continue
		}
		// error NOT expected and NOT present
		if val.Error == nil && err == nil {
			if toString(res) != val.Expected {
				t.Errorf("case %v: expected value: %v, got: %v", toString(val), val.Expected, res)
				continue
			}
		}
	}
}
