package cmd

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"reflect"
	"strings"
	"time"
)

const KB = 1024

type Engine struct {
	APP_STORAGE string `env:"APP_STORAGE" env-required:"true" env-description:"storage driver" validate:"oneof=mem wal"`
	APP_INPUT   string `env:"APP_INPUT" env-default:"tcp4" env-description:"how server accept commands" validate:"oneof=stdin tcp4"`
}

type Replication struct {
	REP_TYPE     string        `env:"REP_TYPE" env-required:"true" env-description:"storage driver" validate:"oneof=master slave"`
	REP_ADDR     string        `env:"REP_ADDR" env-description:"address of the master server. For REP_TYPE=master might be empty" validate:"required_if=REP_TYPE slave,omitempty,ip4_addr"`
	REP_PORT     int           `env:"REP_PORT" env-required:"true" env-description:"opens this port on master for the replication. For replica this means port to connect to" validate:"numeric,gt=0,lt=65536"`
	REP_INTERVAL time.Duration `env:"REP_INTERVAL" env-default:"1s" env-description:"replication requests will happen this often. min 1s. Required for REP_TYPE=slave" validate:"required_if=REP_TYPE slave,omitempty,min=1s"`
}

type Network struct {
	NET_ADDR         string        `env:"NET_ADDR" env-required:"true" env-description:"address to listen" validate:"ip4_addr"`
	NET_PORT         int           `env:"NET_PORT" env-required:"true" env-description:"port to listen" validate:"numeric,gt=0,lt=65536"`
	NET_MAX_CONN     int           `env:"NET_MAX_CONN" env-default:"100" env-description:"maximum accepted connections" validate:"numeric,gte=0"`
	NET_MESSAGE_SIZE int           `env:"NET_MESSAGE_SIZE" env-default:"4" env-description:"max message size KB" validate:"numeric,gt=0"`
	NET_TIMEOUT      time.Duration `env:"NET_TIMEOUT" env-default:"60s" env-description:"idle connection timeout min 1ms" validate:"min=1ms"`
}

type Logging struct {
	LOG_FORMAT string `env:"LOG_FORMAT" env-default:"text" env-description:"log format" validate:"oneof=text json"`
	LOG_LEVEL  string `env:"LOG_LEVEL" env-default:"debug" env-description:"log level" validate:"oneof=debug info warn error"`
}

type Wal struct {
	WAL_BATCH_SIZE    int           `env:"WAL_BATCH_SIZE" env-default:"10" env-description:"connection amount to trigger flush" validate:"numeric,gt=0"`
	WAL_BATCH_TIMEOUT time.Duration `env:"WAL_BATCH_TIMEOUT" env-default:"1s" env-description:"batch flush timeout min 1s" validate:"min=1ms"`
	WAL_SEG_SIZE      int           `env:"WAL_SEG_SIZE" env-default:"1" env-description:"segment size on disk KB" validate:"numeric,gt=0"`
	WAL_SEG_PATH      string        `env:"WAL_SEG_PATH" env-default:"./" env-description:"segment folder" validate:"dir,dirpath"`
	WAL_SEG_RECOVER   bool          `env:"WAL_SEG_RECOVER" env-default:"true" env-description:"load data from wal to ram" validate:"boolean"`
}

type Config struct {
	Engine Engine
	Rep    Replication
	Net    Network
	Log    Logging
	Wal    Wal
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
		return fmt.Errorf("config validation error: %s", c.handleValidatorError(err))
	}

	c.Net.NET_MESSAGE_SIZE *= KB
	c.Wal.WAL_SEG_SIZE *= KB

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
