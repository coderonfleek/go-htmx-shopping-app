package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductID    uuid.UUID
	ProductName  string
	Price        float64
	Description  string
	ProductImage string
	DateCreated  time.Time
	DateModified time.Time
}
