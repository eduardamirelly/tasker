package cmd

import (
	"database/sql"
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
		defer result.Close()
		for result.Next() {
			var task models.Task
			err := result.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.CompletedAt)
			if err != nil {
				fmt.Printf("Error scanning task: %v\n", err)
				return
			}
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listTasks() (*sql.Rows, error) {
	query := `SELECT id, title, description, done, created_at, completed_at FROM tasks`
	result, err := database.DB.Query(query)
	if err != nil {
		fmt.Printf("Error listing tasks: %v\n", err)
		return nil, err
	}
	return result, err
}
