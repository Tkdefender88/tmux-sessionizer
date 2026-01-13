package app

import (
	"fmt"
	"strings"
)

type Tmux struct {
	runner CommandRunner
}

type OptionFunc func(*Tmux)

func WithRunner(runner CommandRunner) OptionFunc {
	return func(t *Tmux) {
		t.runner = runner
	}
}

func NewTmux(options ...OptionFunc) *Tmux {
	tmux := &Tmux{
		runner: ShellRunner{},
	}

	for _, option := range options {
		option(tmux)
	}

	return tmux
}

func (t *Tmux) run(name string, args ...string) (CommandResult, error) {
	return t.runner.Run(name, args...)
}

func (t *Tmux) currentSessionName() (string, error) {
	result, err := t.runner.Run("tmux", "display-message", "-p", "[TMUX] #{session_name}")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result.Stdout), nil
}

func (t *Tmux) listSessions() ([]string, error) {
	result, err := t.runner.Run("tmux", "list-sessions", "-F", "[TMUX] #{session_name}")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(strings.TrimSpace(result.Stdout), "\n"), nil
}

func (t *Tmux) hydrate(sessionName, hydrateFile string) error {
	_, err := t.runner.Run("tmux", "send-keys", "-t", sessionName, fmt.Sprintf("source %s", hydrateFile), "c-M")
	return err
}

func (t *Tmux) switchClient(target string) error {
	_, err := t.runner.Run("tmux", "switch-client", "-t", target)
	return err
}

func (t *Tmux) newSession(sessionName, targetPath string) error {
	_, err := t.runner.Run("tmux", "new-session", "-ds", sessionName, "-c", targetPath)
	return err
}

func (t *Tmux) attachSession(target string) error {
	return t.runner.RunInteractive("tmux", "attach-session", "-t", target)
}
