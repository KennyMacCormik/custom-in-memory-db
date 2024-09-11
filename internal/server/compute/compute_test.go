package compute

import (
	"custom-in-memory-db/internal/server/parser"
	_map "custom-in-memory-db/internal/server/storage/map"
	"log/slog"
	"os"
	"strings"
	"testing"
)

type testCase struct {
	Input    string
	Expected string
}

func TestComp_HandleRequest(t *testing.T) {
	comp := Comp{}
	comp.New()

	st := _map.MapStorage{}
	st.New()

	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	lg := slog.New(slog.NewTextHandler(os.Stdin, &slog.HandlerOptions{Level: logLevel}))

	var testCases = []testCase{
		{
			"GET 1",
			"ErrGet",
		},
		{
			"SET 1 2",
			"ok",
		},
		{
			"GET 1",
			"2",
		},
		{
			"DEL 1",
			"ok",
		},
		{
			"DEL 1",
			"ErrDel",
		},
		{
			"QQUIT",
			"ErrParce",
		},
		{
			"QUIT",
			"Stop",
		},
	}

	for _, val := range testCases {
		p := parser.BuffParser{}
		p.New(strings.NewReader(val.Input))

		res := comp.HandleRequest(&p, &st, lg)
		if res != val.Expected {
			t.Errorf("input %v, expected %s, got %s", val, val.Expected, res)
		}
	}

}
