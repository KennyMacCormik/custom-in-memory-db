package cmd

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const KB = 1024

type Engine struct {
	// underlying storage. defaults to wal
	Type string `mapstructure:"storage" validate:"oneof=map wal"`
}

type Logging struct {
	// log format. defaults to text
	Format string `mapstructure:"log_format" validate:"oneof=text json"`
	// log level. defaults to info
	Level string `mapstructure:"log_level" validate:"oneof=debug info warn error"`
}

type Network struct {
	// network protocol to work with. defaults to http
	Endpoint string `mapstructure:"net_proto" validate:"oneof=tcp http"`
	// address to listen. defaults to 0.0.0.0
	Host string `mapstructure:"net_address" validate:"ip4_addr"`
	// port to listen. defaults to 8080
	Port int `mapstructure:"net_port" validate:"numeric,gt=0,lt=65536"`
	// maximum accepted connections. Defaults to runtime.NumCPU()
	MaxConn int `mapstructure:"net_max_conn" validate:"numeric,gte=0"`
	// idle connection timeout min 1ms. defaults to 1s
	Timeout time.Duration `mapstructure:"net_timeout" validate:"min=1ms"`
}

type Wal struct {
	// max conn collected before writing to wal. defaults to runtime.NumCPU()
	BatchMax int `mapstructure:"wal_batch_max" validate:"numeric,gt=0"`
	// batch flush timeout, min 1ms. defaults to 1s
	BatchTimeout time.Duration `mapstructure:"wal_batch_timeout" validate:"min=1ms"`
	// segment size on disk KB, min 1. defaults to 1
	SegSize int `mapstructure:"wal_seg_size" validate:"numeric,gt=0"`
	// segment folder. defaults to os.Getwd
	SegPath string `mapstructure:"wal_seg_path" validate:"dir"`
	// recover from wal on db start. defaults to true
	Recover bool `mapstructure:"wal_replay" validate:"boolean"`
}

type Config struct {
	Engine  Engine  `mapstructure:",squash"`
	Network Network `mapstructure:",squash"`
	Logging Logging `mapstructure:",squash"`
	Wal     Wal     `mapstructure:",squash"`
}

func New() (Config, error) {
	c := Config{}

	err := c.loadEnv()
	if err != nil {
		return Config{}, fmt.Errorf("config unmarshalling error: %w", err)
	}

	err = c.validate()
	if err != nil {
		return Config{}, fmt.Errorf("config validation error: %w", errors.New(c.handleValidatorError(err)))
	}

	c.Wal.SegSize *= KB

	return c, nil
}

func (c *Config) loadEnv() error {
	viper.SetEnvPrefix("ramdb")

	c.setEngineEnv()
	c.setLoggingEnv()
	c.setNetworkEnv()
	c.setWalEnv()

	viper.AutomaticEnv()
	return viper.Unmarshal(c)
}

func (c *Config) setWalEnv() {
	viper.SetDefault("wal_batch_max", strconv.Itoa(runtime.NumCPU()))
	_ = viper.BindEnv("wal_batch_max")

	viper.SetDefault("wal_batch_timeout", "1s")
	_ = viper.BindEnv("wal_batch_timeout")

	viper.SetDefault("wal_seg_size", "1")
	_ = viper.BindEnv("wal_seg_size")

	path, err := os.Getwd()
	if err != nil {
		if runtime.GOOS == "windows" {
			path = ".\\"
		} else {
			path = "./"
		}
	}
	viper.SetDefault("wal_seg_path", path)
	_ = viper.BindEnv("wal_seg_path")

	viper.SetDefault("wal_replay", "true")
	_ = viper.BindEnv("wal_replay")
}

func (c *Config) setNetworkEnv() {
	viper.SetDefault("net_proto", "http")
	_ = viper.BindEnv("net_proto")

	viper.SetDefault("net_address", "0.0.0.0")
	_ = viper.BindEnv("net_address")

	viper.SetDefault("net_port", "8080")
	_ = viper.BindEnv("net_port")

	viper.SetDefault("net_max_conn", strconv.Itoa(runtime.NumCPU()))
	_ = viper.BindEnv("net_max_conn")

	viper.SetDefault("net_timeout", "1s")
	_ = viper.BindEnv("net_timeout")
}

func (c *Config) setLoggingEnv() {
	viper.SetDefault("log_format", "text")
	_ = viper.BindEnv("format")

	viper.SetDefault("log_level", "info")
	_ = viper.BindEnv("level")
}

func (c *Config) setEngineEnv() {
	viper.SetDefault("storage", "wal")
	_ = viper.BindEnv("storage")
}

func (c *Config) validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(c)
	if err != nil {
		return err
	}

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
		if field.Type.Name() != "bool" {
			for j := 0; j < field.Type.NumField(); j++ {
				intFieldName := field.Type.Field(j)
				if intFieldName.Name == sf {
					return intFieldName.Tag.Get("validate")
				}
			}
		}
	}

	return ""
}
