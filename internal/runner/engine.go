package runner

import (
	"context"
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
	"github.com/qualitymd/quality.md/internal/model"
)

// maxUnitAttempts bounds evaluator retries per work unit: the first attempt
// plus two retries for retryable failures.
const maxUnitAttempts = 3

// evaluatorCallTimeout bounds one evaluator call so a hung CLI or connection
// classifies as a timeout instead of stalling the run.
const evaluatorCallTimeout = 10 * time.Minute

// retryableFailures are the failure categories the runner's retry policy
// retries; every other category fails the run at once.
var retryableFailures = map[evaluator.FailureCategory]struct{}{
	evaluator.FailureRateLimited:            {},
	evaluator.FailureTimeout:                {},
	evaluator.FailureInvalidEvaluatorOutput: {},
	evaluator.FailureSchemaInvalidOutput:    {},
}

type engine struct {
	store         *Store
	artifact      *Artifact
	graph         *Graph
	spec          *model.Spec
	selection     *evaluator.Selection
	logs          *runLogs
	stderr        io.Writer
	progressMu    sync.Mutex
	artifactMu    sync.RWMutex
	workspaceRoot string
	runAbs        string
	displayPath   string
	now           func() time.Time
	sleep         func(ctx context.Context, d time.Duration)
	// sourceBundles memoizes each area's packaged source bundle for the run,
	// keyed by area reference. Packaging is deterministic, so the memoized
	// bundle is identical to a re-packaged one (same hash); the memo only
	// removes redundant filesystem work.
	sourceBundles map[string]*SourceBundle
	// harnessResult is the submitted --evaluator-result envelope, consumed by
	// the pending harness work unit.
	harnessResult *evaluator.HarnessResultEnvelope
	// awaitingRequest is the bounded work request staged for an
	// awaiting-evaluator receipt after a harness checkpoint.
	awaitingRequest *EvaluatorRequest
	// awaitingFailure classifies the rejected attempt an awaiting request is
	// retrying, for the receipt.
	awaitingFailure *Failure
}

// payloadFor returns the first payload a work unit produced, or nil.
func (e *engine) payloadFor(unitID string) map[string]any {
	payloads := e.payloadsByWorkUnit(unitID)
	if len(payloads) == 0 {
		return nil
	}
	return payloads[0]
}

func (e *engine) payloadsByWorkUnit(unitID string) []map[string]any {
	e.artifactMu.RLock()
	defer e.artifactMu.RUnlock()
	return e.artifact.Results.ByWorkUnit(unitID)
}

// requirementPayload returns the payload of the given kind that a
// requirement's combined judgment unit persisted, or nil.
func (e *engine) requirementPayload(reqRef string, kind evaluation.DataKind) map[string]any {
	for _, payload := range e.payloadsByWorkUnit(unitID(KindAssessRateRequirement, reqRef)) {
		if payload["kind"] == string(kind) {
			return payload
		}
	}
	return nil
}

func (e *engine) timestamp() string {
	return e.now().UTC().Format(time.RFC3339)
}

func (e *engine) progress(format string, args ...any) {
	if e.stderr == nil {
		return
	}
	e.progressMu.Lock()
	defer e.progressMu.Unlock()
	_, _ = fmt.Fprintf(e.stderr, format+"\n", args...)
}

func (e *engine) save() error {
	e.artifact.State.UpdatedAt = e.timestamp()
	return e.store.Save(e.artifact)
}

// execute runs the work graph in deterministic order. It returns the final
// run status; classified failures land in the artifact state rather than the
// error return.
func (e *engine) execute(ctx context.Context) (string, error) {
	if failure := e.unsupportedSelectorFailure(); failure != nil {
		return e.markRunFailed(failure)
	}
	if e.artifact.Manifest.Concurrency > 1 && !e.harnessBacked() {
		return e.executeConcurrent(ctx)
	}
	for _, unit := range e.graph.Units {
		if ctx.Err() != nil {
			return e.markCancelled()
		}
		var err error
		var failed bool
		switch {
		case unit.Kind == KindBuildReports:
			failed, err = e.runBuildReports(unit)
		case unit.EvaluatorBacked:
			failed, err = e.runEvaluatorUnit(ctx, unit)
		default:
			err = e.runDeterministicUnit(unit)
		}
		if err != nil {
			return StatusFailed, err
		}
		if e.awaitingRequest != nil {
			// The run checkpointed at a harness work request; the artifact
			// already persists the awaiting state.
			return StatusAwaitingEvaluator, nil
		}
		if ctx.Err() != nil {
			return e.markCancelled()
		}
		if failed {
			e.artifact.State.Status = StatusFailed
			e.artifact.State.CompletedAt = e.timestamp()
			if saveErr := e.save(); saveErr != nil {
				return StatusFailed, saveErr
			}
			return StatusFailed, nil
		}
	}
	e.artifact.State.Status = StatusCompleted
	e.artifact.State.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return StatusFailed, err
	}
	e.logs.event("run-completed", map[string]any{"run": e.artifact.Manifest.Run.Label})
	return StatusCompleted, nil
}

// unsupportedSelectorFailure runs the plan-time resolver check: if any
// in-scope area's pinned selector kind has no resolver under the selected
// evaluator, the run fails before any judgment is dispatched, naming the
// selector, its detected kind, and the remedy. Distinct from
// source_unavailable — the material is not missing; this run cannot resolve
// this kind of selector.
func (e *engine) unsupportedSelectorFailure() *Failure {
	if e.selection.Evaluator.Capabilities().SourceResolution {
		return nil
	}
	for _, area := range e.graph.Plan.Areas {
		record := e.artifact.Sources[area.Ref]
		if record == nil || record.Resolver != ResolverHarness {
			continue
		}
		return &Failure{
			Category: evaluator.FailureSelectorUnsupported,
			Detail:   selectorUnsupportedDetail(area.Ref, record.Selector, SourceKind(record.Kind), e.selection.Name),
		}
	}
	return nil
}

// markRunFailed records a run-level classified failure before any work unit
// is dispatched.
func (e *engine) markRunFailed(failure *Failure) (string, error) {
	e.artifact.State.Status = StatusFailed
	e.artifact.State.Failure = failure
	e.artifact.State.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return StatusFailed, err
	}
	e.logs.event("run-failed", map[string]any{
		"category": string(failure.Category),
		"detail":   failure.Detail,
	})
	e.progress("run failed: %s: %s", failure.Category, failure.Detail)
	return StatusFailed, nil
}

// markCancelled records a user interruption, leaving the artifact valid and
// resumable.
func (e *engine) markCancelled() (string, error) {
	e.artifact.State.Status = StatusCancelled
	e.artifact.State.Cancelled = true
	e.artifact.State.Failure = &Failure{Category: evaluator.FailureCancelled, Detail: "run was interrupted; resume with --resume"}
	if err := e.save(); err != nil {
		return StatusCancelled, err
	}
	e.logs.event("run-cancelled", nil)
	e.progress("Cancelled; accepted work is saved. Resume with --resume %s", e.displayPath)
	return StatusCancelled, nil
}

// runDeterministicUnit generates and persists a runner-owned payload.
func (e *engine) runDeterministicUnit(unit *Unit) error {
	payload, err := e.deterministicPayload(unit)
	if err != nil {
		return err
	}
	hash := hashJSON(payload)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == hash && e.payloadFor(unit.ID) != nil {
		return nil
	}
	if err := evaluation.ValidatePayload(unit.DataKind, payload, e.spec); err != nil {
		return fmt.Errorf("deterministic payload for %s: %w", unit.ID, err)
	}
	return e.acceptUnit(unit, state, []map[string]any{payload}, hash)
}

func (e *engine) deterministicPayload(unit *Unit) (map[string]any, error) {
	switch unit.Kind {
	case KindFrameEvaluation:
		return evaluationFramePayload(e.spec, e.artifact.Manifest), nil
	case KindFrameAreaEvaluation:
		return areaEvaluationFramePayload(e.graph.Plan.Area(unit.Subject)), nil
	case KindFrameRequirementEvaluation:
		return requirementEvaluationFramePayload(e.spec, e.graph.Plan.Requirement(unit.Subject)), nil
	case KindFrameFactorAnalysis:
		return factorAnalysisFramePayload(e.graph.Plan.Factor(unit.Subject)), nil
	case KindFrameAreaAnalysis:
		return areaAnalysisFramePayload(e.graph.Plan.Area(unit.Subject)), nil
	default:
		return nil, fmt.Errorf("no deterministic payload for work unit %s", unit.ID)
	}
}

// acceptUnit merges accepted payloads and persists the artifact atomically
// before the unit counts as complete for scheduling or resume.
func (e *engine) acceptUnit(unit *Unit, state *UnitState, payloads []map[string]any, inputHash string) error {
	e.artifactMu.Lock()
	e.artifact.Results.Merge(unit.ID, payloads)
	e.sortResultPayloads()
	e.artifactMu.Unlock()
	state.Status = UnitCompleted
	state.InputHash = inputHash
	state.Failure = nil
	state.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return err
	}
	e.logs.event("work-unit-completed", map[string]any{"workUnit": unit.ID})
	return nil
}

func (e *engine) sortResultPayloads() {
	order := map[string]int{}
	for i, unit := range e.graph.Units {
		order[unit.ID] = i
	}
	sort.SliceStable(e.artifact.Results.Payloads, func(i, j int) bool {
		return order[e.artifact.Results.Payloads[i].WorkUnit] < order[e.artifact.Results.Payloads[j].WorkUnit]
	})
}

// runEvaluatorUnit dispatches one bounded judgment work unit with the
// runner's retry policy. It returns failed=true when the unit (and so the
// run) fails with a classified failure.
func (e *engine) runEvaluatorUnit(ctx context.Context, unit *Unit) (bool, error) {
	if failure := e.resolveSourceGuard(unit); failure != nil {
		return e.failUnit(unit, failure)
	}
	if e.harnessBacked() {
		return e.runHarnessUnit(unit)
	}
	req, err := e.buildWorkRequest(unit)
	if err != nil {
		if failure := sourceUnavailableFailure(err); failure != nil {
			return e.failUnit(unit, failure)
		}
		return false, err
	}
	inputHash := workUnitInputHash(req)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == inputHash && e.unitResultPresent(unit) {
		e.logs.event("work-unit-reused", map[string]any{"workUnit": unit.ID})
		return false, nil
	}
	e.progress("%s...", unit.ID)
	e.logs.event("work-unit-started", map[string]any{"workUnit": unit.ID, "evaluator": e.selection.Name})

	lastFailure, done, err := e.attemptEvaluatorUnit(ctx, unit, state, req, inputHash)
	if err != nil || done {
		return false, err
	}
	state.Status = UnitFailed
	state.Failure = lastFailure
	e.artifact.State.Failure = lastFailure
	if err := e.save(); err != nil {
		return true, err
	}
	e.progress("%s failed: %s: %s", unit.ID, lastFailure.Category, lastFailure.Detail)
	return true, nil
}

// resolveSourceGuard is the per-unit backstop behind the plan-time
// unsupported-selector check: a resolution work unit must never be dispatched
// to an evaluator that cannot serve source resolution requests.
func (e *engine) resolveSourceGuard(unit *Unit) *Failure {
	if unit.Kind != KindResolveSource || e.selection.Evaluator.Capabilities().SourceResolution {
		return nil
	}
	record := e.artifact.Sources[unit.Subject]
	if record == nil {
		return &Failure{
			Category: evaluator.FailureInternal,
			Detail:   fmt.Sprintf("no pinned source record for %s", unit.Subject),
		}
	}
	return &Failure{
		Category: evaluator.FailureSelectorUnsupported,
		Detail:   selectorUnsupportedDetail(unit.Subject, record.Selector, SourceKind(record.Kind), e.selection.Name),
	}
}

// unitResultPresent reports whether a completed unit's persisted effect is
// present in the artifact: accepted payloads for judgment units, a captured
// source bundle for resolution units.
func (e *engine) unitResultPresent(unit *Unit) bool {
	if unit.Kind == KindResolveSource {
		record := e.artifact.Sources[unit.Subject]
		return record != nil && record.BundleHash != "" && len(record.Files) > 0
	}
	return len(e.payloadsByWorkUnit(unit.ID)) > 0
}

// failUnit records one work unit's terminal classified failure and fails the
// run without dispatching the unit.
func (e *engine) failUnit(unit *Unit, failure *Failure) (bool, error) {
	state := e.artifact.State.unit(unit.ID)
	state.Status = UnitFailed
	state.Failure = failure
	e.artifact.State.Failure = failure
	e.logs.event("work-unit-attempt-failed", map[string]any{
		"workUnit": unit.ID,
		"category": string(failure.Category),
		"detail":   failure.Detail,
	})
	if err := e.save(); err != nil {
		return true, err
	}
	e.progress("%s failed: %s: %s", unit.ID, failure.Category, failure.Detail)
	return true, nil
}

// workUnitInputHash fingerprints a work unit's resolved inputs so resume can
// detect dependency-stale completed work and stale pending harness requests.
func workUnitInputHash(req evaluator.WorkRequest) string {
	return hashJSON(map[string]any{
		"instructions":  req.Instructions,
		"sharedContext": req.SharedContext,
		"context":       req.Context,
		"schema":        string(req.ExpectedSchema),
		"source":        req.SourcePackageHash,
	})
}

// attemptEvaluatorUnit runs the retry loop for one work unit. done reports
// that the unit was accepted or the run was cancelled; otherwise the last
// classified failure is returned for the unit to fail with.
func (e *engine) attemptEvaluatorUnit(ctx context.Context, unit *Unit, state *UnitState, req evaluator.WorkRequest, inputHash string) (*Failure, bool, error) {
	var lastFailure *Failure
	for attempt := 1; attempt <= maxUnitAttempts; attempt++ {
		if ctx.Err() != nil {
			return nil, true, nil
		}
		state.Attempts++
		started := e.now()
		callCtx, cancelCall := context.WithTimeout(ctx, evaluatorCallTimeout)
		result, err := e.selection.Evaluator.Evaluate(callCtx, req)
		cancelCall()
		duration := e.now().Sub(started)
		if err != nil {
			return nil, false, fmt.Errorf("evaluator %s: %w", e.selection.Name, err)
		}

		failure := result.Failure
		detail := result.FailureDetail
		var payloads []map[string]any
		if failure == "" {
			payloads, failure, detail = e.acceptResultPayload(unit, result.Payload)
		}
		e.logCall(unit, req, result, attempt, duration, failure, detail)

		if failure == "" {
			if err := e.acceptUnit(unit, state, payloads, inputHash); err != nil {
				return nil, false, err
			}
			return nil, true, nil
		}
		lastFailure = &Failure{Category: failure, Detail: detail}
		e.logs.event("work-unit-attempt-failed", map[string]any{
			"workUnit": unit.ID,
			"attempt":  state.Attempts,
			"category": string(failure),
			"detail":   detail,
		})
		if failure == evaluator.FailureCancelled || ctx.Err() != nil {
			return nil, true, nil
		}
		if _, retryable := retryableFailures[failure]; !retryable || attempt == maxUnitAttempts {
			break
		}
		if failure == evaluator.FailureRateLimited {
			e.sleep(ctx, time.Duration(attempt)*2*time.Second)
		}
		e.progress("%s: retrying after %s", unit.ID, failure)
	}
	return lastFailure, false, nil
}

// acceptResultPayload normalizes and validates an evaluator payload for a
// work unit. A rejected payload is classified for the retry policy.
func (e *engine) acceptResultPayload(unit *Unit, payload map[string]any) ([]map[string]any, evaluator.FailureCategory, string) {
	if payload == nil {
		return nil, evaluator.FailureInvalidEvaluatorOutput, "evaluator returned no payload"
	}
	if unit.Kind == KindResolveSource {
		return e.acceptSourceResolution(unit, payload)
	}
	if unit.Kind == KindAssessRateRequirement {
		return e.acceptAssessRate(unit, payload)
	}
	if unit.Kind == KindRecommend {
		return e.acceptRecommendations(payload)
	}
	normalized, err := e.normalizePayload(unit, payload)
	if err != nil {
		return nil, evaluator.FailureInvalidEvaluatorOutput, err.Error()
	}
	if err := evaluation.ValidatePayload(unit.DataKind, normalized, e.spec); err != nil {
		return nil, evaluator.FailureSchemaInvalidOutput, err.Error()
	}
	if err := e.verifyAdviceCompleteness(unit, normalized); err != nil {
		return nil, evaluator.FailureSchemaInvalidOutput, err.Error()
	}
	return []map[string]any{normalized}, "", ""
}

// acceptSourceResolution validates a resolution result's returned material
// and captures it as the area's bounded, hashed source bundle in the
// artifact's sources record. The capture persists atomically with the unit's
// completion (the shared acceptUnit save), so every dependent judgment
// request is built from persisted data. Resolution units persist no result
// payloads: the captured bundle is the unit's effect of record.
func (e *engine) acceptSourceResolution(unit *Unit, payload map[string]any) ([]map[string]any, evaluator.FailureCategory, string) {
	bundle, err := captureResolvedSource(payload)
	if err != nil {
		return nil, evaluator.FailureInvalidEvaluatorOutput, err.Error()
	}
	record := e.artifact.Sources[unit.Subject]
	if record == nil {
		return nil, evaluator.FailureInternal, fmt.Sprintf("no pinned source record for %s", unit.Subject)
	}
	record.completeFromBundle(bundle, e.timestamp(), true)
	delete(e.sourceBundles, unit.Subject)
	return nil, "", ""
}

// acceptAssessRate splits the combined requirement judgment composite into
// the assessment and rating payloads and validates each against its own data
// kind. A composite missing either half is a retryable unit failure, so a
// partial requirement result is never persisted.
func (e *engine) acceptAssessRate(unit *Unit, payload map[string]any) ([]map[string]any, evaluator.FailureCategory, string) {
	parts := []struct {
		field string
		kind  evaluation.DataKind
	}{
		{"assessment", evaluation.DataKindRequirementAssessment},
		{"rating", evaluation.DataKindRequirementRating},
	}
	out := make([]map[string]any, 0, len(parts))
	for _, part := range parts {
		body, ok := payload[part.field].(map[string]any)
		if !ok {
			return nil, evaluator.FailureInvalidEvaluatorOutput,
				fmt.Sprintf("combined requirement judgment must carry an %s object", part.field)
		}
		normalized, err := normalizeSubjectPayload(body, part.kind, "requirementId", unit.Subject)
		if err != nil {
			return nil, evaluator.FailureInvalidEvaluatorOutput, err.Error()
		}
		if err := evaluation.ValidatePayload(part.kind, normalized, e.spec); err != nil {
			return nil, evaluator.FailureSchemaInvalidOutput, fmt.Sprintf("%s: %s", part.field, err)
		}
		out = append(out, normalized)
	}
	return out, "", ""
}

// normalizePayload enforces the runner-owned envelope fields: schema version,
// kind, and subject identity.
func (e *engine) normalizePayload(unit *Unit, payload map[string]any) (map[string]any, error) {
	subjectField := ""
	switch unit.Kind {
	case KindAnalyzeFactor:
		subjectField = "factorId"
	case KindAnalyzeArea:
		subjectField = "areaId"
	}
	return normalizeSubjectPayload(payload, unit.DataKind, subjectField, unit.Subject)
}

// normalizeSubjectPayload stamps the runner-owned envelope fields on one
// payload and pins its subject identity field to the work unit's subject.
func normalizeSubjectPayload(payload map[string]any, kind evaluation.DataKind, subjectField, subject string) (map[string]any, error) {
	payload["schemaVersion"] = evaluation.SchemaVersion
	payload["kind"] = string(kind)
	if subjectField != "" {
		if existing, ok := payload[subjectField].(string); ok && existing != "" && existing != subject {
			return nil, fmt.Errorf("payload %s = %q, want the work unit subject %q", subjectField, existing, subject)
		}
		payload[subjectField] = subject
	}
	return payload, nil
}

// acceptRecommendations unpacks the composite recommend result, assigns
// canonical recommendation IDs, and validates every recommendation.
func (e *engine) acceptRecommendations(payload map[string]any) ([]map[string]any, evaluator.FailureCategory, string) {
	items, ok := payload["recommendations"].([]any)
	if !ok || len(items) == 0 {
		return nil, evaluator.FailureInvalidEvaluatorOutput, "recommend result must carry a non-empty recommendations array"
	}
	used := map[string]struct{}{}
	payloads := make([]map[string]any, 0, len(items))
	for i, item := range items {
		rec, ok := item.(map[string]any)
		if !ok {
			return nil, evaluator.FailureInvalidEvaluatorOutput, fmt.Sprintf("recommendations[%d] must be an object", i)
		}
		rec["schemaVersion"] = evaluation.SchemaVersion
		rec["kind"] = string(evaluation.DataKindRecommendation)
		id, _ := rec["id"].(string)
		if !evaluation.ValidRecommendationID(id) {
			assigned, err := evaluation.NewRecommendationID(used)
			if err != nil {
				return nil, evaluator.FailureInternal, err.Error()
			}
			id = assigned
			rec["id"] = id
		}
		if _, dup := used[id]; dup {
			return nil, evaluator.FailureInvalidEvaluatorOutput, fmt.Sprintf("recommendations[%d] duplicates id %s", i, id)
		}
		used[id] = struct{}{}
		if err := evaluation.ValidatePayload(evaluation.DataKindRecommendation, rec, e.spec); err != nil {
			return nil, evaluator.FailureSchemaInvalidOutput, err.Error()
		}
		payloads = append(payloads, rec)
	}
	return payloads, "", ""
}

// verifyAdviceCompleteness checks advice payloads cover their full input set
// before acceptance, so a coverage miss retries instead of dead-ending at
// report build.
func (e *engine) verifyAdviceCompleteness(unit *Unit, payload map[string]any) error {
	switch unit.Kind {
	case KindRankFindings:
		return e.verifyFindingCoverage(payload, "orderedFindings")
	case KindRankRecommendations:
		if err := e.verifyRecommendationCoverage(payload); err != nil {
			return err
		}
		return e.verifyFindingCoverage(payload, "findingCoverage")
	default:
		return nil
	}
}

func (e *engine) verifyFindingCoverage(payload map[string]any, field string) error {
	covered := map[string]struct{}{}
	items, _ := payload[field].([]any)
	for _, item := range items {
		entry, _ := item.(map[string]any)
		if entry == nil {
			continue
		}
		covered[hashJSON(entry["findingRef"])] = struct{}{}
	}
	for _, finding := range e.findingIndex() {
		if _, ok := covered[hashJSON(finding.FindingRef)]; !ok {
			return fmt.Errorf("%s is missing finding %s of %s", field, finding.FindingID, finding.RequirementID)
		}
	}
	return nil
}

func (e *engine) verifyRecommendationCoverage(payload map[string]any) error {
	ranked := map[string]struct{}{}
	items, _ := payload["orderedRecommendations"].([]any)
	for _, item := range items {
		entry, _ := item.(map[string]any)
		if entry == nil {
			continue
		}
		if id, _ := entry["recommendationRef"].(string); id != "" {
			ranked[id] = struct{}{}
		}
	}
	for _, rec := range e.payloadsByWorkUnit(string(KindRecommend)) {
		id, _ := rec["id"].(string)
		if _, ok := ranked[id]; !ok {
			return fmt.Errorf("orderedRecommendations is missing recommendation %s", id)
		}
	}
	return nil
}

// runBuildReports renders the Markdown report tree from evaluation.json.
func (e *engine) runBuildReports(unit *Unit) (bool, error) {
	state := e.artifact.State.unit(unit.ID)
	e.artifactMu.RLock()
	payloads := append([]map[string]any{e.artifact.Manifest.ManifestPayload()}, e.artifact.Results.PayloadList()...)
	e.artifactMu.RUnlock()
	result, gaps, err := evaluation.BuildReportFromPayloads(e.runAbs, e.displayPath, payloads)
	if err != nil {
		return false, err
	}
	if len(gaps) > 0 {
		failure := &Failure{
			Category: evaluator.FailureReportBuildFailed,
			Detail:   fmt.Sprintf("%s %s: %s", gaps[0].Kind, gaps[0].Ref, gaps[0].Detail),
		}
		state.Status = UnitFailed
		state.Failure = failure
		e.artifact.State.Failure = failure
		if err := e.save(); err != nil {
			return true, err
		}
		return true, nil
	}
	e.artifact.Outputs = &Outputs{
		ReportMD:         result.Receipt.ReportMD,
		EvaluationOutput: result.Output,
		Rating:           &result.Receipt.RatingResult,
	}
	state.Status = UnitCompleted
	state.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return false, err
	}
	e.logs.event("reports-built", map[string]any{"reportMd": result.Receipt.ReportMD})
	return false, nil
}

func (e *engine) logCall(unit *Unit, req evaluator.WorkRequest, result evaluator.WorkResult, attempt int, duration time.Duration, failure evaluator.FailureCategory, detail string) {
	entry := map[string]any{
		"workUnit":          unit.ID,
		"kind":              string(unit.Kind),
		"evaluator":         e.selection.Name,
		"evaluatorKind":     result.EvaluatorKind,
		"attempt":           attempt,
		"durationMs":        duration.Milliseconds(),
		"concurrency":       e.artifact.Manifest.Concurrency,
		"promptPrefixHash":  req.PromptPrefixHash,
		"sourcePackageHash": req.SourcePackageHash,
		"status":            "accepted",
		"correlationId":     req.CorrelationID,
	}
	if result.Model != "" {
		entry["model"] = result.Model
	}
	if result.Payload != nil {
		entry["outputHash"] = hashJSON(result.Payload)
	}
	if failure != "" {
		entry["status"] = "failed"
		entry["category"] = string(failure)
		entry["detail"] = detail
	}
	if result.Usage != nil {
		entry["usage"] = usageLogFields(result.Usage)
	}
	if len(result.ContextMeta) > 0 {
		entry["contextMeta"] = result.ContextMeta
	}
	e.logs.call(entry)
}

// usageLogFields renders reported usage for the evaluator-call log, keeping
// unavailable counts absent rather than zero.
func usageLogFields(usage *evaluator.Usage) map[string]any {
	fields := map[string]any{}
	if usage.InputTokens != nil {
		fields["inputTokens"] = *usage.InputTokens
	}
	if usage.OutputTokens != nil {
		fields["outputTokens"] = *usage.OutputTokens
	}
	if usage.CachedInputTokens != nil {
		fields["cachedInputTokens"] = *usage.CachedInputTokens
	}
	if usage.CostUSD != nil {
		fields["costUsd"] = *usage.CostUSD
	}
	return fields
}
