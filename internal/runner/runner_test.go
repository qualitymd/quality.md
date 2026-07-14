package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
)

const testModel = `---
title: Test model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: Minimum
    description: Minimum.
    criterion: Barely meets it.
factors:
  reliability:
    title: Reliability
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
    factors: [reliability]
  has-docs:
    title: Has docs
    assessment: Inspect docs.
source: src
---
`

func testRunnerRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(testModel), 0o644); err != nil {
		t.Fatalf("WriteFile(QUALITY.md) error = %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo, "src"), 0o755); err != nil {
		t.Fatalf("MkdirAll(src) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "src", "main.txt"), []byte("package main\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(src/main.txt) error = %v", err)
	}
	return repo
}

func writeRunnerConfig(t *testing.T, repo, config string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(repo, ".quality"), 0o755); err != nil {
		t.Fatalf("MkdirAll(.quality) error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, ".quality", "config.yaml"), []byte(config), 0o644); err != nil {
		t.Fatalf("WriteFile(config.yaml) error = %v", err)
	}
}

// scriptedEvaluator produces schema-valid payloads for every judgment work
// unit, tracking calls per work unit. failUnits maps work-unit IDs to a count
// of leading attempts that return a schema-invalid payload.
type scriptedEvaluator struct {
	mu        sync.Mutex
	calls     map[string]int
	failUnits map[string]int
}

var _ evaluator.Evaluator = (*scriptedEvaluator)(nil)

func newScriptedEvaluator() *scriptedEvaluator {
	return &scriptedEvaluator{calls: map[string]int{}, failUnits: map[string]int{}}
}

func (s *scriptedEvaluator) Kind() string { return "scripted" }

func (s *scriptedEvaluator) Capabilities() evaluator.Capabilities {
	return evaluator.Capabilities{Concurrent: true}
}

func (s *scriptedEvaluator) Evaluate(_ context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	s.mu.Lock()
	s.calls[req.WorkUnitID]++
	remaining := s.failUnits[req.WorkUnitID]
	if remaining > 0 {
		s.failUnits[req.WorkUnitID] = remaining - 1
	}
	s.mu.Unlock()
	result := evaluator.WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: "scripted"}
	if remaining > 0 {
		result.Payload = map[string]any{"broken": true}
		return result, nil
	}
	payload, err := s.payloadFor(req)
	if err != nil {
		return result, err
	}
	result.Payload = roundTripJSON(payload)
	return result, nil
}

//nolint:cyclop // The scripted payload registry mirrors the work-unit kinds.
func (s *scriptedEvaluator) payloadFor(req evaluator.WorkRequest) (map[string]any, error) {
	analysis := map[string]any{"status": "analyzed", "ratingLevelId": "rating:target", "rationale": "Synthesized."}
	switch UnitKind(req.Kind) {
	case KindResolveSource:
		return map[string]any{
			"files": []any{map[string]any{
				"path":    "tickets/T-1",
				"content": "T-1: the API times out under load.",
			}},
		}, nil
	case KindAssessRateRequirement:
		return map[string]any{
			"assessment": map[string]any{
				"requirementId": req.Subject,
				"status":        "assessed",
				"confidence":    "medium",
				"findings": []any{map[string]any{
					"id":         "gap-001",
					"type":       "gap",
					"severity":   "medium",
					"confidence": "medium",
					"statement":  "A gap exists.",
					"condition":  "The gap condition holds.",
					"criteria": []any{map[string]any{
						"requirementId": req.Subject,
						"ratingLevelId": "rating:target",
						"criterion":     "Meets it.",
					}},
					"basis":    map[string]any{"status": "verified", "statement": "Directly observed."},
					"effect":   map[string]any{"statement": "Holds below target."},
					"evidence": []any{map[string]any{"sourceRef": "src/main.txt", "statement": "Observed in source."}},
				}},
			},
			"rating": map[string]any{
				"requirementId": req.Subject,
				"status":        "rated",
				"ratingLevelId": "rating:minimum",
				"rationale":     "The gap holds the rating at minimum.",
			},
		}, nil
	case KindAnalyzeFactor:
		return map[string]any{
			"factorId":                   req.Subject,
			"localAnalysis":              analysis,
			"localAndDescendantAnalysis": analysis,
		}, nil
	case KindAnalyzeArea:
		return map[string]any{
			"areaId":                     req.Subject,
			"localAnalysis":              analysis,
			"localAndDescendantAnalysis": analysis,
		}, nil
	case KindRankFindings:
		return s.rankFindingsPayload(req)
	case KindRecommend:
		return map[string]any{
			"recommendations": []any{map[string]any{
				"title":         "Close the gap",
				"description":   "Close the observed gap.",
				"background":    "One gap was found.",
				"expectedValue": "Higher rating.",
				"doneCriterion": "The gap is closed.",
				"impact":        "medium",
				"confidence":    "medium",
				"traceRefs": []any{map[string]any{
					"kind":    "RequirementAssessmentResult",
					"subject": map[string]any{"requirementId": "requirement:root::has-tests"},
				}},
			}},
		}, nil
	case KindRankRecommendations:
		return s.rankRecommendationsPayload(req)
	default:
		return nil, fmt.Errorf("unexpected work-unit kind %s", req.Kind)
	}
}

type sequentialOnlyEvaluator struct {
	inner *scriptedEvaluator
}

func (s *sequentialOnlyEvaluator) Kind() string { return s.inner.Kind() }

func (s *sequentialOnlyEvaluator) Capabilities() evaluator.Capabilities {
	return evaluator.Capabilities{Concurrent: false}
}

func (s *sequentialOnlyEvaluator) Evaluate(ctx context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	return s.inner.Evaluate(ctx, req)
}

type overlappingEvaluator struct {
	inner  *scriptedEvaluator
	mu     sync.Mutex
	active int
	max    int
	gate   chan struct{}
	once   sync.Once
}

func newOverlappingEvaluator() *overlappingEvaluator {
	return &overlappingEvaluator{inner: newScriptedEvaluator(), gate: make(chan struct{})}
}

func (e *overlappingEvaluator) Kind() string { return e.inner.Kind() }

func (e *overlappingEvaluator) Capabilities() evaluator.Capabilities {
	return evaluator.Capabilities{Concurrent: true}
}

func (e *overlappingEvaluator) Evaluate(ctx context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	if UnitKind(req.Kind) != KindAssessRateRequirement {
		return e.inner.Evaluate(ctx, req)
	}
	e.mu.Lock()
	e.active++
	if e.active > e.max {
		e.max = e.active
	}
	if e.active >= 2 {
		e.once.Do(func() { close(e.gate) })
	}
	e.mu.Unlock()
	defer func() {
		e.mu.Lock()
		e.active--
		e.mu.Unlock()
	}()
	select {
	case <-e.gate:
	case <-ctx.Done():
		return evaluator.WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: "scripted", Failure: evaluator.FailureCancelled, FailureDetail: "cancelled"}, nil
	case <-time.After(250 * time.Millisecond):
	}
	return e.inner.Evaluate(ctx, req)
}

func (e *overlappingEvaluator) maxActive() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.max
}

func (s *scriptedEvaluator) rankFindingsPayload(req evaluator.WorkRequest) (map[string]any, error) {
	findings, err := contextFindings(req)
	if err != nil {
		return nil, err
	}
	ordered := make([]any, 0, len(findings))
	for i, finding := range findings {
		ordered = append(ordered, map[string]any{
			"rank":       i + 1,
			"findingRef": finding["findingRef"],
			"tier":       "P1",
			"rationale":  "Ranked by severity.",
		})
	}
	return map[string]any{"orderedFindings": ordered}, nil
}

func (s *scriptedEvaluator) rankRecommendationsPayload(req evaluator.WorkRequest) (map[string]any, error) {
	findings, err := contextFindings(req)
	if err != nil {
		return nil, err
	}
	recs, ok := req.Context["recommendations"].([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("recommendations context missing")
	}
	orderedRecs := make([]any, 0, len(recs))
	firstID := ""
	for i, rec := range recs {
		id, _ := rec["id"].(string)
		if firstID == "" {
			firstID = id
		}
		orderedRecs = append(orderedRecs, map[string]any{
			"rank":              i + 1,
			"recommendationRef": id,
			"impact":            "medium",
			"confidence":        "medium",
			"rationale":         "Highest expected value.",
		})
	}
	coverage := make([]any, 0, len(findings))
	for _, finding := range findings {
		coverage = append(coverage, map[string]any{
			"findingRef":         finding["findingRef"],
			"disposition":        "addressed_by_recommendation",
			"recommendationRefs": []any{firstID},
		})
	}
	return map[string]any{
		"orderedRecommendations": orderedRecs,
		"findingCoverage":        coverage,
	}, nil
}

// contextFindings reads the findings index from a work-request context in its
// JSON shape.
func contextFindings(req evaluator.WorkRequest) ([]map[string]any, error) {
	raw, err := json.Marshal(req.Context["findings"])
	if err != nil {
		return nil, err
	}
	var findings []map[string]any
	if err := json.Unmarshal(raw, &findings); err != nil {
		return nil, err
	}
	return findings, nil
}

func roundTripJSON(payload map[string]any) map[string]any {
	raw, err := json.Marshal(payload)
	if err != nil {
		return payload
	}
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.UseNumber()
	var out map[string]any
	if err := decoder.Decode(&out); err != nil {
		return payload
	}
	return out
}

func scriptedSelect(scripted evaluator.Evaluator) func(evaluator.Options) (*evaluator.Selection, error) {
	return func(_ evaluator.Options) (*evaluator.Selection, error) {
		return &evaluator.Selection{Name: "scripted", Evaluator: scripted, Reason: "test"}, nil
	}
}

func TestRunCompletesEndToEnd(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	runAbs := filepath.Join(repo, ".quality", "evaluations", "0001-full-eval")
	for _, name := range []string{ArtifactFile, "model-snapshot.md", "report.md", "findings.md", "recommendations.md", "logs/events.jsonl", "logs/evaluator-calls.jsonl"} {
		if _, err := os.Stat(filepath.Join(runAbs, filepath.FromSlash(name))); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(runAbs, "data")); !os.IsNotExist(err) {
		t.Errorf("runner run must not create the multi-file data tree: %v", err)
	}
	if result.ReportMD == "" || result.RatingResult == nil {
		t.Errorf("receipt = %+v, want reportMd and ratingResult", result)
	}
	if result.Concurrency != defaultConcurrency() {
		t.Errorf("concurrency = %d, want %d", result.Concurrency, defaultConcurrency())
	}
	if result.WorkUnits.Completed != result.WorkUnits.Total {
		t.Errorf("workUnits = %+v, want all completed", result.WorkUnits)
	}

	artifact, err := NewStore(runAbs).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if artifact.SchemaVersion != ArtifactSchemaVersion || artifact.Kind != ArtifactKind {
		t.Errorf("artifact header = %d/%s", artifact.SchemaVersion, artifact.Kind)
	}
	if artifact.State.Status != StatusCompleted {
		t.Errorf("artifact status = %q", artifact.State.Status)
	}
	if artifact.Outputs == nil || artifact.Outputs.EvaluationOutput == nil {
		t.Errorf("artifact outputs missing")
	}
	// Status over the runner artifact reports reportable with no gaps.
	status, err := Status(runAbs, "")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if !status.Reportable || len(status.Gaps) != 0 {
		t.Errorf("status = %+v, want reportable", status)
	}
}

func TestRunUsesConfiguredConcurrencyForConcurrentEvaluator(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 2\n")
	scripted := newOverlappingEvaluator()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	if result.Concurrency != 2 {
		t.Errorf("concurrency = %d, want configured 2", result.Concurrency)
	}
	if scripted.maxActive() < 2 {
		t.Errorf("max active assessRateRequirement calls = %d, want at least 2", scripted.maxActive())
	}
}

func TestRunResolvesUnsupportedEvaluatorConcurrencyToOne(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 8\n")
	scripted := &sequentialOnlyEvaluator{inner: newScriptedEvaluator()}
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Concurrency != 1 {
		t.Errorf("concurrency = %d, want 1 for evaluator without concurrent support", result.Concurrency)
	}
}

func TestDryRunRejectsInvalidConfiguredConcurrency(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 0\n")
	scripted := newScriptedEvaluator()
	_, err := DryRun(Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err == nil {
		t.Fatal("DryRun() error = nil, want usage error for concurrency 0")
	}
	var usage *UsageError
	if !errors.As(err, &usage) {
		t.Fatalf("DryRun() error = %T, want UsageError", err)
	}
}

func TestRunRetriesSchemaInvalidOutput(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	scripted.failUnits["assessRateRequirement:requirement:root::has-docs"] = 1
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q, want completed after retry", result.Status)
	}
	if got := scripted.calls["assessRateRequirement:requirement:root::has-docs"]; got != 2 {
		t.Errorf("attempts = %d, want 2", got)
	}
}

// partialEvaluator strips the rating from the combined requirement judgment
// composite for a count of leading attempts on one work unit.
type partialEvaluator struct {
	inner   *scriptedEvaluator
	unit    string
	partial int
}

func (p *partialEvaluator) Kind() string                         { return p.inner.Kind() }
func (p *partialEvaluator) Capabilities() evaluator.Capabilities { return p.inner.Capabilities() }

func (p *partialEvaluator) Evaluate(ctx context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	result, err := p.inner.Evaluate(ctx, req)
	if err == nil && req.WorkUnitID == p.unit && p.partial > 0 {
		p.partial--
		delete(result.Payload, "rating")
	}
	return result, err
}

func TestRunRetriesPartialCombinedResult(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	unit := "assessRateRequirement:requirement:root::has-tests"
	wrapped := &partialEvaluator{inner: scripted, unit: unit, partial: 1}
	result, err := Run(context.Background(), Options{
		RepoRoot: repo,
		Model:    filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: func(_ evaluator.Options) (*evaluator.Selection, error) {
			return &evaluator.Selection{Name: "scripted", Evaluator: wrapped, Reason: "test"}, nil
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed after retrying the partial composite", result.Status, result.Failure)
	}
	if got := scripted.calls[unit]; got != 2 {
		t.Errorf("attempts = %d, want 2 (partial composite must retry, not persist)", got)
	}
}

func TestRunFailsThenResumesWithoutRepeatingAcceptedWork(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	scripted.failUnits["analyzeFactor:factor:root::reliability"] = maxUnitAttempts
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusFailed {
		t.Fatalf("status = %q, want failed", result.Status)
	}
	if result.Failure == nil || result.Failure.Category != evaluator.FailureSchemaInvalidOutput {
		t.Fatalf("failure = %+v, want schema_invalid_output", result.Failure)
	}
	assessCalls := scripted.calls["assessRateRequirement:requirement:root::has-tests"]

	resumed, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          result.Path,
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run(resume) error = %v", err)
	}
	if resumed.Status != StatusCompleted {
		t.Fatalf("resumed status = %q (failure: %+v), want completed", resumed.Status, resumed.Failure)
	}
	if got := scripted.calls["assessRateRequirement:requirement:root::has-tests"]; got != assessCalls {
		t.Errorf("accepted requirement judgment was re-invoked on resume: calls = %d, want %d", got, assessCalls)
	}
}

func TestResumeRefusesEvaluatorConflict(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	_, err = Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          result.Path,
		Evaluator:       "claude",
		SelectEvaluator: scriptedSelect(scripted),
	})
	var usageErr *UsageError
	if err == nil || !errors.As(err, &usageErr) {
		t.Fatalf("Run(resume, --evaluator claude) error = %v, want evaluator-conflict usage error", err)
	}
}

func TestResumeMissingArtifactFails(t *testing.T) {
	repo := testRunnerRepo(t)
	runDir := filepath.Join(repo, ".quality", "evaluations", "0001-full-eval")
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		t.Fatalf("MkdirAll error = %v", err)
	}
	_, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          ".quality/evaluations/0001-full-eval",
		SelectEvaluator: scriptedSelect(newScriptedEvaluator()),
	})
	var runErr *RunError
	if err == nil || !errors.As(err, &runErr) || runErr.Category != evaluator.FailureRunStateInvalid {
		t.Fatalf("Run(resume) error = %v, want run_state_invalid", err)
	}
}

func TestRunCancellationLeavesResumableArtifact(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	ctx, cancel := context.WithCancel(context.Background())
	cancelAfter := "assessRateRequirement:requirement:root::has-docs"
	wrapped := &cancellingEvaluator{inner: scripted, cancel: cancel, after: cancelAfter}
	result, err := Run(ctx, Options{
		RepoRoot: repo,
		Model:    filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: func(_ evaluator.Options) (*evaluator.Selection, error) {
			return &evaluator.Selection{Name: "scripted", Evaluator: wrapped, Reason: "test"}, nil
		},
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCancelled {
		t.Fatalf("status = %q, want cancelled", result.Status)
	}
	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if artifact.State.Status != StatusCancelled || !artifact.State.Cancelled {
		t.Errorf("artifact state = %+v, want cancelled", artifact.State)
	}

	resumed, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          result.Path,
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run(resume) error = %v", err)
	}
	if resumed.Status != StatusCompleted {
		t.Fatalf("resumed status = %q, want completed", resumed.Status)
	}
}

// cancellingEvaluator cancels the run context after completing a named unit.
type cancellingEvaluator struct {
	inner  *scriptedEvaluator
	cancel context.CancelFunc
	after  string
}

func (c *cancellingEvaluator) Kind() string                         { return c.inner.Kind() }
func (c *cancellingEvaluator) Capabilities() evaluator.Capabilities { return c.inner.Capabilities() }

func (c *cancellingEvaluator) Evaluate(ctx context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	result, err := c.inner.Evaluate(ctx, req)
	if req.WorkUnitID == c.after {
		c.cancel()
	}
	return result, err
}

// harnessSelect resolves the real built-in harness evaluator selection.
func harnessSelect() func(evaluator.Options) (*evaluator.Selection, error) {
	return func(_ evaluator.Options) (*evaluator.Selection, error) {
		return evaluator.Select(evaluator.Options{Name: "harness"})
	}
}

// workRequestFor converts an awaiting receipt's bounded request back into the
// evaluator work-request shape the scripted payload registry consumes.
func workRequestFor(req *runnerEvaluatorRequest) evaluator.WorkRequest {
	return evaluator.WorkRequest{
		WorkUnitID: req.WorkUnitID,
		Kind:       req.Kind,
		Subject:    req.Subject,
		Context:    req.Context,
	}
}

type runnerEvaluatorRequest = EvaluatorRequest

// soleRequest returns the receipt's single outstanding work request.
func soleRequest(t *testing.T, result *Result) *EvaluatorRequest {
	t.Helper()
	if result.Status != StatusAwaitingEvaluator {
		t.Fatalf("status = %q (failure: %+v), want awaiting_evaluator", result.Status, result.Failure)
	}
	if len(result.EvaluatorRequests) != 1 {
		t.Fatalf("outstanding requests = %d, want exactly 1", len(result.EvaluatorRequests))
	}
	return result.EvaluatorRequests[0]
}

// submitHarnessResult resumes an awaiting run with one result envelope,
// exercising the single-object submission shape.
func submitHarnessResult(t *testing.T, repo, runPath string, envelope evaluator.HarnessResultEnvelope) (*Result, error) {
	t.Helper()
	raw, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("marshaling envelope: %v", err)
	}
	return resumeWithResult(repo, runPath, raw)
}

// submitHarnessResults resumes an awaiting run with a batch of result
// envelopes in one submission, exercising the JSON-array shape.
func submitHarnessResults(t *testing.T, repo, runPath string, envelopes []evaluator.HarnessResultEnvelope) (*Result, error) {
	t.Helper()
	raw, err := json.Marshal(envelopes)
	if err != nil {
		t.Fatalf("marshaling envelopes: %v", err)
	}
	return resumeWithResult(repo, runPath, raw)
}

func resumeWithResult(repo, runPath string, raw []byte) (*Result, error) {
	return Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          runPath,
		EvaluatorResult: "-",
		Stdin:           bytes.NewReader(raw),
		SelectEvaluator: harnessSelect(),
	})
}

// envelopeFor builds the scripted result envelope answering one outstanding
// request.
func envelopeFor(t *testing.T, scripted *scriptedEvaluator, req *EvaluatorRequest) evaluator.HarnessResultEnvelope {
	t.Helper()
	payload, err := scripted.payloadFor(workRequestFor(req))
	if err != nil {
		t.Fatalf("scripted payload for %s: %v", req.WorkUnitID, err)
	}
	return evaluator.HarnessResultEnvelope{
		RequestID: req.RequestID,
		InputHash: req.InputHash,
		Evaluator: evaluator.HarnessIdentity{Runtime: "claude-code"},
		Payload:   roundTripJSON(payload),
	}
}

// serviceHarnessRun services awaiting receipts until the run is terminal,
// answering every outstanding request as one batch per resume. It returns the
// terminal receipt and the number of awaiting receipts serviced.
func serviceHarnessRun(t *testing.T, repo string, result *Result, scripted *scriptedEvaluator) (*Result, int) {
	t.Helper()
	checkpoints := 0
	for result.Status == StatusAwaitingEvaluator {
		checkpoints++
		if checkpoints > 20 {
			t.Fatalf("run did not complete after %d checkpoints", checkpoints)
		}
		if len(result.EvaluatorRequests) == 0 {
			t.Fatalf("awaiting receipt carries no outstanding requests")
		}
		envelopes := make([]evaluator.HarnessResultEnvelope, 0, len(result.EvaluatorRequests))
		for _, req := range result.EvaluatorRequests {
			if req.RequestID == "" || req.InputHash == "" || len(req.ExpectedSchema) == 0 {
				t.Fatalf("awaiting receipt lacks a complete bounded request: %+v", req)
			}
			envelopes = append(envelopes, envelopeFor(t, scripted, req))
		}
		var err error
		result, err = submitHarnessResults(t, repo, result.Path, envelopes)
		if err != nil {
			t.Fatalf("submit batch error = %v", err)
		}
	}
	return result, checkpoints
}

func startHarnessRun(t *testing.T, repo string) *Result {
	t.Helper()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Evaluator:       "harness",
		SelectEvaluator: harnessSelect(),
	})
	if err != nil {
		t.Fatalf("Run(harness) error = %v", err)
	}
	return result
}

func TestHarnessRunCheckpointsToCompletion(t *testing.T) {
	repo := testRunnerRepo(t)
	// Concurrency 1 preserves the single-request checkpoint loop: every
	// awaiting receipt carries exactly one outstanding request.
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 1\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if result.Concurrency != 1 {
		t.Fatalf("concurrency = %d, want configured 1", result.Concurrency)
	}
	checkpoints := 0
	for result.Status == StatusAwaitingEvaluator {
		checkpoints++
		if checkpoints > 20 {
			t.Fatalf("run did not complete after %d checkpoints", checkpoints)
		}
		req := soleRequest(t, result)
		if req.RequestID == "" || req.InputHash == "" || len(req.ExpectedSchema) == 0 {
			t.Fatalf("awaiting receipt lacks the bounded request: %+v", req)
		}
		envelope := envelopeFor(t, scripted, req)
		envelope.Evaluator.Model = "test-model"
		var err error
		result, err = submitHarnessResult(t, repo, result.Path, envelope)
		if err != nil {
			t.Fatalf("submit for %s error = %v", req.WorkUnitID, err)
		}
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	// The test model has 7 evaluator-backed units: one checkpoint each.
	if checkpoints != 7 {
		t.Errorf("checkpoints = %d, want 7", checkpoints)
	}
	if result.Evaluator != "harness" || result.EvaluatorKind != "harness" {
		t.Errorf("receipt evaluator = %s/%s, want harness", result.Evaluator, result.EvaluatorKind)
	}
	if result.ReportMD == "" || result.RatingResult == nil {
		t.Errorf("receipt = %+v, want reportMd and ratingResult", result)
	}
	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if artifact.State.HarnessIdentity == nil || artifact.State.HarnessIdentity.Runtime != "claude-code" {
		t.Errorf("harness identity = %+v, want claude-code", artifact.State.HarnessIdentity)
	}
	if len(artifact.State.PendingEvaluatorCalls) != 0 {
		t.Errorf("completed run still has pending evaluator calls")
	}
}

func TestHarnessWindowServicesConcurrentRequests(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 4\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if result.Concurrency != 4 {
		t.Fatalf("concurrency = %d, want configured 4 (harness is no longer clamped to 1)", result.Concurrency)
	}
	// Both requirement judgments are dependency-ready and independent, so the
	// first receipt carries both outstanding requests at once.
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("outstanding requests = %d, want the 2 independent requirement judgments", len(result.EvaluatorRequests))
	}
	for _, req := range result.EvaluatorRequests {
		if UnitKind(req.Kind) != KindAssessRateRequirement {
			t.Errorf("outstanding request kind = %s, want assessRateRequirement", req.Kind)
		}
	}
	result, checkpoints := serviceHarnessRun(t, repo, result, scripted)
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	// 7 evaluator units serviced over 5 receipts: the first batches the two
	// independent requirement judgments, the second batches the factor
	// analysis with the finding ranking (both depend only on the requirement
	// judgments), and the rest are DAG-serial.
	if checkpoints != 5 {
		t.Errorf("checkpoints = %d, want 5", checkpoints)
	}
}

func TestHarnessWindowMatchesSequentialResults(t *testing.T) {
	accepted := func(t *testing.T, config string) (*Results, *evaluation.RatingResult) {
		t.Helper()
		repo := testRunnerRepo(t)
		writeRunnerConfig(t, repo, config)
		scripted := newScriptedEvaluator()
		result, _ := serviceHarnessRun(t, repo, startHarnessRun(t, repo), scripted)
		if result.Status != StatusCompleted {
			t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
		}
		artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		return &artifact.Results, result.RatingResult
	}
	sequential, sequentialRating := accepted(t, "evaluation:\n  concurrency: 1\n")
	windowed, windowedRating := accepted(t, "evaluation:\n  concurrency: 4\n")
	if len(sequential.Payloads) != len(windowed.Payloads) {
		t.Fatalf("accepted payloads = %d vs %d, want identical sets", len(windowed.Payloads), len(sequential.Payloads))
	}
	// Recommendation IDs are randomly assigned per run, so payloads carrying
	// them are compared after mapping each distinct qrec token, in first-seen
	// order, to an ordinal placeholder.
	qrecToken := regexp.MustCompile(`qrec_[a-z0-9]+`)
	normalized := func(payload map[string]any) string {
		raw, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshaling payload: %v", err)
		}
		seen := map[string]string{}
		return qrecToken.ReplaceAllStringFunc(string(raw), func(id string) string {
			if _, ok := seen[id]; !ok {
				seen[id] = fmt.Sprintf("qrec_%04d", len(seen))
			}
			return seen[id]
		})
	}
	for i, want := range sequential.Payloads {
		got := windowed.Payloads[i]
		if got.WorkUnit != want.WorkUnit {
			t.Fatalf("payload %d work unit = %s, want %s (order must be graph order)", i, got.WorkUnit, want.WorkUnit)
		}
		// The evaluation frame embeds run identity (evaluationId, createdAt),
		// which legitimately differs between two runs.
		if want.WorkUnit == string(KindFrameEvaluation) {
			continue
		}
		if normalized(got.Payload) != normalized(want.Payload) {
			t.Errorf("payload for %s differs between window sizes", want.WorkUnit)
		}
	}
	if sequentialRating == nil || windowedRating == nil || windowedRating.Level != sequentialRating.Level {
		t.Errorf("rating = %+v vs %+v, want identical", windowedRating, sequentialRating)
	}
}

func TestHarnessAwaitingRunStatusSurfacesPendingRequest(t *testing.T) {
	repo := testRunnerRepo(t)
	result := startHarnessRun(t, repo)
	if result.Status != StatusAwaitingEvaluator {
		t.Fatalf("status = %q, want awaiting_evaluator", result.Status)
	}
	status, err := Status(filepath.Join(repo, result.Path), result.Path)
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	if status.Reportable {
		t.Errorf("awaiting run must be incomplete, got reportable")
	}
	if status.Lifecycle != StatusAwaitingEvaluator {
		t.Errorf("lifecycle = %q, want awaiting_evaluator", status.Lifecycle)
	}
	if len(status.AwaitingEvaluator) != len(result.EvaluatorRequests) {
		t.Fatalf("awaitingEvaluator = %+v, want the %d outstanding requests", status.AwaitingEvaluator, len(result.EvaluatorRequests))
	}
	for i, awaiting := range status.AwaitingEvaluator {
		request := result.EvaluatorRequests[i]
		if awaiting.WorkUnitID != request.WorkUnitID || awaiting.RequestID != request.RequestID {
			t.Errorf("awaitingEvaluator[%d] = %+v, want the outstanding request %+v", i, awaiting, request)
		}
	}
	if len(status.NextActions) == 0 || status.NextActions[0].ID != "evaluation-run-reemit" {
		t.Errorf("nextActions = %+v, want the resume continuation first", status.NextActions)
	}
}

func TestHarnessResumeWithoutResultReemitsSameRequests(t *testing.T) {
	repo := testRunnerRepo(t)
	first := startHarnessRun(t, repo)
	if first.Status != StatusAwaitingEvaluator {
		t.Fatalf("status = %q, want awaiting_evaluator", first.Status)
	}
	again, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          first.Path,
		SelectEvaluator: harnessSelect(),
	})
	if err != nil {
		t.Fatalf("Run(resume) error = %v", err)
	}
	if again.Status != StatusAwaitingEvaluator {
		t.Fatalf("resumed status = %q, want awaiting_evaluator", again.Status)
	}
	if len(again.EvaluatorRequests) != len(first.EvaluatorRequests) {
		t.Fatalf("re-emitted set = %d requests, want %d", len(again.EvaluatorRequests), len(first.EvaluatorRequests))
	}
	for i, request := range again.EvaluatorRequests {
		original := first.EvaluatorRequests[i]
		if request.RequestID != original.RequestID || request.InputHash != original.InputHash {
			t.Errorf("re-emitted request %+v differs from original %+v", request, original)
		}
	}
}

func TestHarnessRejectsMismatchedResult(t *testing.T) {
	repo := testRunnerRepo(t)
	result := startHarnessRun(t, repo)
	_, err := submitHarnessResult(t, repo, result.Path, evaluator.HarnessResultEnvelope{
		RequestID: "req_other",
		InputHash: result.EvaluatorRequests[0].InputHash,
		Evaluator: evaluator.HarnessIdentity{Runtime: "claude-code"},
		Payload:   map[string]any{"any": true},
	})
	var runErr *RunError
	if err == nil || !errors.As(err, &runErr) || runErr.Category != evaluator.FailureRunStateInvalid {
		t.Fatalf("submit(mismatched) error = %v, want run_state_invalid", err)
	}
	// The outstanding requests are untouched and recoverable.
	again, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          result.Path,
		SelectEvaluator: harnessSelect(),
	})
	if err != nil || again.Status != StatusAwaitingEvaluator ||
		len(again.EvaluatorRequests) != len(result.EvaluatorRequests) ||
		again.EvaluatorRequests[0].RequestID != result.EvaluatorRequests[0].RequestID {
		t.Fatalf("outstanding requests not recoverable after rejected submit: %+v, %v", again, err)
	}
}

func TestHarnessRejectsStaleInput(t *testing.T) {
	repo := testRunnerRepo(t)
	result := startHarnessRun(t, repo)
	if err := os.WriteFile(filepath.Join(repo, "src", "main.txt"), []byte("changed content\n"), 0o644); err != nil {
		t.Fatalf("WriteFile error = %v", err)
	}
	_, err := submitHarnessResult(t, repo, result.Path, evaluator.HarnessResultEnvelope{
		RequestID: result.EvaluatorRequests[0].RequestID,
		InputHash: result.EvaluatorRequests[0].InputHash,
		Evaluator: evaluator.HarnessIdentity{Runtime: "claude-code"},
		Payload:   map[string]any{"any": true},
	})
	var runErr *RunError
	if err == nil || !errors.As(err, &runErr) || runErr.Category != evaluator.FailureRunStateInvalid {
		t.Fatalf("submit(stale source) error = %v, want run_state_invalid", err)
	}
}

func TestHarnessRefusesMixedRuntimes(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	result, err := submitHarnessResult(t, repo, result.Path,
		envelopeFor(t, scripted, result.EvaluatorRequests[0]))
	if err != nil || result.Status != StatusAwaitingEvaluator {
		t.Fatalf("first submit: %+v, %v", result, err)
	}
	_, err = submitHarnessResult(t, repo, result.Path, evaluator.HarnessResultEnvelope{
		RequestID: result.EvaluatorRequests[0].RequestID,
		InputHash: result.EvaluatorRequests[0].InputHash,
		Evaluator: evaluator.HarnessIdentity{Runtime: "codex"},
		Payload:   map[string]any{"any": true},
	})
	var runErr *RunError
	if err == nil || !errors.As(err, &runErr) || runErr.Category != evaluator.FailureRunStateInvalid {
		t.Fatalf("submit(mixed runtime) error = %v, want run_state_invalid", err)
	}
}

func TestHarnessRetriesInvalidPayloadThenFails(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 1\n")
	result := startHarnessRun(t, repo)
	for attempt := 1; attempt <= maxUnitAttempts; attempt++ {
		request := soleRequest(t, result)
		if got := request.Attempt; got != attempt {
			t.Fatalf("request attempt = %d, want %d", got, attempt)
		}
		if attempt > 1 && request.LastFailure == nil {
			t.Errorf("retry request lacks the classified previous-attempt failure")
		}
		var err error
		result, err = submitHarnessResult(t, repo, result.Path, evaluator.HarnessResultEnvelope{
			RequestID: request.RequestID,
			InputHash: request.InputHash,
			Evaluator: evaluator.HarnessIdentity{Runtime: "claude-code"},
			Payload:   map[string]any{"broken": true},
		})
		if err != nil {
			t.Fatalf("submit attempt %d error = %v", attempt, err)
		}
	}
	if result.Status != StatusFailed {
		t.Fatalf("status = %q, want failed after exhausting retries", result.Status)
	}
	if result.Failure == nil || result.Failure.Category != evaluator.FailureInvalidEvaluatorOutput {
		t.Fatalf("failure = %+v, want invalid_evaluator_output", result.Failure)
	}
}

func TestHarnessPartialSubmissionKeepsRemainderOutstanding(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 4\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("outstanding requests = %d, want 2", len(result.EvaluatorRequests))
	}
	// Answer the has-tests judgment: accepting it makes the reliability
	// factor analysis dependency-ready, so the window must top up without the
	// unanswered member being touched.
	var held, answered *EvaluatorRequest
	for _, request := range result.EvaluatorRequests {
		if strings.HasSuffix(request.WorkUnitID, "has-tests") {
			answered = request
		} else {
			held = request
		}
	}
	if held == nil || answered == nil {
		t.Fatalf("outstanding requests = %+v, want both requirement judgments", result.EvaluatorRequests)
	}
	result, err := submitHarnessResult(t, repo, result.Path, envelopeFor(t, scripted, answered))
	if err != nil {
		t.Fatalf("partial submit error = %v", err)
	}
	if result.Status != StatusAwaitingEvaluator || len(result.EvaluatorRequests) != 2 {
		t.Fatalf("receipt after partial submit = %+v, want the held request plus the topped-up analysis", result)
	}
	var remainder, toppedUp *EvaluatorRequest
	for _, request := range result.EvaluatorRequests {
		if request.WorkUnitID == held.WorkUnitID {
			remainder = request
		} else {
			toppedUp = request
		}
	}
	if remainder == nil || remainder.RequestID != held.RequestID || remainder.Attempt != held.Attempt {
		t.Fatalf("remainder = %+v, want the unanswered request unchanged %+v", remainder, held)
	}
	if remainder.LastFailure != nil {
		t.Errorf("unanswered request carries a failure: %+v", remainder.LastFailure)
	}
	if toppedUp == nil || UnitKind(toppedUp.Kind) != KindAnalyzeFactor {
		t.Errorf("topped-up request = %+v, want the newly-ready factor analysis", toppedUp)
	}
	// A not-yet-judged member consumes no retry budget.
	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if attempts := artifact.State.WorkUnits[held.WorkUnitID].Attempts; attempts != 0 {
		t.Errorf("unanswered unit attempts = %d, want 0 (not submitted is not a failed attempt)", attempts)
	}
	if state := artifact.State.WorkUnits[answered.WorkUnitID]; state.Status != UnitCompleted {
		t.Errorf("answered unit state = %+v, want completed and durable", state)
	}
}

func TestHarnessMixedSubmissionAcceptsValidAndRetriesFailed(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 4\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("outstanding requests = %d, want 2", len(result.EvaluatorRequests))
	}
	valid := result.EvaluatorRequests[0]
	broken := result.EvaluatorRequests[1]
	result, err := submitHarnessResults(t, repo, result.Path, []evaluator.HarnessResultEnvelope{
		envelopeFor(t, scripted, valid),
		{
			RequestID: broken.RequestID,
			InputHash: broken.InputHash,
			Evaluator: evaluator.HarnessIdentity{Runtime: "claude-code"},
			Payload:   map[string]any{"broken": true},
		},
	})
	if err != nil {
		t.Fatalf("mixed submit error = %v", err)
	}
	retry := soleRequest(t, result)
	if retry.WorkUnitID != broken.WorkUnitID || retry.Attempt != 2 {
		t.Fatalf("retry request = %+v, want %s attempt 2", retry, broken.WorkUnitID)
	}
	if retry.LastFailure == nil || retry.LastFailure.Category != evaluator.FailureInvalidEvaluatorOutput {
		t.Errorf("retry lastFailure = %+v, want invalid_evaluator_output", retry.LastFailure)
	}
	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if state := artifact.State.WorkUnits[valid.WorkUnitID]; state.Status != UnitCompleted {
		t.Errorf("valid member state = %+v, want completed despite the failed sibling", state)
	}
}

func TestHarnessDuplicateEnvelopeRejectedAfterValidMembers(t *testing.T) {
	repo := testRunnerRepo(t)
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 4\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("outstanding requests = %d, want 2", len(result.EvaluatorRequests))
	}
	answered := result.EvaluatorRequests[0]
	envelope := envelopeFor(t, scripted, answered)
	_, err := submitHarnessResults(t, repo, result.Path, []evaluator.HarnessResultEnvelope{envelope, envelope})
	var runErr *RunError
	if err == nil || !errors.As(err, &runErr) || runErr.Category != evaluator.FailureRunStateInvalid {
		t.Fatalf("submit(duplicate) error = %v, want run_state_invalid", err)
	}
	// The first member's acceptance is durable; the duplicate touched nothing.
	again, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          result.Path,
		SelectEvaluator: harnessSelect(),
	})
	if err != nil {
		t.Fatalf("Run(resume) error = %v", err)
	}
	for _, request := range again.EvaluatorRequests {
		if request.WorkUnitID == answered.WorkUnitID {
			t.Errorf("accepted member re-emitted after duplicate rejection: %+v", request)
		}
	}
	artifact, err := NewStore(filepath.Join(repo, again.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if state := artifact.State.WorkUnits[answered.WorkUnitID]; state.Status != UnitCompleted {
		t.Errorf("accepted member state = %+v, want completed", state)
	}
}

func TestHarnessWindowBoundedByResolvedConcurrency(t *testing.T) {
	const threeRequirementModel = `---
title: Window bound model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: Minimum
    description: Minimum.
    criterion: Barely meets it.
factors:
  reliability:
    title: Reliability
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
    factors: [reliability]
  has-docs:
    title: Has docs
    assessment: Inspect docs.
    factors: [reliability]
  has-ci:
    title: Has CI
    assessment: Inspect CI.
    factors: [reliability]
source: src
---
`
	repo := testRunnerRepo(t)
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(threeRequirementModel), 0o644); err != nil {
		t.Fatalf("WriteFile(QUALITY.md) error = %v", err)
	}
	writeRunnerConfig(t, repo, "evaluation:\n  concurrency: 2\n")
	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	// Three requirement judgments are dependency-ready, but the window is
	// capped at the resolved concurrency.
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("outstanding requests = %d, want the window capped at 2", len(result.EvaluatorRequests))
	}
	result, err := submitHarnessResult(t, repo, result.Path,
		envelopeFor(t, scripted, result.EvaluatorRequests[0]))
	if err != nil {
		t.Fatalf("submit error = %v", err)
	}
	// Accepting one member frees capacity, and the window tops up with the
	// deferred third judgment without waiting for the whole wave.
	if len(result.EvaluatorRequests) != 2 {
		t.Fatalf("topped-up requests = %d, want 2", len(result.EvaluatorRequests))
	}
	terminal, _ := serviceHarnessRun(t, repo, result, scripted)
	if terminal.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", terminal.Status, terminal.Failure)
	}
}

func TestHarnessResultRequiresResumeAndPendingRequest(t *testing.T) {
	repo := testRunnerRepo(t)
	_, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		EvaluatorResult: "-",
		Stdin:           strings.NewReader("{}"),
		SelectEvaluator: harnessSelect(),
	})
	var usageErr *UsageError
	if err == nil || !errors.As(err, &usageErr) {
		t.Fatalf("Run(--evaluator-result without --resume) error = %v, want usage error", err)
	}

	// A completed non-harness run refuses a harness result.
	scripted := newScriptedEvaluator()
	completed, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	_, err = Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          completed.Path,
		EvaluatorResult: "-",
		Stdin:           strings.NewReader(`{"requestId":"r","inputHash":"h","evaluator":{"runtime":"claude-code"},"payload":{}}`),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err == nil || !errors.As(err, &usageErr) {
		t.Fatalf("Run(result for non-harness run) error = %v, want usage error", err)
	}
}

func TestDryRunPreviewsWithoutWriting(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	preview, err := DryRun(Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if preview.ExpectedRunPath != ".quality/evaluations/0001-full-eval" {
		t.Errorf("expectedRunPath = %q", preview.ExpectedRunPath)
	}
	if preview.Concurrency != defaultConcurrency() {
		t.Errorf("concurrency = %d, want %d", preview.Concurrency, defaultConcurrency())
	}
	if preview.WorkUnits.Total == 0 || preview.WorkUnits.EvaluatorUnits == 0 {
		t.Errorf("workUnits = %+v", preview.WorkUnits)
	}
	if len(scripted.calls) != 0 {
		t.Errorf("dry run invoked the evaluator: %v", scripted.calls)
	}
	if _, err := os.Stat(filepath.Join(repo, ".quality")); !os.IsNotExist(err) {
		t.Errorf("dry run wrote workspace state: %v", err)
	}
}

func TestBuildGraphOrderAndDependencies(t *testing.T) {
	repo := testRunnerRepo(t)
	scripted := newScriptedEvaluator()
	preview, err := DryRun(Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	// Root scope of the test model: 1 evaluation frame, 1 area frame,
	// 2 requirements x2 units (frame + combined assess-and-rate), 1 factor x2,
	// 1 area analysis x2, 3 advice, 1 report build.
	if preview.WorkUnits.Total != 14 {
		t.Errorf("total units = %d, want 14", preview.WorkUnits.Total)
	}
	if preview.WorkUnits.EvaluatorUnits != 7 {
		t.Errorf("evaluator units = %d, want 7", preview.WorkUnits.EvaluatorUnits)
	}
}
