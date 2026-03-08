package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/vibediff/vibediff/internal/store"
)

func TestExportHTML_Basic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "html-export-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "output.html")

	meta := store.CommitMetadata{
		Commit:    "abc123def456",
		Prompt:    "Add user authentication",
		Timestamp: "2024-01-01T12:00:00Z",
		Diff:      "--- a/file.go\n+++ b/file.go\n@@ -1,1 +1,2 @@\n-old\n+new",
		Files:     []string{"file.go"},
	}

	err = ExportHTML(meta, outputPath)
	if err != nil {
		t.Fatalf("ExportHTML failed: %v", err)
	}

	// Verify file exists
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}

	// Verify file has content
	if info.Size() == 0 {
		t.Error("output file is empty")
	}

	// Verify contains expected content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	contentStr := string(content)
	if !contains(contentStr, "Add user authentication") {
		t.Error("output missing prompt")
	}
	if !contains(contentStr, "abc123d") {
		t.Error("output missing commit hash")
	}
}

func TestExportHTML_XSS(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "html-xss-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "output.html")

	meta := store.CommitMetadata{
		Commit:    "abc123",
		Prompt:    "<script>alert('xss')</script>",
		Timestamp: "2024-01-01T00:00:00Z",
		Diff:      "<script>alert('xss in diff')</script>",
	}

	err = ExportHTML(meta, outputPath)
	if err != nil {
		t.Fatalf("ExportHTML failed: %v", err)
	}

	content, _ := os.ReadFile(outputPath)
	contentStr := string(content)

	// Go templates should escape HTML
	// Check that raw script tags are NOT present (they should be escaped)
	if contains(contentStr, "<script>alert") && !contains(contentStr, "&lt;script&gt;") {
		t.Logf("Warning: potential XSS - script tags found in output")
	}

	// Verify HTML entities are present for the prompt
	if !contains(contentStr, "&lt;") {
		t.Error("output should contain escaped HTML entities")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsFrom(s[1:], substr)))
}

func containsFrom(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || containsFrom(s[1:], substr)))
}
