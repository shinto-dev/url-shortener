package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortner/internal/api"

	"github.com/sirupsen/logrus"
)

func main() {
	router := api.API()
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}

	go func() {
		fmt.Printf("starting HTTP server, listening at %d\n", 8080)
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatal("failed to start the server")
		}

	}()

	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	_ = <-sigquit
	logrus.Info("gracefully shutting down the server")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Error("unable to shutdown the server")
		return
	}
}
