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
	cmd.AddCommand(newEvaluationDataCmd())
	cmd.AddCommand(newEvaluationReportCmd())
	return cmd
}

func newEvaluationDataCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "data",
		Short: "Work with Evaluation v2 structured data",
		RunE:  runGroupHelpOrUnknown,
	}
	cmd.AddCommand(newEvaluationDataSetCmd())
	cmd.AddCommand(newEvaluationDataListCmd())
	cmd.AddCommand(newEvaluationDataGetCmd())
	cmd.AddCommand(newEvaluationDataKindsCmd())
	cmd.AddCommand(newEvaluationDataExampleCmd())
	return cmd
}

func newEvaluationDataSetCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var file string
	var jsonOutput bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "set <run>",
		Short: "Validate and persist one Evaluation v2 JSON payload",
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
			result, err := evaluation.SetData(runPath, raw, evaluation.DataSetOptions{DryRun: dryRun})
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			verb := "Wrote"
			if dryRun {
				verb = "Would write"
			}
			_, err = fmt.Fprintf(cmd.ErrOrStderr(), "%s %s\n", verb, result.Path)
			return err
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().StringVar(&file, "file", "", "read Evaluation v2 JSON from path, or - for stdin")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "validate and report intended write without persisting")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable write receipt")
	return cmd
}

func newEvaluationDataListCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var kind string
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list <run>",
		Short: "List stored Evaluation v2 JSON payloads",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			result, err := evaluation.ListData(runPath, evaluation.DataKind(kind))
			if err != nil {
				return mapEvaluationError(err)
			}
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			for _, artifact := range result.Artifacts {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\n", artifact.Kind, artifact.Path); err != nil {
					return err
				}
			}
			return nil
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().StringVar(&kind, "kind", "", "filter by Evaluation v2 data kind")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable data list")
	return cmd
}

func newEvaluationDataGetCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var kind string
	var areaRef string
	var factorRef string
	var requirementRef string
	var selector string
	var jsonFlag bool
	cmd := &cobra.Command{
		Use:   "get <run>",
		Short: "Print one stored Evaluation v2 JSON payload",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonFlag {
				return usageError(fmt.Errorf("evaluation data get already emits JSON on stdout; rerun without --json"))
			}
			runPath, err := resolveRunArg(args, runFlags)
			if err != nil {
				return mapEvaluationError(err)
			}
			raw, _, err := evaluation.GetData(runPath, evaluation.DataQuery{
				Kind:           evaluation.DataKind(kind),
				AreaRef:        areaRef,
				FactorRef:      factorRef,
				RequirementRef: requirementRef,
				Selector:       selector,
				AllowCLIOwned:  true,
			})
			if err != nil {
				return mapEvaluationError(err)
			}
			_, err = cmd.OutOrStdout().Write(raw)
			return err
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().StringVar(&kind, "kind", "", "Evaluation v2 data kind")
	cmd.Flags().StringVar(&areaRef, "area", "", "Area ref for Area-scoped payloads")
	cmd.Flags().StringVar(&factorRef, "factor", "", "Factor ref for Factor-scoped payloads")
	cmd.Flags().StringVar(&requirementRef, "requirement", "", "Requirement ref for Requirement-scoped payloads")
	cmd.Flags().StringVar(&selector, "selector", "", "optional sub-result selector")
	cmd.Flags().BoolVar(&jsonFlag, "json", false, "not supported: data get already emits JSON")
	return cmd
}

func newEvaluationDataKindsCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "kinds",
		Short: "List Evaluation v2 data kinds",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			result := evaluation.EvaluationDataKinds()
			if jsonOutput {
				return writeJSON(cmd.OutOrStdout(), result)
			}
			for _, kind := range result.Kinds {
				writable := "cli-owned"
				if kind.AgentWritable {
					writable = "agent-writable"
				}
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", kind.Kind, writable, kind.Description); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable data kind list")
	return cmd
}

func newEvaluationDataExampleCmd() *cobra.Command {
	var jsonFlag bool
	cmd := &cobra.Command{
		Use:   "example <kind>",
		Short: "Print a complete Evaluation v2 example JSON payload",
		Args:  usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if jsonFlag {
				return usageError(fmt.Errorf("evaluation data example already emits JSON on stdout; rerun without --json"))
			}
			raw, err := evaluation.DataExample(evaluation.DataKind(args[0]))
			if err != nil {
				return mapEvaluationError(err)
			}
			_, err = cmd.OutOrStdout().Write(raw)
			return err
		},
	}
	cmd.Flags().BoolVar(&jsonFlag, "json", false, "not supported: data example already emits JSON")
	return cmd
}

func newEvaluationCreateCmd() *cobra.Command {
	var opts evaluation.Options
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "create [model]",
		Short: "Create a numbered evaluation run folder",
		Args:  usage(cobra.RangeArgs(0, 1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				if opts.Model != "" {
					return usageError(fmt.Errorf("pass a model argument or --model, not both"))
				}
				opts.Model = args[0]
			}
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
	cmd.Flags().StringVar(&opts.Model, "model", "", "QUALITY.md file to snapshot")
	cmd.Flags().StringVar(&opts.ResolveDir, "evaluation-dir", "", "override the evaluation directory")
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
			run, err := evaluation.Inspect(runPath)
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
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable status document")
	return cmd
}

func newEvaluationReportCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Build evaluation reports",
		RunE:  runGroupHelpOrUnknown,
	}
	cmd.AddCommand(newEvaluationReportBuildCmd())
	return cmd
}

func newEvaluationReportBuildCmd() *cobra.Command {
	var runFlags evaluationRunFlags
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "build <run>",
		Short: "Build deterministic evaluation reports",
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
			if result.EvaluationOutputResult != "" {
				_, err = fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s and %s\n", result.EvaluationOutputResult, result.ReportMD)
				return err
			}
			_, err = fmt.Fprintf(cmd.ErrOrStderr(), "Wrote %s, %s, and %s\n", result.ReportSummaryMD, result.ReportMD, result.ReportJSON)
			return err
		},
	}
	bindRunFlags(cmd, &runFlags)
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable build receipt")
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

func renderRunList(w io.Writer, result *evaluation.RunList) error {
	if len(result.Runs) == 0 {
		_, err := fmt.Fprintln(w, "No evaluation runs found.")
		return err
	}
	for _, run := range result.Runs {
		if _, err := fmt.Fprintf(w, "%s\t%s\treportable=%v\trecords=%d/%d/%d\n", run.Path, run.RootArea, run.Reportable, run.Counts.AssessmentResults, run.Counts.Analyses, run.Counts.Recommendations); err != nil {
			return err
		}
	}
	return nil
}

func renderEvaluationStatus(cmd *cobra.Command, status evaluation.RunStatus) error {
	out := cmd.OutOrStdout()
	reportable := "false"
	if status.Reportable {
		reportable = "true"
	}
	if _, err := fmt.Fprintf(out, "Run: %s\nReportable: %s\nRecords: %d assessments, %d analyses, %d recommendations\n",
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
