package handlers

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/apperror"
	"github.com/shinto-dev/url-shortener/foundation/observation"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/foundation/web"
	"github.com/shinto-dev/url-shortener/internal/business/shorturl"

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

	return web.HandleRequest("create_short_url",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var req CreateShortURLRequest
			if err := web.Decode(r, &req); err != nil {
				return apperror.NewError(ErrInvalidInput, "request body is invalid")
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
