package sqlconnect

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"schoolManagement/internal/models"
	"schoolManagement/pkg/utils"
	"strings"
)

// ******** General Helper Functions ********

// getFilters - Gets the filters from query params;
func getFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"firstname": "first_name",
		"lastname":  "last_name",
		"class":     "class",
		"email":     "email",
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

// sortQueryParams - Looks for any sortBy params in the request and updates the query string;
func sortQueryParams(r *http.Request, query string) string {

	// Takes the sortBy query param from request;
	sortParams := r.URL.Query()["sortBy"]
	if len(sortParams) > 0 {
		// Adds order by to sql query;
		query += " order by"

		// Loops through the sort params;
		for i, val := range sortParams {

			// Accepts the sortBy params as param:order format;
			// Splits based on the : and extracts the key value pairs;
			parts := strings.Split(val, ":")

			// Skips the iteration if no sort order provided;
			if len(parts) != 2 {
				continue
			}

			// Stores the field name and sort order in 2 variables;
			field, order := parts[0], parts[1]

			// Validates to see if the provided sort and field values are valid;
			if !isValidSortType(order) || !isFieldValid(field) {
				continue
			}

			// Updates the query string appropriately;
			if i > 0 {
				query += ", "
			}

			query += " " + field + " " + order
		}
	}
	// Return the final query string;
	return query
}

// isValidSortType - Validates the sort order;
func isValidSortType(order string) bool {
	return order == "asc" || order == "desc"
}

// isFieldValid - Validates the query fields;
func isFieldValid(field string) bool {
	fields := map[string]bool{
		"firstname": true,
		"lastname":  true,
		"class":     true,
		"email":     true,
	}

	return fields[field]
}

// getInsertQuery - Generates the insert query for students;
func getInsertQuery(model interface{}) string {
	types := reflect.TypeOf(model)
	var cols, placeholders string
	for i := 0; i < types.NumField(); i++ {
		dbTag := types.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id" {
			if cols != "" {
				cols += ", "
				placeholders += ", "
			}
			cols += dbTag
			placeholders += "?"
		}
	}

	fmt.Printf("Generated query : insert into students (%s) values (%s)", cols, placeholders)
	return fmt.Sprintf("INSERT INTO students (%s) VALUES (%s)", cols, placeholders)
}

// getFieldNames - Return the list of fields values based on the struct passed;
func getFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd)
	}

	return fields
}

// getFieldValues - Returns the field values;
func getFieldValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	var values []interface{}

	// Loops through the modelTypes and extracts the "db" tag values from struct;
	// Then stores the value index i to values array;
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id,omitempty" {
			values = append(values, modelValue.Field(i).Interface())
		}
	}

	// Returns the final values array;
	return values
}

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
	query := "select id, first_name, last_name, email, class from students where 1=1"
	var args []interface{}

	query, args = getFilters(r, query, args)
	query = sortQueryParams(r, query)

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
		err = rows.Scan(&student.Id, student.FirstName, &student.LastName, &student.Email, &student.Class)
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
	statement, err := db.Prepare(getInsertQuery(models.Student{}))
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
	keys := getFieldNames(models.Student{})

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
		values := getFieldValues(student)
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
