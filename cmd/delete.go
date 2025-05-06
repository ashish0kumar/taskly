package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a task by its ID",
	Long:  `Permanently removes a task from the database using its unique ID.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbConn == nil {
			return fmt.Errorf("database connection not initialized")
		}
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ID %q: %w", idStr, err)
		}

		// Get task details *before* deleting for a better confirmation message.
		taskToDelete, getErr := dbConn.GetTask(uint(id))

		// Attempt to delete the task
		err = dbConn.Delete(uint(id))
		if err != nil {
			return fmt.Errorf("failed to delete task %d: %w", id, err)
		}

		// Show confirmation using details if fetching them succeeded.
		if getErr == nil && taskToDelete.Name != "" {
			fmt.Printf("Task ('%s') deleted.\n", taskToDelete.Name)
		} else {
			// Fallback message if GetTask failed or name was empty
			fmt.Printf("Task deleted.\n")
		}
		return nil
	},
}
