package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

var DB *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	// Get current working directory (project root)
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create database file path in project directory
	dbPath := filepath.Join(currentDir, "tasker.db")

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	DB = db

	// Create tasks table if it doesn't exist
	return createTables()
}

// createTables creates the necessary database tables
func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		done BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME
	);`

	_, err := DB.Exec(query)
	return err
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
