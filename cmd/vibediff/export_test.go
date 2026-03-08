package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogCommand_Empty(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Log should not error on empty repo
	rootCmd.SetArgs([]string{"log"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("log failed: %v", err)
	}
}

func TestLogCommand_WithCommits(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Make vibe commits
	for i := 1; i <= 3; i++ {
		os.WriteFile("test.txt", []byte(string(rune('0'+i))), 0644)
		runGit(t, "add", ".")
		rootCmd.SetArgs([]string{"commit", "Commit " + string(rune('0'+i))})
		rootCmd.Execute()
	}

	// Log should succeed with commits
	rootCmd.SetArgs([]string{"log"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("log failed: %v", err)
	}
}

func TestExportCommand_Basic(t *testing.T) {
	_, cleanup := setupTestRepo(t)
	defer cleanup()

	// Make a vibe commit
	os.WriteFile("test.txt", []byte("hello"), 0644)
	runGit(t, "add", ".")
	rootCmd.SetArgs([]string{"commit", "Add test file"})
	rootCmd.Execute()

	// Export
	rootCmd.SetArgs([]string{"export"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Verify HTML file was created
	wd, _ := os.Getwd()
	matches, _ := filepath.Glob(filepath.Join(wd, "vibe-*.html"))
	if len(matches) == 0 {
		t.Error("HTML export file was not created")
	}
}
