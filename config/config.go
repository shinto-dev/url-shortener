package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/iamolegga/enviper"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Environment string

type HTTPServer struct {
	Port int
}

type Database struct {
	Hostname string
	Port     int

	DatabaseName string
	Username     string
	Password     string
	DebugLog     bool
}

type Config struct {
	HTTPServer  HTTPServer
	Environment Environment
	Database    Database
}

func (c *Config) Load() error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "unable to get current working directory")
	}

	e := enviper.New(viper.New())
	e.AddConfigPath(fmt.Sprintf("%s/resources", pwd))
	e.SetConfigName(".config")

	e.AutomaticEnv()
	// enable viper to handle env values for nested structs
	e.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := e.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// 	do nothing
		default:
			return errors.Wrap(err, "error reading config file")
		}
	}

	if err := e.Unmarshal(c); err != nil {
		return errors.WithMessage(err, "error unmarshalling config")
	}
	return nil
}
