package cmd

import (
	"fmt"

	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/satoruhiga/wtree/internal/terminal"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <id>",
	Short: "Open an existing worktree in Windows Terminal",
	Long: `Open an existing worktree in Windows Terminal.
The ID can be a partial match (e.g., 'a3f8' for 'a3f8c2d1').

Examples:
  wtree open a3f8       # Open worktree with ID starting with a3f8
  wtree open a3f8c2d1   # Open worktree with exact ID`,
	Args: cobra.ExactArgs(1),
	RunE: runOpen,
}

var openPane bool

func init() {
	openCmd.Flags().BoolVar(&openPane, "pane", false, "Open in split pane instead of new tab")
	rootCmd.AddCommand(openCmd)
}

func runOpen(cmd *cobra.Command, args []string) error {
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

	// Determine terminal mode
	mode := terminal.ModeTab
	if openPane {
		mode = terminal.ModePane
	} else if cfg.Terminal.Mode == "pane" {
		mode = terminal.ModePane
	} else if cfg.Terminal.Mode == "window" {
		mode = terminal.ModeWindow
	}

	// Open in Windows Terminal
	if terminal.IsAvailable() {
		fmt.Printf("Opening %s in Windows Terminal...\n", sess.ID)
		if err := terminal.OpenInTerminal(sess.AbsPath, mode, cfg.Terminal.Exec); err != nil {
			return fmt.Errorf("failed to open terminal: %w", err)
		}
	} else {
		fmt.Printf("Windows Terminal not found.\nPath: %s\n", sess.AbsPath)
	}

	return nil
}
