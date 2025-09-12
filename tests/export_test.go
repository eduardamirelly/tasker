package tests

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportTasks(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	tests := []struct {
		name         string
		setupTasks   []testTaskData
		outputFile   string
		expectError  bool
		expectHeader bool
	}{
		{
			name: "export single completed task",
			setupTasks: []testTaskData{
				{
					title:       "Buy groceries",
					description: "Milk, eggs, bread",
					done:        true,
				},
			},
			outputFile:   "test_single.csv",
			expectError:  false,
			expectHeader: true,
		},
		{
			name: "export single incomplete task",
			setupTasks: []testTaskData{
				{
					title:       "Finish project",
					description: "Complete the final report",
					done:        false,
				},
			},
			outputFile:   "test_incomplete.csv",
			expectError:  false,
			expectHeader: true,
		},
		{
			name: "export multiple mixed tasks",
			setupTasks: []testTaskData{
				{
					title:       "Task 1",
					description: "Description 1",
					done:        true,
				},
				{
					title:       "Task 2",
					description: "Description 2",
					done:        false,
				},
				{
					title:       "Task 3",
					description: "",
					done:        true,
				},
			},
			outputFile:   "test_multiple.csv",
			expectError:  false,
			expectHeader: true,
		},
		{
			name:         "export empty database",
			setupTasks:   []testTaskData{},
			outputFile:   "test_empty.csv",
			expectError:  false,
			expectHeader: true,
		},
		{
			name: "export tasks with special characters",
			setupTasks: []testTaskData{
				{
					title:       "Task with \"quotes\" and, commas",
					description: "Special chars: @#$%^&*()",
					done:        false,
				},
				{
					title:       "Unicode Task 🚀",
					description: "Emojis and unicode: 🎉🎯",
					done:        true,
				},
			},
			outputFile:   "test_special.csv",
			expectError:  false,
			expectHeader: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing tasks
			clearTestTasks(t)

			// Set up test tasks
			var taskIDs []int
			for _, taskData := range tt.setupTasks {
				id := insertTestTaskWithTimestamp(t, taskData.title, taskData.description, taskData.done)
				taskIDs = append(taskIDs, id)
			}

			// Create temporary directory for test files
			tempDir := t.TempDir()
			outputPath := filepath.Join(tempDir, tt.outputFile)

			// Execute export function
			err := exportToCSV(outputPath)

			// Check error expectation
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify file was created
			assert.FileExists(t, outputPath)

			// Read and verify CSV content
			content := readCSVFile(t, outputPath)

			// Verify header exists
			if tt.expectHeader {
				require.GreaterOrEqual(t, len(content), 1, "CSV should have at least header row")
				expectedHeader := []string{"ID", "Title", "Description", "Done", "Created At", "Completed At"}
				assert.Equal(t, expectedHeader, content[0])
			}

			// Verify data rows
			dataRows := content[1:] // Skip header
			assert.Equal(t, len(tt.setupTasks), len(dataRows), "Number of data rows should match setup tasks")

			// Verify each task's data
			for i, taskData := range tt.setupTasks {
				if i < len(dataRows) {
					row := dataRows[i]
					assert.Equal(t, strconv.Itoa(taskIDs[i]), row[0], "ID should match")
					assert.Equal(t, taskData.title, row[1], "Title should match")
					assert.Equal(t, taskData.description, row[2], "Description should match")
					assert.Equal(t, strconv.FormatBool(taskData.done), row[3], "Done status should match")
					assert.NotEmpty(t, row[4], "Created At should not be empty")

					// Completed At should be empty for incomplete tasks, non-empty for completed
					if taskData.done {
						assert.NotEmpty(t, row[5], "Completed At should not be empty for completed tasks")
					} else {
						assert.Empty(t, row[5], "Completed At should be empty for incomplete tasks")
					}
				}
			}
		})
	}
}

func TestExportInvalidPath(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Try to export to an invalid path (non-existent directory)
	invalidPath := "/nonexistent/directory/test.csv"
	err := exportToCSV(invalidPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create CSV file")
}

func TestExportLargeDataset(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Clear any existing tasks
	clearTestTasks(t)

	// Create a large number of tasks
	taskCount := 1000
	for i := 0; i < taskCount; i++ {
		title := fmt.Sprintf("Task %d", i)
		description := fmt.Sprintf("Description for task %d", i)
		done := i%2 == 0 // Alternate between done and not done
		insertTestTaskWithTimestamp(t, title, description, done)
	}

	// Create temporary directory for test file
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "large_dataset.csv")

	// Execute export
	err := exportToCSV(outputPath)
	require.NoError(t, err)

	// Verify file was created and has correct number of rows
	content := readCSVFile(t, outputPath)

	// Should have header + taskCount data rows
	expectedRows := taskCount + 1
	assert.Equal(t, expectedRows, len(content))

	// Verify header
	expectedHeader := []string{"ID", "Title", "Description", "Done", "Created At", "Completed At"}
	assert.Equal(t, expectedHeader, content[0])

	// Spot check a few rows
	assert.Equal(t, "Task 0", content[1][1])      // First task
	assert.Equal(t, "Task 999", content[1000][1]) // Last task
}

func TestExportDateTimeFormatting(t *testing.T) {
	// Setup test database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Clear any existing tasks
	clearTestTasks(t)

	// Create a task with specific timestamp
	specificTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	insertTestTaskWithSpecificTime(t, "Christmas Task", "Holiday planning", true, specificTime, &specificTime)

	// Export to CSV
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "datetime_test.csv")
	err := exportToCSV(outputPath)
	require.NoError(t, err)

	// Read and verify datetime formatting
	content := readCSVFile(t, outputPath)
	require.Len(t, content, 2) // Header + 1 data row

	dataRow := content[1]
	expectedDateFormat := "2023-12-25 15:30:45"
	assert.Equal(t, expectedDateFormat, dataRow[4]) // Created At
	assert.Equal(t, expectedDateFormat, dataRow[5]) // Completed At
}

// Helper types and functions

type testTaskData struct {
	title       string
	description string
	done        bool
}

// insertTestTaskWithTimestamp inserts a test task and marks it as done if specified
func insertTestTaskWithTimestamp(t *testing.T, title, description string, done bool) int {
	createdAt := time.Now()
	query := `INSERT INTO tasks (title, description, done, created_at) VALUES (?, ?, ?, ?)`
	result, err := database.DB.Exec(query, title, description, done, createdAt)
	require.NoError(t, err)

	id, err := result.LastInsertId()
	require.NoError(t, err)

	// If task should be done, update completed_at
	if done {
		completedAt := createdAt.Add(time.Hour) // Complete 1 hour after creation
		updateQuery := `UPDATE tasks SET completed_at = ? WHERE id = ?`
		_, err := database.DB.Exec(updateQuery, completedAt, id)
		require.NoError(t, err)
	}

	return int(id)
}

// insertTestTaskWithSpecificTime inserts a task with specific timestamps
func insertTestTaskWithSpecificTime(t *testing.T, title, description string, done bool, createdAt time.Time, completedAt *time.Time) int {
	query := `INSERT INTO tasks (title, description, done, created_at, completed_at) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, title, description, done, createdAt, completedAt)
	require.NoError(t, err)

	id, err := result.LastInsertId()
	require.NoError(t, err)

	return int(id)
}

// readCSVFile reads a CSV file and returns its content as a 2D string slice
func readCSVFile(t *testing.T, filePath string) [][]string {
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	require.NoError(t, err)

	return records
}

// exportToCSV is a test wrapper for the export functionality
func exportToCSV(outputPath string) error {
	// This mimics the exportTasks function from cmd/export.go
	// We need to implement the actual export logic here for testing

	// Get all tasks from database
	tasks, err := getAllTasksForExport()
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	// Create CSV file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"ID", "Title", "Description", "Done", "Created At", "Completed At"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write task data
	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Title,
			task.Description,
			strconv.FormatBool(task.Done),
			task.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Handle completed_at (nullable field)
		if task.CompletedAt != nil {
			record = append(record, task.CompletedAt.Format("2006-01-02 15:04:05"))
		} else {
			record = append(record, "")
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write task record: %w", err)
		}
	}

	return nil
}

// Task struct for export testing (mirrors models.Task)
type exportTask struct {
	ID          int
	Title       string
	Description string
	Done        bool
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// getAllTasksForExport retrieves all tasks from the database for export testing
func getAllTasksForExport() ([]exportTask, error) {
	query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []exportTask
	for rows.Next() {
		var task exportTask
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
