package tuxychat

import (
	"code.google.com/p/gorilla/mux"
	"github.com/tuxychandru/gomesh"
	"net/http"
	"text/template"
)

func init() {
	parseTemplates()

	dec := gomesh.New(template.Must(template.ParseFiles("templates/decorator.html")))

	r := mux.NewRouter()
	r.HandleFunc("/", ensureLogin(home)).Methods("GET")
	r.HandleFunc("/new", ensureLogin(new)).Methods("GET")
	r.HandleFunc("/join", ensureLogin(join)).Methods("POST")
	r.HandleFunc("/{id}", ensureLogin(chat)).Methods("GET")
	r.HandleFunc("/msg/{id}", ensureLogin(msg)).Methods("POST")

	http.Handle("/", dec.Wrap(r))
}
