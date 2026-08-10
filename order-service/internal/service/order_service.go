package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"order-service/internal/client"
	"order-service/internal/dto"
	"order-service/internal/repository"
)

type OrderService struct {
	repo          *repository.OrderRepository
	productClient *client.ProductClient
}

func NewOrderService(repo *repository.OrderRepository, productClient *client.ProductClient) *OrderService {
	return &OrderService{repo: repo, productClient: productClient}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, jwtCookie *http.Cookie, req dto.CreateOrderRequest) (dto.OrderResponse, error) {
	// TODO: перед созданием заказа нужно сходить в product-service Get запросом
	// и убедиться что товар с req.ProductID существует (и подтянуть его цену).
	// Пока просто доверяем ProductID, который прислал клиент.
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
