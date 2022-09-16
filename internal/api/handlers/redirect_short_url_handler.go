package handlers

import (
	"net/http"
	"url-shortener/internal/shorturl"
)

func RedirectURL(shortURLService shorturl.Core) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}
