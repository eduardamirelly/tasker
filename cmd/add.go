/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/eduardamirelly/tasker/database"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new task",
	Long: `Add a new task to your task list. 

Examples:
  tasker add "Buy groceries"
  tasker add "Finish project" --description "Complete the final report"`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		description, _ := cmd.Flags().GetString("description")

		if err := addTask(title, description); err != nil {
			fmt.Printf("Error adding task: %v\n", err)
			return
		}

		fmt.Printf("✓ Task added: %s\n", title)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Add description flag
	addCmd.Flags().StringP("description", "d", "", "Task description")
}

// addTask adds a new task to the database
func addTask(title, description string) error {
	query := `INSERT INTO tasks (title, description, created_at) VALUES (?, ?, ?)`
	_, err := database.DB.Exec(query, title, description, time.Now())
	return err
}
