package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

// CompressionMiddleware is a middleware that compresses HTTP responses using gzip
func CompressionMiddleware(next http.Handler) http.Handler {
	fmt.Println("COMPRESSION MIDDLEWARE STARTED")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Step 1: Check if the client supports gzip encoding
		// If not, just forward the request to the next handler and return
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return // Important! Prevents continuing to write gzip when client doesn't support it
		}

		// Step 2: Set the Content-Encoding header to tell the client we are using gzip
		w.Header().Set("Content-Encoding", "gzip")

		// Step 3: Create a new gzip writer that wraps the original ResponseWriter
		gz := gzip.NewWriter(w)

		// Step 4: Make sure to close the gzip writer when the handler completes
		defer func() {
			err := gz.Close()
			if err != nil {
				// If closing the gzip writer fails, return a 500 error
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()

		// Step 5: Wrap the ResponseWriter with our custom gzipResponseWriter
		// This ensures all writes go through gzip.Writer
		gw := &gzipResponseWriter{ResponseWriter: w, writer: gz}

		// Step 6: Call the next handler, passing in our wrapped writer
		next.ServeHTTP(gw, r)

		// Optional logging for debug purposes
		fmt.Println("Sent response from compression middleware")
	})
}

// gzipResponseWriter is a custom writer that wraps the original ResponseWriter
// and compresses the response using gzip.
type gzipResponseWriter struct {
	http.ResponseWriter // Embedded interface: allows our struct to behave like a ResponseWriter

	writer *gzip.Writer // Pointer to gzip.Writer: handles actual compression.
	// We use a pointer so that we are writing to the same gzip stream instance
	// and can later close it. This avoids copying and allows shared state.
}

// Write overrides the default Write method of http.ResponseWriter
// It writes the data through gzip.Writer instead of sending it directly.
func (g gzipResponseWriter) Write(b []byte) (int, error) {
	// Data gets compressed before being written to the underlying response
	return g.writer.Write(b)
}
