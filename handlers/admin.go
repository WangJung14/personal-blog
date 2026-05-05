package handlers

import (
	"blog/models"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	authors, _ := models.GetAllAuthors()
	tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Authors": authors})
}

func handleCoverImageUpload(r *http.Request) (string, error) {
	file, header, err := r.FormFile("cover_image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil // No file uploaded
		}
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("cover_%d%s", time.Now().UnixNano(), ext)
	uploadPath := filepath.Join("static", "uploads", filename)

	dst, err := os.Create(uploadPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return "/static/uploads/" + filename, nil
}

func AdminPostNewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]string{"Error": "Form parse error"})
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]string{"Error": "Title and content are required"})
		return
	}

	imageURL, err := handleCoverImageUpload(r)
	if err != nil {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": "Image upload failed", "Title": title, "Content": template.HTML(content)})
		return
	}

	authorIDStr := r.FormValue("author_id")
	authorID, _ := strconv.Atoi(authorIDStr)

	_, err = models.Create(title, content, imageURL, authorID)
	if err != nil {
		authors, _ := models.GetAllAuthors()
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": err.Error(), "Title": title, "Content": template.HTML(content), "Authors": authors})
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

	authors, _ := models.GetAllAuthors()
	tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Post": post, "Authors": authors, "SafeContent": template.HTML(post.Content)})
}

func AdminPostEditPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	err := r.ParseMultipartForm(10 << 20)
	authors, _ := models.GetAllAuthors()
	if err != nil {
		post, _ := models.GetByID(id)
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": "Form parse error", "Post": post, "Authors": authors})
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	authorIDStr := r.FormValue("author_id")
	authorID, _ := strconv.Atoi(authorIDStr)

	post, _ := models.GetByID(id)

	if title == "" || content == "" {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": "Title and content are required", "Post": post, "Authors": authors})
		return
	}

	imageURL, err := handleCoverImageUpload(r)
	if err != nil {
		tmpl.ExecuteTemplate(w, "editor.html", map[string]interface{}{"Error": "Image upload failed", "Post": post})
		return
	}
	if imageURL == "" {
		imageURL = post.ImageURL
	}

	err = models.Update(id, title, content, imageURL, authorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
}

func AdminAuthorsGet(w http.ResponseWriter, r *http.Request) {
	authors, _ := models.GetAllAuthors()
	tmpl.ExecuteTemplate(w, "authors.html", map[string]interface{}{"Authors": authors})
}

func AdminAuthorsPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name != "" {
		models.CreateAuthor(name)
	}
	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
}

func AdminAuthorsDelete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	models.DeleteAuthor(id)
	http.Redirect(w, r, "/admin/authors", http.StatusSeeOther)
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

func AdminUploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate a unique filename using timestamp
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	uploadPath := filepath.Join("static", "uploads", filename)

	dst, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	// Return URL
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url": "/static/uploads/" + filename,
	})
}
