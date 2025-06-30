package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	mw "schoolManagement/internal/api/middlewares"
	"strconv"
	"strings"
	"sync"
)

// Even though the user struct is private (not starting with uppercase), the field values after made public (Name, Age and City).
// This is because, while unmarshalling the field values are accessed by another package (encoding/json), so instead of struct, the field values are made public;
type user struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

// rootHandler - Handler for root route;
func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	_, err := w.Write([]byte("Hello Root Route"))
	if err != nil {
		fmt.Println("Error from root handler ", err)
	}
}

// studentsHandler - Handler for students route;
func studentsHandler(w http.ResponseWriter, r *http.Request) {
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

type teacher struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

var (
	teachers = make(map[int]teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = teacher{
		Id:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Mathematics",
	}
	nextId++
	teachers[nextId] = teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Physics",
	}
	nextId++
	teachers[nextId] = teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Physics",
	}
	nextId++
}

// getTeachersHandler - this will handle the business logic for get teachers;
func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	// Extract path params;
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	// Extract query params;
	firstName := r.URL.Query().Get("first_name")
	lastName := r.URL.Query().Get("last_name")

	teacherList := make([]teacher, 0, len(teachers))

	var emptyVal = true

	if idStr != "" {
		for _, value := range teachers {
			if idStr == strconv.Itoa(value.Id) {
				emptyVal = false
				teacherList = append(teacherList, value)
			}
		}
	} else {
		for _, value := range teachers {
			if (firstName == "" || value.FirstName == firstName) && (lastName == "" || value.LastName == lastName) {
				emptyVal = false
				teacherList = append(teacherList, value)
			}
		}
	}

	if emptyVal {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachers),
		Data:   teacherList,
	}

	// Sets the content type as JSON;
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error from getTeachersHandler ", err)
	}

}

// addTeachersHandler - handles the incoming post requests;
func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var requestData []teacher
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	addedTeachers := make([]teacher, len(requestData))
	for i, requestDatum := range requestData {
		requestDatum.Id = nextId
		teachers[nextId] = requestDatum
		addedTeachers[i] = requestDatum
		nextId++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// teachersHandler - Handler for teachers route;
func teachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	switch r.Method {
	case http.MethodGet:
		// Call the get handler function;
		getTeachersHandler(w, r)
		return
	case http.MethodPost:
		// Call the post handler function;
		addTeachersHandler(w, r)
		return
	case http.MethodPatch:
		_, err := w.Write([]byte("Hello this is a PATCH method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodPut:
		_, err := w.Write([]byte("Hello this is a PUT method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	case http.MethodDelete:
		_, err := w.Write([]byte("Hello this is a DELETE method call to teachers api!"))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	default:
		var msg string = "This a " + r.Method + " call to teachers"
		_, err := w.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error from teachers handler ", err)
		}
		return
	}

}

// execsHandler - Handler for execs route;
func execsHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	const port string = ":3000"

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/students", studentsHandler)
	mux.HandleFunc("/teachers/", teachersHandler)
	mux.HandleFunc("/execs", execsHandler)

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}

	// Should be uncommented while using all middlewares;
	//rl := mw.NewRateLimiter(5, time.Minute)
	//
	//hppOptions := mw.HPPOptions{
	//	CheckQuery:                  true,
	//	CheckBody:                   true,
	//	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	//	WhiteList:                   []string{"name", "age"},
	//}

	// All the other middlewares are passed as argument;
	//secureMux := mw.Cors(rl.Middleware(mw.ResponseTimeMiddleware(mw.SecurityHandler(mw.CompressionMiddleware(mw.Hpp(hppOptions)(mw.Cors(mux)))))))
	// An enhanced and efficient way to apply middlewares;
	//secureMux := applyMiddleWares(mux, mw.Hpp(hppOptions), mw.CompressionMiddleware, mw.SecurityHandler, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)

	// For this server we will use mw.SecurityHandler alone now;
	secureMux := mw.SecurityHandler(mux)

	server := &http.Server{
		Addr:      port,
		TLSConfig: tlsConfig,
		Handler:   secureMux,
	}
	fmt.Println("Server is running on port ", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatal("Error starting the server : ", err)
	}
}

// Middleware is a function that wraps a http.Handler with additional functionality
type Middleware func(http.Handler) http.Handler

func applyMiddleWares(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, middleware := range middleware {
		handler = middleware(handler)
	}
	return handler
}
