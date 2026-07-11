// Package status assembles the read-only workspace status snapshot emitted by
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
	"github.com/qualitymd/quality.md/internal/workspace"
)

// SchemaVersion is the current status snapshot schema version.
const SchemaVersion = 2

// Readiness identifies the current model/evaluation lifecycle state.
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

// WorkspaceSnapshot is the JSON contract emitted by qualitymd status.
type WorkspaceSnapshot struct {
	SchemaVersion int               `json:"schemaVersion"`
	Path          string            `json:"path"`
	Workspace     *WorkspaceStatus  `json:"workspace,omitempty"`
	Readiness     Readiness         `json:"readiness"`
	Model         ModelStatus       `json:"model"`
	Evaluations   EvaluationHistory `json:"evaluations"`
	NextActions   []receipt.Action  `json:"nextActions"`
}

// WorkspaceStatus summarizes the resolved workspace paths for the selected
// model. All paths are repository-relative or workspace-relative, never
// absolute host paths.
type WorkspaceStatus struct {
	Root          string `json:"root"`
	Model         string `json:"model"`
	Config        string `json:"config"`
	ConfigPresent bool   `json:"configPresent"`
	DataDir       string `json:"dataDir"`
	EvaluationDir string `json:"evaluationDir"`
	ChangelogDir  string `json:"changelogDir"`
	LogDir        string `json:"logDir"`
}

// ModelStatus summarizes whether the model is present, valid, and evaluable.
type ModelStatus struct {
	Present        bool             `json:"present"`
	Valid          bool             `json:"valid"`
	Lint           LintStatus       `json:"lint"`
	Shape          *ModelShape      `json:"shape,omitempty"`
	SourceCoverage []SourceCoverage `json:"sourceCoverage,omitempty"`
}

// LintStatus embeds lint summary and findings in status output.
type LintStatus struct {
	Summary  lint.Summary   `json:"summary"`
	Findings []lint.Finding `json:"findings,omitempty"`
}

// ModelShape counts the model's structural elements.
type ModelShape struct {
	Areas        int `json:"areas"`
	Factors      int `json:"factors"`
	Requirements int `json:"requirements"`
	RatingLevels int `json:"ratingLevels"`
}

// SourceCoverage summarizes source declarations for one area. The source
// state comes from the shared model resolver (model.EffectiveSource).
type SourceCoverage struct {
	AreaPath     []string          `json:"areaPath"`
	Label        string            `json:"label"`
	SourceState  model.SourceState `json:"sourceState"`
	Source       string            `json:"source,omitempty"`
	Factors      int               `json:"factors"`
	Requirements int               `json:"requirements"`
	ChildAreas   int               `json:"childAreas"`
}

// EvaluationHistory summarizes discovered evaluation runs for the model.
type EvaluationHistory struct {
	Path    string                 `json:"path,omitempty"`
	Runs    int                    `json:"runs"`
	Latest  *EvaluationRunSummary  `json:"latest,omitempty"`
	Items   []EvaluationRunSummary `json:"items"`
	Summary EvaluationSummary      `json:"summary"`
}

// EvaluationSummary aggregates evaluation history counts.
type EvaluationSummary struct {
	Reportable int `json:"reportable"`
	Incomplete int `json:"incomplete"`
	// AwaitingEvaluator counts incomplete runs checkpointed at a pending
	// harness work request. They are included in Incomplete as well.
	AwaitingEvaluator int `json:"awaitingEvaluator,omitempty"`
	Stale             int `json:"stale"`
	Problems          int `json:"problems"`
}

// EvaluationRunSummary summarizes one evaluation run in status output.
type EvaluationRunSummary struct {
	Path       string `json:"path"`
	Reportable bool   `json:"reportable"`
	// Lifecycle is the runner lifecycle status of an artifact-backed run;
	// empty for manual multi-file runs. awaiting_evaluator means the run is
	// resumable with harness judgment as the pending action, not failed or
	// generically incomplete.
	Lifecycle     string `json:"lifecycle,omitempty"`
	Stale         bool   `json:"stale"`
	DataArtifacts int    `json:"dataArtifacts"`
	Gaps          int    `json:"gaps"`
	Problem       string `json:"problem,omitempty"`
}

// Snapshot assembles a deterministic workspace status snapshot.
func Snapshot(opts Options) (*WorkspaceSnapshot, error) {
	path := opts.Path
	if path == "" {
		path = "QUALITY.md"
	}
	result := &WorkspaceSnapshot{
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

	ws, err := workspace.Resolve(workspace.Options{
		Model: path,
	})
	if err != nil {
		return nil, err
	}
	result.Workspace = workspaceStatus(ws)

	history, err := evaluationHistory(ws, modelBytes)
	if err != nil {
		return nil, err
	}
	result.Evaluations = history
	result.Readiness = readiness(history)
	result.NextActions = nextActions(path, history, result.Readiness)
	return result, nil
}

func workspaceStatus(ws *workspace.Workspace) *WorkspaceStatus {
	return &WorkspaceStatus{
		Root:          ws.WorkspaceRoot.RepoRel,
		Model:         ws.Model.Rel,
		Config:        ws.Config.Rel,
		ConfigPresent: ws.ConfigPresent,
		DataDir:       ws.DataDir.Rel,
		EvaluationDir: ws.Evaluations.Rel,
		ChangelogDir:  ws.Log.Rel,
		LogDir:        ws.FeedbackLog.Rel,
	}
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
	return modelShape(spec), sourceCoverage(spec)
}

// modelShape counts the model's structural elements from the shared projection,
// so status and the model command read one model-tree walk.
func modelShape(spec *model.Spec) ModelShape {
	shape := ModelShape{RatingLevels: len(spec.RatingScale)}
	for _, element := range model.Flatten(model.Project(spec)) {
		switch element.Kind {
		case model.KindArea:
			shape.Areas++
		case model.KindFactor:
			shape.Factors++
		case model.KindRequirement:
			shape.Requirements++
		}
	}
	return shape
}

// sourceCoverage builds the per-area source provenance rows. Source state is a
// status-only concern the identity projection deliberately omits, so it keeps
// its own source-aware walk over the shared model resolver.
func sourceCoverage(spec *model.Spec) []SourceCoverage {
	label := spec.Title
	if label == "" {
		label = "Model"
	}
	rows := []SourceCoverage{sourceCoverageRow(spec, nil, label, len(spec.Factors), len(spec.Requirements), len(spec.Areas))}
	return appendAreaCoverage(rows, spec, spec.Areas, nil)
}

func appendAreaCoverage(rows []SourceCoverage, spec *model.Spec, areas map[string]model.Area, parentPath []string) []SourceCoverage {
	for _, name := range sortedKeys(areas) {
		area := areas[name]
		path := appendString(parentPath, name)
		label := area.Title
		if label == "" {
			label = name
		}
		rows = append(rows, sourceCoverageRow(spec, path, label, len(area.Factors), len(area.Requirements), len(area.Areas)))
		rows = appendAreaCoverage(rows, spec, area.Areas, path)
	}
	return rows
}

func sourceCoverageRow(spec *model.Spec, path []string, label string, factors, requirements, children int) SourceCoverage {
	row := SourceCoverage{
		AreaPath:     append([]string{}, path...),
		Label:        label,
		Factors:      factors,
		Requirements: requirements,
		ChildAreas:   children,
	}
	selector, state := model.EffectiveSource(spec, model.AreaPath(path))
	row.SourceState = state
	if state != model.SourceStateDefault {
		// The default selector is implicit (the document's directory); the row
		// reports it through the state alone.
		row.Source = selector
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

func evaluationHistory(ws *workspace.Workspace, modelBytes []byte) (EvaluationHistory, error) {
	history := EvaluationHistory{Items: []EvaluationRunSummary{}}
	evalDirAbs := ws.Evaluations.Abs
	evalDirRel := ws.Evaluations.Rel
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
			if summary.Lifecycle == evaluation.RunStatusAwaitingEvaluator {
				history.Summary.AwaitingEvaluator++
			}
		}
		if summary.Stale {
			history.Summary.Stale++
		}
	}
	if len(history.Items) > 0 {
		latest := history.Items[len(history.Items)-1]
		history.Latest = &latest
	}
	return history, nil
}

func inspectRun(runDir evaluation.RunDir, modelBytes []byte) EvaluationRunSummary {
	summary := EvaluationRunSummary{Path: runDir.Rel}
	runModel, err := os.ReadFile(filepath.Join(runDir.Abs, evaluation.ModelSnapshotFile))
	if err != nil {
		summary.Problem = fmt.Sprintf("reading %s: %v", evaluation.ModelSnapshotFile, err)
		return summary
	}
	summary.Stale = string(runModel) != string(modelBytes)

	run, err := evaluation.Inspect(runDir.Abs)
	if err != nil {
		summary.Problem = err.Error()
		return summary
	}
	runStatus := run.Status()
	summary.Reportable = runStatus.Reportable
	summary.Lifecycle = runStatus.Lifecycle
	summary.DataArtifacts = runStatus.Data.Artifacts
	summary.Gaps = len(runStatus.Gaps)
	return summary
}

func readiness(history EvaluationHistory) Readiness {
	if history.Runs == 0 {
		return ReadinessReadyToEvaluate
	}
	if history.Summary.Incomplete > 0 || history.Summary.Stale > 0 || history.Summary.Problems > 0 {
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
			Command: "qualitymd evaluation create --model " + path,
		}}
	case ReadinessNeedsEvaluationReconcile:
		return reconciliationActions(path, history)
	case ReadinessHasEvaluationHistory:
		if history.Latest != nil && history.Latest.Reportable {
			return []receipt.Action{{
				ID:      "report-build",
				Label:   "Build the latest evaluation report",
				Command: "qualitymd evaluation report build --model " + path + " " + history.Latest.Path,
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
	case latest.Problem == "" && latest.Lifecycle == evaluation.RunStatusAwaitingEvaluator:
		// Awaiting harness judgment is a normal checkpoint: the continuation
		// is resuming the run to recover the pending work request.
		return []receipt.Action{{
			ID:      "evaluation-run-reemit",
			Label:   "Resume the awaiting evaluation run to recover its pending work request",
			Command: "qualitymd evaluation run --model " + path + " --resume " + latest.Path + " --json",
		}}
	case latest.Problem != "" || !latest.Reportable:
		return []receipt.Action{{
			ID:      "evaluation-status-latest",
			Label:   "Inspect the latest evaluation run",
			Command: "qualitymd evaluation status --model " + path + " " + latest.Path,
		}}
	case latest.Stale:
		return []receipt.Action{{
			ID:      "evaluation-create",
			Label:   "Create a fresh evaluation run",
			Command: "qualitymd evaluation create --model " + path,
		}}
	default:
		return nil
	}
}
