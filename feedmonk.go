package feedmonk

import (
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	http.HandleFunc("/", handler)
}

var tpl = template.Must(template.ParseGlob("templates/*.html"))

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var email, logout string
	if u := user.Current(c); u != nil {
		logout, _ = user.LogoutURL(c, "/")
		email = u.Email
		data := struct {
			Logout, Email string
		}{
			Logout: logout,
			Email:  email,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.ExecuteTemplate(w, "index.html", data); err != nil {
			log.Errorf(c, "%v", err)
		}
		return

	} else {
		url, err := user.LoginURL(c, r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}
}
