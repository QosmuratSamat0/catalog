package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/QosmuratSamat0/catalog/internal/domain/models"
	"github.com/QosmuratSamat0/catalog/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(storagePath string) (*Storage, error) {
	const op = "repository.postgresql.New"

	db, err := pgxpool.New(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateProduct(ctx context.Context, name string, description string, category string, price float64) (int64, error) {
	const op = "repository.postgresql.CreateProduct"
	var id int64
	err := s.db.QueryRow(ctx,
		"INSERT INTO products(name, description, category, price) VALUES ($1, $2, $3, $4) RETURNING id",
		name, description, category, price,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return 0, fmt.Errorf("%s: уникальное ограничение нарушено: %w", op, err)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetProductById(ctx context.Context, id int64) (models.Product, error) {
	const op = "repository.postgresql.GetProductById"

	var product models.Product

	err := s.db.QueryRow(
		ctx, "SELECT id, name, description, category, price, created_at FROM products WHERE id = $1", id,
	).Scan(&product.ID, &product.Name, &product.Description, &product.Category, &product.Price, &product.CreatedAt)
	if err != nil {
		// not found
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("%s: product not found", op)
		}

		// pg error
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return models.Product{}, fmt.Errorf("%s: duplicate constraint: %w", op, err)
			}
		}

		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}
	return product, nil
}

func (s *Storage) SearchProducts(ctx context.Context, f repository.SearchFilter) ([]models.Product, error) {
	const op = "repository.postgresql.SearchProducts"

	query := `
        SELECT id, name, description, category, price, created_at
        FROM products
        WHERE category ILIKE $1
          AND price >= $2
          AND price <= $3
    `

	if f.SortBy == "price_asc" {
		query += " ORDER BY price ASC"
	}
	if f.SortBy == "price_desc" {
		query += " ORDER BY price DESC"
	}
	if f.SortBy == "date_desc" {
		query += " ORDER BY created_at DESC"
	}

	rows, err := s.db.Query(ctx, query, f.Category+"%", f.PriceMin, f.PriceMax)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.Price,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}
