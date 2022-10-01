package repo

import (
	"context"
	"time"
	"url-shortener/platform/data"
	"url-shortener/platform/observation/apm"

	"gorm.io/gorm"
)

type ShortURL struct {
	ID          int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	ShortPath   string
	OriginalURL string
}

type Repo struct {
	*data.Repository[ShortURL, int64]
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db, Repository: data.NewRepository[ShortURL, int64](db)}
}

func (r *Repo) Create(ctx context.Context, shortURLRecord *ShortURL) error {
	segment := apm.StartDataStoreSegment(ctx, "shorturl-create")
	defer segment.ObserveDuration()

	return r.Repository.Create(ctx, shortURLRecord)
}

func (r *Repo) GetByShortPath(ctx context.Context, shortPath string) (ShortURL, error) {
	segment := apm.StartDataStoreSegment(ctx, "shorturl-get-by-short-path")
	defer segment.ObserveDuration()

	var shortURL ShortURL
	err := r.db.WithContext(ctx).
		First(&shortURL, "short_path = ?", shortPath).Error
	if err != nil {
		return r.handleErr(err)
	}

	return shortURL, nil
}

func (r *Repo) handleErr(err error) (ShortURL, error) {
	if err == gorm.ErrRecordNotFound {
		return ShortURL{}, nil
	}

	return ShortURL{}, err
}
