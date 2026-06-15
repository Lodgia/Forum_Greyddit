package database

import (
	"greyddit/models"
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