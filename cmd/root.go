package cmd

import (
	"fmt"
	"os"

	"github.com/eduardamirelly/tasker/database"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tasker",
	Short: "A simple CLI task manager",
	Long: `Tasker is a command-line task management tool that helps you:
- Add new tasks
- List all tasks  
- Mark tasks as done
- Export tasks to CSV

Store your tasks locally in a SQLite database.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}

	// Ensure database is closed when program exits
	defer database.CloseDB()

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tasker.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
