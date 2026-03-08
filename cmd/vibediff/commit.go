package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibediff/vibediff/internal/git"
	"github.com/vibediff/vibediff/internal/store"
)

var (
	commitDetails  string
	noDetailsLimit bool   // When true, bypass the 500KB details limit
	commitTags     string // Comma-separated tags
	commitModel    string // AI model used (e.g., claude-4, gpt-4o)
)

// vibeREADME is the content of .vibediff/README.md
const vibeREADME = `# VibeDiff Metadata

This directory stores AI prompts and context for commits made with VibeDiff.

## What is VibeDiff?

VibeDiff is a Git-native tool that attaches your AI prompts to your commits.
It captures the **intent** behind your code changes.

## Commands

 vibediff show           # Show the most recent vibe commit
 vibediff show <hash>    # Show a specific commit
 vibediff log            # List all vibe commits
 vibediff export         # Export to HTML for sharing

## File Structure

 .vibediff/
 ├── README.md          # This file
 └── commits.json       # Stores prompts + diffs (0600 permissions)

## Privacy

 - All data is stored locally (no cloud)
 - Private file permissions (0600)
 - Can be added to .gitignore if desired

## Learn More

 https://github.com/vibediff/vibediff
 Run 'vibediff --help' for all commands
`

func ensureVibeReadme(dir string) error {
	readmePath := filepath.Join(dir, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		return nil // Already exists
	}
	return os.WriteFile(readmePath, []byte(vibeREADME), 0644)
}

func validatePrompt(prompt string) (string, error) {
	trimmed := strings.TrimSpace(prompt)
	if trimmed == "" {
		return "", fmt.Errorf("prompt cannot be empty")
	}
	// Limit prompt size to prevent DoS via disk space
	const maxPromptSize = 100000 // 100KB
	if len(trimmed) > maxPromptSize {
		return "", fmt.Errorf("prompt too large (max %d bytes)", maxPromptSize)
	}
	return trimmed, nil
}

func validateDetails(details string) (string, error) {
	trimmed := strings.TrimSpace(details)
	if trimmed == "" {
		return "", nil // Empty details is fine
	}
	// Limit details size to prevent DoS via disk space
	// Can be bypassed with --no-details-limit flag
	const maxDetailsSize = 500000 // 500KB
	if !noDetailsLimit && len(trimmed) > maxDetailsSize {
		return "", fmt.Errorf("details too large (max %d bytes, use --no-details-limit to bypass)", maxDetailsSize)
	}
	return trimmed, nil
}

var commitCmd = &cobra.Command{
	Use:   "commit <prompt>",
	Short: "Stage files, commit, and attach prompt",
	Long: `Stage all changes, create a commit, and attach the AI prompt
that generated these changes.

The prompt becomes the commit message (prefixed with "vibe:").
Use the --details flag to add extended context that won't appear
in the git commit message but will be stored with the commit.

By default, --details is limited to 500KB. Use --no-details-limit
to bypass this restriction for very large AI conversation logs.

Examples:
  vibediff commit "Add user authentication"
  vibediff commit "Add auth" --details "Implemented OAuth2 with Google, added refresh token rotation, and secure session management"
  vibediff commit "Refactor" --details "$(cat full-ai-conversation.txt)" --no-details-limit`,
	Args: cobra.ExactArgs(1),
	RunE: runCommit,
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return nil
	}
	tags := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func runCommit(cmd *cobra.Command, args []string) error {
	prompt, err := validatePrompt(args[0])
	if err != nil {
		return err
	}

	details, err := validateDetails(commitDetails)
	if err != nil {
		return err
	}

	tags := parseTags(commitTags)

	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Initialize git repo
	repo, err := git.Init(wd)
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// Stage all changes
	if err := repo.Add("."); err != nil {
		return fmt.Errorf("failed to stage files: %w", err)
	}

	// Get the diff
	diff, err := repo.Diff("staged")
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	// Get changed files
	files, err := repo.GetChangedFiles("staged")
	if err != nil {
		return fmt.Errorf("failed to get changed files: %w", err)
	}

	// Create commit
	commitMsg := fmt.Sprintf("vibe: %s\n\n[View details: vibediff show]", prompt)
	if err := repo.Commit(commitMsg); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Get the commit hash
	hash, err := repo.GetCurrentHead()
	if err != nil {
		return fmt.Errorf("failed to get commit hash: %w", err)
	}

	// Store metadata - use git root instead of cwd
	gitRoot, err := repo.RootPath()
	if err != nil {
		return fmt.Errorf("failed to get git root: %w", err)
	}

	vibeDir := filepath.Join(gitRoot, ".vibediff")
	if err := os.MkdirAll(vibeDir, 0700); err != nil {
		return fmt.Errorf("failed to create .vibediff directory: %w", err)
	}

	// Create README.md for discoverability
	if err := ensureVibeReadme(vibeDir); err != nil {
		return fmt.Errorf("failed to create README: %w", err)
	}

	dataFile := filepath.Join(vibeDir, "commits.json")
	js := store.NewJSONStore(dataFile)

	meta := store.CommitMetadata{
		Commit:    strings.TrimSpace(hash),
		Prompt:    prompt,
		Details:   details,
		Timestamp: time.Now().Format(time.RFC3339),
		Diff:      diff,
		Files:     files,
		Tags:      tags,
		Model:     commitModel,
	}

	if err := js.Save(meta); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	fmt.Printf("Committed %s\n", shortHash(hash))
	fmt.Printf("Prompt: %s\n", prompt)
	if details != "" {
		fmt.Printf("Details: %s\n", details)
	}
	if len(tags) > 0 {
		fmt.Printf("Tags: %s\n", strings.Join(tags, ", "))
	}

	return nil
}

func init() {
	commitCmd.Flags().StringVar(&commitDetails, "details", "", "Extended details about the commit (optional)")
	commitCmd.Flags().BoolVar(&noDetailsLimit, "no-details-limit", false, "Bypass the 500KB limit on details (use for large AI conversations)")
	commitCmd.Flags().StringVar(&commitTags, "tags", "", "Comma-separated tags for categorization (e.g., 'auth,security,api')")
	commitCmd.Flags().StringVar(&commitModel, "model", "", "AI model used (e.g., claude-4, gpt-4o, o1)")
	rootCmd.AddCommand(commitCmd)
}
