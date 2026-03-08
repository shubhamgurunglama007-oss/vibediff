package main

import (
	"testing"
)

func TestValidatePrompt(t *testing.T) {
	tests := []struct {
		name    string
		prompt  string
		wantErr bool
	}{
		{"valid prompt", "Add user authentication", false},
		{"valid with extra spaces", "  Add user authentication  ", false},
		{"empty string", "", true},
		{"only spaces", "   ", true},
		{"only tabs", "\t\t", true},
		{"only newlines", "\n\n", true},
		{"mixed whitespace", " \t\n ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validatePrompt(tt.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePrompt_SizeLimit(t *testing.T) {
	// Test prompt size limit (100KB)
	largePrompt := string(make([]byte, 100001))
	_, err := validatePrompt(largePrompt)
	if err == nil {
		t.Error("Expected error for oversized prompt, got nil")
	}
}

func TestValidateDetails(t *testing.T) {
	tests := []struct {
		name     string
		details  string
		wantErr  bool
		wantSkip bool
	}{
		{"valid details", "Implemented OAuth2 with Google", false, false},
		{"valid with extra spaces", "  Implemented OAuth2  ", false, false},
		{"empty string", "", false, true},
		{"only spaces", "   ", false, true},
		{"multiline details", "Line 1\nLine 2\nLine 3", false, false},
		{"only tabs", "\t\t", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateDetails(tt.details)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateDetails() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantSkip && result != "" {
				t.Errorf("validateDetails() should skip empty details, got = %v", result)
			}
		})
	}
}

func TestValidateDetails_SizeLimit(t *testing.T) {
	// Test details size limit (500KB)
	largeDetails := string(make([]byte, 500001))
	_, err := validateDetails(largeDetails)
	if err == nil {
		t.Error("Expected error for oversized details, got nil")
	}
}

func TestValidateDetails_TrimsWhitespace(t *testing.T) {
	input := "  Details with spaces  "
	result, err := validateDetails(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Details with spaces" {
		t.Errorf("validateDetails() trim = %q, want %q", result, "Details with spaces")
	}
}
