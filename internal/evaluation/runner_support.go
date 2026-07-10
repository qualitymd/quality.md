package evaluation

// Runner support: exported entry points the deterministic evaluation runner
// (internal/runner) uses to prepare run folders, validate evaluator payloads,
// and render reports from in-memory results. Runner-created runs keep one
// authoritative evaluation.json at the run root and never write the
// historical multi-file data tree.

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/workspace"
)

// PreparedRun is a runner-owned run folder seeded with the model snapshot
// only. The runner persists all structured run state in evaluation.json.
type PreparedRun struct {
	RunAbs    string
	RunRel    string
	Manifest  EvaluationManifest
	ModelSpec *model.Spec
	Workspace *workspace.Workspace
}

// PrepareRun resolves scope, allocates the next numbered run folder, and
// writes the model snapshot. Unlike CreateRun it writes no data tree; the
// caller owns all further run persistence.
func PrepareRun(opts Options) (*PreparedRun, error) {
	ws, spec, manifest, err := resolveRunSetup(opts)
	if err != nil {
		return nil, err
	}
	manifest, err = completeRunIdentity(manifest)
	if err != nil {
		return nil, err
	}
	modelRaw, err := modelSnapshot(ws.Model.Abs, ws.Model.Rel)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(ws.Evaluations.Abs, 0o755); err != nil {
		return nil, fmt.Errorf("creating evaluation directory: %w", err)
	}
	number, name, err := nextRunName(ws.Evaluations.Abs, manifest.PlannedScope)
	if err != nil {
		return nil, err
	}
	manifest.Run.Number = number
	manifest.Run.Label = name
	manifest.Model = ws.Model.Rel
	runAbs := filepath.Join(ws.Evaluations.Abs, name)
	runRel := filepath.ToSlash(filepath.Join(ws.Evaluations.Rel, name))
	if err := os.Mkdir(runAbs, 0o755); err != nil {
		return nil, fmt.Errorf("creating run folder %s: %w", runRel, err)
	}
	if err := os.WriteFile(filepath.Join(runAbs, ModelSnapshotFile), modelRaw, 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", ModelSnapshotFile, err)
	}
	return &PreparedRun{RunAbs: runAbs, RunRel: runRel, Manifest: manifest, ModelSpec: spec, Workspace: ws}, nil
}

// RunPlan is the deterministic preview of the run a PrepareRun call with the
// same options would create. It carries no run identity (evaluation IDs are
// time-based) and writes nothing.
type RunPlan struct {
	RunRel    string
	Manifest  EvaluationManifest
	ModelSpec *model.Spec
	Workspace *workspace.Workspace
}

// PlanRun resolves scope and predicts the next run folder without writing
// anything.
func PlanRun(opts Options) (*RunPlan, error) {
	ws, spec, manifest, err := resolveRunSetup(opts)
	if err != nil {
		return nil, err
	}
	number := 1
	if _, statErr := os.Stat(ws.Evaluations.Abs); statErr == nil {
		number, err = nextRunNumber(ws.Evaluations.Abs)
		if err != nil {
			return nil, err
		}
	}
	name := fmt.Sprintf("%04d-%s-eval", number, runScopeSlug(manifest.PlannedScope))
	manifest.Run.Number = number
	manifest.Run.Label = name
	manifest.Model = ws.Model.Rel
	return &RunPlan{
		RunRel:    filepath.ToSlash(filepath.Join(ws.Evaluations.Rel, name)),
		Manifest:  manifest,
		ModelSpec: spec,
		Workspace: ws,
	}, nil
}

func resolveRunSetup(opts Options) (*workspace.Workspace, *model.Spec, EvaluationManifest, error) {
	modelPath := opts.Model
	if modelPath == "" {
		modelPath = workspace.DefaultModelPath
	}
	if filepath.Clean(modelPath) == "." {
		return nil, nil, EvaluationManifest{}, usagef("--model %q must name a QUALITY.md file, not a directory", modelPath)
	}
	ws, err := workspace.Resolve(workspace.Options{
		RepoRoot:              opts.RepoRoot,
		Model:                 opts.Model,
		EvaluationDirOverride: opts.ResolveDir,
	})
	if err != nil {
		return nil, nil, EvaluationManifest{}, err
	}
	spec, err := loadModel(ws.Model.Abs)
	if err != nil {
		return nil, nil, EvaluationManifest{}, err
	}
	manifest, err := resolveCreateScope(spec, opts)
	if err != nil {
		return nil, nil, EvaluationManifest{}, err
	}
	return ws, spec, manifest, nil
}

// LoadRunModel loads the frozen model snapshot for a run folder.
func LoadRunModel(runAbs string) (*model.Spec, error) {
	return loadRunModel(runAbs)
}

// ValidatePayload validates one evaluation data payload against its kind
// contract, the run model's references, and payload semantics.
func ValidatePayload(kind DataKind, payload map[string]any, spec *model.Spec) error {
	return validateDataPayloadForModel(kind, payload, spec)
}

// NewRecommendationID returns an unused recommendation ID in the canonical
// qrec_<token> form.
func NewRecommendationID(used map[string]struct{}) (string, error) {
	return nextRecommendationID(used)
}

// ValidRecommendationID reports whether id is a canonical recommendation ID.
func ValidRecommendationID(id string) bool {
	return validRecommendationID(id)
}

// FindingSelector returns the routine-ref selector addressing one finding
// inside a RequirementAssessmentResult payload.
func FindingSelector(findingID string) string {
	return "findings[" + findingID + "]"
}

// RunArtifactState summarizes the runner lifecycle state of an
// artifact-backed run for status surfaces.
type RunArtifactState struct {
	// Status is the run lifecycle status from evaluation.json state.
	Status string
	// AwaitingEvaluator summarizes the pending harness work request when
	// Status is awaiting_evaluator.
	AwaitingEvaluator *AwaitingEvaluatorCall
}

// AwaitingEvaluatorCall summarizes one pending harness work request.
type AwaitingEvaluatorCall struct {
	RequestID  string `json:"requestId"`
	WorkUnitID string `json:"workUnitId"`
	Attempt    int    `json:"attempt"`
}

// runArtifactPayloads loads the authoritative payload list of a
// runner-created run — the manifest payload followed by every accepted
// routine payload — plus its lifecycle state from its evaluation.json. ok is
// false when the run is not artifact-backed.
func runArtifactPayloads(runAbs string) ([]map[string]any, *RunArtifactState, bool, error) {
	raw, err := os.ReadFile(filepath.Join(runAbs, RunArtifactFile))
	if os.IsNotExist(err) {
		return nil, nil, false, nil
	}
	if err != nil {
		return nil, nil, true, fmt.Errorf("reading %s: %w", RunArtifactFile, err)
	}
	var doc struct {
		Manifest map[string]any `json:"manifest"`
		State    struct {
			Status               string `json:"status"`
			PendingEvaluatorCall *struct {
				RequestID  string `json:"requestId"`
				WorkUnitID string `json:"workUnitId"`
				Attempt    int    `json:"attempt"`
			} `json:"pendingEvaluatorCall"`
		} `json:"state"`
		Results struct {
			Payloads []struct {
				Payload map[string]any `json:"payload"`
			} `json:"payloads"`
		} `json:"results"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, nil, true, fmt.Errorf("decoding %s: %w", RunArtifactFile, err)
	}
	payloads := make([]map[string]any, 0, len(doc.Results.Payloads)+1)
	if doc.Manifest != nil {
		manifest := doc.Manifest
		manifest["schemaVersion"] = SchemaVersion
		manifest["kind"] = string(DataKindEvaluationManifest)
		payloads = append(payloads, manifest)
	}
	for _, entry := range doc.Results.Payloads {
		if entry.Payload != nil {
			payloads = append(payloads, entry.Payload)
		}
	}
	state := &RunArtifactState{Status: doc.State.Status}
	if pending := doc.State.PendingEvaluatorCall; pending != nil {
		state.AwaitingEvaluator = &AwaitingEvaluatorCall{
			RequestID:  pending.RequestID,
			WorkUnitID: pending.WorkUnitID,
			Attempt:    pending.Attempt,
		}
	}
	return payloads, state, true, nil
}

// RunStatusAwaitingEvaluator is the runner lifecycle status of a run
// checkpointed at a pending harness work request.
const RunStatusAwaitingEvaluator = "awaiting_evaluator"

// runArtifactStatus computes the RunStatus of an artifact-backed run.
func (r *Run) runArtifactStatus(payloads []map[string]any, state *RunArtifactState) RunStatus {
	status := RunStatus{
		SchemaVersion: SchemaVersion,
		Path:          r.Path,
		Data:          DataStatus{Artifacts: max(len(payloads)-1, 0)},
	}
	if state != nil {
		status.Lifecycle = state.Status
		if state.Status == RunStatusAwaitingEvaluator {
			status.AwaitingEvaluator = state.AwaitingEvaluator
		}
	}
	artifacts, err := collectArtifactsFromPayloads(r.AbsPath, payloads)
	if err != nil {
		status.Gaps = []RunGap{{Kind: GapMalformedEvaluationData, Ref: RunArtifactFile, Detail: err.Error()}}
	} else {
		status.Gaps, _ = payloadArtifactGaps(r.Model, artifacts)
	}
	if status.Gaps == nil {
		status.Gaps = []RunGap{}
	}
	status.Reportable = len(status.Gaps) == 0
	switch {
	case status.Lifecycle == RunStatusAwaitingEvaluator:
		// A checkpointed run is resumable and incomplete, not failed: the
		// pending action is recovering the work request and submitting the
		// harness judgment.
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-run-reemit",
			Label:   "Recover the pending harness work request",
			Command: "qualitymd evaluation run --resume " + r.Path + " --json",
		}, {
			ID:      "evaluation-evaluator-result",
			Label:   "Submit the harness judgment result",
			Command: "qualitymd evaluation run --resume " + r.Path + " --evaluator-result - --json",
		}}
	case status.Reportable:
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-report-build",
			Label:   "Build evaluation report",
			Command: "qualitymd evaluation report build " + r.Path,
		}}
	default:
		status.NextActions = []receipt.Action{{
			ID:      "evaluation-run-resume",
			Label:   "Resume the evaluation run",
			Command: "qualitymd evaluation run --resume " + r.Path,
		}}
	}
	return status
}

// NonReportableRunError formats the standard non-reportable run error for a
// gap.
func NonReportableRunError(runPath string, gap RunGap) error {
	return nonReportableRunError(runPath, gap)
}

// PayloadReportResult is the outcome of a payload-driven report build.
type PayloadReportResult struct {
	Receipt *BuildReportReceipt
	// Output is the CLI-owned EvaluationOutputResult payload describing the
	// generated report tree. Runner runs persist it inside evaluation.json.
	Output map[string]any
}

// PayloadGaps reports the reportability gaps for an in-memory payload set,
// using the same plan and coverage rules as multi-file runs.
func PayloadGaps(runAbs string, payloads []map[string]any) ([]RunGap, error) {
	spec, err := loadRunModel(runAbs)
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", ModelSnapshotFile, err)
	}
	artifacts, err := collectArtifactsFromPayloads(runAbs, payloads)
	if err != nil {
		return nil, err
	}
	gaps, _ := payloadArtifactGaps(spec, artifacts)
	return gaps, nil
}

// BuildReportFromPayloads renders the deterministic Markdown report tree for
// a runner-created run from in-memory payloads. The payload set must include
// the EvaluationManifest and EvaluationFrame payloads. Reports are written
// under runAbs; the EvaluationOutputResult payload is returned, not written.
func BuildReportFromPayloads(runAbs, displayPath string, payloads []map[string]any) (*PayloadReportResult, []RunGap, error) {
	spec, err := loadRunModel(runAbs)
	if err != nil {
		return nil, nil, fmt.Errorf("loading %s: %w", ModelSnapshotFile, err)
	}
	if displayPath == "" {
		displayPath = displayRunPath(runAbs)
	}
	artifacts, err := collectArtifactsFromPayloads(runAbs, payloads)
	if err != nil {
		return nil, nil, err
	}
	gaps, plan := payloadArtifactGaps(spec, artifacts)
	if len(gaps) > 0 {
		return nil, gaps, nil
	}
	runRel := workspaceRelativeRunPath(runAbs, artifacts.Manifest)
	reports := renderEvaluationReportTree(spec, artifacts, plan, runRel)
	for _, report := range reports {
		reportAbs := filepath.Join(runAbs, report.Path)
		if err := os.MkdirAll(filepath.Dir(reportAbs), 0o755); err != nil {
			return nil, nil, fmt.Errorf("creating report directory: %w", err)
		}
		if err := writeReportFile(reportAbs, []byte(report.Content)); err != nil {
			return nil, nil, err
		}
	}
	output := evaluationOutputResult(artifacts, plan, reports)
	receipt := &BuildReportReceipt{
		SchemaVersion: SchemaVersion,
		Path:          displayPath,
		ReportMD:      filepath.ToSlash(filepath.Join(displayPath, "report.md")),
		// Runner runs persist the output index inside evaluation.json rather
		// than a data-tree file.
		EvaluationOutputResult: filepath.ToSlash(filepath.Join(displayPath, "evaluation.json")),
		RatingResult:           evaluationReceiptRating(plan.ScopedAreaAnalysis),
	}
	return &PayloadReportResult{Receipt: receipt, Output: output}, nil, nil
}

func payloadArtifactGaps(spec *model.Spec, artifacts *evaluationArtifacts) ([]RunGap, *evaluationReportPlan) {
	plan, gap := resolveEvaluationReportPlan(artifacts)
	if gap != nil {
		return []RunGap{*gap}, nil
	}
	if gaps := plannedCoverageGaps(spec, artifacts, plan); len(gaps) > 0 {
		return gaps, nil
	}
	if gaps := adviceCoverageGaps(artifacts); len(gaps) > 0 {
		return gaps, nil
	}
	return nil, plan
}

func collectArtifactsFromPayloads(runAbs string, payloads []map[string]any) (*evaluationArtifacts, error) {
	out := &evaluationArtifacts{
		RunLabel:        filepath.Base(runAbs),
		Recommendations: map[string]map[string]any{},
		Areas:           map[string]*evaluationAreaArtifacts{},
		Factors:         map[string]*evaluationFactorArtifacts{},
		Requirements:    map[string]*evaluationRequirementArtifacts{},
	}
	for _, payload := range payloads {
		kind, err := payloadKind(payload)
		if err != nil {
			return nil, err
		}
		collector := evaluationPayloadCollectors[kind]
		if collector == nil {
			continue
		}
		if err := collector(out, payload); err != nil {
			return nil, err
		}
	}
	return out, nil
}
