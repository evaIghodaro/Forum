package handlers

import (
	"html/template"
	"net/http"
)

var postTemplates *template.Template

func InitializePostHandlers(t *template.Template) {
	postTemplates = t
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	postTemplates.ExecuteTemplate(w, "post.html", nil)
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	postTemplates.ExecuteTemplate(w, "create_post.html", nil)
}
