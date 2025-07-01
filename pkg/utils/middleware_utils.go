package utils

import "net/http"

// Middleware is a function that wraps a http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler

func ApplyMiddleWares(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}
	return handler
}
