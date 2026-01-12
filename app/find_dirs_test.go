package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parsePath(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		defaultDepth  int
		expectedPath  string
		expectedDepth int
	}{
		{
			name:          "path with custom depth",
			path:          "~/projects:3",
			defaultDepth:  5,
			expectedPath:  "~/projects",
			expectedDepth: 3,
		},
		{
			name:          "path without custom depth",
			path:          "~/workspace",
			defaultDepth:  5,
			expectedPath:  "~/workspace",
			expectedDepth: 5,
		},
		{
			name:          "absolute path with custom depth",
			path:          "/home/user/code:2",
			defaultDepth:  10,
			expectedPath:  "/home/user/code",
			expectedDepth: 2,
		},
		{
			name:          "path with depth of 0",
			path:          "/tmp:0",
			defaultDepth:  5,
			expectedPath:  "/tmp",
			expectedDepth: 0,
		},
		{
			name:          "path with invalid depth number returns defaults",
			path:          "~/projects:abc",
			defaultDepth:  5,
			expectedPath:  "~/projects:abc",
			expectedDepth: 5,
		},
		{
			name:          "empty path",
			path:          "",
			defaultDepth:  5,
			expectedPath:  "",
			expectedDepth: 5,
		},
		{
			name:          "path with only colon and number",
			path:          ":5",
			defaultDepth:  3,
			expectedPath:  ":5",
			expectedDepth: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotPath, gotDepth := parsePath(tc.path, tc.defaultDepth)
			assert.Equal(t, tc.expectedPath, gotPath, "path mismatch")
			assert.Equal(t, tc.expectedDepth, gotDepth, "depth mismatch")
		})
	}
}
