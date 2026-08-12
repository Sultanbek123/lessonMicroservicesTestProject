package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"order-service/internal/dto"
	"order-service/internal/middleware"
	"order-service/internal/service"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler(svc *service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

func (h *OrderHandler) HandleOrders(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// Оборачиваем контекст запроса в контекст с таймаутом и передаем его во все слои
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	if r.Method == http.MethodGet {
		resp, err := h.svc.GetOrders(ctx, userID)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}
		return json.NewEncoder(w).Encode(resp)
	}

	if r.Method == http.MethodPost {
		jwtCookie, err := r.Cookie("jwt")
		if err != nil {
			return &AppError{Err: err, Message: "не аутентицирован", StatusCode: http.StatusUnauthorized}
		}
		var req dto.CreateOrderRequest
		if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
			return &AppError{Err: err, Message: "неверный формат данных", StatusCode: http.StatusBadRequest}
		}
		defer r.Body.Close()

		resp, err := h.svc.CreateOrder(ctx, userID, jwtCookie, req)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(resp)
	}

	return &AppError{Err: errors.New("method not allowed"), Message: "метод не поддерживается", StatusCode: http.StatusMethodNotAllowed}
}
