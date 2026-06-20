// Package status assembles the read-only project-state snapshot emitted by
// `qualitymd status`.
package status

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/lint"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
)

const SchemaVersion = 1

type Readiness string

const (
	ReadinessMissingModel             Readiness = "missing-model"
	ReadinessInvalidModel             Readiness = "invalid-model"
	ReadinessReadyToEvaluate          Readiness = "ready-to-evaluate"
	ReadinessHasEvaluationHistory     Readiness = "has-evaluation-history"
	ReadinessNeedsEvaluationReconcile Readiness = "needs-evaluation-reconciliation"
)

// Options configures a status snapshot.
type Options struct {
	Path string
}

// ProjectStatusSnapshot is the JSON contract emitted by qualitymd status.
type ProjectStatusSnapshot struct {
	SchemaVersion int               `json:"schemaVersion"`
	Path          string            `json:"path"`
	Readiness     Readiness         `json:"readiness"`
	Model         ModelStatus       `json:"model"`
	Evaluations   EvaluationHistory `json:"evaluations"`
	NextActions   []receipt.Action  `json:"nextActions"`
}

type ModelStatus struct {
	Present        bool             `json:"present"`
	Valid          bool             `json:"valid"`
	Lint           LintStatus       `json:"lint"`
	Shape          *ModelShape      `json:"shape,omitempty"`
	SourceCoverage []SourceCoverage `json:"sourceCoverage,omitempty"`
}

type LintStatus struct {
	Summary  lint.Summary   `json:"summary"`
	Findings []lint.Finding `json:"findings,omitempty"`
}

type ModelShape struct {
	Targets      int `json:"targets"`
	Factors      int `json:"factors"`
	Requirements int `json:"requirements"`
	RatingLevels int `json:"ratingLevels"`
}

type SourceCoverage struct {
	TargetPath   []string `json:"targetPath"`
	Label        string   `json:"label"`
	SourceState  string   `json:"sourceState"`
	Source       string   `json:"source,omitempty"`
	Factors      int      `json:"factors"`
	Requirements int      `json:"requirements"`
	ChildTargets int      `json:"childTargets"`
}

type EvaluationHistory struct {
	Path    string                 `json:"path,omitempty"`
	Runs    int                    `json:"runs"`
	Latest  *EvaluationRunSummary  `json:"latest,omitempty"`
	Items   []EvaluationRunSummary `json:"items"`
	Summary EvaluationSummary      `json:"summary"`
}

type EvaluationSummary struct {
	Reportable            int `json:"reportable"`
	Incomplete            int `json:"incomplete"`
	Stale                 int `json:"stale"`
	ActiveRecommendations int `json:"activeRecommendations"`
	Problems              int `json:"problems"`
}

type EvaluationRunSummary struct {
	Path                  string                            `json:"path"`
	Reportable            bool                              `json:"reportable"`
	Stale                 bool                              `json:"stale"`
	Counts                evaluation.EvaluationRecordCounts `json:"counts"`
	Gaps                  int                               `json:"gaps"`
	ActiveRecommendations int                               `json:"activeRecommendations"`
	Problem               string                            `json:"problem,omitempty"`
}

// Snapshot assembles a deterministic project-state snapshot.
func Snapshot(opts Options) (*ProjectStatusSnapshot, error) {
	path := opts.Path
	if path == "" {
		path = "QUALITY.md"
	}
	result := &ProjectStatusSnapshot{
		SchemaVersion: SchemaVersion,
		Path:          path,
		Evaluations: EvaluationHistory{
			Items: []EvaluationRunSummary{},
		},
	}

	modelBytes, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			result.Readiness = ReadinessMissingModel
			result.Model = ModelStatus{Present: false}
			result.NextActions = []receipt.Action{{
				ID:      "init",
				Label:   "Create a starter QUALITY.md",
				Command: "qualitymd init " + path,
			}}
			return result, nil
		}
		return nil, err
	}

	lintResult, err := lint.Check(path)
	if err != nil {
		return nil, err
	}
	result.Model = ModelStatus{
		Present: true,
		Valid:   lintResult.Valid,
		Lint: LintStatus{
			Summary:  lintResult.Summary,
			Findings: lintResult.Findings,
		},
	}
	if !lintResult.Valid {
		result.Readiness = ReadinessInvalidModel
		result.NextActions = invalidModelActions(path, lintResult.Summary)
		return result, nil
	}

	spec, err := lint.Load(path)
	if err != nil {
		return nil, err
	}
	shape, coverage := modelSnapshot(spec)
	result.Model.Shape = &shape
	result.Model.SourceCoverage = coverage

	history, err := evaluationHistory(path, modelBytes)
	if err != nil {
		return nil, err
	}
	result.Evaluations = history
	result.Readiness = readiness(history)
	result.NextActions = nextActions(path, history, result.Readiness)
	return result, nil
}

func invalidModelActions(path string, summary lint.Summary) []receipt.Action {
	if summary.Fixable > 0 {
		return []receipt.Action{{
			ID:      "fix",
			Label:   "Apply deterministic lint repairs",
			Command: "qualitymd lint --fix " + path,
		}}
	}
	return []receipt.Action{{
		ID:      "lint",
		Label:   "Review lint findings",
		Command: "qualitymd lint " + path,
	}}
}

func modelSnapshot(spec *model.Spec) (ModelShape, []SourceCoverage) {
	acc := modelAccumulator{
		shape: ModelShape{Targets: 1, RatingLevels: len(spec.RatingScale)},
	}
	label := spec.Title
	if label == "" {
		label = "Model"
	}
	acc.coverage = append(acc.coverage, sourceCoverageRow(nil, label, spec.Source, "", len(spec.Factors), len(spec.Requirements), len(spec.Targets)))
	acc.walkFactors(spec.Factors)
	acc.shape.Requirements += len(spec.Requirements)
	acc.walkTargets(spec.Targets, nil, spec.Source)
	return acc.shape, acc.coverage
}

type modelAccumulator struct {
	shape    ModelShape
	coverage []SourceCoverage
}

func (a *modelAccumulator) walkTargets(targets map[string]model.Target, parentPath []string, inheritedSource string) {
	for _, name := range sortedKeys(targets) {
		target := targets[name]
		path := appendString(parentPath, name)
		label := target.Title
		if label == "" {
			label = name
		}
		a.shape.Targets++
		a.coverage = append(a.coverage, sourceCoverageRow(path, label, target.Source, inheritedSource, len(target.Factors), len(target.Requirements), len(target.Targets)))
		a.walkFactors(target.Factors)
		a.shape.Requirements += len(target.Requirements)
		nextSource := inheritedSource
		if target.Source != "" {
			nextSource = target.Source
		}
		a.walkTargets(target.Targets, path, nextSource)
	}
}

func (a *modelAccumulator) walkFactors(factors map[string]model.Factor) {
	for _, name := range sortedKeys(factors) {
		factor := factors[name]
		a.shape.Factors++
		a.shape.Requirements += len(factor.Requirements)
		a.walkFactors(factor.Factors)
	}
}

func sourceCoverageRow(path []string, label, declaredSource, inheritedSource string, factors, requirements, children int) SourceCoverage {
	row := SourceCoverage{
		TargetPath:   append([]string{}, path...),
		Label:        label,
		Factors:      factors,
		Requirements: requirements,
		ChildTargets: children,
	}
	switch {
	case declaredSource != "":
		row.SourceState = "declared"
		row.Source = declaredSource
	case inheritedSource != "":
		row.SourceState = "inherited"
		row.Source = inheritedSource
	default:
		row.SourceState = "missing"
	}
	return row
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func appendString(path []string, value string) []string {
	out := make([]string, 0, len(path)+1)
	out = append(out, path...)
	out = append(out, value)
	return out
}

func evaluationHistory(modelPath string, modelBytes []byte) (EvaluationHistory, error) {
	history := EvaluationHistory{Items: []EvaluationRunSummary{}}
	repoRoot, err := evaluation.FindRepoRoot(modelPath)
	if err != nil {
		return history, err
	}
	evalDirAbs, evalDirRel, err := evaluation.EvaluationDir(repoRoot, "")
	if err != nil {
		return history, err
	}
	history.Path = evalDirRel
	runs, err := evaluation.ListRunDirs(evalDirAbs, evalDirRel)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return history, nil
		}
		return history, err
	}
	for _, runDir := range runs {
		summary := inspectRun(runDir, modelBytes)
		history.Items = append(history.Items, summary)
		history.Runs++
		if summary.Problem != "" {
			history.Summary.Problems++
		}
		if summary.Reportable {
			history.Summary.Reportable++
		} else {
			history.Summary.Incomplete++
		}
		if summary.Stale {
			history.Summary.Stale++
		}
		history.Summary.ActiveRecommendations += summary.ActiveRecommendations
	}
	if len(history.Items) > 0 {
		latest := history.Items[len(history.Items)-1]
		history.Latest = &latest
	}
	return history, nil
}

func inspectRun(runDir evaluation.RunDir, modelBytes []byte) EvaluationRunSummary {
	summary := EvaluationRunSummary{Path: runDir.Rel}
	runModel, err := os.ReadFile(filepath.Join(runDir.Abs, "model.md"))
	if err != nil {
		summary.Problem = fmt.Sprintf("reading model.md: %v", err)
		return summary
	}
	summary.Stale = string(runModel) != string(modelBytes)

	run, err := evaluation.Load(runDir.Abs)
	if err != nil {
		summary.Problem = err.Error()
		return summary
	}
	runStatus := run.EvaluationRunStatus()
	summary.Reportable = runStatus.Reportable
	summary.Counts = run.RecordCounts()
	summary.Gaps = len(runStatus.Gaps)
	summary.ActiveRecommendations = run.ActiveRecommendationCount()
	return summary
}

func readiness(history EvaluationHistory) Readiness {
	if history.Runs == 0 {
		return ReadinessReadyToEvaluate
	}
	if history.Summary.Incomplete > 0 || history.Summary.Stale > 0 || history.Summary.ActiveRecommendations > 0 || history.Summary.Problems > 0 {
		return ReadinessNeedsEvaluationReconcile
	}
	return ReadinessHasEvaluationHistory
}

func nextActions(path string, history EvaluationHistory, state Readiness) []receipt.Action {
	switch state {
	case ReadinessReadyToEvaluate:
		return []receipt.Action{{
			ID:      "evaluation-create",
			Label:   "Create an evaluation run",
			Command: "qualitymd evaluation create --subject " + path,
		}}
	case ReadinessNeedsEvaluationReconcile:
		return reconciliationActions(path, history)
	case ReadinessHasEvaluationHistory:
		if history.Latest != nil && history.Latest.Reportable {
			return []receipt.Action{{
				ID:      "report-build",
				Label:   "Build the latest evaluation report",
				Command: "qualitymd evaluation report build " + history.Latest.Path,
			}}
		}
	}
	return []receipt.Action{}
}

func reconciliationActions(path string, history EvaluationHistory) []receipt.Action {
	if history.Latest == nil {
		return nil
	}
	latest := history.Latest
	switch {
	case latest.Problem != "" || !latest.Reportable:
		return []receipt.Action{{
			ID:      "evaluation-status-latest",
			Label:   "Inspect the latest evaluation run",
			Command: "qualitymd evaluation status " + latest.Path,
		}}
	case latest.Stale:
		return []receipt.Action{{
			ID:      "evaluation-create",
			Label:   "Create a fresh evaluation run",
			Command: "qualitymd evaluation create --subject " + path,
		}}
	case latest.ActiveRecommendations > 0:
		return []receipt.Action{{
			ID:      "report-build",
			Label:   "Build the latest evaluation report",
			Command: "qualitymd evaluation report build " + latest.Path,
		}}
	default:
		return nil
	}
}
