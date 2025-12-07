package repository

import (
	"context"

	"github.com/QosmuratSamat0/catalog/internal/domain/models"
)

type SearchFilter struct {
	Category string
	PriceMin float64
	PriceMax float64
	SortBy   string
}

type Repository interface {
	CreateProduct(ctx context.Context, name string, description string, category string, price float64) (int64, error)
	GetProductById(ctx context.Context, id int64) (models.Product, error)
	SearchProducts(ctx context.Context, f SearchFilter) ([]models.Product, error)
}
