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
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/v1/short-url", handlers.Create(nil)).
		Methods(http.MethodPost)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = web.JSON(w, http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	}).Methods(http.MethodGet)

	return r
}
