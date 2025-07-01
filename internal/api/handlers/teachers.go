package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolManagement/internal/models"
	"strconv"
	"strings"
	"sync"
)

var (
	teachers = make(map[int]models.Teacher)
	mutex    = &sync.Mutex{}
	nextId   = 1
)

func init() {
	teachers[nextId] = models.Teacher{
		Id:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Mathematics",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Physics",
	}
	nextId++
	teachers[nextId] = models.Teacher{
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

	teacherList := make([]models.Teacher, 0, len(teachers))

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
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
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

	var requestData []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, len(requestData))
	for i, requestDatum := range requestData {
		requestDatum.Id = nextId
		teachers[nextId] = requestDatum
		addedTeachers[i] = requestDatum
		nextId++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
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

// TeachersHandler - Handler for teachers route;
func TeachersHandler(w http.ResponseWriter, r *http.Request) {
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
