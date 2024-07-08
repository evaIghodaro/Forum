package handlers

import (
	auth "forum/internal/authentification"
	"html/template"
	"net/http"
)

var tmpl *template.Template

func InitializeAuthHandlers(t *template.Template) {
	tmpl = t
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth.RegisterHandler(w, r)
	} else {
		tmpl.ExecuteTemplate(w, "register.html", nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		auth.LoginHandler(w, r)
	} else {
		tmpl.ExecuteTemplate(w, "login.html", nil)
	}
}
