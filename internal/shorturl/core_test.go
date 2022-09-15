package shorturl_test

import (
	"context"
	"testing"
	"url-shortner/internal/shorturl"
	"url-shortner/internal/shorturl/repo"
	"url-shortner/internal/test"
	"url-shortner/platform/apperror"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db := test.ConnectTestDB(t)
	core := shorturl.NewShortURLCore(db)

	t.Run("should create a short url", func(t *testing.T) {
		const originalURL = "https://www.google.com"

		testCtx := test.NewTestCtx(db)
		defer testCtx.DeleteShortURLByOriginalURL(t, originalURL)

		req := shorturl.CreateRequest{
			OriginalURL: originalURL,
		}

		shortURL, err := core.Create(context.Background(), req)
		assert.NoError(t, err)

		assert.NotEmpty(t, shortURL.ShortURLPath)

		shortURLRecord := testCtx.GetShortURLByShortPath(t, shortURL.ShortURLPath)
		assert.Equal(t, req.OriginalURL, shortURLRecord.OriginalURL)
		assert.NotEmpty(t, shortURLRecord.ShortPath)
		assert.Equal(t, shortURLRecord.ShortPath, shortURL.ShortURLPath)
	})
}

func TestGet(t *testing.T) {
	db := test.ConnectTestDB(t)
	core := shorturl.NewShortURLCore(db)
	testCtx := test.NewTestCtx(db)

	t.Run("should get a short url", func(t *testing.T) {
		const originalURL = "https://www.google1.com"
		const shortPath = "short"

		testCtx.CreateShortURL(t, repo.ShortURL{
			OriginalURL: originalURL,
			ShortPath:   shortPath,
		})
		defer testCtx.DeleteShortURLByOriginalURL(t, originalURL)

		shortURL, err := core.Get(context.Background(), shortPath)
		assert.NoError(t, err)

		assert.Equal(t, shortPath, shortURL.ShortURLPath)
		assert.Equal(t, originalURL, shortURL.OriginalURL)
	})

	t.Run("should return error if short url not found", func(t *testing.T) {
		_, err := core.Get(context.Background(), "not-found")
		assert.Error(t, err)
		assert.True(t, apperror.Is(err, shorturl.ErrCodeShortURLNotFound))
	})
}
