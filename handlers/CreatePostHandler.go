package forum

import (
    "forum/database"
    "html/template"
    "net/http"
    "strings"
    "time"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the user is authenticated
    session, err := GetSession(r)
    if err != nil {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    db := database.DB

    if r.Method == http.MethodPost {
        title := strings.TrimSpace(r.FormValue("title"))
        content := strings.TrimSpace(r.FormValue("content"))
        categoryIDs := r.Form["categories"]

        // Title validation
        if title == "" {
            http.Error(w, "Title cannot be empty", http.StatusBadRequest)
            return
        }
        if len(title) > 100 { // Set a title length limit
            http.Error(w, "Title must be less than 100 characters", http.StatusBadRequest)
            return
        }

        // Content validation (check for empty content, ignoring spaces or newlines)
        if content == "" {
            http.Error(w, "Content cannot be empty", http.StatusBadRequest)
            return
        }

        // Ensure at least one category is selected
        if len(categoryIDs) == 0 {
            http.Error(w, "You must select at least one category", http.StatusBadRequest)
            return
        }

        // Insert the post into the database
        result, err := db.Exec(`INSERT INTO Posts (UserID, Title, Content, CreatedAt) VALUES (?, ?, ?, ?)`,
            session.UserID, title, content, time.Now())
        if err != nil {
            http.Error(w, "Failed to create post", http.StatusInternalServerError)
            return
        }

        postID, err := result.LastInsertId()
        if err != nil {
            http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
            return
        }

        // Insert the selected categories for the post
        for _, categoryID := range categoryIDs {
            _, err := db.Exec(`INSERT INTO PostCategories (PostID, CategoryID) VALUES (?, ?)`, postID, categoryID)
            if err != nil {
                http.Error(w, "Failed to assign category", http.StatusInternalServerError)
                return
            }
        }

        http.Redirect(w, r, "/posts", http.StatusSeeOther)
    } else {
        // Fetch available categories
        rows, err := db.Query("SELECT CategoryID, CategoryName FROM Categories")
        if err != nil {
            http.Error(w, "Failed to load categories", http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var categories []Category
        for rows.Next() {
            var category Category
            if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
                http.Error(w, "Failed to scan category", http.StatusInternalServerError)
                return
            }
            categories = append(categories, category)
        }

        tmpl := template.Must(template.ParseFiles("templates/createPost.html"))
        tmpl.Execute(w, struct {
            Categories []Category
        }{
            Categories: categories,
        })
    }
}
