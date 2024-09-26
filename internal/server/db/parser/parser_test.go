package parser

import (
	"custom-in-memory-db/mocks/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"testing"
)

func TestRead_GetPositive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "GET 1\n",
		expected: Command{
			Command: "GET",
			Args:    []string{"1"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_GetNegative_ZeroArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "GET\n",
		expected: Command{},
		err:      "argument validation error: expected 1 argument, got 0",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_GetNegative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "GET 1 2\n",
		expected: Command{},
		err:      "argument validation error: expected 1 argument, got 2",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_SetPositive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "SET 1 2\n",
		expected: Command{
			Command: "SET",
			Args:    []string{"1", "2"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_SetNegative_NoArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET\n",
		expected: Command{},
		err:      "argument validation error: expected 2 arguments, got 0",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_SetNegative_InsufficientArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET 1\n",
		expected: Command{},
		err:      "argument validation error: expected 2 arguments, got 1",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_SetNegative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "SET 1 2 3\n",
		expected: Command{},
		err:      "argument validation error: expected 2 arguments, got 3",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_DelPositive(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL 1\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"1"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_DelNegative_NoArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL\n",
		expected: Command{},
		err:      "argument validation error: expected 1 argument, got 0",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_DelNegative_ExcessiveArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL 1 2\n",
		expected: Command{},
		err:      "argument validation error: expected 1 argument, got 2",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_BogusCommand_WithoutArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "QWERTY\n",
		expected: Command{},
		err:      "argument validation error: invalid command: QWERTY",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_BogusCommand_WithArgs(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "HELLO 1 2\n",
		expected: Command{},
		err:      "argument validation error: invalid command: HELLO",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_OneSymbolNewLine_String(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "\n",
		expected: Command{},
		err:      "empty command",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_ZeroString(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "",
		expected: Command{},
		err:      "multiple Read calls return no data or error",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

// printascii,containsany=*_/|alphanum|numeric|alpha

func TestRead_CorrectArgs_alphanum(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL q1\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"q1"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_CorrectArgs_numeric(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL 1\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"1"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_CorrectArgs_alpha(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL qwe\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"qwe"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_CorrectArgs_underscore(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL test_a\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"test_a"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_CorrectArgs_slash(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL test/a\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"test/a"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_CorrectArgs_star(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
	}{
		ioInput: "DEL test*a\n",
		expected: Command{
			Command: "DEL",
			Args:    []string{"test*a"},
		},
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.NoError(t, err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_IncorrectArgs_Unicode(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL val_\u2318\n",
		expected: Command{},
		err:      "argument validation error: invalid argument 1: expected printascii,containsany=*_/|alphanum|numeric|alpha",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}

func TestRead_IncorrectArgs_UnexpectedSymbol(t *testing.T) {
	testCase := struct {
		ioInput  string
		expected Command
		err      string
	}{
		ioInput:  "DEL val\\1\n",
		expected: Command{},
		err:      "argument validation error: invalid argument 1: expected printascii,containsany=*_/|alphanum|numeric|alpha",
	}

	r := ioMock.NewMockReader(t)
	r.On("Read", mock.Anything).Run(func(args mock.Arguments) {
		bytes := args[0].([]byte)
		copy(bytes, testCase.ioInput)
	}).Return(len(testCase.ioInput), nil)

	val, err := Read(r, slog.New(slog.NewTextHandler(io.Discard, nil)))

	assert.EqualError(t, err, testCase.err)
	assert.Equal(t, testCase.expected, val)
}
