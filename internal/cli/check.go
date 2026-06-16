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
// STUB: `check` is placeholder behavior that predates the final CLI surface. It
// loads the current frontmatter shape, traverses assessments, and records them as
// not assessed rather than performing the judgment-based evaluation described in
// SPECIFICATION.md.
func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check [path]",
		Short: "Traverse a QUALITY.md file with placeholder results",
		Long:  "Loads a QUALITY.md file, traverses every requirement assessment, and prints placeholder results.",
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
