package evaluation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/receipt"
)

// BuildReportReceipt is the JSON contract emitted after building Evaluation v2
// report files.
type BuildReportReceipt struct {
	SchemaVersion          int              `json:"schemaVersion"`
	Path                   string           `json:"path"`
	ReportMD               string           `json:"reportMd"`
	EvaluationOutputResult string           `json:"evaluationOutputResult"`
	RatingResult           RatingResult     `json:"ratingResult"`
	NextActions            []receipt.Action `json:"nextActions,omitempty"`
}

// BuildReport renders the Evaluation v2 report tree and output result for a
// run.
func BuildReport(path string) (*BuildReportReceipt, error) {
	return buildV2Report(path)
}

func nonReportableRunError(runPath string, gap RunGap) error {
	return fmt.Errorf("run is not reportable: %s %s: %s (run `qualitymd evaluation status %s` for all gaps)", gap.Kind, gap.Ref, gap.Detail, runPath)
}

func writeReportFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating report directory: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", filepath.ToSlash(path), err)
	}
	return nil
}
