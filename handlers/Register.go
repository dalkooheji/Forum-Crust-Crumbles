package forum

import (
	"database/sql"
	"fmt"
	"forum/database"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser saves a new user to the database
func RegisterUser(db *sql.DB, user User) error {
	query := `INSERT INTO Users (Username, Email, PasswordHash) VALUES (?, ?, ?)`
	_, err := db.Exec(query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

// RegHandler handles the registration form submission
func RegHandler(w http.ResponseWriter, r *http.Request) {
	db := database.DB

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validation
		errorMessages := []string{}

		// Check if username is empty or exceeds length
		if len(username) == 0 {
			errorMessages = append(errorMessages, "Username cannot be empty.")
		} else if len(username) > 20 {
			errorMessages = append(errorMessages, "Username cannot be longer than 20 characters.")
		} else if strings.Contains(username, " "){
			errorMessages = append(errorMessages, "Username should not contain spaces.")
		}

		// Check if email is empty
		if len(email) == 0 {
			errorMessages = append(errorMessages, "Email cannot be empty.")
		}

		emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
		if !emailRegex.MatchString(email) {
			errorMessages = append(errorMessages, "Please enter a valid email.")
		}

		// Check if password contains spaces or is too short
		if strings.Contains(password, " ") {
			errorMessages = append(errorMessages, "Password cannot contain spaces.")
		}
		if len(password) < 8 {
			errorMessages = append(errorMessages, "Password must be at least 8 characters long.")
		}

		// If there are errors, re-render the page with errors
		if len(errorMessages) > 0 {
			tmpl := template.Must(template.ParseFiles("templates/register.html"))
			tmpl.Execute(w, struct {
				Errors []string
			}{
				Errors: errorMessages,
			})
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user := User{
			Username:     username,
			Email:        email,
			PasswordHash: string(hashedPassword),
		}

		// Register the user
		err = RegisterUser(db, user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to register user: %v", err), http.StatusInternalServerError)
			return
		}

		// Retrieve UserID
		var userID int
		query := `SELECT UserID FROM Users WHERE Username = ?`
		err = db.QueryRow(query, username).Scan(&userID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found after registration", http.StatusInternalServerError)
				return
			}
			http.Error(w, "Failed to retrieve user ID", http.StatusInternalServerError)
			return
		}

		// Create session
		_, err = CreateSession(w, userID)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}

		// Redirect to home after successful registration
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/register.html"))
		tmpl.Execute(w, nil)
	}
}
