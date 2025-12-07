package models

import "time"

type Product struct {
	ID          int64
	Name        string
	Description string
	Category    string
	Price       float64
	CreatedAt   time.Time
}
