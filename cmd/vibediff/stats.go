package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vibediff/vibediff/internal/git"
	"github.com/vibediff/vibediff/internal/store"
)

var (
	statsByTag bool
	statsByModel bool
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show statistics about vibe commits",
	Long: `Display statistics about vibe commits including tag distribution,
model usage, and time-based metrics.

Examples:
  vibediff stats
  vibediff stats --by-tag
  vibediff stats --by-model`,
	RunE: runStats,
}

type tagCount struct {
	tag   string
	count int
}

type modelCount struct {
	model string
	count int
}

func runStats(cmd *cobra.Command, args []string) error {
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

	// Overall stats
	fmt.Printf("VibeDiff Statistics\n")
	fmt.Printf("===================\n\n")
	fmt.Printf("Total commits: %d\n\n", len(commits))

	// Time range
	if len(commits) > 0 {
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].Timestamp < commits[j].Timestamp
		})
		first := commits[0].Timestamp
		last := commits[len(commits)-1].Timestamp

		firstTime, _ := time.Parse(time.RFC3339, first)
		lastTime, _ := time.Parse(time.RFC3339, last)
		duration := lastTime.Sub(firstTime)

		fmt.Printf("Time range: %s to %s\n", firstTime.Format("2006-01-02"), lastTime.Format("2006-01-02"))
		fmt.Printf("Duration: %s\n\n", formatDuration(duration))
	}

	// Tag distribution
	tagCounts := make(map[string]int)
	modelCounts := make(map[string]int)
	uncategorized := 0

	for _, meta := range commits {
		if len(meta.Tags) == 0 {
			uncategorized++
		}
		for _, tag := range meta.Tags {
			tagCounts[tag]++
		}
		if meta.Model != "" {
			modelCounts[meta.Model]++
		}
	}

	// Display by tag
	if statsByTag || len(tagCounts) > 0 {
		fmt.Printf("Tags (%d unique):\n", len(tagCounts))
		if len(tagCounts) == 0 {
			fmt.Println("  No tags used yet")
		} else {
			// Sort by count descending
			sortedTags := make([]tagCount, 0, len(tagCounts))
			for tag, count := range tagCounts {
				sortedTags = append(sortedTags, tagCount{tag, count})
			}
			sort.Slice(sortedTags, func(i, j int) bool {
				return sortedTags[i].count > sortedTags[j].count
			})

			for _, tc := range sortedTags {
				bar := strings.Repeat("█", tc.count)
				if len(bar) > 20 {
					bar = bar[:20]
				}
				fmt.Printf("  %-20s %3d %s\n", tc.tag, tc.count, bar)
			}
		}
		if uncategorized > 0 {
			fmt.Printf("  (uncategorized)       %3d\n", uncategorized)
		}
		fmt.Println()
	}

	// Display by model
	if statsByModel || len(modelCounts) > 0 {
		fmt.Printf("Models (%d unique):\n", len(modelCounts))
		if len(modelCounts) == 0 {
			fmt.Println("  No models tracked yet (use --model flag)")
		} else {
			sortedModels := make([]modelCount, 0, len(modelCounts))
			for model, count := range modelCounts {
				sortedModels = append(sortedModels, modelCount{model, count})
			}
			sort.Slice(sortedModels, func(i, j int) bool {
				return sortedModels[i].count > sortedModels[j].count
			})

			for _, mc := range sortedModels {
				pct := float64(mc.count) / float64(len(commits)) * 100
				fmt.Printf("  %-20s %3d (%.0f%%)\n", mc.model, mc.count, pct)
			}
		}
		fmt.Println()
	}

	// Recent activity (last 7 days)
	fmt.Printf("Recent activity (last 7 days):\n")
	weekAgo := time.Now().AddDate(0, 0, -7)
	recentCount := 0
	for _, meta := range commits {
		t, _ := time.Parse(time.RFC3339, meta.Timestamp)
		if t.After(weekAgo) {
			recentCount++
		}
	}
	fmt.Printf("  %d commit(s)\n", recentCount)

	return nil
}

func formatDuration(d time.Duration) string {
	if d < 24*time.Hour {
		return d.String()
	}
	days := int(d.Hours() / 24)
	if days == 1 {
		return "1 day"
	}
	return fmt.Sprintf("%d days", days)
}

func init() {
	statsCmd.Flags().BoolVar(&statsByTag, "by-tag", false, "Show detailed tag breakdown")
	statsCmd.Flags().BoolVar(&statsByModel, "by-model", false, "Show detailed model breakdown")
	rootCmd.AddCommand(statsCmd)
}
