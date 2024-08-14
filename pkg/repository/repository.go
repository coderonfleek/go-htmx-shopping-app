package repository

import (
	"database/sql"
)

type Repository struct {
	Product *ProductRepository
	Order   *OrderRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Product: NewProductRepository(db),
		Order:   NewOrderRepository(db),
	}
}
