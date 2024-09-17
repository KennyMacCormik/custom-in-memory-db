package server

import (
	"errors"
	"os"
	"testing"
	"time"
)

type Env struct {
	Storage        string
	Address        string
	Port           string
	MaxConn        string
	MaxMessageSize string
	IdleTimeout    string
	Format         string
}

type TestCase struct {
	Envs     Env
	Expected Config
	Error    error
}

// ENV vars present but contains wrong data
func TestConfigAllEnvPresent(t *testing.T) {
	var testCases = []TestCase{
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			nil,
		},
		{
			Env{
				"map1",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("STORAGE expected oneof=map, got map1"),
		},
		{
			Env{
				"map",
				"127.0.0.1111",
				"8080",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("ADDR expected ip4_addr, got 127.0.0.1111"),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080a",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("PORT expected number, got 8080a"),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100a",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("MAX_CONN expected number, got 100a"),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4a",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("MESSAGE_SIZE expected number, got 4a"),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("TIMEOUT expected duration, got 60"),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60s",
				"texts",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("LOG_FORMAT expected oneof=text, got texts"),
		},
	}

	for _, val := range testCases {
		setEnv(val)
		conf := Config{}
		err := conf.New()
		validate(err, val, conf, t)
		unsetEnv()
	}
}

// ENV vars missing
func TestConfigEnvMissing(t *testing.T) {
	var testCases = []TestCase{
		{
			Env{
				"",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("STORAGE expected oneof=map, got "),
		},
		{
			Env{
				"map",
				"",
				"8080",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("ADDR expected ip4_addr, got "),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"",
				"100",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			errors.New("PORT expected number, got "),
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"",
				"4",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			nil,
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"",
				"60s",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			nil,
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"",
				"text",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			nil,
		},
		{
			Env{
				"map",
				"127.0.0.1",
				"8080",
				"100",
				"4",
				"60s",
				"",
			},
			Config{
				Engine{"map"},
				Network{
					"127.0.0.1",
					8080,
					100,
					4,
					60 * time.Second,
				},
				Logging{"text"},
			},
			nil,
		},
	}

	for _, val := range testCases {
		setEnv(val)
		conf := Config{}
		err := conf.New()
		validate(err, val, conf, t)
		unsetEnv()
	}
}

func setEnv(t TestCase) {
	if t.Envs.Storage != "" {
		_ = os.Setenv("STORAGE", t.Envs.Storage)
	}
	if t.Envs.Address != "" {
		_ = os.Setenv("ADDR", t.Envs.Address)
	}
	if t.Envs.Port != "" {
		_ = os.Setenv("PORT", t.Envs.Port)
	}
	if t.Envs.MaxConn != "" {
		_ = os.Setenv("MAX_CONN", t.Envs.MaxConn)
	}
	if t.Envs.MaxMessageSize != "" {
		_ = os.Setenv("MESSAGE_SIZE", t.Envs.MaxMessageSize)
	}
	if t.Envs.IdleTimeout != "" {
		_ = os.Setenv("TIMEOUT", t.Envs.IdleTimeout)
	}
	if t.Envs.Format != "" {
		_ = os.Setenv("LOG_FORMAT", t.Envs.Format)
	}
}

func unsetEnv() {
	_ = os.Unsetenv("STORAGE")
	_ = os.Unsetenv("ADDR")
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("MAX_CONN")
	_ = os.Unsetenv("MESSAGE_SIZE")
	_ = os.Unsetenv("TIMEOUT")
	_ = os.Unsetenv("LOG_FORMAT")
}

func validate(err error, val TestCase, conf Config, t *testing.T) {
	// Err - Expected error -
	if err == nil && val.Error == nil {
		if val.Expected != conf {
			t.Errorf("Expected value %v, got value %v", val.Expected, conf)
		}
	}
	// Err - Expected error +
	if err == nil && val.Error != nil {
		t.Errorf("Expected error %v, got value %v", val.Error, val)
	}
	// Err + Expected error -
	if err != nil && val.Error == nil {
		t.Errorf("Expected value %v, got error %v", val.Expected, err)
	}
	// Err + Expected error +
	if err != nil && val.Error != nil {
		if err.Error() != val.Error.Error() {
			t.Errorf("Expected error %v, got error %v", val.Error, err)
		}
	}
}
