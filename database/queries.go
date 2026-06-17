package database

import (
	"Forum_Greyddit/models"
	"time"
)



func CreateUser(u models.User) error {
	_, err := DB.Exec(
		`INSERT INTO users (id, email, username, password) VALUES (?, ?, ?, ?)`,
		u.ID, u.Email, u.Username, u.Password,
	)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	row := DB.QueryRow(`SELECT id, email, username, password FROM users WHERE email = ?`, email)
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Username, &u.Password)
	return u, err
}

func GetUserByID(id string) (models.User, error) {
	row := DB.QueryRow(`SELECT id, email, username FROM users WHERE id = ?`, id)
	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Username)
	return u, err
}



func CreateSession(s models.Session) error {
	_, err := DB.Exec(
		`INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)`,
		s.ID, s.UserID, s.ExpiresAt,
	)
	return err
}

func GetSession(id string) (models.Session, error) {
	row := DB.QueryRow(
		`SELECT id, user_id, expires_at FROM sessions WHERE id = ? AND expires_at > CURRENT_TIMESTAMP`,
		id,
	)
	var s models.Session
	err := row.Scan(&s.ID, &s.UserID, &s.ExpiresAt)
	return s, err
}

func DeleteSession(id string) error {
	_, err := DB.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

func CleanExpiredSessions() {
	DB.Exec(`DELETE FROM sessions WHERE expires_at <= CURRENT_TIMESTAMP`)
}



func GetAllCategories() ([]models.Category, error) {
	rows, err := DB.Query(`SELECT id, name, slug FROM categories ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name, &c.Slug)
		cats = append(cats, c)
	}
	return cats, nil
}



func CreatePost(p models.Post, categoryIDs []int) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`INSERT INTO posts (id, user_id, title, content) VALUES (?, ?, ?, ?)`,
		p.ID, p.UserID, p.Title, p.Content,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, cid := range categoryIDs {
		_, err = tx.Exec(
			`INSERT INTO post_categories (post_id, category_id) VALUES (?, ?)`,
			p.ID, cid,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func GetPosts(filter string, userID string) ([]models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       COALESCE(SUM(CASE WHEN pl.value = 1 THEN 1 ELSE 0 END), 0) AS likes,
		       COALESCE(SUM(CASE WHEN pl.value = -1 THEN 1 ELSE 0 END), 0) AS dislikes,
		       COUNT(DISTINCT c.id) AS comment_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		LEFT JOIN comments c ON p.id = c.post_id`

	args := []interface{}{}

	switch filter {
	case "mine":
		query += ` WHERE p.user_id = ?`
		args = append(args, userID)
	case "liked":
		query += ` WHERE p.id IN (SELECT post_id FROM post_likes WHERE user_id = ? AND value = 1)`
		args = append(args, userID)
	}

	query += ` GROUP BY p.id ORDER BY p.created_at DESC`

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		rows.Scan(
			&post.ID, &post.UserID, &post.Username,
			&post.Title, &post.Content, &post.CreatedAt,
			&post.Likes, &post.Dislikes, &post.CommentCount,
		)
		post.Categories, _ = GetCategoriesForPost(post.ID)
		if userID != "" {
			post.UserVote, _ = GetUserVoteOnPost(userID, post.ID)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetPostsByCategory(slug string) ([]models.Post, error) {
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       COALESCE(SUM(CASE WHEN pl.value = 1 THEN 1 ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN pl.value = -1 THEN 1 ELSE 0 END), 0),
		       COUNT(DISTINCT c.id)
		FROM posts p
		JOIN users u ON p.user_id = u.id
		JOIN post_categories pc ON p.id = pc.post_id
		JOIN categories cat ON pc.category_id = cat.id
		LEFT JOIN post_likes pl ON p.id = pl.post_id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE cat.slug = ?
		GROUP BY p.id ORDER BY p.created_at DESC`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt,
			&p.Likes, &p.Dislikes, &p.CommentCount)
		posts = append(posts, p)
	}
	return posts, nil
}

func GetPostByID(id string) (models.Post, error) {
	row := DB.QueryRow(`
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at, p.updated_at
		FROM posts p JOIN users u ON p.user_id = u.id WHERE p.id = ?`, id)
	var p models.Post
	err := row.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return p, err
	}
	p.Categories, _ = GetCategoriesForPost(p.ID)
	return p, nil
}

func UpdatePost(id, title, content string) error {
	_, err := DB.Exec(
		`UPDATE posts SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		title, content, id,
	)
	return err
}

func DeletePost(id string) error {
	_, err := DB.Exec(`DELETE FROM posts WHERE id = ?`, id)
	return err
}

func GetCategoriesForPost(postID string) ([]models.Category, error) {
	rows, err := DB.Query(`
		SELECT cat.id, cat.name, cat.slug FROM categories cat
		JOIN post_categories pc ON cat.id = pc.category_id
		WHERE pc.post_id = ?`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []models.Category
	for rows.Next() {
		var c models.Category
		rows.Scan(&c.ID, &c.Name, &c.Slug)
		cats = append(cats, c)
	}
	return cats, nil
}



func CreateComment(c models.Comment) error {
	_, err := DB.Exec(
		`INSERT INTO comments (id, post_id, user_id, content) VALUES (?, ?, ?, ?)`,
		c.ID, c.PostID, c.UserID, c.Content,
	)
	return err
}

func GetCommentsForPost(postID, userID string) ([]models.Comment, error) {
	rows, err := DB.Query(`
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
		       COALESCE(SUM(CASE WHEN cl.value = 1 THEN 1 ELSE 0 END), 0),
		       COALESCE(SUM(CASE WHEN cl.value = -1 THEN 1 ELSE 0 END), 0)
		FROM comments c
		JOIN users u ON c.user_id = u.id
		LEFT JOIN comment_likes cl ON c.id = cl.comment_id
		WHERE c.post_id = ?
		GROUP BY c.id ORDER BY c.created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var cm models.Comment
		rows.Scan(&cm.ID, &cm.PostID, &cm.UserID, &cm.Username,
			&cm.Content, &cm.CreatedAt, &cm.Likes, &cm.Dislikes)
		if userID != "" {
			cm.UserVote, _ = GetUserVoteOnComment(userID, cm.ID)
		}
		comments = append(comments, cm)
	}
	return comments, nil
}

func DeleteComment(id string) error {
	_, err := DB.Exec(`DELETE FROM comments WHERE id = ?`, id)
	return err
}




func VotePost(userID, postID string, value int) error {
	existing, _ := GetUserVoteOnPost(userID, postID)
	if existing == value {
		
		_, err := DB.Exec(`DELETE FROM post_likes WHERE user_id = ? AND post_id = ?`, userID, postID)
		return err
	}
	_, err := DB.Exec(
		`INSERT INTO post_likes (user_id, post_id, value) VALUES (?, ?, ?)
		 ON CONFLICT(user_id, post_id) DO UPDATE SET value = excluded.value`,
		userID, postID, value,
	)
	return err
}

func VoteComment(userID, commentID string, value int) error {
	existing, _ := GetUserVoteOnComment(userID, commentID)
	if existing == value {
		_, err := DB.Exec(`DELETE FROM comment_likes WHERE user_id = ? AND comment_id = ?`, userID, commentID)
		return err
	}
	_, err := DB.Exec(
		`INSERT INTO comment_likes (user_id, comment_id, value) VALUES (?, ?, ?)
		 ON CONFLICT(user_id, comment_id) DO UPDATE SET value = excluded.value`,
		userID, commentID, value,
	)
	return err
}

func GetUserVoteOnPost(userID, postID string) (int, error) {
	var v int
	err := DB.QueryRow(`SELECT value FROM post_likes WHERE user_id = ? AND post_id = ?`, userID, postID).Scan(&v)
	return v, err
}

func GetUserVoteOnComment(userID, commentID string) (int, error) {
	var v int
	err := DB.QueryRow(`SELECT value FROM comment_likes WHERE user_id = ? AND comment_id = ?`, userID, commentID).Scan(&v)
	return v, err
}


func CountUsers() int {
	var count int
	DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count
}


var _ = time.Now
