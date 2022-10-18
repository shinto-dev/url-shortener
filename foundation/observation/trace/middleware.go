package trace

import (
	"net/http"

	"github.com/google/uuid"
)

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = WithTraceID(ctx, getTraceID(r))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTraceID(r *http.Request) string {
	return firstNonEmpty(r.Header.Get("X-Request-ID"), uuid.NewString())
}

func firstNonEmpty(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}
