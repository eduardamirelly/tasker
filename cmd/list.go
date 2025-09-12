package cmd

import (
	"fmt"

	"github.com/eduardamirelly/tasker/database"
	"github.com/eduardamirelly/tasker/models"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Long:  `List all tasks saved in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := listTasks()
		if err != nil {
			fmt.Printf("Error listing tasks: %v\n", err)
			return
		}
		if len(result) == 0 {
			emptyTasks()
			return
		}
		printTasks(result)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listTasks() ([]models.Task, error) {
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

func emptyTasks() {
	fmt.Println("No tasks found")
}

func printTasks(tasks []models.Task) {
	for _, task := range tasks {
		done := "✅"
		if !task.Done {
			done = "❌"
		}
		createdAt := task.CreatedAt.Format("2006-01-02 15:04:05")
		completedAt := "N/A"
		if task.CompletedAt != nil {
			completedAt = task.CompletedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Printf("%v %v - %v\n", done, task.ID, task.Title)
		fmt.Printf("Description: %v\n", task.Description)
		fmt.Printf("Created At: %v\n", createdAt)
		fmt.Printf("Completed At: %v\n", completedAt)
		fmt.Println("--------------------------------")
	}
}
