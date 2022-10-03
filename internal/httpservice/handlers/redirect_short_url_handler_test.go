package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/shinto-dev/url-shortener/foundation/apperror"
	shorturl2 "github.com/shinto-dev/url-shortener/internal/core/shorturl"
	"github.com/shinto-dev/url-shortener/internal/httpservice/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedirectURL(t *testing.T) {
	t.Parallel()

	t.Run("should return not found if short url token is not found", func(t *testing.T) {
		mockShortURLCore := shorturl2.NewMockCore(t)
		mockShortURLCore.On("Get", mock.Anything, "shorturl").
			Return(shorturl2.ShortURL{}, apperror.NewError(
				shorturl2.ErrCodeShortURLNotFound, "short url not found"))

		req := httptest.NewRequest(http.MethodGet, "/shorturl", strings.NewReader(""))
		w := httptest.NewRecorder()

		handleRequest(mockShortURLCore, w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("should return internal server error if unexpected error occurs", func(t *testing.T) {
		mockShortURLCore := shorturl2.NewMockCore(t)
		mockShortURLCore.On("Get", mock.Anything, "shorturl").
			Return(shorturl2.ShortURL{}, errors.New("some error"))

		req := httptest.NewRequest(http.MethodGet, "/shorturl", strings.NewReader(""))
		w := httptest.NewRecorder()

		handleRequest(mockShortURLCore, w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("should return redirect if short url token is found", func(t *testing.T) {
		mockShortURLCore := shorturl2.NewMockCore(t)
		mockShortURLCore.On("Get", mock.Anything, "shorturl").
			Return(shorturl2.ShortURL{ShortURLPath: "https://google.com"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/shorturl", strings.NewReader(""))
		w := httptest.NewRecorder()

		handleRequest(mockShortURLCore, w, req)

		assert.Equal(t, http.StatusPermanentRedirect, w.Code)
		assert.Equal(t, "https://google.com", w.Header().Get("Location"))
	})
}

func handleRequest(mockShortURLCore *shorturl2.MockCore, w *httptest.ResponseRecorder, req *http.Request) {
	r := mux.NewRouter()
	r.HandleFunc("/{short_url_token}", handlers.HandleRedirectURL(mockShortURLCore))
	r.ServeHTTP(w, req)
}
