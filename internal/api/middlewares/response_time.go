package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request in ResponseTime middleware")
		startTime := time.Now()

		// Create a custom response writer;
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)
		// Calculate the duration;
		duration := time.Since(startTime)
		// Log the request details;
		fmt.Printf("ResponseTimeMiddleware took %s, status %d\n", duration, wrappedWriter.status)
		fmt.Println("Sent response from ResponseTime middleware")
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
