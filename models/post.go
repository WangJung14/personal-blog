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
	rows, err := storage.DB.Query("SELECT id, title, slug, content, COALESCE(image_url, ''), created_at, updated_at FROM posts ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}

func GetBySlug(slug string) (Post, error) {
	var p Post
	err := storage.DB.QueryRow("SELECT id, title, slug, content, COALESCE(image_url, ''), created_at, updated_at FROM posts WHERE slug = $1", slug).
		Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return p, ErrPostNotFound
	}
	return p, err
}

func GetByID(id int) (Post, error) {
	var p Post
	err := storage.DB.QueryRow("SELECT id, title, slug, content, COALESCE(image_url, ''), created_at, updated_at FROM posts WHERE id = $1", id).
		Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return p, ErrPostNotFound
	}
	return p, err
}

func Create(title, content, imageURL string) (int, error) {
	baseSlug := generateSlug(title)
	slug := ensureUniqueSlug(baseSlug, 0)

	var id int
	err := storage.DB.QueryRow(
		"INSERT INTO posts (title, slug, content, image_url) VALUES ($1, $2, $3, $4) RETURNING id",
		title, slug, content, imageURL,
	).Scan(&id)

	return id, err
}

func Update(id int, title, content, imageURL string) error {
	baseSlug := generateSlug(title)
	slug := ensureUniqueSlug(baseSlug, id)

	res, err := storage.DB.Exec(
		"UPDATE posts SET title = $1, slug = $2, content = $3, image_url = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5",
		title, slug, content, imageURL, id,
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
