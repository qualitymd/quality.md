package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/lint"
)

func newLintCmd() *cobra.Command {
	var jsonOutput bool
	var fix bool
	cmd := &cobra.Command{
		Use:   "lint [path]",
		Short: "Validate a QUALITY.md file",
		Args:  usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "QUALITY.md"
			if len(args) == 1 {
				path = args[0]
			}
			if path == "-" {
				return fmt.Errorf("lint does not read from stdin yet; pass a file path")
			}

			var (
				result lint.Result
				err    error
			)
			if fix {
				result, err = lint.Fix(path)
			} else {
				result, err = lint.Check(path)
			}
			if err != nil {
				return err
			}

			if jsonOutput {
				data, err := result.JSON()
				if err != nil {
					return err
				}
				if _, err := cmd.OutOrStdout().Write(append(data, '\n')); err != nil {
					return err
				}
			} else if err := renderLintHuman(cmd, result); err != nil {
				return err
			}
			if err := result.Err(); err != nil {
				return silentProblems(err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable JSON lint result")
	cmd.Flags().BoolVar(&fix, "fix", false, "apply deterministic in-place repairs")
	return cmd
}

func renderLintHuman(cmd *cobra.Command, result lint.Result) error {
	out := cmd.OutOrStdout()
	if result.Summary.Fixed > 0 {
		if _, err := fmt.Fprintf(out, "Applied %d repair(s).\n", result.Summary.Fixed); err != nil {
			return err
		}
	}
	if len(result.Findings) == 0 {
		_, err := fmt.Fprintf(out, "%s is valid.\n", result.Path)
		return err
	}
	for _, finding := range result.Findings {
		if _, err := fmt.Fprintf(out, "%s %s: %s (%s)\n", finding.Severity, finding.RuleID, finding.Message, finding.Location.Label); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(out, "\n%d error(s), %d warning(s).\n", result.Summary.Errors, result.Summary.Warnings)
	return err
}
