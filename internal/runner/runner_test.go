package runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// scriptedEvaluator produces schema-valid payloads for every judgment work
// unit, tracking calls per work unit. failUnits maps work-unit IDs to a count
// of leading attempts that return a schema-invalid payload.
type scriptedEvaluator struct {
	calls     map[string]int
	failUnits map[string]int
}

var _ evaluator.Evaluator = (*scriptedEvaluator)(nil)

func newScriptedEvaluator() *scriptedEvaluator {
	return &scriptedEvaluator{calls: map[string]int{}, failUnits: map[string]int{}}
}

func (s *scriptedEvaluator) Kind() string { return "scripted" }

func (s *scriptedEvaluator) Capabilities() evaluator.Capabilities {
	return evaluator.Capabilities{Strategies: []evaluator.Strategy{evaluator.StrategySequential}}
}

func (s *scriptedEvaluator) Evaluate(_ context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	s.calls[req.WorkUnitID]++
	result := evaluator.WorkResult{WorkUnitID: req.WorkUnitID, EvaluatorKind: "scripted", Strategy: evaluator.StrategySequential}
	if remaining := s.failUnits[req.WorkUnitID]; remaining > 0 {
		s.failUnits[req.WorkUnitID] = remaining - 1
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

func scriptedSelect(scripted *scriptedEvaluator) func(evaluator.Options) (*evaluator.Selection, error) {
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
	if result.ExecutionStrategy != "sequential" || result.Concurrency != 1 {
		t.Errorf("strategy = %q concurrency = %d, want sequential/1", result.ExecutionStrategy, result.Concurrency)
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
	if preview.ExecutionStrategy != "sequential" || preview.Concurrency != 1 {
		t.Errorf("strategy = %q/%d, want sequential/1", preview.ExecutionStrategy, preview.Concurrency)
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
