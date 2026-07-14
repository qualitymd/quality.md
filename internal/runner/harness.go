package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qualitymd/quality.md/internal/evaluator"
)

// Harness dispatch: the runner owns the work graph and every accepted result;
// the invoking agent harness supplies judgment for bounded work requests. The
// runner keeps a rolling window of dependency-ready requests outstanding —
// bounded by the run's resolved concurrency — persisting an
// awaiting-evaluator checkpoint and returning the outstanding set in the
// receipt. A later resume submits one or more correlated result envelopes;
// each enters the same validation, retry, logging, and persistence paths as
// every other evaluator result, and the window tops up from the ready
// frontier as members are accepted.

// EvaluatorRequest is one complete bounded work request carried by an
// awaiting-evaluator receipt.
type EvaluatorRequest struct {
	RequestID      string                 `json:"requestId"`
	WorkUnitID     string                 `json:"workUnitId"`
	Kind           string                 `json:"kind"`
	Subject        string                 `json:"subject,omitempty"`
	Attempt        int                    `json:"attempt"`
	Instructions   string                 `json:"instructions"`
	SharedContext  map[string]any         `json:"sharedContext,omitempty"`
	Context        map[string]any         `json:"context,omitempty"`
	Source         []evaluator.SourceFile `json:"source,omitempty"`
	ExpectedSchema json.RawMessage        `json:"expectedSchema"`
	InputHash      string                 `json:"inputHash"`
	CorrelationID  string                 `json:"correlationId"`
	// LastFailure classifies the rejected attempt this request is retrying,
	// when the request re-emits after a failed attempt.
	LastFailure *Failure `json:"lastFailure,omitempty"`
}

// harnessRequestID derives the deterministic pending-request identifier, so a
// resume without a result re-emits the identical request identity.
func harnessRequestID(evaluationID, workUnitID, inputHash string, attempt int) string {
	sum := hashJSON(map[string]any{
		"evaluationId": evaluationID,
		"workUnit":     workUnitID,
		"inputHash":    inputHash,
		"attempt":      attempt,
	})
	return "req_" + sum[:16]
}

// loadHarnessResults reads and validates the --evaluator-result submission
// from a file or stdin: one result envelope or a JSON array of them, covering
// any subset of the outstanding requests. Structural problems are usage
// errors: nothing has been judged against the run yet.
func loadHarnessResults(path string, stdin io.Reader) ([]*evaluator.HarnessResultEnvelope, error) {
	var raw []byte
	var err error
	if path == "-" {
		if stdin == nil {
			stdin = os.Stdin
		}
		raw, err = io.ReadAll(stdin)
	} else {
		raw, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, &UsageError{Err: fmt.Errorf("reading --evaluator-result: %w", err)}
	}
	var envelopes []*evaluator.HarnessResultEnvelope
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	if trimmed := bytes.TrimSpace(raw); len(trimmed) > 0 && trimmed[0] == '[' {
		if err := decoder.Decode(&envelopes); err != nil {
			return nil, &UsageError{Err: fmt.Errorf("decoding --evaluator-result: %w", err)}
		}
	} else {
		var envelope evaluator.HarnessResultEnvelope
		if err := decoder.Decode(&envelope); err != nil {
			return nil, &UsageError{Err: fmt.Errorf("decoding --evaluator-result: %w", err)}
		}
		envelopes = append(envelopes, &envelope)
	}
	if len(envelopes) == 0 {
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result must carry at least one result envelope")}
	}
	for _, envelope := range envelopes {
		if err := validateHarnessEnvelope(envelope); err != nil {
			return nil, &UsageError{Err: err}
		}
	}
	return envelopes, nil
}

// validateHarnessEnvelope checks one submitted envelope's structural
// completeness before anything is judged against the run.
func validateHarnessEnvelope(envelope *evaluator.HarnessResultEnvelope) error {
	switch {
	case envelope == nil:
		return fmt.Errorf("--evaluator-result must not carry a null envelope")
	case envelope.RequestID == "":
		return fmt.Errorf("--evaluator-result envelopes must carry the outstanding requestId")
	case envelope.InputHash == "":
		return fmt.Errorf("--evaluator-result envelopes must carry the outstanding inputHash")
	case envelope.Evaluator.Runtime == "":
		return fmt.Errorf("--evaluator-result envelopes must identify the harness runtime in evaluator.runtime")
	case envelope.Failure == "" && envelope.Payload == nil:
		return fmt.Errorf("--evaluator-result envelopes must carry a payload or a classified failure")
	}
	return nil
}

// harnessBacked reports whether the run dispatches judgment to the invoking
// harness.
func (e *engine) harnessBacked() bool {
	return e.selection.Evaluator.Kind() == "harness"
}

// executeHarness runs a harness-backed invocation: apply any submitted result
// envelopes, top the outstanding window up from the dependency-ready
// frontier, and either finish the run or persist the awaiting checkpoint
// carrying the current outstanding set.
func (e *engine) executeHarness(ctx context.Context) (string, error) {
	failed, err := e.applyHarnessResults()
	if err != nil {
		return StatusFailed, err
	}
	if failed {
		return e.finishRunFailed()
	}
	return e.runHarnessGraph(ctx)
}

// finishRunFailed persists the terminal failed run status after a work unit
// failed terminally.
func (e *engine) finishRunFailed() (string, error) {
	e.artifact.State.Status = StatusFailed
	e.artifact.State.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return StatusFailed, err
	}
	return StatusFailed, nil
}

// pendingCall returns the outstanding pending call for a work unit, or nil.
func (s *State) pendingCall(workUnitID string) *PendingEvaluatorCall {
	for _, pending := range s.PendingEvaluatorCalls {
		if pending.WorkUnitID == workUnitID {
			return pending
		}
	}
	return nil
}

// clearPendingCall removes one outstanding pending call, preserving the
// emission order of the remainder.
func (s *State) clearPendingCall(requestID string) {
	kept := s.PendingEvaluatorCalls[:0]
	for _, pending := range s.PendingEvaluatorCalls {
		if pending.RequestID != requestID {
			kept = append(kept, pending)
		}
	}
	if len(kept) == 0 {
		s.PendingEvaluatorCalls = nil
		return
	}
	s.PendingEvaluatorCalls = kept
}

// applyHarnessResults binds each submitted envelope to an outstanding pending
// call by request identity and input hash and passes it through the shared
// acceptance path. Members are applied in the pending set's emission order,
// never submission order, so persisted state does not depend on how the
// harness ordered its reply. Accepted members clear their pending entry;
// failed members consume their unit's retry budget and stay outstanding for a
// retry attempt; an envelope matching no outstanding request is rejected
// after the valid members are applied, without altering any accepted result.
// failed=true means a member failed terminally and the run failed.
func (e *engine) applyHarnessResults() (bool, error) {
	envelopes := e.harnessResults
	e.harnessResults = nil
	if len(envelopes) == 0 {
		return false, nil
	}
	if len(e.artifact.State.PendingEvaluatorCalls) == 0 {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail:   "no work request is awaiting a harness result for this run",
		}
	}
	remaining := make([]*evaluator.HarnessResultEnvelope, len(envelopes))
	copy(remaining, envelopes)
	terminal := false
	outstanding := append([]*PendingEvaluatorCall(nil), e.artifact.State.PendingEvaluatorCalls...)
	for _, pending := range outstanding {
		matched := -1
		for i, envelope := range remaining {
			if envelope != nil && envelope.RequestID == pending.RequestID && envelope.InputHash == pending.InputHash {
				matched = i
				break
			}
		}
		if matched == -1 {
			continue
		}
		envelope := remaining[matched]
		remaining[matched] = nil
		failed, err := e.applyHarnessResult(pending, envelope)
		if err != nil {
			return false, err
		}
		if failed {
			terminal = true
		}
	}
	for _, envelope := range remaining {
		if envelope != nil {
			return terminal, &RunError{
				Category: evaluator.FailureRunStateInvalid,
				Detail: fmt.Sprintf("the submitted result %s does not correlate with an outstanding work request; "+
					"resume without --evaluator-result to recover the outstanding requests", envelope.RequestID),
			}
		}
	}
	return terminal, nil
}

// applyHarnessResult validates one submitted envelope against its pending
// request and passes its judgment through the shared acceptance path.
// failed=true means the member failed terminally (non-retryable failure or
// exhausted retry budget).
func (e *engine) applyHarnessResult(pending *PendingEvaluatorCall, envelope *evaluator.HarnessResultEnvelope) (bool, error) {
	unit := e.graph.Unit(pending.WorkUnitID)
	if unit == nil {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("the pending work request %s targets unknown work unit %s; start a new run",
				pending.RequestID, pending.WorkUnitID),
		}
	}
	if bound := e.artifact.State.HarnessIdentity; bound != nil && bound.Runtime != envelope.Evaluator.Runtime {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("run judgment is bound to harness runtime %q; a result from %q would mix harnesses — start a new run",
				bound.Runtime, envelope.Evaluator.Runtime),
		}
	}
	req, err := e.buildWorkRequest(unit)
	if err != nil {
		if failure := sourceUnavailableFailure(err); failure != nil {
			return e.failUnit(unit, failure)
		}
		return false, err
	}
	inputHash := workUnitInputHash(req)
	if inputHash != pending.InputHash {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("the model or source for %s changed after the work request was emitted; "+
				"a result must not be attached to different evidence — start a new run", unit.ID),
		}
	}

	state := e.artifact.State.unit(unit.ID)
	state.Attempts++
	result := evaluator.WorkResult{
		WorkUnitID:    unit.ID,
		EvaluatorKind: "harness",
		Model:         envelope.Evaluator.Model,
		ContextMeta:   map[string]string{"harnessRuntime": envelope.Evaluator.Runtime},
		Usage:         envelope.Usage,
		Payload:       envelope.Payload,
		Failure:       envelope.Failure,
		FailureDetail: envelope.Detail,
	}
	failure := result.Failure
	detail := result.FailureDetail
	var payloads []map[string]any
	if failure == "" {
		payloads, failure, detail = e.acceptResultPayload(unit, result.Payload)
	}
	e.logCall(unit, req, result, state.Attempts, 0, failure, detail)

	if failure == "" {
		if e.artifact.State.HarnessIdentity == nil {
			e.artifact.State.HarnessIdentity = &HarnessIdentity{Runtime: envelope.Evaluator.Runtime}
		}
		e.attributeResolutionRuntime(unit, envelope.Evaluator.Runtime)
		e.artifact.State.clearPendingCall(pending.RequestID)
		e.artifact.State.Status = StatusRunning
		if err := e.acceptUnit(unit, state, payloads, inputHash); err != nil {
			return false, err
		}
		return false, nil
	}

	return e.rejectHarnessAttempt(unit, state, pending, inputHash, &Failure{Category: failure, Detail: detail}), nil
}

// rejectHarnessAttempt records one member's failed attempt: a retryable
// failure within budget re-arms the pending call under the next attempt's
// request identity; otherwise the unit fails terminally. terminal=true means
// the run failed.
func (e *engine) rejectHarnessAttempt(unit *Unit, state *UnitState, pending *PendingEvaluatorCall, inputHash string, lastFailure *Failure) bool {
	e.logs.event("work-unit-attempt-failed", map[string]any{
		"workUnit": unit.ID,
		"attempt":  state.Attempts,
		"category": string(lastFailure.Category),
		"detail":   lastFailure.Detail,
	})
	if _, retryable := retryableFailures[lastFailure.Category]; retryable && state.Attempts < maxUnitAttempts {
		e.progress("%s: retrying after %s", unit.ID, lastFailure.Category)
		pending.Attempt = state.Attempts + 1
		pending.RequestID = harnessRequestID(e.artifact.Manifest.EvaluationID, unit.ID, inputHash, pending.Attempt)
		e.retryFailures[pending.RequestID] = lastFailure
		return false
	}
	e.artifact.State.clearPendingCall(pending.RequestID)
	state.Status = UnitFailed
	state.Failure = lastFailure
	e.artifact.State.Failure = lastFailure
	e.progress("%s failed: %s: %s", unit.ID, lastFailure.Category, lastFailure.Detail)
	return true
}

// runHarnessGraph advances the work graph for a harness-backed run: ready
// deterministic units run inline, and dependency-ready judgment work is
// emitted into the outstanding window up to the resolved concurrency, reusing
// the concurrent scheduler's frontier computation. The run finishes when the
// graph completes; otherwise the awaiting checkpoint persists atomically
// before the outstanding set is staged for the receipt.
func (e *engine) runHarnessGraph(ctx context.Context) (string, error) {
	window := e.artifact.Manifest.Concurrency
	scheduled := map[string]bool{}
	done := map[string]bool{}
	var requests []*EvaluatorRequest
	for {
		if ctx.Err() != nil {
			return e.markCancelled()
		}
		unit := nextReadyEvaluationStep(e.graph.Units, scheduled, done)
		if unit == nil {
			break
		}
		scheduled[unit.ID] = true
		request, completed, failed, err := e.advanceHarnessUnit(unit, window)
		if err != nil {
			return StatusFailed, err
		}
		if failed {
			return e.finishRunFailed()
		}
		if completed {
			done[unit.ID] = true
		}
		if request != nil {
			requests = append(requests, request)
		}
		// A cap-deferred judgment unit stays scheduled-but-not-done, so its
		// dependents stay blocked while later ready units are considered.
	}
	if len(e.artifact.State.PendingEvaluatorCalls) > 0 {
		return e.emitAwaitingWindow(requests, window)
	}
	if len(done) != len(e.graph.Units) {
		return StatusFailed, fmt.Errorf("evaluation work graph made no progress")
	}
	e.artifact.State.Status = StatusCompleted
	e.artifact.State.CompletedAt = e.timestamp()
	if err := e.save(); err != nil {
		return StatusFailed, err
	}
	e.logs.event("run-completed", map[string]any{"run": e.artifact.Manifest.Run.Label})
	return StatusCompleted, nil
}

// advanceHarnessUnit executes one dependency-ready unit for the harness loop:
// deterministic units run inline; judgment units are staged into the
// outstanding window. completed=true marks the unit done for the frontier;
// failed=true means the unit (and run) failed terminally.
func (e *engine) advanceHarnessUnit(unit *Unit, window int) (*EvaluatorRequest, bool, bool, error) {
	switch {
	case unit.Kind == KindBuildReports:
		failed, err := e.runBuildReports(unit)
		return nil, err == nil && !failed, failed, err
	case !unit.EvaluatorBacked:
		err := e.runDeterministicUnit(unit)
		return nil, err == nil, false, err
	default:
		return e.stageHarnessUnit(unit, window)
	}
}

// emitAwaitingWindow atomically persists the awaiting checkpoint, then logs
// and stages the outstanding request set for the receipt.
func (e *engine) emitAwaitingWindow(requests []*EvaluatorRequest, window int) (string, error) {
	e.artifact.State.Status = StatusAwaitingEvaluator
	e.artifact.State.Failure = nil
	if err := e.save(); err != nil {
		return StatusFailed, err
	}
	for _, request := range requests {
		e.logs.event("work-unit-awaiting-evaluator", map[string]any{
			"workUnit":  request.WorkUnitID,
			"requestId": request.RequestID,
			"attempt":   request.Attempt,
		})
		e.progress("%s awaiting harness judgment (request %s, attempt %d)", request.WorkUnitID, request.RequestID, request.Attempt)
	}
	e.progress("%d of %d permitted work requests outstanding", len(requests), window)
	e.awaitingRequests = requests
	return StatusAwaitingEvaluator, nil
}

// stageHarnessUnit handles one dependency-ready judgment work unit for the
// harness window: a completed unit with matching inputs is reused; an
// already-outstanding unit is re-staged with its persisted identity; a new
// unit is checkpointed while the window has capacity, and deferred otherwise.
// completed=true marks the unit done for the frontier; failed=true means the
// unit (and run) failed terminally before dispatch.
func (e *engine) stageHarnessUnit(unit *Unit, window int) (*EvaluatorRequest, bool, bool, error) {
	if failure := e.resolveSourceGuard(unit); failure != nil {
		failed, err := e.failUnit(unit, failure)
		return nil, false, failed, err
	}
	req, err := e.buildWorkRequest(unit)
	if err != nil {
		if failure := sourceUnavailableFailure(err); failure != nil {
			failed, failErr := e.failUnit(unit, failure)
			return nil, false, failed, failErr
		}
		return nil, false, false, err
	}
	inputHash := workUnitInputHash(req)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == inputHash && e.unitResultPresent(unit) {
		e.logs.event("work-unit-reused", map[string]any{"workUnit": unit.ID})
		return nil, true, false, nil
	}
	pending := e.artifact.State.pendingCall(unit.ID)
	if pending != nil && pending.InputHash != inputHash {
		return nil, false, false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("the model or source for %s changed after the work request was emitted; "+
				"a result must not be attached to different evidence — start a new run", unit.ID),
		}
	}
	if pending == nil {
		if len(e.artifact.State.PendingEvaluatorCalls) >= window {
			return nil, false, false, nil
		}
		pending = &PendingEvaluatorCall{
			RequestID:     harnessRequestID(e.artifact.Manifest.EvaluationID, unit.ID, inputHash, state.Attempts+1),
			WorkUnitID:    unit.ID,
			InputHash:     inputHash,
			CorrelationID: req.CorrelationID,
			Attempt:       state.Attempts + 1,
		}
		e.artifact.State.PendingEvaluatorCalls = append(e.artifact.State.PendingEvaluatorCalls, pending)
	}
	return &EvaluatorRequest{
		RequestID:      pending.RequestID,
		WorkUnitID:     unit.ID,
		Kind:           req.Kind,
		Subject:        req.Subject,
		Attempt:        pending.Attempt,
		Instructions:   req.Instructions,
		SharedContext:  req.SharedContext,
		Context:        req.Context,
		Source:         req.Source,
		ExpectedSchema: req.ExpectedSchema,
		InputHash:      inputHash,
		CorrelationID:  req.CorrelationID,
		LastFailure:    e.retryFailures[pending.RequestID],
	}, false, false, nil
}

// attributeResolutionRuntime records which harness runtime served an accepted
// resolution result in the area's source provenance record.
func (e *engine) attributeResolutionRuntime(unit *Unit, runtime string) {
	if unit.Kind != KindResolveSource {
		return
	}
	if record := e.artifact.Sources[unit.Subject]; record != nil {
		record.HarnessRuntime = runtime
	}
}
