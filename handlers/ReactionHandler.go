package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"forum/database"
	"net/http"
	"strconv"
)

// ToggleCommentReactionHandler handles the like/unlike functionality for comments
func ToggleCommentReactionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	logged := true
	if err != nil {
		logged = false;
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println(logged)
	var requestData struct {
		CommentID interface{} `json:"commentID"`
		IsLike    bool        `json:"isLike"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var commentID int
	switch v := requestData.CommentID.(type) {
	case string:
		commentID, err = strconv.Atoi(v)
		if err != nil {
			http.Error(w, "Invalid CommentID", http.StatusBadRequest)
			return
		}
	case float64:
		commentID = int(v)
	default:
		http.Error(w, "Invalid CommentID type", http.StatusBadRequest)
		return
	}

	db := database.DB

	// Check for existing reaction
	var existingReactionID int
	query := `SELECT ReactionID FROM Reactions WHERE UserID = ? AND CommentID = ?`
	err = db.QueryRow(query, session.UserID, commentID).Scan(&existingReactionID)

	if err == nil && existingReactionID != 0 {
		// Remove existing reaction
		var currentReaction bool
		query = `SELECT IsLike FROM Reactions WHERE ReactionID = ?`
		err = db.QueryRow(query, existingReactionID).Scan(&currentReaction)
		if currentReaction != requestData.IsLike {
			query = `INSERT INTO Reactions (UserID, CommentID, IsLike) VALUES (?, ?, ?)`
			_, err = db.Exec(query, session.UserID, commentID, requestData.IsLike)
			if err != nil {
				http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
				return
			}
		}
		query := `DELETE FROM Reactions WHERE ReactionID = ?`
		_, err = db.Exec(query, existingReactionID)
		if err != nil {
			http.Error(w, "Failed to remove reaction", http.StatusInternalServerError)
			return
		}

		
	} else if err == sql.ErrNoRows {
		// No existing reaction, so insert a new reaction
		query = `INSERT INTO Reactions (UserID, CommentID, IsLike) VALUES (?, ?, ?)`
		_, err = db.Exec(query, session.UserID, commentID, requestData.IsLike)
		if err != nil {
			http.Error(w, "Failed to add reaction", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Failed to toggle reaction", http.StatusInternalServerError)
		return
	}
	// Return the updated counts
	var likeCount, dislikeCount int
	query = `SELECT COUNT(*) FROM Reactions WHERE CommentID = ? AND IsLike = 1`
	err = db.QueryRow(query, commentID).Scan(&likeCount)
	if err != nil {
		http.Error(w, "Failed to retrieve like count", http.StatusInternalServerError)
		return
	}
	query = `SELECT COUNT(*) FROM Reactions WHERE CommentID = ? AND IsLike = 0`
	err = db.QueryRow(query, commentID).Scan(&dislikeCount)
	if err != nil {
		http.Error(w, "Failed to retrieve dislike count", http.StatusInternalServerError)
		return
	}

	responseData := map[string]interface{}{
		"logged":		logged,
		"success":      true,
		"likeCount":    likeCount,
		"dislikeCount": dislikeCount,
		"isLiked":      requestData.IsLike,
	}
	json.NewEncoder(w).Encode(responseData)
}
