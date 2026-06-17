package handlers

import (
	"greyddit/database"
	"greyddit/middleware"
	"greyddit/models"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func Login(w http.ResponseWriter, r *http.Request) {
	if middleware.GetCurrentUser(r) != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		renderTemplate(w, r, "login.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := database.GetUserByEmail(email)
	if err != nil {
		renderTemplate(w, r, "login.html", map[string]string{"Error": "Email ou mot de passe incorrect."})
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		renderTemplate(w, r, "login.html", map[string]string{"Error": "Email ou mot de passe incorrect."})
		return
	}
	createSession(w, user.ID)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		database.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session_id", Value: "", MaxAge: -1, Path: "/"})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func createSession(w http.ResponseWriter, userID string) {
	session := models.Session{
		ID: uuid.NewString(), UserID: userID, ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := database.CreateSession(session); err != nil {
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "session_id", Value: session.ID, Expires: session.ExpiresAt,
		Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
}


func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	tmpl, err := template.New("").Funcs(funcMap).ParseFiles(
		"templates/base.html",
		"templates/"+name,
	)
	if err != nil {
		http.Error(w, "Erreur template : "+err.Error(), http.StatusInternalServerError)
		return
	}
	type PageData struct {
		CurrentUser interface{}
		Data        interface{}
	}
	if err = tmpl.ExecuteTemplate(w, "base", PageData{
		CurrentUser: middleware.GetCurrentUser(r),
		Data:        data,
	}); err != nil {
		http.Error(w, "Erreur rendu : "+err.Error(), http.StatusInternalServerError)
	}
}


func renderError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.WriteHeader(code)
	renderTemplate(w, r, "error.html", map[string]interface{}{
		"Code": code, "Message": msg,
	})
}
