package cmd

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"reflect"
	"strings"
	"time"
)

const KB = 1024

type Engine struct {
	Type string `env:"APP_STORAGE" env-required:"true" env-description:"storage driver" validate:"oneof=map wal"`
}

type Logging struct {
	Format string `env:"LOG_FORMAT" env-default:"text" env-description:"log format" validate:"oneof=text json"`
	Level  string `env:"LOG_LEVEL" env-default:"debug" env-description:"log level" validate:"oneof=debug info warn error"`
}

type Network struct {
	Endpoint string        `env:"NET_ENDPOINT" env-default:"http" env-description:"network protocol to work with" validate:"oneof=tcp http"`
	Host     string        `env:"NET_ADDR" env-required:"true" env-description:"address to listen" validate:"ip4_addr"`
	Port     int           `env:"NET_PORT" env-required:"true" env-description:"port to listen" validate:"numeric,gt=0,lt=65536"`
	MaxConn  int           `env:"NET_MAX_CONN" env-default:"1" env-description:"maximum accepted connections" validate:"numeric,gte=0"`
	Timeout  time.Duration `env:"NET_TIMEOUT" env-default:"60s" env-description:"idle connection timeout min 1ms" validate:"min=1ms"`
}

type Wal struct {
	BatchMax     int           `env:"WAL_BATCH_SIZE" env-default:"10" env-description:"max conn collected before writing to wal" validate:"numeric,gt=0"`
	BatchTimeout time.Duration `env:"WAL_BATCH_TIMEOUT" env-default:"1s" env-description:"batch flush timeout, min 1ms" validate:"min=1ms"`
	SegSize      int           `env:"WAL_SEG_SIZE" env-default:"1" env-description:"segment size on disk KB, min 1" validate:"numeric,gt=0"`
	SegPath      string        `env:"WAL_SEG_PATH" env-default:"./" env-description:"segment folder" validate:"dir,dirpath"`
	Recover      bool          `env:"WAL_SEG_RECOVER" env-default:"true" env-description:"recover from wal on db start" validate:"boolean"`
}

// Parser struct contains args for Parser interface.
// Change this at your own risk
type Parser struct {
	Eol            byte   `env:"PARSER_EOL" env-default:"10" env-description:"symbol representing end of command"`
	Trim           string `env:"PARSER_TRIM" env-default:" \t\n" env-description:"trim set for each arg and command itself"`
	Sep            string `env:"PARSER_SEP" env-default:" " env-description:"separator between args and command itself"`
	ToReplaceBySep string `env:"PARSER_REPBYSEP" env-default:"\t" env-description:"replace set for separator"`
	Tag            string "env:\"PARSER_TAG\" env-default:\"alphanum|numeric|alpha|containsany=*_/,excludesall=!\"#$%&'()+0x2C-.:;<=>?@[]^`{}0x7C~,printascii\" env-description:\"default tag for validator\""
}

type Config struct {
	Engine  Engine
	Network Network
	Logging Logging
	Wal     Wal
	Parser  Parser
}

func (c *Config) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) New() error {
	err := cleanenv.ReadEnv(c)
	if err != nil {
		return fmt.Errorf("could not read config from ENV: %w", err)
	}

	err = c.validate()
	if err != nil {
		return fmt.Errorf("config validation error: %w", errors.New(c.handleValidatorError(err)))
	}

	c.Wal.SegSize *= KB

	return nil
}

func (c *Config) handleValidatorError(err error) string {
	valErr := err.(validator.ValidationErrors)
	errStr := ""

	for _, v := range valErr {
		tag := c.reflectActualTag(v.StructField())
		if tag == "" {
			tag = "err reflect tag"
		}
		errStr += fmt.Sprintf("field '%s' value '%s' invalid, '%s' expected; ",
			v.StructField(), v.Value(), tag)
	}
	errStr = strings.Trim(errStr, " ")

	return errStr
}

func (c *Config) reflectActualTag(sf string) string {
	ref := reflect.TypeOf(*c)

	for i := 0; i < ref.NumField(); i++ {
		fieldName := ref.Field(i).Name
		field, _ := ref.FieldByName(fieldName)
		for j := 0; j < field.Type.NumField(); j++ {
			intFieldName := field.Type.Field(j)
			if intFieldName.Name == sf {
				return intFieldName.Tag.Get("validate")
			}
		}
	}

	return ""
}
