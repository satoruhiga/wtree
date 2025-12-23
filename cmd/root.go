package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wtree",
	Short: "A docker-ps style worktree manager",
	Long: `wtree is a worktree manager that provides docker-ps style management
for git worktrees. It helps you create, list, open, and remove worktrees
with 8-character random IDs.

Example workflow:
  1. wtree new              # Create a new worktree
  2. (work on your changes)
  3. wtree merge <id>       # Merge and remove the worktree

Or use GUI tools like Fork to merge, then:
  4. wtree rm <id>          # Remove the worktree`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// exitWithError prints an error message and exits
func exitWithError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
