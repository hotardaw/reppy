package middleware

import "net/http"

func MaxBodySizeMiddleware(maxBodySize int64) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch { // only methods w/ body
				r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
			}
			next(w, r)
		}
	}
}
