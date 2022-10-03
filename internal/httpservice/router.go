package httpservice

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shinto-dev/url-shortener/foundation/web"
	"github.com/shinto-dev/url-shortener/internal/httpservice/appcontext"
	handlers2 "github.com/shinto-dev/url-shortener/internal/httpservice/handlers"
)

func API(appCtx appcontext.AppContext) *mux.Router {
	r := mux.NewRouter()
	r.Use(web.PanicHandlerMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", handlers2.Health()).Methods(http.MethodGet)

	r.HandleFunc("/v1/short-url", handlers2.HandleShortURLCreate(appCtx.ShortURLCore)).
		Methods(http.MethodPost)
	r.HandleFunc("/{short_url_token}", handlers2.HandleRedirectURL(appCtx.ShortURLCore)).
		Methods(http.MethodGet)
	return r
}
