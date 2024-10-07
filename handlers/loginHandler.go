package forum

import (
	"database/sql"
	"fmt"
	"forum/database"
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// FetchUserByUsername retrieves a user by their username
func FetchUserByUsername(db *sql.DB, username string) (User, int, error) {
	var user User
	var userID int
	query := `SELECT UserID, Username, PasswordHash FROM Users WHERE Username = ?`
	err := db.QueryRow(query, username).Scan(&userID, &user.Username, &user.PasswordHash)
	if err != nil {
		return user, 0, fmt.Errorf("failed to fetch user: %v", err)
	}
	return user, userID, nil
}

// LoginHandler handles the login form submission
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db := database.DB

		username := r.FormValue("username")
		password := r.FormValue("password")
		var errorMessage string

		user, userID, err := FetchUserByUsername(db, username)
		if err != nil {
			log.Printf("Error fetching user: %v", err)
			errorMessage = "Invalid username or password."
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
			if err != nil {
				log.Printf("Password mismatch: %v", err)
				errorMessage = "Invalid username or password."
			}
		}

		if errorMessage != "" {
			tmpl := template.Must(template.ParseFiles("templates/login.html"))
			tmpl.Execute(w, struct {
				ErrorMessage string
			}{
				ErrorMessage: errorMessage,
			})
			return
		}
		// Check for an existing session for the new userID
		existingSession, err := FetchSessionByUserID(db, userID)
		if err == nil {
			// If an existing session is found for this userID, delete it
			if err := DeleteSession(w, existingSession.SessionID); err != nil {
				log.Printf("Error deleting existing session: %v", err)
			}
		} else if err != sql.ErrNoRows {
			log.Printf("Error fetching existing session: %v", err)
			http.Error(w, "Failed to validate session", http.StatusInternalServerError)
			return
		}

		// Create a new session
		_, err = CreateSession(w, userID)
		if err != nil {
			log.Printf("Error creating session: %v", err)
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		log.Println("Login successful, redirecting to home")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	}
}

