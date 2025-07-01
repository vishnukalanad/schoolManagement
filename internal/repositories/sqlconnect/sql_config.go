package sqlconnect

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func ConnectDb() (*sql.DB, error) {

	usr := os.Getenv("DB_USER")
	pwd := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	//conString := "root:root@tcp(127.0.0.1:3306)/" + dbName
	conString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", usr, pwd, host, port, dbName)
	db, err := sql.Open("mysql", conString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database!")
	return db, nil
}
