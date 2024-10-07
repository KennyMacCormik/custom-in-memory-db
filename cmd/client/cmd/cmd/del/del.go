package del

import (
	"custom-in-memory-db/cmd/client/cmd/conf"
	"custom-in-memory-db/cmd/client/shared"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const validatorTag = "printascii,containsany=*_/|alphanum|numeric|alpha"
const argsExpected = 1

func Init(cfg *conf.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "del <key>",
		Short: "Deletes <key> from the database",
		Long:  fmt.Sprintf("Deletes <key> and the associated value from the database\n<key> expected to match %s pattern", validatorTag),
		Args:  args,
		Run:   run,
	}

	viper.BindPFlag("timeout", cmd.PersistentFlags().Lookup("port"))
	cmd.PersistentFlags().DurationVarP(&cfg.Timeout, "timeout", "t", 1*time.Second, "connection timeout")

	return cmd
}

func run(cmd *cobra.Command, args []string) {
	server, _ := cmd.Flags().GetString("server")
	port, _ := cmd.Flags().GetInt("port")
	timeout, _ := cmd.Flags().GetDuration("timeout")
	verbose, _ := cmd.Flags().GetBool("verbose")
	cfg := conf.Config{
		Server:  server,
		Port:    port,
		Timeout: timeout,
		Verbose: verbose,
	}

	fmt.Println(string(shared.InvokeTcp(cfg, strings.Join([]string{"DEL ", args[0], "\n"}, ""))))
}

func args(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(argsExpected)(cmd, args); err != nil {
		return err
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	for _, arg := range args {
		if err := validate.Var(args[0], validatorTag); err != nil {
			return fmt.Errorf("arg [%s] expected to match [%s] tag", arg, validatorTag)
		}
	}

	return nil
}
