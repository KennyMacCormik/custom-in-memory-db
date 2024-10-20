package compute

import (
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/mocks/storage"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"testing"
)

const getKey = "1"
const getKeyNegative = "2"
const getValue = "2"
const nilResult = ""
const setKey = "1"
const setValue = "2"
const success = "OK\n"
const delKey = "1"

var nilLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestComp_New(t *testing.T) {
	st := storage.NewMockStorage(t)
	comp := New(st)
	assert.NotNil(t, comp)
}

func TestComp_Close_Positive_NotCloser(t *testing.T) {
	st := storage.NewMockStorage(t)
	comp := New(st)
	err := comp.Close()
	assert.NoError(t, err)
}

func TestComp_GetPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "GET",
			Arg1:    getKey,
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Get(getKey).Return(getValue, nil)

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.Equal(t, getValue, result)
	assert.Nil(t, err)
}

func TestComp_GetNegative(t *testing.T) {
	testCase := struct {
		input parser.Command
		err   string
	}{
		input: parser.Command{
			Command: "GET",
			Arg1:    getKeyNegative,
		},
		err: "error getting value: key 2 not found",
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Get(getKeyNegative).Return("", errors.New(fmt.Sprintf("key %s not found", getKeyNegative)))

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilResult, result)
}

func TestComp_SetPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "SET",
			Arg1:    setKey,
			Arg2:    setValue,
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Set(setKey, setValue).Return(nil)

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.Equal(t, success, result)
	assert.Nil(t, err)
}

func TestComp_DelPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "DEL",
			Arg1:    delKey,
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Del(delKey).Return(nil)

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.Equal(t, success, result)
	assert.Nil(t, err)
}

func TestComp_DelNegative(t *testing.T) {
	testCase := struct {
		input parser.Command
		err   string
	}{
		input: parser.Command{
			Command: "DEL",
			Arg1:    delKey,
		},
		err: "error deleting value: key 1 not found",
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Del(delKey).Return(errors.New(fmt.Sprintf("key %s not found", delKey)))

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilResult, result)
}

func TestComp_BogusCommand(t *testing.T) {
	testCase := struct {
		input parser.Command
		err   string
	}{
		input: parser.Command{
			Command: "QWE",
			Arg1:    "",
		},
		err: "unknown command",
	}

	st := storage.NewMockStorage(t)

	comp := New(st)

	result, err := comp.Exec(testCase.input, nilLogger)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilResult, result)
}
