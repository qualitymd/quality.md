package runner

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qualitymd/quality.md/internal/evaluator"
)

// Harness dispatch: the runner owns the work graph and every accepted result;
// the invoking agent harness supplies judgment one bounded request at a time.
// When a harness work unit becomes ready the engine persists an
// awaiting-evaluator checkpoint and returns the complete bounded request in
// the receipt; a later resume submits the correlated result envelope, which
// enters the same validation, retry, logging, and persistence paths as every
// other evaluator result.

// EvaluatorRequest is the complete bounded work request carried by an
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

// loadHarnessResult reads and validates the --evaluator-result envelope from
// a file or stdin. Structural problems are usage errors: nothing has been
// judged against the run yet.
func loadHarnessResult(path string, stdin io.Reader) (*evaluator.HarnessResultEnvelope, error) {
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
	var envelope evaluator.HarnessResultEnvelope
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	if err := decoder.Decode(&envelope); err != nil {
		return nil, &UsageError{Err: fmt.Errorf("decoding --evaluator-result: %w", err)}
	}
	switch {
	case envelope.RequestID == "":
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result must carry the pending requestId")}
	case envelope.InputHash == "":
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result must carry the pending inputHash")}
	case envelope.Evaluator.Runtime == "":
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result must identify the harness runtime in evaluator.runtime")}
	case envelope.Failure == "" && envelope.Payload == nil:
		return nil, &UsageError{Err: fmt.Errorf("--evaluator-result must carry a payload or a classified failure")}
	}
	return &envelope, nil
}

// harnessBacked reports whether the run dispatches judgment to the invoking
// harness.
func (e *engine) harnessBacked() bool {
	return e.selection.Evaluator.Kind() == "harness"
}

// runHarnessUnit executes one harness-backed judgment work unit: it either
// consumes the submitted result for the pending request or checkpoints the
// run and hands the bounded request back to the caller. failed=true means the
// unit (and run) failed terminally.
func (e *engine) runHarnessUnit(unit *Unit) (bool, error) {
	req, err := e.buildWorkRequest(unit)
	if err != nil {
		return false, err
	}
	inputHash := workUnitInputHash(req)
	state := e.artifact.State.unit(unit.ID)
	if state.Status == UnitCompleted && state.InputHash == inputHash && len(e.payloadsByWorkUnit(unit.ID)) > 0 {
		e.logs.event("work-unit-reused", map[string]any{"workUnit": unit.ID})
		return false, nil
	}
	pending := e.artifact.State.PendingEvaluatorCall
	if pending != nil && pending.WorkUnitID != unit.ID {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("the pending work request targets %s but %s is the next ready work unit; start a new run",
				pending.WorkUnitID, unit.ID),
		}
	}
	attempt := state.Attempts + 1
	if pending != nil {
		if pending.InputHash != inputHash {
			return false, &RunError{
				Category: evaluator.FailureRunStateInvalid,
				Detail: fmt.Sprintf("the model or source for %s changed after the work request was emitted; "+
					"a result must not be attached to different evidence — start a new run", unit.ID),
			}
		}
		attempt = pending.Attempt
	}
	if e.harnessResult == nil {
		return false, e.checkpointHarness(unit, req, inputHash, attempt, nil)
	}
	return e.consumeHarnessResult(unit, state, req, pending, inputHash)
}

// consumeHarnessResult validates the submitted envelope against the pending
// request and passes its judgment through the shared acceptance path.
func (e *engine) consumeHarnessResult(unit *Unit, state *UnitState, req evaluator.WorkRequest, pending *PendingEvaluatorCall, inputHash string) (bool, error) {
	envelope := e.harnessResult
	e.harnessResult = nil
	if pending == nil {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail:   "no work request is awaiting a harness result for this run",
		}
	}
	if envelope.RequestID != pending.RequestID || envelope.InputHash != pending.InputHash {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("the submitted result does not correlate with pending request %s; "+
				"resume without --evaluator-result to recover the pending request", pending.RequestID),
		}
	}
	if bound := e.artifact.State.HarnessIdentity; bound != nil && bound.Runtime != envelope.Evaluator.Runtime {
		return false, &RunError{
			Category: evaluator.FailureRunStateInvalid,
			Detail: fmt.Sprintf("run judgment is bound to harness runtime %q; a result from %q would mix harnesses — start a new run",
				bound.Runtime, envelope.Evaluator.Runtime),
		}
	}

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
		e.artifact.State.PendingEvaluatorCall = nil
		e.artifact.State.Status = StatusRunning
		if err := e.acceptUnit(unit, state, payloads, inputHash); err != nil {
			return false, err
		}
		return false, nil
	}

	lastFailure := &Failure{Category: failure, Detail: detail}
	e.logs.event("work-unit-attempt-failed", map[string]any{
		"workUnit": unit.ID,
		"attempt":  state.Attempts,
		"category": string(failure),
		"detail":   detail,
	})
	if _, retryable := retryableFailures[failure]; retryable && state.Attempts < maxUnitAttempts {
		e.progress("%s: retrying after %s", unit.ID, failure)
		return false, e.checkpointHarness(unit, req, inputHash, state.Attempts+1, lastFailure)
	}
	e.artifact.State.PendingEvaluatorCall = nil
	state.Status = UnitFailed
	state.Failure = lastFailure
	e.artifact.State.Failure = lastFailure
	if err := e.save(); err != nil {
		return true, err
	}
	e.progress("%s failed: %s: %s", unit.ID, lastFailure.Category, lastFailure.Detail)
	return true, nil
}

// checkpointHarness atomically persists the awaiting-evaluator checkpoint and
// stages the bounded request for the receipt. lastFailure classifies the
// attempt the request is retrying, if any.
func (e *engine) checkpointHarness(unit *Unit, req evaluator.WorkRequest, inputHash string, attempt int, lastFailure *Failure) error {
	pending := &PendingEvaluatorCall{
		RequestID:     harnessRequestID(e.artifact.Manifest.EvaluationID, unit.ID, inputHash, attempt),
		WorkUnitID:    unit.ID,
		InputHash:     inputHash,
		CorrelationID: req.CorrelationID,
		Attempt:       attempt,
	}
	e.artifact.State.PendingEvaluatorCall = pending
	e.artifact.State.Status = StatusAwaitingEvaluator
	e.artifact.State.Failure = nil
	if err := e.save(); err != nil {
		return err
	}
	e.logs.event("work-unit-awaiting-evaluator", map[string]any{
		"workUnit":  unit.ID,
		"requestId": pending.RequestID,
		"attempt":   attempt,
	})
	e.awaitingRequest = &EvaluatorRequest{
		RequestID:      pending.RequestID,
		WorkUnitID:     unit.ID,
		Kind:           req.Kind,
		Subject:        req.Subject,
		Attempt:        attempt,
		Instructions:   req.Instructions,
		SharedContext:  req.SharedContext,
		Context:        req.Context,
		Source:         req.Source,
		ExpectedSchema: req.ExpectedSchema,
		InputHash:      inputHash,
		CorrelationID:  req.CorrelationID,
	}
	e.awaitingFailure = lastFailure
	e.progress("%s awaiting harness judgment (request %s, attempt %d)", unit.ID, pending.RequestID, attempt)
	return nil
}
