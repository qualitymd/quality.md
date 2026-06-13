package cli

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/eval"
	"github.com/qualitymd/quality.md/internal/report"
	"github.com/qualitymd/quality.md/internal/spec"
)

// newCheckCmd builds the `check` command.
//
// STUB: `check` is placeholder behavior that predates the current spec and has
// no counterpart in specs/cli.md. It collapses what the spec splits into two
// distinct tiers — the deterministic structural `lint` and the agentic
// `evaluate` — into one undifferentiated pass/fail run. It is kept only to give
// the binary something to do; expect it to be removed in favor of the spec'd
// commands.
func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [path]",
		Short: "Evaluate a QUALITY.md specification",
		Long:  "Loads a QUALITY.md spec, evaluates every requirement, and prints a report.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("file")
			if len(args) == 1 {
				path = args[0]
			}

			s, err := spec.Load(path)
			if err != nil {
				return err
			}

			results := eval.Run(cmd.Context(), s)
			report.Print(cmd.OutOrStdout(), results)

			if results.Failed() {
				return errors.New("quality check failed")
			}
			return nil
		},
	}
	cmd.Flags().StringP("file", "f", "QUALITY.md", "path to the QUALITY.md spec")
	return cmd
}
