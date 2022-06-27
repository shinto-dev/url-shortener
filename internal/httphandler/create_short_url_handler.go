package httphandler

import (
	"encoding/json"
	"net/http"
	"url-shortner/internal/shorturl"
	"url-shortner/platform/apperror"
	"url-shortner/platform/web"
)

func Create(shortURLService shorturl.Service) http.HandlerFunc {
	type CreateShortURLRequest struct {
		OriginalURL string `json:"original_url"`
		CustomAlias string `json:"custom_alias,omitempty"`
	}

	type CreateShortURLResponse struct {
		ShortURLToken string `json:"short_url_token"`
	}

	parseCreateShortURLRequest := func(w http.ResponseWriter, r *http.Request) (CreateShortURLRequest, error) {
		var req CreateShortURLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return CreateShortURLRequest{}, apperror.NewError(ErrInvalidInput, "request body is invalid")
		}

		if req.OriginalURL == "" {
			return CreateShortURLRequest{}, apperror.NewError(ErrInvalidInput, "original url is empty")
		}

		return req, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseCreateShortURLRequest(w, r)
		if err != nil {
			web.HandleError(w, err)
			return
		}

		shortURL, err := shortURLService.Create(shorturl.CreateRequest{
			OriginalURL: req.OriginalURL,
			CustomAlias: req.CustomAlias,
		})
		if err != nil {
			web.HandleError(w, err)
			return
		}

		_ = web.JSON(w, http.StatusCreated, CreateShortURLResponse{
			ShortURLToken: shortURL.ShortURLToken,
		})
	}
}
