package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}

	log.Printf("[ERROR] %s %s: %v", r.Method, r.URL.Path, err)

	var appErr *AppError
	if errors.As(err, &appErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(appErr.StatusCode)
		json.NewEncoder(w).Encode(map[string]string{
			"error": appErr.Message,
		})
		return
	}

	// Fallback to 500
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Internal server error",
	})
}
