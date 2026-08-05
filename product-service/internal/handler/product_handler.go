package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"product-service/internal/dto"
	"product-service/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// Оборачиваем контекст запроса в контекст с таймаутом и передаем его во все слои
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if r.Method == http.MethodGet {
		resp, err := h.svc.GetProducts(ctx)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}
		return json.NewEncoder(w).Encode(resp)
	}

	if r.Method == http.MethodPost {
		var req dto.CreateProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return &AppError{Err: err, Message: "неверный формат данных", StatusCode: http.StatusBadRequest}
		}
		defer r.Body.Close()

		resp, err := h.svc.CreateProduct(ctx, req)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(resp)
	}

	return &AppError{Err: errors.New("method not allowed"), Message: "метод не поддерживается", StatusCode: http.StatusMethodNotAllowed}
}
