package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"notification-service/internal/middleware"
	"notification-service/internal/service"
)

type NotificationHandler struct {
	svc *service.NotificationService
}

func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) HandleNotifications(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	// Оборачиваем контекст запроса в контекст с таймаутом и передаем его во все слои
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	userID, _ := r.Context().Value(middleware.UserIDKey).(string)

	if r.Method == http.MethodGet {
		resp, err := h.svc.GetNotifications(ctx, userID)
		if err != nil {
			return &AppError{Err: err, Message: "внутренняя ошибка сервера", StatusCode: http.StatusInternalServerError}
		}
		return json.NewEncoder(w).Encode(resp)
	}

	return &AppError{Err: errors.New("method not allowed"), Message: "метод не поддерживается", StatusCode: http.StatusMethodNotAllowed}
}
