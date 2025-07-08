package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolManagement/internal/models"
	"schoolManagement/internal/repositories/sqlconnect"
)

// Students Handlers;

// GetStudentsHandler - Handler to handle get students list route;
func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Student array to hold the fetched students from DB;
	var students []models.Student

	// Calls the DB handler to perform query and get data;
	err, students := sqlconnect.GetStudentDbHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Prints the students list returned from DB;
	fmt.Println("Query results : ", students)

	// Prepares the response;
	response := struct {
		Status   string           `json:"status"`
		Students []models.Student `json:"students"`
		Count    int              `json:"count"`
	}{
		Status:   "Success",
		Students: students,
		Count:    len(students),
	}

	// Sends the response;
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var students []models.Student
	err, students := sqlconnect.AddStudentsDbHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status   string           `json:"status"`
		Students []models.Student `json:"students"`
		Count    int              `json:"count"`
	}{
		Status:   "Success",
		Students: students,
		Count:    len(students),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {

}

// Students By ID Handlers;

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {

}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {

}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {

}
