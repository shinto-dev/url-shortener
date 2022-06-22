package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"url-shortner/platform/web"
)

func API() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		_ = web.JSON(w, map[string]string{
			"message": "Hello World",
		}, http.StatusOK)
	}).Methods(http.MethodGet)

	return r
}
