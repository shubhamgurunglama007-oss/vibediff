package main

import (
	"os"
	"testing"
)

func TestShowCommand_Basic(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// First make a vibe commit
	os.WriteFile("test.txt", []byte("hello"), 0644)
	runGit(t, "add", ".")
	rootCmd.SetArgs([]string{"commit", "Add test file"})
	rootCmd.Execute()

	// Reset args for show
	rootCmd.SetArgs([]string{"show"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("show failed: %v", err)
	}
}

func TestShowCommand_InvalidRef(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	rootCmd.SetArgs([]string{"show", "../../../etc/passwd"})
	err := rootCmd.Execute()

	if err == nil {
		t.Error("expected error for invalid ref, got nil")
	}
}
