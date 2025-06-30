package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

// HPPOptions defines the configuration for the HPP middleware
type HPPOptions struct {
	CheckQuery                  bool     // Whether to sanitize query parameters (from URL)
	CheckBody                   bool     // Whether to sanitize form body parameters (from POST/PUT)
	CheckBodyOnlyForContentType string   // Only apply body check if this Content-Type is present (e.g., "application/x-www-form-urlencoded")
	WhiteList                   []string // List of allowed parameters (others will be removed)
}

// Hpp returns a middleware that removes duplicate and/or unapproved HTTP parameters
// based on the provided options.
func Hpp(options HPPOptions) func(handler http.Handler) http.Handler {
	fmt.Println("HPP MIDDLEWARE STARTED")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Check body parameters if enabled and conditions match
			if options.CheckBody &&
				r.Method == http.MethodPost &&
				isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				filterBodyParams(r, options.WhiteList) // Clean up body/form parameters
			}

			// Check query parameters if enabled
			if options.CheckQuery && r.URL.Query() != nil {
				filterQueryParams(r, options.WhiteList) // Clean up URL parameters
			}

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
			fmt.Println("HPP MIDDLEWARE ENDED")

		})
	}
}

// isCorrectContentType checks if the request's Content-Type header contains the expected type
func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

// filterBodyParams ensures:
// 1. Only one value per form key is kept (first one).
// 2. Parameters not in the whitelist are deleted.
func filterBodyParams(r *http.Request, whiteList []string) {
	err := r.ParseForm() // Parses form and populates r.Form and r.PostForm
	if err != nil {
		fmt.Println("Error parsing form", err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 0 {
			// Only keep the first value for repeated parameters
			r.Form.Set(k, v[0])
		}

		// If the key is not whitelisted, delete it
		if !isWhiteList(k, whiteList) {
			delete(r.Form, k)
		}
	}
}

// filterQueryParams does the same as filterBodyParams, but for URL query params
func filterQueryParams(r *http.Request, whiteList []string) {
	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 0 {
			query.Set(k, v[0]) // Keep only the first value
		}

		if !isWhiteList(k, whiteList) {
			query.Del(k) // Delete parameter if not whitelisted
		}
	}

	// Rewrite the raw query string with cleaned params
	r.URL.RawQuery = query.Encode()
}

// isWhiteList checks if a given parameter is in the whitelist
func isWhiteList(param string, whiteList []string) bool {
	for _, v := range whiteList {
		if param == v {
			return true
		}
	}
	return false
}
