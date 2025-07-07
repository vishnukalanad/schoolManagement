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

func GetTeachersDbHandler(r *http.Request, teachersList []models.Teacher) (error, []models.Teacher) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return utils.HandleError(err, "Err : DB connection failed!"), nil
	}
	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
		}
	}()

	query := "SELECT id, first_name, last_name, class, subject, email  FROM teachers WHERE 1=1"
	var args []interface{}

	query, args = addFilters(r, query, args)
	query = sortByQueryParams(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error at query execution : ", err)
		return utils.HandleError(err, "Err : DB connection failed!"), nil
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()

	for rows.Next() {
		var teacher models.Teacher
		err = rows.Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject, &teacher.Email)
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Error", err)
			return utils.HandleError(err, "Err : Internal server error!"), nil
		}
		teachersList = append(teachersList, teacher)
	}
	return err, teachersList
}

func GetTeacherDbHandler(idStr string) (error, models.Teacher) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return utils.HandleError(err, "Err : DB connection failed!"), models.Teacher{}
	}
	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
		}
	}()
	var teacher models.Teacher
	// Handling param based query;
	err = db.QueryRow("SELECT id, first_name, last_name, class, subject, email  FROM teachers WHERE id=?", idStr).Scan(&teacher.Id, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject, &teacher.Email)
	if errors.Is(err, sql.ErrNoRows) {
		//http.Error(w, "Do records not found", http.StatusNotFound)
		fmt.Println("Error", err)
		return utils.HandleError(err, "Err : DB records not found!"), models.Teacher{}
	} else if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Error", err)
		return utils.HandleError(err, "Err : Internal server error!"), models.Teacher{}
	}
	return err, teacher
}

func AddTeacherDbHandler(w http.ResponseWriter, r *http.Request, newTeachers []models.Teacher) (error, []models.Teacher) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return utils.HandleError(err, "Err : DB connection failed!"), nil
	}
	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
			return
		}
	}()

	var rawTeachers []map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error at reading body", err)
		return utils.HandleError(err, "Error at reading body!"), nil
	}

	defer r.Body.Close()

	//stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?,?,?,?,?)")
	stmt, err := db.Prepare(generateInsertQuery(models.Teacher{}))
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Insert failed", err)
		return utils.HandleError(err, "Err : Insert failed!"), nil
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("2:Close failed", err)
			return
		}
	}()

	//err = json.NewDecoder(r.Body).Decode(&rawTeachers)
	err = json.Unmarshal(body, &rawTeachers)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Invalid Request", err)
		return utils.HandleError(err, "Err : Invalid Request!"), nil
	}

	// Handling any unwanted additional fields sent in request body;
	// Send invalid request body response in such cases to block unwanted stuff;
	fields := GetFieldNames(models.Teacher{})

	// Make a map of allowed tags;
	allowedTags := make(map[string]struct{})
	for _, field := range fields {
		allowedTags[field] = struct{}{}
	}

	for _, teacher := range rawTeachers {
		for key := range teacher {
			_, ok := allowedTags[key]
			if !ok {
				http.Error(w, "Invalid request body!", http.StatusBadRequest)
				return err, nil
			}
		}
	}

	//err = json.NewDecoder(r.Body).Decode(&newTeachers)
	err = json.Unmarshal(body, &newTeachers)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Invalid request", err)
		return utils.HandleError(err, "Err : Invalid request!"), nil
	}

	fmt.Println(newTeachers)
	// This is the handle the validation for empty values passed;
	for _, teacher := range newTeachers {
		//if teacher.FirstName == "" || teacher.LastName == "" || teacher.Class == "" || teacher.Subject == "" || teacher.Email == "" {
		//	http.Error(w, "All fields are required", http.StatusBadRequest)
		//	return err, nil
		//}

		val := reflect.ValueOf(teacher)
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				fmt.Println("field.Kind() : ", field.Kind())
				fmt.Println("reflect.String : ", reflect.String)
				fmt.Println("field.String()", field.String())
				http.Error(w, "All fields are required", http.StatusBadRequest)
				return err, nil
			}
		}
	}

	addedTeachers := make([]models.Teacher, len(newTeachers))

	for i, teacher := range newTeachers {
		//res, err := stmt.Exec(teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject)
		values := getStructValues(teacher)
		res, err := stmt.Exec(values...)
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Err : Data insertion to DB failed!", err)
			return utils.HandleError(err, "Err : Data insertion to DB failed!"), nil
		}
		lastId, err := res.LastInsertId()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Err : Get ID from DB failed!", err)
			return utils.HandleError(err, "Err : Get ID from DB failed!"), nil
		}

		teacher.Id = int(lastId)
		addedTeachers[i] = teacher

	}
	return err, addedTeachers
}

func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd) // Getting json tags
	}
	return fields
}

func generateInsertQuery(model interface{}) string {
	modelType := reflect.TypeOf(model)
	var columns, placeholders string
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		//fmt.Println(dbTag)
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		if dbTag != "" && dbTag != "id" {
			if columns != "" {
				columns += ", "
				placeholders += ", "
			}
			columns += dbTag
			placeholders += "?"
		}
	}

	fmt.Printf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
	return fmt.Sprintf("INSERT INTO teachers (%s) VALUES (%s)", columns, placeholders)
}

func getStructValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	var values []interface{}

	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		if dbTag != "" && dbTag != "id,omitempty" {
			//fmt.Println("\n", dbTag, modelValue, modelValue.Field(i), modelValue.Field(i).Interface())
			values = append(values, modelValue.Field(i).Interface())
		}
	}
	fmt.Println("\nValues : ", values)
	return values
}

func UpdateTeachersDbHandler(id int, updatedTeachers models.Teacher) error {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Internal server error!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//http.Error(w, "Teacher not found!", http.StatusNotFound)
			fmt.Println("Err : Teacher not found", err)
			return utils.HandleError(err, "Err : Teacher not found")
		}
		//http.Error(w, "Query failed!", http.StatusInternalServerError)
		fmt.Println("Err : Teacher not found", err)
		return utils.HandleError(err, "Err : Teacher not found")
	}

	updatedTeachers.Id = existingTeacher.Id
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, class = ?, subject = ?, email = ? WHERE id = ?", updatedTeachers.FirstName, updatedTeachers.LastName, updatedTeachers.Class, updatedTeachers.Subject, updatedTeachers.Email, updatedTeachers.Id)
	if err != nil {
		//http.Error(w, "Update failed!", http.StatusInternalServerError)
		fmt.Println("Err : Update failed", err)
		return utils.HandleError(err, "Err : Update failed")
	}
	return utils.HandleError(err, "Internal server error!")
}

func PatchTeachersDbHandler(updates []map[string]interface{}) error {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	// DB Transaction Beginning;
	tx, err := db.Begin()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB transaction failed!", err)
		return err
	}

	for _, update := range updates {
		id := fmt.Sprintf("%v", update["id"])
		log.Println("TEACHER ID : ", id)

		var teacherFromDb models.Teacher
		err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacherFromDb.Id,
			&teacherFromDb.FirstName, &teacherFromDb.LastName, &teacherFromDb.Email, &teacherFromDb.Class, &teacherFromDb.Subject)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				//http.Error(w, "Teacher not found!", http.StatusNotFound)
				fmt.Println("Err : Teacher not found", err)
				return err
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
								//http.Error(w, "Invalid field value!", http.StatusBadRequest)
								fmt.Println("Err : Invalid field value", err)
								return err
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
			//http.Error(w, "Error updating teacher!", http.StatusInternalServerError)
			fmt.Println("Err : update failed!", err)
			return err
		}

	}

	// Apply the commit;
	err = tx.Commit()
	if err != nil {
		//http.Error(w, "Err : Commit failed!", http.StatusInternalServerError)
		fmt.Println("Err : Commit failed!", err)
		return err
	}
	return nil
}

func PatchTeacherDbHandler(id int, updates map[string]interface{}) (error, models.Teacher) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.Id, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//http.Error(w, "Teacher not found!", http.StatusNotFound)
			fmt.Println("Err : Teacher not found", err)
			return nil, models.Teacher{}
		}
		//http.Error(w, "Query failed!", http.StatusInternalServerError)
		fmt.Println("Err : Teacher not found", err)
		return nil, models.Teacher{}
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
		//http.Error(w, "Update failed!", http.StatusInternalServerError)
		fmt.Println("Err : Update failed", err)
		return nil, models.Teacher{}
	}
	return err, existingTeacher
}

func DeleteTeacherDbHandler(id int) error {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
		}
	}()

	res, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Delete failed", err)
		return err
	}

	fmt.Println(res.RowsAffected())
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		//http.Error(w, "Err retrieving deleted row!", http.StatusInternalServerError)
		fmt.Println("Err : RowsAffected failed", err)
		return err
	}

	if rowsAffected == 0 {
		//http.Error(w, "Teacher not found!", http.StatusNotFound)
		fmt.Println("Err : Teacher not found", err)
	}
	return nil
}

func DeleteTeachersDbHandler(ids []int) (error, []int) {
	db, err := ConnectDb()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : DB connection failed!", err)
		return utils.HandleError(err, "Err : DB connection failed!"), nil
	}

	defer func() {
		err := db.Close()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Println("Close failed", err)
			return
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Begin failed", err)
		return utils.HandleError(err, "Err : Query failed!"), nil
	}

	statement, err := tx.Prepare("DELETE FROM teachers WHERE id = ?")
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		tx.Rollback()
		fmt.Println("Err : Prepare failed", err)
		return utils.HandleError(err, "Err : Query prepare failed!"), nil
	}

	defer statement.Close()

	deletedIds := []int{}
	for _, id := range ids {
		res, err := statement.Exec(id)
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			tx.Rollback()
			fmt.Println("Err : Delete failed", err)
			return utils.HandleError(err, "Err : Delete failed"), nil
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			tx.Rollback()
			fmt.Println("Err : RowsAffected failed", err)
			return utils.HandleError(err, "Err : RowsAffected failed"), nil
		}

		// If teacher was deleted, then push the id to deletedIds[];
		if rowsAffected > 0 {
			deletedIds = append(deletedIds, id)
		}

	}

	err = tx.Commit()
	if err != nil {
		//http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		fmt.Println("Err : Commit failed", err)
		return utils.HandleError(err, "Err : Commit failed"), nil
	}

	if len(deletedIds) == 0 {

		tx.Rollback()
		//http.Error(w, fmt.Sprintf("ID(s) %d does not exist", ids), http.StatusNotFound)
		fmt.Println("Err : Teachers not found", err)
		return utils.HandleError(err, "Err : Teachers not found"), nil

	}
	return err, deletedIds
}
