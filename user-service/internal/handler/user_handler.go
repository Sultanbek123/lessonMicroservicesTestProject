package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"user-service/internal/dto"
	"user-service/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) HandleUsers(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// Оборачиваем контекст запроса в контекст с таймаутом и передаем его во все слои
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if r.Method == http.MethodPost {
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return &AppError{Err: err, Message: "неверный формат данных", StatusCode: http.StatusBadRequest}
		}
		defer r.Body.Close()

		resp, err := h.svc.CreateUser(ctx, req)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}

		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(resp)
	}

	return &AppError{Err: errors.New("method not allowed"), Message: "метод не поддерживается", StatusCode: http.StatusMethodNotAllowed}
}
