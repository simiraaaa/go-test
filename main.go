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

// CRUD http のメソッドによって呼び出し方を変更する。実装されていない場合は 404。 POST: create, GET: read, PUT: update, DELETE: delete
type CRUD struct {
	create,
	read,
	update,
	delete func(w http.ResponseWriter, r *http.Request)
}

func callCRUD(crud CRUD, w http.ResponseWriter, r *http.Request) {
	var f func(w http.ResponseWriter, r *http.Request)
	switch r.Method {
	case http.MethodGet:
		f = crud.read
	case http.MethodPost:
		f = crud.create
	case http.MethodPut:
		f = crud.update
	case http.MethodDelete:
		f = crud.delete
	}
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		f(w, r)
	}
}

// User は DB から users テーブルのデータを受け取ります
type User struct {
	ID        int
	Name      string
	CreatedAt int
	UpdatedAt int
}

var db *sql.DB

func main() {
	// database module化したけど、 Open しかないので今の所あまりmoduleの意味がない
	db = database.Open()
	log.Println("connect")
	http.HandleFunc("/users", handler) // ハンドラを登録してウェブページを表示させる
	log.Println("http://localhost:8085")
	http.ListenAndServe(":8085", nil)
}

var usersCRUD = CRUD{
	create: func(w http.ResponseWriter, r *http.Request) {
		dec := json.NewDecoder(r.Body)
		var user User
		err := dec.Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Println(user)
		// validation はどうやるのが最適?
		if user.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "name is required\n")
			return
		}

		now := strconv.FormatInt(time.Now().Unix(), 10)

		insert, err := db.Query("insert into users (name, created_at, updated_at) values ('" + user.Name + "', " + now + "," + now + ")")

		log.Println("insert")

		if err != nil {
			panic(err.Error())
		}

		defer insert.Close()
	},
	read: func(w http.ResponseWriter, r *http.Request) {

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
	},
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("access /")
	log.Println(r.Method)
	// `show status like 'Threads_connected';` で確認したらプログラム終了時にコネクションが切れているのであえて close しなくても良さそう
	// defer db.Close()
	callCRUD(usersCRUD, w, r)
}
