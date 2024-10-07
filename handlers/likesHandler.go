package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/database"
	"net/http"
	"strconv"
)

// ToggleLikeHandler handles the like/unlike functionality
func ToggleLikeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current session (ensure user is logged in)
	session, err := GetSession(r)
	logged := true
	if err != nil {
		logged = false
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println("logged in:", logged)
	// Parse the incoming request body to get the postID
	var requestData struct {
		PostID interface{} `json:"postID"` // Change type to interface{} to handle both string and int
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Convert the PostID to an integer if it's a string
	var postID int
	switch v := requestData.PostID.(type) {
	case string:
		postID, err = strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Invalid PostID", http.StatusBadRequest)
			return
		}
	case float64:
		postID = int(v)
	default:
		http.Error(w, "Invalid PostID type", http.StatusBadRequest)
		return
	}

	// Open the database
	db := database.DB

	// Check if the user already disliked this post
	var existingDislikeID, existingLikeID int
	query := `SELECT DislikeID FROM Dislikes WHERE UserID = ? AND PostID = ? AND IsDislike = 1`
	err = db.QueryRow(query, session.UserID, postID).Scan(&existingDislikeID)

	if err == nil && existingDislikeID != 0 {
		// Remove existing dislike
		query = `DELETE FROM Dislikes WHERE DislikeID = ?`
		_, err = db.Exec(query, existingDislikeID)
		if err != nil {
			http.Error(w, "Failed to remove dislike", http.StatusInternalServerError)
			return
		}
	}

	// Check if the user already liked this post
	query = `SELECT LikeID FROM Likes WHERE UserID = ? AND PostID = ? AND IsLike = 1`
	err = db.QueryRow(query, session.UserID, postID).Scan(&existingLikeID)

	if err == sql.ErrNoRows {
		// No existing like, so insert a new like
		query = `INSERT INTO Likes (UserID, PostID, IsLike) VALUES (?, ?, 1)`
		_, err = db.Exec(query, session.UserID, postID)
		if err != nil {
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}
	} else if err == nil {
		// Like exists, so remove it
		query = `DELETE FROM Likes WHERE LikeID = ?`
		_, err = db.Exec(query, existingLikeID)
		if err != nil {
			http.Error(w, "Failed to unlike post", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Failed to toggle like", http.StatusInternalServerError)
		return
	}

	// Return the updated like count
	var likeCount int
	query = `SELECT COUNT(*) FROM Likes WHERE PostID = ? AND IsLike = 1`
	err = db.QueryRow(query, postID).Scan(&likeCount)
	if err != nil {
		http.Error(w, "Failed to retrieve like count", http.StatusInternalServerError)
		return
	}
	var dislikeCount int
	query = `SELECT COUNT(*) FROM Dislikes WHERE PostID = ? AND IsDislike = 1`
	err = db.QueryRow(query, postID).Scan(&dislikeCount)
	if err != nil {
		http.Error(w, "Failed to retrieve dislike count", http.StatusInternalServerError)
		return
	}

	// Send the updated like count back to the client
	responseData := map[string]interface{}{
		"success":      true,
		"likeCount":    likeCount,
		"dislikeCount": dislikeCount,
		"isLiked":      existingLikeID > 0,
	}
	json.NewEncoder(w).Encode(responseData)
}
