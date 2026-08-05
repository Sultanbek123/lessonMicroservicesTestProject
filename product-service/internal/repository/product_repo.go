package repository

import (
	"context"
	"database/sql"
	"product-service/internal/dto"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (int, error) {
	var id int
	query := `INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, req.Name, req.Price).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	query := `SELECT id, name, price FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []dto.ProductResponse
	for rows.Next() {
		var p dto.ProductResponse
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	if products == nil {
		products = []dto.ProductResponse{}
	}

	return products, nil
}
