package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/simiraaaa/go-test/lib/database"
)

// User は DB から users テーブルのデータを受け取ります
type User struct {
	ID        int
	Name      string
	CreatedAt int
	UpdatedAt int
}

func main() {
	db, err := database.Open()
	log.Println("connect")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	log.Println("select")
	rows, err := db.Query("select * from users")
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		fmt.Println(user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
	}

	now := strconv.FormatInt(time.Now().Unix(), 10)

	insert, err := db.Query("insert into users (name, created_at, updated_at) values ('TEST_USER_" + now + "', " + now + "," + now + ")")

	log.Println("insert")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}
