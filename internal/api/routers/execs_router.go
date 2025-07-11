package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func ExecsRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/execs", handlers.ExecsHandler)

	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)

	// By ID handlers for students route;
	mux.HandleFunc("GET /execs/{id}", handlers.GetExecByIdHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecByIdHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecByIdHandler)

	// Auth routes
	mux.HandleFunc("POST /execs/login", handlers.LoginHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutHandler)
	mux.HandleFunc("POST /execs/forgot-password", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/reset-password/reset/{resetCode}", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs/{id}/update-password", handlers.UpdatePasswordHandler)
	return mux
}
