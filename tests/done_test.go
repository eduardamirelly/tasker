package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindTaskById(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("find existing task", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a test task
		taskID := insertTestTask(t, "Test Task", "Test Description", false)

		// Test finding task by ID directly in database
		query := `SELECT id, title, description, done, created_at FROM tasks WHERE id = ?`
		var id int
		var title, description string
		var done bool
		var createdAt time.Time

		err := database.DB.QueryRow(query, taskID).Scan(&id, &title, &description, &done, &createdAt)
		require.NoError(t, err)

		assert.Equal(t, taskID, id)
		assert.Equal(t, "Test Task", title)
		assert.Equal(t, "Test Description", description)
		assert.False(t, done)
		assert.NotZero(t, createdAt)
	})

	t.Run("find non-existent task", func(t *testing.T) {
		clearTestTasks(t)

		// Try to find a task that doesn't exist
		query := `SELECT id, title, description, done FROM tasks WHERE id = ?`
		var id int
		var title, description string
		var done bool

		err := database.DB.QueryRow(query, 999).Scan(&id, &title, &description, &done)
		assert.Error(t, err) // Should be sql.ErrNoRows
	})

	t.Run("find completed task", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a completed task
		taskID := insertTestTask(t, "Completed Task", "This is done", true)

		// Set completed_at timestamp
		completedTime := time.Now()
		updateQuery := `UPDATE tasks SET completed_at = ? WHERE id = ?`
		_, err := database.DB.Exec(updateQuery, completedTime, taskID)
		require.NoError(t, err)

		// Find the task
		query := `SELECT id, title, done, completed_at FROM tasks WHERE id = ?`
		var id int
		var title string
		var done bool
		var completedAt time.Time

		err = database.DB.QueryRow(query, taskID).Scan(&id, &title, &done, &completedAt)
		require.NoError(t, err)

		assert.Equal(t, taskID, id)
		assert.True(t, done)
		assert.WithinDuration(t, completedTime, completedAt, time.Second)
	})
}

func TestMarkTaskAsDone(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("mark incomplete task as done", func(t *testing.T) {
		clearTestTasks(t)

		// Insert an incomplete task
		taskID := insertTestTask(t, "Todo Task", "Need to complete this", false)

		// Verify task is initially incomplete
		task := getTaskByID(t, taskID)
		require.NotNil(t, task)
		assert.False(t, task.Done)

		// Mark as done using direct database update (simulating the cmd function)
		completedTime := time.Now()
		updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
		_, err := database.DB.Exec(updateQuery, completedTime, taskID)
		require.NoError(t, err)

		// Verify task is now marked as done
		updatedTask := getTaskByID(t, taskID)
		require.NotNil(t, updatedTask)
		assert.True(t, updatedTask.Done)

		// Verify completed_at timestamp was set
		var completedAt *time.Time
		query := `SELECT completed_at FROM tasks WHERE id = ?`
		err = database.DB.QueryRow(query, taskID).Scan(&completedAt)
		require.NoError(t, err)
		assert.NotNil(t, completedAt)
		assert.WithinDuration(t, time.Now(), *completedAt, 5*time.Second)
	})

	t.Run("mark already completed task as done", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a completed task
		taskID := insertTestTask(t, "Already Done", "This is complete", true)

		// Set completed_at timestamp
		originalTime := time.Now().Add(-time.Hour)
		updateQuery := `UPDATE tasks SET completed_at = ? WHERE id = ?`
		_, err := database.DB.Exec(updateQuery, originalTime, taskID)
		require.NoError(t, err)

		// Verify task is completed
		task := getTaskByID(t, taskID)
		require.NotNil(t, task)
		assert.True(t, task.Done)

		// Try to mark as done again (update timestamp)
		newCompletedTime := time.Now()
		updateQuery = `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
		_, err = database.DB.Exec(updateQuery, newCompletedTime, taskID)
		require.NoError(t, err)

		// Verify task is still marked as done and timestamp was updated
		var completedAt time.Time
		query := `SELECT completed_at FROM tasks WHERE id = ?`
		err = database.DB.QueryRow(query, taskID).Scan(&completedAt)
		require.NoError(t, err)

		// The timestamp should be updated
		assert.WithinDuration(t, time.Now(), completedAt, 5*time.Second)
	})
}

func TestDoneCommandIntegration(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("complete workflow - find and mark task as done", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a test task
		taskID := insertTestTask(t, "Integration Test Task", "Complete workflow test", false)

		// Step 1: Find the task
		query := `SELECT id, title, description, done FROM tasks WHERE id = ?`
		var id int
		var title, description string
		var done bool

		err := database.DB.QueryRow(query, taskID).Scan(&id, &title, &description, &done)
		require.NoError(t, err)
		assert.Equal(t, taskID, id)
		assert.False(t, done)

		// Step 2: Mark as done
		completedTime := time.Now()
		updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
		_, err = database.DB.Exec(updateQuery, completedTime, taskID)
		require.NoError(t, err)

		// Step 3: Verify the task is now completed
		err = database.DB.QueryRow(query, taskID).Scan(&id, &title, &description, &done)
		require.NoError(t, err)
		assert.True(t, done)
	})

	t.Run("try to complete non-existent task", func(t *testing.T) {
		clearTestTasks(t)

		// Try to find a non-existent task
		query := `SELECT id FROM tasks WHERE id = ?`
		var id int
		err := database.DB.QueryRow(query, 999).Scan(&id)
		assert.Error(t, err) // Should be sql.ErrNoRows

		// No tasks should exist in database
		count := getTaskCount(t)
		assert.Equal(t, 0, count)
	})
}

func TestMultipleTaskOperations(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Create multiple tasks
	taskIDs := make([]int, 5)
	for i := 0; i < 5; i++ {
		title := fmt.Sprintf("Task %d", i+1)
		description := fmt.Sprintf("Description for task %d", i+1)
		taskIDs[i] = insertTestTask(t, title, description, false)
	}

	// Mark every other task as done
	completedTime := time.Now()
	for i := 0; i < len(taskIDs); i += 2 {
		updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
		_, err := database.DB.Exec(updateQuery, completedTime, taskIDs[i])
		require.NoError(t, err)
	}

	// Verify the correct tasks are marked as done
	for i, taskID := range taskIDs {
		task := getTaskByID(t, taskID)
		require.NotNil(t, task)

		if i%2 == 0 {
			assert.True(t, task.Done, "Task %d should be completed", i+1)
		} else {
			assert.False(t, task.Done, "Task %d should not be completed", i+1)
		}
	}
}

func TestDoneWithDatabaseError(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Insert a test task
	taskID := insertTestTask(t, "Test Task", "Test Description", false)

	// Close the database to simulate an error
	database.DB.Close()

	// Try to find task - should get an error
	query := `SELECT id FROM tasks WHERE id = ?`
	var id int
	err := database.DB.QueryRow(query, taskID).Scan(&id)
	assert.Error(t, err)
}
