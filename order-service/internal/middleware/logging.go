package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

type key int

const (
	RequestIDKey key = iota
)

// LoggingMiddleware logs the request and injects a Request ID into the context.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := time.Now().UnixNano()

		// Inject Request ID into context
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		r = r.WithContext(ctx)

		log.Printf("[INFO] Started %s %s (ReqID: %d)", r.Method, r.URL.Path, reqID)

		next.ServeHTTP(w, r)

		log.Printf("[INFO] Completed %s %s in %v (ReqID: %d)", r.Method, r.URL.Path, time.Since(start), reqID)
	})
}
