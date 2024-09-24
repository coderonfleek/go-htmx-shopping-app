package repository

import (
	"database/sql"
	"time"

	"shopping-app/pkg/models"

	"github.com/google/uuid"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) GetProductByID(productID uuid.UUID) (*models.Product, error) {
	query := `SELECT product_id, product_name, price, description, product_image, date_created, date_modified 
              FROM products WHERE product_id = ?`
	row := r.DB.QueryRow(query, productID)

	var product models.Product
	err := row.Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Price,
		&product.Description,
		&product.ProductImage,
		&product.DateCreated,
		&product.DateModified,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := `INSERT INTO products (product_id, product_name, price, description, product_image, date_created, date_modified) 
              VALUES (?, ?, ?, ?, ?, ?, ?)`

	product.ProductID = uuid.New()
	product.DateCreated = time.Now()
	product.DateModified = time.Now()

	_, err := r.DB.Exec(query,
		product.ProductID,
		product.ProductName,
		product.Price,
		product.Description,
		product.ProductImage,
		product.DateCreated,
		product.DateModified,
	)
	return err
}

func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := `UPDATE products SET product_name = ?, price = ?, description = ?, date_modified = ? 
              WHERE product_id = ?`

	product.DateModified = time.Now()

	_, err := r.DB.Exec(query,
		product.ProductName,
		product.Price,
		product.Description,
		product.DateModified,
		product.ProductID,
	)
	return err
}

func (r *ProductRepository) DeleteProduct(productID uuid.UUID) error {
	query := `DELETE FROM products WHERE product_id = ?`
	_, err := r.DB.Exec(query, productID)
	return err
}

func (r *ProductRepository) ListProducts(limit, offset int) ([]models.Product, error) {
	query := `SELECT product_id, product_name, price, description, product_image, date_created, date_modified 
              FROM products ORDER BY date_created DESC LIMIT ? OFFSET ?`

	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Price,
			&product.Description,
			&product.ProductImage,
			&product.DateCreated,
			&product.DateModified,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *ProductRepository) GetTotalProductsCount() (int, error) {
	var count int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
