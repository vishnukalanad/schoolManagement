package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool
	CheckBody                   bool
	CheckBodyOnlyForContentType string
	WhiteList                   []string
}

func Hpp(options HPPOptions) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				// Filter the body params;
				filterBodyParams(r, options.WhiteList)
			}

			if options.CheckQuery && r.URL.Query() != nil {
				// Filter the body params;
				filterQueryParams(r, options.WhiteList)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, whiteList []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form", err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 0 {
			r.Form.Set(k, v[0])
		}

		if !isWhiteList(k, whiteList) {
			delete(r.Form, k)
		}
	}

}

func filterQueryParams(r *http.Request, whiteList []string) {
	query := r.URL.Query()

	for k, v := range query {
		if len(v) > 0 {
			query.Set(k, v[0])
		}

		if !isWhiteList(k, whiteList) {
			query.Del(k)
		}
	}

	r.URL.RawQuery = query.Encode()
}

func isWhiteList(param string, whiteList []string) bool {
	for _, v := range whiteList {
		if param == v {
			return true
		}
	}
	return false
}
