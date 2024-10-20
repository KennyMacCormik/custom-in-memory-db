package _map

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const firstKey = "1"
const firstVal = "2"
const firstErrKey = "3"
const firstSetVal = "5"

func TestMapStorage_New(t *testing.T) {
	st := New()
	assert.NotNil(t, st)
}

func TestMapStorage_GetPositive(t *testing.T) {
	var st = Storage{
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	val, err := st.Get(firstKey)
	assert.NoError(t, err)
	assert.Equal(t, firstVal, val)
}

func TestMapStorage_GetNegative(t *testing.T) {
	var st = Storage{
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	_, err := st.Get(firstErrKey)
	assert.EqualError(t, err, fmt.Sprintf("key %s not found", firstErrKey))
}

func TestMapStorage_SetPositive(t *testing.T) {
	var st = Storage{
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	assert.NoError(t, st.Set(firstKey, firstSetVal))
}

func TestMapStorage_DelPositive(t *testing.T) {
	var st = Storage{
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	assert.NoError(t, st.Del(firstKey))
}

func TestMapStorage_DelNegative(t *testing.T) {
	var st = Storage{
		m: map[string]string{
			firstKey: firstVal,
		},
	}

	err := st.Del(firstErrKey)
	assert.EqualError(t, err, fmt.Sprintf("key %s not found", firstErrKey))
}
