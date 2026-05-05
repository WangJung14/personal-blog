package storage

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}

	// Migration
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	// Add image_url column if not exists
	_, err = DB.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS image_url TEXT;`)
	if err != nil {
		log.Fatalf("Failed to add image_url column: %v", err)
	}

	// Authors Migration
	authorMigrationQuery := `
	CREATE TABLE IF NOT EXISTS authors (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL UNIQUE
	);
	INSERT INTO authors (name) VALUES ('Admin') ON CONFLICT (name) DO NOTHING;
	ALTER TABLE posts ADD COLUMN IF NOT EXISTS author_id INTEGER REFERENCES authors(id);
	UPDATE posts SET author_id = (SELECT id FROM authors WHERE name = 'Admin' LIMIT 1) WHERE author_id IS NULL;
	`
	_, err = DB.Exec(authorMigrationQuery)
	if err != nil {
		log.Fatalf("Failed to migrate authors: %v", err)
	}

	log.Println("Database connected and migrated successfully")
}
