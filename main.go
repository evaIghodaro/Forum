package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	auth "forum/internal/authentification"
	filters "forum/internal/filters"
	moderation "forum/internal/moderation"
	notification "forum/internal/notification"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialiser la base de données
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

	// Gestionnaires pour les pages d'authentification
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)

	// Gestionnaires pour les pages du forum
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/post", postHandler)
	http.Handle("/create_post", auth.AuthMiddleware(http.HandlerFunc(createPostHandler)))
	http.HandleFunc("/profile", profileHandler)

	// Middleware pour protéger certaines routes
	http.Handle("/protected", auth.AuthMiddleware(http.HandlerFunc(protectedHandler)))

	log.Println("Starting server on :8081")
	fmt.Println("Server is running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/index.html")
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/categories.html")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/post.html")
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/create_post.html")
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/templates/profile.html")
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected page")
}
