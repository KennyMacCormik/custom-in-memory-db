package wal

import (
	"custom-in-memory-db/internal/server/cmd"
	"custom-in-memory-db/mocks/storage"
	walMock "custom-in-memory-db/mocks/wal"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"time"

	"testing"
)

func TestWal_New(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}

	st := storage.NewMockStorage(t)
	writer := walMock.NewMockWriterInterface(t)

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)
}

func TestWal_Close_Positive(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}

	st := storage.NewMockStorage(t)
	st.On("Close").Return(nil)

	writer := walMock.NewMockWriterInterface(t)
	writer.On("WriteAndRotate", mock.Anything).Return(nil)
	writer.On("Close").Return(nil)

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)

	// if I uncomment this sometimes
	// fatal error: all goroutines are asleep - deadlock!
	// happens.
	// Why?
	//wal.Start()

	err := wal.Close()
	assert.NoError(t, err)
}

func TestWal_Close_Negative_All(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}
	var errStr = "wal.Close().storage.Close() and wal.Close().barrier.Close() joined error: barrier.Close().writeWal() and barrier.Close().writer.writeWal() joined error: writer.Close error; barrier.writeWal().writer.WriteAndRotate() failed: writer.WriteAndRotate error; st.Close error"

	st := storage.NewMockStorage(t)
	st.On("Close").Return(errors.New("st.Close error"))

	writer := walMock.NewMockWriterInterface(t)
	writer.On("WriteAndRotate", mock.Anything).Return(errors.New("writer.WriteAndRotate error"))
	writer.On("Close").Return(errors.New("writer.Close error"))

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)

	err := wal.Close()
	assert.EqualError(t, err, errStr)
}

func TestWal_Recover_Positive(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}

	st := storage.NewMockStorage(t)

	writer := walMock.NewMockWriterInterface(t)
	writer.On("Recover", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	nilLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)

	err := wal.Recover(conf, nilLogger)
	assert.NoError(t, err)
}

func TestWal_Recover_Negative(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}
	var errString = "wal.New().barrier.New() failed:: test error"

	st := storage.NewMockStorage(t)

	writer := walMock.NewMockWriterInterface(t)
	writer.On("Recover", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test error"))

	nilLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)

	err := wal.Recover(conf, nilLogger)
	assert.EqualError(t, err, errString)
}

func TestWal_Start(t *testing.T) {
	var walEnv = cmd.Wal{
		WAL_BATCH_SIZE:    10,
		WAL_BATCH_TIMEOUT: 1 * time.Second,
		WAL_SEG_SIZE:      1,
		WAL_SEG_RECOVER:   false,
	}
	var conf = cmd.Config{
		Wal: walEnv,
	}

	st := storage.NewMockStorage(t)
	writer := walMock.NewMockWriterInterface(t)

	wal := Wal{}
	wal.New(conf, st, writer)

	assert.NotNil(t, wal.st)
	assert.NotNil(t, wal.bar)
	assert.NotNil(t, wal.sendToBarrier)
	assert.NotNil(t, wal.w)

	wal.Start()
}
