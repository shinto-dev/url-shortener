package web

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/observation"
	"github.com/shinto-dev/url-shortener/foundation/observation/apm"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/foundation/observation/trace"
)

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func HandleRequest(apiContext string, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithLogger(r.Context())
		ctx = trace.WithTraceID(ctx, TraceID(r))
		ctx = apm.WithAPM(ctx, apiContext)
		observation.Add(ctx, logging.LField("context", apiContext))

		segment := apm.StartSegment(ctx, "handler")
		defer segment.ObserveDuration()

		if err := fn(ctx, w, r.WithContext(ctx)); err != nil {
			HandleError(ctx, w, err)
		}

		logging.FromContext(ctx).Info("Request completed")
	}
}
