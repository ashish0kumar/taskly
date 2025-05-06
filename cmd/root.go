package cmd

import (
	"fmt"
	"os"

	"github.com/ashish0kumar/taskly/internal/db"

	"github.com/spf13/cobra"
)

// dbConn holds the database connection for use by commands within this package
var dbConn *db.TaskDB

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "taskly",
	Short: "A CLI task management tool.",
	Long: `Taskly helps you manage your tasks efficiently from the command line.
You can add, list, update, delete, and view tasks on a Kanban board.`,
	Args: cobra.NoArgs,
	// PersistentPreRunE runs before any command's RunE. Sets up DB connection
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip database setup for built-in commands.
		// Also skip 'where' command as it doesn't need a live DB connection
		if cmd.Name() == "help" || cmd.Name() == "version" || cmd.Name() == "where" ||
			cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd {
			return nil
		}

		var err error
		// Open DB connection and store it in the package variable.
		dbConn, err = db.OpenDB()
		if err != nil {
			return fmt.Errorf("failed to initialize database: %w", err)
		}
		return nil
	},
	// PersistentPostRunE runs after command's RunE. Closes DB connection.
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		// Only close if dbConn was initialized (i.e., not skipped in PreRun)
		if dbConn != nil {
			if err := dbConn.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Error closing database: %v\n", err)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show help if root command is called without arguments.
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// init registers child commands and flags.
func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(whereCmd)
	rootCmd.AddCommand(kanbanCmd)
}
