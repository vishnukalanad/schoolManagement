package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schoolManagement/internal/models"
	"schoolManagement/internal/repositories/sqlconnect"
	"strconv"
)

// GetTeachersHandler - this will handle the business logic for get teachers;
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var teachers []models.Teacher
	err, teachers := sqlconnect.GetTeachersDbHandler(r, teachers)
	if err != nil {
		return
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teachers),
		Data:   teachers,
	}

	// Sets the content type as JSON;
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error from getTeachersHandler ", err)
	}

}

// GetTeacherHandler - gets details of a single teacher based on ID;
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// Extract path params;
	idStr := r.PathValue("id")
	err, teacher := sqlconnect.GetTeacherDbHandler(idStr)
	if err != nil {
		return
	}

	response := struct {
		Status string         `json:"status"`
		Count  int            `json:"count"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  1,
		Data:   teacher,
	}

	// Sets the content type as JSON;
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error from getTeachersHandler ", err)
	}

}

// AddTeachersHandler - handles the incoming post requests;
func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var newTeachers []models.Teacher

	err, addedTeachers := sqlconnect.AddTeacherDbHandler(w, r, newTeachers)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	var updatedTeachers models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	err = sqlconnect.UpdateTeachersDbHandler(id, updatedTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeachers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PatchTeachersHandler - Patches multiple teachers details in a go;
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload!", http.StatusInternalServerError)
		fmt.Println("Err : Invalid request body!", err)
		return
	}

	err = sqlconnect.PatchTeachersDbHandler(updates)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// PatchTeacherHandler - Patches single teacher details based on the ID;
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	err, existingTeacher := sqlconnect.PatchTeacherDbHandler(id, updates)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	err = sqlconnect.DeleteTeacherDbHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Status string
		Id     int
	}{
		Status: "Teacher deleted!",
		Id:     id,
	})
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {

	var ids []int
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Incorrect request!", http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	err, deletedIds := sqlconnect.DeleteTeachersDbHandler(ids)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Status string `json:"status"`
		Id     []int  `json:"deleted_ids"`
	}{
		Status: "Teacher deleted!",
		Id:     deletedIds,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
}

//----------------------

func GetStudentsByTeacherHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")
	var students []models.Student

	err, students := sqlconnect.GetStudentsByTeacherDbHandler(w, teacherId, students)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status   string           `json:"status"`
		Students []models.Student `json:"students"`
		Count    int              `json:"count"`
	}{
		Status:   "Success",
		Students: students,
		Count:    len(students),
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetStudentsCountByTeacherHandler(w http.ResponseWriter, r *http.Request) {
	teacherId := r.PathValue("id")

	err, count := sqlconnect.GetStudentsCountByTeacherDbHandler(w, teacherId)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Count  int    `json:"count"`
	}{
		Status: "Success",
		Count:  count,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
