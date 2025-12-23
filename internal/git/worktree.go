package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetRepoRoot returns the root directory of the main git repository.
// If called from within a worktree, it returns the main repository root,
// not the worktree root.
func GetRepoRoot() (string, error) {
	// First, get the common git directory (shared across worktrees)
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}

	gitCommonDir := strings.TrimSpace(string(output))

	// If it's just ".git", we're in the main repo
	if gitCommonDir == ".git" {
		cmd = exec.Command("git", "rev-parse", "--show-toplevel")
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("not a git repository: %w", err)
		}
		return strings.TrimSpace(string(output)), nil
	}

	// Otherwise, gitCommonDir points to the main repo's .git directory
	// e.g., "/path/to/main-repo/.git" or "../main-repo/.git"
	absGitDir, err := filepath.Abs(gitCommonDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve git directory: %w", err)
	}

	// The repo root is the parent of .git
	repoRoot := filepath.Dir(absGitDir)
	return repoRoot, nil
}

// GetCurrentWorktreeRoot returns the root of the current worktree (or main repo)
func GetCurrentWorktreeRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// IsInWorktree returns true if the current directory is inside a worktree (not the main repo)
func IsInWorktree() bool {
	cmd := exec.Command("git", "rev-parse", "--git-common-dir")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	gitCommonDir := strings.TrimSpace(string(output))
	return gitCommonDir != ".git"
}

// AddWorktree creates a new worktree with a new branch
func AddWorktree(path, branch, baseBranch string) error {
	cmd := exec.Command("git", "worktree", "add", "-b", branch, path, baseBranch)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create worktree: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// RemoveWorktree removes a worktree
func RemoveWorktree(path string, force bool) error {
	args := []string{"worktree", "remove", path}
	if force {
		args = append(args, "--force")
	}
	cmd := exec.Command("git", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

// WorktreeExists checks if a worktree exists at the given path
func WorktreeExists(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	// Normalize path separators for comparison (git uses forward slashes on Windows)
	normalizedPath := filepath.ToSlash(absPath)

	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "worktree ") {
			wtPath := strings.TrimPrefix(line, "worktree ")
			// Normalize git's path as well
			normalizedWtPath := filepath.ToSlash(wtPath)
			if strings.EqualFold(normalizedPath, normalizedWtPath) {
				return true
			}
		}
	}
	return false
}

// PruneWorktrees runs git worktree prune to clean up stale worktree entries
func PruneWorktrees() error {
	cmd := exec.Command("git", "worktree", "prune")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to prune worktrees: %s", strings.TrimSpace(string(output)))
	}
	return nil
}
