package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// Config de l'application
type Application struct {
	DB *sql.DB
}

func main() {
	// 1. Connexion à la base de données
	connStr := "user=postgres password=secret dbname=reddit_clone sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Impossible de se connecter à la DB: %v", err)
	}
	defer db.Close()

	app := &Application{DB: db}

	// 2. Définition des routes
	mux := http.NewServeMux()
	
	// Servir les fichiers statiques (CSS, JS)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Routes applicatives
	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/r/", app.subredditHandler)
	mux.HandleFunc("/post/create", app.postCreateHandler)

	// 3. Lancement du serveur
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Serveur lancé sur http://localhost:8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// Un handler de test pour la page d'accueil
func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Bienvenue sur le clone de Reddit !"))
}

func (app *Application) subredditHandler(w http.ResponseWriter, r *http.Request) {}
func (app *Application) postCreateHandler(w http.ResponseWriter, r *http.Request) {}