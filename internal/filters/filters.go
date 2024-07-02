package filters

import (
	"database/sql"
)

type Post struct {
	ID       int
	Title    string
	Content  string
	Category string
	AuthorID int
	Likes    int
}

var db *sql.DB

func Initialize(database *sql.DB) {
	db = database
}

func GetPostsByCategory(category string) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, category, author_id, likes FROM posts WHERE category = ?", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.AuthorID, &post.Likes); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByUser(userID int) ([]Post, error) {
	rows, err := db.Query("SELECT id, title, content, category, author_id, likes FROM posts WHERE author_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.AuthorID, &post.Likes); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetLikedPostsByUser(userID int) ([]Post, error) {
	query := `
	SELECT p.id, p.title, p.content, p.category, p.author_id, p.likes
	FROM posts p
	INNER JOIN likes l ON p.id = l.post_id
	WHERE l.user_id = ?
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &post.AuthorID, &post.Likes); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
