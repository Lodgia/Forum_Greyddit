package handlers

import (
	"greyddit/database"
	"/static/template"
	"golang.org/x/crypto/bcrypt"
)

var funcMap = template.FuncMap{
	"sub": func(a, b int) int { return a - b },
}

func Register(w http.ResponseWriter, r *http.Request) {
	if middleware.GetCurrentUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		renderTemplate(w, r, "register.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")
	if email == "" || username == "" || password == "" {
		renderTemplate(w, r, "register.html", map[string]string{"Error": "Tous les champs sont requis."})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Erreur serveur.")
		return
	}
	user := models.User{
		ID: uuid.NewString(), Email: email, Username: username, Password: string(hash),
	}
	if err = database.CreateUser(user); err != nil {
		renderTemplate(w, r, "register.html", map[string]string{"Error": "Email ou nom d'utilisateur déjà pris."})
		return
	}
	createSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}