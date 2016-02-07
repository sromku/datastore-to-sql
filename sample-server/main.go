package sromku

import (
	"net/http"
	"github.com/gorilla/mux"
	"appengine"
	"appengine/datastore"
	"fmt"
)

func init() {
	r := mux.NewRouter()
	r.HandleFunc("/", handle).Methods("GET")
	http.Handle("/", r)
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	putUser("Alban Lestat", "alban.lestat@gmail.com", 0, ctx)
	putUser("Erik Erich", "erik.erich@gmail.com", 0, ctx)
	putUser("Aindrea Tim", "aindrea.tim@gmail.com", 1, ctx)
	putUser("Luisita Karolina", "luisita.karolina@gmail.com", 1, ctx)
	fmt.Fprintln(w, "4 users saved in datastore")
}

func putUser(name, email string, gender int, ctx appengine.Context) {
	key := datastore.NewKey(ctx, "Profile", email, 0, nil)
	datastore.Put(ctx, key, &Profile{
		Name:name,
		Email:email,
		Gender:gender,
	})
}

type Profile struct {
	Name   string `datastore:"name, noindex"`
	Email  string `datastore:"email"`
	Gender int `datastore:"gender, noindex"`
}
