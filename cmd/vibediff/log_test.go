package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "vibediff-cli-test-*")
	if err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)

	// Initialize git repo
	runGit(t, "init")
	runGit(t, "config", "user.email", "test@example.com")
	runGit(t, "config", "user.name", "Test User")

	// Create initial commit
	os.WriteFile("initial.txt", []byte("initial"), 0644)
	runGit(t, "add", ".")
	runGit(t, "commit", "-m", "initial")

	cleanup := func() {
		os.Chdir(origDir)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func runGit(t *testing.T, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %v: %v", args, err)
	}
}

func TestCommitCommand_Basic(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Make a change
	os.WriteFile("test.txt", []byte("hello world"), 0644)
	runGit(t, "add", ".")

	// Run commit command
	rootCmd.SetArgs([]string{"commit", "Add test file"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("commit failed: %v", err)
	}

	// Verify .vibediff was created at git root
	vibeDir := filepath.Join(".", ".vibediff")
	if _, err := os.Stat(vibeDir); os.IsNotExist(err) {
		t.Error(".vibediff directory was not created")
	}
}

func TestCommitCommand_EmptyPrompt(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	os.WriteFile("test.txt", []byte("hello"), 0644)
	runGit(t, "add", ".")

	rootCmd.SetArgs([]string{"commit", "   "})
	err := rootCmd.Execute()

	if err == nil {
		t.Error("expected error for empty prompt, got nil")
	}
}
