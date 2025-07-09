package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func ExecsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	mux.HandleFunc("GET /execs", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.ExecsHandler)

	// By ID handlers for students route;
	mux.HandleFunc("GET /execs/{id}", handlers.ExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.ExecsHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.ExecsHandler)

	// Auth routes
	mux.HandleFunc("POST /execs/login", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/logout", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/forgot-password", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/reset-password/{resetCode}", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/{id}/update-password", handlers.ExecsHandler)
	return mux
}
