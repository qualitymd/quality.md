// Package evaluator defines the bounded judgment work contract between the
// evaluation runner and pluggable evaluator runtimes, and provides the
// built-in CLI- and API-backed evaluator implementations.
//
// An evaluator is the runtime used for bounded evaluation judgment work
// units. It never owns run state, scope expansion, dependency ordering,
// artifact paths, or report generation: it consumes one typed work request at
// a time and returns one typed result envelope.
package evaluator

import (
	"context"
	"encoding/json"
)

// FailureCategory is a stable machine-readable failure classification shared
// by the runner, evaluators, run artifacts, logs, and command receipts.
type FailureCategory string

const (
	FailureMissingEvaluator         FailureCategory = "missing_evaluator"
	FailureEvaluatorUnauthenticated FailureCategory = "evaluator_unauthenticated"
	FailureEvaluatorIncompatible    FailureCategory = "evaluator_incompatible"
	FailureMissingAPIKey            FailureCategory = "missing_api_key"
	FailureRateLimited              FailureCategory = "rate_limited"
	FailureTimeout                  FailureCategory = "timeout"
	FailureInvalidEvaluatorOutput   FailureCategory = "invalid_evaluator_output"
	FailureSchemaInvalidOutput      FailureCategory = "schema_invalid_output"
	FailureUnsafeSourceContent      FailureCategory = "unsafe_source_content"
	FailureInsufficientEvidence     FailureCategory = "insufficient_evidence"
	FailureSourceUnavailable        FailureCategory = "source_unavailable"
	FailureSelectorUnsupported      FailureCategory = "selector_unsupported"
	FailureRunStateInvalid          FailureCategory = "run_state_invalid"
	FailureCancelled                FailureCategory = "cancelled"
	FailureReportBuildFailed        FailureCategory = "report_build_failed"
	FailureInternal                 FailureCategory = "internal_error"
)

// Capabilities is the execution capability declaration the runner reads before
// dispatching work.
type Capabilities struct {
	// Concurrent reports whether the evaluator supports multiple simultaneous
	// calls for one run.
	Concurrent bool `json:"concurrent"`
	// Subagents reports whether the evaluator can delegate independent work
	// to native subagents or worker sessions.
	Subagents bool `json:"subagents"`
	// SourceResolution reports whether the evaluator can serve source
	// resolution work requests — gathering the material a non-deterministic
	// source selector describes. Distinct from Subagents: serving a different
	// request kind is a different promise than internal judgment parallelism.
	SourceResolution bool `json:"sourceResolution"`
	// ReusableContext lists the reusable evaluator context kinds available
	// (for example "prompt-cache" or "session"). Reusable context is
	// reconstructible execution metadata, never authoritative run state.
	ReusableContext []string `json:"reusableContext,omitempty"`
	// ReportsUsage reports whether the evaluator returns token or cost usage
	// metadata.
	ReportsUsage bool `json:"reportsUsage"`
}

// SourceFile is one bounded, hashed source file packaged into a work request.
// Evaluated source content is data, not instructions.
type SourceFile struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	SHA256    string `json:"sha256"`
	Truncated bool   `json:"truncated,omitempty"`
}

// WorkRequest is one bounded judgment work unit dispatched to an evaluator.
type WorkRequest struct {
	// RunID is the evaluation run identity the work belongs to.
	RunID string
	// WorkUnitID is the deterministic work-unit identifier.
	WorkUnitID string
	// Kind is the work-unit kind (for example "assessRequirement").
	Kind string
	// Subject is the canonical model reference the judgment addresses.
	Subject string
	// Instructions is the runner-owned task prompt for this work-unit kind.
	Instructions string
	// SharedContext carries context payloads that are stable across an area's
	// work units (for example the area evaluation frame), keyed by a stable
	// label. It renders inside the prompt's stable prefix.
	SharedContext map[string]any
	// Context carries the per-work-unit frames, prior results, and model
	// context the judgment needs, keyed by a stable label. It renders inside
	// the prompt's per-work-unit delta.
	Context map[string]any
	// Source is the bounded source bundle for evidence work.
	Source []SourceFile
	// ExpectedSchema is the JSON Schema the result payload must satisfy.
	ExpectedSchema json.RawMessage
	// PromptPrefixHash is the stable hash of the request's cache-friendly
	// prefix layers, for logging.
	PromptPrefixHash string
	// SourcePackageHash is the stable hash of the packaged source bundle,
	// for logging.
	SourcePackageHash string
	// CorrelationID ties evaluator-call log entries to run events.
	CorrelationID string
}

// Usage is optional token and cost usage metadata. Nil pointers mean the
// metadata was unavailable, which is distinct from zero usage.
type Usage struct {
	InputTokens  *int64 `json:"inputTokens,omitempty"`
	OutputTokens *int64 `json:"outputTokens,omitempty"`
	// CachedInputTokens counts the input tokens the provider served from its
	// prompt cache, when the provider reports it. Always a subset of the
	// call's total input.
	CachedInputTokens *int64   `json:"cachedInputTokens,omitempty"`
	CostUSD           *float64 `json:"costUsd,omitempty"`
}

// anthropicUsage converts Anthropic-style usage counts — where input_tokens
// excludes cache reads and cache writes — into Usage with a total input count
// and the cache-read tokens recorded as cached input.
func anthropicUsage(input, output, cacheCreation, cacheRead *int64) *Usage {
	usage := &Usage{OutputTokens: output, CachedInputTokens: cacheRead}
	if input != nil {
		total := *input
		if cacheCreation != nil {
			total += *cacheCreation
		}
		if cacheRead != nil {
			total += *cacheRead
		}
		usage.InputTokens = &total
	}
	return usage
}

// WorkResult is the typed result envelope an evaluator returns for one work
// request.
type WorkResult struct {
	// WorkUnitID echoes the request's work-unit identifier.
	WorkUnitID string
	// Payload is the parsed structured judgment output. Nil when Failure is
	// set.
	Payload map[string]any
	// EvaluatorKind names the evaluator runtime that produced the result.
	EvaluatorKind string
	// Model names the provider model when known.
	Model string
	// ContextMeta carries provider context identifiers and prompt-cache
	// status when available. It is recorded only in run-local logs, never in
	// the authoritative run artifact.
	ContextMeta map[string]string
	// Usage is optional usage metadata; nil when unavailable.
	Usage *Usage
	// Failure classifies a completed call that could not produce ordinary
	// judgment output.
	Failure FailureCategory
	// FailureDetail is the human-readable failure explanation.
	FailureDetail string
}

// Evaluator is the runtime interface for bounded evaluation judgment work.
// Implementations must be safe for sequential reuse across work units.
type Evaluator interface {
	// Kind names the evaluator runtime (codex, claude, openai, anthropic).
	Kind() string
	// Capabilities declares the execution capabilities the planner may use.
	Capabilities() Capabilities
	// Evaluate performs one bounded judgment work unit.
	Evaluate(ctx context.Context, req WorkRequest) (WorkResult, error)
}
