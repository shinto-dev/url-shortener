package data

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository[T, ID any] struct {
	db *gorm.DB
}

func NewRepository[T, ID any](db *gorm.DB) *Repository[T, ID] {
	return &Repository[T, ID]{db: db}
}

func (r *Repository[T, ID]) Create(ctx context.Context, t *T) error {
	result := r.db.WithContext(ctx).Create(t)

	return result.Error
}

func (r *Repository[T, ID]) GetByID(ctx context.Context, id ID) (T, error) {
	var item T
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return item, nil
	}

	return item, err
}

func (r *Repository[T, ID]) GetByFieldName(ctx context.Context, fieldName string, value any) (T, error) {
	var item T
	err := r.db.WithContext(ctx).First(&item, fmt.Sprintf("%s=?", fieldName), value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return item, nil
	}

	return item, err
}

func (r *Repository[T, ID]) Delete(ctx context.Context, t *T) error {
	return r.db.WithContext(ctx).Delete(t).Error
}
