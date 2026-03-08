package git

import (
	"fmt"
	"regexp"
)

var commitRefPattern = regexp.MustCompile(`^[a-zA-Z0-9/_\.\-:]+$`)
var parentDirPattern = regexp.MustCompile(`^(\.\./|/\.\.|/\.\./|\.\./)`)

// ValidateCommitRef checks if a commit reference is safe to pass to git.
// It only allows alphanumeric characters, slashes, dots, dashes, and colons
// to prevent command injection attacks.
func ValidateCommitRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("commit ref cannot be empty")
	}
	if !commitRefPattern.MatchString(ref) {
		return fmt.Errorf("invalid commit ref format: %q", ref)
	}
	// Prevent parent directory traversal (../)
	if parentDirPattern.MatchString(ref) {
		return fmt.Errorf("invalid commit ref format: %q (parent directory traversal not allowed)", ref)
	}
	return nil
}
