package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

  //nom de la base de données

func HashPassword(password string) string {
	return password
}

func VerifyPassword(user, password string) int {
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	passwordHash := HashPassword(password)
	for rows.Next() {
    	var id int
		var name string
		var password string

		err := rows.Scan(&id, &name, &password)
		if err != nil {
			panic(err)
		}
		if name == user && password == passwordHash {
			return id
		}
	}
	return -1
}

func main() {
	db, err = sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Création de la table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		password TEXT NOT NULL
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err)
	}

	// Insertion de données
	insertQuery := `
	INSERT INTO users (id, name, password)
	VALUES (?, ?, ?)
	`

	users := []struct {
		Id       int
		Name     string
		Password string
	}{
		{1, "Alice", HashPassword("password1")},
		{2, "Bob", HashPassword("password2")},
	}

	for _, user := range users {
		_, err := db.Exec(insertQuery, user.Id, user.Name, user.Password)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(VerifyPassword("Alice", "password1")) // Devrait afficher 1
	fmt.Println(VerifyPassword("Bob", "password2"))   // Devrait afficher 2
	fmt.Println(VerifyPassword("Alice", "wrongpass")) // Devrait afficher -1
}