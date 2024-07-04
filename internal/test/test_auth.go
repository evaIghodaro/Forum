package test

import (
	"database/sql"
	auth "forum/internal/authentification" // Assurez-vous que le chemin d'importation est correct
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
		`CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT);`,
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

	req, err := http.NewRequest("POST", "/register", strings.NewReader("username=test&email=test@test.com&password=123456"))
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
	err = db.QueryRow("SELECT username FROM users WHERE email = ?", "test@test.com").Scan(&username)
	if err != nil {
		t.Errorf("Failed to find user in DB: %v", err)
	}

	if username != "test" {
		t.Errorf("handler returned wrong username: got %v want %v", username, "test")
	}
}

func TestLoginHandler(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test DB: %v", err)
	}
	auth.Initialize(db)

	password, err := auth.HashPassword("123456")
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "test", "test@test.com", password)
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

func TestLogoutHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(auth.LogoutHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}
