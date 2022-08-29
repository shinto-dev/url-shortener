package handlers

import (
	"net/http"
	"url-shortner/internal/shorturl"
)

func RedirectURL(shortURLService shorturl.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
