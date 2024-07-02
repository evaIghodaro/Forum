package notification

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

// Initialize with database connection
func Initialize(database *sql.DB) {
	db = database
}

// SendNotification sends a notification to a user
func SendNotification(userID int, message string) error {
	_, err := db.Exec("INSERT INTO notifications (user_id, message) VALUES (?, ?)", userID, message)
	return err
}

// NotifyPostLiked notifies a user that their post was liked
func NotifyPostLiked(postID int) error {
	var userID int
	err := db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID).Scan(&userID)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Your post with ID %d was liked", postID)
	return SendNotification(userID, message)
}

// NotifyPostDisliked notifies a user that their post was disliked
func NotifyPostDisliked(postID int) error {
	var userID int
	err := db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID).Scan(&userID)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Your post with ID %d was disliked", postID)
	return SendNotification(userID, message)
}

// NotifyPostCommented notifies a user that their post was commented on
func NotifyPostCommented(postID int) error {
	var userID int
	err := db.QueryRow("SELECT author_id FROM posts WHERE id = ?", postID).Scan(&userID)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("Your post with ID %d was commented on", postID)
	return SendNotification(userID, message)
}
