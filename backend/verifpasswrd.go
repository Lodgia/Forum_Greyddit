package backend

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/mattn/go-sqlite3"
)

func VerifyPassword(user, password string) int {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, name, password FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	passwordHash := bcrypt.HashPassword(password, bcrypt.DefaultCost)
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

func GetUserId(name string) int {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT id FROM users WHERE name = ?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}
		return id
	}
	return -1
}

func CreateUser(name, password string) int {
	db, err := sql.Open("sqlite3", "./base.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var existe bool
    err = db.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM utilisateurs WHERE name = ?)",
        name,
    ).Scan(&existe)
	if err != nil {
		panic(err)
	}
	if existe {
		return -1 // L'utilisateur existe déjà
	}

	passwordHash := bcrypt.HashPassword(password, bcrypt.DefaultCost)
	_, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", name, passwordHash)
	if err != nil {
		panic(err)
	}
	return GetUserId(name)
}

func maintestpwd() {
	fmt.Println(VerifyPassword("Alice", "password1")) // Devrait afficher 1
	fmt.Println(VerifyPassword("Bob", "password2"))   // Devrait afficher 2
	fmt.Println(VerifyPassword("Alice", "wrongpass")) // Devrait afficher -1
}