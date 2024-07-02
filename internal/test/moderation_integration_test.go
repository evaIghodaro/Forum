package test

import (
	"forum/internal/moderation"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestApproveOrRejectPostIntegration(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	moderation.Initialize(db)

	tests := []struct {
		postID        int
		shouldApprove bool
	}{
		{1, true},  // Should be approved
		{2, false}, // Should be rejected due to offensive content
		{3, false}, // Should be rejected due to high report count
	}

	for _, test := range tests {
		err := moderation.ApproveOrRejectPost(test.postID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var approved bool
		err = db.QueryRow("SELECT approved FROM posts WHERE id = ?", test.postID).Scan(&approved)
		if err != nil {
			t.Fatal(err)
		}
		if approved != test.shouldApprove {
			t.Fatalf("expected post %d approval status to be %v", test.postID, test.shouldApprove)
		}
	}
}
