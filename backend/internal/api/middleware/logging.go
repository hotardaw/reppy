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

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
