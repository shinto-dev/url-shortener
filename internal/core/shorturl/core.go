package shorturl

import (
	"context"

	"github.com/shinto-dev/url-shortener/internal/core/shorturl/repo"
	"github.com/shinto-dev/url-shortener/internal/core/shorturl/shortid"

	"gorm.io/gorm"
)

type CreateRequest struct {
	OriginalURL string
}

type ShortURL struct {
	OriginalURL  string
	ShortURLPath string
}

//go:generate mockery --name Core --inpackage --case underscore
type Core interface {
	Create(ctx context.Context, req CreateRequest) (ShortURL, error)
	Get(ctx context.Context, shortURLPath string) (ShortURL, error)
}

type coreImpl struct {
	repo             *repo.Repo
	shortIDGenerator *shortid.IDGenerator
}

func NewShortURLCore(db *gorm.DB) Core {
	return &coreImpl{
		repo:             repo.NewRepo(db),
		shortIDGenerator: shortid.NewShortIDGenerator(db),
	}
}

func (s *coreImpl) Create(ctx context.Context, req CreateRequest) (ShortURL, error) {
	shortID, err := s.shortIDGenerator.NewBase58ID()
	if err != nil {
		return ShortURL{}, err
	}

	shortURLRecord := repo.ShortURL{
		ShortPath:   shortID,
		OriginalURL: req.OriginalURL,
	}

	if err = s.repo.Create(ctx, &shortURLRecord); err != nil {
		return ShortURL{}, err
	}

	return ShortURL{
		OriginalURL:  shortURLRecord.OriginalURL,
		ShortURLPath: shortID,
	}, nil
}

func (s *coreImpl) Get(ctx context.Context, shortURLPath string) (ShortURL, error) {
	shortURLRecord, err := s.repo.GetByShortPath(ctx, shortURLPath)
	if err != nil {
		return ShortURL{}, err
	}

	if shortURLRecord == (repo.ShortURL{}) {
		return ShortURL{}, ErrShortURLNotFound
	}

	return ShortURL{
		OriginalURL:  shortURLRecord.OriginalURL,
		ShortURLPath: shortURLRecord.ShortPath,
	}, nil
}
