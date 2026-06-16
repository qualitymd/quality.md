// Package cli wires the quality.md command tree (Cobra) and runs it through
// Fang for styled help, errors, version, and shell completion.
//
// STUB: This is an early scaffold. SPECIFICATION.md defines the current file
// format and evaluation semantics; the intended CLI surface is still being
// implemented. Treat the single `check` command here as placeholder behavior,
// not as the final command surface.
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
	root.AddCommand(newCheckCmd())
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
