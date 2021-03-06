package httphandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"url-shortner/internal/httphandler"
	"url-shortner/internal/shorturl"
	"url-shortner/platform/apperror"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {

	errorResponse := func(code, message string) map[string]interface{} {
		return map[string]interface{}{
			"code":    code,
			"message": message,
		}
	}

	successResponse := func(expectedShortURLToken string) map[string]interface{} {
		return map[string]interface{}{
			"short_url_token": expectedShortURLToken,
		}
	}

	t.Run("should return bad request if request body is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/short-url", strings.NewReader(""))
		w := httptest.NewRecorder()

		handler := httphandler.Create(nil)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertJsonEquals(t, errorResponse("ERR-101", "request body is invalid"), w.Body.String())
	})

	t.Run("should return bad request if request body does not contain original_url", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/short-url", strings.NewReader("{}"))
		w := httptest.NewRecorder()

		handler := httphandler.Create(nil)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertJsonEquals(t, errorResponse("ERR-101", "original url is empty"), w.Body.String())
	})

	t.Run("should return short id if request body contain original_url", func(t *testing.T) {
		const originalURL = "https://github.com/shinto-dev/linear-programming"
		const expectedShortURLToken = "abcde"

		req := newHttpTestRequest(t, http.MethodPost, map[string]interface{}{
			"original_url": originalURL,
		})
		w := httptest.NewRecorder()

		mockShortURLService := shorturl.NewMockService(t)
		mockShortURLCreate(mockShortURLService, originalURL, "", expectedShortURLToken)
		handler := httphandler.Create(mockShortURLService)
		handler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assertJsonEquals(t, successResponse(expectedShortURLToken), w.Body.String())
		mockShortURLService.AssertExpectations(t)
	})

	t.Run("should return custom alias if request body contain original_url and custom alias", func(t *testing.T) {
		const originalURL = "https://github.com/shinto-dev/linear-programming"
		const customAlias = "abcde"

		req := newHttpTestRequest(t, http.MethodPost, map[string]interface{}{
			"original_url": originalURL,
			"custom_alias": customAlias,
		})
		w := httptest.NewRecorder()

		mockShortURLService := shorturl.NewMockService(t)
		mockShortURLCreate(mockShortURLService, originalURL, customAlias, customAlias)
		handler := httphandler.Create(mockShortURLService)
		handler(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assertJsonEquals(t, successResponse(customAlias), w.Body.String())
		mockShortURLService.AssertExpectations(t)
	})

	t.Run("should return error if custom alias already exists", func(t *testing.T) {
		const originalURL = "https://github.com/shinto-dev/linear-programming"
		const customAlias = "abcde"

		req := newHttpTestRequest(t, http.MethodPost, map[string]interface{}{
			"original_url": originalURL,
			"custom_alias": customAlias,
		})
		w := httptest.NewRecorder()

		mockShortURLService := shorturl.NewMockService(t)
		mockShortURLService.On("Create", mock.Anything).
			Return(shorturl.ShortURL{}, apperror.NewError(shorturl.ErrCustomURLAlreadyExists, "custom url already exists"))
		handler := httphandler.Create(mockShortURLService)
		handler(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assertJsonEquals(t, errorResponse("ERR-201", "custom url already exists"), w.Body.String())
		mockShortURLService.AssertExpectations(t)
	})

	t.Run("should return 500 error if unexpected error occurs while creating short url", func(t *testing.T) {
		const originalURL = "https://github.com/shinto-dev/linear-programming"
		req := newHttpTestRequest(t, http.MethodPost, map[string]interface{}{
			"original_url": originalURL,
		})
		w := httptest.NewRecorder()

		mockShortURLService := shorturl.NewMockService(t)
		mockShortURLService.On("Create", mock.Anything).
			Return(shorturl.ShortURL{}, errors.New("unexpected error"))
		handler := httphandler.Create(mockShortURLService)
		handler(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assertJsonEquals(t, errorResponse("", "internal error"), w.Body.String())
		mockShortURLService.AssertExpectations(t)
	})

}

func mockShortURLCreate(mockShortURLService *shorturl.MockService, originalURL, customAlias, expectedShortURLToken string) {
	mockShortURLService.On("Create", shorturl.CreateRequest{
		OriginalURL: originalURL,
		CustomAlias: customAlias,
	}).Return(shorturl.ShortURL{
		ShortURLToken: expectedShortURLToken,
	}, nil)
}

func newHttpTestRequest(t *testing.T, httpMethod string, req map[string]interface{}) *http.Request {
	requestBytes, err := json.Marshal(req)
	assert.NoError(t, err, "invalid expected argument")

	return httptest.NewRequest(httpMethod, "/", bytes.NewReader(requestBytes))
}

func assertJsonEquals(t *testing.T, expected map[string]interface{}, actual string) {
	responseBytes, err := json.Marshal(expected)
	assert.NoError(t, err, "invalid expected argument")

	assert.JSONEq(t, string(responseBytes), actual)
}
