package repository

import (
	"context"
	"database/sql"
	"notification-service/internal/dto"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(ctx context.Context, userID string, orderID int, message string) (int, error) {
	var id int
	query := `INSERT INTO notifications (user_id, order_id, message, status) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, userID, orderID, message, "created").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *NotificationRepository) GetNotifications(ctx context.Context, userID string) ([]dto.NotificationResponse, error) {
	query := `SELECT id, user_id, order_id, message, status FROM notifications WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []dto.NotificationResponse
	for rows.Next() {
		var n dto.NotificationResponse
		if err := rows.Scan(&n.ID, &n.UserID, &n.OrderID, &n.Message, &n.Status); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if notifications == nil {
		notifications = []dto.NotificationResponse{}
	}

	return notifications, nil
}
