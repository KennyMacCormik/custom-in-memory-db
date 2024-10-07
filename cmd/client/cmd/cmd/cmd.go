package cmd

import (
	"custom-in-memory-db/cmd/client/cmd/cmd/del"
	"custom-in-memory-db/cmd/client/cmd/cmd/get"
	"custom-in-memory-db/cmd/client/cmd/cmd/set"
	"custom-in-memory-db/cmd/client/cmd/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Init(cfg *conf.Config) *cobra.Command {
	var cmd = cobra.Command{
		Use:   "cmd [COMMANDS]",
		Short: "Executes db command",
		Long:  `Executes db command qq ww ee`,
	}

	viper.BindPFlag("server", cmd.PersistentFlags().Lookup("server"))
	cmd.PersistentFlags().StringVarP(&cfg.Server, "server", "s", "127.0.0.1", "database server name")

	viper.BindPFlag("port", cmd.PersistentFlags().Lookup("port"))
	cmd.PersistentFlags().IntVarP(&cfg.Port, "port", "p", 8080, "database port")

	cmd.AddCommand(get.Init(cfg), set.Init(cfg), del.Init(cfg))

	return &cmd
}
