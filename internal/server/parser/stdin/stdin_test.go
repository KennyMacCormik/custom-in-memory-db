package stdin

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	Query    string
	Expected string
	Error    error
}

func toString(E interface{}) string {
	return fmt.Sprintf("%v", E)
}

func TestBuffParser(t *testing.T) {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))

	var buf BuffParser
	var out []byte
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
			errors.New("parsing error: argument validation error: invalid command: QUIT"),
		},
		{
			"EXIT",
			"{EXIT []}",
			errors.New("parsing error: argument validation error: invalid command: EXIT"),
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
			errors.New("parsing error: argument validation error: invalid command: QUIT"),
		},
		{
			"EXIT 2 2",
			"",
			errors.New("parsing error: argument validation error: invalid command: EXIT"),
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

	for _, val := range testCases {
		//t.Logf("case %d: %v", i, val)
		buf.New(strings.NewReader(val.Query), bytes.NewBuffer(out))
		res, wc, err := buf.Read([]string{"GET", "SET", "DEL"}, lg)
		// error expected and present
		if val.Error != nil && err != nil {
			if err.Error() != val.Error.Error() {
				t.Errorf("case %v: expected error: %v, got: %v", toString(val), val.Error, err)
				continue
			}
		}
		// error expected and NOT present
		if val.Error != nil && err == nil {
			wc.Close()
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
			wc.Close()
			if toString(res) != val.Expected {
				t.Errorf("case %v: expected value: %v, got: %v", toString(val), val.Expected, res)
				continue
			}
		}
	}
}
