package _map

import (
	"custom-in-memory-db/internal/server/cmd"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"sync"
	"testing"
)

const firstKey = "1"
const firstVal = "2"
const firstErrKey = "3"
const firstSetVal = "5"

func TestMapStorage_New(t *testing.T) {
	var st MapStorage
	st.New()
}

func TestMapStorage_Close(t *testing.T) {
	var st MapStorage
	err := st.Close()
	assert.NoError(t, err)
}

func TestMapStorage_Recover(t *testing.T) {
	var conf cmd.Config
	var st MapStorage

	// Recover is an empty method for MapStorage. It only exists to implement Storage interface
	err := st.Recover(conf, slog.New(slog.NewTextHandler(io.Discard, nil)))
	assert.NoError(t, err)
}

func TestMapStorage_GetPositive(t *testing.T) {
	var st = MapStorage{
		mtx: sync.Mutex{},
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	val, err := st.Get(firstKey)
	assert.NoError(t, err)
	assert.Equal(t, firstVal, val)
}

func TestMapStorage_GetNegative(t *testing.T) {
	var st = MapStorage{
		mtx: sync.Mutex{},
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	_, err := st.Get(firstErrKey)
	assert.EqualError(t, err, fmt.Sprintf("key %s not found", firstErrKey))
}

func TestMapStorage_SetPositive(t *testing.T) {
	var st = MapStorage{
		mtx: sync.Mutex{},
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	assert.NoError(t, st.Set(firstKey, firstSetVal))
}

func TestMapStorage_DelPositive(t *testing.T) {
	var st = MapStorage{
		mtx: sync.Mutex{},
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	assert.NoError(t, st.Del(firstKey))
}

func TestMapStorage_DelNegative(t *testing.T) {
	var st = MapStorage{
		mtx: sync.Mutex{},
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	err := st.Del(firstErrKey)
	assert.EqualError(t, err, fmt.Sprintf("key %s not found", firstErrKey))
}
