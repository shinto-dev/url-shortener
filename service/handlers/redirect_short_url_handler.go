package handlers

import (
	"net/http"
	"url-shortener/business/shorturl"
	"url-shortener/platform/apperror"
	"url-shortener/platform/observation"
	"url-shortener/platform/observation/logging"
	"url-shortener/platform/web"

	"github.com/gorilla/mux"
)

func HandleRedirectURL(shortURLService shorturl.Core) http.HandlerFunc {
	return web.HandleRequest("redirect_short_url", func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		shortPath := mux.Vars(r)["short_url_token"]
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
