package database

import "net/http"

// CRUD http のメソッドによって呼び出し方を変更する。実装されていない場合は 404。 POST: create, GET: read, PUT: update, DELETE: delete
type CRUD struct {
	Create,
	Read,
	Update,
	Delete func(http.ResponseWriter, *http.Request)
}

func CallCRUD(crud CRUD, w http.ResponseWriter, r *http.Request) {
	var f func(http.ResponseWriter, *http.Request)
	switch r.Method {
	case http.MethodGet:
		f = crud.Read
	case http.MethodPost:
		f = crud.Create
	case http.MethodPut:
		f = crud.Update
	case http.MethodDelete:
		f = crud.Delete
	}
	if f == nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		f(w, r)
	}
}
