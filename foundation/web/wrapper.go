package web

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/observation"
	"github.com/shinto-dev/url-shortener/foundation/observation/apm"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
)

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func HandleRequest(apiContext string, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := observation.WithObservation(r.Context(), observation.Config{
			Context:        apiContext,
			TraceID:        r.Header.Get("X-Trace-ID"),
			SupportLogging: true,
			SupportAPM:     true,
		})

		segment := apm.StartSegment(ctx, "handler")
		defer segment.ObserveDuration()

		if err := fn(ctx, w, r.WithContext(ctx)); err != nil {
			HandleError(ctx, w, err)
		}

		logging.FromContext(ctx).Info("Request completed")
	}
}
