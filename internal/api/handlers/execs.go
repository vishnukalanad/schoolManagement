package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"schoolManagement/internal/models"
	"schoolManagement/internal/repositories/sqlconnect"
	"schoolManagement/pkg/utils"
	"strconv"
	"time"
)

// ******** GENERAL HANDLERS ********

// GetExecsHandler - Handles the get route of execs;
func GetExecsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GET EXECS ROUTE")
	// Execs array to hold the fetched students from DB;
	var execs []models.Exec

	// Calls the DB handler to perform query and get data;
	err, execs := sqlconnect.GetExecsDbHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Prints the execs list returned from DB;
	fmt.Println("Query results : ", execs)

	// Prepares the response;
	response := struct {
		Status string        `json:"status"`
		Execs  []models.Exec `json:"execs"`
		Count  int           `json:"count"`
	}{
		Status: "Success",
		Execs:  execs,
		Count:  len(execs),
	}

	// Sends the response;
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// AddExecsHandler - Onboarding of execs;
func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	var execs []models.Exec
	err, execs := sqlconnect.AddExecsDbHandler(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := struct {
		Status string        `json:"status"`
		Execs  []models.Exec `json:"execs"`
		Count  int           `json:"count"`
	}{
		Status: "Success",
		Execs:  execs,
		Count:  len(execs),
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// PatchExecsHandler - Handles the update of execs (PATCH method);
func PatchExecsHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		fmt.Println("Error: Failed to decode response body!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sqlconnect.PatchExecsDbHandler(updates)
	if err != nil {
		fmt.Println("Error: Failed to patch execs!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// DeleteExecsHandler - Deleting execs;
func DeleteExecsHandler(w http.ResponseWriter, r *http.Request) {
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

// ******** BY ID HANDLERS ********

func GetExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err, exec := sqlconnect.GetExecsByIdHandler(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status  string      `json:"status"`
		Student models.Exec `json:"exec"`
	}{
		Status:  "Success",
		Student: exec,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PatchExecByIdHandler(w http.ResponseWriter, r *http.Request) {
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

	err, _ = sqlconnect.PatchExecByIdDbHandler(id, updates)
	if err != nil {
		fmt.Println("Error: Failed to patch students!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func DeleteExecByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Err : ID parsing failed!")
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = sqlconnect.DeleteExecByIdDbHandler(id)
	if err != nil {
		fmt.Println("Error: Failed to delete exec!")
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Data validation;

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Err: Username or password is empty!", http.StatusBadRequest)
		return
	}

	// Search for user;
	user := &models.Exec{}
	err = sqlconnect.LoginDbHandler(req.Username, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Is user active;
	if user.Inactive {
		http.Error(w, "Err: User is inactive!", http.StatusForbidden)
		return
	}

	// Verify password;
	log.Println(req.Username, user.Password)
	err = utils.PasswordValidate(user.Password, req.Password)
	if err != nil {
		_ = utils.HandleError(err, "Err: Error from password validate!")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate JWT token;
	usrId := strconv.Itoa(user.Id)
	token, err := utils.SignToken(usrId, req.Username, user.Role)
	if err != nil {
		http.Error(w, "Err: Token generation failed!", http.StatusInternalServerError)
	}

	// Send token as a response/cookie;
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
		SameSite: http.SameSiteStrictMode,
	})

	// Response;
	response := struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}{
		Status: "Success",
		Token:  token,
	}

	err = json.NewEncoder(w).Encode(response)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte("Logged out!"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var request models.UpdatePasswordRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if request.CurrentPassword == "" || request.NewPassword == "" {
		http.Error(w, "Err: Current or new password is empty!", http.StatusBadRequest)
		return
	}

	err, token := sqlconnect.UpdatePasswordDbHandler(userId, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string `json:"status"`
		Token  string `json:"token"`
	}{
		Status: "Success",
		Token:  token,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// ExecsHandler - Handler for execs route;
func ExecsHandler(w http.ResponseWriter, r *http.Request) {

}
