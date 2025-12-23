package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Merge merges the given branch into the current branch
func Merge(branch string) error {
	cmd := exec.Command("git", "merge", branch)
	if output, err := cmd.CombinedOutput(); err != nil {
		outputStr := strings.TrimSpace(string(output))
		if strings.Contains(outputStr, "CONFLICT") || strings.Contains(outputStr, "Automatic merge failed") {
			return fmt.Errorf("merge conflict detected")
		}
		return fmt.Errorf("failed to merge: %s", outputStr)
	}
	return nil
}

// MergeCommitCount returns the number of commits that would be merged
func MergeCommitCount(baseBranch, branch string) (int, error) {
	return GetAheadCount(baseBranch, branch)
}

// HasMergeConflict checks if there's currently a merge conflict
func HasMergeConflict() bool {
	cmd := exec.Command("git", "ls-files", "--unmerged")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

// AbortMerge aborts an ongoing merge
func AbortMerge() error {
	cmd := exec.Command("git", "merge", "--abort")
	return cmd.Run()
}
