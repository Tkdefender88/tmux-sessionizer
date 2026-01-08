package app

import (
	"bytes"
	"os/exec"
)

type CommandRunner func(name string, args ...string) (CommandResult, error)

type CommandResult struct {
	Stdout string
	Stderr string
	Code   int
}

func RunCommand(runner CommandRunner, name string, args ...string) (CommandResult, error) {
	return runner(name, args...)
}

func ShellRunner(name string, args ...string) (CommandResult, error) {
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

func MockRunner(result CommandResult) CommandRunner {
	return func(name string, args ...string) (CommandResult, error) {
		return result, nil
	}
}

func ListSessions() []string {
	return []string{}
}
