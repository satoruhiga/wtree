package git

import (
	"os/exec"
	"strconv"
	"strings"
)

// WorktreeStatus represents the status of a worktree
type WorktreeStatus int

const (
	StatusClean WorktreeStatus = iota
	StatusUncommitted
	StatusAhead
	StatusMerged
	StatusStale // Worktree no longer exists
)

// StatusInfo contains detailed status information
type StatusInfo struct {
	Status      WorktreeStatus
	AheadCount  int
	Description string
}

// GetStatus returns the status of a worktree
func GetStatus(worktreePath, baseBranch, branch string) (*StatusInfo, error) {
	// Check for uncommitted changes
	hasChanges, err := HasUncommittedChanges(worktreePath)
	if err != nil {
		return nil, err
	}
	if hasChanges {
		return &StatusInfo{
			Status:      StatusUncommitted,
			Description: "uncommitted",
		}, nil
	}

	// Check if merged
	merged, err := IsMerged(baseBranch, branch)
	if err != nil {
		return nil, err
	}
	if merged {
		return &StatusInfo{
			Status:      StatusMerged,
			Description: "merged",
		}, nil
	}

	// Check ahead count
	aheadCount, err := GetAheadCount(baseBranch, branch)
	if err != nil {
		return nil, err
	}
	if aheadCount > 0 {
		return &StatusInfo{
			Status:      StatusAhead,
			AheadCount:  aheadCount,
			Description: "ahead " + strconv.Itoa(aheadCount),
		}, nil
	}

	return &StatusInfo{
		Status:      StatusClean,
		Description: "clean",
	}, nil
}

// HasUncommittedChanges checks if there are uncommitted changes in the worktree
func HasUncommittedChanges(worktreePath string) (bool, error) {
	cmd := exec.Command("git", "-C", worktreePath, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(string(output))) > 0, nil
}

// GetAheadCount returns the number of commits ahead of base branch
func GetAheadCount(baseBranch, branch string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", baseBranch+".."+branch)
	output, err := cmd.Output()
	if err != nil {
		// Branch might not exist or other error - return 0
		return 0, nil
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, nil
	}
	return count, nil
}

// IsMerged checks if the branch has been merged into base branch
func IsMerged(baseBranch, branch string) (bool, error) {
	cmd := exec.Command("git", "branch", "--merged", baseBranch)
	output, err := cmd.Output()
	if err != nil {
		return false, nil
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		branchName := strings.TrimSpace(strings.TrimPrefix(line, "*"))
		if branchName == branch {
			return true, nil
		}
	}
	return false, nil
}
