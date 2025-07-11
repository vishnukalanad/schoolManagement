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
	"strconv"
	"time"
)

func GetExecsDbHandler(r *http.Request) (error, []models.Exec) {
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

	var execs []models.Exec
	query := "SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE 1=1"
	var args []interface{}

	query, args = utils.GetFilters(r, query, args)
	query = utils.SortQueryParams(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		return utils.HandleError(err, "Err: Query execution failed!"), nil
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			return
		}
	}()

	for rows.Next() {
		var exec models.Exec
		err = rows.Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.CreatedAt, &exec.Inactive, &exec.Role)
		if err != nil {
			return utils.HandleError(err, "Err: Data retrieval failed!"), nil
		}

		execs = append(execs, exec)
	}

	return nil, execs
}

func AddExecsDbHandler(r *http.Request) (error, []models.Exec) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), nil
	}

	var execRaw []map[string]interface{}
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
	statement, err := db.Prepare(utils.GetExecInsertQuery(models.Exec{}))
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

	err = json.Unmarshal(body, &execRaw)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot parse request body!"), nil
	}

	log.Println("\nGenerating keys for struct")

	fmt.Println("Received exec details : ", string(body))

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var execs []models.Exec
	err = json.Unmarshal(body, &execs)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot parse request body!"), nil
	}

	fmt.Println("Processed exec details : ", execs)

	log.Println("\nValidating empty values in request")
	// Validation for empty values in request body;
	for _, exec := range execs {
		if exec.Password == "" {
			return utils.HandleError(errors.New("empty password"), "Err: Please provide a valid password!"), nil
		}

		values := reflect.ValueOf(exec)
		for i := 0; i < values.NumField(); i++ {
			val := values.Field(i)
			if val.Kind() == reflect.String && val.String() == "" {
				return utils.HandleError(err, "Err: Cannot parse request body!"), nil
			}

		}
	}
	log.Println("\nStatement execution begins", execs)
	// Loops through the incoming execs arrays and executes the insert statement for store the values in DB;
	for _, exec := range execs {
		log.Println("\nStarting password hashing")

		encodedHash, err := utils.HashPassword(exec.Password)
		if err != nil {
			return utils.HandleError(err, "Err: Cannot hash password!"), nil
		}

		log.Println("\nPassword : ", exec.Password, encodedHash)

		values := utils.GetFieldValues(exec)
		log.Println("\nField values", values)

		res, err := statement.Exec(values...)
		if err != nil {
			return utils.HandleError(err, "Err: Cannot add exec to database!"), nil
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return utils.HandleError(err, "Err: Cannot add exec to database!"), nil
		}

		exec.Id = int(lastId)
	}

	// Returns the final execs array;
	return nil, execs
}

func PatchExecsDbHandler(execs []map[string]interface{}) error {
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

	log.Println("\nStarting patch execs handler")
	tx, err := db.Begin()
	if err != nil {
		return utils.HandleError(err, "Err: Cannot begin transaction!")
	}

	for _, exec := range execs {
		id := fmt.Sprintf("%v", exec["id"])
		log.Println("\nExec ID: ", id)

		var execsFromDb models.Exec
		err = db.QueryRow("SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?", id).Scan(&execsFromDb.Id, &execsFromDb.FirstName, &execsFromDb.LastName, &execsFromDb.Email, &execsFromDb.Username, &execsFromDb.CreatedAt, &execsFromDb.Inactive, &execsFromDb.Role)
		if err != nil {
			tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return utils.HandleError(err, "Err: No exec found!!")
			}
			return utils.HandleError(err, "Err: Cannot get exec from db!")
		}

		execVal := reflect.ValueOf(&execsFromDb).Elem()
		execType := execVal.Type()

		for k, v := range exec {
			if k == "id" {
				continue // skips updating ID field;
			}

			// Looping through fields of exec model;
			for i := 0; i < execVal.NumField(); i++ {
				field := execType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := execVal.Field(i)

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

			_, err = tx.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, inactive_status = ?, role = ? WHERE id = ?", execsFromDb.FirstName, execsFromDb.LastName, execsFromDb.Email, execsFromDb.Username, execsFromDb.Inactive, execsFromDb.Role, id)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, sql.ErrNoRows) {
					return utils.HandleError(err, "Err: No exec found!!")
				}
				return utils.HandleError(err, "Err: Cannot update exec in db!")
			}
		}

	}
	err = tx.Commit()
	if err != nil {
		return utils.HandleError(err, "Err: Cannot commit transaction!")
	}

	return nil
}

func GetExecsByIdHandler(id int) (error, models.Exec) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), models.Exec{}
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM exec WHERE id = ?", id).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.CreatedAt, &exec.Inactive, &exec.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!"), models.Exec{}
		}
		return utils.HandleError(err, "Err: Data retrieval failed!"), models.Exec{}
	}
	return nil, exec
}

func PatchExecByIdDbHandler(id int, updatedStudent map[string]interface{}) (error, models.Exec) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), models.Exec{}
	}

	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	var exec models.Exec
	err = db.QueryRow("SELECT id, first_name, last_name, email, username, user_created_at, inactive_status, role FROM execs WHERE id = ?", id).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.CreatedAt, &exec.Inactive, &exec.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No exec found!!"), models.Exec{}
		}
		return utils.HandleError(err, "Err: Cannot get exec from db!"), models.Exec{}
	}

	execVal := reflect.ValueOf(&exec).Elem()
	execType := execVal.Type()

	for k, v := range updatedStudent {
		for i := 0; i < execVal.NumField(); i++ {
			field := execType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					execVal.Field(i).Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	_, err = db.Exec("UPDATE execs SET first_name = ?, last_name = ?, email = ?, username = ?, inactive_status = ?, role = ? WHERE id = ?", exec.FirstName, exec.LastName, exec.Email, exec.Username, exec.Inactive, exec.Role, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No exec found!!"), models.Exec{}
		}
		return utils.HandleError(err, "Err: Cannot update exec in db!"), models.Exec{}
	}

	return nil, exec
}

func DeleteExecByIdDbHandler(id int) error {
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

	res, err := db.Exec("DELETE FROM execs WHERE id = ?", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Err: No exec found!")
			return utils.HandleError(err, "Err: No exec found!")
		}
		return utils.HandleError(err, "Err: Cannot delete exec from db!")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Err: Failed to get affected rows count!")
		return utils.HandleError(err, "Err: Failed to get affected rows count!")
	}

	if rowsAffected == 0 {
		fmt.Println("Err: No exec found!")
		return utils.HandleError(err, "Err: No exec found!")
	}

	return err
}

func LoginDbHandler(username string, exec *models.Exec) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!")
	}

	defer db.Close()

	err = db.QueryRow("SELECT id, first_name, last_name, email, username, password, inactive_status, role FROM execs WHERE username = ?", username).Scan(&exec.Id, &exec.FirstName, &exec.LastName, &exec.Email, &exec.Username, &exec.Password, &exec.Inactive, &exec.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!!")
		}
		return utils.HandleError(err, "Err: Cannot get records from db!")
	}
	return nil
}

func UpdatePasswordDbHandler(userId int, request models.UpdatePasswordRequest) (error, string) {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!"), ""
	}

	var username, password, role string
	err = db.QueryRow("select username, password, role from execs where id = ?", userId).Scan(&username, &password, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!!"), ""
		}
		return utils.HandleError(err, "Err: Cannot get records from db!"), ""
	}

	err = utils.PasswordValidate(password, request.CurrentPassword)
	if err != nil {
		return utils.HandleError(err, "Err: Current password incorrect!"), ""
	}

	hashedPass, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot hash password!"), ""
	}

	currentTime := time.Now().Format(time.RFC3339)

	_, err = db.Exec("update execs set password = ?, password_changed_at = ? where id = ?", hashedPass, currentTime, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!!"), ""
		}
		return utils.HandleError(err, "Err: Cannot update records from db!"), ""
	}

	idStr := strconv.Itoa(userId)
	token, err := utils.SignToken(idStr, username, role)
	if err != nil {
		return utils.HandleError(err, "Err: Cannot sign token!"), ""
	}

	return nil, token
}

func ForgotPasswordDbHandler(email string, exec *models.Exec) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!")
	}

	defer func() {
		er := db.Close()
		if er != nil {
			return
		}
	}()

	err = db.QueryRow("select id from execs where email = ?", email).Scan(&exec.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!!")
		}
		return utils.HandleError(err, "Err: Cannot get records from db!")
	}
	return nil
}

func ForgotPasswordUpdateDbHandler(exec models.Exec, hashedToken string, expiry string) error {
	db, err := ConnectDb()
	if err != nil {
		return utils.HandleError(err, "Err: Internal server error!")
	}
	defer func() {
		er := db.Close()
		if er != nil {
			return
		}
	}()

	log.Println("Executing query", expiry, hashedToken, exec.Id)
	_, err = db.Exec("update execs set password_reset_expiry = ?, password_reset_token = ? where id = ?", expiry, hashedToken, exec.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.HandleError(err, "Err: No records found!!")
		}
		return utils.HandleError(err, "Err: Cannot update records from db!")
	}

	log.Println("DB update done for password reset expiry time and token!")

	return nil
}
