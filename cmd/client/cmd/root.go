package root

import (
	"custom-in-memory-db/cmd/client/cmd/cmd"
	"custom-in-memory-db/cmd/client/cmd/conf"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

func Execute(cfg *conf.Config) {
	var rootCmd = &cobra.Command{
		Use: "ramdb-cli [OPTIONS] [COMMANDS]",
	}

	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "verbose output")

	rootCmd.AddCommand(cmd.Init(cfg))

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
