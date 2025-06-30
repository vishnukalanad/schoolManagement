package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// RateLimiter is a struct to hold rate limiting logic per visitor (based on IP).
type RateLimiter struct {
	mu        sync.Mutex     // Mutex to protect concurrent access to the visitor map
	visitor   map[string]int // Stores number of requests per IP address
	limit     int            // Max number of requests allowed per reset time
	resetTIme time.Duration  // Time interval after which visitor counts are reset
}

// NewRateLimiter creates and returns a new RateLimiter instance with the given request limit and reset interval.
func NewRateLimiter(limit int, resetTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitor:   make(map[string]int),
		limit:     limit,
		resetTIme: resetTime,
	}

	// Start a background goroutine to reset the visitor map periodically.
	go rl.resetVisitorCount()

	return rl
}

// resetVisitorCount resets the visitor map at regular intervals.
// This allows each IP to start fresh after `resetTime`.
func (rl *RateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTIme)

		// Lock before modifying the shared map
		rl.mu.Lock()
		rl.visitor = make(map[string]int) // Clear all counts
		rl.mu.Unlock()
	}
}

// Middleware is the actual HTTP middleware function that applies rate limiting logic.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	fmt.Println("RATE LIMITER MIDDLEWARE STARTED")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Step 1: Lock the mutex to protect the visitor map from concurrent access
		rl.mu.Lock()
		defer rl.mu.Unlock() // Always unlock at the end of this function

		// Step 2: Use the client's IP address to identify them
		visitorIP := r.RemoteAddr // Note: This includes port; may want to strip it in real-world cases

		// Step 3: Increment request count for this visitor
		rl.visitor[visitorIP]++
		fmt.Printf("Visitor count from %s is %d\n", visitorIP, rl.visitor[visitorIP])

		// Step 4: Check if the visitor has exceeded the request limit
		if rl.visitor[visitorIP] > rl.limit {
			// Respond with HTTP 429 Too Many Requests if limit is exceeded
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		// Step 5: Forward the request to the next handler if within the limit
		next.ServeHTTP(w, r)
		fmt.Println("RATE LIMITER MIDDLEWARE ENDED")

	})
}
