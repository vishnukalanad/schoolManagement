package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"schoolManagement/internal/models"
	"schoolManagement/internal/repositories/sqlconnect"
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

// getTeachersHandler - this will handle the business logic for get teachers;
func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
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

// addTeachersHandler - handles the incoming post requests;
func addTeachersHandler(w http.ResponseWriter, r *http.Request) {
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
