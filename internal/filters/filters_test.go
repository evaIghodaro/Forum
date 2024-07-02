package filters

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
		`CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, category TEXT, author_id INTEGER, likes INTEGER)`,
		`CREATE TABLE likes (user_id INTEGER, post_id INTEGER)`,
	}

	for _, query := range createTableQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	insertDataQueries := []string{
		`INSERT INTO posts (id, title, content, category, author_id, likes) VALUES (1, 'Title 1', 'Content 1', 'Category1', 1, 10)`,
		`INSERT INTO posts (id, title, content, category, author_id, likes) VALUES (2, 'Title 2', 'Content 2', 'Category2', 2, 5)`,
		`INSERT INTO posts (id, title, content, category, author_id, likes) VALUES (3, 'Title 3', 'Content 3', 'Category1', 1, 8)`,
		`INSERT INTO likes (user_id, post_id) VALUES (1, 1)`,
		`INSERT INTO likes (user_id, post_id) VALUES (1, 3)`,
	}

	for _, query := range insertDataQueries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func TestGetPostsByCategory(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	posts, err := GetPostsByCategory("Category1")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

func TestGetPostsByUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	posts, err := GetPostsByUser(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}

func TestGetLikedPostsByUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	Initialize(db)

	posts, err := GetLikedPostsByUser(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}
