package parser

import (
	"strings"
	"testing"
)

type testCase struct {
	Query    string
	Expected string
	Error    error
}

func TestBuffParser_Read(t *testing.T) {
	var buf BuffParser
	var testCases = []testCase{
		{
			"GET 2\n",
			"GET 2\n",
			nil,
		},
	}

	for _, val := range testCases {
		buf.New(strings.NewReader(val.Query))
		res, err := buf.Read()
		// error expected and present
		if val.Error != nil && err != nil {
			if err.Error() != val.Error.Error() {
				t.Errorf("case %v: expected error: %v, got: %v", val, val.Error, err)
				continue
			}
		}
		// error expected and NOT present
		if val.Error != nil && err == nil {
			t.Errorf("case %v: expected error: %v, got value instead: %v", val, val.Error, res)
			continue
		}
		// error NOT expected and present
		if val.Error == nil && err != nil {
			t.Fatal(err)
		}
		// error NOT expected and NOT present
		if val.Error == nil && err == nil {
			t.Fatal(err)
		}
	}
}
