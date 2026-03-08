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

var (
	searchQuery  string
	searchTag    string
	searchAfter  string
	searchBefore string
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search vibe commits by text or tags",
	Long: `Search through vibe commits by text in prompts/details or by tags.

Examples:
  vibediff search "authentication"
  vibediff search --tag "security"
  vibediff search "auth" --tag "api"
  vibediff search --after "2025-01-01"`,
	RunE: runSearch,
}

func matchesQuery(meta store.CommitMetadata, query string) bool {
	if query == "" {
		return true
	}
	lowerQuery := strings.ToLower(query)
	return strings.Contains(strings.ToLower(meta.Prompt), lowerQuery) ||
		strings.Contains(strings.ToLower(meta.Details), lowerQuery)
}

func matchesTags(meta store.CommitMetadata, tag string) bool {
	if tag == "" {
		return true
	}
	lowerTag := strings.ToLower(tag)
	for _, t := range meta.Tags {
		if strings.ToLower(t) == lowerTag {
			return true
		}
	}
	return false
}

func matchesDateRange(meta store.CommitMetadata, after, before string) bool {
	if after == "" && before == "" {
		return true
	}
	timestamp := meta.Timestamp
	if after != "" && timestamp < after {
		return false
	}
	if before != "" && timestamp > before {
		return false
	}
	return true
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := ""
	if len(args) > 0 {
		query = strings.Join(args, " ")
	}

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

	// Filter commits
	var results []store.CommitMetadata
	for _, meta := range commits {
		if matchesQuery(meta, query) &&
			matchesTags(meta, searchTag) &&
			matchesDateRange(meta, searchAfter, searchBefore) {
			results = append(results, meta)
		}
	}

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return nil
	}

	fmt.Printf("Found %d result(s):\n\n", len(results))

	for i, meta := range results {
		fmt.Printf("%d. %s\n", len(results)-i, shortHash(meta.Commit))
		fmt.Printf("  %s\n", meta.Timestamp)

		// Show tags if present
		if len(meta.Tags) > 0 {
			fmt.Printf("  Tags: %s\n", strings.Join(meta.Tags, ", "))
		}

		// Show prompt (highlight query match if possible)
		prompt := meta.Prompt
		if query != "" {
			// Simple highlighting - wrap match in brackets
			lowerPrompt := strings.ToLower(prompt)
			lowerQuery := strings.ToLower(query)
			if idx := strings.Index(lowerPrompt, lowerQuery); idx >= 0 {
				end := idx + len(query)
				if end > len(prompt) {
					end = len(prompt)
				}
				prompt = prompt[:idx] + "[" + prompt[idx:end] + "]" + prompt[end:]
			}
		}
		fmt.Printf("  %s\n", prompt)

		// Show truncated details
		if meta.Details != "" {
			details := meta.Details
			if len(details) > 60 {
				details = details[:57] + "..."
			}
			fmt.Printf("  Details: %s\n", details)
		}
		fmt.Println()
	}

	return nil
}

func init() {
	searchCmd.Flags().StringVar(&searchTag, "tag", "", "Filter by tag")
	searchCmd.Flags().StringVar(&searchAfter, "after", "", "Filter commits after this date (YYYY-MM-DD)")
	searchCmd.Flags().StringVar(&searchBefore, "before", "", "Filter commits before this date (YYYY-MM-DD)")
	rootCmd.AddCommand(searchCmd)
}
