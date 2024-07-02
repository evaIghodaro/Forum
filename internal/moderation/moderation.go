package moderation

import (
	"database/sql"
	"strings"
)

var db *sql.DB

func Initialize(database *sql.DB) {
	db = database
}

func ApproveOrRejectPost(postID int) error {
	var content string
	var reportCount int

	err := db.QueryRow("SELECT content, report_count FROM posts WHERE id = ?", postID).Scan(&content, &reportCount)
	if err != nil {
		return err
	}

	// Logic to reject offensive content
	if containsOffensiveContent(content) || reportCount > 5 {
		_, err = db.Exec("UPDATE posts SET approved = 0 WHERE id = ?", postID)
		return err
	}

	// Approve the post if it's not offensive and doesn't have too many reports
	_, err = db.Exec("UPDATE posts SET approved = 1 WHERE id = ?", postID)
	return err
}

func containsOffensiveContent(content string) bool {
	// Simple check for offensive words
	offensiveWords := []string{"offensive", "spam", "inappropriate"}
	for _, word := range offensiveWords {
		if contains(content, word) {
			return true
		}
	}
	return false
}

func contains(content, word string) bool {
	// Check if the content contains the offensive word
	return strings.Contains(content, word)
}
