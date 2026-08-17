package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"order-service/internal/client"
	"order-service/internal/dto"
	"order-service/internal/event"
	"order-service/internal/kafka"
	"order-service/internal/repository"
)

type OrderService struct {
	repo          *repository.OrderRepository
	productClient *client.ProductClient
	orderProducer *kafka.OrderProducer
}

func NewOrderService(repo *repository.OrderRepository, productClient *client.ProductClient, orderProducer *kafka.OrderProducer) *OrderService {
	return &OrderService{repo: repo, productClient: productClient, orderProducer: orderProducer}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, jwtCookie *http.Cookie, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	products, err := s.productClient.GetProducts(ctx, jwtCookie)
	if err != nil {
		return dto.OrderResponse{}, fmt.Errorf("can not get products %w", err)
	}

	var found bool = false
	for _, p := range products {
		if p.ID == req.ProductID {
			found = true
			break
		}
	}
	if found == false {
		return dto.OrderResponse{}, errors.New("item not found")
	}

	id, err := s.repo.CreateOrder(ctx, userID, req)
	if err != nil {
		return dto.OrderResponse{}, err
	}

	var orderCreated event.OrderCreated
	orderCreated.OrderID = id
	orderCreated.UserID = userID
	orderCreated.ProductID = req.ProductID
	orderCreated.Quantity = req.Quantity

	err = s.orderProducer.PublishOrderCreated(ctx, orderCreated)
	if err != nil {
		log.Printf("publish ивента не удался")
	}

	return dto.OrderResponse{
		ID:        id,
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Status:    "created",
	}, nil
}

func (s *OrderService) GetOrders(ctx context.Context, userID string) ([]dto.OrderResponse, error) {
	return s.repo.GetOrders(ctx, userID)
}
