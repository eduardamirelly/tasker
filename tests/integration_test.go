package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompleteTaskWorkflow tests the full workflow of adding, listing, and completing tasks
func TestCompleteTaskWorkflow(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Step 1: Start with empty task list
	count := getTaskCount(t)
	assert.Equal(t, 0, count, "Should start with no tasks")

	// Step 2: Add first task
	taskID1 := insertTestTask(t, "Buy groceries", "Milk, eggs, bread", false)

	// Step 3: Add second task
	taskID2 := insertTestTask(t, "Walk the dog", "30 minute walk in the park", false)

	// Step 4: Add task without description
	taskID3 := insertTestTask(t, "Call mom", "", false)

	// Step 5: List tasks and verify all were added
	count = getTaskCount(t)
	assert.Equal(t, 3, count, "Should have 3 tasks")

	// Verify task details
	expectedTasks := []struct {
		id          int
		title       string
		description string
		done        bool
	}{
		{taskID1, "Buy groceries", "Milk, eggs, bread", false},
		{taskID2, "Walk the dog", "30 minute walk in the park", false},
		{taskID3, "Call mom", "", false},
	}

	for _, expected := range expectedTasks {
		task := getTaskByID(t, expected.id)
		require.NotNil(t, task)
		assert.Equal(t, expected.title, task.Title)
		assert.Equal(t, expected.description, task.Description)
		assert.Equal(t, expected.done, task.Done)
	}

	// Step 6: Complete the first task (Buy groceries)
	completedTime := time.Now()
	updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
	_, err := database.DB.Exec(updateQuery, completedTime, taskID1)
	require.NoError(t, err)

	// Step 7: Verify task is marked as completed
	updatedTask := getTaskByID(t, taskID1)
	require.NotNil(t, updatedTask)
	assert.True(t, updatedTask.Done)

	// Step 8: List tasks again and verify completion status
	count = getTaskCount(t)
	assert.Equal(t, 3, count, "Should still have 3 tasks")

	// First task should be completed, others incomplete
	task1 := getTaskByID(t, taskID1)
	task2 := getTaskByID(t, taskID2)
	task3 := getTaskByID(t, taskID3)

	assert.True(t, task1.Done, "First task should be completed")
	assert.False(t, task2.Done, "Second task should be incomplete")
	assert.False(t, task3.Done, "Third task should be incomplete")

	// Step 9: Complete remaining tasks
	taskIDs := []int{taskID2, taskID3}
	for _, id := range taskIDs {
		completedTime := time.Now()
		_, err := database.DB.Exec(updateQuery, completedTime, id)
		require.NoError(t, err)
	}

	// Step 10: Verify all tasks are completed
	allTasks := []int{taskID1, taskID2, taskID3}
	for i, id := range allTasks {
		task := getTaskByID(t, id)
		require.NotNil(t, task)
		assert.True(t, task.Done, "Task %d should be completed", i+1)
	}
}

// TestTaskLifecycle tests the complete lifecycle of a single task
func TestTaskLifecycle(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	taskTitle := "Learn Go programming"
	taskDescription := "Complete Go tutorial and build a small project"

	// Phase 1: Creation
	taskID := insertTestTask(t, taskTitle, taskDescription, false)

	// Verify task was created
	count := getTaskCount(t)
	assert.Equal(t, 1, count)

	createdTask := getTaskByID(t, taskID)
	require.NotNil(t, createdTask)
	assert.Equal(t, taskTitle, createdTask.Title)
	assert.Equal(t, taskDescription, createdTask.Description)
	assert.False(t, createdTask.Done)

	// Phase 2: Finding task by ID
	query := `SELECT id, title, description, done FROM tasks WHERE id = ?`
	var id int
	var title, description string
	var done bool

	err := database.DB.QueryRow(query, taskID).Scan(&id, &title, &description, &done)
	require.NoError(t, err)

	assert.Equal(t, taskID, id)
	assert.Equal(t, taskTitle, title)
	assert.Equal(t, taskDescription, description)
	assert.False(t, done)

	// Phase 3: Completion
	completedTime := time.Now()
	updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
	_, err = database.DB.Exec(updateQuery, completedTime, taskID)
	require.NoError(t, err)

	// Verify completion
	err = database.DB.QueryRow(query, taskID).Scan(&id, &title, &description, &done)
	require.NoError(t, err)
	assert.True(t, done)

	// Phase 4: Verify in list
	count = getTaskCount(t)
	assert.Equal(t, 1, count)

	completedTask := getTaskByID(t, taskID)
	require.NotNil(t, completedTask)
	assert.Equal(t, taskID, completedTask.ID)
	assert.True(t, completedTask.Done)
}

// TestMixedTaskStates tests handling tasks in various states
func TestMixedTaskStates(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Create tasks with different characteristics
	testCases := []struct {
		title          string
		description    string
		shouldComplete bool
	}{
		{"Task with description", "This has a description", false},
		{"Task without description", "", true},
		{"Long task title that goes on and on", "Short desc", false},
		{"Unicode task 🚀", "Unicode description with emojis 🎉", true},
		{"Special chars !@#$%", "Description with ^&*()", false},
	}

	// Add all tasks
	taskIDs := make([]int, len(testCases))
	for i, tc := range testCases {
		taskIDs[i] = insertTestTask(t, tc.title, tc.description, false)
	}

	// Get initial list
	count := getTaskCount(t)
	assert.Equal(t, len(testCases), count)

	// All should be incomplete initially
	for _, id := range taskIDs {
		task := getTaskByID(t, id)
		require.NotNil(t, task)
		assert.False(t, task.Done)
	}

	// Complete selected tasks
	completedTime := time.Now()
	updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
	for i, id := range taskIDs {
		if testCases[i].shouldComplete {
			_, err := database.DB.Exec(updateQuery, completedTime, id)
			require.NoError(t, err)
		}
	}

	// Verify final state
	completedCount := 0
	incompleteCount := 0

	for i, id := range taskIDs {
		task := getTaskByID(t, id)
		require.NotNil(t, task)

		if testCases[i].shouldComplete {
			assert.True(t, task.Done, "Task %d should be completed", i)
			completedCount++
		} else {
			assert.False(t, task.Done, "Task %d should be incomplete", i)
			incompleteCount++
		}
	}

	assert.Equal(t, 2, completedCount, "Should have 2 completed tasks")
	assert.Equal(t, 3, incompleteCount, "Should have 3 incomplete tasks")
}

// TestErrorScenarios tests various error conditions in workflow
func TestErrorScenarios(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Scenario 1: Try to complete non-existent task
	query := `SELECT id FROM tasks WHERE id = ?`
	var id int
	err := database.DB.QueryRow(query, 999).Scan(&id)
	assert.Error(t, err, "Non-existent task should return error")

	// Scenario 2: Add task and then try various operations
	taskID := insertTestTask(t, "Test Task", "For error testing", false)

	count := getTaskCount(t)
	assert.Equal(t, 1, count)

	// Complete the task
	completedTime := time.Now()
	updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
	_, err = database.DB.Exec(updateQuery, completedTime, taskID)
	require.NoError(t, err)

	// Try to complete it again (should work without error)
	task := getTaskByID(t, taskID)
	require.NotNil(t, task)
	assert.True(t, task.Done, "Task should already be completed")

	// Update again (should not cause error)
	newCompletedTime := time.Now()
	_, err = database.DB.Exec(updateQuery, newCompletedTime, taskID)
	require.NoError(t, err)

	// Verify task is still completed
	updatedTask := getTaskByID(t, taskID)
	require.NotNil(t, updatedTask)
	assert.True(t, updatedTask.Done, "Task should still be completed")
}

// TestConcurrentOperations tests concurrent access to the database
func TestConcurrentOperations(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	clearTestTasks(t)

	// Add multiple tasks concurrently
	numTasks := 10
	errChan := make(chan error, numTasks)
	taskIDChan := make(chan int, numTasks)

	for i := 0; i < numTasks; i++ {
		go func(index int) {
			title := fmt.Sprintf("Concurrent Task %d", index)
			description := fmt.Sprintf("Description %d", index)

			query := `INSERT INTO tasks (title, description, done) VALUES (?, ?, ?)`
			result, err := database.DB.Exec(query, title, description, false)
			if err != nil {
				errChan <- err
				taskIDChan <- 0
				return
			}

			id, err := result.LastInsertId()
			errChan <- err
			taskIDChan <- int(id)
		}(i)
	}

	// Wait for all additions to complete
	var taskIDs []int
	for i := 0; i < numTasks; i++ {
		err := <-errChan
		id := <-taskIDChan
		assert.NoError(t, err, "Task addition %d should not error", i)
		if id > 0 {
			taskIDs = append(taskIDs, id)
		}
	}

	// Verify all tasks were added
	count := getTaskCount(t)
	assert.Equal(t, numTasks, count, "Should have all tasks")

	// Complete tasks concurrently
	completionErrChan := make(chan error, numTasks)
	completedTime := time.Now()
	updateQuery := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`

	for _, taskID := range taskIDs {
		go func(id int) {
			_, err := database.DB.Exec(updateQuery, completedTime, id)
			completionErrChan <- err
		}(taskID)
	}

	// Wait for all completions
	for i := 0; i < len(taskIDs); i++ {
		err := <-completionErrChan
		assert.NoError(t, err, "Task completion %d should not error", i)
	}

	// Verify all tasks are completed
	for i, id := range taskIDs {
		task := getTaskByID(t, id)
		require.NotNil(t, task)
		assert.True(t, task.Done, "Task %d should be completed", i)
	}
}
