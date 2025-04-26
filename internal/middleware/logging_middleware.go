package middleware

import (
	"log"
	"net/http"
	"time"
)

// loggingResponseWriter wraps http.ResponseWriter to capture status code and response size.
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

// newLoggingResponseWriter constructs a new loggingResponseWriter.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// Default status code is 200 OK until WriteHeader is called
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the status code.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Write captures the size of the response.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseSize += size
	return size, err
}

// LoggingMiddleware logs client IP, HTTP method, URI, status code, response size and duration.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		clientIP := r.RemoteAddr
		method := r.Method
		uri := r.RequestURI
		status := lrw.statusCode
		size := lrw.responseSize

		log.Printf(
			"%s %s %s → %d (%d bytes) in %v",
			clientIP,
			method,
			uri,
			status,
			size,
			duration,
		)
	})
}
