package config

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Worktree: WorktreeConfig{
			WorktreeBaseDir: "../worktree",
			BranchPrefix:    "wt/",
			BaseBranch:      "main",
		},
		Setup: SetupConfig{
			Copy:     []string{},
			Commands: []string{},
		},
		Terminal: TerminalConfig{
			Mode: "tab",
			Exec: "",
		},
	}
}

// ConfigTemplate returns a template configuration with comments
func ConfigTemplate(baseBranch string) string {
	if baseBranch == "" {
		baseBranch = "main"
	}
	return `[worktree]
# Directory where worktrees are created (relative to repo root)
worktree_base_dir = "../worktree"
# Branch name prefix
branch_prefix = "wt/"
# Base branch for new worktrees
base_branch = "` + baseBranch + `"

[setup]
# Files/directories to copy to new worktrees (supports gitignored files)
copy = [
    # ".env",
    # ".claude/",
]
# Commands to run after worktree creation
commands = [
    # "npm install",
]

[terminal]
# How to open Windows Terminal: "tab" | "pane" | "window"
mode = "pane"
# Command to run after opening (optional)
# exec = "claude"
`
}
