package test

import (
	"database/sql"
	auth "forum/internal/authentification"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	queries := []string{
		`CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            username TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestRegisterHandler(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	auth.Initialize(db)

	req, err := http.NewRequest("POST", "/register", strings.NewReader("username=testuser&email=test@example.com&password=testpass"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	var username string
	err = db.QueryRow("SELECT username FROM users WHERE email = ?", "test@example.com").Scan(&username)
	if err != nil {
		t.Errorf("Failed to find user in DB: %v", err)
	}

	if username != "testuser" {
		t.Errorf("handler returned wrong username: got %v want %v", username, "testuser")
	}
}

func TestLoginHandler(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	auth.Initialize(db)

	hashedPassword, _ := auth.HashPassword("123456")
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "test", "test@test.com", hashedPassword)
	if err != nil {
		t.Fatalf("Failed to insert user into DB: %v", err)
	}

	req, err := http.NewRequest("POST", "/login", strings.NewReader("email=test@test.com&password=123456"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}
