package test

import (
	"forum/internal/notification"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestNotificationsIntegration(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	notification.Initialize(db)

	tests := []struct {
		notifyFunc  func(int) error
		expectedMsg string
	}{
		{notification.NotifyPostLiked, "Your post with ID 1 was liked"},
		{notification.NotifyPostDisliked, "Your post with ID 1 was disliked"},
		{notification.NotifyPostCommented, "Your post with ID 1 was commented on"},
	}

	for _, test := range tests {
		err := test.notifyFunc(1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		var message string
		err = db.QueryRow("SELECT message FROM notifications WHERE user_id = 1 ORDER BY id DESC LIMIT 1").Scan(&message)
		if err != nil {
			t.Fatal(err)
		}
		if message != test.expectedMsg {
			t.Fatalf("expected message '%s', got '%s'", test.expectedMsg, message)
		}
	}
}
