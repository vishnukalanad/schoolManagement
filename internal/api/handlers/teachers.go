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
	"strings"
)

var (
	teachers = make(map[int]models.Teacher)
)

func isValidSort(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidField(field string) bool {
	fields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}

	return fields[field]
}

// GetTeachersHandler - this will handle the business logic for get teachers;
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
		}
	}()

	// Extract path params;
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	teacherList := make([]models.Teacher, 0, len(teachers))

	var teacher models.Teacher
	if idStr != "" {
		// Handling param based query;
		err = db.QueryRow("SELECT id, first_name, last_name, class, subject, email  FROM teachers WHERE id=?", idStr).Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject, &teacher.Email)
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Do records not found", http.StatusNotFound)
			fmt.Println("Error", err)
			return
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error", err)
			return
		}

		teacherList = append(teacherList, teacher)

	} else {

		query := "SELECT id, first_name, last_name, class, subject, email  FROM teachers WHERE 1=1"

		var args []interface{}

		query, args = addFilters(r, query, args)

		query = sortByQueryParams(r, query)

		rows, err := db.Query(query, args...)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error at query execution : ", err)
			return
		}

		defer func() {
			err := rows.Close()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		for rows.Next() {
			var teacher models.Teacher
			err = rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject, &teacher.Email)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Println("Error", err)
				return
			}
			teacherList = append(teacherList, teacher)
		}
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "Success",
		Count:  len(teacherList),
		Data:   teacherList,
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
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
		}
	}()

	// Extract path params;
	idStr := r.PathValue("id")

	var teacher models.Teacher
	// Handling param based query;
	err = db.QueryRow("SELECT id, first_name, last_name, class, subject, email  FROM teachers WHERE id=?", idStr).Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject, &teacher.Email)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Do records not found", http.StatusNotFound)
		fmt.Println("Error", err)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error", err)
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

func sortByQueryParams(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortBy"]
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, val := range sortParams {
			parts := strings.Split(val, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]

			if !isValidSort(order) || !isValidField(field) {
				continue
			}

			// if more than one condition for sorting, then separate them by comma (,);
			if i > 0 {
				query += ","
			}

			query += " " + field + " " + order

		}
	}
	return query
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}

	for param, dbField := range params {
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

// AddTeachersHandler - handles the incoming post requests;
func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDb()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return
	}
	defer func() {
		err := db.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
		}
	}()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Insert failed", err)
		return
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
			return
		}
	}()

	var newTeachers []models.Teacher
	err = json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Insert failed", err)
		return
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, teacher := range newTeachers {
		res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Err : Data insertion to DB failed!", err)
			return
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Err : Get ID from DB failed!", err)
			return
		}

		teacher.Id = int(lastId)
		addedTeachers[i] = teacher

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

	updatedTeachers.Id = existingTeacher.Id
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, class = ?, subject = ?, email = ? WHERE id = ?", updatedTeachers.FirstName, updatedTeachers.LastName, updatedTeachers.Class, updatedTeachers.Subject, updatedTeachers.Email, updatedTeachers.Id)
	if err != nil {
		http.Error(w, "Update failed!", http.StatusInternalServerError)
		fmt.Println("Err : Update failed", err)
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
