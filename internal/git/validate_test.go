package git

import "testing"

func TestValidateCommitRef(t *testing.T) {
	tests := []struct {
		name    string
		ref     string
		wantErr bool
	}{
		{"valid HEAD", "HEAD", false},
		{"valid hash", "abc123def456", false},
		{"valid tag", "v1.0.0", false},
		{"valid with slash", "origin/main", false},
		{"valid with dash", "feature-branch", false},
		{"valid with dot", "refs/tags/v1.0.0", false},
		{"empty string", "", true},
		{"with space", "HEAD~1", true},  // ~ not allowed
		{"pipe injection", "HEAD|rm -rf /", true},
		{"command injection", "--upload-pack=;touch /tmp/pwned", true},
		{"path traversal", "../../../etc/passwd", true},
		{"semicolon injection", "HEAD;echo pwned", true},
		{"backtick injection", "HEAD`whoami`", true},
		{"dollar sign", "$HOME", true},
		{"newline injection", "HEAD\nrm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommitRef(tt.ref)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCommitRef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
