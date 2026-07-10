package cli

import (
	"errors"
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/evaluator"
	"github.com/qualitymd/quality.md/internal/runner"
)

func newEvaluationRunCmd() *cobra.Command {
	var opts runner.Options
	var jsonOutput bool
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Execute a complete evaluation run with the deterministic runner",
		Long: "Execute a complete evaluation run: the runner owns the work graph, " +
			"invokes the selected evaluator for bounded judgment work units, persists " +
			"the authoritative evaluation.json run artifact, and builds the Markdown reports.",
		Example: "  qualitymd evaluation run\n" +
			"  qualitymd evaluation run --area area:docs --evaluator claude\n" +
			"  qualitymd evaluation run --dry-run --json\n" +
			"  qualitymd evaluation run --resume .quality/evaluations/0007-full-eval",
		Args: usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if dryRun {
				return runEvaluationDryRun(cmd, opts, jsonOutput)
			}
			return runEvaluationRun(cmd, opts, jsonOutput)
		},
	}
	cmd.Flags().StringVar(&opts.Model, "model", "", "QUALITY.md file to evaluate")
	cmd.Flags().StringVar(&opts.EvaluationDir, "evaluation-dir", "", "override the model-relative evaluation directory")
	cmd.Flags().StringVar(&opts.Area, "area", "", "canonical area reference for the evaluation scope")
	cmd.Flags().StringArrayVar(&opts.Factors, "factor", nil, "canonical factor reference for a scoped evaluation; repeatable")
	cmd.Flags().StringVar(&opts.Evaluator, "evaluator", "", "evaluator to use: auto (default), a built-in name, or a configured profile")
	cmd.Flags().StringVar(&opts.Resume, "resume", "", "resume an existing run from its evaluation.json")
	cmd.Flags().StringVar(&opts.EvaluatorResult, "evaluator-result", "",
		"submit a harness result envelope for the awaiting work request, from a file or - for stdin (requires --resume)")
	cmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "preview the resolved run without invoking an evaluator or writing evaluation data")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable run receipt")
	return cmd
}

func runEvaluationDryRun(cmd *cobra.Command, opts runner.Options, jsonOutput bool) error {
	preview, err := runner.DryRun(opts)
	if err != nil {
		return mapRunnerError(cmd, err, jsonOutput)
	}
	if jsonOutput {
		return writeJSON(cmd.OutOrStdout(), preview)
	}
	out := cmd.ErrOrStderr()
	if _, err := fmt.Fprintf(out, "Would evaluate %s at %s\n", preview.Model, preview.ExpectedRunPath); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(out, "Evaluator: %s (%s); concurrency: %d; work units: %d (%d evaluator-backed)\n",
		preview.Evaluator, preview.EvaluatorReason, preview.Concurrency,
		preview.WorkUnits.Total, preview.WorkUnits.EvaluatorUnits); err != nil {
		return err
	}
	return renderNextActions(out, preview.NextActions)
}

func runEvaluationRun(cmd *cobra.Command, opts runner.Options, jsonOutput bool) error {
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	opts.Stderr = cmd.ErrOrStderr()
	opts.Stdin = cmd.InOrStdin()
	result, err := runner.Run(ctx, opts)
	if err != nil {
		return mapRunnerError(cmd, err, jsonOutput)
	}
	if jsonOutput {
		if err := writeJSON(cmd.OutOrStdout(), result); err != nil {
			return err
		}
	} else {
		if err := renderEvaluationRunResult(cmd, result); err != nil {
			return err
		}
	}
	// An awaiting-evaluator checkpoint is expected progress, not a failure:
	// the command exits successfully with the pending request in the receipt.
	if result.Status != runner.StatusCompleted && result.Status != runner.StatusAwaitingEvaluator {
		return silentProblems(fmt.Errorf("evaluation run %s", result.Status))
	}
	return nil
}

func renderEvaluationRunResult(cmd *cobra.Command, result *runner.Result) error {
	out := cmd.ErrOrStderr()
	switch result.Status {
	case runner.StatusAwaitingEvaluator:
		if err := renderAwaitingEvaluator(out, result); err != nil {
			return err
		}
	case runner.StatusCompleted:
		rating := "not assessed"
		if result.RatingResult != nil && result.RatingResult.Level != "" {
			rating = result.RatingResult.Level
		}
		if _, err := fmt.Fprintf(out, "Evaluation completed: %s (rating: %s)\nReport: %s\n",
			result.Path, rating, result.ReportMD); err != nil {
			return err
		}
	default:
		if _, err := fmt.Fprintf(out, "Evaluation %s: %s\n", result.Status, result.Path); err != nil {
			return err
		}
		if result.Failure != nil {
			if _, err := fmt.Fprintf(out, "Failure: %s: %s\n", result.Failure.Category, result.Failure.Detail); err != nil {
				return err
			}
		}
	}
	return renderNextActions(out, result.NextActions)
}

func renderAwaitingEvaluator(out io.Writer, result *runner.Result) error {
	if result.EvaluatorRequest != nil {
		if _, err := fmt.Fprintf(out, "Awaiting harness judgment: %s (request %s, attempt %d)\n",
			result.EvaluatorRequest.WorkUnitID, result.EvaluatorRequest.RequestID, result.EvaluatorRequest.Attempt); err != nil {
			return err
		}
	}
	if result.Failure != nil {
		if _, err := fmt.Fprintf(out, "Previous attempt: %s: %s\n", result.Failure.Category, result.Failure.Detail); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(out, "Run with --json to receive the bounded work request.")
	return err
}

// mapRunnerError classifies runner errors onto CLI exit-code categories and,
// under --json, emits a machine-readable error receipt carrying the failure
// category.
func mapRunnerError(cmd *cobra.Command, err error, jsonOutput bool) error {
	var usageErr *runner.UsageError
	if errors.As(err, &usageErr) {
		return usageError(err)
	}
	var selErr *evaluator.SelectionError
	if errors.As(err, &selErr) {
		return emitRunnerFailure(cmd, string(selErr.Category), selErr.Error(), err, jsonOutput)
	}
	var runErr *runner.RunError
	if errors.As(err, &runErr) {
		return emitRunnerFailure(cmd, string(runErr.Category), runErr.Detail, err, jsonOutput)
	}
	return mapEvaluationError(err)
}

func emitRunnerFailure(cmd *cobra.Command, category, detail string, err error, jsonOutput bool) error {
	if jsonOutput {
		receipt := struct {
			SchemaVersion int    `json:"schemaVersion"`
			Status        string `json:"status"`
			Failure       struct {
				Category string `json:"category"`
				Detail   string `json:"detail"`
			} `json:"failure"`
		}{SchemaVersion: 3, Status: "failed"}
		receipt.Failure.Category = category
		receipt.Failure.Detail = detail
		if writeErr := writeJSON(cmd.OutOrStdout(), receipt); writeErr != nil {
			return writeErr
		}
		return silentProblems(err)
	}
	return &codedError{code: ExitProblems, err: err}
}
