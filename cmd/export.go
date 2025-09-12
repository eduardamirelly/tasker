package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/eduardamirelly/tasker/database"
	"github.com/eduardamirelly/tasker/models"
	"github.com/spf13/cobra"
)

var (
	outputFile string
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export tasks to CSV",
	Long:  `Export tasks to CSV file.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := exportTasks()
		if err != nil {
			fmt.Printf("Error exporting tasks: %v\n", err)
			return
		}
		fmt.Printf("Tasks exported successfully to %s\n", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "tasks.csv", "Output CSV file path")
}

func exportTasks() error {
	// Get all tasks from database
	tasks, err := getAllTasks()
	if err != nil {
		return fmt.Errorf("failed to fetch tasks: %w", err)
	}

	// Create CSV file
	file, err := os.Create(outputFile)
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

// getAllTasks retrieves all tasks from the database
func getAllTasks() ([]models.Task, error) {
	query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
