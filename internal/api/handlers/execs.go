package handlers

import (
	"fmt"
	"net/http"
)

// ExecsHandler - Handler for execs route;
func ExecsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("Hello this is a GET method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPost:
		form := r.Form
		fmt.Println("Form : ", form)
		_, err := w.Write([]byte("Hello this is a POST method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to execs api!"))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to execs"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from execs handler ", err)
		}
		return
	}
}
