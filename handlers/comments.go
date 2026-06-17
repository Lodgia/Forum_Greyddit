package handlers

import (
	"greyddit/database"
	"greyddit/middleware"
	"greyddit/models"
	"net/http"
	"strings"

	"github.com/google/uuid"
)


func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	postID := r.FormValue("post_id")
	content := r.FormValue("content")

	if postID == "" || content == "" {
		http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
		return
	}

	comment := models.Comment{
		ID:      uuid.NewString(),
		PostID:  postID,
		UserID:  user.ID,
		Content: content,
	}

	database.CreateComment(comment)
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}


func DeleteComment(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	id := strings.TrimPrefix(r.URL.Path, "/comment/delete/")
	postID := r.URL.Query().Get("post_id")

	comments, err := database.GetCommentsForPost(postID, user.ID)
	if err != nil {
		renderError(w, r, http.StatusNotFound, "Commentaire introuvable.")
		return
	}

	
	var found bool
	for _, c := range comments {
		if c.ID == id && c.UserID == user.ID {
			found = true
			break
		}
	}

	if !found {
		renderError(w, r, http.StatusForbidden, "Action non autorisée.")
		return
	}

	database.DeleteComment(id)
	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}


func VotePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	postID := r.FormValue("post_id")
	value := parseVoteValue(r.FormValue("value"))

	if postID == "" {
		http.Error(w, "post_id manquant", http.StatusBadRequest)
		return
	}

	database.VotePost(user.ID, postID, value)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}


func VoteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetCurrentUser(r)
	commentID := r.FormValue("comment_id")
	value := parseVoteValue(r.FormValue("value"))

	if commentID == "" {
		http.Error(w, "comment_id manquant", http.StatusBadRequest)
		return
	}

	database.VoteComment(user.ID, commentID, value)
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

func parseVoteValue(s string) int {
	if s == "1" {
		return 1
	}
	return -1
}
