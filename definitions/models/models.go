package models

import (
	"context"

	"github.com/Nemutagk/godb/v2/definitions/models"
)

type Model[T any] interface {
	Get(ctx context.Context, filters models.GroupFilter, options *models.Options) ([]T, error)
	GetOne(ctx context.Context, filters models.GroupFilter) (T, error)
	Create(ctx context.Context, data map[string]any) (T, error)
	Update(ctx context.Context, filters models.GroupFilter, data map[string]any) (T, error)
	Delete(ctx context.Context, filters models.GroupFilter) error
	Count(ctx context.Context, filters models.GroupFilter) (int64, error)
}
