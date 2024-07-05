package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	auth "forum/internal/authentification"
	filters "forum/internal/filters"
	moderation "forum/internal/moderation"
	notification "forum/internal/notification"

	_ "github.com/mattn/go-sqlite3"
)

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob(filepath.Join("templates", "*.html"))
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
}

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

func main() {
	setupDatabase() // Appel à la fonction pour initialiser la base de données

	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	auth.Initialize(db)
	filters.Initialize(db)
	moderation.Initialize(db)
	notification.Initialize(db)

	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, db)
	})
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/post", postHandler)
	http.Handle("/create_post", auth.AuthMiddleware(http.HandlerFunc(createPostHandler)))
	http.HandleFunc("/profile", profileHandler)

	http.Handle("/protected", auth.AuthMiddleware(http.HandlerFunc(protectedHandler)))

	// Servir les fichiers statiques
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./image"))))

	log.Println("Starting server on :8082")
	fmt.Println("Server is running at http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /register")
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		auth.RegisterHandler(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /login")
	if r.Method == http.MethodGet {
		err := templates.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		auth.LoginHandler(w, r)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	posts := []Post{}
	categories := []Category{}

	// Ignorer l'avertissement pour le paramètre non utilisé
	_ = r

	// Récupérer les derniers posts
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

	// Récupérer les catégories populaires
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

	if err := templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "categories.html"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "post.html"))
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "create_post.html"))
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join("templates", "profile.html"))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This is a protected page")
}
