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
          reports-errors:
            title: Reports errors
            assessment: Inspect errors.
requirements:
  starts:
    title: Starts
    factors: [reliability]
    assessment: Run it.
areas:
  api:
    title: API
    requirements:
      responds:
        title: Responds
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

func TestSnapshotLeanRootUsesDefaultSourceState(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`factors:
  reliability:
    title: Reliability
    description: Reliability.
areas:
  api:
    title: API
    requirements:
      responds:
        title: Responds
        factors: [reliability]
        assessment: Call it.
`))
	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if len(snapshot.Model.SourceCoverage) != 2 {
		t.Fatalf("sourceCoverage = %#v, want root plus one child area", snapshot.Model.SourceCoverage)
	}
	root := snapshot.Model.SourceCoverage[0]
	if root.SourceState != SourceStateDefault || root.Source != "" {
		t.Fatalf("root source coverage = %#v, want default source state and no declared source", root)
	}
	child := snapshot.Model.SourceCoverage[1]
	if child.SourceState != SourceStateDefault {
		t.Fatalf("child source coverage = %#v, want default source state inherited from document default", child)
	}
}

func TestSnapshotEvaluationHistoryStaleAndLatestRun(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  starts:
    title: Starts
    factors: [reliability]
    assessment: Run it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	first, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(first) error = %v", err)
	}
	if !strings.Contains(first.Path, "0001-full-eval") {
		t.Fatalf("first path = %q, want run 0001", first.Path)
	}
	writeFile(t, repo, validModel(`requirements:
  starts-well:
    title: Starts well
    factors: [reliability]
    assessment: Run it again.
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

func TestSnapshotUsesWorkspaceEvaluationDir(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, strings.Replace(validModel(`requirements:
  starts:
    title: Starts
    factors: [reliability]
    assessment: Run it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`), "---\n", "---\nconfig: .quality/custom-config.yaml\n", 1))
	configPath := filepath.Join(repo, ".quality", "custom-config.yaml")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatalf("MkdirAll(config dir) error = %v", err)
	}
	if err := os.WriteFile(configPath, []byte("evaluationDir: tmp/evals\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(config) error = %v", err)
	}
	run, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if !strings.HasPrefix(run.Path, "tmp/evals/") {
		t.Fatalf("run path = %q, want configured evaluation dir", run.Path)
	}

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Evaluations.Path != "tmp/evals" || snapshot.Evaluations.Runs != 1 {
		t.Fatalf("evaluation history = %#v, want configured path with one run", snapshot.Evaluations)
	}
}

func TestSnapshotMalformedCurrentRunDoesNotHideLaterRuns(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  starts:
    title: Starts
    factors: [reliability]
    assessment: Run it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	if err := os.MkdirAll(filepath.Join(repo, ".quality", "evaluations", "0001-full-eval"), 0o755); err != nil {
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

func TestSnapshotIgnoresUnrecognizedRunFolders(t *testing.T) {
	repo := newRepo(t)
	path := writeFile(t, repo, validModel(`requirements:
  starts:
    title: Starts
    factors: [reliability]
    assessment: Run it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
`))
	if err := os.MkdirAll(filepath.Join(repo, ".quality", "evaluations", "0001-quality-eval"), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	snapshot, err := Snapshot(Options{Path: path})
	if err != nil {
		t.Fatalf("Snapshot() error = %v", err)
	}
	if snapshot.Readiness != ReadinessReadyToEvaluate {
		t.Fatalf("readiness = %q, want ready to evaluate", snapshot.Readiness)
	}
	if snapshot.Evaluations.Runs != 0 || snapshot.Evaluations.Summary.Problems != 0 {
		t.Fatalf("evaluation summary = %#v, want unrecognized folder ignored", snapshot.Evaluations)
	}
	if snapshot.Evaluations.Latest != nil {
		t.Fatalf("latest = %#v, want no recognized latest run", snapshot.Evaluations.Latest)
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
