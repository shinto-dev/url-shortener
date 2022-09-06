package handlers

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

	return web.HandleRequest("create_short_url", func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := parseCreateShortURLRequest(w, r)
		if err != nil {
			return err
		}

		shortURL, err := shortURLService.Create(ctx, shorturl.CreateRequest{
			OriginalURL: req.OriginalURL,
			CustomAlias: req.CustomAlias,
		})
		if err != nil {
			return err
		}

		_ = web.JSON(w, http.StatusCreated, CreateShortURLResponse{
			ShortURLToken: shortURL.ShortURLToken,
		})
		return nil
	})
}
