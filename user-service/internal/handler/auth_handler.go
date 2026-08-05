package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"user-service/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error" : "неверный формат данных"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, `{"error" : "неверное имя пользователя или пароль"}`, http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // Set to true in production
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "успешный вход"})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error" : "метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "вы успешно вышли"})
}
