package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

// GetFilters - Gets the filters from query params;
func GetFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"class":      "class",
		"email":      "email",
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

// SortQueryParams - Looks for any sortBy params in the request and updates the query string;
func SortQueryParams(r *http.Request, query string) string {

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

// GetInsertQuery - Generates the insert query for students;
func GetInsertQuery(model interface{}) string {
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

// GetExecInsertQuery - Generates the insert query for execs;
func GetExecInsertQuery(model interface{}) string {
	types := reflect.TypeOf(model)
	nullSqlType := reflect.TypeOf(sql.NullString{})
	var cols, placeholders string
	for i := 0; i < types.NumField(); i++ {
		dbTag := types.Field(i).Tag.Get("db")
		dbTag = strings.TrimSuffix(dbTag, ",omitempty")
		fieldType := types.Field(i).Type
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id" && fieldType != nullSqlType {
			if cols != "" {
				cols += ", "
				placeholders += ", "
			}
			cols += dbTag
			placeholders += "?"
		}
	}

	fmt.Printf("Generated query : insert into students (%s) values (%s)", cols, placeholders)
	return fmt.Sprintf("INSERT INTO execs (%s) VALUES (%s)", cols, placeholders)
}

// GetFieldNames - Return the list of fields values based on the struct passed;
func GetFieldNames(model interface{}) []string {
	val := reflect.TypeOf(model)
	var fields []string
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd)
	}

	return fields
}

// GetFieldValues - Returns the field values;
func GetFieldValues(model interface{}) []interface{} {
	modelValue := reflect.ValueOf(model)
	modelType := reflect.TypeOf(model)

	var values []interface{}

	nullSqlType := reflect.TypeOf(sql.NullString{})
	// Loops through the modelTypes and extracts the "db" tag values from struct;
	// Then stores the value index i to values array;
	for i := 0; i < modelType.NumField(); i++ {
		dbTag := modelType.Field(i).Tag.Get("db")
		fieldType := modelType.Field(i).Type
		fmt.Println(dbTag)
		if dbTag != "" && dbTag != "id,omitempty" && fieldType != nullSqlType {
			values = append(values, modelValue.Field(i).Interface())
		}
	}

	// Returns the final values array;
	return values
}
