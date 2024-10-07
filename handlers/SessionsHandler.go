package forum

import (
	"database/sql"
	"forum/database"
	"time"
	// "time"
)

// SaveSession saves a new session to the database
func SaveSession(session Session) error {
	db := database.DB
	var err error
	query := `INSERT INTO Sessions (SessionID, UserID, ExpiresAt) VALUES (?, ?, ?)`
	_, err = db.Exec(query, session.SessionID, session.UserID, session.ExpiresAt)
	return err
}

// FetchSessionFromDB retrieves a session by its ID from the database
func FetchSessionFromDB(sessionID string) (Session, error) {
	db := database.DB
	var err error
	var session Session
	query := `SELECT SessionID, UserID, ExpiresAt FROM Sessions WHERE SessionID = ?`
	err = db.QueryRow(query, sessionID).Scan(&session.SessionID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

// RemoveSessionFromDB deletes a session from the database
func RemoveSessionFromDB(sessionID string) error {
	db := database.DB
	var err error
	query := `DELETE FROM Sessions WHERE SessionID = ?`
	_, err = db.Exec(query, sessionID)
	return err
}

// FetchUsernameByID retrieves the username by user ID
func FetchUsernameByID(userID int) (string, error) {
	db := database.DB

	var username string
	query := `SELECT Username FROM Users WHERE UserID = ?`
	err := db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

// FetchSessionByUserID retrieves a session by user ID
func FetchSessionByUserID(db *sql.DB, userID int) (Session, error) {
	var session Session
	query := `SELECT SessionID, UserID, ExpiresAt FROM Sessions WHERE UserID = ? AND ExpiresAt > ?`
	err := db.QueryRow(query, userID, time.Now()).Scan(&session.SessionID, &session.UserID, &session.ExpiresAt)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}
