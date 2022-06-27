package httphandler

import (
	"net/http"
	"url-shortner/platform/web"

	"github.com/gorilla/mux"
)

func API() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/v1/short-url", Create(nil)).
		Methods(http.MethodPost)

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = web.JSON(w, http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	}).Methods(http.MethodGet)

	return r
}
