package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/satoruhiga/wtree/internal/session"
	"github.com/satoruhiga/wtree/internal/ui"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all worktrees",
	Long: `List all managed worktrees with their status in a docker-ps style format.

Status types:
  clean       - No changes
  uncommitted - Has uncommitted changes
  ahead N     - N commits ahead of base branch
  merged      - Already merged to base branch
  stale       - Worktree no longer exists (use 'wtree rm' to clean up)`,
	RunE: runLs,
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func runLs(cmd *cobra.Command, args []string) error {
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

	sessions := store.All()
	if len(sessions) == 0 {
		if !config.Exists(repoRoot) {
			cmd.Println("No .wtree found. Run 'wtree init' to initialize.")
		} else {
			cmd.Println("No worktrees found.")
		}
		return nil
	}

	// Sort by creation time (newest first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	// Build table data
	headers := []string{"ID", "BRANCH", "CREATED", "STATUS", "PATH"}
	var rows [][]string

	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	for _, sess := range sessions {
		var statusStr string

		// First check if worktree exists
		if !git.WorktreeExists(sess.AbsPath) {
			statusStr = gray("stale")
		} else {
			// Get status
			statusInfo, err := git.GetStatus(sess.AbsPath, cfg.Worktree.BaseBranch, sess.Branch)
			if err != nil {
				statusStr = gray("unknown")
			} else {
				switch statusInfo.Status {
				case git.StatusClean:
					statusStr = green(statusInfo.Description)
				case git.StatusUncommitted:
					statusStr = red(statusInfo.Description)
				case git.StatusAhead:
					statusStr = yellow(statusInfo.Description)
				case git.StatusMerged:
					statusStr = cyan(statusInfo.Description)
				default:
					statusStr = gray("unknown")
				}
			}
		}

		rows = append(rows, []string{
			sess.ID,
			sess.Branch,
			sess.RelativeTime(),
			statusStr,
			sess.Path,
		})
	}

	ui.PrintTable(headers, rows)
	return nil
}
