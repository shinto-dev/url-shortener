package handlers

import (
	"context"
	"net/http"
	"url-shortener/business/shorturl"
	"url-shortener/foundation/apperror"
	"url-shortener/foundation/observation"
	"url-shortener/foundation/observation/logging"
	"url-shortener/foundation/web"
)

func HandleRedirectURL(shortURLService shorturl.Core) http.HandlerFunc {
	return web.HandleRequest("redirect_short_url",
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			shortPath := web.Param(r, "short_url_token")
			observation.Add(ctx, logging.LField("short_path", shortPath))

			shortURL, err := shortURLService.Get(ctx, shortPath)
			if err != nil {
				if apperror.Is(err, shorturl.ErrCodeShortURLNotFound) {
					return apperror.NewError(web.ErrCodeRecordNotFound, err.Error())
				}

				return err
			}

			http.Redirect(w, r, shortURL.ShortURLPath, http.StatusPermanentRedirect)
			return nil
		})
}
