package forum

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// PostDetailsHandler fetches the details of a single post and renders it using a template
func PostDetailsHandler(w http.ResponseWriter, r *http.Request) {

	_, errr := GetSession(r)
	logged := true
	if errr != nil {
		logged = false
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
		// return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get postID from the URL query parameters
	parameters := strings.Split(r.URL.Path, "/")
	postString := ""

	if len(parameters) == 3 && parameters[2] != "" {
		postString = parameters[2]
	} else {
		renderErrorTemplate(w, http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(postString)
	if err != nil || !isIdValid(id) {
		renderErrorTemplate(w, http.StatusBadRequest)
		// http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}
	// Fetch post details from the database
	post, err := fetchPostDetails(id, logged)
	if err != nil {
		http.Error(w, "Error fetching post details", http.StatusInternalServerError)
		return
	}

	// post.Logged = logged

	tmpl := template.Must(template.ParseFiles("templates/singlePost.html"))
	err = tmpl.ExecuteTemplate(w, "singlePost.html", struct {
		Post
		Errors []string
	}{
		Post:   post,
		Errors: nil,
	})
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the user is authenticated
	session, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Extract postID from URL path
	parameters := strings.Split(r.URL.Path, "/")
	postString := ""
	if len(parameters) == 4 && parameters[2] != "" && parameters[3] == "comment" {
		postString = parameters[2]
	} else {
		http.Error(w, "Post ID not found", http.StatusNotFound)
		return
	}

	postID, err := strconv.Atoi(postString)
	if err != nil {
		http.Error(w, "Invalid postID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodPost {
		content := r.FormValue("content")
		var errorMessage []string

		if strings.TrimSpace(content) == "" {
			errorMessage = append(errorMessage, "Comment cannot be empty.")
		} else {
			// Insert the comment
			_, err := db.Exec(`INSERT INTO Comments (PostID, UserID, Content, CreatedAt) VALUES (?, ?, ?, ?)`,
				postID, session.UserID, content, time.Now())
			if err != nil {
				http.Error(w, "Failed to create comment", http.StatusInternalServerError)
				return
			}
			redirectURL := fmt.Sprintf("/posts/%d", postID)
			http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			return
		}

		// Fetch post details to render the template with error messages
		post, err := fetchPostDetails(postID, true) // Pass `true` since the user is authenticated
		if err != nil {
			http.Error(w, "Error fetching post details", http.StatusInternalServerError)
			return
		}

		// If there are errors, re-render the page with errors
		tmpl := template.Must(template.ParseFiles("templates/singlePost.html"))
		tmpl.ExecuteTemplate(w, "singlePost.html", struct {
			Post
			Errors []string
		}{
			Post:   post,
			Errors: errorMessage,
		})

	}
}

func fetchPostDetails(postID int, logged bool) (Post, error) {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return Post{}, err
	}
	defer db.Close()

	query := `
		SELECT p.PostID, p.Title, p.Content, u.Username, p.CreatedAt, 
			GROUP_CONCAT(c.CategoryName, ', ') as Categories,
			(SELECT COUNT(*) FROM Likes WHERE PostID = p.PostID AND IsLike = 1) as LikeCount,
			(SELECT COUNT(*) FROM Dislikes WHERE PostID = p.PostID AND IsDislike = 1) as DislikeCount,
			(SELECT COUNT(*) FROM Comments WHERE PostID = p.PostID) as CommentCount
		FROM Posts p 
		JOIN Users u ON p.UserID = u.UserID 
		LEFT JOIN PostCategories pc ON p.PostID = pc.PostID
		LEFT JOIN Categories c ON pc.CategoryID = c.CategoryID
		WHERE p.PostID = ?
	`

	var post Post
	var categories sql.NullString
	err = db.QueryRow(query, postID).Scan(&post.PostID, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &categories, &post.LikeCount, &post.DislikeCount, &post.CommentsCount)
	if err != nil {
		return Post{}, err
	}

	if categories.Valid {
		post.Categories = strings.Split(categories.String, ", ")
	}

	post.Logged = logged

	commentsQuery := `
	SELECT c.CommentID, c.Content, u.Username, c.CreatedAt,
		(SELECT COUNT(*) FROM Reactions WHERE CommentID = c.CommentID AND IsLike = 1) as LikeCount,
		(SELECT COUNT(*) FROM Reactions WHERE CommentID = c.CommentID AND IsLike = 0) as DislikeCount
	FROM Comments c
	JOIN Users u ON c.UserID = u.UserID
	WHERE c.PostID = ?
	ORDER BY c.CreatedAt DESC
	`
	rows, err := db.Query(commentsQuery, postID)
	if err != nil {
		return Post{}, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.CommentID, &comment.Content, &comment.Username, &comment.CreatedAt, &comment.LikeCount, &comment.DislikeCount); err != nil {
			return Post{}, err
		}
		comment.Logged = logged
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return Post{}, err
	}

	post.Comments = comments

	return post, nil
}

func isIdValid(postID int) bool {
	db, err := sql.Open("sqlite3", "./database/forum.db")
	if err != nil {
		return false
	}
	defer db.Close()

	query := `
		SELECT p.PostID, p.Title, p.Content, p.CreatedAt
		FROM Posts p 
		WHERE p.PostID = ?
	`

	var post Post
	err = db.QueryRow(query, postID).Scan(&post.PostID, &post.Title, &post.Content, &post.CreatedAt)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}
	return true
}
