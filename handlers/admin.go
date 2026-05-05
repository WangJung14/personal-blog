package handlers

import (
	"blog/models"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

var tmpl *template.Template

func SetTemplates(t *template.Template) {
	tmpl = t
}

func AdminLoginGet(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	expectedUser := os.Getenv("ADMIN_USER")
	expectedPass := os.Getenv("ADMIN_PASS")

	if username != expectedUser || password != expectedPass {
		tmpl.ExecuteTemplate(w, "login.html", map[string]string{"Error": "Invalid credentials"})
		return
	}

	secret := []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		Subject:   "admin",
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func AdminLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func AdminDashboard(w http.ResponseWriter, r *http.Request) {
	posts, err := models.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "dashboard.html", map[string]interface{}{"Posts": posts})
}

func AdminPostNewGet(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "editor.html", nil)
}

func AdminPostNewPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]string{"Error": "Title and content are required"})
		return
	}

	_, err := models.Create(title, content)
	if err != nil {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": err.Error(), "Title": title, "Content": template.HTML(content)})
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func AdminPostEditGet(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	post, err := models.GetByID(id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Post": post, "SafeContent": template.HTML(post.Content)})
}

func AdminPostEditPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		post, _ := models.GetByID(id)
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": "Title and content are required", "Post": post})
		return
	}

	err := models.Update(id, title, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func AdminPostDeletePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	err := models.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}
