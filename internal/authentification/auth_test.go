package authentification

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	createTable := `
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        email TEXT,
        password TEXT
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	return db
}

func TestRegisterHandler(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	req, err := http.NewRequest("POST", "/register", strings.NewReader("username=testuser&email=test@example.com&password=testpass"))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}

func TestLoginHandler(t *testing.T) {
	db = setupTestDB()
	defer db.Close()

	// Insert a test user into the database
	password, err := HashPassword("123456")
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
	handler := http.HandlerFunc(LoginHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}

	cookie := rr.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, "session=") {
		t.Errorf("Set-Cookie header does not contain session")
	}
}

func TestLogoutHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LogoutHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusSeeOther)
	}
}
