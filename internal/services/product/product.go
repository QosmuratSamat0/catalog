package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/QosmuratSamat0/catalog/internal/cache"
	"github.com/QosmuratSamat0/catalog/internal/domain/models"
	"github.com/QosmuratSamat0/catalog/internal/pkg/logger/sl"
	"github.com/QosmuratSamat0/catalog/internal/repository"
)

type Service struct {
	log   *slog.Logger
	repo  repository.Repository
	cache cache.Cache
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrProductExists      = errors.New("product already exists")
)

type Cache interface {
	Set(ctx context.Context, key string, value string, tll time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

func New(repo repository.Repository, cache *cache.Cache) *Service {
	return &Service{repo: repo, cache: *cache}
}

func (s *Service) CreateProduct(ctx context.Context, name, description, category string, price float64) (int64, error) {
	const op = "services.product.CreateProduct"

	log := slog.With("op", op)

	log.Info("creating product")

	id, err := s.repo.CreateProduct(ctx, name, description, category, price)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			s.log.Warn("Product not found", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, err)
		}
		s.log.Error("failed to create product", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	p := models.Product{
		ID:        id,
		Name:      name,
		Category:  category,
		Price:     price,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(p)
	if err != nil {
		return id, nil
	}

	key := fmt.Sprintf("product:%d", id)

	err = s.cache.Set(ctx, key, string(data), 10*time.Minute)
	if err != nil {
		s.log.Error("failed to cache product", sl.Err(err))
	}

	return id, nil
}

func (s *Service) GetProductById(ctx context.Context, id int64) (product models.Product, err error) {
	const op = "services.product.GetProductById"

	log := slog.With(
		"op", op,
		"id", id,
	)

	log.Info("getting product")

	key := fmt.Sprintf("product:%d", id)

	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var p models.Product
		if err := json.Unmarshal([]byte(cached), &p); err != nil {
			return p, nil
		}
	}

	product, err = s.repo.GetProductById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrProductNotFound) {
			s.log.Warn("Product not found", sl.Err(err))

			return models.Product{}, fmt.Errorf("%s: %w", op, err)
		}
		s.log.Error("failed to get product", sl.Err(err))

		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	data, _ := json.Marshal(product)
	s.cache.Set(ctx, key, string(data), 10*time.Minute)
	return product, nil
}

func (s *Service) SearchProducts(ctx context.Context, category string, priceMin, priceMax float64, sortBy string) (products []models.Product, err error) {
	const op = "services.product.SearchProducts"

	log := slog.With("op", op)

	log.Info("searching product")

	key := fmt.Sprintf("search:%s:%f:%f:%s", category, priceMin, priceMax, sortBy)

	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var products []models.Product
		if err := json.Unmarshal([]byte(cached), &products); err != nil {
			return products, fmt.Errorf("%s: %w", op, err)
		}
	}

	filter := repository.SearchFilter{
		Category: category,
		PriceMin: priceMin,
		PriceMax: priceMax,
		SortBy:   sortBy,
	}

	products, err = s.repo.SearchProducts(ctx, filter)
	if err != nil {
		log.Error("failed to get products", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	data, err := json.Marshal(products)
	if err == nil {
		_ = s.cache.Set(ctx, key, string(data), time.Minute*5)
	}

	return products, nil

}
