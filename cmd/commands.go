package cmd

import (
	"url-shortner/config"
	"url-shortner/platform/observation/logging"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCLI() *cobra.Command {
	cli := &cobra.Command{
		Use:   "url-shortener",
		Short: "A simple URL shortener",
	}

	var conf config.Config
	if err := conf.Load(); err != nil {
		logging.WithFields(zap.Error(err)).
			Fatal("error while reading config")
	}

	cli.AddCommand(newStartServerCommand(conf))
	cli.AddCommand(newMigrateUpCommand(conf))

	return cli
}
