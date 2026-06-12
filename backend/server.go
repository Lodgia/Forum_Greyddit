package backend

import (
	"html/template"
	"net/http"
	"database/sql"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// Config de l'application
type Application struct {
	DB *sql.DB
}
// Un handler de test pour la page d'accueil
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//tmplPath := filepath.Join("ui", "template", "register.html")

	tmpl, err := template.ParseFiles("ui/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("ui/templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {}