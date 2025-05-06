package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/kancli"
	"github.com/spf13/cobra"

	"github.com/ashish0kumar/taskly/internal/task"
)

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "View tasks on an interactive Kanban board",
	Long: `Displays tasks visually categorized by status (todo, in progress, done)
on an interactive Kanban board. Use arrow keys to navigate, Enter/Space
to potentially interact (depending on kancli features).`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbConn == nil {
			return fmt.Errorf("database connection not initialized")
		}

		allTasks, err := dbConn.GetTasks() // Fetch all tasks once
		if err != nil {
			return fmt.Errorf("failed to get tasks for kanban: %w", err)
		}

		if len(allTasks) == 0 {
			fmt.Println("No tasks found to display on the board.")
			return nil
		}

		todosMap := []list.Item{}
		iprMap := []list.Item{}
		finishedMap := []list.Item{}

		for i := range allTasks {
			// Pass task value directly assuming task methods use value receivers
			tsk := allTasks[i]
			switch tsk.Status {
			case task.Todo.String():
				todosMap = append(todosMap, tsk)
			case task.InProgress.String():
				iprMap = append(iprMap, tsk)
			case task.Done.String():
				finishedMap = append(finishedMap, tsk)
			}
		}

		// Create Kanban columns using the status enum values from task package
		todoCol := kancli.NewColumn(todosMap, task.Todo, true)
		iprCol := kancli.NewColumn(iprMap, task.InProgress, false)
		doneCol := kancli.NewColumn(finishedMap, task.Done, false)

		// Note: Persisting changes (moving tasks between columns) requires
		// more advanced integration or a custom Bubble Tea model.
		// Kancli's default board might not automatically save changes back to the DB.

		board := kancli.NewDefaultBoard([]kancli.Column{todoCol, iprCol, doneCol})

		p := tea.NewProgram(board)

		// Run the Bubble Tea program (blocking)
		// Use p.Start() for non-blocking if needed elsewhere.
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("kanban board error: %w", err)
		}

		// Placeholder for potential post-run processing if model state needs saving
		_ = finalModel // Avoid unused variable error

		fmt.Println("\nKanban board closed.")
		return nil
	},
}
