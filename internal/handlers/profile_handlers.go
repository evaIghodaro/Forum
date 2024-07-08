package handlers

import (
	"html/template"
	"net/http"
)

var profileTemplates *template.Template

func InitializeProfileHandlers(t *template.Template) {
	profileTemplates = t
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	profileTemplates.ExecuteTemplate(w, "profile.html", nil)
}
