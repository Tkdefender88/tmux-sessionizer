package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app"
)

func Test_FindDirs(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		depth int
		want  []string
	}{
		{
			name:  "Find dirs in current directory",
			paths: []string{"~/workspace/personal/tmux-sessionizer"},
			depth: 1,
			want: []string{
				"/Users/justin.bak/workspace/personal/tmux-sessionizer/app",
			},
		},
		{
			name:  "Find dirs in current directory depth 2",
			paths: []string{"~/workspace/personal/tmux-sessionizer"},
			depth: 2,
			want: []string{
				"/Users/justin.bak/workspace/personal/tmux-sessionizer/app",
				"/Users/justin.bak/workspace/personal/tmux-sessionizer/app/filepicker",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := app.FindDirs(tc.paths, tc.depth)
			assert.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
