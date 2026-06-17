// Package cli wires the QUALITY.md command tree (Cobra) and runs it through
// Fang for styled help, errors, version, and shell completion.
package cli

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

// Build-time metadata. goreleaser overrides these via -ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "qualitymd",
		Short: "Work with QUALITY.md files",
		Long: "qualitymd works with QUALITY.md files: Markdown documents whose YAML " +
			"frontmatter declares a quality model with a ratingScale, targets, factors, " +
			"requirements, and one assessment per requirement.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(newInitCmd())
	return root
}

// Execute builds the command tree and runs it. It exits non-zero on error;
// Fang renders the error.
func Execute() {
	if err := fang.Execute(
		context.Background(),
		newRootCmd(),
		fang.WithVersion(version),
		fang.WithCommit(commit),
	); err != nil {
		os.Exit(1)
	}
}
