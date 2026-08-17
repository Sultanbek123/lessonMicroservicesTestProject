package service

import (
	"context"
	"fmt"

	"notification-service/internal/dto"
	"notification-service/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification создаёт запись об уведомлении для пользователя.
//
// Сейчас этот метод дергается только "руками" (через будущий Kafka consumer,
// который слушает топик order.created из order-service). Сам consumer
// в проект намеренно не добавлен — это часть лекции про Kafka.
func (s *NotificationService) CreateNotification(ctx context.Context, userID string, orderID int) (dto.NotificationResponse, error) {
	message := fmt.Sprintf("Ваш заказ #%d принят в обработку", orderID)

	id, err := s.repo.CreateNotification(ctx, userID, orderID, message)
	if err != nil {
		return dto.NotificationResponse{}, err
	}

	return dto.NotificationResponse{
		ID:      id,
		UserID:  userID,
		OrderID: orderID,
		Message: message,
		Status:  "created",
	}, nil
}

func (s *NotificationService) GetNotifications(ctx context.Context, userID string) ([]dto.NotificationResponse, error) {
	return s.repo.GetNotifications(ctx, userID)
}
