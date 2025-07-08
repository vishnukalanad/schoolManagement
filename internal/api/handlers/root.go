package handlers

import (
	"net/http"
)

// RootHandler - Handler for root route;
func RootHandler(w http.ResponseWriter, r *http.Request) {

	_, err := w.Write([]byte("Welcome to the school API!"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
