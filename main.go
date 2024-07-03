package main

import (
	"database/sql"
	auth "forum/internal/authentification"
	filters "forum/internal/filters"
	moderation "forum/internal/moderation"
	notification "forum/internal/notification"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Call setupDatabase function from setup_database.go
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

	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/login", auth.LoginHandler)
	// Ajoutez d'autres gestionnaires ici

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
