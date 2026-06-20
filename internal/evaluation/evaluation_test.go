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
    title: 🟢 Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: 🟡 Minimum
    description: Minimum.
    criterion: Barely meets it.
  - level: unacceptable
    title: 🔴 Unacceptable
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

func TestInspectReportsCompatibilityGapsAndDiscoveredCounts(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	writeRunFile(t, runPath, "assessments/001-malformed.json", `{`)
	writeRunFile(t, runPath, "assessments/002-missing-schema.json", `{
  "targetPath": [],
  "requirement": "Has tests",
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Looks good."
  },
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": []
}`)
	writeRunFile(t, runPath, "assessments/003-future-schema.json", `{
  "schemaVersion": 99,
  "targetPath": [],
  "requirement": "Has tests",
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Looks good."
  },
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": []
}`)
	writeRunFile(t, runPath, "analysis/root.json", `{
  "schemaVersion": 1,
  "targetPath": []
}`)
	writeRunFile(t, runPath, "recommendations/001-bad.md", `# No runtime frontmatter
`)

	if _, err := Load(runPath); err == nil {
		t.Fatal("Load() error = nil, want strict load to reject malformed record")
	}
	inspected, err := Inspect(runPath)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	status := inspected.EvaluationRunStatus()
	if status.Reportable {
		t.Fatalf("status.Reportable = true, gaps = %#v", status.Gaps)
	}
	if status.Counts.AssessmentResults != 3 || status.Counts.Analyses != 1 || status.Counts.Recommendations != 1 {
		t.Fatalf("counts = %#v, want discovered record-file counts", status.Counts)
	}
	for _, kind := range []EvaluationRunGapKind{
		GapMalformedEvaluationRecord,
		GapMissingRecordSchemaVersion,
		GapUnsupportedRecordSchemaVersion,
		GapIncompleteEvaluationRecord,
	} {
		if !hasGap(status.Gaps, kind) {
			t.Fatalf("status gaps = %#v, want %s", status.Gaps, kind)
		}
	}
	if _, err := BuildReport(runPath); err == nil || !strings.Contains(err.Error(), "qualitymd evaluation status") || !strings.Contains(err.Error(), string(GapMalformedEvaluationRecord)) {
		t.Fatalf("BuildReport() error = %v, want reportability diagnostic with status command", err)
	}
	if _, err := GateReport(runPath, "target"); err == nil || !strings.Contains(err.Error(), "qualitymd evaluation status") || !strings.Contains(err.Error(), string(GapMalformedEvaluationRecord)) {
		t.Fatalf("GateReport() error = %v, want reportability diagnostic with status command", err)
	}

	records, err := ListRecords(KindAssessmentResult, runPath)
	if err != nil {
		t.Fatalf("ListRecords() error = %v", err)
	}
	if len(records.Records) != 3 || records.Records[0] != "assessments/001-malformed.json" {
		t.Fatalf("records = %#v, want discovered assessment filenames", records.Records)
	}
}

func TestListRunsIncludesReportableAndIncompatibleRuns(t *testing.T) {
	repo := testRepo(t)
	first, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(first) error = %v", err)
	}
	completeReportableRun(t, filepath.Join(repo, first.Path))
	second, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(second) error = %v", err)
	}
	writeRunFile(t, filepath.Join(repo, second.Path), "assessments/001-bad.json", `{`)

	list, err := ListRuns(repo, "", "all")
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if len(list.Runs) != 2 {
		t.Fatalf("runs = %#v, want reportable plus incompatible run", list.Runs)
	}
	if !list.Runs[0].Reportable || list.Runs[0].Gaps != 0 {
		t.Fatalf("first run = %#v, want reportable", list.Runs[0])
	}
	if list.Runs[1].Reportable || list.Runs[1].Gaps == 0 {
		t.Fatalf("second run = %#v, want incompatible incomplete run", list.Runs[1])
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
  "evidenceLocators": [
    "tests/example_test.go:1"
  ],
  "remediationOptions": [
    "Add tests"
  ],
  "recommendedOption": "Add tests",
  "doneCriterion": "The requirement reaches target.",
  "assessmentResultRecords": []
}`))
	if err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if !strings.HasSuffix(rec.Path, "recommendations/001-fix-the-test-gap.md") {
		t.Fatalf("recommendation path = %q, want numbered markdown", rec.Path)
	}

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "category": "coverage",
      "severity": "medium",
      "evidence": [
        {
          "kind": "source",
          "ref": "tests/example_test.go:1"
        }
      ],
      "observation": "Only a smoke test exists."
    },
    {
      "locator": "docs/testing.md:1",
      "category": "supporting-evidence",
      "severity": "info",
      "observation": "A testing note exists as supporting evidence."
    }
  ],
  "recommendations": [
    "001-fix-the-test-gap"
  ],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Some evidence exists but coverage is thin."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}

	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "Thin coverage binds the local rating."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "The root local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
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
	for _, want := range []string{"## Verdict", "## Scope", "## Selected Findings and Limitations", "## Evidence Basis", "## Next Action", "## Target Summary", "- **Evaluation verdict:** 🟡 Minimum", "| Test model | 🟡 Minimum | 🟡 Minimum |"} {
		if !strings.Contains(string(reportMD), want) {
			t.Fatalf("report.md missing %q:\n%s", want, reportMD)
		}
	}
	for _, want := range []string{`"scope": {`, `"evidenceBasis": [`, `"targetSummary": [`, `"recommendations": [`, `"ratingResult": {`, `"level": "minimum"`, `"severity": {`, `"level": "medium"`, `"title": "Medium"`} {
		if !strings.Contains(string(first), want) {
			t.Fatalf("report.json missing %q:\n%s", want, first)
		}
	}
	if strings.Contains(string(first), "🟡 Minimum") {
		t.Fatalf("report.json contains rating title:\n%s", first)
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
	if build.RatingResult.Kind != "rated" || build.RatingResult.Level != "minimum" {
		t.Fatalf("build ratingResult = %#v, want rated minimum", build.RatingResult)
	}
	if build.ReportSummaryMD == "" {
		t.Fatal("build.ReportSummaryMD is empty")
	}
	summaryMD, err := os.ReadFile(filepath.Join(runPath, "report-summary.md"))
	if err != nil {
		t.Fatalf("reading report-summary.md: %v", err)
	}
	for _, want := range []string{"# Quality Evaluation Summary", "| Run |", "| Scope | Full evaluation |", "| Evaluation verdict | 🟡 Minimum |", "[report.md](report.md)", "[report.json](report.json)", "## Verdict", "## Selected Findings", "**Medium**", "Only a smoke test exists.", "## Recommended Actions", "| Recommendation ID | Priority | Recommendation | Done criterion |", "`001-fix-the-test-gap`", "## Scope & Limitations"} {
		if !strings.Contains(string(summaryMD), want) {
			t.Fatalf("report-summary.md missing %q:\n%s", want, summaryMD)
		}
	}
	if strings.Contains(string(summaryMD), "| Evaluation verdict | minimum |") {
		t.Fatalf("report-summary.md rendered level id as evaluation verdict:\n%s", summaryMD)
	}
	if strings.Contains(string(summaryMD), "A testing note exists as supporting evidence.") {
		t.Fatalf("report-summary.md selected an info finding:\n%s", summaryMD)
	}
	for _, notWant := range []string{"**Root rating:**", "## Rating Summary", "## Summary", "## Top Issues", "## Recommendations", "## Limitations", "## Next Action", "| Overall rating |"} {
		if strings.Contains(string(summaryMD), notWant) {
			t.Fatalf("report-summary.md contains old summary section or label %q:\n%s", notWant, summaryMD)
		}
	}
	for _, notWant := range []string{"## Requirements", "## Findings", "### Has tests"} {
		if strings.Contains(string(summaryMD), notWant) {
			t.Fatalf("report-summary.md contains detailed section %q:\n%s", notWant, summaryMD)
		}
	}
	beforeSummary := append([]byte(nil), summaryMD...)
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("third BuildReport() error = %v", err)
	}
	afterSummary, err := os.ReadFile(filepath.Join(runPath, "report-summary.md"))
	if err != nil {
		t.Fatalf("reading third report-summary.md: %v", err)
	}
	if string(beforeSummary) != string(afterSummary) {
		t.Fatal("report-summary.md changed across idempotent render")
	}
}

func TestAddAssessmentResultRequiresCanonicalFindingSeverity(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	for _, tc := range []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "missing",
			raw: `{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "category": "coverage",
      "observation": "Only a smoke test exists."
    }
  ],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Some evidence exists but coverage is thin."
  }
}`,
			want: "must include locator, observation, category, and severity",
		},
		{
			name: "unknown",
			raw: `{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "category": "coverage",
      "severity": "blocker",
      "observation": "Only a smoke test exists."
    }
  ],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Some evidence exists but coverage is thin."
  }
}`,
			want: "severity must be one of critical, high, medium, low, or info",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AddRecord(KindAssessmentResult, runPath, []byte(tc.raw))
			if err == nil {
				t.Fatal("AddRecord() error = nil, want severity validation error")
			}
			if _, ok := err.(*UsageError); !ok {
				t.Fatalf("error type = %T, want *UsageError (%v)", err, err)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func TestLoadedAssessmentResultWithInvalidFindingSeverityIsNotReportable(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	if err := os.WriteFile(filepath.Join(runPath, "assessments", "001-root-has-tests.json"), []byte(`{
  "schemaVersion": 1,
  "targetPath": [],
  "requirement": "Has tests",
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Some evidence exists but coverage is thin."
  },
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "category": "coverage",
      "severity": "blocker",
      "observation": "Only a smoke test exists."
    }
  ],
  "recommendations": []
}`), 0o644); err != nil {
		t.Fatalf("write invalid assessment result: %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	if !hasGap(status.Gaps, GapInvalidFindingSeverity) {
		t.Fatalf("status.Gaps = %#v, want invalid-finding-severity", status.Gaps)
	}
	if _, err := BuildReport(runPath); err == nil || !strings.Contains(err.Error(), "invalid-finding-severity") {
		t.Fatalf("BuildReport() error = %v, want invalid-finding-severity", err)
	}
}

func TestLoadedAssessmentResultWithInvalidRatingResultIsNotReportable(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	if err := os.WriteFile(filepath.Join(runPath, "assessments", "001-root-has-tests.json"), []byte(`{
  "schemaVersion": 1,
  "targetPath": [],
  "requirement": "Has tests",
  "factorPaths": [],
  "ratingResult": {
    "kind": "not-assessed",
    "level": "minimum",
    "rationale": "A not-assessed result must not carry a level."
  },
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": []
}`), 0o644); err != nil {
		t.Fatalf("write invalid assessment result: %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	if !hasGap(status.Gaps, GapInvalidRatingResult) {
		t.Fatalf("status.Gaps = %#v, want invalid-rating-result", status.Gaps)
	}
	if _, err := BuildReport(runPath); err == nil || !strings.Contains(err.Error(), "invalid-rating-result") {
		t.Fatalf("BuildReport() error = %v, want invalid-rating-result", err)
	}
}

func hasGap(gaps []EvaluationRunGap, kind EvaluationRunGapKind) bool {
	for _, gap := range gaps {
		if gap.Kind == kind {
			return true
		}
	}
	return false
}

func TestBuildReportUsesModelTitlesForHumanTargetAndFactorLabels(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Root Quality Model
ratingScale:
  - level: target
    title: Good
    description: Good.
    criterion: Meets it.
  - level: unacceptable
    title: Bad
    description: Bad.
    criterion: Does not meet it.
factors:
  reliability:
    title: Operational reliability
    description: Reliable operation.
targets:
  api-service:
    title: API Service
    factors:
      automation-compatibility:
        title: Automation compatibility
        description: Automation-safe behavior.
    requirements:
      Handles API requests:
        factors:
          - automation-compatibility
          - reliability
        assessment: Inspect API handling.
---
`)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [
    "api-service"
  ],
  "requirement": "Handles API requests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [
    [
      "automation-compatibility"
    ],
    [
      "reliability"
    ]
  ],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "The assessed requirement reaches the target level."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [
    "api-service"
  ],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local requirement reaches the target level."
  },
  "factorRatingResults": [
    {
      "factorPath": [
        "automation-compatibility"
      ],
      "ratingResult": {
        "kind": "rated",
        "level": "target",
        "rationale": "The automation factor reaches the target level."
      }
    },
    {
      "factorPath": [
        "reliability"
      ],
      "ratingResult": {
        "kind": "rated",
        "level": "target",
        "rationale": "The inherited reliability factor reaches the target level."
      }
    }
  ],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local rating binds the aggregate."
  },
  "assessmentResultRecords": [
    "assessments/001-api-service-handles-api-requests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [
    "analysis/api-service.json"
  ],
  "localRatingResult": null,
  "factorRatingResults": [
    {
      "factorPath": [
        "reliability"
      ],
      "ratingResult": {
        "kind": "rated",
        "level": "target",
        "rationale": "The root reliability factor reaches the target level."
      }
    }
  ],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child aggregate rating binds the root."
  },
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(root analysis) error = %v", err)
	}

	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	summaryMD, err := os.ReadFile(filepath.Join(runPath, "report-summary.md"))
	if err != nil {
		t.Fatalf("reading report-summary.md: %v", err)
	}
	reportJSON, err := os.ReadFile(filepath.Join(runPath, "report.json"))
	if err != nil {
		t.Fatalf("reading report.json: %v", err)
	}

	for _, want := range []string{
		"- **Subject:** Root Quality Model",
		"| Root Quality Model | n/a (structural) | Good |",
		"| API Service | Good | Good |",
		"### API Service",
		"- **Path:** api-service",
		"- **Factor Automation compatibility:** Good",
		"- **Factor Operational reliability:** Good",
		"- **Target:** API Service",
	} {
		if !strings.Contains(string(reportMD), want) {
			t.Fatalf("report.md missing %q:\n%s", want, reportMD)
		}
	}
	for _, notWant := range []string{
		"| api-service |",
		"### api-service",
		"- **Factor automation-compatibility:**",
		"- **Factor reliability:**",
		"- **Target:** api-service",
	} {
		if strings.Contains(string(reportMD), notWant) {
			t.Fatalf("report.md contains raw identifier label %q:\n%s", notWant, reportMD)
		}
	}
	for _, want := range []string{
		"| Subject | Root Quality Model |",
		"| Root Quality Model | n/a (structural) | Good |",
		"| API Service | Good | Good |",
	} {
		if !strings.Contains(string(summaryMD), want) {
			t.Fatalf("report-summary.md missing %q:\n%s", want, summaryMD)
		}
	}
	for _, want := range []string{`"targetPath": [`, `"api-service"`, `"factorPath": [`, `"automation-compatibility"`, `"reliability"`, `"ratingResult": {`, `"level": "target"`} {
		if !strings.Contains(string(reportJSON), want) {
			t.Fatalf("report.json missing stable identifier %q:\n%s", want, reportJSON)
		}
	}
	for _, notWant := range []string{"Root Quality Model", "API Service", "Automation compatibility", "Operational reliability", "Good"} {
		if strings.Contains(string(reportJSON), notWant) {
			t.Fatalf("report.json contains human display label %q:\n%s", notWant, reportJSON)
		}
	}
}

func TestDisplayRatingUsesTitleAndPreservesNonRatingStates(t *testing.T) {
	result := RatingResult{Kind: RatingResultRated, Level: "minimum", Rationale: "Meets the floor."}
	if got := displayRatingResult(result, map[string]string{"minimum": "🟡 Minimum"}); got != "🟡 Minimum" {
		t.Fatalf("displayRatingResult() = %q, want title", got)
	}
	if got := displayRatingResult(result, map[string]string{}); got != "minimum" {
		t.Fatalf("displayRatingResult() = %q, want level fallback", got)
	}
	if got := displayRatingResult(RatingResult{Kind: RatingResultNotAssessed, Rationale: "No evidence."}, map[string]string{"minimum": "🟡 Minimum"}); got != "not assessed" {
		t.Fatalf("displayRatingResult(not assessed) = %q", got)
	}
	if got := displayLocalRatingState(localRatingStateFromResult(nil), map[string]string{"minimum": "🟡 Minimum"}); got != "n/a (structural)" {
		t.Fatalf("displayLocalRatingState(structural) = %q", got)
	}
}

func TestBuildReportRendersStructuralTargetAndEmptyRecommendations(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Narrowing: "child-quick", Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	if err := os.WriteFile(filepath.Join(runPath, "design.md"), []byte("# Evaluation design\n\n## Resolved parameters\n\n- Mode: evaluate\n- Altitude: subject\n- Target file: `QUALITY.md`\n- Scope description: child target quick pass\n- Narrowing: `child-quick`\n- Rigor: quick\n\n## Out of Scope\n\n- Browser tests.\n- Package-manager integration.\n"), 0o644); err != nil {
		t.Fatalf("write design.md: %v", err)
	}
	if err := os.WriteFile(filepath.Join(runPath, "plan.md"), []byte("# Evaluation plan\n\n## Rigor\n\nQuick pass over one child target.\n\n## Planned limitations\n\n- Does not execute the full CI matrix.\n- Defers performance benchmarks.\n\n## Deferred Areas\n\n- Release packaging.\n"), 0o644); err != nil {
		t.Fatalf("write plan.md: %v", err)
	}

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the requirement. Does not execute the full CI matrix."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child requirement reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-child-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [
    "analysis/child.json"
  ],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child aggregate rating binds the root."
  },
  "assessmentResultRecords": []
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
	for _, want := range []string{`"recommendations": []`, `"structural": true`, `"localRating": {`, `"kind": "structural"`, `"narrowing": "child-quick"`, `"rigor": "quick"`} {
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
	if len(report.Targets) == 0 || report.Targets[0].LocalRating.Kind != LocalRatingStructural {
		t.Fatalf("root local rating = %#v, want structural", report.Targets)
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
  "evidenceLocators": [
    "src/config.js:4"
  ],
  "remediationOptions": [
    "Move the value to runtime configuration and rotate it."
  ],
  "recommendedOption": "Move the value to runtime configuration and rotate it.",
  "doneCriterion": "The safety requirement reaches target and a credential scan over src/config.js passes.",
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "src/config.js:4",
      "category": "secret-handling",
      "severity": "critical",
      "evidence": [
        {
          "kind": "source",
          "ref": "src/config.js:4"
        }
      ],
      "attributes": {
        "credentialType": "payment-gateway API key"
      },
      "observation": "A credential-like API key is committed; cite the credential type and locator only."
    },
    {
      "locator": "docs/evaluator-notes.md:5",
      "category": "prompt-injection",
      "severity": "high",
      "evidence": [
        {
          "kind": "source",
          "ref": "docs/evaluator-notes.md:5"
        }
      ],
      "observation": "Hostile evaluator-directed source text is present and was treated as data, not followed."
    }
  ],
  "recommendations": [
    "001-remove-committed-credential-like-value"
  ],
  "factorPaths": [
    [
      "safety"
    ]
  ],
  "ratingResult": {
    "kind": "rated",
    "level": "unacceptable",
    "rationale": "Safety is unacceptable because hostile source content and a committed credential-like value are present."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "unacceptable",
      "rationale": "The safety assessment result binds the root."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "unacceptable",
      "rationale": "The committed credential-like value and hostile source text bind the aggregate rating."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
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
  "evidenceLocators": [
    "docs/production-telemetry.md"
  ],
  "remediationOptions": [
    "Add production telemetry evidence."
  ],
  "recommendedOption": "Add production telemetry evidence.",
  "doneCriterion": "The requirement becomes assessable and reaches at least minimum.",
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "docs/production-telemetry.md",
      "category": "missing-evidence",
      "severity": "high",
      "evidence": [
        {
          "kind": "missing-path",
          "ref": "docs/production-telemetry.md"
        }
      ],
      "observation": "Required production telemetry evidence is absent."
    }
  ],
  "recommendations": [
    "001-add-production-telemetry-evidence"
  ],
  "factorPaths": [
    [
      "observability"
    ]
  ],
  "ratingResult": {
    "kind": "not-assessed",
    "rationale": "Observability cannot be assessed because docs/production-telemetry.md is absent and the model says not to infer from local tests."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "not-assessed",
      "rationale": "The required production telemetry evidence is missing."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "not-assessed",
      "rationale": "The root cannot be assessed because docs/production-telemetry.md is absent."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
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
	if report.RatingResult.Kind != "not-assessed" || report.RatingResult.Level != "" {
		t.Fatalf("root ratingResult = %#v, want not assessed with empty level", report.RatingResult)
	}
	if len(report.Targets) != 1 || report.Targets[0].LocalRating.Kind != LocalRatingNotAssessed {
		t.Fatalf("local rating = %#v, want not assessed", report.Targets)
	}
	if len(report.AssessmentResults) != 1 || report.AssessmentResults[0].RatingResult.Kind != "not-assessed" || report.AssessmentResults[0].RatingResult.Level != "" {
		t.Fatalf("assessment result = %#v, want not assessed with empty level", report.AssessmentResults)
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

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the requirement."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child requirement reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-child-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
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

func TestStatusRejectsDuplicateAssessmentResults(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	payload := []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the requirement."
  }
}`)
	if _, err := AddRecord(KindAssessmentResult, runPath, payload); err != nil {
		t.Fatalf("AddRecord(first assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, payload); err != nil {
		t.Fatalf("AddRecord(duplicate assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The first assessment result reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The root local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "duplicate-assessment-result" {
		t.Fatalf("status.Gaps = %#v, want duplicate-assessment-result", status.Gaps)
	}
	if status.Gaps[0].Ref != "assessments/002-root-has-tests.json" || !strings.Contains(status.Gaps[0].Detail, "assessments/001-root-has-tests.json") {
		t.Fatalf("duplicate gap = %#v, want duplicate ref and prior detail", status.Gaps[0])
	}
	if len(status.NextActions) != 1 || status.NextActions[0].ID != "review-gaps" {
		t.Fatalf("status.NextActions = %#v, want review-gaps", status.NextActions)
	}
	if _, err := BuildReport(runPath); err == nil {
		t.Fatal("BuildReport() error = nil, want duplicate-assessment-result")
	}
}

func TestAssessmentResultSupersedingRequiresActiveAnalysisReference(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Only thin evidence exists."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(first assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "supersedes": [
    "001-root-has-tests"
  ],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Corrected evidence reaches target."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(corrected assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "The stale assessment result still binds."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "The stale assessment result still binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(stale analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable with stale analysis = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "superseded-assessment-result-reference" {
		t.Fatalf("status.Gaps = %#v, want superseded-assessment-result-reference", status.Gaps)
	}

	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The corrected assessment result reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The corrected assessment result binds."
  },
  "assessmentResultRecords": [
    "assessments/002-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(corrected analysis) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after corrected analysis) error = %v", err)
	}
	if status := loaded.EvaluationRunStatus(); !status.Reportable {
		t.Fatalf("status.Reportable after corrected analysis = false, gaps = %#v", status.Gaps)
	}
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	report, err := loaded.Report()
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if len(report.AssessmentResults) != 2 {
		t.Fatalf("assessment result count = %d, want 2", len(report.AssessmentResults))
	}
	if report.AssessmentResults[0].Active {
		t.Fatalf("first assessment result active = true, want superseded")
	}
	if !report.AssessmentResults[1].Active || len(report.AssessmentResults[1].Supersedes) != 1 {
		t.Fatalf("second assessment result = %#v, want active superseding assessment result", report.AssessmentResults[1])
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.Contains(string(reportMD), "- **State:** superseded") || !strings.Contains(string(reportMD), "- **State:** active") {
		t.Fatalf("report.md missing assessment result states:\n%s", reportMD)
	}
}

func TestAssessmentResultSupersedingRejectsMissingOrDifferentRequirement(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "supersedes": [
    "999-missing-assessment-result"
  ],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the requirement."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(missing supersedes assessment result) error = %v", err)
	}
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if status := loaded.EvaluationRunStatus(); len(status.Gaps) != 2 || status.Gaps[0].Kind != "missing-superseded-assessment-result" {
		t.Fatalf("status.Gaps = %#v, want missing-superseded-assessment-result plus missing-analysis", status.Gaps)
	}

	repo = testRepo(t)
	run, err = CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun(second) error = %v", err)
	}
	runPath = filepath.Join(repo, run.Path)
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Thin evidence exists."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(first assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has docs",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "supersedes": [
    "001-root-has-tests"
  ],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Documentation exists."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(different requirement supersedes assessment result) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(second) error = %v", err)
	}
	found := false
	for _, gap := range loaded.EvaluationRunStatus().Gaps {
		if gap.Kind == GapInvalidAssessmentResultSupersedes {
			found = true
		}
	}
	if !found {
		t.Fatalf("status.Gaps = %#v, want invalid-assessment-result-supersedes", loaded.EvaluationRunStatus().Gaps)
	}
}

func TestRecommendationSupersedingSelectsActiveNextStep(t *testing.T) {
	repo := testRepo(t)
	run, err := CreateRun(Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	if _, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Original recommendation",
  "gap": "The initial remediation route is too broad.",
  "evidenceLocators": [
    "tests/example_test.go:1"
  ],
  "remediationOptions": [
    "Do the broad fix"
  ],
  "recommendedOption": "Do the broad fix",
  "doneCriterion": "The broad fix is complete.",
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(first recommendation) error = %v", err)
	}
	second, err := AddRecord(KindRecommendation, runPath, []byte(`{
  "title": "Corrected recommendation",
  "gap": "The remediation route should be narrower.",
  "evidenceLocators": [
    "tests/example_test.go:1"
  ],
  "remediationOptions": [
    "Do the narrow fix"
  ],
  "recommendedOption": "Do the narrow fix",
  "doneCriterion": "The narrow fix is complete.",
  "supersedes": [
    "001-original-recommendation"
  ],
  "assessmentResultRecords": []
}`))
	if err != nil {
		t.Fatalf("AddRecord(second recommendation) error = %v", err)
	}
	if filepath.IsAbs(second.Path) {
		t.Fatalf("second recommendation path = %q, want repository-relative receipt path", second.Path)
	}
	secondRaw, err := os.ReadFile(filepath.Join(repo, second.Path))
	if err != nil {
		t.Fatalf("reading second recommendation: %v", err)
	}
	if !strings.Contains(string(secondRaw), "## Supersedes") || !strings.Contains(string(secondRaw), "001-original-recommendation") {
		t.Fatalf("second recommendation missing supersedes body:\n%s", secondRaw)
	}

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [
    {
      "locator": "tests/example_test.go:1",
      "category": "coverage",
      "severity": "low",
      "observation": "Only a smoke test exists."
    }
  ],
  "recommendations": [
    "001-original-recommendation",
    "002-corrected-recommendation"
  ],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "minimum",
    "rationale": "Some evidence exists but coverage is thin."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "Thin coverage binds the local rating."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "minimum",
      "rationale": "The root local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
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
	if report.Recommendations[0].State != RecordLifecycleSuperseded {
		t.Fatalf("first recommendation state = %q, want superseded", report.Recommendations[0].State)
	}
	if report.Recommendations[1].Path != "recommendations/002-corrected-recommendation.md" || !report.Recommendations[1].Active {
		t.Fatalf("second recommendation = %#v, want active corrected recommendation", report.Recommendations[1])
	}
	if report.Recommendations[1].State != RecordLifecycleActive {
		t.Fatalf("second recommendation state = %q, want active", report.Recommendations[1].State)
	}
	if report.NextStep.RecommendationID != "002-corrected-recommendation" || report.NextStep.Summary != "The narrow fix is complete." {
		t.Fatalf("next step = %#v, want corrected recommendation", report.NextStep)
	}
	reportMD, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.Contains(string(reportMD), "[superseded]") || !strings.Contains(string(reportMD), "[active]") {
		t.Fatalf("report.md missing recommendation states:\n%s", reportMD)
	}
	summaryMD, err := os.ReadFile(filepath.Join(runPath, "report-summary.md"))
	if err != nil {
		t.Fatalf("reading report-summary.md: %v", err)
	}
	if strings.Contains(string(summaryMD), "001-original-recommendation") {
		t.Fatalf("report-summary.md links superseded recommendation:\n%s", summaryMD)
	}
	if !strings.Contains(string(summaryMD), "`002-corrected-recommendation`") || !strings.Contains(string(summaryMD), "[Corrected Recommendation](recommendations/002-corrected-recommendation.md)") {
		t.Fatalf("report-summary.md missing active recommendation:\n%s", summaryMD)
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
  "evidenceLocators": [
    "tests/example_test.go:1"
  ],
  "remediationOptions": [
    "Do the narrow fix"
  ],
  "recommendedOption": "Do the narrow fix",
  "doneCriterion": "The narrow fix is complete.",
  "supersedes": [
    "999-missing-recommendation"
  ],
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(recommendation) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	found := false
	for _, gap := range status.Gaps {
		if gap.Kind == GapMissingSupersededRecommendation && gap.Ref == "999-missing-recommendation" {
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
	if err := os.WriteFile(filepath.Join(runPath, "plan.md"), []byte(`---
coverage:
  assessmentResults:
    - targetPath: []
      requirement: Has tests
    - targetPath: [Child]
      requirement: Has tests
  analyses:
    - targetPath: []
    - targetPath: [Child]
---
# Evaluation plan
`), 0o644); err != nil {
		t.Fatalf("write plan.md: %v", err)
	}
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the root requirement."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(root assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The root assessment result reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The root local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(root analysis) error = %v", err)
	}

	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	status := loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable = true, want false")
	}
	gotKinds := gapKinds(status.Gaps)
	wantKinds := []string{"missing-planned-assessment-result", "missing-planned-analysis"}
	if !slices.Equal(gotKinds, wantKinds) {
		t.Fatalf("gap kinds = %#v, want %#v; gaps = %#v", gotKinds, wantKinds, status.Gaps)
	}

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Focused tests cover the child requirement."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(child assessment result) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [
    "Child"
  ],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child assessment result reaches target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The child local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/002-child-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(child analysis) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after complete plan) error = %v", err)
	}
	if status := loaded.EvaluationRunStatus(); !status.Reportable {
		t.Fatalf("completed planned coverage reportable = false, gaps = %#v", status.Gaps)
	}

	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [
    "Extra"
  ],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "An unplanned extra assessment result exists."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(extra assessment result) error = %v", err)
	}
	loaded, err = Load(runPath)
	if err != nil {
		t.Fatalf("Load(after extra record) error = %v", err)
	}
	status = loaded.EvaluationRunStatus()
	if status.Reportable {
		t.Fatal("status.Reportable with extra record = true, want false")
	}
	if len(status.Gaps) != 1 || status.Gaps[0].Kind != "unexpected-assessment-result" {
		t.Fatalf("status.Gaps = %#v, want unexpected-assessment-result", status.Gaps)
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
	_, err = AddRecord(KindAssessmentResult, filepath.Join(repo, run.Path), []byte(`{
  "schemaVersion": 1,
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "ok"
  }
}`))
	if err == nil {
		t.Fatal("AddRecord() error = nil, want CLI-owned field rejection")
	}
	if _, ok := err.(*UsageError); !ok {
		t.Fatalf("error type = %T, want *UsageError (%v)", err, err)
	}
}

func TestSparrowExampleReportFixtureIsGenerated(t *testing.T) {
	fixture := filepath.Join("..", "..", "specs", "skills", "quality-skill", "examples", "0001-subject-quality-eval")
	tempRun := filepath.Join(t.TempDir(), "0001-subject-quality-eval")
	copyDir(t, fixture, tempRun)

	loaded, err := Load(tempRun)
	if err != nil {
		t.Fatalf("Load(fixture copy) error = %v", err)
	}
	if status := loaded.EvaluationRunStatus(); !status.Reportable {
		t.Fatalf("fixture status.Reportable = false, gaps = %#v", status.Gaps)
	}
	if _, err := BuildReport(tempRun); err != nil {
		t.Fatalf("BuildReport(fixture copy) error = %v", err)
	}

	for _, name := range []string{"report-summary.md", "report.md", "report.json"} {
		want, err := os.ReadFile(filepath.Join(fixture, name))
		if err != nil {
			t.Fatalf("read fixture %s: %v", name, err)
		}
		got, err := os.ReadFile(filepath.Join(tempRun, name))
		if err != nil {
			t.Fatalf("read generated %s: %v", name, err)
		}
		if string(got) != string(want) {
			t.Fatalf("generated %s differs from checked-in fixture", name)
		}
	}
}

func copyDir(t *testing.T, src, dst string) {
	t.Helper()
	entries, err := os.ReadDir(src)
	if err != nil {
		t.Fatalf("read dir %s: %v", src, err)
	}
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dst, err)
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			copyDir(t, srcPath, dstPath)
			continue
		}
		data, err := os.ReadFile(srcPath)
		if err != nil {
			t.Fatalf("read file %s: %v", srcPath, err)
		}
		if err := os.WriteFile(dstPath, data, 0o644); err != nil {
			t.Fatalf("write file %s: %v", dstPath, err)
		}
	}
}

func gapKinds(gaps []EvaluationRunGap) []string {
	kinds := make([]string, 0, len(gaps))
	for _, gap := range gaps {
		kinds = append(kinds, string(gap.Kind))
	}
	return kinds
}

func testRepo(t *testing.T) string {
	t.Helper()
	return testRepoWithModel(t, testModel)
}

func writeRunFile(t *testing.T, runPath, name, content string) {
	t.Helper()
	path := filepath.Join(runPath, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}

func completeReportableRun(t *testing.T, runPath string) {
	t.Helper()
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(`{
  "targetPath": [],
  "requirement": "Has tests",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {
    "kind": "rated",
    "level": "target",
    "rationale": "Tests cover the requirement."
  }
}`)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "targetPath": [],
  "childAnalysisRecords": [],
  "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local rating is target."
  },
  "factorRatingResults": [],
  "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local rating binds."
  },
  "assessmentResultRecords": [
    "assessments/001-root-has-tests.json"
  ]
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}
}

func testRepoWithModel(t *testing.T, content string) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write QUALITY.md: %v", err)
	}
	return repo
}
