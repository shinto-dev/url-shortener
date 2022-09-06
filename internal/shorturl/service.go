package shorturl

import (
	"context"
)

type CreateRequest struct {
	OriginalURL string
	CustomAlias string
}

type ShortURL struct {
	OriginalURL   string
	ShortURLToken string
}

//go:generate mockery --name Service --inpackage --case underscore
type Service interface {
	Create(ctx context.Context, req CreateRequest) (ShortURL, error)
}
