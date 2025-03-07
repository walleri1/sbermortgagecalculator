// Package middleware implements the logic of measuring the execution time of service requests.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logging.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(begin)
		log.Printf("status_code: %d, duration: %s\n", rw.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
