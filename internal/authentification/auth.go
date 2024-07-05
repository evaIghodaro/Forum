package authentification

import (
	"database/sql"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func Initialize(database *sql.DB) {
	db = database
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		hashedPassword, err := HashPassword(password)
		if err != nil {
			http.Error(w, "Unable to create your account", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", username, email, hashedPassword)
		if err != nil {
			http.Error(w, "Unable to create your account", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	} else {
		tmpl, _ := template.ParseFiles("templates/register.html")
		tmpl.Execute(w, nil)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		var hash string
		err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&hash)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if !CheckPasswordHash(password, hash) {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		tmpl, _ := template.ParseFiles("templates/login.html")
		tmpl.Execute(w, nil)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Example: Check if the user is authenticated
		// In a real application, you should check session or token here
		if r.URL.Path == "/login" || r.URL.Path == "/register" {
			next.ServeHTTP(w, r)
			return
		}

		// Redirect to login page if not authenticated
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}
