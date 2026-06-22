package evaluation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// nestedBreakdownModel exercises nested areas at depth >= 3 and nested factors.
// Factor map keys are lowercase; analysis factor paths reference them so the
// path-aware factor display resolves each prefix to its title.
const nestedBreakdownModel = `---
title: Atlas Commerce
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: 🟡 Minimum
    description: Minimum.
    criterion: Barely meets it.
  - level: unacceptable
    title: 🔴 Unacceptable
    description: Does not meet it.
    criterion: Fails it.
factors:
  trust:
    title: Trust
    description: Trust lens.
  operability:
    title: Operability
    description: Operability lens.
    factors:
      observability:
        title: Observability
        description: Observability lens.
areas:
  services:
    title: Services
    areas:
      fulfillment:
        title: Fulfillment
        areas:
          inventory:
            title: Inventory
            areas:
              reservations:
                title: Reservations
                requirements:
                  Reservations hold stock:
                    factors:
                      - trust
                    assessment: Inspect reservation holds.
  operations:
    title: Operations
    areas:
      incident-response:
        title: Incident Response
        requirements:
          Incidents are triaged:
            assessment: Inspect incident triage.
---
`

// addBreakdownAssessment writes one assessment result for the nested model.
func addBreakdownAssessment(t *testing.T, runPath, raw string) {
	t.Helper()
	if _, err := AddRecord(KindAssessmentResult, runPath, []byte(raw)); err != nil {
		t.Fatalf("AddRecord(assessment) error = %v", err)
	}
}

// addBreakdownAnalysis writes one analysis record for the nested model.
func addBreakdownAnalysis(t *testing.T, runPath, raw string) {
	t.Helper()
	if _, err := AddRecord(KindAnalysis, runPath, []byte(raw)); err != nil {
		t.Fatalf("AddRecord(analysis) error = %v", err)
	}
}

// buildNestedBreakdownRun assembles a reportable run over nestedBreakdownModel
// that exercises: a structural root area group; a depth-4 nested area with a
// nested factor; a rated mid-level area; a not-assessed area and factor; and an
// area with no factor ratings.
func buildNestedBreakdownRun(t *testing.T) string {
	t.Helper()
	repo := testRepoWithModel(t, nestedBreakdownModel)
	run, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	addBreakdownAssessment(t, runPath, `{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "requirement": "Reservations hold stock",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [["trust"]],
  "ratingResult": {"kind": "rated", "level": "minimum", "rationale": "Holds usually clear."}
}`)
	addBreakdownAssessment(t, runPath, `{
  "areaPath": ["operations", "incident-response"],
  "requirement": "Incidents are triaged",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [],
  "ratingResult": {"kind": "not-assessed", "rationale": "Incident triage was not assessed."}
}`)

	// Deepest area: nested factor path operability/observability + nested area path.
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "childAnalysisRecords": [],
  "localRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Reservation holds are minimum."},
  "factorRatingResults": [
    {"factorPath": ["trust"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Trust is target."}},
    {"factorPath": ["operability", "observability"], "ratingResult": {"kind": "rated", "level": "minimum", "rationale": "Observability is minimum."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Reservations roll up to minimum."},
  "assessmentResultRecords": ["assessments/001-services-fulfillment-inventory-reservations-reservations-hold-stock.json"]
}`)
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["services", "fulfillment", "inventory"],
  "childAnalysisRecords": ["analysis/services-fulfillment-inventory-reservations.json"],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Inventory rolls up to minimum."},
  "assessmentResultRecords": []
}`)
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["services", "fulfillment"],
  "childAnalysisRecords": ["analysis/services-fulfillment-inventory.json"],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Fulfillment rolls up to minimum."},
  "assessmentResultRecords": []
}`)
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["services"],
  "childAnalysisRecords": ["analysis/services-fulfillment.json"],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Services rolls up to minimum."},
  "assessmentResultRecords": []
}`)
	// Not-assessed area, with a not-assessed factor and another not-assessed factor.
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["operations", "incident-response"],
  "childAnalysisRecords": [],
  "localRatingResult": {"kind": "not-assessed", "rationale": "Incident response was not assessed."},
  "factorRatingResults": [
    {"factorPath": ["operability", "observability"], "ratingResult": {"kind": "not-assessed", "rationale": "Observability not assessed."}}
  ],
  "aggregateRatingResult": {"kind": "not-assessed", "rationale": "Incident response was not assessed."},
  "assessmentResultRecords": ["assessments/002-operations-incident-response-incidents-are-triaged.json"]
}`)
	// operations area group: child areas, no direct requirements, no factor ratings.
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["operations"],
  "childAnalysisRecords": ["analysis/operations-incident-response.json"],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {"kind": "not-assessed", "rationale": "Operations was not assessed."},
  "assessmentResultRecords": []
}`)
	// Root area group: structural, but carries roll-up + factor ratings.
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": [],
  "childAnalysisRecords": ["analysis/services.json", "analysis/operations.json"],
  "localRatingResult": null,
  "factorRatingResults": [
    {"factorPath": ["trust"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Trust is target."}},
    {"factorPath": ["operability"], "ratingResult": {"kind": "rated", "level": "minimum", "rationale": "Operability is minimum."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "The model rolls up to minimum."},
  "assessmentResultRecords": []
}`)
	return runPath
}

func TestAreaBreakdownRendersNestedPathsStatesAndEmptyFactors(t *testing.T) {
	runPath := buildNestedBreakdownRun(t)
	loaded, err := Load(runPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if status := loaded.Status(); !status.Reportable {
		t.Fatalf("status.Reportable = false, gaps = %#v", status.Gaps)
	}
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	summaryMD := readRunFile(t, runPath, "report-summary.md")
	reportMD := readRunFile(t, runPath, "report.md")

	for _, want := range []string{
		"## Area Breakdown",
		"| Area | Path | Area Rating | Area + Sub-Areas Rating | Factors |",
		// Structural root area group: Path column renders the display value; roll-up + factors preserved.
		"| Atlas Commerce | `/` | (area group) | 🟡 Minimum | Trust: 🔵 Target; Operability: 🟡 Minimum |",
		// Depth-4 absolute area path with a nested factor display path.
		"| Reservations | `services/fulfillment/inventory/reservations` | 🟡 Minimum | 🟡 Minimum | Trust: 🔵 Target; Operability / Observability: 🟡 Minimum |",
		// Not-assessed area and not-assessed factor render distinctly, never as a level.
		"| Incident Response | `operations/incident-response` | not assessed | not assessed | Operability / Observability: not assessed |",
		// Area group with no recorded factor ratings renders the explicit empty-state.
		"| Operations | `operations` | (area group) | not assessed | (no factor ratings) |",
	} {
		if !strings.Contains(summaryMD, want) {
			t.Fatalf("report-summary.md missing %q:\n%s", want, summaryMD)
		}
		if !strings.Contains(reportMD, want) {
			t.Fatalf("report.md missing shared breakdown row %q:\n%s", want, reportMD)
		}
	}
	// not-assessed must not leak as a rating-scale level title.
	if strings.Contains(summaryMD, "🟡 Minimum | not assessed | 🟡 Minimum") {
		t.Fatalf("not-assessed rendered as a level:\n%s", summaryMD)
	}
}

func TestAreaBreakdownJSONIdentifiersAndStructuralState(t *testing.T) {
	runPath := buildNestedBreakdownRun(t)
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	reportJSON := readRunFile(t, runPath, "report.json")

	for _, want := range []string{
		`"areaRatingState": {`,
		`"areaRatingResult": null`,
		`"areaWithDescendantsRatingResult": {`,
		`"factorRatingResults": [`,
		`"factorPath": [`,
		`"observability"`,
		`"reservations"`,
		`"level": "minimum"`,
	} {
		if !strings.Contains(reportJSON, want) {
			t.Fatalf("report.json missing %q:\n%s", want, reportJSON)
		}
	}
	for _, notWant := range []string{
		`"structural":`,
		`"localRating":`,
		`"localRatingResult":`,
		`"aggregateRatingResult":`,
	} {
		if strings.Contains(reportJSON, notWant) {
			t.Fatalf("report.json retains old field %q:\n%s", notWant, reportJSON)
		}
	}
	// factorPath/areaPath are element arrays, not joined display strings.
	if strings.Contains(reportJSON, `"factorPath": "operability/observability"`) {
		t.Fatalf("report.json joined a factor path instead of an element array:\n%s", reportJSON)
	}
}

func TestAreaBreakdownTitleFallbackToIdentifier(t *testing.T) {
	// A model with a titled root but an area lacking a title falls back to the
	// stable key for that area while still resolving the titled prefix.
	repo := testRepoWithModel(t, `---
title: Titled Root
ratingScale:
  - level: target
    title: 🔵 Target
    description: Target.
    criterion: Meets it.
factors:
  trust:
    title: Trust
    description: Trust lens.
areas:
  untitled-area:
    requirements:
      Does the thing:
        assessment: Inspect the thing.
---
`)
	run, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	addBreakdownAssessment(t, runPath, `{
  "areaPath": ["untitled-area"],
  "requirement": "Does the thing",
  "criterionSource": "rating-scale",
  "findings": [],
  "recommendations": [],
  "factorPaths": [["trust"]],
  "ratingResult": {"kind": "rated", "level": "target", "rationale": "Meets it."}
}`)
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": ["untitled-area"],
  "childAnalysisRecords": [],
  "localRatingResult": {"kind": "rated", "level": "target", "rationale": "Meets it."},
  "factorRatingResults": [
    {"factorPath": ["trust"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Trust is target."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "target", "rationale": "Meets it."},
  "assessmentResultRecords": ["assessments/001-untitled-area-does-the-thing.json"]
}`)
	addBreakdownAnalysis(t, runPath, `{
  "areaPath": [],
  "childAnalysisRecords": ["analysis/untitled-area.json"],
  "localRatingResult": null,
  "factorRatingResults": [],
  "aggregateRatingResult": {"kind": "rated", "level": "target", "rationale": "Meets it."},
  "assessmentResultRecords": []
}`)
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("BuildReport() error = %v", err)
	}
	summaryMD := readRunFile(t, runPath, "report-summary.md")
	if !strings.Contains(summaryMD, "| Titled Root | `/` | (area group) | 🔵 Target |") {
		t.Fatalf("root row missing titled fallback:\n%s", summaryMD)
	}
	if !strings.Contains(summaryMD, "| untitled-area | `untitled-area` | 🔵 Target | 🔵 Target | Trust: 🔵 Target |") {
		t.Fatalf("untitled area row did not fall back to the stable key:\n%s", summaryMD)
	}
}

func TestAreaBreakdownReportGenerationIsIdempotent(t *testing.T) {
	runPath := buildNestedBreakdownRun(t)
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("first BuildReport() error = %v", err)
	}
	first := map[string]string{}
	for _, name := range []string{"report-summary.md", "report.md", "report.json"} {
		first[name] = readRunFile(t, runPath, name)
	}
	if _, err := BuildReport(runPath); err != nil {
		t.Fatalf("second BuildReport() error = %v", err)
	}
	for _, name := range []string{"report-summary.md", "report.md", "report.json"} {
		if got := readRunFile(t, runPath, name); got != first[name] {
			t.Fatalf("%s changed across idempotent render", name)
		}
	}
}

func TestAnalysisValidationRejectsDuplicateAndUnresolvableFactorPaths(t *testing.T) {
	repo := testRepoWithModel(t, nestedBreakdownModel)
	run, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
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
			name: "duplicate factor path",
			raw: `{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "childAnalysisRecords": [],
  "localRatingResult": null,
  "factorRatingResults": [
    {"factorPath": ["trust"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Trust target."}},
    {"factorPath": ["trust"], "ratingResult": {"kind": "rated", "level": "minimum", "rationale": "Trust again."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Rolls up."},
  "assessmentResultRecords": []
}`,
			want: "duplicate within the analysis record",
		},
		{
			name: "unresolvable factor path",
			raw: `{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "childAnalysisRecords": [],
  "localRatingResult": null,
  "factorRatingResults": [
    {"factorPath": ["nonexistent"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Bogus."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Rolls up."},
  "assessmentResultRecords": []
}`,
			want: "does not resolve against the area's declared or inherited factor vocabulary",
		},
		{
			name: "unresolvable nested factor path",
			raw: `{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "childAnalysisRecords": [],
  "localRatingResult": null,
  "factorRatingResults": [
    {"factorPath": ["operability", "nope"], "ratingResult": {"kind": "rated", "level": "target", "rationale": "Bogus."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Rolls up."},
  "assessmentResultRecords": []
}`,
			want: "does not resolve against the area's declared or inherited factor vocabulary",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AddRecord(KindAnalysis, runPath, []byte(tc.raw))
			if err == nil {
				t.Fatalf("AddRecord() error = nil, want %q", tc.want)
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

func TestAnalysisValidationAcceptsInheritedFactorPath(t *testing.T) {
	repo := testRepoWithModel(t, nestedBreakdownModel)
	run, err := CreateRun(Options{RepoRoot: repo, Model: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)
	// operability/observability is inherited from the root model factor
	// vocabulary by a deeply nested area; it must be accepted.
	if _, err := AddRecord(KindAnalysis, runPath, []byte(`{
  "areaPath": ["services", "fulfillment", "inventory", "reservations"],
  "childAnalysisRecords": [],
  "localRatingResult": null,
  "factorRatingResults": [
    {"factorPath": ["operability", "observability"], "ratingResult": {"kind": "rated", "level": "minimum", "rationale": "Observability is minimum."}}
  ],
  "aggregateRatingResult": {"kind": "rated", "level": "minimum", "rationale": "Rolls up."},
  "assessmentResultRecords": []
}`)); err != nil {
		t.Fatalf("AddRecord(analysis) with inherited factor path error = %v", err)
	}
}

func readRunFile(t *testing.T, runPath, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(runPath, name))
	if err != nil {
		t.Fatalf("reading %s: %v", name, err)
	}
	return string(data)
}
