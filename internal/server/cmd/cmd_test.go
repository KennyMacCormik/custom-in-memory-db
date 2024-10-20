package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type testCase = struct {
	env map[string]string
}

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
	test := testCase{
		env: map[string]string{
			// type Engine struct
			"RAMDB_STORAGE": "map",
			// type Logging struct
			"RAMDB_LOG_FORMAT": "json",
			"RAMDB_LOG_LEVEL":  "warn",
			// type Network struct
			"RAMDB_NET_PROTO":    "http",
			"RAMDB_NET_ADDRESS":  "1.2.3.4",
			"RAMDB_NET_PORT":     "8081",
			"RAMDB_NET_MAX_CONN": "1001",
			"RAMDB_NET_TIMEOUT":  "601s",
			// type Wal struct
			"RAMDB_WAL_BATCH_MAX":     "100",
			"RAMDB_WAL_BATCH_TIMEOUT": "101s",
			"RAMDB_WAL_SEG_SIZE":      "41",
			"RAMDB_WAL_SEG_PATH":      "./",
			"RAMDB_WAL_REPLAY":        "false",
		},
	}

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.NoError(t, err)

	assert.Equal(t, test.env["RAMDB_STORAGE"], conf.Engine.Type)
	// LOG
	assert.Equal(t, test.env["RAMDB_LOG_FORMAT"], conf.Logging.Format)
	assert.Equal(t, test.env["RAMDB_LOG_LEVEL"], conf.Logging.Level)
	// NET
	assert.Equal(t, test.env["RAMDB_NET_PROTO"], conf.Network.Endpoint)
	assert.Equal(t, test.env["RAMDB_NET_ADDRESS"], conf.Network.Host)

	i, err := strconv.Atoi(os.Getenv("RAMDB_NET_PORT"))
	assert.Equal(t, i, conf.Network.Port)
	assert.NoError(t, err)

	i, err = strconv.Atoi(os.Getenv("RAMDB_NET_MAX_CONN"))
	assert.Equal(t, i, conf.Network.MaxConn)
	assert.NoError(t, err)

	d, err := time.ParseDuration(test.env["RAMDB_NET_TIMEOUT"])
	assert.Equal(t, d, conf.Network.Timeout)
	assert.NoError(t, err)
	// WAL
	i, err = strconv.Atoi(os.Getenv("RAMDB_WAL_BATCH_MAX"))
	assert.Equal(t, i, conf.Wal.BatchMax)
	assert.NoError(t, err)

	d, err = time.ParseDuration(test.env["RAMDB_WAL_BATCH_TIMEOUT"])
	assert.Equal(t, d, conf.Wal.BatchTimeout)
	assert.NoError(t, err)

	i, err = strconv.Atoi(test.env["RAMDB_WAL_SEG_SIZE"])
	assert.Equal(t, i*KB, conf.Wal.SegSize)
	assert.NoError(t, err)

	assert.Equal(t, test.env["RAMDB_WAL_SEG_PATH"], conf.Wal.SegPath)

	b, err := strconv.ParseBool(test.env["RAMDB_WAL_REPLAY"])
	assert.Equal(t, b, conf.Wal.Recover)

}

func TestConfig_Positive_AllOptionalMissing(t *testing.T) {
	test := testCase{
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
			// type Wal struct
			"PARSER_EOL":      "",
			"PARSER_TRIM":     "",
			"PARSER_SEP":      "",
			"PARSER_REPBYSEP": "",
			"PARSER_TAG":      "",
		},
	}

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.NoError(t, err)

	assert.Equal(t, "wal", conf.Engine.Type)
	// LOG
	assert.Equal(t, "text", conf.Logging.Format)
	assert.Equal(t, "info", conf.Logging.Level)
	// NET
	assert.Equal(t, "http", conf.Network.Endpoint)
	assert.Equal(t, "0.0.0.0", conf.Network.Host)
	assert.Equal(t, 8080, conf.Network.Port)
	assert.Equal(t, runtime.NumCPU(), conf.Network.MaxConn)

	d, err := time.ParseDuration("1s")
	assert.Equal(t, d, conf.Network.Timeout)
	assert.NoError(t, err)
	// WAL
	assert.Equal(t, runtime.NumCPU(), conf.Wal.BatchMax)

	d, err = time.ParseDuration("1s")
	assert.Equal(t, d, conf.Wal.BatchTimeout)
	assert.NoError(t, err)

	assert.Equal(t, 1*KB, conf.Wal.SegSize)

	path, err := os.Getwd()
	assert.Equal(t, path, conf.Wal.SegPath)
	assert.NoError(t, err)

	b, err := strconv.ParseBool("true")
	assert.Equal(t, b, conf.Wal.Recover)
}

// Engine
func TestConfig_Positive_RAMDB_STORAGE_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_STORAGE": "map",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_STORAGE": "wal",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf := Config{}
		err := conf.New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["RAMDB_STORAGE"], conf.Engine.Type)
		unsetEnv(test.env)
	}
}

func TestConfig_Negative_BogusArg_RAMDB_STORAGE(t *testing.T) {
	test := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			// type Engine struct
			"RAMDB_STORAGE": "mapa",
		},
		err: "config validation error: field 'Type' value 'mapa' invalid, 'oneof=map wal' expected;",
	}

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, test.err)
}

// LOG
func TestConfig_Positive_RAMDB_LOG_FORMAT_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_FORMAT": "text",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_FORMAT": "json",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf := Config{}
		err := conf.New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["RAMDB_LOG_FORMAT"], conf.Logging.Format)
		unsetEnv(test.env)
	}
}

func TestConfig_Positive_RAMDB_LOG_LEVEL_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_LEVEL": "debug",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_LEVEL": "info",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_LEVEL": "warn",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_LOG_LEVEL": "error",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf := Config{}
		err := conf.New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["RAMDB_LOG_LEVEL"], conf.Logging.Level)
		unsetEnv(test.env)
	}
}

func TestConfig_Negative_BogusArg_RAMDB_LOG_FORMAT(t *testing.T) {
	test := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"RAMDB_LOG_FORMAT": "jsona",
		},
	}
	expectedError := "config validation error: field 'Format' value 'jsona' invalid, 'oneof=text json' expected;"

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, expectedError)
}

func TestConfig_Negative_BogusArg_RAMDB_LOG_LEVEL(t *testing.T) {
	test := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"RAMDB_LOG_LEVEL": "warna",
		},
	}
	expectedError := "config validation error: field 'Level' value 'warna' invalid, 'oneof=debug info warn error' expected;"

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, expectedError)
}

// Network
func TestConfig_Positive_RAMDB_NET_PROTO_AllValid(t *testing.T) {
	tests := []testCase{
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_NET_PROTO": "tcp",
			},
		},
		{
			env: map[string]string{
				// type Engine struct
				"RAMDB_NET_PROTO": "http",
			},
		},
	}

	for _, test := range tests {
		setEnv(test.env)
		conf := Config{}
		err := conf.New()
		assert.NoError(t, err)
		assert.Equal(t, test.env["RAMDB_NET_PROTO"], conf.Network.Endpoint)
		unsetEnv(test.env)
	}
}

func TestConfig_Negative_BogusArg_RAMDB_NET_PROTO(t *testing.T) {
	test := struct {
		env map[string]string
	}{
		env: map[string]string{
			// type Engine struct
			"RAMDB_NET_PROTO": "tcpp",
		},
	}
	expectedError := "config validation error: field 'Endpoint' value 'tcpp' invalid, 'oneof=tcp http' expected;"

	setEnv(test.env)
	defer unsetEnv(test.env)

	conf := Config{}
	err := conf.New()
	assert.EqualError(t, err, expectedError)
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
