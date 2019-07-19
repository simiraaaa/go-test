package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var db *sql.DB

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("access /")
	log.Println(r.Method)
	if r.Method == http.MethodGet {
		getUsers(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	http.HandleFunc("/users", handler) // ハンドラを登録してウェブページを表示させる
	log.Println("http://localhost:8085")
	http.ListenAndServe(":8085", nil)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// database module化したけど、 Open しかないので今の所あまりmoduleの意味がない
	db = database.Open()
	log.Println("connect")
	// TODO: プログラム (サーバー ?) が終了したときに close したい
	defer db.Close()
	log.Println("select")
	const limit = 10
	// 新しい順に10件取得する
	rows, err := db.Query("select * from users order by created_at desc limit " + strconv.Itoa(limit))
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	outputs := [limit]User{}
	for i := 0; rows.Next(); i++ {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			panic(err.Error())
		}
		outputs[i] = user
		// fmt.Fprintf(w, "%d, %s, %d, %d\n", user.ID, user.Name, user.CreatedAt, user.UpdatedAt)
	}
	json, err := json.Marshal(outputs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	fmt.Fprintf(w, "%s", json)
}

func insertUser() {
	db = database.Open()
	log.Println("connect")
	// TODO: プログラム (サーバー ?) が終了したときに close したい
	defer db.Close()
	now := strconv.FormatInt(time.Now().Unix(), 10)

	insert, err := db.Query("insert into users (name, created_at, updated_at) values ('TEST_USER_" + now + "', " + now + "," + now + ")")

	log.Println("insert")

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}
