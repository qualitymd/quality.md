package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

type evaluationRunFlags struct {
	latest        bool
	evaluationDir string
}

func newEvaluationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "evaluation",
		Short: "Work with QUALITY.md evaluation runs",
		RunE:  runGroupHelpOrUnknown,
	}
	cmd.AddCommand(newEvaluationCreateCmd())
	cmd.AddCommand(newEvaluationListCmd())
	cmd.AddCommand(newEvaluationStatusCmd())
	cmd.AddCommand(newEvaluationRecordNounCmd(evaluation.KindAssessmentResult, "add"))
	cmd.AddCommand(newEvaluationRecordNounCmd(evaluation.KindAnalysis, "set"))
	cmd.AddCommand(newEvaluationRecordNounCmd(evaluation.KindRecommendation, "add"))
	cmd.AddCommand(newEvaluationReportCmd())
	return cmd
}

func newEvaluationCreateCmd() *cobra.Command {
	var opts evaluation.Options
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create",
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
					*evaluation.CreateRunReceipt
				}{SchemaVersion: evaluation.SchemaVersion, CreateRunReceipt: result})
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

func newEvaluationListCmd() *cobra.Command {
	var evaluationDir string
	var state string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List evaluation runs",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := evaluation.ValidateRunState(state); err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.ListRuns("", evaluationDir, state)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			return renderRunList(cmd.OutOrStdout(), result)
		},
	}
	cmd.Flags().StringVar(&evaluationDir, "evaluation-dir", "", "override the evaluation directory")
	cmd.Flags().StringVar(&state, "state", "all", "filter runs: all, reportable, incomplete")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable run list")
	return cmd
}

func newEvaluationStatusCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "status <run>",
		Short: "Show whether an evaluation run is reportable",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			run, err := evaluation.Load(runPath)
			if err != nil {
				return mapEvaluationError(err)
			}
			status := run.EvaluationRunStatus()
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), status)
			}
			return renderEvaluationStatus(cmd, status)
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable status document")
	return cmd
}

func newEvaluationRecordNounCmd(kind evaluation.EvaluationRecordKind, writeVerb string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   string(kind),
		Short: "Work with " + string(kind) + " records",
		RunE:  runGroupHelpOrUnknown,
	}
	cmd.AddCommand(newEvaluationRecordWriteCmd(kind, writeVerb))
	cmd.AddCommand(newEvaluationRecordListCmd(kind))
	return cmd
}

func newEvaluationRecordWriteCmd(kind evaluation.EvaluationRecordKind, verb string) *cobra.Command {
	var runFlags evaluationRunFlags
	var file string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   verb + " <run>",
		Short: titleVerb(verb) + " " + string(kind) + " record payloads",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			raw, err := readPayload(cmd, file)
			if err != nil {
				return usageError(err)
			}
			result, err := evaluation.WriteRecords(kind, runPath, raw)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			_, err = fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s\n", strings.Join(result.Paths, ", "))
			return err
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().StringVar(&file, "file", "", "read judgment JSON from path, or - for stdin")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable write receipt")
	return cmd
}

func newEvaluationRecordListCmd(kind evaluation.EvaluationRecordKind) *cobra.Command {
	var runFlags evaluationRunFlags
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list <run>",
		Short: "List " + string(kind) + " records",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.ListRecords(kind, runPath)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			return renderRecordList(cmd.OutOrStdout(), result)
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable record list")
	return cmd
}

func newEvaluationReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Build and gate evaluation reports",
		RunE:  runGroupHelpOrUnknown,
	}
	cmd.AddCommand(newEvaluationReportBuildCmd())
	cmd.AddCommand(newEvaluationReportGateCmd())
	return cmd
}

func newEvaluationReportBuildCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "build <run>",
		Short: "Build report-summary.md, report.md, and report.json from evaluation records",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.BuildReport(runPath)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			_, err = fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s, %s, and %s\n", result.ReportSummaryMD, result.ReportMD, result.ReportJSON)
			return err
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable build receipt")
	return cmd
}

func newEvaluationReportGateCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var threshold string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "gate <run>",
		Short: "Gate an already-built report.json against a rating threshold",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if threshold == "" {
				return usageError(fmt.Errorf("--at-or-below is required"))
			}
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.GateReport(runPath, threshold)
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				if err := writeJSON(cmd.OutOrStdout(), result); err != nil {
					return err
				}
			} else if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Gate compared %s at or below %s: %v\n", ratingLabel(result.RatingResult), threshold, result.Pass); err != nil {
				return err
			}
			if !result.Pass {
				return silentProblems(fmt.Errorf("evaluation rating did not clear %s", threshold))
			}
			return nil
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().StringVar(&threshold, "at-or-below", "", "exit 1 when the evaluation verdict is this level or worse")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable gate result")
	return cmd
}

func bindRunFlags(cmd *cobra.Command, flags *evaluationRunFlags) {
	cmd.Flags().BoolVar(&flags.latest, "latest", false, "use the most recent evaluation run")
	cmd.Flags().StringVar(&flags.evaluationDir, "evaluation-dir", "", "override the evaluation directory when using --latest")
}

func resolveRunArg(args []string, flags evaluationRunFlags) (string, error) {
	runArg := ""
	if len(args) == 1 {
		runArg = args[0]
	}
	return evaluation.ResolveRun("", flags.evaluationDir, runArg, flags.latest)
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

func renderRunList(w io.Writer, result *evaluation.EvaluationRunList) error {
	if len(result.Runs) == 0 {
		_, err := fmt.Fprintln(w, "No evaluation runs found.")
		return err
	}
	for _, run := range result.Runs {
		if _, err := fmt.Fprintf(w, "%s\t%s\treportable=%v\trecords=%d/%d/%d\n", run.Path, run.Subject, run.Reportable, run.Counts.AssessmentResults, run.Counts.Analyses, run.Counts.Recommendations); err != nil {
			return err
		}
	}
	return nil
}

func renderRecordList(w io.Writer, result *evaluation.EvaluationRecordList) error {
	for _, record := range result.Records {
		if _, err := fmt.Fprintln(w, record); err != nil {
			return err
		}
	}
	return nil
}

func renderEvaluationStatus(cmd *cobra.Command, status evaluation.EvaluationRunStatus) error {
	out := cmd.OutOrStdout()
	reportable := "false"
	if status.Reportable {
		reportable = "true"
	}
	if _, err := fmt.Fprintf(out, "Run: %s\nReportable: %s\nRecords: %d assessment-results, %d analyses, %d recommendations\n",
		status.Path, reportable, status.Counts.AssessmentResults, status.Counts.Analyses, status.Counts.Recommendations); err != nil {
		return err
	}
	for _, gap := range status.Gaps {
		if _, err := fmt.Fprintf(out, "- %s %s: %s\n", gap.Kind, gap.Ref, gap.Detail); err != nil {
			return err
		}
	}
	return renderNextActions(cmd.ErrOrStderr(), status.NextActions)
}

func ratingLabel(result evaluation.RatingResult) string {
	if result.Kind == "not-assessed" || result.Level == "" {
		return "not assessed"
	}
	return result.Level
}

func titleVerb(verb string) string {
	if verb == "" {
		return ""
	}
	return strings.ToUpper(verb[:1]) + verb[1:]
}

func runGroupHelpOrUnknown(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return usageError(fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath()))
	}
	return cmd.Help()
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
