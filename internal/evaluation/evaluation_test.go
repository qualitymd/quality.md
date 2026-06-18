package evaluation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testModel = `---
title: Test model
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: minimum
    description: Minimum.
    criterion: Barely meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  Has tests:
    assessment: Inspect tests.
---
`

func TestCreateRunUsesSharedNumberingAndSeedsLayout(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Altitude: "subject", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if result.Path != "quality/evaluations/0001-subject-quality-eval" {
		t.Fatalf("path = %q, want default numbered run", result.Path)
	}
	for _, name := range []string{"model.md", "design.md", "plan.md", "assessments", "analysis", "recommendations"} {
		if _, err := os.Stat(filepath.Join(repo, result.Path, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}

	second, err := CreateRun(Options{RepoRoot: repo, Altitude: "model", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("second CreateRun() error = %v", err)
	}
	if second.Path != "quality/evaluations/0002-model-quality-eval" {
		t.Fatalf("second path = %q, want shared numbering across altitudes", second.Path)
	}
}

func TestAddRecordStatusAndBuildReport(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Altitude: "subject", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	rec, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Fix the test gap",
  "gap": "Tests are missing.",
  "evidenceLocators": ["tests/example_test.go:1"],
  "assessmentRecords": [],
  "remediationOptions": ["Add tests"],
  "recommendedOption": "Add tests",
  "doneCriterion": "The requirement reaches target."
}`))
	if err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if !strings.HasSuffix(rec.Path, "recommendations/001-fix-the-test-gap.md") {
		t.Fatalf("recommendation path = %q, want numbered markdown", rec.Path)
	}

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "minimum",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "observation": "Only a smoke test exists.",
      "category": "coverage",
      "evidence": [{ "kind": "source", "ref": "tests/example_test.go:1" }]
    }
  ],
  "rationale": "Some evidence exists but coverage is thin.",
  "recommendations": ["001-fix-the-test-gap"]
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}

	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "minimum",
    "notAssessed": false,
    "rationale": "Thin coverage binds the local rating."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "minimum",
    "notAssessed": false,
    "rationale": "The root local rating binds."
  },
  "assessmentRecords": ["assessments/001-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.Status()
	if !status.Reportable {
		t.Fatalf("status.Reportable = false, gaps = %#v", status.Gaps)
	}

	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	first, err := os.ReadFile(filepath.Join(runPath, "report.json"))
	if err != nil {
		t.Fatalf("reading report.json: %v", err)
	}
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("second BuildReport() error = %v", err)
	}
	second, err := os.ReadFile(filepath.Join(runPath, "report.json"))
	if err != nil {
		t.Fatalf("reading second report.json: %v", err)
	}
	if string(first) != string(second) {
		t.Fatal("report.json changed across idempotent render")
	}
	if build.Rating == nil || *build.Rating != "minimum" {
		t.Fatalf("build rating = %#v, want minimum", build.Rating)
	}
}

func TestAddRecordRejectsCLIOwnedFields(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Altitude: "subject", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	_, err = AddRecord(KindAssessment, filepath.Join(repo, run.Path), []byte(`{
  "schemaVersion": 1,
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "ok",
  "recommendations": []
}`))
	if err == nil {
		t.Fatal("AddRecord() error = nil, want CLI-owned field rejection")
	}
	if _, ok := err.(*UsageError); !ok {
		t.Fatalf("error type = %T, want *UsageError (%v)", err, err)
	}
}

func testRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(testModel), 0o644); err != nil {
		t.Fatalf("write QUALITY.md: %v", err)
	}
	return repo
}
