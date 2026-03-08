package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vibediff/vibediff/internal/git"
	"github.com/vibediff/vibediff/internal/store"
)

var showCmd = &cobra.Command{
	Use:   "show [commit-ref]",
	Short: "Show prompt and diff for a commit",
	Long: `Display the AI prompt and resulting diff for a commit.

If no commit is specified, shows the most recent vibe commit.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runShow,
}

func runShow(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Get commit ref
	ref := "HEAD"
	if len(args) > 0 {
		ref = args[0]
		if err := git.ValidateCommitRef(ref); err != nil {
			return fmt.Errorf("invalid commit ref: %w", err)
		}
	}

	// Resolve to hash
	repo, err := git.Init(wd)
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	hash, err := repo.RunCommand("rev-parse", ref)
	if err != nil {
		return fmt.Errorf("invalid commit ref: %w", err)
	}
	hash = strings.TrimSpace(hash)

	// Get git root for consistent .vibediff path
	gitRoot, err := repo.RootPath()
	if err != nil {
		return fmt.Errorf("failed to get git root: %w", err)
	}

	// Load metadata
	dataFile := filepath.Join(gitRoot, ".vibediff", "commits.json")
	js := store.NewJSONStore(dataFile)

	meta, err := js.Load(hash)
	if err != nil {
		return fmt.Errorf("no vibe metadata found for commit %s: %w", shortHash(hash), err)
	}

	// Display
	fmt.Printf("Commit: %s\n", shortHash(hash))
	fmt.Printf("When:   %s\n", meta.Timestamp)
	fmt.Printf("\nPrompt:\n")
	fmt.Printf("  %s\n\n", meta.Prompt)

	if meta.Details != "" {
		fmt.Printf("Details:\n")
		fmt.Printf("  %s\n\n", meta.Details)
	}

	if len(meta.Files) > 0 {
		fmt.Printf("Files: %s\n\n", strings.Join(meta.Files, ", "))
	}

	fmt.Printf("Diff:\n%s\n", meta.Diff)

	return nil
}

func init() {
	rootCmd.AddCommand(showCmd)
}
