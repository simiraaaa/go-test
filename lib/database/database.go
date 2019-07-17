package database

import (
	"database/sql"
	"fmt"
	"os"
)

func Open() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/test", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASS")))
	if err != nil {
		panic(err.Error())
	}
	return db
}
