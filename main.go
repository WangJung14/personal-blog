package main

import (
	"blog/handlers"
	"blog/storage"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it")
	}

	// Initialize DB
	storage.InitDB()

	// Ensure uploads directory exists
	os.MkdirAll("static/uploads", 0755)

	// Parse templates
	tmpl := template.New("")
	tmpl, err = tmpl.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing public templates: %v", err)
	}
	tmpl, err = tmpl.ParseGlob("templates/admin/*.html")
	if err != nil {
		log.Fatalf("Error parsing admin templates: %v", err)
	}
	handlers.SetTemplates(tmpl)

	// Setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Public routes
	r.Get("/", handlers.PublicIndex)
	r.Get("/post/{slug}", handlers.PublicPost)

	// Admin routes
	r.Route("/admin", func(r chi.Router) {
		// Unprotected
		r.Get("/login", handlers.AdminLoginGet)
		r.Post("/login", handlers.AdminLoginPost)

		// Protected
		r.Group(func(r chi.Router) {
			r.Use(handlers.AdminAuthMiddleware)

			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			})
			r.Get("/dashboard", handlers.AdminDashboard)
			r.Get("/logout", handlers.AdminLogout)

			r.Get("/post/new", handlers.AdminPostNewGet)
			r.Post("/post/new", handlers.AdminPostNewPost)

			r.Get("/post/{id}/edit", handlers.AdminPostEditGet)
			r.Post("/post/{id}/edit", handlers.AdminPostEditPost)

			r.Post("/post/{id}/delete", handlers.AdminPostDeletePost)

			r.Get("/authors", handlers.AdminAuthorsGet)
			r.Post("/authors", handlers.AdminAuthorsPost)
			r.Post("/authors/{id}/delete", handlers.AdminAuthorsDelete)

			r.Post("/upload", handlers.AdminUploadImage)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
