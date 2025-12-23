package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/satoruhiga/wtree/internal/ui"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <id>",
	Short: "Remove a worktree",
	Long: `Remove a worktree and its associated branch.
Shows a warning if there are uncommitted or unmerged changes.

Examples:
  wtree rm a3f8         # Remove with confirmation
  wtree rm a3f8 --force # Skip confirmation`,
	Args: cobra.ExactArgs(1),
	RunE: runRm,
}

var rmForce bool

func init() {
	rmCmd.Flags().BoolVarP(&rmForce, "force", "f", false, "Skip confirmation")
	rootCmd.AddCommand(rmCmd)
}

func runRm(cmd *cobra.Command, args []string) error {
	partialID := args[0]

	// Get repository root
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	// Load configuration
	cfg, err := config.Load(repoRoot)
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

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Check if worktree still exists
	worktreeExists := git.WorktreeExists(sess.AbsPath)

	if worktreeExists {
		// Check status and warn if necessary
		if !rmForce {
			statusInfo, err := git.GetStatus(sess.AbsPath, cfg.Worktree.BaseBranch, sess.Branch)
			if err == nil {
				switch statusInfo.Status {
				case git.StatusUncommitted:
					fmt.Printf("%s %s has uncommitted changes.\n", yellow("Warning:"), sess.ID)
					if !ui.Confirm("Continue?") {
						fmt.Println("Cancelled.")
						return nil
					}
				case git.StatusAhead:
					fmt.Printf("%s %s has %d unmerged commits.\n", yellow("Warning:"), sess.ID, statusInfo.AheadCount)
					if !ui.Confirm("Continue?") {
						fmt.Println("Cancelled.")
						return nil
					}
				}
			}
		}

		// Remove worktree
		if err := git.RemoveWorktree(sess.AbsPath, rmForce); err != nil {
			fmt.Printf("%s Failed to remove worktree: %v\n", yellow("Warning:"), err)
		}

		// Delete branch
		if err := git.DeleteBranch(sess.Branch, rmForce); err != nil {
			fmt.Printf("%s Failed to delete branch %s: %v\n", yellow("Warning:"), sess.Branch, err)
		}
	} else {
		// Worktree doesn't exist, just clean up session and try to delete branch
		fmt.Printf("Worktree %s no longer exists, cleaning up session...\n", sess.ID)

		// Try to delete branch anyway (it might still exist)
		if git.BranchExists(sess.Branch) {
			if err := git.DeleteBranch(sess.Branch, true); err != nil {
				fmt.Printf("%s Failed to delete branch %s: %v\n", yellow("Warning:"), sess.Branch, err)
			}
		}
	}

	// Remove from sessions
	store.Remove(sess.ID)
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to update sessions: %w", err)
	}

	fmt.Printf("%s Removed %s\n", green("✓"), sess.ID)

	return nil
}
