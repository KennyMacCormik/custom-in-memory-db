package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			_ = os.Setenv(k, v)
		}
	}
}

func unsetEnv(list map[string]string) {
	for k, v := range list {
		if v != "" {
			_ = os.Unsetenv(k)
		}
	}
}

func TestConfig_Positive_AllPresent(t *testing.T) {
	testCase := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "map",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ENDPOINT":     "http",
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
			// type Wal struct
			"PARSER_EOL":      "10",
			"PARSER_TRIM":     " \\t\\n",
			"PARSER_SEP":      " ",
			"PARSER_REPBYSEP": "\\t",
			"PARSER_TAG":      "alphanum|numeric|alpha|printascii,containsany=*_/",
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
			"APP_STORAGE": "map",
			"APP_INPUT":   "",
			// type Network struct
			"NET_ENDPOINT":     "",
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
			// type Wal struct
			"PARSER_EOL":      "",
			"PARSER_TRIM":     "",
			"PARSER_SEP":      "",
			"PARSER_REPBYSEP": "",
			"PARSER_TAG":      "",
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
			"NET_ENDPOINT":     "",
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
		err: "could not read config from ENV: field \"Type\" is required but the value is not provided",
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
			"NET_ENDPOINT":     "",
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
		err: "could not read config from ENV: field \"Host\" is required but the value is not provided",
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
			"NET_ENDPOINT":     "",
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
		err: "could not read config from ENV: field \"Port\" is required but the value is not provided",
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
			"NET_ENDPOINT":     "",
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
		err: "config validation error: field 'Type' value 'memo' invalid, 'oneof=map wal' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

// Network
func TestConfig_Negative_BogusArg_NET_ENDPOINT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "map",
			"APP_INPUT":   "tcp4",
			// type Network struct
			"NET_ENDPOINT":     "tcpp",
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
		err: "config validation error: field 'Endpoint' value 'tcpp' invalid, 'oneof=tcp http' expected; field 'Host' value '0.0.0.1111' invalid, 'ip4_addr' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_NET_ADDR(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Host' value '0.0.0.1111' invalid, 'ip4_addr' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Port' value '%!s(int=70000)' invalid, 'numeric,gt=0,lt=65536' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Port' value '%!s(int=-10)' invalid, 'numeric,gt=0,lt=65536' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'MaxConn' value '%!s(int=-10)' invalid, 'numeric,gte=0' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Timeout' value '0s' invalid, 'min=1ms' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Format' value 'textt' invalid, 'oneof=text json' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'Level' value 'debugg' invalid, 'oneof=debug info warn error' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'BatchMax' value '%!s(int=-100)' invalid, 'numeric,gt=0' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'BatchTimeout' value '10ns' invalid, 'min=1ms' expected;",
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
			"APP_STORAGE": "map",
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
		err: "could not read config from ENV: parsing field BatchTimeout env WAL_BATCH_TIMEOUT: time: missing unit in duration \"10\"",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'SegSize' value '%!s(int=-4)' invalid, 'numeric,gt=0' expected;",
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
			"APP_STORAGE": "map",
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
		err: "config validation error: field 'SegPath' value './q' invalid, 'dir,dirpath' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, testCase.err)
}
