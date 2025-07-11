package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-mail/mail/v2"
	"log"
	"net/http"
	"os"
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

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Err: Bad request!", http.StatusBadRequest)
		return
	}
	log.Println("Request : ", request)

	err = r.Body.Close()
	if err != nil {
		http.Error(w, "Err: Bad request!", http.StatusBadRequest)
	}

	var exec models.Exec
	err = sqlconnect.ForgotPasswordDbHandler(request.Email, &exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXPIRY"))
	if err != nil {
		http.Error(w, "Err: Failed to send password reset email!", http.StatusInternalServerError)
		return
	}

	minutes := time.Duration(duration) * time.Minute
	log.Println("Duration : ", minutes)
	expiry := time.Now().Add(minutes).Format(time.RFC3339)
	tokenBytes := make([]byte, 32)

	_, err = rand.Read(tokenBytes)
	if err != nil {
		http.Error(w, "Err: Failed to send password reset email!", http.StatusInternalServerError)
		return
	}

	token := hex.EncodeToString(tokenBytes)
	hashedToken := sha256.Sum256(tokenBytes)
	hashedTokenString := hex.EncodeToString(hashedToken[:])

	err = sqlconnect.ForgotPasswordUpdateDbHandler(exec, hashedTokenString, expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sending email;
	resetUrl := fmt.Sprintf("https://localhost:3000/execs/reset-password/reset/%s", token)
	message := fmt.Sprintf("Forgot your password? Reset your password using the following link,\n%s\nThis reset link is valid for %d minutes!", resetUrl, duration)

	m := mail.NewMessage()
	m.SetHeader("From", "schooladmin@school.com")
	m.SetHeader("To", request.Email)
	m.SetHeader("Subject", "Password reset link")
	m.SetBody("text/html", message)

	log.Println("Email generated : ", message)
	log.Println("Sending email to ", request.Email)

	d := mail.NewDialer("localhost", 1025, "", "")
	err = d.DialAndSend(m)
	if err != nil {
		http.Error(w, "Err: Mail sending failed!", http.StatusInternalServerError)
		return
	}

	log.Println("Dial and send", err)

	fmt.Fprintf(w, "Password reset link has beed shared to %s!", request.Email)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("resetCode")
	type request struct {
		NewPassword     string `json:"newPassword"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req request

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Err: Invalid request body!", http.StatusBadRequest)
		return
	}

	// Checks for empty password;
	if req.NewPassword == "" || req.ConfirmPassword == "" {
		http.Error(w, "Err: Invalid request body!", http.StatusBadRequest)
	}

	// Check if password and confirm password matches;
	if req.NewPassword != req.ConfirmPassword {
		http.Error(w, "Err: Password mismatch!", http.StatusBadRequest)
	}

	var exec models.Exec
	err = sqlconnect.ResetPasswordDbHandler(token, &exec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Hashing the password;
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = sqlconnect.ForgetPasswordResetDbHandler(&exec, hashedPassword)
	if err != nil {
		http.Error(w, "Err: Password reset failed!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "Success",
		Message: "Password updated successfully!",
	})
}

// ExecsHandler - Handler for execs route;
func ExecsHandler(w http.ResponseWriter, r *http.Request) {

}
