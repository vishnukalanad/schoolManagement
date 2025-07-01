package handlers

import (
	"fmt"
	"net/http"
)

// RootHandler - Handler for root route;
func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	_, err := w.Write([]byte("Hello Root Route"))
	if err != nil {
		fmt.Println("Error from root handler ", err)
	}
}
