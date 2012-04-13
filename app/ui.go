package tuxychat

import (
	"appengine"
	"appengine/user"
	"code.google.com/p/gorilla/mux"
	"fmt"
	"html/template"
	"net/http"
)

var tmpls *template.Template

type homeVal struct {
	User    *user.User
	Missing bool
	Unknown bool
}

type chatVal struct {
	User   *user.User
	RoomId string
	Token  string
}

func home(w http.ResponseWriter, r *http.Request) {
	u := user.Current(appengine.NewContext(r))
	tmpls.ExecuteTemplate(w, "home.html", homeVal{u, false, false})
}

func new(w http.ResponseWriter, r *http.Request) {
	roomId := randStr(10)
	if err := createRoom(appengine.NewContext(r), roomId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/"+roomId, http.StatusFound)
}

func join(w http.ResponseWriter, r *http.Request) {
	roomId := r.FormValue("roomId")

	if len(roomId) == 0 {
		u := user.Current(appengine.NewContext(r))
		tmpls.ExecuteTemplate(w, "home.html", homeVal{u, true, false})
	} else {
		http.Redirect(w, r, "/"+roomId, http.StatusFound)
	}
}

func chat(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	roomId := mux.Vars(r)["id"]

	exists, err := roomExists(c, roomId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		tmpls.ExecuteTemplate(w, "home.html", homeVal{u, false, true})
		return
	}

	token, err := joinRoom(c, roomId, u.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpls.ExecuteTemplate(w, "chat.html", chatVal{u, roomId, token})
}

func msg(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	u := user.Current(c)
	roomId := mux.Vars(r)["id"]

	if err := publish(c, roomId, u.Email, r.FormValue("msg")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintln(w, "sent")
}

func parseTemplates() {
	tmpls = template.Must(template.ParseGlob("templates/*"))
}
