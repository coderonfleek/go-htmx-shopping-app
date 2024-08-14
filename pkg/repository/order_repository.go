package repository

import (
	"database/sql"
	"time"

	"shopping-app/pkg/models"

	"github.com/google/uuid"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	query := `INSERT INTO orders (order_id, session_id, order_status, order_date) 
              VALUES (?, ?, ?, ?)`

	order.OrderID = uuid.New()
	order.OrderDate = time.Now()

	_, err := r.DB.Exec(query,
		order.OrderID,
		order.SessionID,
		order.OrderStatus,
		order.OrderDate,
	)
	return err
}

func (r *OrderRepository) AddOrderItem(orderItem *models.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity) 
              VALUES (?, ?, ?)`

	_, err := r.DB.Exec(query,
		orderItem.OrderID,
		orderItem.ProductID,
		orderItem.Quantity,
	)
	return err
}

func (r *OrderRepository) GetOrderWithProducts(orderID uuid.UUID) (*models.Order, error) {
	// First, get the order details
	orderQuery := `SELECT order_id, session_id, order_status, order_date 
                   FROM orders WHERE order_id = ?`

	var order models.Order
	err := r.DB.QueryRow(orderQuery, orderID).Scan(
		&order.OrderID,
		&order.SessionID,
		&order.OrderStatus,
		&order.OrderDate,
	)
	if err != nil {
		return nil, err
	}

	// Then, get all order items with their corresponding products
	itemsQuery := `
        SELECT oi.product_id, oi.quantity,
               p.product_name, p.price, p.description, p.product_image, p.date_created, p.date_modified
        FROM order_items oi
        JOIN products p ON oi.product_id = p.product_id
        WHERE oi.order_id = ?
    `
	rows, err := r.DB.Query(itemsQuery, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ProductID,
			&item.Quantity,
			&item.Product.ProductName,
			&item.Product.Price,
			&item.Product.Description,
			&item.Product.ProductImage,
			&item.Product.DateCreated,
			&item.Product.DateModified,
		)
		if err != nil {
			return nil, err
		}
		item.OrderID = orderID
		item.Product.ProductID = item.ProductID
		order.Items = append(order.Items, item)
	}

	return &order, nil
}
