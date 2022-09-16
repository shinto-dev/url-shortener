package test

import (
	"testing"
	"url-shortener/internal/shorturl/repo"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type TestCtx struct {
	db *gorm.DB
}

func NewTestCtx(db *gorm.DB) *TestCtx {
	return &TestCtx{db: db}
}

func (tc *TestCtx) CreateShortURL(t *testing.T, shortURL repo.ShortURL) {
	err := tc.db.Create(&shortURL).Error
	assert.NoError(t, err)
}

func (tc *TestCtx) GetShortURLByShortPath(t *testing.T, shortPath string) repo.ShortURL {
	var shortURL repo.ShortURL
	err := tc.db.First(&shortURL, "short_path=?", shortPath).Error
	assert.NoError(t, err)
	return shortURL
}

func (tc *TestCtx) DeleteShortURLByOriginalURL(t *testing.T, originalURL string) {
	err := tc.db.Exec("DELETE FROM short_urls WHERE original_url=?", originalURL).Error
	assert.NoError(t, err)
}
