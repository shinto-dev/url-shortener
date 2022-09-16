package web

import (
	"net/http"
	"url-shortener/platform/observation"
	"url-shortener/platform/observation/apm"
	"url-shortener/platform/observation/logging"
	"url-shortener/platform/observation/trace"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func HandleRequest(apiContext string, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithLogger(r.Context())
		ctx = trace.WithTraceID(ctx, firstNonEmpty(r.Header.Get("X-Request-ID"), uuid.NewString()))
		ctx = apm.WithAPM(ctx, apiContext)
		observation.Add(ctx, logging.LField("context", apiContext))

		hist := apm.FromContext(ctx)
		timer := prometheus.NewTimer(hist.WithLabelValues("service", "handler"))
		defer timer.ObserveDuration()

		if err := fn(w, r.WithContext(ctx)); err != nil {
			HandleError(ctx, w, err)
		}

		logging.FromContext(ctx).Info("Request completed")
	}
}

func firstNonEmpty(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}
