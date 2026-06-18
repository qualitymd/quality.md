package evaluation

import (
	"os"
	"path/filepath"
	"slices"
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
	result, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
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

	second, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("second CreateRun() error = %v", err)
	}
	if second.Path != "quality/evaluations/0002-subject-quality-eval" {
		t.Fatalf("second path = %q, want next subject run", second.Path)
	}
}

func TestCreateRunRejectsRemovedModelAltitude(t *testing.T) {
	repo := testRepo(t)
	_, err := CreateRun(Options{RepoRoot: repo, Altitude: "model", Subject: "QUALITY.md"})
	if err == nil {
		t.Fatal("CreateRun() error = nil, want removed model-altitude diagnostic")
	}
	if !strings.Contains(err.Error(), "model-altitude evaluation has been removed") {
		t.Fatalf("CreateRun() error = %v, want removed model-altitude diagnostic", err)
	}
}

func TestCreateRunValidatesSubjectBeforeCreatingRunFolder(t *testing.T) {
	repo := testRepo(t)
	_, err := CreateRun(Options{RepoRoot: repo, Subject: "."})
	if err == nil {
		t.Fatal("CreateRun() error = nil, want invalid subject")
	}
	if !strings.Contains(err.Error(), "--subject") {
		t.Fatalf("CreateRun() error = %v, want subject diagnostic", err)
	}
	if _, err := os.Stat(filepath.Join(repo, "quality")); !os.IsNotExist(err) {
		t.Fatalf("quality directory stat error = %v, want not exist", err)
	}
}

func TestAddRecordStatusAndBuildReport(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
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
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	for _, want := range []string{"## Summary", "## Scope", "## Top Risks and Limitations", "## Evidence Basis", "## Next Action", "## Target Summary"} {
		if !strings.Contains(string(reportMD), want) {
			t.Fatalf("report.md missing %q:\n%s", want, reportMD)
		}
	}
	for _, want := range []string{`"scope": {`, `"evidenceBasis": [`, `"targetSummary": [`, `"recommendations": [`} {
		if !strings.Contains(string(first), want) {
			t.Fatalf("report.json missing %q:\n%s", want, first)
		}
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

func TestBuildReportRendersStructuralTargetAndEmptyRecommendations(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Narrowing: "child-quick", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	if err := os.WriteFile(filepath.Join(runPath, "design.md"), []byte("# Evaluation design\n\n## Resolved parameters\n\n- Mode: evaluate\n- Altitude: subject\n- Target file: `QUALITY.md`\n- Scope description: child target quick pass\n- Narrowing: `child-quick`\n- Effort: quick\n\n## Out of Scope\n\n- Browser tests.\n- Package-manager integration.\n"), 0o644); err != nil {
		t.Fatalf("write design.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runPath, "plan.md"), []byte("# Evaluation plan\n\n## Effort\n\nQuick pass over one child target.\n\n## Planned limitations\n\n- Does not execute the full CI matrix.\n- Defers performance benchmarks.\n\n## Deferred Areas\n\n- Release packaging.\n"), 0o644); err != nil {
		t.Fatalf("write plan.md: %v", err)
	}

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the requirement. Does not execute the full CI matrix.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child requirement reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child local rating binds."
  },
  "assessmentRecords": ["assessments/001-child-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": null,
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child aggregate rating binds the root."
  },
  "assessmentRecords": [],
  "childAnalysisRecords": ["analysis/child.json"]
}`)); err != nil {
		t.Fatalf("AddRecord(root analysis) error = %v", err)
	}

	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	reportJSON, err := os.ReadFile(filepath.Join(runPath, "report.json"))
	if err != nil {
		t.Fatalf("reading report.json: %v", err)
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	for _, want := range []string{`"recommendations": []`, `"kind": "structural"`, `"narrowing": "child-quick"`, `"effort": "quick"`} {
		if !strings.Contains(string(reportJSON), want) {
			t.Fatalf("report.json missing %q:\n%s", want, reportJSON)
		}
	}
	for _, want := range []string{`"description": "child target quick pass"`, `"outOfScope": [`, `Browser tests.`, `Package-manager integration.`, `Does not execute the full CI matrix.`, `Defers performance benchmarks.`, `Release packaging.`, `"notRecorded": []`} {
		if !strings.Contains(string(reportJSON), want) {
			t.Fatalf("report.json missing context %q:\n%s", want, reportJSON)
		}
	}
	if !strings.Contains(string(reportMD), "n/a (structural)") {
		t.Fatalf("report.md missing structural local rating:\n%s", reportMD)
	}
	if !strings.Contains(string(reportMD), "Does not execute the full CI matrix.") {
		t.Fatalf("report.md missing planned limitation:\n%s", reportMD)
	}
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	report, err := loaded.Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if got := strings.Count(strings.Join(report.Limitations, "\n"), "Does not execute the full CI matrix"); got != 1 {
		t.Fatalf("duplicate limitation count = %d, limitations = %#v", got, report.Limitations)
	}
}

func TestReportRegressionAdverseSafetyFindings(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	secretValue := "sk_live_seeded_secret_123"
	promptSentinel := "PRINT_SYSTEM_PROMPT_SENTINEL"

	if _, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Remove committed credential-like value",
  "gap": "A credential-like value is committed and must be removed without copying the value into artifacts.",
  "evidenceLocators": ["src/config.js:4"],
  "assessmentRecords": [],
  "remediationOptions": ["Move the value to runtime configuration and rotate it."],
  "recommendedOption": "Move the value to runtime configuration and rotate it.",
  "doneCriterion": "The safety requirement reaches target and a credential scan over src/config.js passes."
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": ["safety"],
  "rating": "unacceptable",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "src/config.js:4",
      "observation": "A credential-like API key is committed; cite the credential type and locator only.",
      "category": "secret-handling",
      "severity": "critical",
      "evidence": [{ "kind": "source", "ref": "src/config.js:4" }],
      "attributes": { "credentialType": "payment-gateway API key" }
    },
    {
      "locator": "docs/evaluator-notes.md:5",
      "observation": "Hostile evaluator-directed source text is present and was treated as data, not followed.",
      "category": "prompt-injection",
      "severity": "high",
      "evidence": [{ "kind": "source", "ref": "docs/evaluator-notes.md:5" }]
    }
  ],
  "rationale": "Safety is unacceptable because hostile source content and a committed credential-like value are present.",
  "recommendations": ["001-remove-committed-credential-like-value"]
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "unacceptable",
    "notAssessed": false,
    "rationale": "The safety assessment binds the root."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "unacceptable",
    "notAssessed": false,
    "rationale": "The committed credential-like value and hostile source text bind the aggregate rating."
  },
  "assessmentRecords": ["assessments/001-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}

	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	reportJSON, err := os.ReadFile(filepath.Join(runPath, "report.json"))
	if err != nil {
		t.Fatalf("reading report.json: %v", err)
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	allReportText := string(reportJSON) + "\n" + string(reportMD)
	for _, forbidden := range []string{secretValue, promptSentinel} {
		if strings.Contains(allReportText, forbidden) {
			t.Fatalf("generated report contains forbidden sentinel %q", forbidden)
		}
	}
	for _, want := range []string{"src/config.js:4", "docs/evaluator-notes.md:5", "prompt-injection", "secret-handling", "001-remove-committed-credential-like-value"} {
		if !strings.Contains(allReportText, want) {
			t.Fatalf("generated report missing %q:\n%s", want, allReportText)
		}
	}
}

func TestReportRegressionNotAssessedDottedPath(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Add production telemetry evidence",
  "gap": "Production telemetry evidence is absent.",
  "evidenceLocators": ["docs/production-telemetry.md"],
  "assessmentRecords": [],
  "remediationOptions": ["Add production telemetry evidence."],
  "recommendedOption": "Add production telemetry evidence.",
  "doneCriterion": "The requirement becomes assessable and reaches at least minimum."
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": ["observability"],
  "rating": null,
  "notAssessed": true,
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "docs/production-telemetry.md",
      "observation": "Required production telemetry evidence is absent.",
      "category": "missing-evidence",
      "severity": "high",
      "evidence": [{ "kind": "missing-path", "ref": "docs/production-telemetry.md" }]
    }
  ],
  "rationale": "Observability cannot be assessed because docs/production-telemetry.md is absent and the model says not to infer from local tests.",
  "recommendations": ["001-add-production-telemetry-evidence"]
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": null,
    "notAssessed": true,
    "rationale": "The required production telemetry evidence is missing."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": null,
    "notAssessed": true,
    "rationale": "The root cannot be assessed because docs/production-telemetry.md is absent."
  },
  "assessmentRecords": ["assessments/001-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}

	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	report, err := loaded.Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if report.Rating.Rating != nil || !report.Rating.NotAssessed {
		t.Fatalf("root rating = %#v, want not assessed with nil rating", report.Rating)
	}
	if len(report.Assessments) != 1 || report.Assessments[0].Rating != nil || !report.Assessments[0].NotAssessed {
		t.Fatalf("assessment = %#v, want not assessed with nil rating", report.Assessments)
	}
	joinedLimitations := strings.Join(report.Limitations, "\n")
	if !strings.Contains(joinedLimitations, "docs/production-telemetry.md") {
		t.Fatalf("limitations missing dotted path: %#v", report.Limitations)
	}
	for _, limitation := range report.Limitations {
		if strings.HasPrefix(limitation, "md ") || strings.HasPrefix(limitation, "md is absent") {
			t.Fatalf("limitations split dotted path incorrectly: %#v", report.Limitations)
		}
	}
}

func TestStatusRequiresRootAnalysis(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Narrowing: "child-only", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the requirement.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child requirement reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child local rating binds."
  },
  "assessmentRecords": ["assessments/001-child-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.Status()
	if status.Reportable {
		t.Fatalf("status.Reportable = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "missing-root-analysis" {
		t.Fatalf("status.Gaps = %#v, want missing-root-analysis", status.Gaps)
	}
	if _, err := BuildReport(runPath); err == nil {
		t.Fatal("BuildReport() error = nil, want missing-root-analysis")
	}
}

func TestStatusRejectsDuplicateAssessments(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	payload := []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the requirement.",
  "recommendations": []
}`)
	if _, err := AddRecord(KindAssessment, runPath, payload); err != nil {
		t.Fatalf("AddRecord(first assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, payload); err != nil {
		t.Fatalf("AddRecord(duplicate assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The first assessment reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
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
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "duplicate-assessment" {
		t.Fatalf("status.Gaps = %#v, want duplicate-assessment", status.Gaps)
	}
	if status.Gaps[0].Ref != "assessments/002-root-has-tests.json" || !strings.Contains(status.Gaps[0].Detail, "assessments/001-root-has-tests.json") {
		t.Fatalf("duplicate gap = %#v, want duplicate ref and prior detail", status.Gaps[0])
	}
	if len(status.NextActions) != 1 || status.NextActions[0].ID != "review-gaps" {
		t.Fatalf("status.NextActions = %#v, want review-gaps", status.NextActions)
	}
	if _, err := BuildReport(runPath); err == nil {
		t.Fatal("BuildReport() error = nil, want duplicate-assessment")
	}
}

func TestAssessmentSupersedingRequiresActiveAnalysisReference(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "minimum",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Only thin evidence exists.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(first assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Corrected evidence reaches target.",
  "recommendations": [],
  "supersedes": ["001-root-has-tests"]
}`)); err != nil {
		t.Fatalf("AddRecord(corrected assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "minimum",
    "notAssessed": false,
    "rationale": "The stale assessment still binds."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "minimum",
    "notAssessed": false,
    "rationale": "The stale assessment still binds."
  },
  "assessmentRecords": ["assessments/001-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(stale analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.Status()
	if status.Reportable {
		t.Fatal("status.Reportable with stale analysis = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "superseded-assessment-reference" {
		t.Fatalf("status.Gaps = %#v, want superseded-assessment-reference", status.Gaps)
	}

	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The corrected assessment reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The corrected assessment binds."
  },
  "assessmentRecords": ["assessments/002-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(corrected analysis) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after corrected analysis) error = %v", err)
	}
	if status := loaded.Status(); !status.Reportable {
		t.Fatalf("status.Reportable after corrected analysis = false, gaps = %#v", status.Gaps)
	}
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	report, err := loaded.Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if len(report.Assessments) != 2 {
		t.Fatalf("assessment count = %d, want 2", len(report.Assessments))
	}
	if report.Assessments[0].Active {
		t.Fatalf("first assessment active = true, want superseded")
	}
	if !report.Assessments[1].Active || len(report.Assessments[1].Supersedes) != 1 {
		t.Fatalf("second assessment = %#v, want active superseding assessment", report.Assessments[1])
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.Contains(string(reportMD), "- **State:** superseded") || !strings.Contains(string(reportMD), "- **State:** active") {
		t.Fatalf("report.md missing assessment states:\n%s", reportMD)
	}
}

func TestAssessmentSupersedingRejectsMissingOrDifferentRequirement(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the requirement.",
  "recommendations": [],
  "supersedes": ["999-missing-assessment"]
}`)); err != nil {
		t.Fatalf("AddRecord(missing supersedes assessment) error = %v", err)
	}
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if status := loaded.Status(); len(status.Gaps) != 2 || status.Gaps[0].Kind != "missing-superseded-assessment" {
		t.Fatalf("status.Gaps = %#v, want missing-superseded-assessment plus missing-analysis", status.Gaps)
	}

	repo = testRepo(t)
	run, err = CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(second) error = %v", err)
	}
	runPath = filepath.Join(repo, run.Path)
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "minimum",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Thin evidence exists.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(first assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has docs",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Documentation exists.",
  "recommendations": [],
  "supersedes": ["001-root-has-tests"]
}`)); err != nil {
		t.Fatalf("AddRecord(different requirement supersedes assessment) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(second) error = %v", err)
	}
	found := false
	for _, gap := range loaded.Status().Gaps {
		if gap.Kind == "invalid-assessment-supersedes" {
			found = true
		}
	}
	if !found {
		t.Fatalf("status.Gaps = %#v, want invalid-assessment-supersedes", loaded.Status().Gaps)
	}
}

func TestRecommendationSupersedingSelectsActiveNextAction(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Original recommendation",
  "gap": "The initial remediation route is too broad.",
  "evidenceLocators": ["tests/example_test.go:1"],
  "assessmentRecords": [],
  "remediationOptions": ["Do the broad fix"],
  "recommendedOption": "Do the broad fix",
  "doneCriterion": "The broad fix is complete."
}`)); err != nil {
		t.Fatalf("AddRecord(first recommendation) error = %v", err)
	}
	second, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Corrected recommendation",
  "gap": "The remediation route should be narrower.",
  "evidenceLocators": ["tests/example_test.go:1"],
  "assessmentRecords": [],
  "remediationOptions": ["Do the narrow fix"],
  "recommendedOption": "Do the narrow fix",
  "doneCriterion": "The narrow fix is complete.",
  "supersedes": ["001-original-recommendation"]
}`))
	if err != nil {
		t.Fatalf("AddRecord(second recommendation) error = %v", err)
	}
	secondRaw, err := os.ReadFile(second.Path)
	if err != nil {
		t.Fatalf("reading second recommendation: %v", err)
	}
	if !strings.Contains(string(secondRaw), "## Supersedes") || !strings.Contains(string(secondRaw), "001-original-recommendation") {
		t.Fatalf("second recommendation missing supersedes body:\n%s", secondRaw)
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
      "category": "coverage"
    }
  ],
  "rationale": "Some evidence exists but coverage is thin.",
  "recommendations": ["001-original-recommendation", "002-corrected-recommendation"]
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

	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	report, err := loaded.Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if len(report.Recommendations) != 2 {
		t.Fatalf("recommendation count = %d, want 2", len(report.Recommendations))
	}
	if report.Recommendations[0].Path != "recommendations/001-original-recommendation.md" || report.Recommendations[0].Active {
		t.Fatalf("first recommendation = %#v, want superseded original recommendation", report.Recommendations[0])
	}
	if report.Recommendations[1].Path != "recommendations/002-corrected-recommendation.md" || !report.Recommendations[1].Active {
		t.Fatalf("second recommendation = %#v, want active corrected recommendation", report.Recommendations[1])
	}
	if report.NextAction.RecommendationID != "002-corrected-recommendation" || report.NextAction.Summary != "The narrow fix is complete." {
		t.Fatalf("next action = %#v, want corrected recommendation", report.NextAction)
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.Contains(string(reportMD), "[superseded]") || !strings.Contains(string(reportMD), "[active]") {
		t.Fatalf("report.md missing recommendation states:\n%s", reportMD)
	}
}

func TestStatusRejectsMissingSupersededRecommendation(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Corrected recommendation",
  "gap": "The remediation route should be narrower.",
  "evidenceLocators": ["tests/example_test.go:1"],
  "assessmentRecords": [],
  "remediationOptions": ["Do the narrow fix"],
  "recommendedOption": "Do the narrow fix",
  "doneCriterion": "The narrow fix is complete.",
  "supersedes": ["999-missing-recommendation"]
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.Status()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	found := false
	for _, gap := range status.Gaps {
		if gap.Kind == "missing-superseded-recommendation" && gap.Ref == "999-missing-recommendation" {
			found = true
		}
	}
	if !found {
		t.Fatalf("status.Gaps = %#v, want missing-superseded-recommendation", status.Gaps)
	}
	if len(status.NextActions) != 1 || status.NextActions[0].ID != "review-gaps" {
		t.Fatalf("status.NextActions = %#v, want review-gaps", status.NextActions)
	}
}

func TestPlannedCoverageStatusGaps(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	if _, err := SetPlannedCoverage(runPath, []byte(`{
  "schemaVersion": 1,
  "assessments": [
    { "targetPath": [], "requirement": "Has tests" },
    { "targetPath": ["Child"], "requirement": "Has tests" }
  ],
  "analyses": [
    { "targetPath": [] },
    { "targetPath": ["Child"] }
  ]
}`)); err != nil {
		t.Fatalf("SetPlannedCoverage() error = %v", err)
	}
	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the root requirement.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(root assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Root",
  "targetPath": [],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The root assessment reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The root local rating binds."
  },
  "assessmentRecords": ["assessments/001-root-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(root analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.Status()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	gotKinds := gapKinds(status.Gaps)
	wantKinds := []string{"missing-planned-assessment", "missing-planned-analysis"}
	if !slices.Equal(gotKinds, wantKinds) {
		t.Fatalf("gap kinds = %#v, want %#v; gaps = %#v", gotKinds, wantKinds, status.Gaps)
	}

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "Focused tests cover the child requirement.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(child assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "target": "Child",
  "targetPath": ["Child"],
  "localRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child assessment reaches target."
  },
  "factorRatings": [],
  "aggregateRating": {
    "rating": "target",
    "notAssessed": false,
    "rationale": "The child local rating binds."
  },
  "assessmentRecords": ["assessments/002-child-has-tests.json"],
  "childAnalysisRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after complete plan) error = %v", err)
	}
	if status := loaded.Status(); !status.Reportable {
		t.Fatalf("completed planned coverage reportable = false, gaps = %#v", status.Gaps)
	}

	if _, err := AddRecord(KindAssessment, runPath, []byte(`{
  "target": "Extra",
  "targetPath": ["Extra"],
  "requirement": "Has tests",
  "factors": [],
  "rating": "target",
  "notAssessed": false,
  "criterionSource": "rating-scale",
  "findings": [],
  "rationale": "An unplanned extra assessment exists.",
  "recommendations": []
}`)); err != nil {
		t.Fatalf("AddRecord(extra assessment) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after extra record) error = %v", err)
	}
	status = loaded.Status()
	if status.Reportable {
		t.Fatal("status.Reportable with extra record = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "unexpected-assessment" {
		t.Fatalf("status.Gaps = %#v, want unexpected-assessment", status.Gaps)
	}
	if len(status.NextActions) != 1 || status.NextActions[0].ID != "review-gaps" {
		t.Fatalf("status.NextActions = %#v, want review-gaps", status.NextActions)
	}
}

func TestLimitationSentencesPreservesDottedPaths(t *testing.T) {
	got := limitationSentences("Observability cannot be assessed because docs/production-telemetry.md is absent and the model says not to infer from tests. The file is absent, so the requirement is not assessed rather than rated.")
	want := []string{
		"Observability cannot be assessed because docs/production-telemetry.md is absent and the model says not to infer from tests",
		"The file is absent, so the requirement is not assessed rather than rated",
	}
	if !slices.Equal(got, want) {
		t.Fatalf("limitationSentences() = %#v, want %#v", got, want)
	}
}

func TestAddRecordRejectsCLIOwnedFields(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
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

func gapKinds(gaps []Gap) []string {
	kinds := make([]string, 0, len(gaps))
	for _, gap := range gaps {
		kinds = append(kinds, gap.Kind)
	}
	return kinds
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
