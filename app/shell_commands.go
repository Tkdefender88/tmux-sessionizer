package app

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

type CommandResult struct {
	Stdout string
	Stderr string
	Code   int
}

type CommandRunner interface {
	Run(name string, args ...string) (CommandResult, error)
	RunInteractive(name string, args ...string) error
}

type ShellRunner struct{}

func (s ShellRunner) RunInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	slog.Info("running command in interactive mode", "command", name, "args", args)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s ShellRunner) Run(name string, args ...string) (CommandResult, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return CommandResult{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
			Code:   cmd.ProcessState.ExitCode(),
		}, err
	}

	return CommandResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Code:   cmd.ProcessState.ExitCode(),
	}, nil
}

type MockRunner struct {
	runInteractive func(name string, args ...string) error
	run            func(name string, args ...string) (CommandResult, error)
}

func (m *MockRunner) Run(name string, args ...string) (CommandResult, error) {
	if m.run != nil {
		return m.run(name, args...)
	}
	return CommandResult{}, fmt.Errorf("mock not implemented")
}

func (m *MockRunner) RunInteractive(name string, args ...string) error {
	if m.runInteractive != nil {
		return m.runInteractive(name, args...)
	}
	return fmt.Errorf("mock not implemented")
}

func ListSessions() []string {
	return []string{}
}
