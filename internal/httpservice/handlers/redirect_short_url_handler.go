package handlers

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/apperror"
	"github.com/shinto-dev/url-shortener/foundation/observation"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/foundation/web"
	shorturl2 "github.com/shinto-dev/url-shortener/internal/core/shorturl"
)

func HandleRedirectURL(shortURLService shorturl2.Core) http.HandlerFunc {
	return web.HandleRequest("redirect_short_url",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			shortPath := web.Param(r, "short_url_token")
			observation.Add(ctx, logging.LField("short_path", shortPath))

			shortURL, err := shortURLService.Get(ctx, shortPath)
			if err != nil {
				if apperror.Is(err, shorturl2.ErrCodeShortURLNotFound) {
					return apperror.NewError(web.ErrCodeRecordNotFound, err.Error())
				}

				return err
			}

			observation.Add(ctx, logging.LField("original_url", shortURL.OriginalURL))
			http.Redirect(w, r, shortURL.OriginalURL, http.StatusPermanentRedirect)
			return nil
		})
}
