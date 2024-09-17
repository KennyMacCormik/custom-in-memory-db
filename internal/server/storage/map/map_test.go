package _map

import (
	"custom-in-memory-db/internal/server/parser"
	"errors"
	"fmt"
	"testing"
)

// Init exists only for testing purpose
func (s *MapStorage) init() {
	s.m["1"] = "one"
	s.m["2"] = "two"
}

type testCase struct {
	Comm     parser.Command
	Expected string
	Error    error
}

func TestMapStorage(t *testing.T) {
	var storage MapStorage
	storage.New()
	storage.init()

	var testCases = []testCase{
		{
			parser.Command{
				"GET",
				[]string{"1"},
			},
			"one",
			nil,
		},
		{
			parser.Command{
				"GET",
				[]string{"3"},
			},
			"",
			errors.New("key 3 not found"),
		},
		{
			parser.Command{
				"SET",
				[]string{"3", "three"},
			},
			"three",
			nil,
		},
		{
			parser.Command{
				"SET",
				[]string{"1", "one1"},
			},
			"one1",
			nil,
		},
		// error because we check existence of key after deletion
		{
			parser.Command{
				"DEL",
				[]string{"2"},
			},
			"",
			errors.New("key 2 not found"),
		},
		{
			parser.Command{
				"DEL",
				[]string{"5"},
			},
			"",
			errors.New("key 5 not found"),
		},
	}
	for _, val := range testCases {
		switch val.Comm.Command {
		case "GET":
			res, err := storage.Get(val.Comm)
			validateTest(val, res, err, t)
		case "SET":
			err := storage.Set(val.Comm)
			res, err := storage.Get(val.Comm)
			validateTest(val, res, err, t)
		case "DEL":
			err := storage.Del(val.Comm)
			res, err := storage.Get(val.Comm)
			validateTest(val, res, err, t)
		}
	}
}

func validateTest(val testCase, res string, err error, t *testing.T) {
	// error expected and present
	if val.Error != nil && err != nil {
		if err.Error() != val.Error.Error() {
			t.Errorf("case %v: expected error: %v, got: %v", toString(val), val.Error, err)
		}
	}
	// error expected and NOT present
	if val.Error != nil && err == nil {
		t.Errorf("case %v: expected error: %v, got no error", toString(val), val.Error)
	}
	// error NOT expected and present
	if val.Error == nil && err != nil {
		t.Errorf("case %v: expected value: %v, got error: %v", toString(val), val.Expected, err)
	}
	// error NOT expected and NOT present
	if val.Error == nil && err == nil {
		if toString(res) != val.Expected {
			t.Errorf("case %v: expected value: %v, got: %v", toString(val), val.Expected, res)
		}
	}
}

func toString(E interface{}) string {
	return fmt.Sprintf("%v", E)
}
