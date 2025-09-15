package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs information about each HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time
		start := time.Now()

		// Create a response writer that captures the status code
		lrw := newLoggingResponseWriter(w)

		// Process the request
		next.ServeHTTP(lrw, r)

		// Calculate the duration
		duration := time.Since(start)

		// Log the request details
		log.Printf(
			"[%s] %s %s %d %s %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			duration,
			r.UserAgent(),
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
