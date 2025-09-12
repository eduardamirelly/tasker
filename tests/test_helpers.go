package tests

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/eduardamirelly/tasker/database"
	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates a temporary test database for testing
func setupTestDB(t *testing.T) func() {
	// Create a temporary database file
	tempDir := t.TempDir()
	testDBPath := filepath.Join(tempDir, "test_tasker.db")

	// Open test database connection
	db, err := sql.Open("sqlite3", testDBPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Store original DB and replace with test DB
	originalDB := database.DB
	database.DB = db

	// Create tables in test database
	err = createTestTables()
	if err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Return cleanup function
	return func() {
		database.DB.Close()
		database.DB = originalDB
		os.Remove(testDBPath)
	}
}

// createTestTables creates the necessary database tables for testing
func createTestTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		done BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		completed_at DATETIME
	);`

	_, err := database.DB.Exec(query)
	return err
}

// insertTestTask is a helper function to insert a test task
func insertTestTask(t *testing.T, title, description string, done bool) int {
	query := `INSERT INTO tasks (title, description, done) VALUES (?, ?, ?)`
	result, err := database.DB.Exec(query, title, description, done)
	if err != nil {
		t.Fatalf("Failed to insert test task: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Failed to get last insert ID: %v", err)
	}

	return int(id)
}

// getTaskCount returns the number of tasks in the database
func getTaskCount(t *testing.T) int {
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to get task count: %v", err)
	}
	return count
}

// getTaskByID retrieves a task by its ID for testing
func getTaskByID(t *testing.T, id int) *testTask {
	query := `SELECT id, title, description, done FROM tasks WHERE id = ?`
	var task testTask
	err := database.DB.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.Done)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		t.Fatalf("Failed to get task by ID: %v", err)
	}
	return &task
}

// testTask is a simplified task struct for testing
type testTask struct {
	ID          int
	Title       string
	Description string
	Done        bool
}

// clearTestTasks removes all tasks from the test database
func clearTestTasks(t *testing.T) {
	_, err := database.DB.Exec("DELETE FROM tasks")
	if err != nil {
		t.Fatalf("Failed to clear test tasks: %v", err)
	}
}
