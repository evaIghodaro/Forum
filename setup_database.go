package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func setupDatabase() {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queries := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL,
        email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT
    );
    CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT,
        category_id INTEGER,
        author_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        approved BOOLEAN DEFAULT 0,
        FOREIGN KEY (category_id) REFERENCES categories (id),
        FOREIGN KEY (author_id) REFERENCES users (id)
    );
    CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT,
        post_id INTEGER,
        author_id INTEGER,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (author_id) REFERENCES users (id)
    );
    CREATE TABLE IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        post_id INTEGER,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (post_id) REFERENCES posts (id)
    );
    CREATE TABLE IF NOT EXISTS notifications (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER,
        message TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

    INSERT INTO categories (name, description) VALUES ('Recettes', 'Partagez vos recettes de pâtisserie préférées');
    INSERT INTO categories (name, description) VALUES ('Techniques de Pâtisserie', 'Discutez des techniques de cuisson, de glaçage, de décoration, etc.');
    INSERT INTO categories (name, description) VALUES ('Matériel et Ingrédients', 'Échangez des conseils sur les meilleurs ustensiles et ingrédients');
    INSERT INTO categories (name, description) VALUES ('Photos de Vos Créations', 'Montrez vos réalisations pâtissières et inspirez les autres');
    INSERT INTO categories (name, description) VALUES ('Questions et Conseils', 'Posez des questions et donnez des conseils sur la pâtisserie');
    `
	_, err = db.Exec(queries)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")
}
