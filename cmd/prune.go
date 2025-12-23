package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/satoruhiga/wtree/internal/ui"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Clean up merged worktrees and empty directories",
	Long: `Remove all merged worktrees and clean up empty directories in the worktree base directory.

This command:
1. Removes all worktrees that have been merged to the base branch
2. Runs 'git worktree prune' to clean up stale entries
3. Removes empty directories in the worktree base directory

Use this after merging via GUI tools (Fork, etc.) or when worktree directories
couldn't be removed due to locked files.

Examples:
  wtree prune          # Clean up with confirmation
  wtree prune --force  # Skip confirmation`,
	RunE: runPrune,
}

var pruneForce bool

func init() {
	pruneCmd.Flags().BoolVarP(&pruneForce, "force", "f", false, "Skip confirmation")
	rootCmd.AddCommand(pruneCmd)
}

func runPrune(cmd *cobra.Command, args []string) error {
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

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Find merged worktrees
	var mergedSessions []*session.Session
	for _, sess := range store.All() {
		statusInfo, err := git.GetStatus(sess.AbsPath, cfg.Worktree.BaseBranch, sess.Branch)
		if err != nil {
			continue
		}
		if statusInfo.Status == git.StatusMerged {
			mergedSessions = append(mergedSessions, sess)
		}
	}

	// Report what will be done
	if len(mergedSessions) > 0 {
		fmt.Printf("Found %d merged worktree(s):\n", len(mergedSessions))
		for _, sess := range mergedSessions {
			fmt.Printf("  - %s (%s)\n", sess.ID, sess.Branch)
		}
	}

	// Check for empty/orphan directories
	worktreeBaseDir := filepath.Join(repoRoot, cfg.Worktree.WorktreeBaseDir)
	emptyDirs := findEmptyOrOrphanDirs(worktreeBaseDir, store)
	if len(emptyDirs) > 0 {
		fmt.Printf("Found %d empty/orphan directory(ies):\n", len(emptyDirs))
		for _, dir := range emptyDirs {
			fmt.Printf("  - %s\n", filepath.Base(dir))
		}
	}

	if len(mergedSessions) == 0 && len(emptyDirs) == 0 {
		fmt.Println("Nothing to clean up.")
		return nil
	}

	// Confirm
	if !pruneForce {
		if !ui.Confirm("Proceed with cleanup?") {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// Remove merged worktrees
	for _, sess := range mergedSessions {
		// Try to remove worktree
		if err := git.RemoveWorktree(sess.AbsPath, true); err != nil {
			fmt.Printf("%s Failed to remove worktree %s: %v\n", yellow("!"), sess.ID, err)
		}

		// Delete branch
		if err := git.DeleteBranch(sess.Branch, false); err != nil {
			fmt.Printf("%s Failed to delete branch %s: %v\n", yellow("!"), sess.Branch, err)
		}

		// Remove from sessions
		store.Remove(sess.ID)
		fmt.Printf("%s Removed %s\n", green("✓"), sess.ID)
	}

	// Save sessions
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to update sessions: %w", err)
	}

	// Run git worktree prune
	if err := git.PruneWorktrees(); err != nil {
		fmt.Printf("%s git worktree prune failed: %v\n", yellow("!"), err)
	} else {
		fmt.Printf("%s Pruned stale worktree entries\n", green("✓"))
	}

	// Remove empty directories
	for _, dir := range emptyDirs {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("%s Failed to remove directory %s: %v\n", yellow("!"), filepath.Base(dir), err)
		} else {
			fmt.Printf("%s Removed directory %s\n", green("✓"), filepath.Base(dir))
		}
	}

	return nil
}

// findEmptyOrOrphanDirs finds directories in worktreeBaseDir that are empty or not tracked in sessions
func findEmptyOrOrphanDirs(worktreeBaseDir string, store *session.Store) []string {
	var result []string

	entries, err := os.ReadDir(worktreeBaseDir)
	if err != nil {
		return result
	}

	// Build a map of known worktree absolute paths
	knownPaths := make(map[string]bool)
	for _, sess := range store.All() {
		knownPaths[sess.AbsPath] = true
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		dirPath := filepath.Join(worktreeBaseDir, entry.Name())
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			continue
		}

		// Check if it's a known worktree
		if knownPaths[absPath] {
			continue
		}

		// Check if it's a valid git worktree
		if git.WorktreeExists(absPath) {
			continue
		}

		// Check if directory is empty or only contains .git file
		isEmpty := isDirEmptyOrOnlyGit(dirPath)
		if isEmpty {
			result = append(result, dirPath)
		}
	}

	return result
}

// isDirEmptyOrOnlyGit checks if a directory is empty or only contains a .git file/folder
func isDirEmptyOrOnlyGit(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.Name() != ".git" {
			return false
		}
	}
	return true
}
