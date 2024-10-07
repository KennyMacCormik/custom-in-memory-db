package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			os.Setenv(k, v)
		}
	}
}

func unsetEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			os.Unsetenv(k)
		}
	}
}

func TestConfig_Positive_AllPresent(t *testing.T) {
	testCase := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.NoError(t, err)
}

func TestConfig_Positive_AllOptionalMissing(t *testing.T) {
	testCase := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "",
			"NET_MESSAGE_SIZE": "",
			"NET_TIMEOUT":      "",
			// type Logging struct
			"LOG_FORMAT": "",
			"LOG_LEVEL":  "",
			// type Wal struct
			"WAL_BATCH_SIZE":    "",
			"WAL_BATCH_TIMEOUT": "",
			"WAL_SEG_SIZE":      "",
			"WAL_SEG_PATH":      "",
		},
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.NoError(t, err)
}

func TestConfig_Negative_AllMissing(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "",
			"APP_INPUT":   "",
			// type Network struct
			"NET_ADDR":         "",
			"NET_PORT":         "",
			"NET_MAX_CONN":     "",
			"NET_MESSAGE_SIZE": "",
			"NET_TIMEOUT":      "",
			// type Logging struct
			"LOG_FORMAT": "",
			"LOG_LEVEL":  "",
			// type Wal struct
			"WAL_BATCH_SIZE":    "",
			"WAL_BATCH_TIMEOUT": "",
			"WAL_SEG_SIZE":      "",
			"WAL_SEG_PATH":      "",
		},
		err: "could not read config from ENV: field \"APP_STORAGE\" is required but the value is not provided",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_NET_ADDR_Missing(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "map",
			"APP_INPUT":   "",
			// type Network struct
			"NET_ADDR":         "",
			"NET_PORT":         "",
			"NET_MAX_CONN":     "",
			"NET_MESSAGE_SIZE": "",
			"NET_TIMEOUT":      "",
			// type Logging struct
			"LOG_FORMAT": "",
			"LOG_LEVEL":  "",
			// type Wal struct
			"WAL_BATCH_SIZE":    "",
			"WAL_BATCH_TIMEOUT": "",
			"WAL_SEG_SIZE":      "",
			"WAL_SEG_PATH":      "",
		},
		err: "could not read config from ENV: field \"NET_ADDR\" is required but the value is not provided",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_NET_PORT_Missing(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "map",
			"APP_INPUT":   "",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "",
			"NET_MAX_CONN":     "",
			"NET_MESSAGE_SIZE": "",
			"NET_TIMEOUT":      "",
			// type Logging struct
			"LOG_FORMAT": "",
			"LOG_LEVEL":  "",
			// type Wal struct
			"WAL_BATCH_SIZE":    "",
			"WAL_BATCH_TIMEOUT": "",
			"WAL_SEG_SIZE":      "",
			"WAL_SEG_PATH":      "",
		},
		err: "could not read config from ENV: field \"NET_PORT\" is required but the value is not provided",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

// Engine

func TestConfig_Negative_BogusArg_APP_STORAGE(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "memo",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'APP_STORAGE' value 'memo' invalid, 'oneof=mem wal' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_APP_INPUT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp44",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'APP_INPUT' value 'tcp44' invalid, 'oneof=stdin tcp4' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

// Network

func TestConfig_Negative_BogusArg_NET_ADDR(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.1111",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_ADDR' value '0.0.0.1111' invalid, 'ip4_addr' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_PORT_More65536(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "70000",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_PORT' value '%!s(int=70000)' invalid, 'numeric,gt=0,lt=65536' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_PORT_Less_0(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "-10",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_PORT' value '%!s(int=-10)' invalid, 'numeric,gt=0,lt=65536' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_MAX_CONN(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "-10",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_MAX_CONN' value '%!s(int=-10)' invalid, 'numeric,gte=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_MESSAGE_SIZE(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "-4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_MESSAGE_SIZE' value '%!s(int=-4)' invalid, 'numeric,gt=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_TIMEOUT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "0s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'NET_TIMEOUT' value '0s' invalid, 'min=1ms' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

// Logging

func TestConfig_Negative_BogusArg_LOG_FORMAT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "textt",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'LOG_FORMAT' value 'textt' invalid, 'oneof=text json' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_LOG_LEVEL(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "10s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debugg",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'LOG_LEVEL' value 'debugg' invalid, 'oneof=debug info warn error' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

// Wal

func TestConfig_Negative_BogusArg_WAL_BATCH_SIZE(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "-100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'WAL_BATCH_SIZE' value '%!s(int=-100)' invalid, 'numeric,gt=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_WAL_BATCH_TIMEOUT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "10s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10ns",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'WAL_BATCH_TIMEOUT' value '10ns' invalid, 'min=1ms' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_WAL_BATCH_TIMEOUT_BogusDuration(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "10s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./",
		},
		err: "could not read config from ENV: parsing field WAL_BATCH_TIMEOUT env WAL_BATCH_TIMEOUT: time: missing unit in duration \"10\"",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_WAL_SEG_SIZE(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "60s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "-4",
			"WAL_SEG_PATH":      "./",
		},
		err: "config validation error: field 'WAL_SEG_SIZE' value '%!s(int=-4)' invalid, 'numeric,gt=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_WAL_SEG_PATH(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "mem",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ADDR":         "0.0.0.0",
			"NET_PORT":         "8080",
			"NET_MAX_CONN":     "100",
			"NET_MESSAGE_SIZE": "4",
			"NET_TIMEOUT":      "10s",
			// type Logging struct
			"LOG_FORMAT": "text",
			"LOG_LEVEL":  "debug",
			// type Wal struct
			"WAL_BATCH_SIZE":    "100",
			"WAL_BATCH_TIMEOUT": "10s",
			"WAL_SEG_SIZE":      "4",
			"WAL_SEG_PATH":      "./q",
		},
		err: "config validation error: field 'WAL_SEG_PATH' value './q' invalid, 'dir,dirpath' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}
