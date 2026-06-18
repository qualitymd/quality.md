package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

func newEvaluationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evaluation",
		Short: "Work with QUALITY.md evaluation runs",
	}
	cmd.AddCommand(newEvaluationCreateRunCmd())
	cmd.AddCommand(newEvaluationAddRecordCmd())
	cmd.AddCommand(newEvaluationSetPlannedCoverageCmd())
	cmd.AddCommand(newEvaluationShowStatusCmd())
	cmd.AddCommand(newEvaluationBuildReportCmd())
	return cmd
}

func newEvaluationCreateRunCmd() *cobra.Command {
	var opts evaluation.Options
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create-run",
		Short: "Create a numbered evaluation run folder",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			result, err := evaluation.CreateRun(opts)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), struct {
					SchemaVersion int `json:"schemaVersion"`
					*evaluation.CreateRunResult
				}{SchemaVersion: evaluation.SchemaVersion, CreateRunResult: result})
			}
			if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "Created %s\n", result.Path); err != nil {
				return err
			}
			return renderNextActions(cmd.ErrOrStderr(), result.NextActions)
		},
	}
	cmd.Flags().StringVar(&opts.Narrowing, "narrowing", "", "optional scope slug for the run folder")
	cmd.Flags().StringVar(&opts.Subject, "subject", "", "QUALITY.md file to snapshot")
	cmd.Flags().StringVar(&opts.EvaluationDir, "evaluation-dir", "", "override the evaluation directory")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable run creation receipt")
	return cmd
}

func newEvaluationAddRecordCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-record",
		Short: "Write one evaluation record",
	}
	cmd.AddCommand(newEvaluationAddRecordKindCmd(evaluation.KindAssessment))
	cmd.AddCommand(newEvaluationAddRecordKindCmd(evaluation.KindAnalysis))
	cmd.AddCommand(newEvaluationAddRecordKindCmd(evaluation.KindRecommendation))
	return cmd
}

func newEvaluationAddRecordKindCmd(kind evaluation.WriteKind) *cobra.Command {
	var file string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   string(kind) + " <run>",
		Short: "Write one " + string(kind) + " record",
		Args:  usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayload(cmd, file)
			if err != nil {
				return usageError(err)
			}
			result, err := evaluation.AddRecord(kind, args[0], raw)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			_, err = fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s\n", result.Path)
			return err
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "read judgment JSON from path, or - for stdin")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable write receipt")
	return cmd
}

func newEvaluationSetPlannedCoverageCmd() *cobra.Command {
	var file string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "set-planned-coverage <run>",
		Short: "Write planned coverage metadata for an evaluation run",
		Args:  usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, err := readPayload(cmd, file)
			if err != nil {
				return usageError(err)
			}
			result, err := evaluation.SetPlannedCoverage(args[0], raw)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s\n", result.Path); err != nil {
				return err
			}
			return renderNextActions(cmd.ErrOrStderr(), result.NextActions)
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "read planned coverage JSON from path, or - for stdin")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable write receipt")
	return cmd
}

func newEvaluationShowStatusCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "show-status <run>",
		Short: "Show whether an evaluation run is reportable",
		Args:  usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			run, err := evaluation.Load(args[0])
			if err != nil {
				return mapEvaluationError(err)
			}
			status := run.Status()
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), status)
			}
			return renderEvaluationStatus(cmd, status)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable status document")
	return cmd
}

func newEvaluationBuildReportCmd() *cobra.Command {
	var jsonOutput bool
	var failAtOrBelow string
	cmd := &cobra.Command{
		Use:   "build-report <run>",
		Short: "Build report.md and report.json from evaluation records",
		Args:  usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			levels, err := evaluation.ScaleLevels(args[0])
			if err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.BuildReport(args[0])
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				if err := writeJSON(cmd.OutOrStdout(), result); err != nil {
					return err
				}
			} else if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s and %s\n", result.ReportMD, result.ReportJSON); err != nil {
				return err
			}
			pass, err := evaluation.Gate(result, levels, failAtOrBelow)
			if err != nil {
				return mapEvaluationError(err)
			}
			if failAtOrBelow != "" {
				rating := "not assessed"
				if result.Rating != nil {
					rating = *result.Rating
				}
				if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "Gate compared %s at or below %s: %v\n", rating, failAtOrBelow, pass); err != nil {
					return err
				}
				if !pass {
					return silentProblems(fmt.Errorf("evaluation rating did not clear %s", failAtOrBelow))
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable build receipt")
	cmd.Flags().StringVar(&failAtOrBelow, "fail-at-or-below", "", "exit 1 when the root rating is this level or worse")
	return cmd
}

func readPayload(cmd *cobra.Command, file string) ([]byte, error) {
	if file != "" {
		if file == "-" {
			return io.ReadAll(cmd.InOrStdin())
		}
		return os.ReadFile(file)
	}
	in := cmd.InOrStdin()
	if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		return nil, fmt.Errorf("pass --file <path> or pipe JSON on standard input")
	}
	return io.ReadAll(in)
}

func renderEvaluationStatus(cmd *cobra.Command, status evaluation.Status) error {
	out := cmd.OutOrStdout()
	reportable := "false"
	if status.Reportable {
		reportable = "true"
	}
	if _, err := fmt.Fprintf(out, "Run: %s\nReportable: %s\nRecords: %d assessments, %d analyses, %d recommendations\n",
		status.Path, reportable, status.Counts.Assessments, status.Counts.Analyses, status.Counts.Recommendations); err != nil {
		return err
	}
	for _, gap := range status.Gaps {
		if _, err := fmt.Fprintf(out, "- %s %s: %s\n", gap.Kind, gap.Ref, gap.Detail); err != nil {
			return err
		}
	}
	return renderNextActions(cmd.ErrOrStderr(), status.NextActions)
}

func writeJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func mapEvaluationError(err error) error {
	var usageErr *evaluation.UsageError
	if errors.As(err, &usageErr) {
		return usageError(usageErr)
	}
	return err
}
