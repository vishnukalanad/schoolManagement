package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
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
	err, addedTeachers := sqlconnect.AddTeacherDbHandler(r, newTeachers)
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
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func UpdateTeachersHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	var updatedTeachers models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeachers)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	err = sqlconnect.UpdateTeachersDbHandler(id, updatedTeachers)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(updatedTeachers)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

// PatchTeachersHandler - Patches multiple teachers details in a go;
func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Patch teachers handler started!")
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	var updates []map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request payload!", http.StatusInternalServerError)
		fmt.Println("Err : Invalid request body!", err)
		return
	}

	// DB Transaction Beginning;
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB transaction failed!", err)
		return
	}

	for _, update := range updates {
		id := fmt.Sprintf("%v", update["id"])
		log.Println("TEACHER ID : ", id)
		//if !ok {
		//	err := tx.Rollback()
		//	if err != nil {
		//		return
		//	}
		//	http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		//	fmt.Println("Err : Invalid Teacher ID", err)
		//	return
		//}

		//id, err := strconv.Atoi(idStr)
		//if err != nil {
		//	err := tx.Rollback()
		//	if err != nil {
		//		http.Error(w, "DB Error!", http.StatusInternalServerError)
		//		return
		//	}
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//	fmt.Println("Err : Error while converting teacher ID", err)
		//}

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.Id,
			&teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Teacher not found!", http.StatusNotFound)
				fmt.Println("Err : Teacher not found", err)
				return
			}
		}

		// Apply updates using reflect;
		teacherVal := reflect.ValueOf(&teacherFromDb).Elem()
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skips updating the id field;
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							err = tx.Rollback()
							if err != nil {
								http.Error(w, "Invalid field value!", http.StatusBadRequest)
								fmt.Println("Err : Invalid field value", err)
								return
							}
							break
						}
					}
				}
			}
		}

		_, err = tx.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", teacherFromDb.FirstName, teacherFromDb.LastName, teacherFromDb.Email, teacherFromDb.Class, teacherFromDb.Subject, id)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Error updating teacher!", http.StatusInternalServerError)
			fmt.Println("Err : update failed!", err)
			return
		}

		log.Println("Patch teachers handler ended!")

	}

	// Apply the commit;
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Err : Commit failed!", http.StatusInternalServerError)
		fmt.Println("Err : Commit failed!", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// PatchTeacherHandler - Patches single teacher details based on the ID;
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Teacher not found!", http.StatusNotFound)
			fmt.Println("Err : Teacher not found", err)
			return
		}
		http.Error(w, "Query failed!", http.StatusInternalServerError)
		fmt.Println("Err : Teacher not found", err)
		return
	}

	// Apply updates using reflect;
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	fmt.Println("Teacher value field @ (0) := ", teacherVal.Type().Field(0))
	fmt.Println("Teacher value field @ (1) := ", teacherVal.Type().Field(1))
	fmt.Println("Teacher value field @ (2) := ", teacherVal.Type().Field(2))

	teacherType := teacherVal.Type() // this will store the teacher type (models.Teacher)
	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			fmt.Println(field.Tag.Get("json"), field.Name)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, class = ?, subject = ?, email = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Class, existingTeacher.Subject, existingTeacher.Email, existingTeacher.Id)
	if err != nil {
		http.Error(w, "Update failed!", http.StatusInternalServerError)
		fmt.Println("Err : Update failed", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Delete failed", err)
		return
	}

	fmt.Println(res.RowsAffected())
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Err retrieving deleted row!", http.StatusInternalServerError)
		fmt.Println("Err : RowsAffected failed", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Teacher not found!", http.StatusNotFound)
		fmt.Println("Err : Teacher not found", err)
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
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	var ids []int
	err = json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		http.Error(w, "Incorrect request!", http.StatusInternalServerError)
		fmt.Println("Err : Invalid Teacher ID", err)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Begin failed", err)
		return
	}

	statement, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		tx.Rollback()
		fmt.Println("Err : Prepare failed", err)
		return
	}

	defer statement.Close()

	deletedIds := []int{}
	for _, id := range ids {
		res, err := statement.Exec(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			tx.Rollback()
			fmt.Println("Err : Delete failed", err)
			return
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			tx.Rollback()
			fmt.Println("Err : RowsAffected failed", err)
			return
		}

		// If teacher was deleted, then push the id to deletedIds[];
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Commit failed", err)
		return
	}

	if len(deletedIds) == 0 {

		tx.Rollback()
		http.Error(w, fmt.Sprintf("ID(s) %d does not exist", ids), http.StatusNotFound)
		fmt.Println("Err : Teachers not found", err)
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
