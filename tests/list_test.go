package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListTasks(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	t.Run("empty task list", func(t *testing.T) {
		clearTestTasks(t)

		// Use reflection or create a wrapper to access the private listTasks function
		// For now, we'll test the database query directly
		query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
		rows, err := database.DB.Query(query)
		require.NoError(t, err)
		defer rows.Close()

		taskCount := 0
		for rows.Next() {
			taskCount++
		}
		assert.Equal(t, 0, taskCount)
	})

	t.Run("single task", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a test task
		insertTestTask(t, "Test Task", "Test Description", false)

		// Query tasks
		query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
		rows, err := database.DB.Query(query)
		require.NoError(t, err)
		defer rows.Close()

		taskCount := 0
		var title, description string
		var done bool
		var id int
		var createdAt, completedAt interface{}

		for rows.Next() {
			err := rows.Scan(&id, &title, &description, &done, &createdAt, &completedAt)
			require.NoError(t, err)
			taskCount++
		}

		assert.Equal(t, 1, taskCount)
		assert.Equal(t, "Test Task", title)
		assert.Equal(t, "Test Description", description)
		assert.False(t, done)
		assert.NotNil(t, createdAt)
		assert.Nil(t, completedAt)
	})

	t.Run("multiple tasks", func(t *testing.T) {
		clearTestTasks(t)

		// Insert multiple test tasks
		testData := []struct {
			title       string
			description string
			done        bool
		}{
			{"Task 1", "Description 1", false},
			{"Task 2", "Description 2", true},
			{"Task 3", "", false},
		}

		for _, td := range testData {
			insertTestTask(t, td.title, td.description, td.done)
		}

		// Query and verify tasks
		query := `SELECT id, title, description, done FROM tasks ORDER BY id`
		rows, err := database.DB.Query(query)
		require.NoError(t, err)
		defer rows.Close()

		var tasks []testTask
		for rows.Next() {
			var task testTask
			err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done)
			require.NoError(t, err)
			tasks = append(tasks, task)
		}

		assert.Len(t, tasks, 3)
		for i, expectedTask := range testData {
			assert.Equal(t, expectedTask.title, tasks[i].Title)
			assert.Equal(t, expectedTask.description, tasks[i].Description)
			assert.Equal(t, expectedTask.done, tasks[i].Done)
		}
	})

	t.Run("tasks with completed_at timestamp", func(t *testing.T) {
		clearTestTasks(t)

		// Insert a completed task with completed_at timestamp
		completedTime := time.Now()
		query := `INSERT INTO tasks (title, description, done, completed_at) VALUES (?, ?, ?, ?)`
		_, err := database.DB.Exec(query, "Completed Task", "This task is done", true, completedTime)
		require.NoError(t, err)

		// Query and verify
		queryCheck := `SELECT completed_at FROM tasks WHERE title = ?`
		var retrievedCompletedAt time.Time
		err = database.DB.QueryRow(queryCheck, "Completed Task").Scan(&retrievedCompletedAt)
		require.NoError(t, err)

		// Check if completed time is approximately correct (within 1 second)
		assert.WithinDuration(t, completedTime, retrievedCompletedAt, time.Second)
	})

	t.Run("mixed completed and incomplete tasks", func(t *testing.T) {
		clearTestTasks(t)

		// Insert mix of completed and incomplete tasks
		insertTestTask(t, "Todo Task 1", "Not done yet", false)
		insertTestTask(t, "Completed Task 1", "Already done", true)
		insertTestTask(t, "Todo Task 2", "Still working", false)

		// Count tasks by status
		var completedCount, incompleteCount int

		err := database.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE done = true").Scan(&completedCount)
		require.NoError(t, err)

		err = database.DB.QueryRow("SELECT COUNT(*) FROM tasks WHERE done = false").Scan(&incompleteCount)
		require.NoError(t, err)

		assert.Equal(t, 1, completedCount)
		assert.Equal(t, 2, incompleteCount)
	})
}

func TestListTasksWithLargeDataset(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Insert many tasks to test performance and correctness
	taskCount := 100
	for i := 0; i < taskCount; i++ {
		title := fmt.Sprintf("Task %d", i+1)
		description := fmt.Sprintf("Description for task %d", i+1)
		done := i%3 == 0 // Every third task is completed
		insertTestTask(t, title, description, done)
	}

	// Verify all tasks were inserted
	count := getTaskCount(t)
	assert.Equal(t, taskCount, count)

	// Verify tasks have proper IDs (should be sequential)
	query := `SELECT id, title FROM tasks ORDER BY id`
	rows, err := database.DB.Query(query)
	require.NoError(t, err)
	defer rows.Close()

	i := 0
	for rows.Next() {
		var id int
		var title string
		err := rows.Scan(&id, &title)
		require.NoError(t, err)

		assert.Equal(t, i+1, id)
		assert.Contains(t, title, fmt.Sprintf("Task %d", i+1))
		i++
	}
}

func TestListTasksErrorHandling(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Close the database to simulate an error
	database.DB.Close()

	// Try to query tasks - should get an error
	query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
	rows, err := database.DB.Query(query)
	assert.Error(t, err)
	assert.Nil(t, rows)
}
