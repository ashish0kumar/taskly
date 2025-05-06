package main

import (
	"fmt"
	"os"

	"github.com/ashish0kumar/taskly/cmd"
)

// main is the simple entry point that executes the root command.
func main() {
	// Execute the root command defined in the cmd package.
	if err := cmd.Execute(); err != nil {
		// Print errors returned by commands or setup to stderr.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
