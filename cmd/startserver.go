package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/internal/config"
	"github.com/shinto-dev/url-shortener/internal/httpservice"
	"github.com/shinto-dev/url-shortener/internal/httpservice/appcontext"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
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
		exporter, err := jaeger.New(
			jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.Observation.JaegerEndpoint)),
		)

		probability := 1.0
		serviceName := "url-shortener"

		traceProvider := trace.NewTracerProvider(
			trace.WithSampler(trace.TraceIDRatioBased(probability)),
			trace.WithBatcher(exporter,
				trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
				trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
				trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			),
			trace.WithResource(
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceNameKey.String(serviceName),
					attribute.String("exporter", "jaeger"),
				),
			),
		)

		otel.SetTracerProvider(traceProvider)

		appCtx, err := appcontext.Get(conf)
		if err != nil {
			logging.WithFields(zap.Error(err)).
				Fatal("error while creating app context")
		}

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", conf.HTTPServer.Port),
			Handler: httpservice.API(appCtx),
		}

		serverErrors := make(chan error, 1)
		sigquit := make(chan os.Signal, 1)
		signal.Notify(sigquit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			logging.WithFields(
				zap.Int("port", conf.HTTPServer.Port),
				zap.String("jaeger_endpoint", conf.Observation.JaegerEndpoint),
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
