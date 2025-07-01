package handlers

import (
	"fmt"
	"net/http"
)

// StudentsHandler - Handler for students route;
func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		_, err := w.Write([]byte("Hello this is a GET method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPost:
		_, err := w.Write([]byte("Hello this is a POST method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to students api!"))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to students"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from students handler ", err)
		}
		return
	}
}
