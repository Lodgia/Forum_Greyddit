package middleware

import (
	"context"
	"Forum_Greyddit/database"
	"Forum_Greyddit/models"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"


func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			
			next.ServeHTTP(w, r)
			return
		}

		session, err := database.GetSession(cookie.Value)
		if err != nil {
			
			http.SetCookie(w, &http.Cookie{
				Name:   "session_id",
				Value:  "",
				MaxAge: -1,
				Path:   "/",
			})
			next.ServeHTTP(w, r)
			return
		}

		user, err := database.GetUserByID(session.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetCurrentUser(r) == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func GetCurrentUser(r *http.Request) *models.User {
	u, ok := r.Context().Value(UserContextKey).(models.User)
	if !ok {
		return nil
	}
	return &u
}
