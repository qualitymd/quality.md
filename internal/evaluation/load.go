package evaluation

import (
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
)

// Run is a loaded Evaluation run folder.
type Run struct {
	Path    string
	AbsPath string
	Model   *model.Spec
	Scale   []model.RatingLevel
}

// RunGap describes a reportability or validation gap in an Evaluation run.
type RunGap struct {
	Kind   RunGapKind `json:"kind"`
	Ref    string     `json:"ref"`
	Detail string     `json:"detail"`
}

// RunGapKind identifies a class of Evaluation run gap.
type RunGapKind string

const (
	GapMissingEvaluationData    RunGapKind = "missing-evaluation-data"
	GapMalformedEvaluationData  RunGapKind = "malformed-evaluation-data"
	GapUnreadableEvaluationData RunGapKind = "unreadable-evaluation-data"
	GapIncompleteEvaluationData RunGapKind = "incomplete-evaluation-data"
)

// DataStatus summarizes persisted Evaluation data artifacts.
type DataStatus struct {
	Artifacts int `json:"artifacts"`
}

// RunStatus is the JSON contract emitted by evaluation status.
type RunStatus struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Reportable    bool             `json:"reportable"`
	Data          DataStatus       `json:"data"`
	Gaps          []RunGap         `json:"gaps"`
	NextActions   []receipt.Action `json:"nextActions"`
}

// Load reads an Evaluation run.
func Load(path string) (*Run, error) {
	return load(path)
}

// Inspect reads an Evaluation run for history/status views.
func Inspect(path string) (*Run, error) {
	return load(path)
}

func load(path string) (*Run, error) {
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
	return &Run{
		Path:    displayRunPath(runAbs),
		AbsPath: filepath.ToSlash(runAbs),
		Model:   spec,
		Scale:   spec.RatingScale,
	}, nil
}

// Status summarizes whether the Evaluation data graph is reportable.
func (r *Run) Status() RunStatus {
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
			Label:   "Build Evaluation report",
			Command: "qualitymd evaluation report build " + r.Path,
		}}
	} else {
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-data-set",
			Label:   "Persist required Evaluation data",
			Command: "qualitymd evaluation data set " + r.Path + " < payload.json",
		}}
	}
	return status
}
