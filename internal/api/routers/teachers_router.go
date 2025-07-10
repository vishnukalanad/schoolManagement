package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func TeachersRouter() *http.ServeMux {

	mux := http.NewServeMux()

	//mux.HandleFunc("GET /", handlers.RootHandler)

	// General handlers for teachers route;
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers", handlers.AddTeachersHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("DELETE /teachers", handlers.DeleteTeachersHandler)

	// By ID handlers for teachers route;
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeachersHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	// Sub routes for teacher;
	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetStudentsByTeacherHandler)
	mux.HandleFunc("GET /teachers/{id}/studentCount", handlers.GetStudentsCountByTeacherHandler)

	return mux
}
