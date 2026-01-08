package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListTmuxSessions(t *testing.T) {
	tests := []struct {
		name         string
		want         []string
		runnerResult CommandResult
	}{
		{
			name: "simple case",
			want: []string{
				"[TMUX] foo_directory",
				"[TMUX] bar_directory",
			},
			runnerResult: CommandResult{
				Stdout: "[TMUX] foo_directory\n[TMUX] bar_directory",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmux := NewTmux()
			tmux.Runner = MockRunner(tc.runnerResult)

			got, err := tmux.ListTmuxSessions()
			assert.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
