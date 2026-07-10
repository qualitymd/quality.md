package evaluation

import (
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
)

// Run is a loaded evaluation run folder.
type Run struct {
	Path    string
	AbsPath string
	Model   *model.Spec
	Scale   []model.RatingLevel
}

// RunGap describes a reportability or validation gap in an evaluation run.
type RunGap struct {
	Kind   RunGapKind `json:"kind"`
	Ref    string     `json:"ref"`
	Detail string     `json:"detail"`
}

// RunGapKind identifies a class of evaluation run gap.
type RunGapKind string

const (
	GapMissingEvaluationData    RunGapKind = "missing-evaluation-data"
	GapMalformedEvaluationData  RunGapKind = "malformed-evaluation-data"
	GapUnreadableEvaluationData RunGapKind = "unreadable-evaluation-data"
	GapIncompleteEvaluationData RunGapKind = "incomplete-evaluation-data"
)

// DataStatus summarizes persisted evaluation data artifacts.
type DataStatus struct {
	Artifacts int `json:"artifacts"`
}

// RunStatus is the JSON contract emitted by evaluation status.
type RunStatus struct {
	SchemaVersion int    `json:"schemaVersion"`
	Path          string `json:"path"`
	Reportable    bool   `json:"reportable"`
	// Lifecycle is the runner lifecycle status of an artifact-backed run
	// (running, awaiting_evaluator, completed, failed, cancelled); empty for
	// manual multi-file runs.
	Lifecycle string `json:"lifecycle,omitempty"`
	// AwaitingEvaluator summarizes the pending harness work request when
	// Lifecycle is awaiting_evaluator: the run is resumable and incomplete,
	// and submitting the harness judgment is the pending action.
	AwaitingEvaluator *AwaitingEvaluatorCall `json:"awaitingEvaluator,omitempty"`
	Data              DataStatus             `json:"data"`
	Gaps              []RunGap               `json:"gaps"`
	NextActions       []receipt.Action       `json:"nextActions"`
}

// Load reads an evaluation run.
func Load(path string) (*Run, error) {
	return load(path)
}

// Inspect reads an evaluation run for history/status views.
func Inspect(path string) (*Run, error) {
	return load(path)
}

// InspectWithDisplay reads an evaluation run and uses displayPath in receipts.
func InspectWithDisplay(path, displayPath string) (*Run, error) {
	return loadWithDisplay(path, displayPath)
}

func load(path string) (*Run, error) {
	return loadWithDisplay(path, "")
}

func loadWithDisplay(path, displayPath string) (*Run, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	doc, err := document.Parse(filepath.Join(runAbs, ModelSnapshotFile))
	if err != nil {
		return nil, err
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, err
	}
	if displayPath == "" {
		displayPath = displayRunPath(runAbs)
	}
	return &Run{
		Path:    displayPath,
		AbsPath: filepath.ToSlash(runAbs),
		Model:   spec,
		Scale:   spec.RatingScale,
	}, nil
}

// Status summarizes whether the evaluation data graph is reportable.
func (r *Run) Status() RunStatus {
	if payloads, state, ok, err := runArtifactPayloads(r.AbsPath); ok {
		if err != nil {
			return RunStatus{
				SchemaVersion: SchemaVersion,
				Path:          r.Path,
				Gaps:          []RunGap{{Kind: GapUnreadableEvaluationData, Ref: RunArtifactFile, Detail: err.Error()}},
			}
		}
		return r.runArtifactStatus(payloads, state)
	}
	gaps := evaluationRenderableGaps(r.AbsPath)
	data := DataStatus{}
	if list, err := ListData(r.AbsPath, ""); err == nil {
		data.Artifacts = len(list.Artifacts)
	}
	status := RunStatus{
		SchemaVersion: SchemaVersion,
		Path:          r.Path,
		Reportable:    len(gaps) == 0,
		Data:          data,
		Gaps:          gaps,
	}
	if status.Reportable {
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-report-build",
			Label:   "Build evaluation report",
			Command: "qualitymd evaluation report build " + r.Path,
		}}
	} else {
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-data-set",
			Label:   "Persist required evaluation data",
			Command: "qualitymd evaluation data set " + r.Path + " < payloads.json",
		}}
	}
	return status
}
