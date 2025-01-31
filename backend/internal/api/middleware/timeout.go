package middleware

import (
	"context"
	"net/http"
	"time"

	"go-fitstat/backend/internal/api/response"
)

func TimeoutMiddleware(timeout time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan bool)
			go func() {
				next.ServeHTTP(w, r)
				done <- true
			}()

			select {
			case <-done: // if handler completes normally
				return
			case <-ctx.Done(): // if timeout expires
				response.SendError(w, "Request timed out", http.StatusRequestTimeout)
				return
			}
		}
	}
}
