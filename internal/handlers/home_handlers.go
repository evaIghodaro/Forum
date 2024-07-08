package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
)

type Post struct {
	ID      int
	Title   string
	Snippet string
}

type Category struct {
	ID          int
	Name        string
	Description string
}

var homeTemplates *template.Template

func InitializeHomeHandlers(t *template.Template) {
	homeTemplates = t
}

func HomeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	posts := []Post{}
	categories := []Category{}

	rows, err := db.Query("SELECT id, title, SUBSTR(content, 1, 100) as snippet FROM posts ORDER BY created_at DESC LIMIT 5")
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Snippet); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	rows, err = db.Query("SELECT id, name, description FROM categories LIMIT 5")
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			http.Error(w, "Failed to scan category", http.StatusInternalServerError)
			return
		}
		categories = append(categories, category)
	}

	data := struct {
		Posts      []Post
		Categories []Category
	}{
		Posts:      posts,
		Categories: categories,
	}

	if err := homeTemplates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
