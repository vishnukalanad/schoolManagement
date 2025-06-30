package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	fmt.Println("RESPONSE TIME CHECK MIDDLEWARE STARTED")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request in ResponseTime middleware")
		startTime := time.Now()

		// Create a custom response writer;
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		// Calculate the duration;
		duration := time.Since(startTime)
		w.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter, r)
		duration = time.Since(startTime)
		// Log the request details;
		fmt.Printf("ResponseTimeMiddleware took %s, status %d\n", duration, wrappedWriter.status)
		fmt.Println("Sent response from ResponseTime middleware")
		fmt.Println("RESPONSE TIME CHECK MIDDLEWARE ENDED")

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
