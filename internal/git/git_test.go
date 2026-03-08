package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetStatus_GetDiff(t *testing.T) {
	// Create a test repo
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git init: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git: %v", err)
	}

	repo, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	// Create initial commit for baseline
	initialFile := filepath.Join(tmpDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.Add("initial.txt"); err != nil {
		t.Fatalf("failed to add initial: %v", err)
	}
	if err := repo.Commit("initial"); err != nil {
		t.Fatalf("failed to commit initial: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	// Stage the file
	if err := repo.Add("test.txt"); err != nil {
		t.Fatalf("failed to add: %v", err)
	}

	// Get diff
	diff, err := repo.Diff("staged")
	if err != nil {
		t.Fatalf("failed to get diff: %v", err)
	}

	if diff == "" {
		t.Error("expected non-empty diff")
	}

	// Verify diff contains expected content
	if !strings.Contains(diff, "hello") {
		t.Errorf("expected diff to contain 'hello', got: %s", diff)
	}
}

func TestGetChangedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git init: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git: %v", err)
	}
	cmd = exec.Command("git", "config", "user.name", "Test")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to configure git: %v", err)
	}

	repo, err := Init(tmpDir)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	// Create initial commit
	initialFile := filepath.Join(tmpDir, "initial.txt")
	if err := os.WriteFile(initialFile, []byte("initial"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := repo.Add("initial.txt"); err != nil {
		t.Fatalf("failed to add initial: %v", err)
	}
	if err := repo.Commit("initial"); err != nil {
		t.Fatalf("failed to commit initial: %v", err)
	}

	// Create test files
	testFile1 := filepath.Join(tmpDir, "test1.txt")
	testFile2 := filepath.Join(tmpDir, "test2.txt")
	if err := os.WriteFile(testFile1, []byte("hello1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile2, []byte("hello2"), 0644); err != nil {
		t.Fatal(err)
	}

	// Stage files
	if err := repo.Add("test1.txt", "test2.txt"); err != nil {
		t.Fatalf("failed to add: %v", err)
	}

	// Get changed files
	files, err := repo.GetChangedFiles("staged")
	if err != nil {
		t.Fatalf("failed to get changed files: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 changed files, got %d", len(files))
	}
}
