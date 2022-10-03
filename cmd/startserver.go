package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/internal/config"
	"github.com/shinto-dev/url-shortener/internal/httpservice"
	"github.com/shinto-dev/url-shortener/internal/httpservice/appcontext"
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

		router := httpservice.API(appCtx)
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HTTPServer.Port),
			Handler: router,
		}

		serverErrors := make(chan error, 1)
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			logging.WithFields(
				zap.Int("port", conf.HTTPServer.Port),
			).Info("starting HTTP server")
			serverErrors <- server.ListenAndServe()
		}()

		select {
		case err := <-serverErrors:
			logging.WithFields(zap.Error(err)).
				Fatal("error while running HTTP server")
		case <-sigquit:
			logging.WithFields(zap.Any("signal", sigquit)).
				Info("gracefully shutting down the server")

			defer logging.WithFields(zap.Any("signal", sigquit)).
				Info("shutdown completed")

			ctx, cancel := context.WithTimeout(context.Background(), conf.HTTPServer.ShutdownTimeout)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logging.WithFields(zap.Error(err)).
					Error("unable to shutdown the server")
				server.Close()
				return
			}
		}
	}
}
