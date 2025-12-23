package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	worktreeDir = ".wtree"
	configFile  = "config.toml"
)

// Config represents the configuration file structure
type Config struct {
	Worktree WorktreeConfig `toml:"worktree"`
	Setup    SetupConfig    `toml:"setup"`
	Terminal TerminalConfig `toml:"terminal"`
}

// WorktreeConfig contains worktree-related settings
type WorktreeConfig struct {
	WorktreeBaseDir string `toml:"worktree_base_dir"`
	BranchPrefix    string `toml:"branch_prefix"`
	BaseBranch      string `toml:"base_branch"`
}

// SetupConfig contains setup-related settings
type SetupConfig struct {
	Copy     []string `toml:"copy"`
	Commands []string `toml:"commands"`
}

// TerminalConfig contains terminal-related settings
type TerminalConfig struct {
	Mode string `toml:"mode"`
	Exec string `toml:"exec"`
}

// Load reads the configuration from config.toml
func Load(repoRoot string) (*Config, error) {
	configPath := filepath.Join(repoRoot, worktreeDir, configFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Fill in defaults for empty values
	config.fillDefaults()

	return &config, nil
}

// Save writes the configuration to config.toml
func Save(repoRoot string, config *Config) error {
	configDir := filepath.Join(repoRoot, worktreeDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, configFile)
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Exists checks if the config file exists
func Exists(repoRoot string) bool {
	configPath := filepath.Join(repoRoot, worktreeDir, configFile)
	_, err := os.Stat(configPath)
	return err == nil
}

// fillDefaults fills in default values for empty config fields
func (c *Config) fillDefaults() {
	defaults := DefaultConfig()

	if c.Worktree.WorktreeBaseDir == "" {
		c.Worktree.WorktreeBaseDir = defaults.Worktree.WorktreeBaseDir
	}
	if c.Worktree.BranchPrefix == "" {
		c.Worktree.BranchPrefix = defaults.Worktree.BranchPrefix
	}
	if c.Worktree.BaseBranch == "" {
		c.Worktree.BaseBranch = defaults.Worktree.BaseBranch
	}
	if c.Terminal.Mode == "" {
		c.Terminal.Mode = defaults.Terminal.Mode
	}
}
