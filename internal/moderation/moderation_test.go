package moderation

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
		`CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, category TEXT, author_id INTEGER, approved BOOLEAN, report_count INTEGER)`,
	}

	for _, query := range createTableQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	insertDataQueries := []string{
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (1, 'Title 1', 'Content 1', 'Category1', 1, 0, 0)`,
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (2, 'Title 2', 'Content offensiveWord1', 'Category2', 2, 0, 0)`,
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (3, 'Title 3', 'Spam Content', 'Category1', 3, 0, 6)`,
	}

	for _, query := range insertDataQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestApproveOrRejectPost(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	// Cas où le post doit être approuvé
	err = ApproveOrRejectPost(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var approved bool
	err = db.QueryRow("SELECT approved FROM posts WHERE id = 1").Scan(&approved)
	if err != nil {
		t.Fatal(err)
	}
	if !approved {
		t.Fatalf("expected post to be approved")
	}

	// Cas où le post contient du contenu offensant et doit être rejeté
	err = ApproveOrRejectPost(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = db.QueryRow("SELECT approved FROM posts WHERE id = 2").Scan(&approved)
	if err != nil {
		t.Fatal(err)
	}
	if approved {
		t.Fatalf("expected post to be rejected due to offensive content")
	}

	// Cas où le post a beaucoup de signalements et doit être rejeté
	err = ApproveOrRejectPost(3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = db.QueryRow("SELECT approved FROM posts WHERE id = 3").Scan(&approved)
	if err != nil {
		t.Fatal(err)
	}
	if approved {
		t.Fatalf("expected post to be rejected due to high report count")
	}
}
