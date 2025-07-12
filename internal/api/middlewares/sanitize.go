package middlewares

import (
	"errors"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"log"
	"net/http"
	"net/url"
	"schoolManagement/pkg/utils"
)

func XSSMiddleware(next http.Handler) http.Handler {
	log.Println("----( Starting XSS middleware )----")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Sanitize URL path;
		sanitizePath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		fmt.Println(sanitizePath)

		// Sanitize query params;
		params := r.URL.Query()
		sanitizedQuery := make(map[string][]string)
		for k, v := range params {
			sanitizedKey, err := clean(k)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var sanitizedValues []string
			for _, value := range v {
				cleanValue, err := clean(value)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizedQuery[sanitizedKey.(string)] = sanitizedValues
			log.Println(sanitizedQuery)
		}

		r.URL.Path = sanitizePath.(string)
		r.URL.RawQuery = url.Values(sanitizedQuery).Encode()
		next.ServeHTTP(w, r)
	})
}

func clean(data interface{}) (interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		for i, val := range v {
			v[i] = sanitizeVal(val)
		}
		return v, nil
	case []interface{}:
		for i, val := range v {
			v[i] = sanitizeVal(val)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil
	default:
		return nil, utils.HandleError(errors.New("unsupported type"), "Unsupported type")
	}

}

func sanitizeVal(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return sanitizeString(v)
	case map[string]interface{}:
		for i, val := range v {
			v[i] = sanitizeVal(val)
		}
		return v
	case []interface{}:
		for i, val := range v {
			v[i] = sanitizeVal(val)
		}
		return v
	default:
		return v
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
