package shorturl

import (
	"github.com/shinto-dev/url-shortener/foundation/apperror"
)

const (
	ErrCodeShortURLNotFound = "short_url_not_found"
)

var ErrShortURLNotFound = apperror.NewError(ErrCodeShortURLNotFound, "short url not found")
