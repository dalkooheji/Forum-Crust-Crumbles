package forum

import (
	"database/sql"
	"forum/database"
	"html/template"
	"log"
	"net/http"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {

    _, errr := GetSession(r)
	logged := true
	if errr != nil {
		logged = false
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
		// return
	}

    db := database.DB

    category := r.URL.Query().Get("category")

    var rows *sql.Rows
    var err error
    if category != "" {
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
        WHERE c.CategoryName = ?
        GROUP BY p.PostID
        ORDER BY p.CreatedAt DESC`
        rows, err = db.Query(query, category)
    } else {
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
        GROUP BY p.PostID
        ORDER BY p.CreatedAt DESC`
        rows, err = db.Query(query)
    }

    if err != nil {
        log.Printf("Failed to fetch posts: %v", err)
        renderErrorTemplate(w, http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        var categories sql.NullString
        if err := rows.Scan(&post.PostID, &post.Title, &post.Content, &post.Username, &post.CreatedAt, &categories, &post.LikeCount, &post.DislikeCount, &post.CommentsCount); err != nil {
            http.Error(w, "Failed to scan post", http.StatusInternalServerError)
            return
        }
        if categories.Valid {
            post.Categories = append(post.Categories, categories.String)
        }
        post.Logged = logged
        posts = append(posts, post)
    }

    if category != "" {
        tmpl := template.Must(template.ParseFiles("templates/post_category.html"))
        tmpl.ExecuteTemplate(w, "post_category.html", posts)
    } else {
        tmpl := template.Must(template.ParseFiles("templates/posts.html"))
        tmpl.ExecuteTemplate(w, "posts.html", posts)
    }
}
