package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ashish0kumar/taskly/internal/task"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a task's details (name, project, status)",
	Long: `Updates the specified task's name, project, or status.
Provide the task ID and use flags for the fields you want to change.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbConn == nil {
			return fmt.Errorf("database connection not initialized")
		}

		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ID %q: %w", idStr, err)
		}

		// Use pointers to detect which flags were actually set
		var name *string
		if cmd.Flags().Changed("name") {
			n, _ := cmd.Flags().GetString("name")
			name = &n
		}

		var project *string
		if cmd.Flags().Changed("project") {
			p, _ := cmd.Flags().GetString("project")
			project = &p
		}

		var statusStrPtr *string
		if cmd.Flags().Changed("status") {
			sInt, _ := cmd.Flags().GetInt("status")
			var sStr string
			isValid := false
			// Validate against the defined status enum values from task package
			allStatuses := task.AllStatuses() // Get all defined statuses
			for _, s := range allStatuses {
				if sInt == s.Int() {
					sStr = s.String()
					isValid = true
					break
				}
			}

			if !isValid {
				validOptions := []string{}
				for _, s := range allStatuses {
					validOptions = append(validOptions, fmt.Sprintf("%d=%s", s.Int(), s.String()))
				}
				return fmt.Errorf("invalid status value: %d. Use %s", sInt, strings.Join(validOptions, ", "))
			}
			statusStrPtr = &sStr
		}

		// Call the exported Update method
		updatedTask, err := dbConn.Update(uint(id), name, project, statusStrPtr)
		if err != nil {
			return fmt.Errorf("failed to update task %d: %w", id, err)
		}
		fmt.Printf("Task ('%s') updated.\n", updatedTask.Name)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringP("name", "n", "", "Update the name of the task")
	updateCmd.Flags().StringP("project", "p", "", "Update the project of the task")
	updateCmd.Flags().IntP("status", "s", -1, "Update status: 0=todo, 1=in progress, 2=done")
}
