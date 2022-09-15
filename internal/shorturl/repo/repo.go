package repo

import (
	"context"
	"time"
	"url-shortner/platform/data"

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

func (r *Repo) GetByShortPath(ctx context.Context, shortPath string) (ShortURL, error) {
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
