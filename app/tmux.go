package app

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Tmux struct {
	Runner CommandRunner
}

func NewTmux() *Tmux {
	return &Tmux{
		Runner: ShellRunner,
	}
}

func (t *Tmux) run(name string, args ...string) (CommandResult, error) {
	return t.Runner(name, args...)
}

func isInsideTmuxSession() bool {
	return os.Getenv("TMUX") != ""
}

func (t *Tmux) sessionExists(target string) bool {
	if strings.HasPrefix(target, "[TMUX] ") {
		return true
	}

	sessions, err := t.ListTmuxSessions()
	if err != nil {
		return false
	}

	sessionName := fmt.Sprintf("[TMUX] %s", pathToSessionName(target))

	return slices.Contains(sessions, sessionName)
}

func pathToSessionName(path string) string {
	if sessionName, ok := strings.CutPrefix(path, "[TMUX] "); ok {
		return sessionName
	}
	return filepath.Base(path)
}

func (t *Tmux) OpenTmuxSession(target string) error {
	sessionName := pathToSessionName(target)

	if !t.sessionExists(target) {
		if err := t.createNewSession(sessionName, target); err != nil {
			return err
		}
	}

	return t.switchTo(sessionName)
}

func (t *Tmux) createNewSession(sessionName, targetPath string) error {
	_, err := t.run("tmux", "new-session", "-ds", sessionName, "-c", targetPath)
	if err != nil {
		return err
	}

	return t.hydrateSession(sessionName, targetPath)
}

func (t *Tmux) hydrateSession(sessionName, targetPath string) error {
	localHydrateFile := filepath.Join(targetPath, ".tmux-sessionizer")
	if _, err := os.Stat(localHydrateFile); err == nil {
		return t.hydrate(sessionName, localHydrateFile)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("couldn't hydrate, can't find home directory: %v\n", err)
	}

	globalHydrateFile := filepath.Join(homeDir, ".tmux-sessionizer")
	if _, err := os.Stat(globalHydrateFile); err == nil {
		return t.hydrate(sessionName, globalHydrateFile)
	}
	return nil
}

func (t *Tmux) hydrate(sessionName, hydrateFile string) error {
	_, err := t.run("tmux", "send-keys", "-t", sessionName, fmt.Sprintf("source %s", hydrateFile), "c-M")
	return err
}

func (t *Tmux) switchTo(target string) error {
	if isInsideTmuxSession() {
		_, err := t.run("tmux", "switch-client", "-t", target)
		return err
	}
	_, err := t.run("tmux", "attach-session", "-t", target)
	return err
}

func (t *Tmux) ListTmuxSessions() ([]string, error) {
	if isInsideTmuxSession() {
		curentSessionResult, err := t.run("tmux", "display-message", "-p", "[TMUX] #{session_name}")
		if err != nil {
			return []string{}, err
		}
		currentSession := strings.TrimSpace(curentSessionResult.Stdout)

		result, err := t.run("tmux", "list-sessions", "-F", "[TMUX] #{session_name}")
		if err != nil {
			return []string{}, err
		}

		output := strings.Split(strings.TrimSpace(result.Stdout), "\n")
		sessions := make([]string, 0)

		for _, line := range output {
			if line == currentSession {
				continue
			}
			sessions = append(sessions, line)
		}

		return sessions, nil
	}

	result, err := t.run("tmux", "list-sessions", "-F", "[TMUX] #{session_name}")
	if err != nil {
		return []string{}, err
	}

	sessions := strings.Split(result.Stdout, "\n")
	return sessions, nil
}
