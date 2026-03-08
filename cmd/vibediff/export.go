package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vibediff/vibediff/internal/git"
	"github.com/vibediff/vibediff/internal/store"
	"github.com/vibediff/vibediff/internal/ui"
)

var exportCmd = &cobra.Command{
	Use:   "export [commit-ref]",
	Short: "Export vibe commit to HTML",
	Long: `Generate a shareable HTML file showing the prompt and diff.

If no commit is specified, exports the most recent vibe commit.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	ref := "HEAD"
	if len(args) > 0 {
		ref = args[0]
		if err := git.ValidateCommitRef(ref); err != nil {
			return fmt.Errorf("invalid commit ref: %w", err)
		}
	}

	repo, err := git.Init(wd)
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	hash, err := repo.RunCommand("rev-parse", ref)
	if err != nil {
		return fmt.Errorf("invalid commit ref: %w", err)
	}
	hash = strings.TrimSpace(hash)

	gitRoot, err := repo.RootPath()
	if err != nil {
		return fmt.Errorf("failed to get git root: %w", err)
	}

	dataFile := filepath.Join(gitRoot, ".vibediff", "commits.json")
	js := store.NewJSONStore(dataFile)

	meta, err := js.Load(hash)
	if err != nil {
		return fmt.Errorf("no vibe metadata found: %w", err)
	}

	outputPath := fmt.Sprintf("vibe-%s.html", shortHash(hash))
	if err := ui.ExportHTML(meta, outputPath); err != nil {
		return fmt.Errorf("failed to export HTML: %w", err)
	}

	fmt.Printf("Exported to %s\n", outputPath)

	return nil
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
