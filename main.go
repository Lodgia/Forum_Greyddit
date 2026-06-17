package main

import (
	"Forum_Greyddit/database"
	"Forum_Greyddit/handlers"
	"Forum_Greyddit/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	
	database.Init("forum.db")
	defer database.DB.Close()

	
	go func() {
		for {
			time.Sleep(1 * time.Hour)
			database.CleanExpiredSessions()
		}
	}()

	mux := http.NewServeMux()

	
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	
	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/post/", handlers.ViewPost)       
	mux.HandleFunc("/category/", handlers.ByCategory) 
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/logout", handlers.Logout)

	
	mux.Handle("/post/create", middleware.RequireAuth(http.HandlerFunc(handlers.CreatePost)))
	mux.Handle("/post/edit/", middleware.RequireAuth(http.HandlerFunc(handlers.EditPost)))
	mux.Handle("/post/delete/", middleware.RequireAuth(http.HandlerFunc(handlers.DeletePost)))
	mux.Handle("/comment/create", middleware.RequireAuth(http.HandlerFunc(handlers.CreateComment)))
	mux.Handle("/comment/delete/", middleware.RequireAuth(http.HandlerFunc(handlers.DeleteComment)))
	mux.Handle("/vote/post", middleware.RequireAuth(http.HandlerFunc(handlers.VotePost)))
	mux.Handle("/vote/comment", middleware.RequireAuth(http.HandlerFunc(handlers.VoteComment)))

	
	withSession := middleware.SessionMiddleware(mux)

	log.Println("Greyddit démarre sur http://localhost:5050")
	if err := http.ListenAndServe(":5050", withSession); err != nil {
		log.Fatal(err)
	}
}
