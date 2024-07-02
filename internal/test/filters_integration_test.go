package test

import (
	"forum/internal/filters"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetPostsByCategoryIntegration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, title, content, category, author_id, likes FROM posts WHERE category = ?").
		WithArgs("Category1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content", "category", "author_id", "likes"}).
			AddRow(1, "Title 1", "Content 1", "Category1", 1, 10).
			AddRow(3, "Title 3", "Content 3", "Category1", 1, 8))

	filters.Initialize(db)

	posts, err := filters.GetPostsByCategory("Category1")
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
}
