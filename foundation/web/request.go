package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
