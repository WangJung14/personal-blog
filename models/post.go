package models

import (
	"blog/storage"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	ImageURL  string    `json:"image_url"`
	AuthorID  int       `json:"author_id"`
	AuthorName string   `json:"author_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrPostNotFound = errors.New("post not found")

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "post"
	}
	return slug
}

func ensureUniqueSlug(baseSlug string, excludeID int) string {
	slug := baseSlug
	counter := 2
	for {
		var id int
		err := storage.DB.QueryRow("SELECT id FROM posts WHERE slug = $1 AND id != $2", slug, excludeID).Scan(&id)
		if err == sql.ErrNoRows {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return slug
}

func GetAll() ([]Post, error) {
	rows, err := storage.DB.Query(`
		SELECT p.id, p.title, p.slug, p.content, COALESCE(p.image_url, ''), p.author_id, a.name, p.created_at, p.updated_at 
		FROM posts p 
		JOIN authors a ON p.author_id = a.id 
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetBySlug(slug string) (Post, error) {
	var p Post
	err := storage.DB.QueryRow(`
		SELECT p.id, p.title, p.slug, p.content, COALESCE(p.image_url, ''), p.author_id, a.name, p.created_at, p.updated_at 
		FROM posts p 
		JOIN authors a ON p.author_id = a.id 
		WHERE p.slug = $1
	`, slug).Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return p, ErrPostNotFound
	}
	return p, err
}

func GetByID(id int) (Post, error) {
	var p Post
	err := storage.DB.QueryRow(`
		SELECT p.id, p.title, p.slug, p.content, COALESCE(p.image_url, ''), p.author_id, a.name, p.created_at, p.updated_at 
		FROM posts p 
		JOIN authors a ON p.author_id = a.id 
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.AuthorID, &p.AuthorName, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return p, ErrPostNotFound
	}
	return p, err
}

func Create(title, content, imageURL string, authorID int) (int, error) {
	baseSlug := generateSlug(title)
	slug := ensureUniqueSlug(baseSlug, 0)

	var id int
	err := storage.DB.QueryRow(
		"INSERT INTO posts (title, slug, content, image_url, author_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		title, slug, content, imageURL, authorID,
	).Scan(&id)

	return id, err
}

func Update(id int, title, content, imageURL string, authorID int) error {
	baseSlug := generateSlug(title)
	slug := ensureUniqueSlug(baseSlug, id)

	res, err := storage.DB.Exec(
		"UPDATE posts SET title = $1, slug = $2, content = $3, image_url = $4, author_id = $5, updated_at = CURRENT_TIMESTAMP WHERE id = $6",
		title, slug, content, imageURL, authorID, id,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrPostNotFound
	}
	return nil
}

func Delete(id int) error {
	res, err := storage.DB.Exec("DELETE FROM posts WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return ErrPostNotFound
	}
	return nil
}
