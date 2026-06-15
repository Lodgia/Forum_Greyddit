package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB


func Init(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Impossible d'ouvrir la base de données : %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Base de données inaccessible : %v", err)
	}

	applySchema()
	log.Println("Base de données initialisée :", path)
}

func applySchema() {
	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatalf("Impossible de lire schema.sql : %v", err)
	}
	if _, err = DB.Exec(string(schema)); err != nil {
		log.Fatalf("Erreur lors de l'application du schéma : %v", err)
	}
}
