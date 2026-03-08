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

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Show all vibe commits with their prompts",
	Long: `Display a log of all commits with attached vibe prompts.`,
	Args: cobra.NoArgs,
	RunE: runLog,
}

func runLog(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	repo, err := git.Init(wd)
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	gitRoot, err := repo.RootPath()
	if err != nil {
		return fmt.Errorf("failed to get git root: %w", err)
	}

	dataFile := filepath.Join(gitRoot, ".vibediff", "commits.json")
	js := store.NewJSONStore(dataFile)

	commits, err := js.List()
	if err != nil {
		return fmt.Errorf("failed to load commits: %w", err)
	}

	if len(commits) == 0 {
		fmt.Println("No vibe commits found.")
		return nil
	}

	fmt.Printf("Found %d vibe commit(s):\n\n", len(commits))

	for i, meta := range commits {
		fmt.Printf("%d. %s\n", len(commits)-i, shortHash(meta.Commit))
		fmt.Printf("  %s\n", meta.Timestamp)
		if len(meta.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(meta.Tags, ", "))
		}
		fmt.Printf("  %s\n", meta.Prompt)
		if meta.Details != "" {
			// Truncate details to 80 chars for log view
			details := meta.Details
			if len(details) > 80 {
				details = details[:77] + "..."
			}
			fmt.Printf("  Details: %s\n", details)
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(logCmd)
}
