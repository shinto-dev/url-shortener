package web

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TraceID(r *http.Request) string {
	return firstNonEmpty(r.Header.Get("X-Request-ID"), uuid.NewString())
}

func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	//decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}

func Param(r *http.Request, key string) string {
	m := mux.Vars(r)
	return m[key]
}

func firstNonEmpty(s1, s2 string) string {
	if s1 != "" {
		return s1
	}
	return s2
}
