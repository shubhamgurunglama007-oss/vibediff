package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestJSONStore_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := filepath.Join(tmpDir, "commits.json")
	js := NewJSONStore(dataFile)

	meta := CommitMetadata{
		Commit:    "abc123",
		Prompt:    "Test prompt",
		Timestamp: time.Now().Format(time.RFC3339),
		Files:     []string{"test.go"},
	}

	// Save
	err := js.Save(meta)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	// Load
	loaded, err := js.Load("abc123")
	if err != nil {
		t.Fatalf("failed to load: %v", err)
	}

	if loaded.Prompt != "Test prompt" {
		t.Errorf("expected prompt 'Test prompt', got %s", loaded.Prompt)
	}
}

func TestJSONStore_FilePermissions(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := filepath.Join(tmpDir, "commits.json")
	js := NewJSONStore(dataFile)

	meta := CommitMetadata{
		Commit:    "abc123",
		Prompt:    "test prompt",
		Timestamp: "2024-01-01T00:00:00Z",
	}

	err := js.Save(meta)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Check file permissions
	info, err := os.Stat(dataFile)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	mode := info.Mode()
	if mode.Perm() != 0600 {
		t.Errorf("File has wrong permissions: got %o, want 0600", mode.Perm())
	}
}

func TestJSONStore_CorruptJSON(t *testing.T) {
	tmpDir := t.TempDir()
	dataFile := filepath.Join(tmpDir, "commits.json")

	// Write invalid JSON
	err := os.WriteFile(dataFile, []byte("{invalid json}"), 0600)
	if err != nil {
		t.Fatal(err)
	}

	js := NewJSONStore(dataFile)

	// Loading existing data should fail on corrupt file
	meta := CommitMetadata{
		Commit:    "abc123",
		Prompt:    "test prompt",
		Timestamp: "2024-01-01T00:00:00Z",
	}

	err = js.Save(meta)
	if err == nil {
		t.Error("Expected error on corrupt JSON, got nil")
	}
	// Error message should mention corruption
	if err != nil && !contains(err.Error(), "corrupt") {
		t.Errorf("Expected error to mention 'corrupt', got: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsFrom(s[1:], substr)))
}

func containsFrom(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsFrom(s[1:], substr)))
}
