package handlers

import (
	"net/http"
	"url-shortener/platform/web"
)

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = web.JSON(w, http.StatusOK, map[string]string{
			"message": "Hello World",
		})
	}
}
