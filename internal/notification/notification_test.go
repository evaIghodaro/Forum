package notification

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	createTableQueries := []string{
		`CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, category TEXT, author_id INTEGER, approved BOOLEAN)`,
		`CREATE TABLE notifications (id INTEGER PRIMARY KEY, user_id INTEGER, message TEXT)`,
	}

	for _, query := range createTableQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	insertDataQueries := []string{
		`INSERT INTO posts (id, title, content, category, author_id, approved) VALUES (1, 'Title 1', 'Content 1', 'Category1', 1, 0)`,
	}

	for _, query := range insertDataQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestNotifyPostLiked(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	err = NotifyPostLiked(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var message string
	err = db.QueryRow("SELECT message FROM notifications WHERE user_id = 1").Scan(&message)
	if err != nil {
		t.Fatal(err)
	}
	expectedMessage := "Your post with ID 1 was liked"
	if message != expectedMessage {
		t.Fatalf("expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestNotifyPostDisliked(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	err = NotifyPostDisliked(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var message string
	err = db.QueryRow("SELECT message FROM notifications WHERE user_id = 1").Scan(&message)
	if err != nil {
		t.Fatal(err)
	}
	expectedMessage := "Your post with ID 1 was disliked"
	if message != expectedMessage {
		t.Fatalf("expected message '%s', got '%s'", expectedMessage, message)
	}
}

func TestNotifyPostCommented(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	err = NotifyPostCommented(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var message string
	err = db.QueryRow("SELECT message FROM notifications WHERE user_id = 1").Scan(&message)
	if err != nil {
		t.Fatal(err)
	}
	expectedMessage := "Your post with ID 1 was commented on"
	if message != expectedMessage {
		t.Fatalf("expected message '%s', got '%s'", expectedMessage, message)
	}
}
