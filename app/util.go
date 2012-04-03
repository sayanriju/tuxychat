package tuxychat

import (
	"appengine"
	"appengine/user"
	"math/rand"
	"net/http"
	"time"
)

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	src = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
)

func ensureLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := appengine.NewContext(r)
		u := user.Current(c)

		if u == nil {
			url, err := user.LoginURL(c, r.URL.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, url, http.StatusFound)
			return
		}

		handler(w, r)
	}
}

func randStr(size int) string {
	str := ""

	for i := 0; i < size; i++ {
		str += string(src[rnd.Int()%len(src)])
	}
	return str
}
