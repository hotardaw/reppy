package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// written w/ 3 layers so we can initialize it in same fashion as other middlewares in main.go
func LoggingMiddleware() func(http.HandlerFunc) http.HandlerFunc { // layer 1: factory func returns a middleware func
	return func(next http.HandlerFunc) http.HandlerFunc { // layer 2: middleware func takes next handler & returns new handler
		return func(w http.ResponseWriter, r *http.Request) { // layer 3: the actual logging handler
			// skip logging favicon reqs
			if r.URL.Path == "/favicon.ico" {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()
			rw := NewResponseWriter(w)

			// process req
			next.ServeHTTP(rw, r)

			// log after req is processed
			duration := time.Since(start)

			log.Printf(
				"%s %s %s %d %v",
				r.Method,
				r.RequestURI,
				r.RemoteAddr,
				rw.statusCode,
				duration,
			)
		}
	}
}
