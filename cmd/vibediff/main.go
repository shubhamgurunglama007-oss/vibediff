package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vibediff",
	Short: "Git-native versioning layer for prompts and AI outputs",
	Long: `VibeDiff captures AI prompts alongside the code changes they generate.

Use "vibediff commit" to attach a prompt to your commit.`,
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
