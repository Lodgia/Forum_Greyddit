package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Ouvre (ou crée) la base de données
	fmt.Println("Avant Open")
	db, err := sql.Open("sqlite3", "./demo.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Après Open")
	fmt.Println("Avant Ping")

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}	
	fmt.Println("Après Ping")

	// Création de la table
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Table créée.")

	// Insertion de données
	insertQuery := `
	INSERT INTO users (name, age)
	VALUES (?, ?)
	`

	users := []struct {
		Name string
		Age  int
	}{
		{"Alice", 25},
		{"Bob", 32},
		{"Charlie", 41},
	}

	for _, user := range users {
		_, err := db.Exec(insertQuery, user.Name, user.Age)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Données insérées.")

	// Lecture des données
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("\nContenu de la table users :")

	for rows.Next() {
		var id int
		var name string
		var age int

		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID=%d, Name=%s, Age=%d\n", id, name, age)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}