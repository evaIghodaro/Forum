package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	auth "forum/internal/authentification"
	filters "forum/internal/filters"
	handlers "forum/internal/handlers"
	moderation "forum/internal/moderation"
	notification "forum/internal/notification"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

func main() {
	setupDatabase()

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auth.Initialize(db)
	filters.Initialize(db)
	moderation.Initialize(db)
	notification.Initialize(db)

	// Initialize handlers with templates
	handlers.InitializeAuthHandlers(templates)
	handlers.InitializePostHandlers(templates)
	handlers.InitializeHomeHandlers(templates)
	handlers.InitializeProfileHandlers(templates)

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		handlers.HomeHandler(w, r, db)
	})
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.Handle("/create_post", auth.AuthMiddleware(http.HandlerFunc(handlers.CreatePostHandler)))
	http.HandleFunc("/profile", handlers.ProfileHandler)

	http.Handle("/protected", auth.AuthMiddleware(http.HandlerFunc(protectedHandler)))

	// Serve static files
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./image"))))

	log.Println("Starting server on :8082")
	fmt.Println("Server is running at http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/categories.html")
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected page")
}
