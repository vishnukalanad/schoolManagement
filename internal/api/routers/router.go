package routers

import (
	"net/http"
	"schoolManagement/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/students", handlers.StudentsHandler)
	mux.HandleFunc("/teachers/", handlers.TeachersHandler)
	mux.HandleFunc("/execs", handlers.ExecsHandler)

	return mux
}
