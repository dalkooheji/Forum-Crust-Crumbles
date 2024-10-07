package forum

import (
	"time"
)

// User represents a user in the forum
type User struct {
	Username     string
	Email        string
	PasswordHash string
}

type Post struct {
	PostID        int       `json:"postID"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"createdAt"`
	Categories    []string  `json:"categories"`
	Comments      []Comment `json:"comments"`
	CommentsCount int       `json:"commentsCount"`
	LikeCount     int       `json:"likeCount"`
	DislikeCount  int       `json:"dislikeCount"`
	Logged        bool      `json:"logged"`
}

type Comment struct {
	CommentID    int       `json:"commentID"`
	PostID       int       `json:"postID"`
	Username     string    `json:"username"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"createdAt"`
	LikeCount    int       `json:"likeCount"`
	DislikeCount int       `json:"dislikeCount"`
	Logged       bool      `json:"logged"`
}

// Define the Category struct
type Category struct {
	CategoryID   int
	CategoryName string
}

type ErrorPage struct {
	Title       string
	Code        int
	Description string
}
