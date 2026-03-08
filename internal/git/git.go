// Package git provides a wrapper around git commands for VibeDiff.
// It handles common git operations like staging, committing, and retrieving diffs.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Repo represents a Git repository at a specific path.
type Repo struct {
	path string
}

// Init creates a new Repo for the given path after verifying
// it is a valid git repository.
func Init(path string) (*Repo, error) {
	// Verify it's a git repo
	_, err := runGit(path, "rev-parse", "--git-dir")
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}
	return &Repo{path: path}, nil
}

// Add stages files for commit
func (r *Repo) Add(files ...string) error {
	_, err := runGit(r.path, append([]string{"add"}, files...)...)
	return err
}

// Commit creates a new commit
func (r *Repo) Commit(message string) error {
	_, err := runGit(r.path, "commit", "-m", message)
	return err
}

// Diff returns the diff for staged or unstaged changes
func (r *Repo) Diff(which string) (string, error) {
	var args []string
	if which == "staged" {
		args = []string{"diff", "--staged"}
	} else {
		args = []string{"diff"}
	}
	return runGit(r.path, args...)
}

// GetCurrentHead returns the current commit hash
func (r *Repo) GetCurrentHead() (string, error) {
	return runGit(r.path, "rev-parse", "HEAD")
}

// GetChangedFiles returns list of changed files
func (r *Repo) GetChangedFiles(which string) ([]string, error) {
	var args []string
	if which == "staged" {
		args = []string{"diff", "--staged", "--name-only"}
	} else {
		args = []string{"diff", "--name-only"}
	}
	output, err := runGit(r.path, args...)
	if err != nil {
		return nil, err
	}
	if output == "" {
		return []string{}, nil
	}
	return strings.Split(strings.TrimSpace(output), "\n"), nil
}

// RunCommand executes a git command and returns output
func (r *Repo) RunCommand(args ...string) (string, error) {
	return runGit(r.path, args...)
}

// RootPath returns the absolute path to the git repository root
func (r *Repo) RootPath() (string, error) {
	root, err := runGit(r.path, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", fmt.Errorf("failed to get git root: %w", err)
	}
	return strings.TrimSpace(root), nil
}

func runGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git %v: %s: %w", args, errOut.String(), err)
	}
	return out.String(), nil
}
