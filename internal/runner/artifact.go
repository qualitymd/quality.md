// Package runner is the CLI-owned deterministic evaluation runner. It owns
// the evaluation work graph, concurrency resolution, evaluator invocation,
// validation, the authoritative evaluation.json run artifact, run-local logs,
// and report generation. Evaluators (internal/evaluator) perform bounded
// judgment work units and nothing else.
package runner

import (
	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
)

const (
	// ArtifactSchemaVersion versions the evaluation.json run artifact.
	// Versions 1-3 belong to the historical multi-file data tree.
	ArtifactSchemaVersion = 5
	// ArtifactKind marks the evaluation.json document kind.
	ArtifactKind = "EvaluationRun"
	// ArtifactFile is the run-root artifact file name.
	ArtifactFile = evaluation.RunArtifactFile
)

// Run status values persisted in evaluation.json state.
const (
	StatusRunning   = "running"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
	// StatusAwaitingEvaluator marks a run checkpointed at a harness work
	// request: resumable, incomplete, and not failed — the pending action is
	// submitting the harness judgment result.
	StatusAwaitingEvaluator = "awaiting_evaluator"
)

// Work-unit status values persisted in evaluation.json state.
const (
	UnitPending   = "pending"
	UnitCompleted = "completed"
	UnitFailed    = "failed"
)

// Artifact is the authoritative structured run artifact at
// <run>/evaluation.json. It carries everything needed to validate, resume,
// review, and render the run; run-local logs stay in separate files.
type Artifact struct {
	SchemaVersion int      `json:"schemaVersion"`
	Kind          string   `json:"kind"`
	Manifest      Manifest `json:"manifest"`
	State         State    `json:"state"`
	Results       Results  `json:"results"`
	Outputs       *Outputs `json:"outputs,omitempty"`
}

// Manifest carries the run's immutable identity and setup.
type Manifest struct {
	EvaluationID   string                     `json:"evaluationId"`
	CreatedAt      string                     `json:"createdAt"`
	Model          string                     `json:"model"`
	RequestedScope evaluation.RunScope        `json:"requestedScope"`
	PlannedScope   evaluation.PlannedRunScope `json:"plannedScope"`
	Run            evaluation.RunMetadata     `json:"run"`
	// Evaluator is the selected evaluator or profile name.
	Evaluator string `json:"evaluator"`
	// EvaluatorKind is the selected evaluator's runtime kind.
	EvaluatorKind string `json:"evaluatorKind"`
	// Concurrency is the resolved concurrency cap for the run.
	Concurrency int `json:"concurrency"`
}

// ManifestPayload renders the manifest as the EvaluationManifest data payload
// the report renderer consumes.
func (m Manifest) ManifestPayload() map[string]any {
	requested := map[string]any{}
	if m.RequestedScope.AreaID != "" {
		requested["areaId"] = m.RequestedScope.AreaID
	}
	if len(m.RequestedScope.FactorFilter) > 0 {
		requested["factorFilter"] = anyStrings(m.RequestedScope.FactorFilter)
	}
	return map[string]any{
		"schemaVersion":  evaluation.SchemaVersion,
		"kind":           "EvaluationManifest",
		"evaluationId":   m.EvaluationID,
		"createdAt":      m.CreatedAt,
		"model":          m.Model,
		"requestedScope": requested,
		"plannedScope": map[string]any{
			"areaId":       m.PlannedScope.AreaID,
			"factorFilter": anyStrings(m.PlannedScope.FactorFilter),
		},
		"run": map[string]any{
			"number": m.Run.Number,
			"label":  m.Run.Label,
		},
	}
}

func anyStrings(values []string) []any {
	out := make([]any, 0, len(values))
	for _, value := range values {
		out = append(out, value)
	}
	return out
}

// State carries the run's execution lifecycle. Provider context identifiers
// and prompt-cache status live only in the evaluator-call log, never here.
type State struct {
	Status      string                `json:"status"`
	Failure     *Failure              `json:"failure,omitempty"`
	WorkUnits   map[string]*UnitState `json:"workUnits"`
	StartedAt   string                `json:"startedAt,omitempty"`
	UpdatedAt   string                `json:"updatedAt,omitempty"`
	CompletedAt string                `json:"completedAt,omitempty"`
	// Cancelled records a user interruption observed mid-run.
	Cancelled bool `json:"cancelled,omitempty"`
	// PendingEvaluatorCall is the persisted harness checkpoint: correlation
	// metadata for the one work request awaiting its result. It carries no
	// raw prompt, source, or result bodies; the request itself is rebuilt
	// deterministically on resume.
	PendingEvaluatorCall *PendingEvaluatorCall `json:"pendingEvaluatorCall,omitempty"`
	// HarnessIdentity is the harness runtime the run's judgment is bound to,
	// set by the first accepted harness result. Later results from a
	// different runtime are refused.
	HarnessIdentity *HarnessIdentity `json:"harnessIdentity,omitempty"`
}

// PendingEvaluatorCall is one awaiting harness work request's persisted
// correlation metadata.
type PendingEvaluatorCall struct {
	RequestID     string `json:"requestId"`
	WorkUnitID    string `json:"workUnitId"`
	InputHash     string `json:"inputHash"`
	CorrelationID string `json:"correlationId"`
	Attempt       int    `json:"attempt"`
}

// HarnessIdentity is the stable harness runtime attribution for a run.
type HarnessIdentity struct {
	Runtime string `json:"runtime"`
}

// UnitState is the persisted execution state of one work unit.
type UnitState struct {
	Status   string `json:"status"`
	Attempts int    `json:"attempts,omitempty"`
	// InputHash fingerprints the unit's resolved inputs so resume can detect
	// dependency-stale completed work.
	InputHash   string   `json:"inputHash,omitempty"`
	Failure     *Failure `json:"failure,omitempty"`
	StartedAt   string   `json:"startedAt,omitempty"`
	CompletedAt string   `json:"completedAt,omitempty"`
}

// Failure is a classified run or work-unit failure.
type Failure struct {
	Category evaluator.FailureCategory `json:"category"`
	Detail   string                    `json:"detail,omitempty"`
}

// Results carries the run's structured evaluation judgment: the accepted
// routine payloads in deterministic work-graph order.
type Results struct {
	Payloads []ResultPayload `json:"payloads"`
}

// ResultPayload is one accepted routine payload attributed to the work unit
// that produced it.
type ResultPayload struct {
	WorkUnit string         `json:"workUnit"`
	Payload  map[string]any `json:"payload"`
}

// Merge replaces every payload previously produced by workUnit with the given
// payloads, preserving first-insertion order for unchanged units.
func (r *Results) Merge(workUnit string, payloads []map[string]any) {
	kept := make([]ResultPayload, 0, len(r.Payloads)+len(payloads))
	inserted := false
	for _, existing := range r.Payloads {
		if existing.WorkUnit == workUnit {
			if !inserted {
				for _, payload := range payloads {
					kept = append(kept, ResultPayload{WorkUnit: workUnit, Payload: payload})
				}
				inserted = true
			}
			continue
		}
		kept = append(kept, existing)
	}
	if !inserted {
		for _, payload := range payloads {
			kept = append(kept, ResultPayload{WorkUnit: workUnit, Payload: payload})
		}
	}
	r.Payloads = kept
}

// ByWorkUnit returns the payloads the given work unit produced.
func (r *Results) ByWorkUnit(workUnit string) []map[string]any {
	var out []map[string]any
	for _, entry := range r.Payloads {
		if entry.WorkUnit == workUnit {
			out = append(out, entry.Payload)
		}
	}
	return out
}

// PayloadList returns every stored payload in artifact order.
func (r *Results) PayloadList() []map[string]any {
	out := make([]map[string]any, 0, len(r.Payloads))
	for _, entry := range r.Payloads {
		out = append(out, entry.Payload)
	}
	return out
}

// Outputs carries generated report references after report build.
type Outputs struct {
	ReportMD string `json:"reportMd"`
	// EvaluationOutput is the CLI-owned EvaluationOutputResult payload.
	EvaluationOutput map[string]any `json:"evaluationOutput"`
	// Rating is the scoped area rating carried into command receipts.
	Rating *evaluation.RatingResult `json:"rating,omitempty"`
}

// unit returns the state entry for id, creating it when absent.
func (s *State) unit(id string) *UnitState {
	if s.WorkUnits == nil {
		s.WorkUnits = map[string]*UnitState{}
	}
	if existing, ok := s.WorkUnits[id]; ok {
		return existing
	}
	created := &UnitState{Status: UnitPending}
	s.WorkUnits[id] = created
	return created
}
