package store

import (
	"testing"
	"time"
)

func TestCommitMetadata_Valid(t *testing.T) {
	meta := CommitMetadata{
		Commit:    "abc123",
		Prompt:    "Add user authentication",
		Timestamp: time.Now().Format(time.RFC3339),
		Files:     []string{"auth.go"},
	}

	if meta.Commit != "abc123" {
		t.Errorf("expected commit abc123, got %s", meta.Commit)
	}
	if meta.Prompt != "Add user authentication" {
		t.Errorf("unexpected prompt: %s", meta.Prompt)
	}
}
