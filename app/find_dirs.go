package app

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	extraPathPattern = "([^:]+):([0-9]+)"
)

var (
	extraPathRegex *regexp.Regexp
)

func init() {
	extraPathRegex = regexp.MustCompile(extraPathPattern)
}

// FindSessionTargets will find and list directories and running tmux sessions from the listed paths while ignoring
// git directories (.git/)
func FindSessionTargets(paths []string, maxDepth int) ([]string, error) {
	targets := []string{}

	tmux := NewTmux()
	sessions, err := tmux.ListTmuxSessions()
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

func parsePath(path string, depth int) (string, int) {
	matches := extraPathRegex.FindStringSubmatch(path)
	if len(matches) == 3 {
		// matches[0] is the full match, matches[1] is the path, matches[2] is the depth
		pathPart := matches[1]
		depthLimit, err := strconv.Atoi(matches[2])
		if err != nil {
			return path, depth
		}
		return pathPart, depthLimit
	}

	return path, depth
}

func findDirectories(paths []string, maxDepth int) ([]string, error) {
	if len(paths) == 0 {
		return []string{}, nil
	}

	directories := []string{}

	// find directories
	for _, p := range paths {
		path, depthLimit := parsePath(p, maxDepth)
		searchPath, err := expandPath(path)
		if err != nil {
			return []string{}, err
		}
		startDepth := len(strings.Split(searchPath, string(os.PathSeparator)))
		err = filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			currentDepth := len(strings.Split(path, string(os.PathSeparator))) - startDepth
			if path == searchPath {
				return nil
			}

			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}

				if currentDepth > depthLimit {
					return filepath.SkipDir
				}

				directories = append(directories, path)
			}

			return nil
		})
		if err != nil {
			return []string{}, err
		}
	}

	return directories, nil
}

func expandPath(path string) (string, error) {
	// Expand environment variables like $HOME
	path = os.ExpandEnv(path)

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[2:]), nil
	}

	return path, nil
}
