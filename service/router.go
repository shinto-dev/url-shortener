package service

import (
	"net/http"
	"url-shortener/platform/web"
	"url-shortener/service/appcontext"
	"url-shortener/service/handlers"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func API(appCtx appcontext.AppContext) *mux.Router {
	r := mux.NewRouter()
	r.Use(web.PanicHandlerMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", handlers.Health()).Methods(http.MethodGet)

	r.HandleFunc("/v1/short-url", handlers.HandleShortURLCreate(appCtx.ShortURLCore)).
		Methods(http.MethodPost)
	r.HandleFunc("/{short_url_token}", handlers.HandleRedirectURL(appCtx.ShortURLCore)).
		Methods(http.MethodGet)
	return r
}
