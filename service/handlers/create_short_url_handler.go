package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortener/business/shorturl"
	"url-shortener/foundation/apperror"
	"url-shortener/foundation/observation"
	"url-shortener/foundation/observation/logging"
	"url-shortener/foundation/web"

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

func HandleShortURLCreate(shortURLService shorturl.Core) http.HandlerFunc {
	type CreateShortURLResponse struct {
		ShortURLToken string `json:"short_url_token"`
	}

	parseCreateShortURLRequest := func(r *http.Request) (CreateShortURLRequest, error) {
		var req CreateShortURLRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			return CreateShortURLRequest{}, apperror.NewError(ErrInvalidInput, "request body is invalid")
		}

		return req, nil
	}

	return web.HandleRequest("create_short_url",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			req, err := parseCreateShortURLRequest(r)
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
