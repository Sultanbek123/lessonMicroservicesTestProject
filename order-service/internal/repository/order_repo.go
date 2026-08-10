package repository

import (
	"context"
	"database/sql"
	"order-service/internal/dto"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, userID string, req dto.CreateOrderRequest) (int, error) {
	var id int
	query := `INSERT INTO orders (user_id, product_id, quantity, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, userID, req.ProductID, req.Quantity, "created").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context, userID string) ([]dto.OrderResponse, error) {
	query := `SELECT id, user_id, product_id, quantity, status FROM orders WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []dto.OrderResponse
	for rows.Next() {
		var o dto.OrderResponse
		if err := rows.Scan(&o.ID, &o.UserID, &o.ProductID, &o.Quantity, &o.Status); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if orders == nil {
		orders = []dto.OrderResponse{}
	}

	return orders, nil
}
