package conf

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server  string
	Port    int
	Timeout time.Duration
	Verbose bool
}

func InitConf() {
	viper.SetDefault("server", "127.0.0.1")
	viper.SetDefault("port", 8080)
	viper.SetDefault("timeout", "1s")
	viper.SetDefault("Verbose", false)

	viper.AddConfigPath(".")
	viper.SetConfigName("conf")

	viper.SetEnvPrefix("MEMDB_CLT")
	viper.AutomaticEnv()
}
