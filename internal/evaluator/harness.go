package evaluator

import (
	"context"
	"errors"
)

// harnessEvaluator is the reserved checkpoint transport for judgment supplied
// by the invoking agent harness. It is selected explicitly (by the quality
// skill or the caller), never by auto discovery, and it is never invoked
// synchronously: the runner persists an awaiting-evaluator checkpoint and
// receives the result through a later resume instead of calling Evaluate.
type harnessEvaluator struct{}

var _ Evaluator = (*harnessEvaluator)(nil)

func (e *harnessEvaluator) Kind() string { return "harness" }

func (e *harnessEvaluator) Capabilities() Capabilities {
	return Capabilities{
		Strategies:      []Strategy{StrategySequential},
		ReusableContext: []string{"session"},
		ReportsUsage:    true,
	}
}

func (e *harnessEvaluator) Evaluate(context.Context, WorkRequest) (WorkResult, error) {
	return WorkResult{}, errors.New("the harness evaluator is checkpointed; the runner never calls Evaluate on it")
}

// HarnessIdentity names the harness runtime that produced a result (for
// example "claude-code" or "codex") and, when the harness reports one, the
// model used. The runtime is stable attribution; the model is optional
// per-call metadata, never a correctness dependency.
type HarnessIdentity struct {
	Runtime string `json:"runtime"`
	Model   string `json:"model,omitempty"`
}

// HarnessResultEnvelope is the transport envelope the invoking harness
// submits for one pending work request. It wraps the same judgment payload an
// in-process evaluator would return, plus the correlation the runner uses to
// bind the result to the request it emitted.
type HarnessResultEnvelope struct {
	// RequestID echoes the pending request identifier from the awaiting
	// receipt.
	RequestID string `json:"requestId"`
	// InputHash echoes the pending request's input hash.
	InputHash string `json:"inputHash"`
	// Evaluator identifies the harness runtime supplying the judgment.
	Evaluator HarnessIdentity `json:"evaluator"`
	// Payload is the structured judgment output. Empty when Failure is set.
	Payload map[string]any `json:"payload,omitempty"`
	// Failure classifies a harness attempt that could not produce ordinary
	// judgment output.
	Failure FailureCategory `json:"failure,omitempty"`
	// Detail is the human-readable failure explanation.
	Detail string `json:"detail,omitempty"`
	// Usage is optional usage metadata; nil when unavailable.
	Usage *Usage `json:"usage,omitempty"`
}
