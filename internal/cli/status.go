package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/status"
)

func newStatusCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "status [path]",
		Short: "Show a QUALITY.md project-state snapshot",
		Example: "  qualitymd status\n" +
			"  qualitymd status docs/QUALITY.md\n" +
			"  qualitymd status --json",
		Args: usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "QUALITY.md"
			if len(args) == 1 {
				path = args[0]
			}
			if path == "-" {
				return usageError(fmt.Errorf("status does not read from stdin; pass a file path"))
			}
			snapshot, err := status.Snapshot(status.Options{Path: path})
			if err != nil {
				return err
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), snapshot)
			}
			if err := renderStatusHuman(cmd, snapshot); err != nil {
				return err
			}
			return renderNextActions(cmd.ErrOrStderr(), snapshot.NextActions)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable status snapshot")
	return cmd
}

func renderStatusHuman(cmd *cobra.Command, snapshot *status.SnapshotResult) error {
	out := cmd.OutOrStdout()
	model := "absent"
	if snapshot.Model.Present {
		if snapshot.Model.Valid {
			model = "present, valid"
		} else {
			model = fmt.Sprintf("present, invalid (%d error(s), %d warning(s))", snapshot.Model.Lint.Summary.Errors, snapshot.Model.Lint.Summary.Warnings)
		}
	}
	if _, err := fmt.Fprintf(out, "Status\n- QUALITY.md: %s\n", model); err != nil {
		return err
	}
	if snapshot.Model.Shape != nil {
		shape := snapshot.Model.Shape
		if _, err := fmt.Fprintf(out, "- Model: %d target(s), %d factor(s), %d requirement(s)\n", shape.Targets, shape.Factors, shape.Requirements); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(out, "- Evaluation history: %d run(s), %d incomplete, %d stale, %d active recommendation(s)\n",
		snapshot.Evaluations.Runs,
		snapshot.Evaluations.Summary.Incomplete,
		snapshot.Evaluations.Summary.Stale,
		snapshot.Evaluations.Summary.ActiveRecommendations,
	); err != nil {
		return err
	}
	_, err := fmt.Fprintf(out, "- Readiness: %s\n", snapshot.Readiness)
	return err
}
