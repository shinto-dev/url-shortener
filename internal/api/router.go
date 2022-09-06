package api

import (
	"net/http"
	"url-shortner/internal/api/handlers"
	"url-shortner/platform/web"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func API() *mux.Router {
	r := mux.NewRouter()
	r.Use(web.PanicHandlerMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", handlers.Health()).Methods(http.MethodGet)

	r.HandleFunc("/v1/short-url", handlers.Create(nil)).
		Methods(http.MethodPost)
	return r
}
