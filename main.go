package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/database"
	forum "forum/handlers"

	// other necessary imports

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()
	// Serve static files (e.g., CSS)
	statichandler := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", statichandler))
	mux.HandleFunc("/", forum.HandleRequest)
	mux.HandleFunc("/register", forum.RegHandler)
	mux.HandleFunc("/home", forum.HomeHandler)
	mux.HandleFunc("/profile", forum.ProfileHandler)
	mux.HandleFunc("/login", forum.LoginHandler)
	mux.HandleFunc("/error", forum.RenderErrorPage)
	mux.HandleFunc("/logout", forum.LogoutHandler)
	mux.HandleFunc("/posts", forum.PostsHandler)
	mux.HandleFunc("/createPost", forum.CreatePostHandler)
	mux.HandleFunc("/toggle-like", forum.ToggleLikeHandler)
	mux.HandleFunc("/toggle-dislike", forum.ToggleDislikeHandler)
	mux.HandleFunc("/toggle-comment-reaction", forum.ToggleCommentReactionHandler)
	mux.HandleFunc("/posts/", forum.PostDetailsHandler)
	mux.HandleFunc("/posts/{id}/comment", forum.CreateCommentHandler) 

	// http.HandleFunc("/like_dislike", forum.LikeDislikeHandler)
	// http.HandleFunc("/comment", forum.CommentHandler)

	s := &http.Server{
		Addr:    ":8989",
		Handler: mux,
	}
	fmt.Println("Server is running on port  http://localhost:8989")
	log.Fatal(s.ListenAndServe())
}
