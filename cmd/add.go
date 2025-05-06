package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add NAME",
	Short: "Add a new task",
	Long:  `Add a new task to your list. You can optionally assign it to a project.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbConn == nil {
			return fmt.Errorf("database connection not initialized")
		}
		taskName := args[0]
		// Get project flag value
		project, _ := cmd.Flags().GetString("project")

		// Use the exported Insert method from the db package via dbConn
		newTask, err := dbConn.Insert(taskName, project)
		if err != nil {
			return fmt.Errorf("failed to add task: %w", err)
		}
		fmt.Printf("Task ('%s') added.\n", newTask.Name)
		return nil
	},
}

// init registers flags specific to the add command.
func init() {
	addCmd.Flags().StringP("project", "p", "", "Assign task to a specific project")
}
