package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge <id>",
	Short: "Merge a worktree branch and remove the worktree",
	Long: `Merge the worktree's branch into the current branch and remove the worktree.
If there are conflicts, the worktree is kept for manual resolution.

Examples:
  wtree merge a3f8      # Merge and remove worktree`,
	Args: cobra.ExactArgs(1),
	RunE: runMerge,
}

func init() {
	rootCmd.AddCommand(mergeCmd)
}

func runMerge(cmd *cobra.Command, args []string) error {
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

	// Check if we're in the worktree itself
	currentBranch, _ := git.GetCurrentBranch()
	if currentBranch == sess.Branch {
		return fmt.Errorf("cannot merge from within the worktree itself. Please run from the main repository")
	}

	// Check for uncommitted changes in worktree
	hasChanges, err := git.HasUncommittedChanges(sess.AbsPath)
	if err != nil {
		return err
	}
	if hasChanges {
		return fmt.Errorf("worktree %s has uncommitted changes. Please commit or stash them first", sess.ID)
	}

	// Get ahead count for display
	aheadCount, _ := git.GetAheadCount(cfg.Worktree.BaseBranch, sess.Branch)

	// Merge
	fmt.Printf("Merging %s into %s...\n", sess.Branch, currentBranch)
	if err := git.Merge(sess.Branch); err != nil {
		if err.Error() == "merge conflict detected" {
			yellow := color.New(color.FgYellow).SprintFunc()
			fmt.Printf("%s Conflict detected. Resolve manually.\n", yellow("!"))
			fmt.Printf("Worktree kept at: %s\n", sess.Path)
			return nil
		}
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	if aheadCount > 0 {
		fmt.Printf("%s Merged %d commits\n", green("✓"), aheadCount)
	} else {
		fmt.Printf("%s Merged\n", green("✓"))
	}

	// Remove worktree
	if err := git.RemoveWorktree(sess.AbsPath, false); err != nil {
		fmt.Printf("Warning: failed to remove worktree: %v\n", err)
	}

	// Delete branch
	if err := git.DeleteBranch(sess.Branch, false); err != nil {
		fmt.Printf("Warning: failed to delete branch: %v\n", err)
	}

	// Remove from sessions
	store.Remove(sess.ID)
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to update sessions: %w", err)
	}

	fmt.Printf("%s Removed worktree %s\n", green("✓"), sess.ID)

	return nil
}
