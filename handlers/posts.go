package handlers

import (
	"greyddit/database"
	"greyddit/middleware"
	"greyddit/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		renderError(w, r, http.StatusNotFound, "Page introuvable.")
		return
	}
	user := middleware.GetCurrentUser(r)
	filter := r.URL.Query().Get("filter")
	userID := ""
	if user != nil {
		userID = user.ID
	} else {
		filter = ""
	}
	posts, err := database.GetPosts(filter, userID)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Impossible de charger les posts.")
		return
	}
	cats, _ := database.GetAllCategories()
	renderTemplate(w, r, "index.html", map[string]interface{}{
		"Posts": posts, "Categories": cats, "Filter": filter,
	})
}

func ByCategory(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/category/")
	if slug == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	posts, err := database.GetPostsByCategory(slug)
	if err != nil {
		renderError(w, r, http.StatusInternalServerError, "Catégorie introuvable.")
		return
	}
	cats, _ := database.GetAllCategories()
	renderTemplate(w, r, "index.html", map[string]interface{}{
		"Posts": posts, "Categories": cats, "Category": slug,
	})
}

func ViewPost(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/post/create") ||
		strings.HasPrefix(path, "/post/edit") ||
		strings.HasPrefix(path, "/post/delete") {
		http.NotFound(w, r)
		return
	}
	id := strings.TrimPrefix(path, "/post/")
	if id == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	post, err := database.GetPostByID(id)
	if err != nil {
		renderError(w, r, http.StatusNotFound, "Post introuvable.")
		return
	}
	user := middleware.GetCurrentUser(r)
	userID := ""
	if user != nil {
		userID = user.ID
		post.UserVote, _ = database.GetUserVoteOnPost(userID, id)
	}
	comments, _ := database.GetCommentsForPost(id, userID)
	renderTemplate(w, r, "post.html", map[string]interface{}{
		"Post": post, "Comments": comments,
	})
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	cats, _ := database.GetAllCategories()
	if r.Method == http.MethodGet {
		renderTemplate(w, r, "create.html", map[string]interface{}{"Categories": cats})
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	catIDs := parseIntList(r.Form["categories"])
	if title == "" || content == "" {
		renderTemplate(w, r, "create.html", map[string]interface{}{
			"Categories": cats, "Error": "Le titre et le contenu sont requis.",
		})
		return
	}
	post := models.Post{ID: uuid.NewString(), UserID: user.ID, Title: title, Content: content}
	if err := database.CreatePost(post, catIDs); err != nil {
		renderError(w, r, http.StatusInternalServerError, "Impossible de créer le post.")
		return
	}
	http.Redirect(w, r, "/post/"+post.ID, http.StatusSeeOther)
}

func EditPost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	id := strings.TrimPrefix(r.URL.Path, "/post/edit/")
	post, err := database.GetPostByID(id)
	if err != nil {
		renderError(w, r, http.StatusNotFound, "Post introuvable.")
		return
	}
	if post.UserID != user.ID {
		renderError(w, r, http.StatusForbidden, "Vous ne pouvez modifier que vos propres posts.")
		return
	}
	if r.Method == http.MethodGet {
		renderTemplate(w, r, "edit.html", map[string]interface{}{"Post": post})
		return
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	if title == "" || content == "" {
		renderTemplate(w, r, "edit.html", map[string]interface{}{
			"Post": post, "Error": "Le titre et le contenu sont requis.",
		})
		return
	}
	database.UpdatePost(id, title, content)
	http.Redirect(w, r, "/post/"+id, http.StatusSeeOther)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	id := strings.TrimPrefix(r.URL.Path, "/post/delete/")
	post, err := database.GetPostByID(id)
	if err != nil {
		renderError(w, r, http.StatusNotFound, "Post introuvable.")
		return
	}
	if post.UserID != user.ID {
		renderError(w, r, http.StatusForbidden, "Vous ne pouvez supprimer que vos propres posts.")
		return
	}
	database.DeletePost(id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func parseIntList(values []string) []int {
	var result []int
	for _, v := range values {
		n, err := strconv.Atoi(v)
		if err == nil {
			result = append(result, n)
		}
	}
	return result
}
