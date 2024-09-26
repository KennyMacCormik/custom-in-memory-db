package compute

import (
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/mocks/storage"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const getKey = "1"
const getKeyNegative = "2"
const getValue = "2"
const nilResult = ""
const setKey = "1"
const setValue = "2"
const success = "OK"
const delKey = "1"

func TestComp_New(t *testing.T) {
	st := storage.NewMockStorage(t)
	comp := Comp{}
	comp.New(st)
}

func TestComp_GetPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "GET",
			Args:    []string{getKey},
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Get(getKey).Return(getValue, nil)

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

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
			Args:    []string{getKeyNegative},
		},
		err: "error getting value: key 2 not found",
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Get(getKeyNegative).Return("", errors.New(fmt.Sprintf("key %s not found", getKeyNegative)))

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilResult, result)
}

func TestComp_SetPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "SET",
			Args:    []string{setKey, setValue},
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Set(setKey, setValue).Return(nil)

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

	assert.Equal(t, success, result)
	assert.Nil(t, err)
}

func TestComp_DelPositive(t *testing.T) {
	testCase := struct {
		input parser.Command
	}{
		input: parser.Command{
			Command: "DEL",
			Args:    []string{delKey},
		},
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Del(delKey).Return(nil)

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

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
			Args:    []string{delKey},
		},
		err: "error deleting value: key 1 not found",
	}

	st := storage.NewMockStorage(t)
	st.EXPECT().Del(delKey).Return(errors.New(fmt.Sprintf("key %s not found", delKey)))

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

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
			Args:    []string{},
		},
		err: "unknown command",
	}

	st := storage.NewMockStorage(t)

	comp := Comp{}
	comp.New(st)

	result, err := comp.Exec(testCase.input)

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilResult, result)
}
