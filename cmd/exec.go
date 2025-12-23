package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute terminal.exec command",
	Long: `Execute the command specified in terminal.exec config.

This is useful for starting your development environment (e.g., claude, vim)
in the current worktree.

Examples:
  wtree exec  # Run terminal.exec command`,
	Args: cobra.NoArgs,
	RunE: runExec,
}

func init() {
	rootCmd.AddCommand(execCmd)
}

func runExec(cmd *cobra.Command, args []string) error {
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

	if cfg.Terminal.Exec == "" {
		return fmt.Errorf("terminal.exec not configured in .wtree/config.toml")
	}

	// Execute command
	var execCommand *exec.Cmd
	if runtime.GOOS == "windows" {
		execCommand = exec.Command("cmd", "/c", cfg.Terminal.Exec)
	} else {
		execCommand = exec.Command("sh", "-c", cfg.Terminal.Exec)
	}

	execCommand.Stdin = os.Stdin
	execCommand.Stdout = os.Stdout
	execCommand.Stderr = os.Stderr

	return execCommand.Run()
}
