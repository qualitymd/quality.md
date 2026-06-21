package status

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

func TestSnapshotMissingModelSucceeds(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Readiness != ReadinessMissingModel || snapshot.Model.Present {
		t.Fatalf("snapshot = %#v, want missing model", snapshot)
	}
	if len(snapshot.NextActions) != 1 || snapshot.NextActions[0].ID != "init" {
		t.Fatalf("nextActions = %#v, want init action", snapshot.NextActions)
	}
}

func TestSnapshotInvalidModelCarriesLintFindings(t *testing.T) {
	path := writeFile(t, t.TempDir(), `---
title: Invalid
---
`)
	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Readiness != ReadinessInvalidModel || snapshot.Model.Valid {
		t.Fatalf("snapshot readiness = %q valid = %v, want invalid", snapshot.Readiness, snapshot.Model.Valid)
	}
	if snapshot.Model.Shape != nil {
		t.Fatalf("shape = %#v, want nil for invalid model", snapshot.Model.Shape)
	}
	if len(snapshot.Model.Lint.Findings) == 0 {
		t.Fatal("lint findings are empty, want invalid-model details")
	}
}

func TestSnapshotValidModelShapeAndSourceCoverage(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`source: .
factors:
  reliability:
    title: Reliability
    description: Reliability.
    factors:
      errors:
        title: Error handling
        description: Error handling.
        requirements:
          "reports errors":
            assessment: Inspect errors.
requirements:
  "starts":
    factors: [reliability]
    assessment: Run it.
areas:
  api:
    title: API
    requirements:
      "responds":
        factors: [reliability]
        assessment: Call it.
  cli:
    title: CLI
    source: ./cmd
    factors:
      usability:
        title: Usability
        description: Usability.
`))
	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Readiness != ReadinessReadyToEvaluate {
		t.Fatalf("readiness = %q, want ready", snapshot.Readiness)
	}
	shape := snapshot.Model.Shape
	if shape == nil {
		t.Fatal("shape is nil")
	}
	if shape.Areas != 3 || shape.Factors != 3 || shape.Requirements != 3 || shape.RatingLevels != 2 {
		t.Fatalf("shape = %#v, want 3 areas, 3 factors, 3 requirements, 2 levels", *shape)
	}
	if len(snapshot.Model.SourceCoverage) != 3 {
		t.Fatalf("sourceCoverage = %#v, want root plus 2 child areas", snapshot.Model.SourceCoverage)
	}
	if snapshot.Model.SourceCoverage[1].Label != "API" || snapshot.Model.SourceCoverage[1].SourceState != "inherited" || snapshot.Model.SourceCoverage[1].Source != "." {
		t.Fatalf("api source coverage = %#v, want inherited root source", snapshot.Model.SourceCoverage[1])
	}
	if snapshot.Model.SourceCoverage[2].Label != "CLI" || snapshot.Model.SourceCoverage[2].SourceState != "declared" || snapshot.Model.SourceCoverage[2].Source != "./cmd" {
		t.Fatalf("cli source coverage = %#v, want declared source", snapshot.Model.SourceCoverage[2])
	}
}

func TestSnapshotEvaluationHistoryStaleAndLatestRun(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  "starts":
    factors: [reliability]
    assessment: EvaluationRun it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	first, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(first) error = %v", err)
	}
	if !strings.Contains(first.Path, "0001-quality-eval") {
		t.Fatalf("first path = %q, want run 0001", first.Path)
	}
	writeFile(t, repo, validModel(`requirements:
  "starts well":
    factors: [reliability]
    assessment: EvaluationRun it again.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	second, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(second) error = %v", err)
	}

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Evaluations.Runs != 2 || snapshot.Evaluations.Summary.Stale != 1 || snapshot.Evaluations.Summary.Incomplete != 2 {
		t.Fatalf("evaluation history = %#v, want 2 runs, 1 stale, 2 incomplete", snapshot.Evaluations)
	}
	if snapshot.Evaluations.Latest == nil || snapshot.Evaluations.Latest.Path != second.Path {
		t.Fatalf("latest = %#v, want %s", snapshot.Evaluations.Latest, second.Path)
	}
	if snapshot.Readiness != ReadinessNeedsEvaluationReconcile {
		t.Fatalf("readiness = %q, want reconciliation", snapshot.Readiness)
	}
}

func TestSnapshotMalformedRunDoesNotHideLaterRuns(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  "starts":
    factors: [reliability]
    assessment: EvaluationRun it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	if err := os.MkdirAll(filepath.Join(repo, "quality", "evaluations", "0001-quality-eval"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	later, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Evaluations.Runs != 2 || snapshot.Evaluations.Summary.Problems != 1 {
		t.Fatalf("evaluation history = %#v, want one malformed run and one later run", snapshot.Evaluations)
	}
	if snapshot.Evaluations.Latest == nil || snapshot.Evaluations.Latest.Path != later.Path {
		t.Fatalf("latest = %#v, want later run %s", snapshot.Evaluations.Latest, later.Path)
	}
}

func TestSnapshotIncompatibleRunRecordIsHistoryGap(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  "starts":
    factors: [reliability]
    assessment: EvaluationRun it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	run, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, run.Path, "assessments", "001-bad.json"), []byte(`{`), 0o644); err != nil {
		t.Fatalf("write bad assessment: %v", err)
	}

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Readiness != ReadinessNeedsEvaluationReconcile {
		t.Fatalf("readiness = %q, want reconciliation", snapshot.Readiness)
	}
	if snapshot.Evaluations.Summary.Problems != 0 || snapshot.Evaluations.Summary.Incomplete != 1 {
		t.Fatalf("evaluation summary = %#v, want incompatible record as incomplete history gap", snapshot.Evaluations.Summary)
	}
	if snapshot.Evaluations.Latest == nil || snapshot.Evaluations.Latest.Gaps == 0 {
		t.Fatalf("latest = %#v, want latest run with gaps", snapshot.Evaluations.Latest)
	}
}

func TestSnapshotActiveRecommendationCountHonorsSuperseding(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  "starts":
    factors: [reliability]
    assessment: EvaluationRun it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	run, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	addRecommendation(t, filepath.Join(repo, run.Path), `{
  "title": "First",
  "gap": "Gap one.",
  "evidenceLocators": ["QUALITY.md"],
  "assessmentResultRecords": [],
  "remediationOptions": ["Do one."],
  "recommendedOption": "Do one.",
  "doneCriterion": "Done."
}`)
	addRecommendation(t, filepath.Join(repo, run.Path), `{
  "title": "Second",
  "gap": "Gap two.",
  "evidenceLocators": ["QUALITY.md"],
  "assessmentResultRecords": [],
  "remediationOptions": ["Do two."],
  "recommendedOption": "Do two.",
  "doneCriterion": "Done.",
  "supersedes": ["001-first"]
}`)

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Evaluations.Summary.ActiveRecommendations != 1 {
		t.Fatalf("active recommendations = %d, want 1", snapshot.Evaluations.Summary.ActiveRecommendations)
	}
}

func addRecommendation(t *testing.T, runPath, payload string) {
	t.Helper()
	if _, err := evaluation.AddRecord(evaluation.KindRecommendation, runPath, []byte(payload)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
}

func newRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	return repo
}

func writeFile(t *testing.T, dir, content string) string {
	t.Helper()
	path := filepath.Join(dir, "QUALITY.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func validModel(body string) string {
	return `---
title: Example
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
` + body + `---
`
}
