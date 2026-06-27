package evaluation

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/model"
)

const testModel = `---
title: Test model
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: Minimum
    description: Minimum.
    criterion: Barely meets it.
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
---
`

func TestCreateRunSeedsEvaluationLayoutOnly(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if result.Path != ".quality/evaluations/0001-full-eval" {
		t.Fatalf("path = %q, want default numbered run", result.Path)
	}
	for _, name := range []string{ModelSnapshotFile, "data"} {
		if _, err := os.Stat(filepath.Join(repo, result.Path, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	for _, name := range []string{"model.md", "design.md", "plan.md", "assessments", "analysis", "recommendations", "debug-log.md"} {
		if _, err := os.Stat(filepath.Join(repo, result.Path, name)); !os.IsNotExist(err) {
			t.Fatalf("%s should not be seeded for evaluation runs: %v", name, err)
		}
	}
	if got := result.NextActions[0].Command; strings.Contains(got, "--file") || !strings.Contains(got, "< payloads.json") {
		t.Fatalf("next action = %q, want stdin data set command", got)
	}
}

func TestCreateRunUsesModelRelativeWorkspace(t *testing.T) {
	repo := testRepoWithModelAt(t, "packages/api/QUALITY.md", testModel)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "packages/api/QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if result.Path != ".quality/evaluations/0001-full-eval" {
		t.Fatalf("path = %q, want model-relative default numbered run", result.Path)
	}
	runPath := filepath.Join(repo, "packages", "api", result.Path)
	for _, name := range []string{ModelSnapshotFile, "data"} {
		if _, err := os.Stat(filepath.Join(runPath, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repo, ".quality", "evaluations", "0001-full-eval")); !os.IsNotExist(err) {
		t.Fatalf("repo-root evaluation run should not be created: %v", err)
	}
	want := "qualitymd evaluation data set --model packages/api/QUALITY.md .quality/evaluations/0001-full-eval < payloads.json"
	if got := result.NextActions[0].Command; got != want {
		t.Fatalf("next action = %q, want %q", got, want)
	}
}

func TestCreateRunUsesScopePathNaming(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Narrowing: "security-reliability-latency"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if result.Path != ".quality/evaluations/0001-security-reliability-latency-eval" {
		t.Fatalf("path = %q, want numbered scope-path run", result.Path)
	}
	if strings.Contains(filepath.Base(result.Path), "quality") {
		t.Fatalf("path = %q, want no quality segment in new run name", result.Path)
	}
}

func TestCreateRunRejectsReservedQualityScopeSegment(t *testing.T) {
	repo := testRepo(t)
	_, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Narrowing: "format-quality"})
	if err == nil || !strings.Contains(err.Error(), "reserved segment") {
		t.Fatalf("CreateRun() error = %v, want reserved segment diagnostic", err)
	}
}

func TestRunNameRecognitionNormalizesCurrentNarrowing(t *testing.T) {
	tests := []struct {
		name          string
		wantNumber    int
		wantNarrowing string
		wantOK        bool
	}{
		{name: "0007-full-eval", wantNumber: 7, wantNarrowing: "", wantOK: true},
		{name: "0007-security-eval", wantNumber: 7, wantNarrowing: "security", wantOK: true},
		{name: "0007-security-network-eval", wantNumber: 7, wantNarrowing: "security-network", wantOK: true},
		{name: "0007-security-reliability-eval", wantNumber: 7, wantNarrowing: "security-reliability", wantOK: true},
		{name: "0007-security-reliability-latency-eval", wantNumber: 7, wantNarrowing: "security-reliability-latency", wantOK: true},
		{name: "0007-quality-format-eval", wantOK: false},
		{name: "0006-quality-eval", wantOK: false},
		{name: "0005-subject-quality-eval", wantOK: false},
		{name: "0005-model-quality-eval", wantOK: false},
		{name: "0005-subject-security-network-quality-eval", wantOK: false},
		{name: "0005-model-security-reliability-quality-eval", wantOK: false},
		{name: "0005-security-network-quality-eval", wantOK: false},
		{name: "0005-quality-report", wantOK: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := parseRunName(tc.name)
			if ok != tc.wantOK {
				t.Fatalf("parseRunName() ok = %v, want %v", ok, tc.wantOK)
			}
			if !ok {
				return
			}
			if got.number != tc.wantNumber || got.narrowing != tc.wantNarrowing {
				t.Fatalf("parseRunName() = %#v, want number %d narrowing %q", got, tc.wantNumber, tc.wantNarrowing)
			}
			if narrowingFromRunName(tc.name) != tc.wantNarrowing {
				t.Fatalf("narrowingFromRunName() = %q, want %q", narrowingFromRunName(tc.name), tc.wantNarrowing)
			}
		})
	}
}

func TestListRunsIgnoresUnrecognizedRunFolders(t *testing.T) {
	repo := testRepo(t)
	evalDir := filepath.Join(repo, ".quality", "evaluations")
	unrecognized := filepath.Join(evalDir, "0006-quality-eval")
	if err := os.MkdirAll(filepath.Join(unrecognized, "data"), 0o755); err != nil {
		t.Fatalf("mkdir unrecognized data: %v", err)
	}
	modelRaw, err := os.ReadFile(filepath.Join(repo, "QUALITY.md"))
	if err != nil {
		t.Fatalf("read model: %v", err)
	}
	if err := os.WriteFile(filepath.Join(unrecognized, ModelSnapshotFile), modelRaw, 0o644); err != nil {
		t.Fatalf("write unrecognized model snapshot: %v", err)
	}
	if _, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Narrowing: "security-network"}); err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runs, err := ListRuns(repo, "", "")
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if len(runs.Runs) != 1 {
		t.Fatalf("ListRuns() = %#v, want one recognized run", runs.Runs)
	}
	if runs.Runs[0].Path != ".quality/evaluations/0001-security-network-eval" || runs.Runs[0].Narrowing != "security-network" {
		t.Fatalf("current run = %#v, want scope-path narrowing", runs.Runs[0])
	}
}

func TestStatusReportsMissingEvaluationDataAndNextAction(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	run, err := Inspect(filepath.Join(repo, result.Path))
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	status := run.Status()
	if status.Reportable {
		t.Fatalf("status.Reportable = true, gaps = %#v", status.Gaps)
	}
	if len(status.Gaps) == 0 || status.Gaps[0].Kind != GapMissingEvaluationData {
		t.Fatalf("status.Gaps = %#v, want missing evaluation data", status.Gaps)
	}
	if got := status.NextActions[0].Command; strings.Contains(got, "--file") || !strings.Contains(got, "< payloads.json") {
		t.Fatalf("next action = %q, want stdin data set command", got)
	}
}

func TestSetDataAndBuildEvaluationReport(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	if _, err := SetData(runPath, batchPayloads(completeRootEvaluationPayloads()...), DataSetOptions{}); err != nil {
		t.Fatalf("SetData(complete root batch) error = %v", err)
	}
	run, err := Inspect(runPath)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if status := run.Status(); !status.Reportable || status.Data.Artifacts != 3 {
		t.Fatalf("status = %#v, want reportable with three data artifacts", status)
	}
	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	if build.ReportMD == "" || build.EvaluationOutputResult == "" {
		t.Fatalf("BuildReport() = %#v, want evaluation report and output result", build)
	}
	for _, name := range []string{"report-summary.md", "report.json"} {
		if _, err := os.Stat(filepath.Join(runPath, name)); !os.IsNotExist(err) {
			t.Fatalf("%s should not be generated by evaluation report build: %v", name, err)
		}
	}
	report, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.HasPrefix(string(report), "# Area: Test model\n\nArea: [Test model](report.md)") {
		t.Fatalf("report.md = %s, want kind-prefixed title before Area trail", report)
	}
	if !strings.Contains(string(report), "| 🔵 Target | 🔵 Target | 🟢 High / 🟢 High | [area-analysis-result.json](data/areas/root/area-analysis-result.json) |") {
		t.Fatalf("report.md = %s, want target rating title", report)
	}
	if !strings.Contains(string(report), "## Legend\n\n- `—` - not applicable or not recorded.") {
		t.Fatalf("report.md = %s, want empty-cell legend", report)
	}
	if !strings.Contains(string(report), "Area: [Test model](report.md)") {
		t.Fatalf("report.md = %s, want Area trail", report)
	}
	if !strings.Contains(string(report), "| Overall Rating | Local Rating | Confidence | Data |") {
		t.Fatalf("report.md = %s, want confidence display titles", report)
	}
	if strings.Contains(string(report), "Overall Rating | target") {
		t.Fatalf("report.md = %s, want rating title instead of level ID", report)
	}
	if build.RatingResult.Level != "target" {
		t.Fatalf("BuildReport().RatingResult.Level = %q, want stable rating level ID", build.RatingResult.Level)
	}
	outputRaw, err := os.ReadFile(filepath.Join(runPath, "data", "evaluation-output-result.json"))
	if err != nil {
		t.Fatalf("reading EvaluationOutputResult: %v", err)
	}
	var output map[string]any
	if err := json.Unmarshal(outputRaw, &output); err != nil {
		t.Fatalf("EvaluationOutputResult JSON error: %v", err)
	}
	if output["kind"] != string(DataKindEvaluationOutput) {
		t.Fatalf("EvaluationOutputResult kind = %v", output["kind"])
	}
	if output["kind"] == dataKindTitle(DataKindEvaluationOutput) {
		t.Fatalf("EvaluationOutputResult kind = %v, want stable data kind", output["kind"])
	}
}

func batchPayloads(payloads ...string) []byte {
	return []byte("[\n" + strings.Join(payloads, ",\n") + "\n]")
}

func testFindingCore(id, typ, severity, confidence, statement string) map[string]any {
	return map[string]any{
		"id":         id,
		"type":       typ,
		"severity":   severity,
		"confidence": confidence,
		"statement":  statement,
		"condition":  statement + " condition.",
		"criteria": []any{map[string]any{
			"requirementId": "requirement:root::has-tests",
			"ratingLevelId": "rating:target",
			"criterion":     "Meets it.",
			"rationale":     "The target criterion is the relevant bar.",
		}},
		"cause": map[string]any{
			"status":    "not_assessed",
			"statement": "Cause was not assessed.",
		},
		"effect": map[string]any{
			"statement":    statement + " effect.",
			"ratingEffect": "constrains target",
		},
		"evidence": []any{map[string]any{
			"sourceRef": "QUALITY.md",
			"statement": "Evidence is available.",
		}},
	}
}

func testFindingJSON(finding map[string]any) string {
	raw, err := json.Marshal(finding)
	if err != nil {
		panic(err)
	}
	return string(raw)
}

func testRequirementFindingJSON(overrides map[string]any) string {
	finding := testFindingCore("gap-1", "gap", "high", "high", "Gap")
	for k, v := range overrides {
		finding[k] = v
	}
	return testFindingJSON(finding)
}

func testAreaFindingJSON(overrides map[string]any) string {
	finding := testFindingCore("gap-1", "gap", "high", "high", "Gap")
	finding["inputRefs"] = []any{map[string]any{
		"kind": "RequirementAssessmentResult",
		"subject": map[string]any{
			"requirementId": "requirement:root::has-tests",
		},
	}}
	for k, v := range overrides {
		finding[k] = v
	}
	return testFindingJSON(finding)
}

func TestEvaluationReportNavigationHeadersAndSubjectLinks(t *testing.T) {
	repo := testRepoWithModel(t, navigationReportModel)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	if _, err := SetData(runPath, batchPayloads(navigationReportPayloads()...), DataSetOptions{}); err != nil {
		t.Fatalf("SetData(navigation batch) error = %v", err)
	}
	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	rootReport := readReport(t, runPath, "report.md")
	assertContains(t, rootReport, "# Area: Navigation model\n\nArea: [Navigation model](report.md)")
	assertContains(t, rootReport, "Area: [Navigation model](report.md)")
	assertNotContains(t, rootReport, "Path: `/`")
	assertContains(t, rootReport, "| Overall Rating | Local Rating | Confidence | Data |")
	assertContains(t, rootReport, "## Findings\n\n| ID | Statement | Type | Severity | Confidence | Effect | Cause |")
	assertContains(t, rootReport, "| `root-risk` | Root reliability is the highest concern. | ⚠️ Risk | 🔴 High | 🟢 High | Reliability constrains the root aggregate rating. | Verified: Reliability analysis is the binding roll-up driver. |")
	assertContains(t, rootReport, "| `root-note` | Root note is contextual. | ℹ️ Note | 🔵 Low | 🔵 Medium | The note preserves context without driving the rating. | Not Assessed: Cause was not assessed. |")
	assertContains(t, rootReport, "#### root-risk Root reliability is the highest concern.\n\n##### Condition")
	assertContains(t, rootReport, "##### Criteria")
	assertContains(t, rootReport, "##### Cause")
	assertContains(t, rootReport, "##### Effect")
	assertContains(t, rootReport, "##### Evidence")
	assertContains(t, rootReport, "##### Relationships\n\n[Reliability Primary Driver](factors/reliability/reliability-factor.md)")
	assertContains(t, rootReport, "| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |")
	assertContains(t, rootReport, "| [Reliability](factors/reliability/reliability-factor.md) | `reliability` | 🔵 Target | 🔴 Below | [Latency](factors/reliability/factors/latency/latency-factor.md) 🔵 Target |")
	assertContains(t, rootReport, "| [Payments](areas/payments/payments-area.md) | `/payments` | 🔵 Target | — | — |")
	assertContains(t, rootReport, "| [Has tests](requirements/has-tests/has-tests-requirement.md) | 🔵 Target | ✅ Assessed | [reliability](factors/reliability/reliability-factor.md) |")
	assertContains(t, rootReport, "## Legend\n\n- `—` - not applicable or not recorded.")
	assertNotContains(t, rootReport, "Breadcrumb:")
	assertNotContains(t, rootReport, "Parent Area:")
	assertNotContains(t, rootReport, "| Details |")

	factorReport := readReport(t, runPath, "factors/reliability/reliability-factor.md")
	assertContains(t, factorReport, "# Factor: Reliability\n\nArea: [Navigation model](../../report.md)")
	assertContains(t, factorReport, "Area: [Navigation model](../../report.md)")
	assertContains(t, factorReport, "Factor: [Reliability](reliability-factor.md)")
	assertNotContains(t, factorReport, "Path: `reliability`")
	assertContains(t, factorReport, "| Overall Rating | Local Rating | Status | Confidence | Data |")
	assertContains(t, factorReport, "## Findings\n\n| ID | Statement | Type | Severity | Confidence | Effect | Cause |")
	assertContains(t, factorReport, "| `root-risk` | Root reliability is the highest concern. | ⚠️ Risk | 🔴 High | 🟢 High | Reliability constrains the root aggregate rating. | Verified: Reliability analysis is the binding roll-up driver. |")
	assertContains(t, factorReport, "#### root-risk Root reliability is the highest concern.\n\n##### Condition")
	assertNotContains(t, factorReport, "root-note")
	assertContains(t, factorReport, "| [Has tests](../../requirements/has-tests/has-tests-requirement.md) | 🔵 Target | ✅ Assessed |")
	assertContains(t, factorReport, "| [Latency](factors/latency/latency-factor.md) | `reliability/latency` | 🔵 Target | — |")
	assertNotContains(t, factorReport, "Parent Factor:")
	assertNotContains(t, factorReport, "| Details |")

	childFactorReport := readReport(t, runPath, "factors/reliability/factors/latency/latency-factor.md")
	assertContains(t, childFactorReport, "Factor: [Reliability](../../reliability-factor.md) / [Latency](latency-factor.md)")

	requirementReport := readReport(t, runPath, "requirements/has-tests/has-tests-requirement.md")
	assertContains(t, requirementReport, "# Requirement: Has tests\n\nArea: [Navigation model](../../report.md)")
	assertContains(t, requirementReport, "Area: [Navigation model](../../report.md)")
	assertContains(t, requirementReport, "Factors: [reliability](../../factors/reliability/reliability-factor.md)")
	assertContains(t, requirementReport, "# Requirement: Has tests\n\nArea: [Navigation model](../../report.md)\n\nFactors: [reliability](../../factors/reliability/reliability-factor.md)")
	assertNotContains(t, requirementReport, "Name: `has-tests`")
	assertContains(t, requirementReport, "| Rating | Assessment | Confidence | Data |")
	assertNotContains(t, requirementReport, "| Rating | Assessment | Factors | Confidence | Data |")
	assertContains(t, requirementReport, "[requirement-assessment-result.json](")
	assertContains(t, requirementReport, "[requirement-rating-result.json](")
	assertNotContains(t, requirementReport, "\nFactor:")
	assertNotContains(t, requirementReport, "Parent Area:")
	assertOnlyRootReportMD(t, runPath)

	outputRaw, err := os.ReadFile(filepath.Join(runPath, "data", "evaluation-output-result.json"))
	if err != nil {
		t.Fatalf("reading EvaluationOutputResult: %v", err)
	}
	if !strings.Contains(string(outputRaw), `"path": "factors/reliability/reliability-factor.md"`) {
		t.Fatalf("EvaluationOutputResult = %s, want subject-aware report refs", outputRaw)
	}
	if build.RatingResult.Level != "target" {
		t.Fatalf("BuildReport().RatingResult.Level = %q, want stable rating level ID", build.RatingResult.Level)
	}
}

func TestSetDataRejectsCLIOwnedOutput(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	_, err = SetData(filepath.Join(repo, result.Path), batchPayloads(`{"schemaVersion":2,"kind":"EvaluationOutputResult"}`), DataSetOptions{})
	if err == nil || !strings.Contains(err.Error(), "generated by evaluation report build") {
		t.Fatalf("SetData(EvaluationOutputResult) error = %v, want CLI-owned diagnostic", err)
	}
}

func TestSetDataRejectsInvalidEvaluationShapes(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	cases := []string{
		`{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":[],"localAnalysis":{"status":"analyzed","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"rated","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":2,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingDrivers":["finding-1"]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":["finding-1"]}`,
	}
	for _, raw := range cases {
		if _, err := SetData(runPath, batchPayloads(raw), DataSetOptions{}); err == nil {
			t.Fatalf("SetData(%s) error = nil, want validation error", raw)
		}
	}
}

func TestSetDataRequiresArrayEnvelope(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	raw := `{"schemaVersion":2,"kind":"EvaluationFrame"}`
	if _, err := SetData(runPath, []byte(raw), DataSetOptions{}); err == nil || !strings.Contains(err.Error(), "invalid JSON payload array") {
		t.Fatalf("SetData(object) error = %v, want array diagnostic", err)
	}
	if _, err := SetData(runPath, []byte(`[]`), DataSetOptions{}); err == nil || !strings.Contains(err.Error(), "at least one") {
		t.Fatalf("SetData(empty array) error = %v, want non-empty diagnostic", err)
	}
}

func TestSetDataAggregatesInvalidBatchAndWritesNothing(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	valid := `{"schemaVersion":2,"kind":"EvaluationFrame"}`
	invalidRequirement := `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::invented","status":"assessed","findings":[]}`
	invalidRating := `{"schemaVersion":2,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:invented","ratingDrivers":[]}`
	_, err = SetData(runPath, batchPayloads(valid, invalidRequirement, invalidRating), DataSetOptions{})
	if err == nil {
		t.Fatal("SetData(invalid batch) error = nil, want aggregate validation error")
	}
	for _, want := range []string{"payload[1]", "payload[2]", "RequirementAssessmentResult.requirementId", "rating:invented"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("SetData() error = %v, want %q", err, want)
		}
	}
	if _, err := os.Stat(filepath.Join(runPath, "data", "frame", "evaluation-frame.json")); !os.IsNotExist(err) {
		t.Fatalf("valid payload was written despite invalid batch: %v", err)
	}
}

func TestSetDataRejectsDuplicateBatchPaths(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	payload := `{"schemaVersion":2,"kind":"EvaluationFrame"}`
	_, err = SetData(runPath, batchPayloads(payload, payload), DataSetOptions{})
	if err == nil || !strings.Contains(err.Error(), "derive the same path") {
		t.Fatalf("SetData(duplicate batch) error = %v, want duplicate path diagnostic", err)
	}
	if _, err := os.Stat(filepath.Join(runPath, "data", "frame", "evaluation-frame.json")); !os.IsNotExist(err) {
		t.Fatalf("duplicate batch wrote data: %v", err)
	}
}

func TestSetDataDryRunReturnsBatchReceiptWithoutWriting(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	payloads := completeRootEvaluationPayloads()[:2]
	receipt, err := SetData(runPath, batchPayloads(payloads...), DataSetOptions{DryRun: true})
	if err != nil {
		t.Fatalf("SetData(dry-run batch) error = %v", err)
	}
	if receipt.Count != 2 || len(receipt.Writes) != 2 || !receipt.DryRun {
		t.Fatalf("receipt = %#v, want dry-run batch receipt", receipt)
	}
	for i, write := range receipt.Writes {
		if write.Index != i || write.Path == "" {
			t.Fatalf("receipt.Writes[%d] = %#v, want input-order path", i, write)
		}
	}
	for _, write := range receipt.Writes {
		if _, err := os.Stat(filepath.Join(runPath, filepath.FromSlash(write.Path))); !os.IsNotExist(err) {
			t.Fatalf("dry-run wrote %s: %v", write.Path, err)
		}
	}
}

func TestSetDataRejectsUnknownFieldsAndUnresolvedModelReferences(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "unknown finding field",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"title": "Wrong field"}) + `]}`,
			want: "unknown field title",
		},
		{
			name: "legacy requirement finding description field",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"description": "Legacy field"}) + `]}`,
			want: "unknown field description",
		},
		{
			name: "legacy candidate actions field",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"actions": []any{map[string]any{"description": "do it"}}}) + `]}`,
			want: "unknown field actions",
		},
		{
			name: "candidate action missing id",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"description": "do it"}}}) + `]}`,
			want: "is missing required field id",
		},
		{
			name: "candidate action missing description",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "rationale": "why"}}}) + `]}`,
			want: "is missing required field description",
		},
		{
			name: "candidate action unknown field",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "do it", "effort": "high"}}}) + `]}`,
			want: "unknown field effort",
		},
		{
			name: "candidate action non-string id",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": 5, "description": "do it"}}}) + `]}`,
			want: "must be a non-empty string",
		},
		{
			name: "candidate action non-string description",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": 5}}}) + `]}`,
			want: "must be a non-empty string",
		},
		{
			name: "requirement finding info severity",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"severity": "info"}) + `]}`,
			want: `severity = "info", want one of critical, high, medium, low`,
		},
		{
			name: "invented requirement",
			raw:  `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::invented","status":"assessed","findings":[]}`,
			want: "does not resolve in the model",
		},
		{
			name: "invented rating",
			raw:  `{"schemaVersion":2,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:invented","ratingDrivers":[]}`,
			want: "does not resolve in the model",
		},
		{
			name: "out-of-vocabulary reference kind",
			raw:  `{"schemaVersion":2,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","ratingDrivers":[{"description":"d","effect":"supports target","ratingLevelId":"rating:target","inputRefs":[{"kind":"RequirementAssessment","subject":{"requirementId":"requirement:root::has-tests"}}]}]}`,
			want: `kind = "RequirementAssessment", want one of`,
		},
		{
			name: "area finding unknown field",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"impact": "high"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "unknown field impact",
		},
		{
			name: "area finding candidate actions",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "do it"}}}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "unknown field candidateActions",
		},
		{
			name: "legacy area finding summary field",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"summary": "Legacy field"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "unknown field summary",
		},
		{
			name: "legacy area finding top-level rationale field",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"rationale": "Legacy field"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "unknown field rationale",
		},
		{
			name: "area finding type enum",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"type": "problem"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: `type = "problem", want one of`,
		},
		{
			name: "area finding info severity",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"id": "note-1", "type": "note", "severity": "info"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: `severity = "info", want one of critical, high, medium, low`,
		},
		{
			name: "area finding empty input refs",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"inputRefs": []any{}}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "inputRefs must contain at least 1 item",
		},
		{
			name: "area finding duplicate id",
			raw:  `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(nil) + `,` + testAreaFindingJSON(map[string]any{"type": "risk", "severity": "medium", "confidence": "medium", "statement": "Risk"}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "is duplicated",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SetData(runPath, batchPayloads(tc.raw), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("SetData() error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestSetDataRejectsAreaFindingFactorFromDifferentArea(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Test model
ratingScale:
  - level: target
    title: Target
    criterion: Meets it.
factors:
  reliability:
    title: Reliability
areas:
  payments:
    title: Payments
    factors:
      reliability:
        title: Payments Reliability
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
---
`)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	raw := `{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + testAreaFindingJSON(map[string]any{"factorRelationships": []any{map[string]any{"factorId": "factor:payments::reliability", "relationship": "primary-driver"}}}) + `],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`
	if _, err := SetData(runPath, batchPayloads(raw), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), "declares a different Area") {
		t.Fatalf("SetData() error = %v, want different-Area Factor diagnostic", err)
	}
}

func TestSetDataAcceptsFindingCandidateActions(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	raw := `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}, map[string]any{"id": "action-002", "description": "Cover error paths.", "rationale": "Lifts coverage."}}}) + `]}`
	if _, err := SetData(runPath, batchPayloads(raw), DataSetOptions{DryRun: true}); err != nil {
		t.Fatalf("SetData() error = %v, want valid candidate actions accepted", err)
	}
}

func TestSetDataRejectsDuplicateCandidateActionIDsWithinFinding(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	raw := `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}, map[string]any{"id": "action-001", "description": "Cover error paths."}}}) + `]}`
	if _, err := SetData(runPath, batchPayloads(raw), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), `candidateActions[1].id "action-001" is duplicated`) {
		t.Fatalf("SetData() error = %v, want duplicate candidate action ID diagnostic", err)
	}
}

func TestFindingDetailsOmitCandidateActions(t *testing.T) {
	var b strings.Builder
	finding := testFindingCore("gap-1", "gap", "high", "high", "Edge cases untested.")
	finding["candidateActions"] = []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}}
	writeEvaluationFindingDetails(&b, map[string]any{
		"findings": []any{finding},
	})
	rendered := b.String()
	if strings.Contains(rendered, "| Actions |") {
		t.Fatalf("finding details = %s, want no Actions row in the v0 report", rendered)
	}
	if strings.Contains(rendered, "Add boundary tests.") {
		t.Fatalf("finding details = %s, want candidate actions kept out of the report", rendered)
	}
	for _, want := range []string{"### gap-1 Edge cases untested.", "#### Condition", "#### Criteria", "#### Cause", "#### Effect", "#### Evidence"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("finding details = %s, want %q", rendered, want)
		}
	}
}

func TestDataSchemaAndExamplesUseContract(t *testing.T) {
	raw, err := EvaluationDataSchema("")
	if err != nil {
		t.Fatalf("EvaluationDataSchema() error = %v", err)
	}
	if !strings.Contains(string(raw), `"$defs": {`) || !strings.Contains(string(raw), `"RequirementAssessmentResult"`) {
		t.Fatalf("schema = %s, want definitions", raw)
	}
	kindSchema, err := EvaluationDataSchema(DataKindRequirementAssessment)
	if err != nil {
		t.Fatalf("EvaluationDataSchema(kind) error = %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(kindSchema, &doc); err != nil {
		t.Fatalf("Unmarshal(kind schema) error = %v", err)
	}
	if _, ok := doc["$defs"]; ok {
		t.Fatalf("kind schema = %s, want no $defs envelope", kindSchema)
	}
	if _, ok := doc["$ref"]; ok {
		t.Fatalf("kind schema = %s, want no root $ref", kindSchema)
	}
	if got, want := doc["$id"], evaluationDataSchemaID+"/"+string(DataKindRequirementAssessment); got != want {
		t.Fatalf("kind schema $id = %v, want %q", got, want)
	}
	required, ok := doc["required"].([]any)
	if !ok {
		t.Fatalf("kind schema required = %#v, want array", doc["required"])
	}
	for _, want := range []string{"schemaVersion", "kind", "requirementId", "status", "findings"} {
		if !jsonArrayContains(required, want) {
			t.Fatalf("kind schema required = %#v, want %q", required, want)
		}
	}
	props, ok := doc["properties"].(map[string]any)
	if !ok {
		t.Fatalf("kind schema properties = %#v, want object", doc["properties"])
	}
	status, ok := props["status"].(map[string]any)
	if !ok {
		t.Fatalf("kind schema status property = %#v, want object", props["status"])
	}
	statusEnum, ok := status["enum"].([]any)
	if !ok {
		t.Fatalf("kind schema status enum = %#v, want array", status["enum"])
	}
	for _, want := range []string{"assessed", "partially_assessed", "not_assessed", "blocked"} {
		if !jsonArrayContains(statusEnum, want) {
			t.Fatalf("kind schema status enum = %#v, want %q", statusEnum, want)
		}
	}
	areaSchema, err := EvaluationDataSchema(DataKindAreaAnalysis)
	if err != nil {
		t.Fatalf("EvaluationDataSchema(area kind) error = %v", err)
	}
	var areaDoc map[string]any
	if err := json.Unmarshal(areaSchema, &areaDoc); err != nil {
		t.Fatalf("Unmarshal(area schema) error = %v", err)
	}
	areaProps, ok := areaDoc["properties"].(map[string]any)
	if !ok {
		t.Fatalf("area schema properties = %#v, want object", areaDoc["properties"])
	}
	areaID, ok := areaProps["areaId"].(map[string]any)
	if !ok {
		t.Fatalf("area schema areaId property = %#v, want object", areaProps["areaId"])
	}
	if pattern, ok := areaID["pattern"].(string); !ok || !strings.Contains(pattern, "area:") {
		t.Fatalf("area schema areaId pattern = %#v, want area:<id> reference pattern", areaID["pattern"])
	}
	findings, ok := areaProps["findings"].(map[string]any)
	if !ok {
		t.Fatalf("area schema findings property = %#v, want object", areaProps["findings"])
	}
	items, ok := findings["items"].(map[string]any)
	if !ok {
		t.Fatalf("area schema findings items = %#v, want object", findings["items"])
	}
	findingProps, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("area schema finding properties = %#v, want object", items["properties"])
	}
	for field, wantValues := range map[string][]string{
		"type":       {"strength", "gap", "risk", "unknown", "note"},
		"severity":   {"critical", "high", "medium", "low"},
		"confidence": {"high", "medium", "low", "none"},
	} {
		prop, ok := findingProps[field].(map[string]any)
		if !ok {
			t.Fatalf("area schema finding %s property = %#v, want object", field, findingProps[field])
		}
		enumValues, ok := prop["enum"].([]any)
		if !ok {
			t.Fatalf("area schema finding %s enum = %#v, want array", field, prop["enum"])
		}
		for _, want := range wantValues {
			if !jsonArrayContains(enumValues, want) {
				t.Fatalf("area schema finding %s enum = %#v, want %q", field, enumValues, want)
			}
		}
	}
	inputRefs, ok := findingProps["inputRefs"].(map[string]any)
	if !ok || inputRefs["minItems"] != float64(1) {
		t.Fatalf("area schema inputRefs = %#v, want minItems 1", findingProps["inputRefs"])
	}
	relationships, ok := findingProps["factorRelationships"].(map[string]any)
	if !ok {
		t.Fatalf("area schema factorRelationships = %#v, want object", findingProps["factorRelationships"])
	}
	relationshipItems, ok := relationships["items"].(map[string]any)
	if !ok {
		t.Fatalf("area schema factorRelationships items = %#v, want object", relationships["items"])
	}
	relationshipProps, ok := relationshipItems["properties"].(map[string]any)
	if !ok {
		t.Fatalf("area schema relationship properties = %#v, want object", relationshipItems["properties"])
	}
	relationship, ok := relationshipProps["relationship"].(map[string]any)
	if !ok {
		t.Fatalf("area schema relationship property = %#v, want object", relationshipProps["relationship"])
	}
	relationshipEnum, ok := relationship["enum"].([]any)
	if !ok {
		t.Fatalf("area schema relationship enum = %#v, want array", relationship["enum"])
	}
	for _, want := range []string{"primary-driver", "contributing-driver", "evidence-limit", "offsetting-strength", "related"} {
		if !jsonArrayContains(relationshipEnum, want) {
			t.Fatalf("area schema relationship enum = %#v, want %q", relationshipEnum, want)
		}
	}
	cause, ok := findingProps["cause"].(map[string]any)
	if !ok {
		t.Fatalf("area schema cause = %#v, want object", findingProps["cause"])
	}
	causeProps, ok := cause["properties"].(map[string]any)
	if !ok {
		t.Fatalf("area schema cause properties = %#v, want object", cause["properties"])
	}
	causeStatus, ok := causeProps["status"].(map[string]any)
	if !ok {
		t.Fatalf("area schema cause status = %#v, want object", causeProps["status"])
	}
	causeStatusEnum, ok := causeStatus["enum"].([]any)
	if !ok {
		t.Fatalf("area schema cause status enum = %#v, want array", causeStatus["enum"])
	}
	for _, want := range []string{"verified", "plausible", "not_assessed", "not_applicable"} {
		if !jsonArrayContains(causeStatusEnum, want) {
			t.Fatalf("area schema cause status enum = %#v, want %q", causeStatusEnum, want)
		}
	}

	example, err := DataExample(DataKindRequirementAssessment)
	if err != nil {
		t.Fatalf("DataExample() error = %v", err)
	}
	for _, want := range []string{`"findings": [`, `"statement": "Focused test coverage is present."`, `"condition": "A focused test covers the requirement's primary path."`, `"cause": {`, `"effect": {`, `"candidateActions": [`, `"id": "action-001"`, `"description": "Add focused tests for the boundary and error paths."`} {
		if !strings.Contains(string(example), want) {
			t.Fatalf("example = %s, want %q", example, want)
		}
	}
}

func TestDataExamplesCoverAllKindsAndReferenceShapes(t *testing.T) {
	var combined strings.Builder
	for _, kind := range dataContractOrder {
		raw, err := DataExample(kind)
		if err != nil {
			t.Fatalf("DataExample(%s) error = %v", kind, err)
		}
		combined.Write(raw)
		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err != nil {
			t.Fatalf("DataExample(%s) JSON error = %v\n%s", kind, err, raw)
		}
		if got := payload["kind"]; got != string(kind) {
			t.Fatalf("DataExample(%s) kind = %v", kind, got)
		}
		if err := validateDataPayload(kind, payload); err != nil {
			t.Fatalf("DataExample(%s) validation error = %v\n%s", kind, err, raw)
		}
	}

	examples := combined.String()
	for _, want := range []string{
		`"areaId": "area:root"`,
		`"areaId": "area:operations"`,
		`"factorId": "factor:root::verification"`,
		`"factorId": "factor:root::verification/coverage"`,
		`"requirementId": "requirement:root::has-tests"`,
		`"ratingLevelId": "rating:target"`,
		`"kind": "area"`,
		`"kind": "factor"`,
		`"kind": "requirement"`,
	} {
		if !strings.Contains(examples, want) {
			t.Fatalf("generated examples missing %q", want)
		}
	}
}

func jsonArrayContains(values []any, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}

func TestEvaluationDataSchemaArtifactIsCurrent(t *testing.T) {
	generated, err := EvaluationDataSchema("")
	if err != nil {
		t.Fatalf("EvaluationDataSchema() error = %v", err)
	}
	second, err := EvaluationDataSchema("")
	if err != nil {
		t.Fatalf("EvaluationDataSchema() second error = %v", err)
	}
	if !bytes.Equal(generated, second) {
		t.Fatal("EvaluationDataSchema() is not deterministic")
	}
	artifact, err := os.ReadFile("evaluation-data.schema.json")
	if err != nil {
		t.Fatalf("ReadFile(evaluation-data.schema.json) error = %v", err)
	}
	if !bytes.Equal(generated, artifact) {
		t.Fatal("evaluation-data.schema.json is stale; regenerate with `qualitymd evaluation data schema`")
	}
}

func TestVerifyDataReportsPersistedPayloadFailures(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	badPath := filepath.Join(runPath, "data", "areas", "root", "requirements", "has-tests", "requirement-assessment-result.json")
	if err := os.MkdirAll(filepath.Dir(badPath), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	badPayload := `{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"title": "Wrong"}) + `]}`
	if err := os.WriteFile(badPath, []byte(badPayload), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	resultReceipt, err := VerifyData(runPath)
	if err != nil {
		t.Fatalf("VerifyData() error = %v", err)
	}
	if resultReceipt.Valid || len(resultReceipt.Failures) != 1 || !strings.Contains(resultReceipt.Failures[0].Reason, "unknown field title") {
		t.Fatalf("VerifyData() = %#v, want one unknown-field failure", resultReceipt)
	}
}

func TestRequirementTitleFindsFactorDeclaredRequirement(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Test model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
factors:
  reliability:
    title: Reliability
    requirements:
      retry-window:
        title: Retry window is bounded
        assessment: Inspect retry behavior.
---
`)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	run, err := Inspect(filepath.Join(repo, result.Path))
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	got := requirementTitle(run.Model, requirementID{Name: "retry-window"})
	if got != "Retry window is bounded" {
		t.Fatalf("requirementTitle() = %q, want Factor-declared Requirement title", got)
	}
}

func TestRatingTitleFallsBackToLevelID(t *testing.T) {
	spec := &model.Spec{RatingScale: []model.RatingLevel{{Level: "target"}}}
	if got := ratingTitle(spec, "target"); got != "target" {
		t.Fatalf("ratingTitle() = %q, want level ID fallback", got)
	}
}

func TestReportDisplayTitles(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"analysis status", analysisStatusTitle(string(AnalysisStatusNotAnalyzed)), "⚪ Not Analyzed"},
		{"assessment status", assessmentStatusTitle(string(AssessmentStatusPartiallyAssessed)), "🟡 Partially Assessed"},
		{"rating status", ratingStatusTitle(string(RatingStatusNotRated)), "⚪ Not Rated"},
		{"confidence", confidenceTitle(string(ConfidenceMedium)), "🔵 Medium"},
		{"run gap kind", runGapKindTitle(GapMissingEvaluationData), "📭 Missing Evaluation Data"},
		{"rating result kind", ratingResultKindTitle(RatingResultNotAssessed), "⚪ Not Assessed"},
		{"report kind", reportKindTitle(string(ReportKindFactor)), "🧩 Factor"},
		{"boolean", boolTitle(true), "✅ Yes"},
		{"known finding severity", findingSeverityTitle("high"), "🔴 High"},
		{"unknown fallback", findingTypeTitle("new_finding_type"), "New Finding Type"},
		{"camel fallback", limitTypeTitle("futureLimitType"), "Future Limit Type"},
	}
	for _, tc := range tests {
		if tc.got != tc.want {
			t.Fatalf("%s title = %q, want %q", tc.name, tc.got, tc.want)
		}
	}
}

func TestReportPairCellsRenderEmptyComponents(t *testing.T) {
	req := &evaluationRequirementArtifacts{}
	if got := evaluationConfidencePair(map[string]any{}, map[string]any{}); got != "— / —" {
		t.Fatalf("evaluationConfidencePair() = %q, want empty components marked", got)
	}
	if got := evaluationAnalysisStatusPair(map[string]any{}, map[string]any{}); got != "— / —" {
		t.Fatalf("evaluationAnalysisStatusPair() = %q, want empty components marked", got)
	}
	if got := evaluationRequirementConfidencePair(req); got != "— / —" {
		t.Fatalf("evaluationRequirementConfidencePair() = %q, want empty components marked", got)
	}
}

func TestDataKindDisplayTitlesCoverEvaluationDataKinds(t *testing.T) {
	kinds := append([]DataKind{}, acceptedDataKinds...)
	kinds = append(kinds, DataKindEvaluationOutput)
	for _, kind := range kinds {
		if got := dataKindTitle(kind); got == "" || got == humanizeEnum(string(kind)) {
			t.Fatalf("dataKindTitle(%s) = %q, want explicit title", kind, got)
		}
	}
}

func completeRootEvaluationPayloads() []string {
	return []string{
		`{"schemaVersion":2,"kind":"EvaluationFrame"}`,
		`{"schemaVersion":2,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}`,
		`{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Local work meets the bar.","ratingDrivers":[],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"The model meets the bar overall.","ratingDrivers":[],"confidence":"high"}}`,
	}
}

const navigationReportModel = `---
title: Navigation model
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
  - level: below
    title: 🔴 Below
    description: Below.
    criterion: Misses it.
factors:
  reliability:
    title: Reliability
    factors:
      latency:
        title: Latency
areas:
  payments:
    title: Payments
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
---
`

func navigationReportPayloads() []string {
	rootNote := testAreaFindingJSON(map[string]any{
		"id":         "root-note",
		"type":       "note",
		"severity":   "low",
		"confidence": "medium",
		"statement":  "Root note is contextual.",
		"effect": map[string]any{
			"statement":    "The note preserves context without driving the rating.",
			"ratingEffect": "context only",
		},
		"inputRefs": []any{map[string]any{"kind": "RequirementRatingResult", "subject": map[string]any{"requirementId": "requirement:root::has-tests"}}},
	})
	rootRisk := testAreaFindingJSON(map[string]any{
		"id":         "root-risk",
		"type":       "risk",
		"severity":   "high",
		"confidence": "high",
		"statement":  "Root reliability is the highest concern.",
		"condition":  "Reliability holds the root roll-up down.",
		"cause": map[string]any{
			"status":    "verified",
			"statement": "Reliability analysis is the binding roll-up driver.",
		},
		"effect": map[string]any{
			"statement":    "Reliability constrains the root aggregate rating.",
			"ratingEffect": "primary roll-up driver",
		},
		"inputRefs": []any{map[string]any{"kind": "FactorAnalysisResult", "subject": map[string]any{"factorId": "factor:root::reliability"}, "selector": "localAndDescendantAnalysis"}},
		"factorRelationships": []any{map[string]any{
			"factorId":     "factor:root::reliability",
			"relationship": "primary-driver",
			"rationale":    "Reliability is the affected factor.",
		}},
	})
	return []string{
		`{"schemaVersion":2,"kind":"EvaluationFrame"}`,
		`{"schemaVersion":2,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}`,
		`{"schemaVersion":2,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments"}}`,
		`{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[` + rootNote + `,` + rootRisk + `],"localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Root local work meets the bar.","ratingDrivers":[],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Root work meets the bar overall.","ratingDrivers":[],"confidence":"high"}}`,
		`{"schemaVersion":2,"kind":"AreaAnalysisResult","areaId":"area:payments","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Payments local work meets the bar.","ratingDrivers":[],"confidence":"medium"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Payments work meets the bar overall.","ratingDrivers":[],"confidence":"medium"}}`,
		`{"schemaVersion":2,"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability"}}`,
		`{"schemaVersion":2,"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability/latency"}}`,
		`{"schemaVersion":2,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability local work meets the bar.","ratingDrivers":[],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:below","rationale":"Reliability work misses the bar once sub-Factors roll up.","ratingDrivers":[],"confidence":"high"}}`,
		`{"schemaVersion":2,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability/latency","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Latency local work meets the bar.","ratingDrivers":[],"confidence":"medium"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Latency work meets the bar overall.","ratingDrivers":[],"confidence":"medium"}}`,
		`{"schemaVersion":2,"kind":"RequirementEvaluationFrame","subject":{"requirementId":"requirement:root::has-tests","factorIds":["factor:root::reliability"]}}`,
		`{"schemaVersion":2,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","confidence":"high","summary":"Tests are present.","factors":["factor:root::reliability"],"findings":[]}`,
		`{"schemaVersion":2,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","confidence":"high","rationale":"Tests meet the bar.","ratingDrivers":[]}`,
	}
}

func readReport(t *testing.T, runPath, rel string) string {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(runPath, rel))
	if err != nil {
		t.Fatalf("reading %s: %v", rel, err)
	}
	return string(raw)
}

func assertOnlyRootReportMD(t *testing.T, runPath string) {
	t.Helper()
	var found []string
	err := filepath.WalkDir(runPath, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || entry.Name() != "report.md" {
			return nil
		}
		rel, err := filepath.Rel(runPath, path)
		if err != nil {
			return err
		}
		found = append(found, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatalf("walking reports: %v", err)
	}
	if len(found) != 1 || found[0] != "report.md" {
		t.Fatalf("report.md files = %v, want only root report.md", found)
	}
}

func assertContains(t *testing.T, got, want string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("got:\n%s\nwant substring:\n%s", got, want)
	}
}

func assertNotContains(t *testing.T, got, unwanted string) {
	t.Helper()
	if strings.Contains(got, unwanted) {
		t.Fatalf("got:\n%s\nunwanted substring:\n%s", got, unwanted)
	}
}

func testRepo(t *testing.T) string {
	t.Helper()
	return testRepoWithModel(t, testModel)
}

func testRepoWithModel(t *testing.T, model string) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeRepoModel(t, repo, "QUALITY.md", model)
	return repo
}

func testRepoWithModelAt(t *testing.T, modelPath, model string) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeRepoModel(t, repo, modelPath, model)
	return repo
}

func writeRepoModel(t *testing.T, repo, modelPath, model string) {
	t.Helper()
	path := filepath.Join(repo, filepath.FromSlash(modelPath))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll(model dir) error = %v", err)
	}
	if err := os.WriteFile(path, []byte(model), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", modelPath, err)
	}
}
