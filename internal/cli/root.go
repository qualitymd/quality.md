// Package cli wires the quality.md command tree (Cobra) and runs it through
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
		Short: "Evaluate quality.md specifications",
		Long: "quality.md evaluates a QUALITY.md specification: a Markdown file whose " +
			"frontmatter declares factors and requirements, each scored by rules, a bash " +
			"command, or a CEL expression.",
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
