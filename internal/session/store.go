package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	worktreeDir  = ".wtree"
	sessionsFile = "sessions.json"
)

// Store manages session persistence
type Store struct {
	repoRoot string
	sessions map[string]*Session
}

// NewStore creates a new Store for the given repository root
func NewStore(repoRoot string) *Store {
	return &Store{
		repoRoot: repoRoot,
		sessions: make(map[string]*Session),
	}
}

// sessionsPath returns the full path to sessions.json
func (s *Store) sessionsPath() string {
	return filepath.Join(s.repoRoot, worktreeDir, sessionsFile)
}

// worktreeDirPath returns the full path to .wtree directory
func (s *Store) worktreeDirPath() string {
	return filepath.Join(s.repoRoot, worktreeDir)
}

// Load reads sessions from sessions.json
func (s *Store) Load() error {
	data, err := os.ReadFile(s.sessionsPath())
	if err != nil {
		if os.IsNotExist(err) {
			s.sessions = make(map[string]*Session)
			return nil
		}
		return fmt.Errorf("failed to read sessions file: %w", err)
	}

	if err := json.Unmarshal(data, &s.sessions); err != nil {
		return fmt.Errorf("failed to parse sessions file: %w", err)
	}

	return nil
}

// Save writes sessions to sessions.json
func (s *Store) Save() error {
	// Ensure .wtree directory exists
	if err := os.MkdirAll(s.worktreeDirPath(), 0755); err != nil {
		return fmt.Errorf("failed to create .wtree directory: %w", err)
	}

	data, err := json.MarshalIndent(s.sessions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal sessions: %w", err)
	}

	if err := os.WriteFile(s.sessionsPath(), data, 0644); err != nil {
		return fmt.Errorf("failed to write sessions file: %w", err)
	}

	return nil
}

// Add adds a new session
func (s *Store) Add(session *Session) {
	s.sessions[session.ID] = session
}

// Remove removes a session by ID
func (s *Store) Remove(id string) {
	delete(s.sessions, id)
}

// Get returns a session by exact ID
func (s *Store) Get(id string) (*Session, bool) {
	session, ok := s.sessions[id]
	return session, ok
}

// FindByPartialID finds a session by partial ID match
func (s *Store) FindByPartialID(partialID string) (*Session, error) {
	var matches []*Session

	for id, session := range s.sessions {
		if strings.HasPrefix(id, partialID) {
			matches = append(matches, session)
		}
	}

	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("worktree not found: %s", partialID)
	case 1:
		return matches[0], nil
	default:
		return nil, fmt.Errorf("ambiguous ID '%s': matches %d worktrees", partialID, len(matches))
	}
}

// All returns all sessions
func (s *Store) All() []*Session {
	result := make([]*Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		result = append(result, session)
	}
	return result
}

// Count returns the number of sessions
func (s *Store) Count() int {
	return len(s.sessions)
}
