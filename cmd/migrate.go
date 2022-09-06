package cmd

import (
	"context"
	"url-shortner/config"
	"url-shortner/platform/data"
	"url-shortner/platform/observation/logging"

	"github.com/spf13/cobra"
)

const (
	migrateLong = "Connects to the database and runs the necessary database migrations configured by the application .yaml file\n" +
		"located at ./.config.yaml, or any file type supported by https://github.com/spf13/viper.\n"
	migrationsFilePath = "file://resources/migrations"
)

func newMigrateUpCommand(conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Executes database migrations.",
		Long:  migrateLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := logging.WithLogger(context.Background())
			if err := data.MigrateUp(ctx, mapDBConfig(conf), migrationsFilePath); err != nil {
				logging.FromContext(ctx).Fatal("error while migration", err)
			}

			return nil
		},
	}
}

func mapDBConfig(conf config.Config) data.DBConfig {
	return data.DBConfig{
		Hostname: conf.Database.Hostname,
		Port:     conf.Database.Port,
		Database: conf.Database.DatabaseName,
		Username: conf.Database.Username,
		Password: conf.Database.Password,
		DebugLog: conf.Database.DebugLog,
	}
}
