package main

import (
    "database/sql"
    "html/template"
    "log"
    "net/http"
    "time"
    "golang.org/x/crypto/bcrypt" // pour hacher les mots de passe
    _ "github.com/lib/pq"
)

type Application struct {
    DB        *sql.DB
    Templates map[string]*template.Template
}

func main() {
    connStr := "user=postgres password=secret dbname=greyddit sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Pré-charger les templates HTML
    templates := map[string]*template.Template{
        "login":    template.Must(template.ParseFiles("ui/templates/login.html")),
        "register": template.Must(template.ParseFiles("ui/templates/register.html")),
    }

    app := &Application{DB: db, Templates: templates}

    mux := http.NewServeMux()

    // Fichiers statiques (CSS, images)
    fs := http.FileServer(http.Dir("./ui/static/"))
    mux.Handle("/static/", http.StripPrefix("/static", fs))

    // Routes HTML (afficher les pages)
    mux.HandleFunc("/login",    app.loginHandler)
    mux.HandleFunc("/register", app.registerHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        IdleTimeout:  time.Minute,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    log.Println("Serveur lancé sur http://localhost:8080")
    log.Fatal(srv.ListenAndServe())
}

// ─── Handler Login ──────────────────────────────────────────────────────────

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Afficher la page login
        app.Templates["login"].Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        // Récupérer les données du formulaire
        username := r.FormValue("userID")
        password := r.FormValue("pwd")

        // Chercher l'utilisateur en base
        var passwordHash string
        err := app.DB.QueryRow(
            "SELECT password_hash FROM users WHERE username = $1",
            username,
        ).Scan(&passwordHash)

        if err == sql.ErrNoRows {
            http.Error(w, "Identifiant ou mot de passe incorrect", http.StatusUnauthorized)
            return
        }
        if err != nil {
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        // Vérifier le mot de passe
        err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
        if err != nil {
            http.Error(w, "Identifiant ou mot de passe incorrect", http.StatusUnauthorized)
            return
        }

        // Succès → rediriger vers l'accueil
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

// ─── Handler Register ────────────────────────────────────────────────────────

func (app *Application) registerHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        app.Templates["register"].Execute(w, nil)
        return
    }

    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        email    := r.FormValue("email")
        password := r.FormValue("newpwd")

        // Hacher le mot de passe (NE JAMAIS stocker en clair)
        hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
        if err != nil {
            http.Error(w, "Erreur serveur", http.StatusInternalServerError)
            return
        }

        // Insérer en base de données
        _, err = app.DB.Exec(
            "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
            username, email, string(hash),
        )
        if err != nil {
            http.Error(w, "Nom d'utilisateur ou email déjà utilisé", http.StatusConflict)
            return
        }

        // Succès → rediriger vers login
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    }
}