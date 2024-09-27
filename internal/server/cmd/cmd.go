package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Engine struct {
	Storage string `env:"STORAGE" env-required:"true" env-description:"storage driver"`
	Input   string `env:"INPUT" env-default:"tcp4" env-description:"how server accept commands"`
}

type Network struct {
	Address        string        `env:"ADDR" env-required:"true" env-description:"address to listen"`
	Port           int           `env:"PORT" env-required:"true" env-description:"port to listen"`
	MaxConn        int           `env:"MAX_CONN" env-default:"100" env-description:"maximum accepted connections"`
	MaxMessageSize int           `env:"MESSAGE_SIZE" env-default:"4" env-description:"max message size KB"`
	IdleTimeout    time.Duration `env:"TIMEOUT" env-default:"60s" env-description:"idle connection timeout"`
}

type Logging struct {
	Format string `env:"LOG_FORMAT" env-default:"text" env-description:"log format"`
}

type Config struct {
	Eng Engine
	Net Network
	Log Logging
}

func validateMandatory() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	// ENV
	errMsg := "%s expected %s, got %s"

	val := "STORAGE"
	tag := "oneof=map"
	env := os.Getenv(val)
	err := validate.Var(env, tag)
	if err != nil {
		return fmt.Errorf(errMsg, val, tag, env)
	}

	val = "ADDR"
	tag = "ip4_addr"
	env = os.Getenv(val)
	err = validate.Var(env, tag)
	if err != nil {
		return fmt.Errorf(errMsg, val, tag, env)
	}

	val = "PORT"
	tag = "number"
	env = os.Getenv(val)
	err = validate.Var(env, tag)
	if err != nil {
		return fmt.Errorf(errMsg, val, tag, env)
	}

	return nil
}

func validateOptional() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	// ENV
	errMsg := "%s expected %s, got %s"
	err := error(nil)

	val := "MAX_CONN"
	tag := "number"
	env := os.Getenv(val)
	if env != "" {
		err = validate.Var(env, tag)
		if err != nil {
			return fmt.Errorf(errMsg, val, tag, env)
		}
	}

	val = "INPUT"
	tag = "oneof=stdin tcp4"
	env = os.Getenv(val)
	if env != "" {
		err = validate.Var(env, tag)
		if err != nil {
			return fmt.Errorf(errMsg, val, tag, env)
		}
	}

	val = "MESSAGE_SIZE"
	tag = "number"
	env = os.Getenv(val)
	if env != "" {
		err = validate.Var(env, tag)
		if err != nil {
			return fmt.Errorf(errMsg, val, tag, env)
		}
	}

	val = "TIMEOUT"
	tag = "duration"
	env = os.Getenv(val)
	if env != "" {
		_, err = time.ParseDuration(env)
		if err != nil {
			return fmt.Errorf(errMsg, val, tag, env)
		}
	}

	val = "LOG_FORMAT"
	tag = "oneof=text"
	env = os.Getenv(val)
	if env != "" {
		err = validate.Var(env, tag)
		if err != nil {
			return fmt.Errorf(errMsg, val, tag, env)
		}
	}

	return nil
}

func (c *Config) New() error {
	err := validateMandatory()
	if err != nil {
		return err
	}
	err = validateOptional()
	if err != nil {
		return err
	}
	err = cleanenv.ReadEnv(c)
	if err != nil {
		return fmt.Errorf("could not read config from ENV: %w", err)
	}
	return nil
}
