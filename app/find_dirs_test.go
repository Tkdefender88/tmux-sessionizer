package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/Tkdefender88/tmux-sessionizer/app"
)

func Test_FindSessionTargets(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		depth int
		want  []string
	}{}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := app.FindSessionTargets(tc.paths, tc.depth)
			assert.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}
