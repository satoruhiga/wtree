package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/id"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/satoruhiga/wtree/internal/terminal"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new worktree and open in terminal",
	Long: `Create a new worktree with an 8-character random ID and open it
in terminal (Windows Terminal on Windows, tmux on macOS/Linux).

Examples:
  wtree new          # Create and open in new tab/window
  wtree new --pane   # Create and open in split pane
  wtree new -q       # Create without opening terminal
  wtree new -n 3     # Create 3 worktrees at once`,
	RunE: runNew,
}

var (
	newPane  bool
	newQuiet bool
	newCount int
)

func init() {
	newCmd.Flags().BoolVar(&newPane, "pane", false, "Open in split pane instead of new tab")
	newCmd.Flags().BoolVarP(&newQuiet, "quiet", "q", false, "Create worktree without opening terminal")
	newCmd.Flags().IntVarP(&newCount, "n", "n", 1, "Number of worktrees to create")
	rootCmd.AddCommand(newCmd)
}

func runNew(cmd *cobra.Command, args []string) error {
	// Get repository root
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return fmt.Errorf("not a git repository")
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

	// Determine terminal mode
	mode := terminal.ModeTab
	if newPane {
		mode = terminal.ModePane
	} else if cfg.Terminal.Mode == "pane" {
		mode = terminal.ModePane
	} else if cfg.Terminal.Mode == "window" {
		mode = terminal.ModeWindow
	}

	// Create worktrees
	for i := 0; i < newCount; i++ {
		if err := createWorktree(repoRoot, cfg, store, mode, newQuiet); err != nil {
			return err
		}
	}

	return nil
}

func createWorktree(repoRoot string, cfg *config.Config, store *session.Store, mode terminal.OpenMode, quiet bool) error {
	// Generate ID
	newID, err := id.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate ID: %w", err)
	}

	// Build paths and names
	branchName := cfg.Worktree.BranchPrefix + newID
	worktreeRelPath := filepath.Join(cfg.Worktree.WorktreeBaseDir, "wt-"+newID)
	worktreeAbsPath := filepath.Join(repoRoot, worktreeRelPath)

	// Create worktree
	if err := git.AddWorktree(worktreeAbsPath, branchName, cfg.Worktree.BaseBranch); err != nil {
		return err
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("Created: %s\n", green(newID))

	// Copy files if configured
	if len(cfg.Setup.Copy) > 0 {
		for _, item := range cfg.Setup.Copy {
			srcPath := filepath.Join(repoRoot, item)
			dstPath := filepath.Join(worktreeAbsPath, item)
			if err := copyPath(srcPath, dstPath); err != nil {
				fmt.Printf("Warning: failed to copy %s: %v\n", item, err)
			}
		}
	}

	// Run setup commands if configured
	if len(cfg.Setup.Commands) > 0 {
		for _, cmdStr := range cfg.Setup.Commands {
			fmt.Printf("Running: %s\n", cmdStr)
			var execCmd *exec.Cmd
			if runtime.GOOS == "windows" {
				execCmd = exec.Command("cmd", "/c", cmdStr)
			} else {
				execCmd = exec.Command("sh", "-c", cmdStr)
			}
			execCmd.Dir = worktreeAbsPath
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			if err := execCmd.Run(); err != nil {
				fmt.Printf("Warning: command failed: %v\n", err)
			}
		}
	}

	// Save session
	sess := session.NewSession(newID, branchName, worktreeRelPath, worktreeAbsPath)
	store.Add(sess)
	if err := store.Save(); err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Open in terminal (unless quiet mode)
	if !quiet {
		if terminal.IsAvailable() {
			fmt.Printf("Opening in %s...\n", terminal.TerminalName())
			if err := terminal.OpenInTerminal(worktreeAbsPath, mode, cfg.Terminal.Exec); err != nil {
				fmt.Printf("Warning: failed to open terminal: %v\n", err)
			}
		} else {
			fmt.Printf("Path: %s\n", worktreeAbsPath)
		}
	}

	return nil
}

// copyPath copies a file or directory
func copyPath(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}
