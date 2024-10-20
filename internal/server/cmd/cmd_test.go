package cmd

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

type testCase struct {
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

	conf, err := New()
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

	conf, err := New()
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
		conf, err := New()
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

	conf, err := New()
	assert.EqualError(t, err, test.err)
	assert.Equal(t, Config{}, conf)
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
		conf, err := New()
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
		conf, err := New()
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

	conf, err := New()
	assert.EqualError(t, err, expectedError)
	assert.Equal(t, Config{}, conf)
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

	conf, err := New()
	assert.Equal(t, Config{}, conf)
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
		conf, err := New()
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

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, expectedError)
}

func TestConfig_Negative_BogusArg_RAMDB_NET_ADDRESS(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_NET_ADDRESS": "0.0.0.1111",
		},
		err: "config validation error: field 'Host' value '0.0.0.1111' invalid, 'ip4_addr' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_NET_PORT_More65536(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_NET_PORT": "70000",
		},
		err: "config validation error: field 'Port' value '%!s(int=70000)' invalid, 'numeric,gt=0,lt=65536' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_NET_PORT_Less_0(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_NET_PORT": "-10",
		},
		err: "config validation error: field 'Port' value '%!s(int=-10)' invalid, 'numeric,gt=0,lt=65536' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_NET_MAX_CONN(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_NET_MAX_CONN": "-10",
		},
		err: "config validation error: field 'MaxConn' value '%!s(int=-10)' invalid, 'numeric,gte=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_NET_TIMEOUT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_NET_TIMEOUT": "0s",
		},
		err: "config validation error: field 'Timeout' value '0s' invalid, 'min=1ms' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

// Wal

func TestConfig_Negative_BogusArg_RAMDB_WAL_BATCH_MAX(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_BATCH_MAX": "-100",
		},
		err: "config validation error: field 'BatchMax' value '%!s(int=-100)' invalid, 'numeric,gt=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_WAL_BATCH_TIMEOUT(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_BATCH_TIMEOUT": "10ns",
		},
		err: "config validation error: field 'BatchTimeout' value '10ns' invalid, 'min=1ms' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_WAL_BATCH_TIMEOUT_BogusDuration(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_BATCH_TIMEOUT": "4",
		},
		err: "config unmarshalling error: 1 error(s) decoding:\n\n* error decoding 'wal_batch_timeout': time: missing unit in duration \"4\"",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_WAL_SEG_SIZE(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_SEG_SIZE": "-4",
		},
		err: "config validation error: field 'SegSize' value '%!s(int=-4)' invalid, 'numeric,gt=0' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_WAL_SEG_PATH(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_SEG_PATH": "./q",
		},
		err: "config validation error: field 'SegPath' value './q' invalid, 'dir' expected;",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}

func TestConfig_Negative_BogusArg_RAMDB_WAL_REPLAY(t *testing.T) {
	testCase := struct {
		env map[string]string
		err string
	}{
		env: map[string]string{
			"RAMDB_WAL_REPLAY": "truee",
		},
		err: "config unmarshalling error: 1 error(s) decoding:\n\n* cannot parse 'wal_replay' as bool: strconv.ParseBool: parsing \"truee\": invalid syntax",
	}

	setEnv(testCase.env)
	defer unsetEnv(testCase.env)

	conf, err := New()
	assert.Equal(t, Config{}, conf)
	assert.EqualError(t, err, testCase.err)
}
