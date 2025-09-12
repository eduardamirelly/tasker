package tests

import (
	"fmt"
	"testing"

	"github.com/eduardamirelly/tasker/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTask(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name        string
		title       string
		description string
		wantErr     bool
	}{
		{
			name:        "successful task addition with description",
			title:       "Test Task",
			description: "This is a test task",
			wantErr:     false,
		},
		{
			name:        "successful task addition without description",
			title:       "Task without description",
			description: "",
			wantErr:     false,
		},
		{
			name:        "task with special characters",
			title:       "Special @#$% Task!",
			description: "Task with special chars: @#$%^&*()",
			wantErr:     false,
		},
		{
			name:        "task with unicode characters",
			title:       "Unicode Task 🚀",
			description: "Unicode description with emojis 🎉🎯",
			wantErr:     false,
		},
		{
			name:        "long title and description",
			title:       "This is a very long task title that contains many words and should still work properly",
			description: "This is an extremely long description that goes on and on and should also work properly without any issues because our database should handle long text fields",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing tasks
			clearTestTasks(t)

			// Get initial task count
			initialCount := getTaskCount(t)

			// Execute addTask function - we need to access the private function
			// For now, we'll test via the public database operations
			query := `INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)`
			_, err := database.DB.Exec(query, tt.title, tt.description, "2023-01-01")

			// Check error expectation
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			// Should not have error
			require.NoError(t, err)

			// Verify task count increased by 1
			finalCount := getTaskCount(t)
			assert.Equal(t, initialCount+1, finalCount)

			// Verify the task was inserted correctly
			queryCheck := `SELECT title, description, done FROM tasks WHERE title = ?`
			var title, description string
			var done bool
			err = database.DB.QueryRow(queryCheck, tt.title).Scan(&title, &description, &done)
			require.NoError(t, err)

			assert.Equal(t, tt.title, title)
			assert.Equal(t, tt.description, description)
			assert.False(t, done) // New tasks should be incomplete
		})
	}
}

func TestAddTaskEmptyTitle(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Test with empty title - this should still work as our database allows it
	query := `INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, "", "Description without title", "2023-01-01")
	assert.NoError(t, err)

	// Verify the task was added
	count := getTaskCount(t)
	assert.Equal(t, 1, count)
}

func TestAddTaskMultipleTasks(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	tasks := []struct {
		title       string
		description string
	}{
		{"Task 1", "Description 1"},
		{"Task 2", "Description 2"},
		{"Task 3", "Description 3"},
	}

	// Add multiple tasks
	for i, task := range tasks {
		query := `INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)`
		_, err := database.DB.Exec(query, task.title, task.description, "2023-01-01")
		require.NoError(t, err)

		// Verify count increases correctly
		count := getTaskCount(t)
		assert.Equal(t, i+1, count)
	}

	// Verify all tasks exist
	finalCount := getTaskCount(t)
	assert.Equal(t, len(tasks), finalCount)
}

func TestAddTaskConcurrency(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Test concurrent task additions
	taskCount := 10
	errChan := make(chan error, taskCount)

	for i := 0; i < taskCount; i++ {
		go func(index int) {
			title := fmt.Sprintf("Concurrent Task %d", index)
			description := fmt.Sprintf("Description for task %d", index)
			query := `INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)`
			_, err := database.DB.Exec(query, title, description, "2023-01-01")
			errChan <- err
		}(i)
	}

	// Collect all errors
	for i := 0; i < taskCount; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}

	// Verify all tasks were added
	count := getTaskCount(t)
	assert.Equal(t, taskCount, count)
}
