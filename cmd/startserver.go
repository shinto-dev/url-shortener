package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shinto-dev/url-shortener/config"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/service"
	"github.com/shinto-dev/url-shortener/service/appcontext"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func newStartServerCommand(conf config.Config) *cobra.Command {
	return &cobra.Command{
		Use:     "startserver",
		Short:   "Start HTTP API server",
		Aliases: []string{"startapp", "runserver"},
		Run:     runHTTServerFunc(conf),
	}
}

func runHTTServerFunc(conf config.Config) func(_ *cobra.Command, _ []string) {
	return func(_ *cobra.Command, _ []string) {
		appCtx, err := appcontext.Get(conf)
		if err != nil {
			logging.WithFields(zap.Error(err)).
				Fatal("error while creating app context")
		}

		router := service.API(appCtx)
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HTTPServer.Port),
			Handler: router,
		}

		go func() {
			logging.WithFields(
				zap.Int("port", conf.HTTPServer.Port),
			).Info("starting HTTP server")
			server.ListenAndServe()
		}()

		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

		<-sigquit
		logging.WithFields().Info("gracefully shutting down the server")

		if err := server.Shutdown(context.Background()); err != nil {
			logging.WithFields(zap.Error(err)).Error("unable to shutdown the server")
			return
		}

		logging.WithFields().Info("server stopped")
	}
}
