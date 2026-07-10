package runner

import (
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

// IsRunnerRun reports whether the run folder is a runner-created run backed
// by evaluation.json.
func IsRunnerRun(runPath string) bool {
	abs, err := filepath.Abs(runPath)
	if err != nil {
		return false
	}
	return NewStore(abs).Exists()
}

// Status computes the RunStatus receipt for a runner-created run. The
// evaluation package reads evaluation.json directly, so this is the same
// status every other command surface reports.
func Status(runPath, displayPath string) (*evaluation.RunStatus, error) {
	run, err := evaluation.InspectWithDisplay(runPath, displayPath)
	if err != nil {
		return nil, err
	}
	status := run.Status()
	return &status, nil
}

// RebuildReports re-renders the Markdown report tree of a runner-created run
// from its evaluation.json and refreshes the artifact's outputs section.
func RebuildReports(runPath, displayPath string) (*evaluation.BuildReportReceipt, error) {
	abs, err := filepath.Abs(runPath)
	if err != nil {
		return nil, err
	}
	store := NewStore(abs)
	artifact, err := store.Load()
	if err != nil {
		return nil, err
	}
	payloads := append([]map[string]any{artifact.Manifest.ManifestPayload()}, artifact.Results.PayloadList()...)
	result, gaps, err := evaluation.BuildReportFromPayloads(abs, displayPath, payloads)
	if err != nil {
		return nil, err
	}
	if len(gaps) > 0 {
		return nil, evaluation.NonReportableRunError(displayPath, gaps[0])
	}
	artifact.Outputs = &Outputs{
		ReportMD:         result.Receipt.ReportMD,
		EvaluationOutput: result.Output,
		Rating:           &result.Receipt.RatingResult,
	}
	if err := store.Save(artifact); err != nil {
		return nil, err
	}
	return result.Receipt, nil
}
