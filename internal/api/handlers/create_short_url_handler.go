package handlers

import (
	"encoding/json"
	"net/http"
	"url-shortner/internal/shorturl"
	"url-shortner/platform/apperror"
	"url-shortner/platform/observation/apm"
	"url-shortner/platform/observation/logging"
	"url-shortner/platform/web"

	"github.com/prometheus/client_golang/prometheus"
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

	hist := apm.NewHistogram("create_short_url")

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithLogger(r.Context())
		ctx = apm.WithValue(ctx, hist)

		timer := prometheus.NewTimer(hist.WithLabelValues("service", "handler"))
		defer timer.ObserveDuration()

		req, err := parseCreateShortURLRequest(w, r)
		if err != nil {
			web.HandleError(ctx, w, err)
			return
		}

		shortURL, err := shortURLService.Create(shorturl.CreateRequest{
			OriginalURL: req.OriginalURL,
			CustomAlias: req.CustomAlias,
		})
		if err != nil {
			web.HandleError(ctx, w, err)
			return
		}

		_ = web.JSON(w, http.StatusCreated, CreateShortURLResponse{
			ShortURLToken: shortURL.ShortURLToken,
		})
	}
}
