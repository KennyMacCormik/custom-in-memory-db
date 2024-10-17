package db

import (
	"bytes"
	"custom-in-memory-db/internal/server/db/parser"
	"custom-in-memory-db/mocks/compute"
	"custom-in-memory-db/mocks/network"
	mockParser "custom-in-memory-db/mocks/parser"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"testing"
)

var nilLogger = slog.New(slog.NewTextHandler(io.Discard, nil))

func TestDatabase_New(t *testing.T) {
	comp := compute.NewMockCompute(t)
	netEndpoint := network.NewMockEndpoint(t)
	pr := mockParser.NewMockParser(t)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)
}

func TestDatabase_HandleRequest_Positive(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		res string
	}{
		cmd: parser.Command{
			Command: "GET",
			Arg1:    "1",
		},
		in:  "GET 1\n",
		res: "2",
	}
	r := bytes.NewBuffer([]byte(testCase.in))

	comp := compute.NewMockCompute(t)
	comp.EXPECT().Exec(testCase.cmd, nilLogger).Return(testCase.res, nil)

	netEndpoint := network.NewMockEndpoint(t)

	pr := mockParser.NewMockParser(t)
	pr.EXPECT().Read(r, nilLogger).Return(testCase.cmd, nil)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	result, err := db.HandleRequest(r, nilLogger)
	assert.NoError(t, err)
	assert.Equal(t, testCase.res, result)
}

func TestDatabase_HandleRequest_Negative_Parser(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		err string
	}{
		cmd: parser.Command{
			Command: "GET",
			Arg1:    "1",
		},
		in:  "GET 1\n",
		err: "test error",
	}
	r := bytes.NewBuffer([]byte(testCase.in))

	comp := compute.NewMockCompute(t)
	netEndpoint := network.NewMockEndpoint(t)

	pr := mockParser.NewMockParser(t)
	pr.EXPECT().Read(r, nilLogger).Return(parser.Command{}, errors.New("test error"))

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	result, err := db.HandleRequest(r, nilLogger)
	assert.Empty(t, result)
	assert.EqualError(t, err, testCase.err)
}

func TestDatabase_HandleRequest_Negative_Compute(t *testing.T) {
	testCase := struct {
		cmd parser.Command
		in  string
		err string
	}{
		cmd: parser.Command{
			Command: "GET",
			Arg1:    "1",
		},
		in:  "GET 1\n",
		err: "test error",
	}
	r := bytes.NewBuffer([]byte(testCase.in))

	comp := compute.NewMockCompute(t)
	comp.EXPECT().Exec(testCase.cmd, nilLogger).Return("", errors.New("test error"))

	netEndpoint := network.NewMockEndpoint(t)

	pr := mockParser.NewMockParser(t)
	pr.EXPECT().Read(r, nilLogger).Return(testCase.cmd, nil)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	result, err := db.HandleRequest(r, nilLogger)
	assert.Empty(t, result)
	assert.EqualError(t, err, testCase.err)
}

func TestDatabase_ListenClient(t *testing.T) {
	comp := compute.NewMockCompute(t)

	netEndpoint := network.NewMockEndpoint(t)
	netEndpoint.EXPECT().Listen(mock.AnythingOfType("network.Handler"))

	pr := mockParser.NewMockParser(t)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	db.ListenClient()
}

func TestDatabase_Close_Positive(t *testing.T) {
	comp := compute.NewMockCompute(t)
	comp.EXPECT().Close().Return(nil)

	netEndpoint := network.NewMockEndpoint(t)

	pr := mockParser.NewMockParser(t)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	err := db.Close()
	assert.NoError(t, err)
}

func TestDatabase_Close_Negative(t *testing.T) {
	comp := compute.NewMockCompute(t)
	comp.EXPECT().Close().Return(fmt.Errorf("error: %w", errors.New("test error")))

	netEndpoint := network.NewMockEndpoint(t)

	pr := mockParser.NewMockParser(t)

	db := Database{}
	db.New(comp, netEndpoint, pr, nilLogger)

	err := db.Close()
	assert.EqualError(t, err, "Database.Close() failed")
}
