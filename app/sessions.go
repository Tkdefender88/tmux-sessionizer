package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type SessionManager struct {
	tmux *Tmux
}

func NewSessionManager() *SessionManager {
	tmux := NewTmux()
	return &SessionManager{
		tmux: tmux,
	}
}

// FindSessionTargets will find and list directories and running tmux sessions from the listed paths while ignoring
// git directories (.git/)
func (s *SessionManager) FindSessionTargets(paths []string, maxDepth int) ([]string, error) {
	targets := []string{}

	sessions, err := s.listActiveSessions()
	if err != nil {
		return []string{}, err
	}
	targets = append(targets, sessions...)

	// find directories
	dirs, err := findDirectories(paths, maxDepth)
	if err != nil {
		return []string{}, err
	}
	targets = append(targets, dirs...)

	return targets, nil
}

func (s *SessionManager) sessionExists(target string) bool {
	if strings.HasPrefix(target, "[TMUX] ") {
		return true
	}

	sessions, err := s.tmux.listSessions()
	if err != nil {
		slog.Error("could not find tmux sessions", "error", err)
		return false
	}

	sessionName := fmt.Sprintf("[TMUX] %s", pathToSessionName(target))

	return slices.Contains(sessions, sessionName)
}

// OpenSession will open a tmux session for the given target
func (s *SessionManager) OpenSession(target string) error {
	sessionName := pathToSessionName(target)

	if !s.sessionExists(target) {
		if err := s.createNewSession(sessionName, target); err != nil {
			slog.Error("failed to create a new tmux session", "error", err)
			return err
		}
	}

	return s.switchTo(sessionName)
}

// createNewSession will create a new tmux session with the given name and target path
// and hydrate it with the given hydrate file
func (s *SessionManager) createNewSession(sessionName, targetPath string) error {
	err := s.tmux.newSession(sessionName, targetPath)
	if err != nil {
		slog.Error("failed to start a new tmux session", "error", err)
		return err
	}

	return s.hydrateSession(sessionName, targetPath)
}

// hydrateSession will hydrate the session with the given hydrate file
// if the hydrate file is in the local target directory it will be used
// otherwise it will look for a global hydration file in the user's home directory
func (s *SessionManager) hydrateSession(sessionName, targetPath string) error {
	localHydrateFile := filepath.Join(targetPath, ".tmux-sessionizer")
	if _, err := os.Stat(localHydrateFile); err == nil {
		slog.Info("local hydration found executing local hydration")
		return s.tmux.hydrate(sessionName, localHydrateFile)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("couldn't hydrate, can't find home directory: %v\n", err)
	}

	globalHydrateFile := filepath.Join(homeDir, ".tmux-sessionizer")
	if _, err := os.Stat(globalHydrateFile); err == nil {
		slog.Info("using global hydration")
		return s.tmux.hydrate(sessionName, globalHydrateFile)
	}
	slog.Info("no hydration found, skipping")
	return nil
}

// switchTo will switch to the target session if inside tmux or attach to it if not
func (s *SessionManager) switchTo(target string) error {
	if isInsideTmuxSession() {
		return s.tmux.switchClient(target)
	}
	return s.tmux.attachSession(target)
}

// listActiveSessions will list all active tmux sessions but exclude the current session if inside tmux
func (s *SessionManager) listActiveSessions() ([]string, error) {
	if isInsideTmuxSession() {
		currentSession, err := s.tmux.currentSessionName()
		if err != nil {
			return []string{}, err
		}

		sessions, err := s.tmux.listSessions()
		if err != nil {
			return []string{}, err
		}

		sessionSelection := make([]string, 0, len(sessions)-1)

		for _, session := range sessions {
			if session == currentSession {
				continue
			}
			sessionSelection = append(sessionSelection, session)
		}

		slog.Info("found the following tmux sessions", "sessions", sessions, "length", len(sessions))
		return sessionSelection, nil
	}

	sessions, err := s.tmux.listSessions()
	if err != nil {
		return []string{}, err
	}

	return sessions, nil
}

// isInsideTmuxSession returns true if the current process is running inside tmux
func isInsideTmuxSession() bool {
	return strings.TrimSpace(os.Getenv("TMUX")) != ""
}

// pathToSessionName will return the session name from a path
func pathToSessionName(path string) string {
	if sessionName, ok := strings.CutPrefix(path, "[TMUX] "); ok {
		return sessionName
	}
	return filepath.Base(path)
}
