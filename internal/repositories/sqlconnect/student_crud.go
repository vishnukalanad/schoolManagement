package sqlconnect

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"schoolManagement/internal/models"
	"schoolManagement/pkg/utils"
)

// ******** DB Crud Handlers ********

// GetStudentDbHandler - Fetches students list from DB;
func GetStudentDbHandler(r *http.Request) (error, []models.Student) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), []models.Student{}
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var students []models.Student
	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	var args []interface{}

	query, args = utils.GetFilters(r, query, args)
	query = utils.SortQueryParams(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return utils.HandleError(err, "Err: Query execution failed!"), []models.Student{}
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			return
		}
	}()

	for rows.Next() {
		var student models.Student
		err = rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
		if err != nil {
			return utils.HandleError(err, "Err: Data retrieval failed!"), []models.Student{}
		}

		students = append(students, student)
	}

	return nil, students
}

// AddStudentsDbHandler - Handles the crud operation to store new student details in table;
func AddStudentsDbHandler(r *http.Request) (error, []models.Student) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), []models.Student{}
	}

	var studentRaw []map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot read request body!"), nil
	}

	defer func() {
		err := r.Body.Close()
		if err != nil {
			return
		}
	}()

	// Query statement prepare;
	statement, err := db.Prepare(utils.GetInsertQuery(models.Student{}))
	if err != nil {
		return utils.HandleError(err, "Err: Cannot prepare statement!"), nil
	}

	// Closing the statement;
	defer func() {
		err := statement.Close()
		if err != nil {
			return
		}
	}()

	log.Println("\nUnmarshalling body")

	err = json.Unmarshal(body, &studentRaw)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot parse request body!"), nil
	}

	log.Println("\nGenerating keys for struct")
	// Handling extra fields passed in request;
	keys := utils.GetFieldNames(models.Student{})

	// Creating a map of allowed keys;
	validKeys := make(map[string]struct{})
	for _, key := range keys {
		validKeys[key] = struct{}{}
	}

	log.Println("\nValidating request body", keys)
	for _, student := range studentRaw {
		for key := range student {
			_, ok := validKeys[key]
			if !ok {
				return utils.HandleError(err, "Err: Internal server error!"), nil
			}
		}
	}

	fmt.Println("Received student details : ", string(body))

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var students []models.Student
	err = json.Unmarshal(body, &students)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot parse request body!"), nil
	}

	fmt.Println("Processed student details : ", students)

	log.Println("\nValidating empty values in request")
	// Validation for empty values in request body;
	for _, student := range students {
		values := reflect.ValueOf(student)
		for i := 0; i < values.NumField(); i++ {
			val := values.Field(i)
			if val.Kind() == reflect.String && val.String() == "" {
				return utils.HandleError(err, "Err: Cannot parse request body!"), nil
			}

		}
	}
	log.Println("\nStatement execution begins")
	// Loops through the incoming students arrays and executes the insert statement for store the values in DB;
	for _, student := range students {
		values := utils.GetFieldValues(student)
		log.Println("\nField values", values)

		res, err := statement.Exec(values...)
		if err != nil {
			return utils.HandleError(err, "Err: Cannot add student to database!"), nil
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return utils.HandleError(err, "Err: Cannot add student to database!"), nil
		}

		student.Id = int(lastId)
	}

	// Returns the final students array;
	return nil, students
}

// UpdateStudentsDbHandler - Handles the update operation of students;
func UpdateStudentsDbHandler(id int, updatedStudent models.Student) (error, []models.Student) {
	// Connect to DB;
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), []models.Student{}
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var student models.Student
	// Fetch student details based on ID;
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No student found!"), []models.Student{}
		}
		return utils.HandleError(err, "Err: Cannot get student from db!"), nil
	}

	// Execute the update query;
	updatedStudent.Id = int(student.Id)
	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", updatedStudent.FirstName, updatedStudent.LastName, updatedStudent.Email, updatedStudent.Class, student.Id)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot update student in db!"), nil
	}
	return nil, []models.Student{updatedStudent}
}

func PatchStudentsDbHandler(students []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!")
	}

	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("DB Close failed!")
			return
		}
	}()

	log.Println("\nStarting patch students handler")
	tx, err := db.Begin()
	if err != nil {
		return utils.HandleError(err, "Err: Cannot begin transaction!")
	}

	for _, student := range students {
		id := fmt.Sprintf("%v", student["id"])
		log.Println("\nStudent ID: ", id)

		var studentFromDb models.Student
		err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&studentFromDb.Id, &studentFromDb.FirstName, &studentFromDb.LastName, &studentFromDb.Email, &studentFromDb.Class)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return utils.HandleError(err, "Err: No student found!!")
			}
			return utils.HandleError(err, "Err: Cannot get student from db!")
		}

		studentVal := reflect.ValueOf(&studentFromDb).Elem()
		studentType := studentVal.Type()

		for k, v := range student {
			if k == "id" {
				continue // skips updating ID field;
			}

			// Looping through fields of student model;
			for i := 0; i < studentVal.NumField(); i++ {
				field := studentType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := studentVal.Field(i)

					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							err = tx.Rollback()
							if err != nil {
								fmt.Println("Rollback failed!")
								return utils.HandleError(err, "Err: Cannot rollback transaction!")
							}

							break
						}
					}
				}
			}

			_, err = tx.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", studentFromDb.FirstName, studentFromDb.LastName, studentFromDb.Email, studentFromDb.Class, id)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, sql.ErrNoRows) {
					return utils.HandleError(err, "Err: No student found!!")
				}
				return utils.HandleError(err, "Err: Cannot update student in db!")
			}
		}

		err = tx.Commit()
		if err != nil {
			return utils.HandleError(err, "Err: Cannot commit transaction!")
		}
	}

	return nil
}

func PatchStudentDbHandler(id int, updatedStudent map[string]interface{}) (error, models.Student) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), models.Student{}
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var student models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).Scan(&student.Id, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No student found!!"), models.Student{}
		}
		return utils.HandleError(err, "Err: Cannot get student from db!"), models.Student{}
	}

	studentVal := reflect.ValueOf(&student).Elem()
	studentType := studentVal.Type()

	for k, v := range updatedStudent {
		for i := 0; i < studentVal.NumField(); i++ {
			field := studentType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					studentVal.Field(i).Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?", student.FirstName, student.LastName, student.Email, student.Class, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No student found!!"), models.Student{}
		}
		return utils.HandleError(err, "Err: Cannot update student in db!"), models.Student{}
	}

	return nil, student
}

func DeleteStudentsDbHandler(ids []int) (error, []int) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), nil
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), nil
	}

	statement, err := tx.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		func() {
			err := tx.Rollback()
			if err != nil {
				fmt.Println("Rollback failed!")
				return
			}
		}()
		return utils.HandleError(err, "Err: Internal server error!"), nil
	}
	defer func() {
		err := statement.Close()
		if err != nil {
			return
		}
	}()

	deletedIds := []int{}
	for _, id := range ids {
		res, err := statement.Exec(id)
		if err != nil {
			func() {
				err := tx.Rollback()
				if err != nil {
					fmt.Println("Rollback failed!")
					return
				}
			}()
			if errors.Is(err, sql.ErrNoRows) {
				return utils.HandleError(err, "Err: No student found!!"), deletedIds
			}
			return utils.HandleError(err, "Err: Cannot delete student from db!"), deletedIds
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			func() {
				err := tx.Rollback()
				if err != nil {
					fmt.Println("Rollback failed!")
					return
				}
			}()
			if errors.Is(err, sql.ErrNoRows) {
				return utils.HandleError(err, "Err: No student found!!"), deletedIds
			}
			return utils.HandleError(err, "Err: Failed to get affected rows count!"), deletedIds
		}

		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}
	}

	err = tx.Commit()
	if err != nil {
		return utils.HandleError(err, "Err: Cannot commit transaction!"), deletedIds
	}

	if len(deletedIds) == 0 {
		func() {
			err := tx.Rollback()
			if err != nil {
				fmt.Println("Rollback failed!")
				return
			}
		}()
	}

	return err, deletedIds
}

func DeleteStudentDbHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!")
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	res, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Err: No student found!")
			return utils.HandleError(err, "Err: No student found!")
		}
		return utils.HandleError(err, "Err: Cannot delete student from db!")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Err: Failed to get affected rows count!")
		return utils.HandleError(err, "Err: Failed to get affected rows count!")
	}

	if rowsAffected == 0 {
		fmt.Println("Err: No student found!")
		return utils.HandleError(err, "Err: No student found!")
	}

	return err
}
