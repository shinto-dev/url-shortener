package handlers

import (
	"encoding/json"
	"net/http"
	"url-shortner/internal/shorturl"
	"url-shortner/platform/apperror"
	"url-shortner/platform/observation"
	"url-shortner/platform/observation/logging"
	"url-shortner/platform/web"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url"`
}

func (c CreateShortURLRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.OriginalURL,
			validation.Required,
			validation.Length(5, 255),
			is.URL),
	)
}

func Create(shortURLService shorturl.Core) http.HandlerFunc {
	type CreateShortURLResponse struct {
		ShortURLToken string `json:"short_url_token"`
	}

	parseCreateShortURLRequest := func(w http.ResponseWriter, r *http.Request) (CreateShortURLRequest, error) {
		var req CreateShortURLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return CreateShortURLRequest{}, apperror.NewError(ErrInvalidInput, "request body is invalid")
		}

		return req, nil
	}

	return web.HandleRequest("create_short_url", func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		req, err := parseCreateShortURLRequest(w, r)
		if err != nil {
			return err
		}

		observation.Add(ctx, logging.LField("original_url", req.OriginalURL))
		if err := req.Validate(); err != nil {
			return apperror.NewErrorWithCause(err, ErrInvalidInput, err.Error())
		}

		shortURL, err := shortURLService.Create(ctx, shorturl.CreateRequest{
			OriginalURL: req.OriginalURL,
		})
		if err != nil {
			return err
		}

		_ = web.JSON(w, http.StatusCreated, CreateShortURLResponse{
			ShortURLToken: shortURL.ShortURLPath,
		})
		return nil
	})
}
