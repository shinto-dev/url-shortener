package handlers

import (
	"context"
	"net/http"

	"github.com/shinto-dev/url-shortener/foundation/apperror"
	"github.com/shinto-dev/url-shortener/foundation/observation/logging"
	"github.com/shinto-dev/url-shortener/foundation/web"
	"github.com/shinto-dev/url-shortener/internal/core/shorturl"
)

func HandleRedirectURL(shortURLService shorturl.Core) http.HandlerFunc {
	errHandler := commonErrHandler(map[apperror.Code]*web.RequestError{
		shorturl.ErrCodeShortURLNotFound: ErrShortURLNotFound,
	})

	return web.HandleRequest("redirect_short_url",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			shortPath := web.Param(r, "short_url_token")
			logging.Add(ctx, logging.LField("short_path", shortPath))

			shortURL, err := shortURLService.Get(ctx, shortPath)
			if err != nil {
				return errHandler(ctx, err)
			}

			logging.Add(ctx, logging.LField("original_url", shortURL.OriginalURL))
			http.Redirect(w, r, shortURL.OriginalURL, http.StatusPermanentRedirect)
			return nil
		})
}
