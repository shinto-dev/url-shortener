package web

import (
	"net/http"

	"github.com/google/uuid"
)

func TraceID(r *http.Request) string {
	return firstNonEmpty(r.Header.Get("X-Request-ID"), uuid.NewString())
}

func firstNonEmpty(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}
