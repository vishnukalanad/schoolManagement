package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeachersHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	return mux
}
