// Package store provides persistent storage for VibeDiff metadata.
// It stores commit prompts, diffs, and associated metadata in JSON format.
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// JSONStore persists metadata to a JSON file with thread-safe access.
type JSONStore struct {
	mu       sync.RWMutex
	dataPath string
}

// NewJSONStore creates a new JSON-based store.
func NewJSONStore(dataPath string) *JSONStore {
	return &JSONStore{dataPath: dataPath}
}

// data represents the JSON file structure.
type data struct {
	Commits map[string]CommitMetadata `json:"commits"`
}

// Save stores a commit metadata entry.
func (js *JSONStore) Save(meta CommitMetadata) error {
	js.mu.Lock()
	defer js.mu.Unlock()

	// Read existing data
	d := data{Commits: make(map[string]CommitMetadata)}
	if content, err := os.ReadFile(js.dataPath); err == nil {
		if err := json.Unmarshal(content, &d); err != nil {
			return fmt.Errorf("corrupt store file %s: %w", js.dataPath, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read store: %w", err)
	}

	// Add new entry
	d.Commits[meta.Commit] = meta

	// Write back
	content, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return os.WriteFile(js.dataPath, content, 0600)
}

// Load retrieves a commit metadata entry by commit hash.
func (js *JSONStore) Load(commit string) (CommitMetadata, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	content, err := os.ReadFile(js.dataPath)
	if err != nil {
		return CommitMetadata{}, err
	}

	var d data
	if err := json.Unmarshal(content, &d); err != nil {
		return CommitMetadata{}, err
	}

	meta, ok := d.Commits[commit]
	if !ok {
		return CommitMetadata{}, os.ErrNotExist
	}

	return meta, nil
}

// List returns all commit metadata entries.
func (js *JSONStore) List() ([]CommitMetadata, error) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	content, err := os.ReadFile(js.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []CommitMetadata{}, nil
		}
		return nil, err
	}

	var d data
	if err := json.Unmarshal(content, &d); err != nil {
		return nil, err
	}

	result := make([]CommitMetadata, 0, len(d.Commits))
	for _, meta := range d.Commits {
		result = append(result, meta)
	}
	return result, nil
}
