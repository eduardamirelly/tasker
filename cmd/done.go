package cmd

import (
	"fmt"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/eduardamirelly/tasker/models"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done [id]",
	Short: "Mark a task as done",
	Long:  `Mark a task as done in the database.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]

		task, err := findTaskById(id)

		if err != nil {
			fmt.Printf("Error finding task: %v\n", err)
			return
		}

		if task.ID == 0 {
			fmt.Printf("❌ Task not found: %s\n", id)
			return
		}

		if task.Done {
			fmt.Printf("✅ Task already done!\n")
			printTask(task)
			return
		}

		markTaskAsDone(task)
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}

func findTaskById(id string) (*models.Task, error) {
	query := `SELECT id, title, description, done, created_at, completed_at FROM tasks WHERE id = ?`
	rows, err := database.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var task models.Task
	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Done, &task.CreatedAt, &task.CompletedAt)
		if err != nil {
			return nil, err
		}
	}
	return &task, nil
}

func markTaskAsDone(task *models.Task) {
	if task == nil {
		fmt.Printf("❌ Task not found!\n")
		return
	}

	query := `UPDATE tasks SET done = TRUE, completed_at = ? WHERE id = ?`
	_, err := database.DB.Exec(query, time.Now(), task.ID)
	if err != nil {
		fmt.Printf("Error marking task as done: %v\n", err)
		return
	}
	fmt.Printf("✓ Task marked as done: %s\n", task.Title)
	printTask(task)
}

func printTask(task *models.Task) {
	if task.Description == "" {
		task.Description = "N/A"
	}
	fmt.Println("--------------------------------")
	fmt.Printf("Title: %s\n", task.Title)
	fmt.Printf("Description: %s\n", task.Description)
	fmt.Printf("Created At: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Completed At: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("--------------------------------")
}
