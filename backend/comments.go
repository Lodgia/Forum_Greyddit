package backend

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)

type Comment struct {
	commentsId int
	commentsName  string
	id  int
	content string
}

var commentsId int = 1

func AddComment(authorid int, namecomment, content string) {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("INSERT INTO comments (commentsId, commentsName, id, content) VALUES (?, ?, ?, ?)", commentsId, namecomment, authorid, content)
	commentsId++
}

func GetComments(postID int) []Comment {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT commentsId, commentsName, id, content FROM comments WHERE id = ?", postID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.commentsId, &comment.commentsName, &comment.id, &comment.content)
		if err != nil {
			panic(err)
		}
		comments = append(comments, comment)
	}
	return comments
}

func DeleteComment(commentID int, userlogged int) {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM comments WHERE commentsId = ? AND id = ?", commentID, userlogged)
	if err != nil {
		panic(err)
	}
}

func maintestcomment() {
	AddComment(1, "Alice", "Ceci est un commentaire d'Alice")
	AddComment(2, "Bob", "Ceci est un commentaire de Bob")

	comments := GetComments(1)
	for _, comment := range comments {
		fmt.Printf("Commentaire ID: %d, Auteur: %s, Contenu: %s\n", comment.commentsId, comment.commentsName, comment.content)
	}

	DeleteComment(1, 1) // Supprime le commentaire d'Alice
	DeleteComment(2, 2) // Ne supprime pas le commentaire de Bob car l'utilisateur connecté n'est pas l'auteur
}