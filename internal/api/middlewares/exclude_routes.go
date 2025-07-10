package middlewares

import (
	"net/http"
	"strings"
)

func MiddleWareExcludePaths(middlewares func(http.Handler) http.Handler, excludePaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, excludePath := range excludePaths {
				if strings.HasPrefix(r.URL.Path, excludePath) {
					next.ServeHTTP(w, r)
					return
				}
			}
			middlewares(next).ServeHTTP(w, r)
		})
	}
}
