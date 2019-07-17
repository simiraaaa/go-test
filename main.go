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
	// database module化したけど、 Open しかないので今の所あまりmoduleの意味がない
	db := database.Open()
	log.Println("connect")
	// TODO: プログラム (サーバー ?) が終了したときに close したい
	defer db.Close()
	log.Println("select")
	// 新しい順に10件取得する特に意味のない処理
	rows, err := db.Query("select * from users order by created_at desc limit 10")
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
