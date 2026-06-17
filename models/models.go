package models

import "time"


type User struct {
	ID        string
	Email     string
	Username  string
	Password  string 
	CreatedAt time.Time
}


type Post struct {
	ID         string
	UserID     string
	Username   string 
	Title      string
	Content    string
	Categories []Category
	Likes      int
	Dislikes   int
	UserVote   int 
	CommentCount int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}


type Comment struct {
	ID        string
	PostID    string
	UserID    string
	Username  string 
	Content   string
	Likes     int
	Dislikes  int
	UserVote  int
	CreatedAt time.Time
}


type Category struct {
	ID   int
	Name string
	Slug string
}


type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}
