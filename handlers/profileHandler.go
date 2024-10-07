package forum

import (
	"forum/database"
	"html/template"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ProfileData struct {
	Icon          string
	Username      string
	Email         string
	DateCreated   time.Time
	CreatedPosts  int
	LikedPosts    int
	DislikedPosts int
	// structs that we'll use to send the posts to the template
	UserPosts         []Post
	LikedPostsList    []Post
	DislikedPostsList []Post
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Fetch user details from the database
	db := database.DB

	var profileData ProfileData

	// Fetch username, email, date created
	query := `SELECT Username, Email, CreatedAt FROM Users WHERE UserID = ?`
	err = db.QueryRow(query, session.UserID).Scan(&profileData.Username, &profileData.Email, &profileData.DateCreated)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the count of created posts
	query = `SELECT COUNT(*) FROM Posts WHERE UserID = ?`
	err = db.QueryRow(query, session.UserID).Scan(&profileData.CreatedPosts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the count of liked posts
	query = `SELECT COUNT(*) FROM Likes WHERE UserID = ? AND IsLike = 1 AND PostID IS NOT NULL`
	err = db.QueryRow(query, session.UserID).Scan(&profileData.LikedPosts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the count of disliked posts
	query = `SELECT COUNT(*) FROM Dislikes WHERE UserID = ? AND IsDislike = 1 AND PostID IS NOT NULL`
	err = db.QueryRow(query, session.UserID).Scan(&profileData.DislikedPosts)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch posts created by the user
	query = `
        SELECT p.PostID, p.Title, p.Content, u.Username, p.CreatedAt, 
               (SELECT COUNT(*) FROM Likes WHERE PostID = p.PostID AND IsLike = 1) as LikeCount
        FROM Posts p 
        JOIN Users u ON p.UserID = u.UserID 
        WHERE p.UserID = ?
        ORDER BY p.CreatedAt DESC`
	rows, err := db.Query(query, session.UserID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &post.LikeCount); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		profileData.UserPosts = append(profileData.UserPosts, post)
	}

	// Fetch posts liked by the user
	query = `
        SELECT p.PostID, p.Title, p.Content, u.Username, p.CreatedAt, 
               (SELECT COUNT(*) FROM Likes WHERE PostID = p.PostID AND IsLike = 1) as LikeCount
        FROM Posts p 
        JOIN Users u ON p.UserID = u.UserID 
        JOIN Likes l ON p.PostID = l.PostID
        WHERE l.UserID = ? AND l.IsLike = 1
        ORDER BY p.CreatedAt DESC`
	rows, err = db.Query(query, session.UserID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &post.LikeCount); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		profileData.LikedPostsList = append(profileData.LikedPostsList, post)
	}
	// Fetch posts disliked by the user
	query = `
			SELECT p.PostID, p.Title, p.Content, u.Username, p.CreatedAt, 
							(SELECT COUNT(*) FROM Likes WHERE PostID = p.PostID AND IsLike = 1) as LikeCount,
							(SELECT COUNT(*) FROM Dislikes WHERE PostID = p.PostID AND IsDislike = 1) as DislikeCount
			FROM Posts p 
			JOIN Users u ON p.UserID = u.UserID 
			JOIN Dislikes d ON p.PostID = d.PostID
			WHERE d.UserID = ? AND d.IsDislike = 1
			ORDER BY p.CreatedAt DESC`
	rows, err = db.Query(query, session.UserID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &post.LikeCount, &post.DislikeCount); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		profileData.DislikedPostsList = append(profileData.DislikedPostsList, post)
	}
	// For now, the icon can be a placeholder image; you can later update this to fetch the actual user-uploaded image
	profileData.Icon = "https://digitalhealthskills.com/wp-content/uploads/2022/11/3da39-no-user-image-icon-27.png"

	// Parse the profile.html template
	tmpl := template.Must(template.ParseFiles("templates/profile.html"))
	tmpl.Execute(w, profileData)
}
