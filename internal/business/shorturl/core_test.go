package shorturl_test

import (
	"context"
	"testing"

	"github.com/shinto-dev/url-shortener/foundation/apperror"
	"github.com/shinto-dev/url-shortener/foundation/observation/apm"
	shorturl2 "github.com/shinto-dev/url-shortener/internal/business/shorturl"
	"github.com/shinto-dev/url-shortener/internal/business/shorturl/repo"
	test2 "github.com/shinto-dev/url-shortener/internal/business/test"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db := test2.ConnectTestDB(t)
	core := shorturl2.NewShortURLCore(db)

	ctx := context.Background()
	ctx = apm.WithAPM(ctx, "test") //fixme our tests should not worry about apm related configs

	t.Run("should create a short url", func(t *testing.T) {
		const originalURL = "https://www.google.com"

		testCtx := test2.NewTestCtx(db)
		defer testCtx.DeleteShortURLByOriginalURL(t, originalURL)

		req := shorturl2.CreateRequest{
			OriginalURL: originalURL,
		}

		shortURL, err := core.Create(ctx, req)
		assert.NoError(t, err)

		assert.NotEmpty(t, shortURL.ShortURLPath)

		shortURLRecord := testCtx.GetShortURLByShortPath(t, shortURL.ShortURLPath)
		assert.Equal(t, req.OriginalURL, shortURLRecord.OriginalURL)
		assert.NotEmpty(t, shortURLRecord.ShortPath)
		assert.Equal(t, shortURLRecord.ShortPath, shortURL.ShortURLPath)
	})
}

func TestGet(t *testing.T) {
	db := test2.ConnectTestDB(t)
	core := shorturl2.NewShortURLCore(db)
	testCtx := test2.NewTestCtx(db)

	ctx := context.Background()
	ctx = apm.WithAPM(ctx, "test") //fixme our tests should not worry about apm related configs

	t.Run("should get a short url", func(t *testing.T) {
		const originalURL = "https://www.google1.com"
		const shortPath = "short"

		testCtx.CreateShortURL(t, repo.ShortURL{
			OriginalURL: originalURL,
			ShortPath:   shortPath,
		})
		defer testCtx.DeleteShortURLByOriginalURL(t, originalURL)

		shortURL, err := core.Get(ctx, shortPath)
		assert.NoError(t, err)

		assert.Equal(t, shortPath, shortURL.ShortURLPath)
		assert.Equal(t, originalURL, shortURL.OriginalURL)
	})

	t.Run("should return error if short url not found", func(t *testing.T) {
		_, err := core.Get(ctx, "not-found")
		assert.Error(t, err)
		assert.True(t, apperror.Is(err, shorturl2.ErrCodeShortURLNotFound))
	})
}
