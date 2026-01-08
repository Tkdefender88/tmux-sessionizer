package app

import (
	"os"
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

func (t *Tmux) OpenTmuxSession(target string) error {

	if strings.Contains(target, "[TMUX]") {
		target = strings.TrimPrefix(target, "[TMUX] ")
		_, err := t.run("tmux", "switch-client", "-t", target)
		if err != nil {
			return err
		}
	}

	_, err := t.run("tmux", "new-sesssion", "-d", "-s", target)
	if err != nil {
		return err
	}

	return nil
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
