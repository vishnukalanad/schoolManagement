package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolManagement/internal/models"
	"schoolManagement/internal/repositories/sqlconnect"
	"strconv"
)

// Students Handlers;

// GetStudentsHandler - Handler to handle get students list route;
func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Student array to hold the fetched students from DB;
	var students []models.Student

	// Calls the DB handler to perform query and get data;
	err, students := sqlconnect.GetStudentsDbHandler(r)
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

// PatchStudentsHandler - Handles patch students operation;
func PatchStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		fmt.Println("Error: Failed to decode response body!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sqlconnect.PatchStudentsDbHandler(updates)
	if err != nil {
		fmt.Println("Error: Failed to patch students!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteStudentsHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		fmt.Println("Error: Failed to decode response body!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err, deletedIds := sqlconnect.DeleteStudentsDbHandler(ids)
	if err != nil {
		fmt.Println("Error: Failed to delete students!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Ids    []int  `json:"ids"`
	}{
		Status: "Success",
		Ids:    deletedIds,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Students By ID Handlers;

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err, student := sqlconnect.GetStudentHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string         `json:"status"`
		Student models.Student `json:"student"`
	}{
		Status:  "Success",
		Student: student,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from query params;
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Println("Err : ID parsing failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updatedStudent models.Student
	// Decodes the request body;
	err = json.NewDecoder(r.Body).Decode(&updatedStudent)
	if err != nil {
		fmt.Println("Err : Student JSON decoding failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update CRUD operation;
	err, student := sqlconnect.UpdateStudentsDbHandler(id, updatedStudent)
	if err != nil {
		fmt.Println("Err : Student Update Failed!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare and send the response;
	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status   string           `json:"status"`
		Message  string           `json:"message"`
		Students []models.Student `json:"students"`
	}{
		Status:   "Success",
		Message:  "Student details updated successfully!",
		Students: student,
	}
	err = json.NewEncoder(w).Encode(response)
}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Err : ID parsing failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		fmt.Println("Error: Failed to decode response body!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err, _ = sqlconnect.PatchStudentDbHandler(id, updates)
	if err != nil {
		fmt.Println("Error: Failed to patch students!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Err : ID parsing failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = sqlconnect.DeleteStudentDbHandler(id)
	if err != nil {
		fmt.Println("Error: Failed to delete students!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Ids    int    `json:"id"`
	}{
		Status: "Success",
		Ids:    id,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
}
