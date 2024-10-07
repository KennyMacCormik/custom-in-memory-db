package db

import (
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/mocks/compute"
	ioMock "custom-in-memory-db/mocks/io"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"testing"
)

const nilValue = ""

func TestDatabase_New(t *testing.T) {
	comp := compute.NewMockCompute(t)

	db := Database{}
	db.New(comp)

	assert.NotNil(t, db.comp)
}

func TestDatabase_HandleRequest_Positive(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		res string
	}{
		cmd: parser.Command{
			Command: "GET",
			Args:    []string{"1"},
		},
		in:  "GET 1\n",
		res: "2",
	}

	comp := compute.NewMockCompute(t)
	comp.EXPECT().Exec(testCase.cmd, slog.New(slog.NewTextHandler(io.Discard, nil))).Return(testCase.res, nil)

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.in)
	}).Return(len(testCase.in), nil)

	db := Database{}
	db.New(comp)

	result, err := db.HandleRequest(r, slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert.NoError(t, err)
	assert.Equal(t, testCase.res, result)
}

func TestDatabase_HandleRequest_NegativeReader(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		res string
		err string
	}{
		cmd: parser.Command{
			Command: "GET",
			Args:    []string{"1"},
		},
		in:  "GET 1\n",
		res: "2",
		err: "mock",
	}

	comp := compute.NewMockCompute(t)

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {}).Return(0, errors.New("mock"))

	db := Database{}
	db.New(comp)

	result, err := db.HandleRequest(r, slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert.Equal(t, nilValue, result)
	assert.EqualError(t, err, testCase.err)
}

func TestDatabase_HandleRequest_NegativeCompute(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		res string
		err string
	}{
		cmd: parser.Command{
			Command: "GET",
			Args:    []string{"1"},
		},
		in:  "GET 1\n",
		res: "2",
		err: "mock",
	}

	comp := compute.NewMockCompute(t)
	comp.EXPECT().Exec(testCase.cmd, slog.New(slog.NewTextHandler(io.Discard, nil))).Return("", errors.New("mock"))

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.in)
	}).Return(len(testCase.in), nil)

	db := Database{}
	db.New(comp)

	result, err := db.HandleRequest(r, slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, nilValue, result)
}
