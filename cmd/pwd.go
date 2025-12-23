package cmd

import (
	"fmt"

	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/spf13/cobra"
)

var pwdCmd = &cobra.Command{
	Use:   "pwd <id>",
	Short: "Print worktree directory path",
	Long: `Print the absolute path of a worktree.

Examples:
  wtree pwd a3f8          # Print worktree path
  cd $(wtree pwd a3f8)    # Change to worktree directory
  pushd $(wtree pwd a3f8) # Push worktree directory`,
	Args: cobra.ExactArgs(1),
	RunE: runPwd,
}

func init() {
	rootCmd.AddCommand(pwdCmd)
}

func runPwd(cmd *cobra.Command, args []string) error {
	partialID := args[0]

	// Get repository root
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	// Load sessions
	store := session.NewStore(repoRoot)
	if err := store.Load(); err != nil {
		return err
	}

	// Find session by partial ID
	sess, err := store.FindByPartialID(partialID)
	if err != nil {
		return err
	}

	// Print path
	fmt.Println(sess.AbsPath)
	return nil
}
