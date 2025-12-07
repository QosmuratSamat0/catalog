package repository

import "errors"

var (
	ErrProductExists   = errors.New("Product already exists")
	ErrProductNotFound = errors.New("Product not found")
)
