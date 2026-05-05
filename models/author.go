package models

import (
	"blog/storage"
	"errors"
)

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllAuthors() ([]Author, error) {
	rows, err := storage.DB.Query("SELECT id, name FROM authors ORDER BY name ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []Author
	for rows.Next() {
		var a Author
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	return authors, nil
}

func CreateAuthor(name string) error {
	if name == "" {
		return errors.New("name is required")
	}
	_, err := storage.DB.Exec("INSERT INTO authors (name) VALUES ($1)", name)
	return err
}

func DeleteAuthor(id int) error {
	var defaultID int
	err := storage.DB.QueryRow("SELECT id FROM authors WHERE name = 'Admin' LIMIT 1").Scan(&defaultID)
	if err != nil {
		return err
	}

	if id == defaultID {
		return errors.New("cannot delete default Admin author")
	}

	// Reassign posts to default author
	storage.DB.Exec("UPDATE posts SET author_id = $1 WHERE author_id = $2", defaultID, id)
	_, err = storage.DB.Exec("DELETE FROM authors WHERE id = $1", id)
	return err
}
