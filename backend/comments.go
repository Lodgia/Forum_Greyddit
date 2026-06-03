package backend

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
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