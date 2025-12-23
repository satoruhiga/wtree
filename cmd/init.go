package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/satoruhiga/wtree/internal/config"
	"github.com/satoruhiga/wtree/internal/git"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wtree configuration",
	Long: `Create a .wtree/config.toml configuration file with default settings.
The base branch is automatically detected from the repository.

Example:
  wtree init`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	// Get repository root
	repoRoot, err := git.GetRepoRoot()
	if err != nil {
		return err
	}

	// Check if config already exists
	if config.Exists(repoRoot) {
		return fmt.Errorf("config file already exists at .wtree/config.toml")
	}

	// Detect default branch
	baseBranch, err := git.GetDefaultBranch()
	if err != nil {
		baseBranch = "main"
	}

	// Create .wtree directory
	wtreeDir := filepath.Join(repoRoot, ".wtree")
	if err := os.MkdirAll(wtreeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .wtree directory: %w", err)
	}

	// Write config template
	configPath := filepath.Join(wtreeDir, "config.toml")
	template := config.ConfigTemplate(baseBranch)
	if err := os.WriteFile(configPath, []byte(template), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s Created .wtree/config.toml\n", green("✓"))
	fmt.Printf("Base branch: %s\n", baseBranch)
	fmt.Println("\nTip: Add .wtree/ to your .gitignore")

	return nil
}
