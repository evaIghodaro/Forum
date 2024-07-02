package test

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	queries := []string{
		`CREATE TABLE posts (id INTEGER PRIMARY KEY, title TEXT, content TEXT, category TEXT, author_id INTEGER, approved BOOLEAN, report_count INTEGER)`,
		`CREATE TABLE notifications (id INTEGER PRIMARY KEY, user_id INTEGER, message TEXT)`,
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (1, 'Title 1', 'Content 1', 'Category1', 1, 0, 0)`,
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (2, 'Title 2', 'offensive content', 'Category2', 2, 0, 0)`,
		`INSERT INTO posts (id, title, content, category, author_id, approved, report_count) VALUES (3, 'Title 3', 'Spam content', 'Category3', 3, 0, 6)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return nil, err
		}
	}

	return db, nil
}
