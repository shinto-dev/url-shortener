package httpservice

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shinto-dev/url-shortener/foundation/web"
	"github.com/shinto-dev/url-shortener/internal/httpservice/appcontext"
	"github.com/shinto-dev/url-shortener/internal/httpservice/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func API(appCtx appcontext.AppContext) http.Handler {
	r := mux.NewRouter()
	r.Use(web.PanicHandlerMiddleware)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/healthz", handlers.Health()).Methods(http.MethodGet)

	r.HandleFunc("/v1/short-url", handlers.HandleShortURLCreate(appCtx.ShortURLCore)).
		Methods(http.MethodPost)
	r.HandleFunc("/{short_url_token}", handlers.HandleRedirectURL(appCtx.ShortURLCore)).
		Methods(http.MethodGet)

	return otelhttp.NewHandler(r, "request")
}
