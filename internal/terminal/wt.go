package terminal

import (
	"fmt"
	"os/exec"
	"runtime"
)

// OpenMode represents how to open the terminal
type OpenMode string

const (
	ModeTab    OpenMode = "tab"
	ModePane   OpenMode = "pane"
	ModeWindow OpenMode = "window"
)

// OpenInTerminal opens the given path in a terminal
// Windows: Windows Terminal (wt.exe)
// macOS/Linux: tmux
func OpenInTerminal(path string, mode OpenMode, execCmd string) error {
	if runtime.GOOS == "windows" {
		return openWindowsTerminal(path, mode, execCmd)
	}
	return openTmux(path, mode, execCmd)
}

// openWindowsTerminal opens Windows Terminal
func openWindowsTerminal(path string, mode OpenMode, execCmd string) error {
	var args []string

	switch mode {
	case ModePane:
		args = []string{"-w", "0", "sp", "-d", path}
	case ModeWindow:
		args = []string{"new-window", "-d", path}
	default:
		args = []string{"-w", "0", "nt", "-d", path}
	}

	if execCmd != "" {
		args = append(args, "cmd", "/k", execCmd)
	}

	cmd := exec.Command("wt.exe", args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open Windows Terminal: %w", err)
	}
	return nil
}

// openTmux opens tmux pane/window
func openTmux(path string, mode OpenMode, execCmd string) error {
	var args []string

	switch mode {
	case ModePane:
		args = []string{"split-window", "-c", path}
	default:
		args = []string{"new-window", "-c", path}
	}

	if execCmd != "" {
		args = append(args, execCmd)
	}

	cmd := exec.Command("tmux", args...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open tmux: %w", err)
	}
	return nil
}

// IsAvailable checks if terminal is available
func IsAvailable() bool {
	if runtime.GOOS == "windows" {
		_, err := exec.LookPath("wt.exe")
		return err == nil
	}
	_, err := exec.LookPath("tmux")
	return err == nil
}

// TerminalName returns the name of the terminal being used
func TerminalName() string {
	if runtime.GOOS == "windows" {
		return "Windows Terminal"
	}
	return "tmux"
}
