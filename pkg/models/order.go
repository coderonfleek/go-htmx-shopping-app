package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderID     uuid.UUID
	SessionID   string
	OrderStatus string
	OrderDate   time.Time
	Items       []OrderItem
}
