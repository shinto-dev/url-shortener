package main

import (
	"url-shortner/cmd"
	"url-shortner/platform/observation/logging"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		logging.WithFields(zap.Error(err)).
			Fatal("error while starting the app")
	}
}

func run() error {
	cli := cmd.NewCLI()
	err := cli.Execute()
	if err != nil {
		return errors.Wrap(err, "error initializing the command")
	}
	return nil
}
