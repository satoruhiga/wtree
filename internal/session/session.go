package session

import (
	"fmt"
	"time"
)

// Session represents a single worktree session
type Session struct {
	ID        string    `json:"id"`
	Branch    string    `json:"branch"`
	Path      string    `json:"path"`
	AbsPath   string    `json:"abs_path"`
	CreatedAt time.Time `json:"created_at"`
}

// NewSession creates a new Session
func NewSession(id, branch, path, absPath string) *Session {
	return &Session{
		ID:        id,
		Branch:    branch,
		Path:      path,
		AbsPath:   absPath,
		CreatedAt: time.Now(),
	}
}

// RelativeTime returns a human-readable relative time string
func (s *Session) RelativeTime() string {
	d := time.Since(s.CreatedAt)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return formatDuration(mins, "min")
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return formatDuration(hours, "hour")
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return formatDuration(days, "day")
	}
}

func formatDuration(n int, unit string) string {
	return fmt.Sprintf("%d %ss ago", n, unit)
}
