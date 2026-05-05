package handlers

import (
	"blog/models"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func PublicIndex(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAll()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "index.html", map[string]interface{}{"Posts": posts})
}

func PublicPost(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	post, err := models.GetBySlug(slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Create a safe HTML content
	safeContent := template.HTML(post.Content)

	tmpl.ExecuteTemplate(w, "post.html", map[string]interface{}{
		"Post":        post,
		"SafeContent": safeContent,
	})
}
