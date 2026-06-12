package backend

import (
	"database/sql"
	"net/http"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// Config de l'application
type Application struct {
	DB *sql.DB
}
// Un handler de test pour la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Bienvenue sur le clone de Reddit !"))
}

func SubredditHandler(w http.ResponseWriter, r *http.Request) {}
func PostCreateHandler(w http.ResponseWriter, r *http.Request) {}