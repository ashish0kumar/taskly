package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"

	"github.com/ashish0kumar/taskly/internal/task"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Long:  `Displays all tasks currently stored in the database, ordered by creation date.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbConn == nil {
			return fmt.Errorf("database connection not initialized")
		}

		tasks, err := dbConn.GetTasks()
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found. Add one with 'taskly add \"My new task\"'")
			return nil
		}

		fmt.Println(setupTable(tasks).String())
		return nil
	},
}

func setupTable(tasks []task.Task) *table.Table {
	columns := []string{"ID", "Name", "Project", "Status", "Created At"}
	var rows [][]string

	for _, t := range tasks {
		// Create raw strings first
		idStr := fmt.Sprintf("%d", t.ID)
		nameStr := t.Name
		projectStr := t.Project
		statusStr := t.Status
		dateStr := t.Created.Format("2006-01-02")

		// Add the row with raw strings
		row := []string{idStr, nameStr, projectStr, statusStr, dateStr}
		rows = append(rows, row)
	}

	// Create the table with raw strings first
	t := table.New().
		Headers(columns...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			baseStyle := lipgloss.NewStyle().Padding(0, 1)
			if row == 0 {
				return baseStyle.Bold(true).Foreground(lipgloss.Color("212"))
			}

			if col == 3 && row > 0 {
				cellValue := rows[row-1][col]
				switch cellValue {
				case "todo":
					return baseStyle.Foreground(lipgloss.Color("240"))
				case "in progress":
					return baseStyle.Foreground(lipgloss.Color("214"))
				case "done":
					return baseStyle.Foreground(lipgloss.AdaptiveColor{Light: "46", Dark: "82"})
				}
			}
			return baseStyle
		})

	return t
}
