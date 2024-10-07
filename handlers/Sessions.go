package forum

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// Session structure to hold session data
type Session struct {
	SessionID string
	UserID    int
	ExpiresAt time.Time
}

// Generate a secure random session ID
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Create a new session for a user
func CreateSession(w http.ResponseWriter, userID int) (Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return Session{}, err
	}

	expiresAt := time.Now().Add(24 * time.Hour) // Session valid for 24 hours

	session := Session{
		SessionID: sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	// Store the session in the database
	// This assumes you have a function `SaveSession` to handle this
	if err := SaveSession(session); err != nil {
		return Session{}, err
	}

	// Set the session ID in a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Expires:  expiresAt,
		HttpOnly: true,
		Path:     "/",
	})

	return session, nil
}

// Retrieve session information based on session ID
func GetSession(r *http.Request) (Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return Session{}, err
	}

	// Fetch the session from the database using the session ID
	session, err := FetchSessionFromDB(cookie.Value)
	if err != nil {
		return Session{}, err
	}

	// Check if the session has expired
	if session.ExpiresAt.Before(time.Now()) {
		return Session{}, http.ErrNoCookie
	}

	return session, nil
}

// Delete session (logout)
func DeleteSession(w http.ResponseWriter, sessionID string) error {
	// Remove session from the database
	if err := RemoveSessionFromDB(sessionID); err != nil {
		return err
	}

	// Invalidate the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	return nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err = DeleteSession(w, session.SessionID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
