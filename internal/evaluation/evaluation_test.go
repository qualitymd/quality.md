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
	manifest, err := loadRunManifest(filepath.Join(repo, result.Path))
	if err != nil {
		t.Fatalf("loadRunManifest() error = %v", err)
	}
	if manifest.Number != 1 || manifest.RequestedScope.AreaID != "" || manifest.PlannedScope.AreaID != "area:root" || len(manifest.PlannedScope.FactorFilter) != 0 {
		t.Fatalf("manifest = %#v, want full root scope", manifest)
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

func TestCreateRunUsesScopeFlagsAndManifest(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Scoped model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
factors:
  security:
    title: Security
    factors:
      reliability:
        title: Reliability
        factors:
          latency:
            title: Latency
---
`)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Factors: []string{"factor:root::security/reliability/latency"}})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	if result.Path != ".quality/evaluations/0001-root-security-reliability-latency-eval" {
		t.Fatalf("path = %q, want numbered scope-path run", result.Path)
	}
	manifest, err := loadRunManifest(filepath.Join(repo, result.Path))
	if err != nil {
		t.Fatalf("loadRunManifest() error = %v", err)
	}
	if manifest.RequestedScope.AreaID != "area:root" || manifest.PlannedScope.AreaID != "area:root" ||
		len(manifest.PlannedScope.FactorFilter) != 1 || manifest.PlannedScope.FactorFilter[0] != "factor:root::security/reliability/latency" {
		t.Fatalf("manifest = %#v, want factor-filtered root scope", manifest)
	}
}

func TestCreateRunRejectsUnknownScopeWithoutCreatingRun(t *testing.T) {
	repo := testRepo(t)
	_, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Factors: []string{"factor:root::missing"}})
	if err == nil || !strings.Contains(err.Error(), "does not resolve in the model") {
		t.Fatalf("CreateRun() error = %v, want unresolved factor diagnostic", err)
	}
	evalDir := filepath.Join(repo, ".quality", "evaluations")
	entries, readErr := os.ReadDir(evalDir)
	if readErr != nil && !os.IsNotExist(readErr) {
		t.Fatalf("ReadDir() error = %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("evaluation dir entries = %d, want no numbered run folder", len(entries))
	}
}

func TestRunNameRecognitionReadsCurrentNumber(t *testing.T) {
	tests := []struct {
		name       string
		wantNumber int
		wantOK     bool
	}{
		{name: "0007-full-eval", wantNumber: 7, wantOK: true},
		{name: "0007-security-eval", wantNumber: 7, wantOK: true},
		{name: "0007-security-network-eval", wantNumber: 7, wantOK: true},
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
			if got.number != tc.wantNumber {
				t.Fatalf("parseRunName() = %#v, want number %d", got, tc.wantNumber)
			}
		})
	}
}

func TestListRunsReadsScopeFromManifestAfterFolderRename(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Area: "area:root"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	renamed := filepath.Join(repo, ".quality", "evaluations", "renamed-run")
	if err := os.Rename(filepath.Join(repo, result.Path), renamed); err != nil {
		t.Fatalf("rename run: %v", err)
	}
	runs, err := ListRuns(repo, "", "")
	if err != nil {
		t.Fatalf("ListRuns() error = %v", err)
	}
	if len(runs.Runs) != 1 {
		t.Fatalf("ListRuns() = %#v, want one recognized run", runs.Runs)
	}
	if runs.Runs[0].Path != ".quality/evaluations/renamed-run" || runs.Runs[0].PlannedScope.AreaID != "area:root" {
		t.Fatalf("current run = %#v, want renamed path and manifest scope", runs.Runs[0])
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

func TestStatusRequiresPlannedExpansionArtifacts(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	if _, err := SetData(runPath, batchPayloads(completeRootEvaluationPayloads()[:3]...), DataSetOptions{}); err != nil {
		t.Fatalf("SetData(partial root batch) error = %v", err)
	}
	run, err := Inspect(runPath)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	status := run.Status()
	if status.Reportable || len(status.Gaps) == 0 {
		t.Fatalf("status = %#v, want planned expansion gap", status)
	}
	if status.Gaps[0].Ref != "data/areas/root/requirements/has-tests/requirement-assessment-result.json" {
		t.Fatalf("status.Gaps = %#v, want missing planned Requirement assessment", status.Gaps)
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
	if status := run.Status(); !status.Reportable || status.Data.Artifacts != 9 {
		t.Fatalf("status = %#v, want reportable with nine data artifacts", status)
	}
	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	if build.ReportMD == "" || build.EvaluationOutputResult == "" {
		t.Fatalf("BuildReport() = %#v, want evaluation report and output result", build)
	}
	if build.ReportMD != filepath.ToSlash(filepath.Join(result.Path, "report.md")) {
		t.Fatalf("BuildReport().ReportMD = %q, want run report path", build.ReportMD)
	}
	for _, name := range []string{"report-summary.md", "report.json"} {
		if _, err := os.Stat(filepath.Join(runPath, name)); !os.IsNotExist(err) {
			t.Fatalf("%s should not be generated by evaluation report build: %v", name, err)
		}
	}
	runReport, err := os.ReadFile(filepath.Join(runPath, "report.md"))
	if err != nil {
		t.Fatalf("reading report.md: %v", err)
	}
	if !strings.Contains(string(runReport), "type: Evaluation Overview Report\n") ||
		!strings.Contains(string(runReport), "title: Test model\n") ||
		!strings.Contains(string(runReport), "# Evaluation Report: Area: Test model") {
		t.Fatalf("report.md = %s, want run report title", runReport)
	}
	if !strings.Contains(string(runReport), "Run: #1 - Created: ") ||
		!strings.Contains(string(runReport), "Report: Overview - [Findings](findings.md) - [Recommendations](recommendations.md)") {
		t.Fatalf("report.md = %s, want run header navigation", runReport)
	}
	if !strings.Contains(string(runReport), "  - data/run-manifest.json\n") ||
		!strings.Contains(string(runReport), "  - data/areas/root/area-analysis-result.json\n") ||
		strings.Contains(string(runReport), "  - data/evaluation-output-result.json\n") {
		t.Fatalf("report.md = %s, want source-data frontmatter without EvaluationOutputResult", runReport)
	}
	if !strings.Contains(string(runReport), "| 🔵 Target | full evaluation | 🟢 High / 🟢 High |") {
		t.Fatalf("report.md = %s, want scoped Area rating row", runReport)
	}
	report, err := os.ReadFile(filepath.Join(runPath, "root-area.md"))
	if err != nil {
		t.Fatalf("reading root-area.md: %v", err)
	}
	if !strings.Contains(string(report), "type: Area Evaluation Report\n") ||
		!strings.Contains(string(report), "title: Test model\n") ||
		!strings.Contains(string(report), "# Area: Test model") ||
		!strings.Contains(string(report), "Area: [Test model](root-area.md)") {
		t.Fatalf("root-area.md = %s, want kind-prefixed title before Area trail", report)
	}
	if !strings.Contains(string(report), "  - data/run-manifest.json\n") ||
		!strings.Contains(string(report), "  - data/areas/root/area-analysis-result.json\n") {
		t.Fatalf("root-area.md = %s, want source-data frontmatter", report)
	}
	if !strings.Contains(string(report), "| 🔵 Target | 🔵 Target | 🟢 High / 🟢 High |") {
		t.Fatalf("root-area.md = %s, want target rating title", report)
	}
	if !strings.Contains(string(report), "## Legend\n\n- `—` - not applicable or not recorded.") {
		t.Fatalf("root-area.md = %s, want empty-cell legend", report)
	}
	if !strings.Contains(string(report), "Area: [Test model](root-area.md)") {
		t.Fatalf("root-area.md = %s, want Area trail", report)
	}
	if !strings.Contains(string(report), "| Overall Rating | Local Rating | Confidence |") ||
		strings.Contains(string(report), "| Overall Rating | Local Rating | Confidence | Data |") {
		t.Fatalf("root-area.md = %s, want confidence display titles", report)
	}
	if strings.Contains(string(report), "Overall Rating | target") {
		t.Fatalf("root-area.md = %s, want rating title instead of level ID", report)
	}
	if !strings.Contains(string(report), "## Area / Factor Breakdown") ||
		!strings.Contains(string(report), "| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |") ||
		!strings.Contains(string(report), "| **[Test model](root-area.md)** | 🔵 Target | 🔵 Target | 1 | 1 |") ||
		strings.Contains(string(report), "## Factors") ||
		strings.Contains(string(report), "## Child Areas") {
		t.Fatalf("root-area.md = %s, want Area / Factor breakdown without legacy Area tables", report)
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
	if output["runReportRef"] == nil || output["scopedAreaAnalysisRef"] == nil || output["headlineResultRef"] != nil || output["headlineReportRef"] != nil {
		t.Fatalf("EvaluationOutputResult = %#v, want run and scoped Area refs without headline refs", output)
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
		"basis": map[string]any{
			"status":    "not_assessed",
			"statement": "Basis was not assessed.",
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
	runReport := readReport(t, runPath, "report.md")
	assertContains(t, runReport, "type: Evaluation Overview Report\n")
	assertContains(t, runReport, "title: Navigation model\n")
	assertContains(t, runReport, "data:\n  - data/run-manifest.json")
	assertNotContains(t, runReport, "  - data/evaluation-output-result.json\n")
	assertContains(t, runReport, "# Evaluation Report: Area: Navigation model")
	assertContains(t, runReport, "Run: #1 - Created: ")
	assertContains(t, runReport, "Report: Overview - [Findings](findings.md) - [Recommendations](recommendations.md)")
	assertContains(t, runReport, "Area: [Navigation model](root-area.md)")
	assertContains(t, runReport, "| Overall Rating | Scope | Confidence |")
	assertNotContains(t, runReport, "| Overall Rating | Scope | Confidence | Data |")
	assertNotContains(t, runReport, "[evaluation-output-result.json](data/evaluation-output-result.json)")
	assertNotContains(t, runReport, "## Rating Drivers")
	assertNotContains(t, runReport, "| Driver | Effect | Inputs |")
	assertContains(t, runReport, "## Top Findings")
	assertContains(t, runReport, "| Rank | Finding | Area | Factors | Type | Severity |")
	assertContains(t, runReport, "| 1 | [Tests are present.](requirements/has-tests/has-tests-requirement.md#finding-strength-1) | [Navigation model](root-area.md) | [Reliability](factors/reliability/reliability-factor.md) | ✅ Strength | 🔵 Low |")
	assertContains(t, runReport, "Full findings index: [findings.md](findings.md)")
	assertContains(t, runReport, "## Top Recommendations")
	assertContains(t, runReport, "| Rank | Recommendation | Area / Factors | Reason |")
	assertContains(t, runReport, "| 1 | [Review the next quality bar](recommendations/001-review-the-next-quality-bar.md) | [Navigation model](root-area.md) / [Reliability](factors/reliability/reliability-factor.md) | The quality model stays aligned with the evaluated evidence and next bar. |")
	assertContains(t, runReport, "[recommendations.md](recommendations.md)")
	assertContains(t, runReport, "## Area / Factor Breakdown")
	assertContains(t, runReport, "| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |")
	assertContains(t, runReport, "| **[Navigation model](root-area.md)** | 🔵 Target | 🔵 Target | 1 | 1 |")
	assertContains(t, runReport, "| ↳ 🧩 [Reliability](factors/reliability/reliability-factor.md) | 🔴 Below | 🔵 Target | 1 | 1 |")
	assertContains(t, runReport, "| ↳ ↳ 🧩 [Latency](factors/reliability/factors/latency/latency-factor.md) | 🔵 Target | 🔵 Target | 0 | 0 |")
	assertContains(t, runReport, "| ↳ [Payments](areas/payments/payments-area.md) | 🔵 Target | 🔵 Target | 0 | 0 |")
	assertNotContains(t, runReport, "## Subject Reports")
	assertNotContains(t, runReport, "| Subject | Kind | Rating | Report |")

	recommendationIndex := readReport(t, runPath, "recommendations.md")
	assertContains(t, recommendationIndex, "type: Recommendation Index Report\n")
	assertContains(t, recommendationIndex, "title: Recommendations\n")
	assertContains(t, recommendationIndex, "data:\n  - data/run-manifest.json")
	assertContains(t, recommendationIndex, "  - data/advice/recommendation-ranking-result.json")
	assertContains(t, recommendationIndex, "Report: [Overview](report.md) - [Findings](findings.md) - Recommendations")
	assertContains(t, recommendationIndex, "| Recommendations | Highest Impact | Coverage |")
	assertNotContains(t, recommendationIndex, "| Recommendations | Highest Impact | Coverage | Data |")
	assertContains(t, recommendationIndex, "| Rank | Recommendation | Area / Factors | Impact | Confidence | Reason | Ranking Rationale |")
	assertContains(t, recommendationIndex, "| 1 | [Review the next quality bar](recommendations/001-review-the-next-quality-bar.md) | [Navigation model](root-area.md) / [Reliability](factors/reliability/reliability-factor.md) | High | 🟢 High | The quality model stays aligned with the evaluated evidence and next bar. | This recommendation addresses the highest-ranked finding. |")

	recommendationReport := readReport(t, runPath, "recommendations/001-review-the-next-quality-bar.md")
	assertContains(t, recommendationReport, "type: Recommendation Report\n")
	assertContains(t, recommendationReport, "title: Review the next quality bar\n")
	assertContains(t, recommendationReport, "data:\n  - data/run-manifest.json")
	assertContains(t, recommendationReport, "  - data/advice/recommendations/rec-001/recommendation-result.json")
	assertContains(t, recommendationReport, "# Recommendation: Review the next quality bar")
	assertContains(t, recommendationReport, "| Rank | Impact | Confidence |")
	assertNotContains(t, recommendationReport, "| Rank | Impact | Confidence | Data |")
	assertContains(t, recommendationReport, "## Description")
	assertContains(t, recommendationReport, "## Background")
	assertContains(t, recommendationReport, "## Expected value")
	assertContains(t, recommendationReport, "## Done criterion")
	assertContains(t, recommendationReport, "## Ranking rationale")
	assertNotContains(t, recommendationReport, "## Why it matters")
	assertNotContains(t, recommendationReport, "## Expected benefit")

	rootReport := readReport(t, runPath, "root-area.md")
	assertContains(t, rootReport, "type: Area Evaluation Report\n")
	assertContains(t, rootReport, "title: Navigation model\n")
	assertContains(t, rootReport, "# Area: Navigation model")
	assertContains(t, rootReport, "Area: [Navigation model](root-area.md)")
	assertNotContains(t, rootReport, "Path: `/`")
	assertContains(t, rootReport, "data:\n  - data/run-manifest.json")
	assertContains(t, rootReport, "  - data/areas/root/area-analysis-result.json")
	assertContains(t, rootReport, "| Overall Rating | Local Rating | Confidence |")
	assertNotContains(t, rootReport, "| Overall Rating | Local Rating | Confidence | Data |")
	assertNotContains(t, rootReport, "## Findings")
	assertNotContains(t, rootReport, "## Rating Drivers")
	assertNotContains(t, rootReport, "| Driver | Effect | Inputs |")
	assertNotContains(t, rootReport, "Reliability analysis is the binding roll-up driver.")
	assertContains(t, rootReport, "## Area / Factor Breakdown")
	assertContains(t, rootReport, "| Area / Factor | Overall Rating | Local Rating | Findings | Recommendations |")
	assertContains(t, rootReport, "| **[Navigation model](root-area.md)** | 🔵 Target | 🔵 Target | 1 | 1 |")
	assertContains(t, rootReport, "| ↳ 🧩 [Reliability](factors/reliability/reliability-factor.md) | 🔴 Below | 🔵 Target | 1 | 1 |")
	assertContains(t, rootReport, "| ↳ ↳ 🧩 [Latency](factors/reliability/factors/latency/latency-factor.md) | 🔵 Target | 🔵 Target | 0 | 0 |")
	assertContains(t, rootReport, "| ↳ [Payments](areas/payments/payments-area.md) | 🔵 Target | 🔵 Target | 0 | 0 |")
	assertNotContains(t, rootReport, "## Factors")
	assertNotContains(t, rootReport, "## Child Areas")
	assertNotContains(t, rootReport, "| Factor | Path | Local Rating | + Sub-Factors Rating | Sub-Factors |")
	assertNotContains(t, rootReport, "| Area | Path | Local Rating | + Child Areas Rating | Factors |")
	assertContains(t, rootReport, "| [Has tests](requirements/has-tests/has-tests-requirement.md) | 🔵 Target | ✅ Assessed | [reliability](factors/reliability/reliability-factor.md) |")
	assertContains(t, rootReport, "## Legend\n\n- `—` - not applicable or not recorded.")
	assertNotContains(t, rootReport, "## Sub-Areas")
	assertNotContains(t, rootReport, "+ Sub-Areas Rating")
	assertNotContains(t, rootReport, "Breadcrumb:")
	assertNotContains(t, rootReport, "Parent Area:")
	assertNotContains(t, rootReport, "| Details |")

	factorReport := readReport(t, runPath, "factors/reliability/reliability-factor.md")
	assertContains(t, factorReport, "type: Factor Evaluation Report\n")
	assertContains(t, factorReport, "title: Reliability\n")
	assertContains(t, factorReport, "# Factor: Reliability")
	assertContains(t, factorReport, "Area: [Navigation model](../../root-area.md)")
	assertContains(t, factorReport, "Factor: [Reliability](reliability-factor.md)")
	assertNotContains(t, factorReport, "Path: `reliability`")
	assertContains(t, factorReport, "data:\n  - data/run-manifest.json")
	assertContains(t, factorReport, "  - data/areas/root/factors/reliability/factor-analysis-result.json")
	assertContains(t, factorReport, "| Overall Rating | Local Rating | Status | Confidence |")
	assertNotContains(t, factorReport, "| Overall Rating | Local Rating | Status | Confidence | Data |")
	assertNotContains(t, factorReport, "## Findings")
	assertNotContains(t, factorReport, "## Rating Drivers")
	assertNotContains(t, factorReport, "| Driver | Effect | Inputs |")
	assertNotContains(t, factorReport, "Latency roll-up constrains reliability.")
	assertContains(t, factorReport, "| [Has tests](../../requirements/has-tests/has-tests-requirement.md) | 🔵 Target | ✅ Assessed |")
	assertContains(t, factorReport, "## Sub-Factors")
	assertContains(t, factorReport, "| [Latency](factors/latency/latency-factor.md) | `reliability/latency` | 🔵 Target | — |")
	assertNotContains(t, factorReport, "## Child Factors")
	assertNotContains(t, factorReport, "Parent Factor:")
	assertNotContains(t, factorReport, "| Details |")

	paymentsReport := readReport(t, runPath, "areas/payments/payments-area.md")
	assertContains(t, paymentsReport, "## Area / Factor Breakdown")
	assertContains(t, paymentsReport, "| **[Payments](payments-area.md)** | 🔵 Target | 🔵 Target | 0 | 0 |")
	assertNotContains(t, paymentsReport, "## Child Areas")
	assertNotContains(t, paymentsReport, "| (no Child Areas) |  |  |  |  |")
	assertNotContains(t, paymentsReport, "## Factors")
	assertNotContains(t, paymentsReport, "Sub-Areas")

	childFactorReport := readReport(t, runPath, "factors/reliability/factors/latency/latency-factor.md")
	assertContains(t, childFactorReport, "Factor: [Reliability](../../reliability-factor.md) / [Latency](latency-factor.md)")
	assertContains(t, childFactorReport, "## Sub-Factors")
	assertContains(t, childFactorReport, "| (no Sub-Factors) |  |  |  |")
	assertNotContains(t, childFactorReport, "Child Factors")

	requirementReport := readReport(t, runPath, "requirements/has-tests/has-tests-requirement.md")
	assertContains(t, requirementReport, "type: Requirement Evaluation Report\n")
	assertContains(t, requirementReport, "title: Has tests\n")
	assertContains(t, requirementReport, "# Requirement: Has tests")
	assertContains(t, requirementReport, "Area: [Navigation model](../../root-area.md)")
	assertContains(t, requirementReport, "Factors: [reliability](../../factors/reliability/reliability-factor.md)")
	assertContains(t, requirementReport, `<a id="finding-strength-1"></a>`)
	assertContains(t, requirementReport, "| Advice Rank | Tier | Ranking Rationale |")
	assertContains(t, requirementReport, "| 1 / 1 | P1 | This finding most directly informs next advice. |")
	assertNotContains(t, requirementReport, "Name: `has-tests`")
	assertContains(t, requirementReport, "data:\n  - data/run-manifest.json")
	assertContains(t, requirementReport, "  - data/areas/root/requirements/has-tests/requirement-assessment-result.json")
	assertContains(t, requirementReport, "  - data/areas/root/requirements/has-tests/requirement-rating-result.json")
	assertContains(t, requirementReport, "| Rating | Assessment | Confidence |")
	assertNotContains(t, requirementReport, "| Rating | Assessment | Confidence | Data |")
	assertNotContains(t, requirementReport, "| Rating | Assessment | Factors | Confidence | Data |")
	assertNotContains(t, requirementReport, "[requirement-assessment-result.json](")
	assertNotContains(t, requirementReport, "[requirement-rating-result.json](")
	assertNotContains(t, requirementReport, "\nFactor:")
	assertNotContains(t, requirementReport, "Parent Area:")
	assertOnlyRootReportMD(t, runPath)

	findingsReport := readReport(t, runPath, "findings.md")
	assertContains(t, findingsReport, "type: Finding Index Report\n")
	assertContains(t, findingsReport, "title: Findings\n")
	assertContains(t, findingsReport, "# Findings")
	assertContains(t, findingsReport, "Report: [Overview](report.md) - Findings - [Recommendations](recommendations.md)")
	assertContains(t, findingsReport, "data:\n  - data/run-manifest.json")
	assertContains(t, findingsReport, "  - data/advice/finding-ranking-result.json")
	assertNotContains(t, findingsReport, "[finding-ranking-result.json](data/advice/finding-ranking-result.json)")
	assertContains(t, findingsReport, "| Rank | Finding | Area | Factors | Type | Severity |")
	assertContains(t, findingsReport, "| 1 | [Tests are present.](requirements/has-tests/has-tests-requirement.md#finding-strength-1) | [Navigation model](root-area.md) | [Reliability](factors/reliability/reliability-factor.md) | ✅ Strength | 🔵 Low |")

	outputRaw, err := os.ReadFile(filepath.Join(runPath, "data", "evaluation-output-result.json"))
	if err != nil {
		t.Fatalf("reading EvaluationOutputResult: %v", err)
	}
	if !strings.Contains(string(outputRaw), `"path": "factors/reliability/reliability-factor.md"`) {
		t.Fatalf("EvaluationOutputResult = %s, want subject-aware report refs", outputRaw)
	}
	if !strings.Contains(string(outputRaw), `"runReportRef":`) || !strings.Contains(string(outputRaw), `"kind": "run"`) || !strings.Contains(string(outputRaw), `"path": "root-area.md"`) || !strings.Contains(string(outputRaw), `"kind": "findings"`) {
		t.Fatalf("EvaluationOutputResult = %s, want run and root Area report refs", outputRaw)
	}
	if build.RatingResult.Level != "target" {
		t.Fatalf("BuildReport().RatingResult.Level = %q, want stable rating level ID", build.RatingResult.Level)
	}
}

func TestScopedAreaEvaluationBuildsRunReportWithoutRootArea(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Scoped model
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
areas:
  payments:
    title: Payments
    areas:
      refunds:
        title: Refunds
---
`)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Area: "area:payments/refunds"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	payloads := []string{
		`{"schemaVersion":3,"kind":"EvaluationFrame","inputs":{},"derivedContext":{}}`,
		`{"schemaVersion":3,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments/refunds"}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:payments/refunds","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Refunds meet the scoped bar.","ratingDrivers":[{"description":"Refunds frame supports local analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments/refunds"}}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Refunds meet the scoped bar overall.","ratingDrivers":[{"description":"Refunds frame supports overall analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments/refunds"}}]}],"confidence":"high"}}`,
	}
	payloads = append(payloads, noFindingAreaAdvicePayloads("area:payments/refunds")...)
	if _, err := SetData(runPath, batchPayloads(payloads...), DataSetOptions{}); err != nil {
		t.Fatalf("SetData(scoped area batch) error = %v", err)
	}
	run, err := Inspect(runPath)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if status := run.Status(); !status.Reportable {
		t.Fatalf("status = %#v, want scoped area run reportable", status)
	}
	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	if build.RatingResult.Level != "target" {
		t.Fatalf("BuildReport() = %#v, want scoped Area rating", build)
	}
	if _, err := os.Stat(filepath.Join(runPath, "root-area.md")); !os.IsNotExist(err) {
		t.Fatalf("root-area.md should not be generated for root-out-of-scope run: %v", err)
	}
	runReport := readReport(t, runPath, "report.md")
	assertContains(t, runReport, "# Evaluation Report: Area: Refunds")
	assertContains(t, runReport, "| Planned Area | `area:payments/refunds` |")
	assertContains(t, runReport, "Root Area was not evaluated in this run.")
	areaReport := readReport(t, runPath, "areas/payments/refunds/refunds-area.md")
	assertContains(t, areaReport, "# Area: Refunds")
	assertContains(t, areaReport, "Area: Scoped model / Payments / [Refunds](refunds-area.md)")
	outputRaw, err := os.ReadFile(filepath.Join(runPath, "data", "evaluation-output-result.json"))
	if err != nil {
		t.Fatalf("reading EvaluationOutputResult: %v", err)
	}
	assertContains(t, string(outputRaw), `"scopedAreaAnalysisRef":`)
	assertContains(t, string(outputRaw), `"areaId": "area:payments/refunds"`)
	assertNotContains(t, string(outputRaw), `"headlineResultRef":`)
	assertNotContains(t, string(outputRaw), "rootAreaAnalysisRef")
}

func TestScopedFactorEvaluationBuildsRunReportWithoutRootArea(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Factor scoped model
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
factors:
  reliability:
    title: Reliability
---
`)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md", Factors: []string{"factor:root::reliability"}})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	payloads := []string{
		`{"schemaVersion":3,"kind":"EvaluationFrame","inputs":{},"derivedContext":{}}`,
		`{"schemaVersion":3,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability meets the scoped Area bar.","ratingDrivers":[{"description":"Reliability supports the scoped Area.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisResult","subject":{"factorId":"factor:root::reliability"},"selector":"localAndDescendantAnalysis"}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability meets the scoped Area bar overall.","ratingDrivers":[{"description":"Reliability supports the scoped Area overall.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisResult","subject":{"factorId":"factor:root::reliability"},"selector":"localAndDescendantAnalysis"}]}],"confidence":"high"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability meets the scoped bar.","ratingDrivers":[{"description":"Reliability frame supports local analysis.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability"}}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability meets the scoped bar overall.","ratingDrivers":[{"description":"Reliability frame supports overall analysis.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability"}}]}],"confidence":"high"}}`,
	}
	payloads = append(payloads, noFindingFactorAdvicePayloads("factor:root::reliability")...)
	if _, err := SetData(runPath, batchPayloads(payloads...), DataSetOptions{}); err != nil {
		t.Fatalf("SetData(scoped factor batch) error = %v", err)
	}
	run, err := Inspect(runPath)
	if err != nil {
		t.Fatalf("Inspect() error = %v", err)
	}
	if status := run.Status(); !status.Reportable {
		t.Fatalf("status = %#v, want scoped factor run reportable", status)
	}
	build, err := BuildReport(runPath)
	if err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	if build.RatingResult.Level != "target" {
		t.Fatalf("BuildReport() = %#v, want factor-filtered Area rating", build)
	}
	runReport := readReport(t, runPath, "report.md")
	assertContains(t, runReport, "# Evaluation Report: Area: Factor scoped model - Reliability")
	assertContains(t, runReport, "| Factor Filter | `factor:root::reliability` Reliability |")
	assertContains(t, runReport, "The scoped Area rating is a partial roll-up, not a complete Area assessment.")
	factorReport := readReport(t, runPath, "factors/reliability/reliability-factor.md")
	assertContains(t, factorReport, "# Factor: Reliability")
	assertContains(t, factorReport, "Area: [Factor scoped model](../../root-area.md)")
	outputRaw, err := os.ReadFile(filepath.Join(runPath, "data", "evaluation-output-result.json"))
	if err != nil {
		t.Fatalf("reading EvaluationOutputResult: %v", err)
	}
	assertContains(t, string(outputRaw), `"scopedAreaAnalysisRef":`)
	assertContains(t, string(outputRaw), `"factorId": "factor:root::reliability"`)
	assertNotContains(t, string(outputRaw), `"headlineResultRef":`)
}

func TestSetDataRejectsCLIOwnedOutput(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	_, err = SetData(filepath.Join(repo, result.Path), batchPayloads(`{"schemaVersion":3,"kind":"EvaluationOutputResult"}`), DataSetOptions{})
	if err == nil || !strings.Contains(err.Error(), "CLI-owned") {
		t.Fatalf("SetData(EvaluationOutputResult) error = %v, want CLI-owned diagnostic", err)
	}
	_, err = SetData(filepath.Join(repo, result.Path), batchPayloads(`{"schemaVersion":3,"kind":"RunManifest"}`), DataSetOptions{})
	if err == nil || !strings.Contains(err.Error(), "CLI-owned") {
		t.Fatalf("SetData(RunManifest) error = %v, want CLI-owned diagnostic", err)
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
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":[],"localAnalysis":{"status":"analyzed","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"rated","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingDrivers":["finding-1"]},"localAndDescendantAnalysis":{"status":"analyzed","ratingDrivers":[]}}`,
		`{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":["finding-1"]}`,
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
	raw := `{"schemaVersion":3,"kind":"EvaluationFrame"}`
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
	valid := `{"schemaVersion":3,"kind":"EvaluationFrame"}`
	invalidRequirement := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::invented","status":"assessed","findings":[]}`
	invalidRating := `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:invented","ratingDrivers":[]}`
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
	payload := `{"schemaVersion":3,"kind":"EvaluationFrame"}`
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
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"title": "Wrong field"}) + `]}`,
			want: "unknown field title",
		},
		{
			name: "legacy requirement finding description field",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"description": "Legacy field"}) + `]}`,
			want: "unknown field description",
		},
		{
			name: "legacy candidate actions field",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"actions": []any{map[string]any{"description": "do it"}}}) + `]}`,
			want: "unknown field actions",
		},
		{
			name: "candidate action missing id",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"description": "do it"}}}) + `]}`,
			want: "is missing required field id",
		},
		{
			name: "candidate action missing description",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "rationale": "why"}}}) + `]}`,
			want: "is missing required field description",
		},
		{
			name: "candidate action unknown field",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "do it", "effort": "high"}}}) + `]}`,
			want: "unknown field effort",
		},
		{
			name: "candidate action non-string id",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": 5, "description": "do it"}}}) + `]}`,
			want: "must be a non-empty string",
		},
		{
			name: "candidate action non-string description",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": 5}}}) + `]}`,
			want: "must be a non-empty string",
		},
		{
			name: "requirement finding info severity",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"severity": "info"}) + `]}`,
			want: `severity = "info", want one of critical, high, medium, low`,
		},
		{
			name: "invented requirement",
			raw:  `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::invented","status":"assessed","findings":[]}`,
			want: "does not resolve in the model",
		},
		{
			name: "invented rating",
			raw:  `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:invented","ratingDrivers":[]}`,
			want: "does not resolve in the model",
		},
		{
			name: "out-of-vocabulary reference kind",
			raw:  `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","ratingDrivers":[{"description":"d","effect":"supports target","ratingLevelId":"rating:target","inputRefs":[{"kind":"RequirementAssessment","subject":{"requirementId":"requirement:root::has-tests"}}]}]}`,
			want: `kind = "RequirementAssessment", want one of`,
		},
		{
			name: "area analysis findings field removed",
			raw:  `{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","findings":[],"localAnalysis":{"status":"analyzed"},"localAndDescendantAnalysis":{"status":"analyzed"}}`,
			want: "unknown field findings",
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

func TestSetDataRejectsRatedRequirementWithoutFindingBackedAssessment(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	rating := `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","ratingDrivers":[{"description":"Driver.","effect":"supports target","inputRefs":[{"kind":"RequirementAssessmentResult","subject":{"requirementId":"requirement:root::has-tests"}}]}]}`
	if _, err := SetData(runPath, batchPayloads(rating), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), "requires paired RequirementAssessmentResult") {
		t.Fatalf("SetData() error = %v, want paired assessment diagnostic", err)
	}
	emptyAssessment := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[]}`
	if _, err := SetData(runPath, batchPayloads(emptyAssessment, rating), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), "at least one finding") {
		t.Fatalf("SetData() error = %v, want finding-backed assessment diagnostic", err)
	}
}

func TestSetDataRejectsRatedOutputsWithoutDriversOrResolvableRefs(t *testing.T) {
	repo := testRepoWithModel(t, `---
title: Test model
ratingScale:
  - level: target
    title: Target
    criterion: Meets it.
factors:
  reliability:
    title: Reliability
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
	assessment := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(nil) + `]}`
	cases := []struct {
		name     string
		payloads []string
		want     string
	}{
		{
			name:     "requirement rating without drivers",
			payloads: []string{assessment, `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","ratingDrivers":[]}`},
			want:     "requires at least one ratingDrivers entry",
		},
		{
			name:     "requirement rating unresolved driver ref",
			payloads: []string{assessment, `{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","ratingDrivers":[{"description":"Driver.","effect":"supports target","inputRefs":[{"kind":"RequirementAssessmentResult","subject":{"requirementId":"requirement:root::invented"}}]}]}`},
			want:     "does not resolve",
		},
		{
			name:     "factor analysis without drivers",
			payloads: []string{`{"schemaVersion":3,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed"}}`},
			want:     "requires at least one ratingDrivers entry",
		},
		{
			name:     "area analysis without drivers",
			payloads: []string{`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","ratingDrivers":[]},"localAndDescendantAnalysis":{"status":"analyzed"}}`},
			want:     "requires at least one ratingDrivers entry",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := SetData(runPath, batchPayloads(tc.payloads...), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("SetData() error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestSetDataAcceptsFindingCandidateActions(t *testing.T) {
	repo := testRepo(t)
	result, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, result.Path)
	raw := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}, map[string]any{"id": "action-002", "description": "Cover error paths.", "rationale": "Lifts coverage."}}}) + `]}`
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
	raw := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"candidateActions": []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}, map[string]any{"id": "action-001", "description": "Cover error paths."}}}) + `]}`
	if _, err := SetData(runPath, batchPayloads(raw), DataSetOptions{DryRun: true}); err == nil || !strings.Contains(err.Error(), `candidateActions[1].id "action-001" is duplicated`) {
		t.Fatalf("SetData() error = %v, want duplicate candidate action ID diagnostic", err)
	}
}

func TestFindingDetailsOmitCandidateActions(t *testing.T) {
	var b strings.Builder
	finding := testFindingCore("gap-1", "gap", "high", "high", "Edge cases untested.")
	finding["candidateActions"] = []any{map[string]any{"id": "action-001", "description": "Add boundary tests."}}
	req := &evaluationRequirementArtifacts{
		ID:         requirementID{Name: "has-tests"},
		Assessment: map[string]any{"findings": []any{finding}},
	}
	writeEvaluationFindingDetails(&b, &evaluationArtifacts{}, req)
	rendered := b.String()
	if strings.Contains(rendered, "| Actions |") {
		t.Fatalf("finding details = %s, want no Actions row in the v0 report", rendered)
	}
	if strings.Contains(rendered, "Add boundary tests.") {
		t.Fatalf("finding details = %s, want candidate actions kept out of the report", rendered)
	}
	for _, want := range []string{`<a id="finding-gap-1"></a>`, "### gap-1 Edge cases untested.", "| (not ranked) | — | — |", "#### Condition", "#### Criteria", "#### Basis", "#### Effect", "#### Evidence"} {
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
	findings, ok := props["findings"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema findings property = %#v, want object", props["findings"])
	}
	items, ok := findings["items"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema findings items = %#v, want object", findings["items"])
	}
	findingProps, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema finding properties = %#v, want object", items["properties"])
	}
	for field, wantValues := range map[string][]string{
		"type":       {"strength", "gap", "risk", "unknown", "note"},
		"severity":   {"critical", "high", "medium", "low"},
		"confidence": {"high", "medium", "low", "none"},
	} {
		prop, ok := findingProps[field].(map[string]any)
		if !ok {
			t.Fatalf("requirement schema finding %s property = %#v, want object", field, findingProps[field])
		}
		enumValues, ok := prop["enum"].([]any)
		if !ok {
			t.Fatalf("requirement schema finding %s enum = %#v, want array", field, prop["enum"])
		}
		for _, want := range wantValues {
			if !jsonArrayContains(enumValues, want) {
				t.Fatalf("requirement schema finding %s enum = %#v, want %q", field, enumValues, want)
			}
		}
	}
	basis, ok := findingProps["basis"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema basis = %#v, want object", findingProps["basis"])
	}
	basisProps, ok := basis["properties"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema basis properties = %#v, want object", basis["properties"])
	}
	basisStatus, ok := basisProps["status"].(map[string]any)
	if !ok {
		t.Fatalf("requirement schema basis status = %#v, want object", basisProps["status"])
	}
	basisStatusEnum, ok := basisStatus["enum"].([]any)
	if !ok {
		t.Fatalf("requirement schema basis status enum = %#v, want array", basisStatus["enum"])
	}
	for _, want := range []string{"verified", "plausible", "not_assessed", "not_applicable"} {
		if !jsonArrayContains(basisStatusEnum, want) {
			t.Fatalf("requirement schema basis status enum = %#v, want %q", basisStatusEnum, want)
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
	if _, ok := areaProps["findings"]; ok {
		t.Fatalf("area schema findings property = %#v, want removed", areaProps["findings"])
	}

	example, err := DataExample(DataKindRequirementAssessment)
	if err != nil {
		t.Fatalf("DataExample() error = %v", err)
	}
	for _, want := range []string{`"findings": [`, `"statement": "Focused test coverage is present."`, `"condition": "A focused test covers the requirement's primary path."`, `"basis": {`, `"effect": {`, `"candidateActions": [`, `"id": "action-001"`, `"description": "Add focused tests for the boundary and error paths."`} {
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

func TestFindingRankingExampleRanksEveryFindingIncludingTail(t *testing.T) {
	raw, err := DataExample(DataKindFindingRanking)
	if err != nil {
		t.Fatalf("DataExample(%s) error = %v", DataKindFindingRanking, err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("DataExample(%s) JSON error = %v\n%s", DataKindFindingRanking, err, raw)
	}
	ordered, ok := payload["orderedFindings"].([]any)
	if !ok || len(ordered) < 2 {
		t.Fatalf("finding-ranking example must rank more than one finding, got %v", payload["orderedFindings"])
	}
	hasTail := false
	for _, entry := range ordered {
		obj, _ := entry.(map[string]any)
		if obj["tier"] == "P4" {
			hasTail = true
		}
	}
	if !hasTail {
		t.Fatalf("finding-ranking example must include a lowest-tier (P4) tail entry to show completeness, got %s", raw)
	}
}

func TestFactorAnalysisExampleModelsUmbrellaFactor(t *testing.T) {
	raw, err := DataExample(DataKindFactorAnalysis)
	if err != nil {
		t.Fatalf("DataExample(%s) error = %v", DataKindFactorAnalysis, err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatalf("DataExample(%s) JSON error = %v\n%s", DataKindFactorAnalysis, err, raw)
	}

	local, _ := payload["localAnalysis"].(map[string]any)
	if local["status"] != "empty" {
		t.Fatalf("umbrella factor-analysis example localAnalysis.status = %v, want empty", local["status"])
	}
	if _, hasInputs := local["inputRefs"]; hasInputs {
		t.Fatalf("umbrella factor-analysis example localAnalysis must have no Requirement inputs, got %s", raw)
	}

	overall, _ := payload["localAndDescendantAnalysis"].(map[string]any)
	if overall["status"] != "analyzed" {
		t.Fatalf("umbrella factor-analysis example localAndDescendantAnalysis.status = %v, want analyzed", overall["status"])
	}
	inputRefs, ok := overall["inputRefs"].([]any)
	if !ok || len(inputRefs) == 0 {
		t.Fatalf("umbrella factor-analysis roll-up must carry inputRefs, got %s", raw)
	}
	ref, _ := inputRefs[0].(map[string]any)
	if ref["kind"] != string(DataKindFactorAnalysis) {
		t.Fatalf("umbrella factor-analysis roll-up input must reference a child FactorAnalysisResult, got %v", ref["kind"])
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
	badPayload := `{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","findings":[` + testRequirementFindingJSON(map[string]any{"title": "Wrong"}) + `]}`
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
		{"findings report kind", reportKindTitle(string(ReportKindFindings)), "🔝 Findings"},
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

func TestReportMarkdownHelpersEscapeTableContent(t *testing.T) {
	if got := markdownCell("first | second\nthird"); got != `first \| second<br>third` {
		t.Fatalf("markdownCell() = %q, want escaped pipe and normalized newline", got)
	}
	if got := reportLink("areas/api/api-area.md", "root-area.md", "Root [Area] | link"); got != `[Root \[Area\] \| link](../../root-area.md)` {
		t.Fatalf("reportLink() = %q, want escaped label and relative target", got)
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
	payloads := []string{
		`{"schemaVersion":3,"kind":"EvaluationFrame"}`,
		`{"schemaVersion":3,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Local work meets the bar.","ratingDrivers":[{"description":"Area frame supports local analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"The model meets the bar overall.","ratingDrivers":[{"description":"Area frame supports overall analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}]}],"confidence":"high"}}`,
		`{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","evidenceSummary":"Tests were inspected.","findings":[{"id":"strength-1","type":"strength","severity":"low","confidence":"high","statement":"Tests exist.","condition":"A focused test exists.","criteria":[{"requirementId":"requirement:root::has-tests","ratingLevelId":"rating:target","criterion":"Meets it."}],"basis":{"status":"not_applicable","statement":"No separate basis beyond cited evidence is claimed."},"effect":{"statement":"The finding supports the target rating."},"evidence":[{"sourceRef":"tests/example_test.go","statement":"A focused test exists."}]}]}`,
		`{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","rationale":"Evidence satisfies the target criterion.","ratingDrivers":[{"description":"Assessment finding supports the target rating.","effect":"supports target","inputRefs":[{"kind":"RequirementAssessmentResult","subject":{"requirementId":"requirement:root::has-tests"}}]}],"confidence":"high"}`,
	}
	return append(payloads, singleFindingAdvicePayloads("requirement:root::has-tests", "strength-1")...)
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
	requirementFinding := testRequirementFindingJSON(map[string]any{
		"id":         "strength-1",
		"type":       "strength",
		"severity":   "low",
		"confidence": "high",
		"statement":  "Tests are present.",
		"condition":  "The inspected source includes focused tests.",
		"effect": map[string]any{
			"statement":    "The finding supports the target rating.",
			"ratingEffect": "supports target",
		},
	})
	payloads := []string{
		`{"schemaVersion":3,"kind":"EvaluationFrame"}`,
		`{"schemaVersion":3,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:root"}}`,
		`{"schemaVersion":3,"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments"}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:root","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Root local work meets the bar.","ratingDrivers":[{"description":"Requirement rating supports root local analysis.","effect":"supports target","inputRefs":[{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:root::has-tests"}}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Root work meets the bar overall.","ratingDrivers":[{"description":"Reliability analysis is the binding roll-up driver.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisResult","subject":{"factorId":"factor:root::reliability"},"selector":"localAndDescendantAnalysis"}]}],"confidence":"high"}}`,
		`{"schemaVersion":3,"kind":"AreaAnalysisResult","areaId":"area:payments","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Payments local work meets the bar.","ratingDrivers":[{"description":"Payments frame supports local analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments"}}]}],"confidence":"medium"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Payments work meets the bar overall.","ratingDrivers":[{"description":"Payments frame supports overall analysis.","effect":"supports target","inputRefs":[{"kind":"AreaEvaluationFrame","subject":{"areaId":"area:payments"}}]}],"confidence":"medium"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability/latency"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Reliability local work meets the bar.","ratingDrivers":[{"description":"Requirement rating supports reliability.","effect":"supports target","inputRefs":[{"kind":"RequirementRatingResult","subject":{"requirementId":"requirement:root::has-tests"}}]}],"confidence":"high"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:below","rationale":"Reliability work misses the bar once Sub-Factors roll up.","ratingDrivers":[{"description":"Latency roll-up constrains reliability.","effect":"constrains below","inputRefs":[{"kind":"FactorAnalysisResult","subject":{"factorId":"factor:root::reliability/latency"},"selector":"localAndDescendantAnalysis"}]}],"confidence":"high"}}`,
		`{"schemaVersion":3,"kind":"FactorAnalysisResult","factorId":"factor:root::reliability/latency","localAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Latency local work meets the bar.","ratingDrivers":[{"description":"Latency frame supports local analysis.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability/latency"}}]}],"confidence":"medium"},"localAndDescendantAnalysis":{"status":"analyzed","ratingLevelId":"rating:target","rationale":"Latency work meets the bar overall.","ratingDrivers":[{"description":"Latency frame supports overall analysis.","effect":"supports target","inputRefs":[{"kind":"FactorAnalysisFrame","subject":{"factorId":"factor:root::reliability/latency"}}]}],"confidence":"medium"}}`,
		`{"schemaVersion":3,"kind":"RequirementEvaluationFrame","subject":{"requirementId":"requirement:root::has-tests","factorIds":["factor:root::reliability"]}}`,
		`{"schemaVersion":3,"kind":"RequirementAssessmentResult","requirementId":"requirement:root::has-tests","status":"assessed","confidence":"high","summary":"Tests are present.","factors":["factor:root::reliability"],"findings":[` + requirementFinding + `]}`,
		`{"schemaVersion":3,"kind":"RequirementRatingResult","requirementId":"requirement:root::has-tests","status":"rated","ratingLevelId":"rating:target","confidence":"high","rationale":"Tests meet the bar.","ratingDrivers":[{"description":"Focused test finding supports target.","effect":"supports target","ratingLevelId":"rating:target","inputRefs":[{"kind":"RequirementAssessmentResult","subject":{"requirementId":"requirement:root::has-tests"},"selector":"findings[strength-1]"}]}]}`,
	}
	return append(payloads, singleFindingAdvicePayloads("requirement:root::has-tests", "strength-1")...)
}

func singleFindingAdvicePayloads(requirementID, findingID string) []string {
	findingRef := `{"kind":"RequirementAssessmentResult","subject":{"requirementId":"` + requirementID + `"},"selector":"findings[` + findingID + `]"}`
	return []string{
		`{"schemaVersion":3,"kind":"FindingRankingResult","orderedFindings":[{"rank":1,"findingRef":` + findingRef + `,"tier":"P1","rationale":"This finding most directly informs next advice."}],"rationale":"Findings were ranked by quality-bar relevance and confidence."}`,
		`{"schemaVersion":3,"kind":"RecommendationResult","id":"rec-001","title":"Review the next quality bar","description":"Review whether the next rating level should be raised or clarified for this requirement.","background":"The evaluation found a quality signal that should inform the next target level.","expectedValue":"The quality model stays aligned with the evaluated evidence and next bar.","doneCriterion":"The requirement criterion or target-level rationale reflects the review decision.","impact":"high","confidence":"high","traceRefs":[` + findingRef + `]}`,
		`{"schemaVersion":3,"kind":"RecommendationRankingResult","orderedRecommendations":[{"rank":1,"recommendationRef":"rec-001","impact":"high","confidence":"high","rationale":"This recommendation addresses the highest-ranked finding."}],"findingCoverage":[{"findingRef":` + findingRef + `,"disposition":"addressed_by_recommendation","recommendationRefs":["rec-001"],"rationale":"The recommendation is traced to this finding."}],"rationale":"Recommendations were ranked by expected quality impact."}`,
	}
}

func noFindingAreaAdvicePayloads(areaID string) []string {
	traceRef := `{"kind":"AreaAnalysisResult","subject":{"areaId":"` + areaID + `"},"selector":"localAndDescendantAnalysis"}`
	return noFindingAdvicePayloads(traceRef)
}

func noFindingFactorAdvicePayloads(factorID string) []string {
	traceRef := `{"kind":"FactorAnalysisResult","subject":{"factorId":"` + factorID + `"},"selector":"localAndDescendantAnalysis"}`
	return noFindingAdvicePayloads(traceRef)
}

func noFindingAdvicePayloads(traceRef string) []string {
	return []string{
		`{"schemaVersion":3,"kind":"FindingRankingResult","orderedFindings":[],"rationale":"No findings were produced; no finding ranking is applicable."}`,
		`{"schemaVersion":3,"kind":"RecommendationResult","id":"rec-001","title":"Review the next quality bar","description":"Review the analyzed area or factor against the next intended quality bar.","background":"The evaluation met the current bar, so the next useful step is deciding whether the quality bar should rise or be clarified.","expectedValue":"The quality model remains current without inventing remediation work.","doneCriterion":"The review either confirms the current bar or records a clearer next bar.","impact":"medium","confidence":"medium","traceRefs":[` + traceRef + `]}`,
		`{"schemaVersion":3,"kind":"RecommendationRankingResult","orderedRecommendations":[{"rank":1,"recommendationRef":"rec-001","impact":"medium","confidence":"medium","rationale":"This is the only recommendation and keeps the evaluation actionable."}],"findingCoverage":[],"rationale":"No finding coverage entries are needed because the evaluation produced no findings."}`,
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
