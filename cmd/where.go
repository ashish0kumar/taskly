package cmd

import (
	"fmt"

	"github.com/ashish0kumar/taskly/internal/db"

	"github.com/spf13/cobra"
)

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the location of the tasks database file",
	Long:  `Displays the full path to the SQLite database file where tasks are stored.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the path using the exported function from the db package
		dbPath, err := db.GetStoragePath()
		if err != nil {
			return fmt.Errorf("could not determine database path: %w", err)
		}
		// Print the path to standard output
		_, err = fmt.Println(dbPath)
		return err
	},
}
