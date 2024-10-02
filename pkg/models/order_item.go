package models

import (
	"github.com/google/uuid"
)

type OrderItem struct {
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Product   Product
	Cost      float64
}
