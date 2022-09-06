package web

import (
	"net/http"
	"url-shortner/platform/observation"
	"url-shortner/platform/observation/apm"
	"url-shortner/platform/observation/logging"

	"github.com/prometheus/client_golang/prometheus"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func HandleRequest(apiContext string, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithLogger(r.Context())
		ctx = apm.WithAPM(ctx, "create_short_url")
		observation.Add(ctx, logging.LField("context", apiContext))

		hist := apm.FromContext(ctx)
		timer := prometheus.NewTimer(hist.WithLabelValues("service", "handler"))
		defer timer.ObserveDuration()

		if err := fn(w, r.WithContext(ctx)); err != nil {
			HandleError(ctx, w, err)
		}
	}
}
