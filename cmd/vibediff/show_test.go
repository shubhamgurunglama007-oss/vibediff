package main

import (
	"testing"
)

func TestShortHash(t *testing.T) {
	tests := []struct {
		name     string
		hash     string
		expected string
	}{
		{"normal hash", "abc123def456", "abc123d"},
		{"exact 7 chars", "abc1234", "abc1234"},
		{"short hash", "abc", "abc"},
		{"empty string", "", ""},
		{"with whitespace", "  abc123def456  ", "abc123d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shortHash(tt.hash)
			if result != tt.expected {
				t.Errorf("shortHash() = %q, want %q", result, tt.expected)
			}
		})
	}
}
