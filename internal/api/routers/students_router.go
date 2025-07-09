package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func StudentsRouter() *http.ServeMux {
	mux := http.NewServeMux()
	// General handlers for students route;
	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)
	mux.HandleFunc("PATCH /students/", handlers.PatchStudentsHandler)
	mux.HandleFunc("DELETE /students/", handlers.DeleteStudentsHandler)

	// By ID handlers for students route;
	mux.HandleFunc("GET /students/{id}", handlers.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentHandler)

	return mux
}
