package evaluation

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
)

// ReportDocument is the JSON report document derived from an
// evaluation run.
type ReportDocument struct {
	SchemaVersion     int                       `json:"schemaVersion"`
	Summary           ReportSummary             `json:"summary"`
	RatingResult      RatingResult              `json:"ratingResult"`
	Scope             ReportScope               `json:"scope"`
	EvidenceBasis     []ReportEvidence          `json:"evidenceBasis"`
	Limitations       []string                  `json:"limitations"`
	NextStep          ReportNextStep            `json:"nextAction"`
	AreaSummary       []AreaRatingSummary       `json:"areaSummary"`
	Areas             []AreaEvaluationDetail    `json:"areas"`
	AssessmentResults []AssessmentResultDigest  `json:"assessmentResults"`
	FindingSummaries  []FindingDigest           `json:"findingSummaries"`
	Recommendations   []RecommendationReference `json:"recommendations"`
}

// ReportJSON is the JSON report document contract.
type ReportJSON = ReportDocument

// ReportSummary captures the headline report metadata and rating.
type ReportSummary struct {
	Run          string       `json:"run,omitempty"`
	RootArea     string       `json:"rootArea"`
	Level        Level        `json:"evaluationLevel"`
	Rigor        Rigor        `json:"rigor"`
	Narrowing    string       `json:"narrowing"`
	RatingResult RatingResult `json:"ratingResult"`
}

// ReportScope describes recorded in-scope, out-of-scope, and missing scope
// metadata.
type ReportScope struct {
	Recorded    bool              `json:"recorded"`
	Description string            `json:"description"`
	Narrowing   string            `json:"narrowing"`
	InScope     []string          `json:"inScope"`
	OutOfScope  []string          `json:"outOfScope"`
	NotRecorded []MissingMetadata `json:"notRecorded"`
}

// ReportEvidence names evidence documents used to derive a report.
type ReportEvidence struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

// ReportNextStep summarizes the next recommended report action.
type ReportNextStep struct {
	Kind               ReportNextStepKind `json:"kind"`
	Summary            string             `json:"summary"`
	RecommendationID   string             `json:"recommendationId,omitempty"`
	RecommendationPath string             `json:"recommendationPath,omitempty"`
}

// AreaRatingSummary summarizes rating state for one area. It is the canonical
// compact Area breakdown layer: it carries the Factor rating results so the
// summary-level Factor breakdown is recoverable without joining the detailed
// Areas array.
type AreaRatingSummary struct {
	AreaPath                        AreaPath             `json:"areaPath"`
	AreaRatingState                 AreaRatingState      `json:"areaRatingState"`
	AreaRatingResult                *RatingResult        `json:"areaRatingResult"`
	AreaWithDescendantsRatingResult RatingResult         `json:"areaWithDescendantsRatingResult"`
	FactorRatingResults             []FactorRatingResult `json:"factorRatingResults"`
	CoveredRequirements             int                  `json:"coveredRequirements"`
	Note                            string               `json:"note,omitempty"`
}

// AreaEvaluationDetail contains the detailed report entry for one area. The
// structural/area-group state lives only in AreaRatingState; there is no
// parallel boolean re-encoding it.
type AreaEvaluationDetail struct {
	AreaPath                        AreaPath             `json:"areaPath"`
	AreaRatingState                 AreaRatingState      `json:"areaRatingState"`
	AreaRatingResult                *RatingResult        `json:"areaRatingResult"`
	AreaWithDescendantsRatingResult RatingResult         `json:"areaWithDescendantsRatingResult"`
	FactorRatingResults             []FactorRatingResult `json:"factorRatingResults"`
	AnalysisRecord                  string               `json:"analysisRecord"`
	NotAssessedRequirements         []string             `json:"notAssessedRequirements"`
}

// AssessmentResultDigest summarizes an assessment result included in a report.
type AssessmentResultDigest struct {
	AssessmentResultRecord string               `json:"assessmentResultRecord"`
	AreaPath               AreaPath             `json:"areaPath"`
	Requirement            string               `json:"requirement"`
	RatingResult           RatingResult         `json:"ratingResult"`
	State                  RecordLifecycleState `json:"state"`
	Active                 bool                 `json:"active"`
	Supersedes             []string             `json:"supersedes,omitempty"`
}

// FindingDigest summarizes a finding for report display.
type FindingDigest struct {
	AssessmentResultRecord string          `json:"assessmentResultRecord"`
	Locator                string          `json:"locator"`
	Category               string          `json:"category"`
	Severity               FindingSeverity `json:"severity"`
	Summary                string          `json:"summary"`
}

// RecommendationReference summarizes a recommendation linked from a report.
type RecommendationReference struct {
	ID            string               `json:"id"`
	Path          string               `json:"path"`
	DoneCriterion string               `json:"doneCriterion"`
	State         RecordLifecycleState `json:"state"`
	Active        bool                 `json:"active"`
	Supersedes    []string             `json:"supersedes,omitempty"`
}

// BuildReportReceipt is the JSON contract emitted after building report files.
type BuildReportReceipt struct {
	SchemaVersion   int              `json:"schemaVersion"`
	Path            string           `json:"path"`
	ReportSummaryMD string           `json:"reportSummaryMd"`
	ReportMD        string           `json:"reportMd"`
	ReportJSON      string           `json:"reportJson"`
	RatingResult    RatingResult     `json:"ratingResult"`
	NextActions     []receipt.Action `json:"nextActions,omitempty"`
}

// GateReceipt is the JSON contract emitted by report gate checks.
type GateReceipt struct {
	SchemaVersion int          `json:"schemaVersion"`
	Path          string       `json:"path"`
	Threshold     string       `json:"threshold"`
	RatingResult  RatingResult `json:"ratingResult"`
	Pass          bool         `json:"pass"`
}

// BuildReport renders report-summary.md, report.md, and report.json for a run.
func BuildReport(path string) (*BuildReportReceipt, error) {
	run, err := Inspect(path)
	if err != nil {
		return nil, err
	}
	if gaps := run.Renderable(); len(gaps) > 0 {
		return nil, nonReportableRunError(run.Path, gaps[0])
	}
	report, err := run.Report()
	if err != nil {
		return nil, err
	}
	labels := reportDisplayLabelsFromModel(run.Model)
	md := renderReportMarkdown(report, labels)
	summaryMD := renderReportSummaryMarkdown(report, labels)
	js, err := marshalJSON(report)
	if err != nil {
		return nil, err
	}
	reportSummaryAbs := filepath.Join(run.AbsPath, "report-summary.md")
	reportAbs := filepath.Join(run.AbsPath, "report.md")
	reportJSONAbs := filepath.Join(run.AbsPath, "report.json")
	if err := writeReportFile(reportSummaryAbs, summaryMD); err != nil {
		return nil, err
	}
	if err := writeReportFile(reportAbs, md); err != nil {
		return nil, err
	}
	if err := writeReportFile(reportJSONAbs, js); err != nil {
		return nil, err
	}
	return &BuildReportReceipt{
		SchemaVersion:   SchemaVersion,
		Path:            run.Path,
		ReportSummaryMD: filepath.ToSlash(filepath.Join(run.Path, "report-summary.md")),
		ReportMD:        filepath.ToSlash(filepath.Join(run.Path, "report.md")),
		ReportJSON:      filepath.ToSlash(filepath.Join(run.Path, "report.json")),
		RatingResult:    report.RatingResult,
	}, nil
}

func nonReportableRunError(runPath string, gap RunGap) error {
	return fmt.Errorf("run is not reportable: %s %s: %s (run `qualitymd evaluation status %s` for all gaps)", gap.Kind, gap.Ref, gap.Detail, runPath)
}

// Report assembles the evaluation report document from the run's analysis
// records. It returns an error when the run has no analysis records, no root
// analysis (one with an empty area path), or more than one root analysis. These
// are plain formatted errors describing the condition, not errors.Is-matchable
// sentinels.
func (r *Run) Report() (ReportDocument, error) {
	analyses, rootAnalysis, err := sortedReportAnalyses(r.Analyses)
	if err != nil {
		return ReportDocument{}, err
	}
	context := r.reportContext()
	report := newReportDocument(r.Path, rootAnalysis, context)
	collector := limitationCollector{seen: map[string]bool{}, out: &report.Limitations}
	collector.addContext(context, rootAnalysis)

	supersededAssessmentResults, _ := r.assessmentResultSupersedingState()
	assessmentResultsByFile := addAssessmentResultSummaries(&report, r.AssessmentResults, supersededAssessmentResults, &collector)
	addReportAreas(&report, analyses, assessmentResultsByFile, context, &collector)
	addReportRecommendations(&report, r.Recommendations, r.supersededRecommendations())
	setReportNextStep(&report)
	return report, nil
}

func sortedReportAnalyses(records []AnalysisRecord) ([]AnalysisRecord, AnalysisRecord, error) {
	analyses := append([]AnalysisRecord(nil), records...)
	if len(analyses) == 0 {
		return nil, AnalysisRecord{}, fmt.Errorf("run has no analysis records")
	}
	slices.SortFunc(analyses, func(a, b AnalysisRecord) int {
		if cmp := len(a.AreaPath) - len(b.AreaPath); cmp != 0 {
			return cmp
		}
		return strings.Compare(strings.Join(a.AreaPath, "/"), strings.Join(b.AreaPath, "/"))
	})
	var root *AnalysisRecord
	for _, analysis := range analyses {
		if len(analysis.AreaPath) == 0 {
			if root != nil {
				return nil, AnalysisRecord{}, fmt.Errorf("run has multiple root analysis records")
			}
			candidate := analysis
			root = &candidate
		}
	}
	if root == nil {
		return nil, AnalysisRecord{}, fmt.Errorf("run has no root analysis record")
	}
	return analyses, *root, nil
}

func newReportDocument(runPath string, rootAnalysis AnalysisRecord, context runContext) ReportDocument {
	rootRating := rootAnalysis.AggregateRatingResult
	report := ReportDocument{
		SchemaVersion: SchemaVersion,
		Summary: ReportSummary{
			Run:          filepath.Base(runPath),
			RootArea:     rootAnalysis.AreaPath.Display(),
			Level:        context.Level,
			Rigor:        context.Rigor,
			Narrowing:    context.Narrowing,
			RatingResult: rootRating,
		},
		RatingResult: rootRating,
		Scope: ReportScope{
			Recorded:    context.Recorded,
			Description: context.ScopeDescription,
			Narrowing:   context.Narrowing,
			InScope:     append([]string(nil), context.InScope...),
			OutOfScope:  append([]string(nil), context.OutOfScope...),
			NotRecorded: []MissingMetadata{},
		},
		EvidenceBasis:     []ReportEvidence{},
		Limitations:       []string{},
		AreaSummary:       []AreaRatingSummary{},
		Areas:             []AreaEvaluationDetail{},
		AssessmentResults: []AssessmentResultDigest{},
		FindingSummaries:  []FindingDigest{},
		Recommendations:   []RecommendationReference{},
	}
	normalizeReportScope(&report)
	return report
}

func normalizeReportScope(report *ReportDocument) {
	if report.Scope.Description == "" {
		report.Scope.Description = "Evaluation scope reconstructed from the run's analysis and assessment result records."
	}
	if report.Scope.InScope == nil {
		report.Scope.InScope = []string{}
	}
	if report.Scope.OutOfScope == nil {
		report.Scope.OutOfScope = []string{}
	}
	if report.Summary.Rigor == "" {
		report.Scope.NotRecorded = append(report.Scope.NotRecorded, missingMetadata(MissingMetadataRigor))
	}
	if len(report.Scope.OutOfScope) == 0 {
		report.Scope.NotRecorded = append(report.Scope.NotRecorded, missingMetadata(MissingMetadataOutOfScope))
	}
}

type limitationCollector struct {
	seen map[string]bool
	out  *[]string
}

func (c *limitationCollector) add(limitation string) {
	limitation = strings.TrimSpace(limitation)
	key := limitationKey(limitation)
	if key == "" || c.seen[key] {
		return
	}
	c.seen[key] = true
	*c.out = append(*c.out, limitation)
}

func (c *limitationCollector) addText(text string) {
	for _, limitation := range limitationSentences(text) {
		c.add(limitation)
	}
}

func (c *limitationCollector) addContext(context runContext, rootAnalysis AnalysisRecord) {
	for _, limitation := range context.Limitations {
		c.addText(limitation)
		c.add(limitation)
	}
	c.addText(rootAnalysis.AggregateRatingResult.Rationale)
}

func addAssessmentResultSummaries(report *ReportDocument, records []AssessmentResultRecord, superseded map[string]bool, collector *limitationCollector) map[string]AssessmentResultRecord {
	recordsByFile := map[string]AssessmentResultRecord{}
	evidenceSeen := map[string]bool{}
	for _, record := range records {
		active := !superseded[record.File]
		recordsByFile[record.File] = record
		report.AssessmentResults = append(report.AssessmentResults, AssessmentResultDigest{
			AssessmentResultRecord: record.File,
			AreaPath:               record.AreaPath,
			Requirement:            record.Requirement,
			RatingResult:           record.RatingResult,
			State:                  lifecycleState(active),
			Active:                 active,
			Supersedes:             append([]string(nil), record.Supersedes...),
		})
		if !active {
			continue
		}
		if isNotAssessed(record.RatingResult) {
			collector.addText(record.Requirement + ": " + record.RatingResult.Rationale)
		}
		collector.addText(record.RatingResult.Rationale)
		addReportFindings(report, record, evidenceSeen)
	}
	return recordsByFile
}

func addReportFindings(report *ReportDocument, record AssessmentResultRecord, evidenceSeen map[string]bool) {
	for _, finding := range record.Findings {
		report.FindingSummaries = append(report.FindingSummaries, FindingDigest{
			AssessmentResultRecord: record.File,
			Locator:                finding.Locator,
			Category:               finding.Category,
			Severity:               finding.Severity.Display(),
			Summary:                finding.Observation,
		})
		for _, evidence := range finding.Evidence {
			key := evidence.Kind + "\x00" + evidence.Ref
			if evidence.Ref == "" || evidenceSeen[key] {
				continue
			}
			evidenceSeen[key] = true
			report.EvidenceBasis = append(report.EvidenceBasis, ReportEvidence(evidence))
		}
	}
}

func addReportAreas(report *ReportDocument, analyses []AnalysisRecord, assessmentResultsByFile map[string]AssessmentResultRecord, context runContext, collector *limitationCollector) {
	for _, analysis := range analyses {
		structural := len(analysis.AssessmentResultRecords) == 0 && len(analysis.ChildAnalysisRecords) > 0
		if analysis.LocalRatingResult != nil {
			collector.addText(analysis.LocalRatingResult.Rationale)
		}
		collector.addText(analysis.AggregateRatingResult.Rationale)
		for _, factor := range analysis.FactorRatingResults {
			collector.addText(factor.RatingResult.Rationale)
		}
		area := areaEvaluationFromAnalysis(analysis, structural)
		for _, ref := range analysis.AssessmentResultRecords {
			if record, ok := assessmentResultsByFile[ref]; ok && isNotAssessed(record.RatingResult) {
				area.NotAssessedRequirements = append(area.NotAssessedRequirements, record.Requirement)
			}
		}
		report.AreaSummary = append(report.AreaSummary, areaEvaluationSummary(area, len(analysis.AssessmentResultRecords)))
		if len(context.InScope) == 0 {
			report.Scope.InScope = append(report.Scope.InScope, area.AreaPath.Display())
		}
		report.Areas = append(report.Areas, area)
	}
}

func areaEvaluationFromAnalysis(analysis AnalysisRecord, structural bool) AreaEvaluationDetail {
	area := AreaEvaluationDetail{
		AreaPath:                        analysis.AreaPath.Clone(),
		AnalysisRecord:                  analysis.File,
		AreaRatingState:                 areaRatingStateFromResult(analysis.LocalRatingResult),
		AreaRatingResult:                cloneRatingResult(analysis.LocalRatingResult),
		AreaWithDescendantsRatingResult: analysis.AggregateRatingResult,
		FactorRatingResults:             nonNilFactorRatingResults(analysis.FactorRatingResults),
		NotAssessedRequirements:         []string{},
	}
	if structural {
		area.AreaRatingResult = nil
		area.AreaRatingState = areaRatingStateFromResult(nil)
	}
	return area
}

func areaEvaluationSummary(area AreaEvaluationDetail, coveredRequirements int) AreaRatingSummary {
	note := ""
	if len(area.NotAssessedRequirements) > 0 {
		note = "has not-assessed requirements"
	}
	return AreaRatingSummary{
		AreaPath:                        area.AreaPath.Clone(),
		AreaRatingState:                 area.AreaRatingState,
		AreaRatingResult:                cloneRatingResult(area.AreaRatingResult),
		AreaWithDescendantsRatingResult: area.AreaWithDescendantsRatingResult,
		FactorRatingResults:             cloneFactorRatingResults(area.FactorRatingResults),
		CoveredRequirements:             coveredRequirements,
		Note:                            note,
	}
}

func addReportRecommendations(report *ReportDocument, recommendations []RecommendationRecord, superseded map[string]bool) {
	for _, rec := range recommendations {
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		report.Recommendations = append(report.Recommendations, RecommendationReference{
			ID:            id,
			Path:          rec.File,
			DoneCriterion: rec.DoneCriterion,
			State:         lifecycleState(!superseded[rec.File] && !superseded[id]),
			Active:        !superseded[rec.File] && !superseded[id],
			Supersedes:    append([]string(nil), rec.Supersedes...),
		})
	}
}

func setReportNextStep(report *ReportDocument) {
	if first, ok := firstActiveRecommendation(report.Recommendations); ok {
		report.NextStep = ReportNextStep{
			Kind:               ReportNextStepRecommendation,
			Summary:            first.DoneCriterion,
			RecommendationID:   first.ID,
			RecommendationPath: first.Path,
		}
	} else {
		report.NextStep = ReportNextStep{
			Kind:    ReportNextStepNone,
			Summary: "No recommendation records exist for this run.",
		}
	}
}

func renderReportMarkdown(report ReportDocument, labels reportDisplayLabels) []byte {
	var out bytes.Buffer
	out.WriteString("# Evaluation Report\n\n")
	writeReportSummarySection(&out, report, labels)
	writeReportScopeSection(&out, report, labels)
	writeReportRisksAndLimitationsSection(&out, report)
	writeReportEvidenceSection(&out, report)
	writeReportNextStepSection(&out, report)
	writeAreaBreakdownSection(&out, report, labels)
	writeReportAreaDetailsSection(&out, report, labels)
	writeReportRequirementsSection(&out, report, labels)
	writeReportFindingsSection(&out, report)
	writeReportAdviceSection(&out, report)
	return out.Bytes()
}

func writeReportSummarySection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Verdict\n\n")
	out.WriteString("- **Root area:** " + labels.Area(nil, report.Summary.RootArea) + "\n")
	if report.Summary.Level != "" {
		out.WriteString("- **Evaluation level:** " + string(report.Summary.Level) + "\n")
	} else {
		out.WriteString("- **Evaluation level:** not recorded\n")
	}
	if report.Summary.Rigor != "" {
		out.WriteString("- **Rigor:** " + report.Summary.Rigor.Display() + "\n")
	} else {
		out.WriteString("- **Rigor:** not recorded\n")
	}
	if report.Summary.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Summary.Narrowing + "\n")
	}
	out.WriteString("- **Evaluation verdict:** " + displayRatingResult(report.Summary.RatingResult, labels.Ratings) + "\n")
	if report.Summary.RatingResult.Rationale != "" {
		out.WriteString("- **Rationale:** " + report.Summary.RatingResult.Rationale + "\n")
	}
}

func writeReportScopeSection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Scope\n\n")
	out.WriteString(report.Scope.Description + "\n\n")
	if report.Scope.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Scope.Narrowing + "\n")
	} else {
		out.WriteString("- **Narrowing:** whole recorded run\n")
	}
	out.WriteString("- **In scope:** " + displayScopeList(report.Scope.InScope, report.AreaSummary, labels) + "\n")
	if len(report.Scope.OutOfScope) > 0 {
		out.WriteString("- **Out of scope:** " + displayList(report.Scope.OutOfScope) + "\n")
	} else {
		out.WriteString("- **Out of scope:** not recorded\n")
	}
	if len(report.Scope.NotRecorded) > 0 {
		out.WriteString("- **Metadata not recorded:** " + displayMissingMetadata(report.Scope.NotRecorded) + "\n")
	}
}

func writeReportRisksAndLimitationsSection(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Selected Findings and Limitations\n\n")
	risks := riskFindings(report.FindingSummaries)
	summaryRisks := firstFindingSummaries(risks, 8)
	summaryLimitations := firstStrings(report.Limitations, 8)
	if len(summaryRisks) == 0 && len(summaryLimitations) == 0 {
		out.WriteString("No selected findings or limitations were recorded in the report data.\n")
		return
	}
	for _, finding := range summaryRisks {
		out.WriteString("- `" + finding.AssessmentResultRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		out.WriteString(" [" + finding.Severity.Title + "]")
		out.WriteString(": " + finding.Summary + "\n")
	}
	for _, limitation := range summaryLimitations {
		out.WriteString("- Limitation: " + limitation + "\n")
	}
	if len(risks) > len(summaryRisks) || len(report.Limitations) > len(summaryLimitations) {
		out.WriteString("- Additional selected findings or limitations are available in `report.json`.\n")
	}
}

func writeReportEvidenceSection(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Evidence Basis\n\n")
	if len(report.EvidenceBasis) == 0 {
		out.WriteString("No command or source evidence basis was recorded in findings.\n")
	} else {
		for _, evidence := range report.EvidenceBasis {
			out.WriteString("- **" + evidence.Kind + ":** `" + evidence.Ref + "`\n")
		}
	}
}

func writeReportNextStepSection(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Next Action\n\n")
	if report.NextStep.RecommendationID != "" {
		out.WriteString("- [" + report.NextStep.RecommendationID + "](" + report.NextStep.RecommendationPath + ") - " + report.NextStep.Summary + "\n")
	} else {
		out.WriteString(report.NextStep.Summary + "\n")
	}
}

// writeAreaBreakdownSection renders the shared compact Area Breakdown table used
// by both report.md and report-summary.md so the two cannot drift.
func writeAreaBreakdownSection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Area Breakdown\n\n")
	if len(report.AreaSummary) == 0 {
		out.WriteString("No area ratings were recorded.\n")
		return
	}
	out.WriteString("| Path | Area | + Sub-Areas | Factors |\n")
	out.WriteString("| --- | --- | --- | --- |\n")
	for _, area := range report.AreaSummary {
		out.WriteString(areaBreakdownRow(area, labels))
	}
}

// areaBreakdownRow renders one compact Area Breakdown row. The Area column shows
// the Area-only rating (area-group state renders as `(area group)`); the
// + Sub-Areas column shows the Area-with-descendants rating; the Factors column
// lists each factor as `<factor display path>: <rating>`.
func areaBreakdownRow(area AreaRatingSummary, labels reportDisplayLabels) string {
	return "| " + tableCell(labels.AreaDisplayPath(area.AreaPath.Elements())) +
		" | " + tableCell(displayAreaRatingState(area.AreaRatingState, labels.Ratings)) +
		" | " + tableCell(displayRatingResult(area.AreaWithDescendantsRatingResult, labels.Ratings)) +
		" | " + tableCell(areaBreakdownFactors(area, labels)) + " |\n"
}

// areaBreakdownFactors renders the Factors cell: a "; "-joined list of
// `<factor display path>: <rating>`. Areas with no recorded factor ratings
// render an explicit empty-state token so an omission cannot read as accidental.
func areaBreakdownFactors(area AreaRatingSummary, labels reportDisplayLabels) string {
	if len(area.FactorRatingResults) == 0 {
		return "(no factor ratings)"
	}
	parts := make([]string, 0, len(area.FactorRatingResults))
	for _, factor := range area.FactorRatingResults {
		label := labels.FactorDisplayPath(area.AreaPath.Elements(), factor.FactorPath.Elements())
		parts = append(parts, label+": "+displayRatingResult(factor.RatingResult, labels.Ratings))
	}
	return strings.Join(parts, "; ")
}

func writeReportAreaDetailsSection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Area Details\n\n")
	for _, area := range report.Areas {
		writeReportAreaDetail(out, area, labels)
	}
}

func writeReportAreaDetail(out *bytes.Buffer, area AreaEvaluationDetail, labels reportDisplayLabels) {
	out.WriteString("### " + labels.Area(area.AreaPath.Elements(), area.AreaPath.Display()) + "\n\n")
	out.WriteString("- **Path:** " + displayPath(area.AreaPath) + "\n")
	out.WriteString("- **Area rating:** " + displayAreaRatingState(area.AreaRatingState, labels.Ratings) + "\n")
	if area.AreaRatingState.RatingResult != nil {
		writeOptionalRationale(out, area.AreaRatingState.RatingResult.Rationale)
	}
	out.WriteString("- **+ Sub-Areas rating:** " + displayRatingResult(area.AreaWithDescendantsRatingResult, labels.Ratings) + "\n")
	writeOptionalRationale(out, area.AreaWithDescendantsRatingResult.Rationale)
	for _, factor := range area.FactorRatingResults {
		out.WriteString("- **Factor " + labels.FactorDisplayPath(area.AreaPath.Elements(), factor.FactorPath.Elements()) + ":** " + displayRatingResult(factor.RatingResult, labels.Ratings) + "\n")
		writeOptionalRationale(out, factor.RatingResult.Rationale)
	}
	out.WriteString("- **Analysis record:** `" + area.AnalysisRecord + "`\n")
	if len(area.NotAssessedRequirements) > 0 {
		out.WriteString("- **Not assessed:** " + strings.Join(area.NotAssessedRequirements, "; ") + "\n")
	}
	out.WriteString("\n")
}

func writeOptionalRationale(out *bytes.Buffer, rationale string) {
	if rationale != "" {
		out.WriteString("  - " + rationale + "\n")
	}
}

func writeReportRequirementsSection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Requirements\n\n")
	for _, result := range report.AssessmentResults {
		out.WriteString("### " + result.Requirement + "\n\n")
		out.WriteString("- **State:** " + string(result.State) + "\n")
		out.WriteString("- **Area:** " + labels.Area(result.AreaPath.Elements(), result.AreaPath.Display()) + "\n")
		out.WriteString("- **Rating:** " + displayRatingResult(result.RatingResult, labels.Ratings) + "\n")
		out.WriteString("- **Assessment result record:** `" + result.AssessmentResultRecord + "`\n")
		if result.RatingResult.Rationale != "" {
			out.WriteString("- **Rationale:** " + result.RatingResult.Rationale + "\n")
		}
		out.WriteString("\n")
	}
}

func writeReportFindingsSection(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("## Findings\n\n")
	for _, finding := range report.FindingSummaries {
		out.WriteString("- `" + finding.AssessmentResultRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		out.WriteString(": " + finding.Summary + "\n")
	}
}

func writeReportAdviceSection(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Advice\n\n")
	if len(report.Recommendations) == 0 {
		out.WriteString("No recommendation records exist for this run.\n")
	} else {
		for _, rec := range report.Recommendations {
			out.WriteString("- [" + rec.ID + "](" + rec.Path + ") [" + string(rec.State) + "] - " + rec.DoneCriterion + "\n")
		}
	}
}

func renderReportSummaryMarkdown(report ReportDocument, labels reportDisplayLabels) []byte {
	var out bytes.Buffer
	out.WriteString("# Quality Evaluation Summary\n\n")
	writeSummaryKeyDetails(&out, report, labels)
	writeSummarySection(&out, report, labels)
	writeAreaBreakdownSection(&out, report, labels)
	writeSummaryTopIssues(&out, report)
	writeSummaryRecommendations(&out, report)
	writeSummaryScopeAndLimitations(&out, report, labels)
	return out.Bytes()
}

func writeSummaryKeyDetails(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("| Field | Value |\n")
	out.WriteString("| --- | --- |\n")
	out.WriteString("| Root area | " + tableCell(labels.Area(nil, report.Summary.RootArea)) + " |\n")
	if report.Summary.Run != "" {
		out.WriteString("| Run | `" + tableCell(report.Summary.Run) + "` |\n")
	}
	out.WriteString("| Scope | " + tableCell(summaryScope(report)) + " |\n")
	out.WriteString("| Rigor | " + tableCell(report.Summary.Rigor.Display()) + " |\n")
	out.WriteString("| Evaluation verdict | " + tableCell(displayRatingResult(report.Summary.RatingResult, labels.Ratings)) + " |\n")
	out.WriteString("| Full report | [report.md](report.md) |\n")
	out.WriteString("| Machine report | [report.json](report.json) |\n\n")
}

func writeSummarySection(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Verdict\n\n")
	if report.Summary.RatingResult.Rationale != "" {
		out.WriteString(report.Summary.RatingResult.Rationale + "\n")
	} else if isNotAssessed(report.Summary.RatingResult) {
		out.WriteString("The evaluation did not produce a rated verdict.\n")
	} else {
		out.WriteString("The evaluation completed with verdict " + displayRatingResult(report.Summary.RatingResult, labels.Ratings) + ".\n")
	}
}

func writeSummaryTopIssues(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Selected Findings\n\n")
	issues := firstFindingSummaries(riskFindings(report.FindingSummaries), 5)
	if len(issues) == 0 {
		out.WriteString("None recorded.\n")
		return
	}
	for i, finding := range issues {
		fmt.Fprintf(out, "%d. ", i+1)
		out.WriteString("**" + finding.Severity.Title + "**  \n   ")
		out.WriteString(finding.Summary)
		out.WriteString("\n")
		if finding.Locator != "" {
			out.WriteString("   `" + finding.Locator + "`\n")
		}
		if finding.AssessmentResultRecord != "" {
			out.WriteString("   Assessment: `" + finding.AssessmentResultRecord + "`\n")
		}
	}
}

func writeSummaryRecommendations(out *bytes.Buffer, report ReportDocument) {
	out.WriteString("\n## Recommended Actions\n\n")
	active := activeRecommendations(report.Recommendations)
	if report.NextStep.RecommendationID != "" {
		out.WriteString("Primary next action: use `" + report.NextStep.RecommendationID + "`.\n\n")
	} else {
		out.WriteString(report.NextStep.Summary + "\n")
	}
	if len(active) == 0 {
		return
	}
	out.WriteString("| Recommendation ID | Priority | Recommendation | Done criterion |\n")
	out.WriteString("| --- | --- | --- | --- |\n")
	for i, rec := range active {
		out.WriteString("| `" + tableCell(rec.ID) + "` | " + fmt.Sprintf("%d", i+1) + " | [" + tableCell(recommendationTitle(rec.ID)) + "](" + rec.Path + ") | " + tableCell(rec.DoneCriterion) + " |\n")
	}
}

func recommendationTitle(id string) string {
	parts := strings.Split(id, "-")
	if len(parts) > 1 && isDigits(parts[0]) {
		parts = parts[1:]
	}
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, " ")
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func writeSummaryScopeAndLimitations(out *bytes.Buffer, report ReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Scope & Limitations\n\n")
	out.WriteString("Scope: **" + summaryScope(report) + "**\n")
	inScope := displayScopeList(report.Scope.InScope, report.AreaSummary, labels)
	if inScope != "" && inScope != "not recorded" {
		out.WriteString("\nIn scope: " + inScope + "\n")
	}
	limitations := firstStrings(report.Limitations, 5)
	if len(limitations) == 0 {
		out.WriteString("\nLimitations: none recorded.\n")
		return
	}
	out.WriteString("\n")
	for _, limitation := range limitations {
		out.WriteString("- " + limitation + "\n")
	}
	if len(report.Limitations) > len(limitations) {
		out.WriteString("- Additional limitations are available in [report.md](report.md).\n")
	}
}

func summaryScope(report ReportDocument) string {
	if report.Summary.Narrowing != "" {
		return report.Summary.Narrowing
	}
	if report.Scope.Narrowing != "" {
		return report.Scope.Narrowing
	}
	if isFullEvaluationScope(report.Scope.Description) {
		return "Full evaluation"
	}
	if report.Scope.Description != "" {
		return report.Scope.Description
	}
	return "Full evaluation"
}

func isFullEvaluationScope(value string) bool {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "", "whole model", "whole root area", "whole recorded run", "full evaluation", "evaluation scope reconstructed from the run's analysis and assessment result records.":
		return true
	default:
		return false
	}
}

func activeRecommendations(recommendations []RecommendationReference) []RecommendationReference {
	active := []RecommendationReference{}
	for _, rec := range recommendations {
		if rec.State.Active() {
			active = append(active, rec)
		}
	}
	return active
}

func (r *Run) supersededRecommendations() map[string]bool {
	known := map[string]string{}
	superseded := map[string]bool{}
	for _, rec := range r.Recommendations {
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		for _, ref := range rec.Supersedes {
			ref = strings.TrimSpace(ref)
			if file, ok := known[ref]; ok {
				superseded[file] = true
				superseded[strings.TrimSuffix(filepath.Base(file), ".md")] = true
				continue
			}
			if strings.HasPrefix(ref, "recommendations/") {
				id := strings.TrimSuffix(filepath.Base(ref), ".md")
				if file, ok := known[id]; ok {
					superseded[file] = true
					superseded[id] = true
				}
				continue
			}
			path := "recommendations/" + ref + ".md"
			if file, ok := known[path]; ok {
				superseded[file] = true
				superseded[ref] = true
			}
		}
		known[id] = rec.File
		known[rec.File] = rec.File
	}
	return superseded
}

func firstActiveRecommendation(recommendations []RecommendationReference) (RecommendationReference, bool) {
	for _, rec := range recommendations {
		if rec.State.Active() {
			return rec, true
		}
	}
	return RecommendationReference{}, false
}

type runContext struct {
	Recorded         bool
	Level            Level
	Narrowing        string
	Rigor            Rigor
	ScopeDescription string
	InScope          []string
	OutOfScope       []string
	Limitations      []string
}

func (r *Run) reportContext() runContext {
	context := inferRunContext(r.Path)
	designParams := markdownBulletMap(markdownSection(r.Design, "Resolved parameters"))
	if value := designParams["Altitude"]; value != "" {
		context.Level = Level(cleanInlineCode(value))
		context.Recorded = true
	}
	if value := designParams["Narrowing slug"]; value != "" {
		context.Narrowing = cleanInlineCode(value)
		context.Recorded = true
	} else if value := designParams["Narrowing"]; value != "" {
		context.Narrowing = cleanInlineCode(value)
		context.Recorded = true
	}
	if value := designParams["Rigor"]; value != "" {
		context.Rigor = Rigor(strings.ToLower(cleanInlineCode(value)))
		context.Recorded = true
	}
	if value := designParams["Scope"]; value != "" {
		context.ScopeDescription = cleanInlineCode(value)
		context.Recorded = true
	} else if value := designParams["Scope description"]; value != "" {
		context.ScopeDescription = cleanInlineCode(value)
		context.Recorded = true
	}
	if rigor := firstParagraph(markdownSection(r.Plan, "Rigor")); rigor != "" {
		context.Recorded = true
		if context.Rigor == "" {
			context.Rigor = inferRigor(strings.ToLower(rigor))
		}
		if context.ScopeDescription == "" {
			context.ScopeDescription = rigor
		}
	}
	if limitations := markdownListItems(markdownSection(r.Plan, "Planned limitations")); len(limitations) > 0 {
		context.Limitations = limitations
		context.OutOfScope = append([]string(nil), limitations...)
		context.Recorded = true
	}
	for _, doc := range []string{r.Design, r.Plan} {
		for _, heading := range []string{"Out of scope", "Deferred areas"} {
			if items := markdownListItems(markdownSection(doc, heading)); len(items) > 0 {
				context.OutOfScope = appendUniqueStrings(context.OutOfScope, items...)
				context.Recorded = true
			}
		}
	}
	return context
}

func appendUniqueStrings(items []string, more ...string) []string {
	seen := map[string]bool{}
	for _, item := range items {
		key := strings.ToLower(strings.TrimSpace(item))
		if key != "" {
			seen[key] = true
		}
	}
	for _, item := range more {
		item = strings.TrimSpace(item)
		key := strings.ToLower(item)
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		items = append(items, item)
	}
	return items
}

func inferRunContext(path string) runContext {
	base := filepath.Base(path)
	name := strings.TrimSuffix(base, "-quality-eval")
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 {
		return runContext{}
	}
	context := runContext{}
	suffix := parts[1]
	if suffix == "model" {
		context.Level = LevelModel
		return context
	}
	if strings.HasPrefix(suffix, "model-") {
		context.Level = LevelModel
		context.Narrowing = strings.TrimPrefix(suffix, "model-")
		context.Rigor = inferRigor(context.Narrowing)
		return context
	}
	if strings.HasPrefix(suffix, "subject-") {
		context.Narrowing = strings.TrimPrefix(suffix, "subject-")
		context.Rigor = inferRigor(context.Narrowing)
		return context
	}
	if suffix != "subject" {
		context.Narrowing = suffix
		context.Rigor = inferRigor(suffix)
	}
	return context
}

func inferRigor(narrowing string) Rigor {
	for _, rigor := range []Rigor{RigorQuick, RigorStandard, RigorDeep} {
		if strings.Contains(narrowing, string(rigor)) {
			return rigor
		}
	}
	return ""
}

func markdownSection(doc, heading string) string {
	lines := strings.Split(doc, "\n")
	area := "## " + heading
	var out []string
	inSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, area) {
			inSection = true
			continue
		}
		if inSection && strings.HasPrefix(trimmed, "## ") {
			break
		}
		if inSection {
			out = append(out, line)
		}
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func markdownBulletMap(section string) map[string]string {
	values := map[string]string{}
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "- ") {
			continue
		}
		key, value, ok := strings.Cut(strings.TrimPrefix(trimmed, "- "), ":")
		if !ok {
			continue
		}
		values[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return values
}

func markdownListItems(section string) []string {
	var items []string
	for _, line := range strings.Split(section, "\n") {
		trimmedRight := strings.TrimRight(line, " \t")
		trimmed := strings.TrimSpace(trimmedRight)
		if !strings.HasPrefix(trimmed, "- ") {
			if len(items) > 0 && trimmed != "" && (strings.HasPrefix(trimmedRight, " ") || strings.HasPrefix(trimmedRight, "\t")) {
				items[len(items)-1] = strings.TrimSpace(items[len(items)-1] + " " + cleanInlineCode(trimmed))
			}
			continue
		}
		item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
		if item != "" {
			items = append(items, cleanInlineCode(item))
		}
	}
	return items
}

func firstParagraph(section string) string {
	var lines []string
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(lines) > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "#") {
			if len(lines) > 0 {
				break
			}
			continue
		}
		lines = append(lines, trimmed)
	}
	return cleanInlineCode(strings.Join(lines, " "))
}

func cleanInlineCode(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Trim(value, "`")
	return strings.TrimSpace(value)
}

func cloneRatingResult(result *RatingResult) *RatingResult {
	if result == nil {
		return nil
	}
	clone := *result
	return &clone
}

func nonNilFactorRatingResults(factors []FactorRatingResult) []FactorRatingResult {
	if factors == nil {
		return []FactorRatingResult{}
	}
	return factors
}

func cloneFactorRatingResults(factors []FactorRatingResult) []FactorRatingResult {
	out := make([]FactorRatingResult, 0, len(factors))
	for _, factor := range factors {
		out = append(out, FactorRatingResult{
			FactorPath:   factor.FactorPath.Clone(),
			RatingResult: factor.RatingResult,
		})
	}
	return out
}

func riskFindings(findings []FindingDigest) []FindingDigest {
	risks := []FindingDigest{}
	for _, finding := range findings {
		if finding.Severity.Level.IsRisk() {
			risks = append(risks, finding)
		}
	}
	return risks
}

func firstFindingSummaries(items []FindingDigest, limit int) []FindingDigest {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func firstStrings(items []string, limit int) []string {
	if len(items) <= limit {
		return items
	}
	return items[:limit]
}

func limitationSentences(text string) []string {
	normalized := strings.Join(strings.Fields(text), " ")
	if normalized == "" {
		return nil
	}
	parts := strings.Split(normalized, ". ")
	out := []string{}
	for i, part := range parts {
		sentence := strings.TrimSpace(part)
		if sentence != "" && i < len(parts)-1 {
			sentence += "."
		}
		if sentence == "" || !looksLikeLimitation(sentence) {
			continue
		}
		out = append(out, strings.TrimSuffix(sentence, "."))
	}
	return out
}

func limitationKey(limitation string) string {
	limitation = strings.ToLower(strings.TrimSpace(limitation))
	limitation = strings.TrimRight(limitation, ".;:")
	limitation = strings.Join(strings.Fields(limitation), " ")
	return limitation
}

func looksLikeLimitation(text string) bool {
	lower := strings.ToLower(text)
	for _, keyword := range []string{
		"not executed",
		"not assessed",
		"not fully assessed",
		"inspected rather than executed",
		"inspected, not executed",
		"limited to",
		"limitation",
		"absent",
		"missing",
		"deferred",
		"out of scope",
	} {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

type reportDisplayLabels struct {
	Ratings map[string]string
	Areas   map[string]string
	Factors map[string]string
}

func reportDisplayLabelsFromModel(spec *model.Spec) reportDisplayLabels {
	labels := reportDisplayLabels{
		Ratings: map[string]string{},
		Areas:   map[string]string{},
		Factors: map[string]string{},
	}
	if spec == nil {
		return labels
	}
	labels.Ratings = ratingDisplayLabels(spec.RatingScale)
	if strings.TrimSpace(spec.Title) != "" {
		labels.Areas[areaLabelKey(nil)] = spec.Title
	}
	labels.addFactors(nil, spec.Factors)
	labels.addAreas(nil, spec.Areas)
	return labels
}

func (l reportDisplayLabels) addAreas(parentPath []string, areas map[string]model.Area) {
	for key, area := range areas {
		path := appendPathElement(parentPath, key)
		if strings.TrimSpace(area.Title) != "" {
			l.Areas[areaLabelKey(path)] = area.Title
		}
		l.addFactors(path, area.Factors)
		l.addAreas(path, area.Areas)
	}
}

func (l reportDisplayLabels) addFactors(areaPath []string, factors map[string]model.Factor) {
	for key, factor := range factors {
		path := appendPathElement(nil, key)
		if strings.TrimSpace(factor.Title) != "" {
			l.Factors[factorLabelKey(areaPath, path)] = factor.Title
		}
		l.addSubFactors(areaPath, path, factor.Factors)
	}
}

func (l reportDisplayLabels) addSubFactors(areaPath, parentFactorPath []string, factors map[string]model.Factor) {
	for key, factor := range factors {
		path := appendPathElement(parentFactorPath, key)
		if strings.TrimSpace(factor.Title) != "" {
			l.Factors[factorLabelKey(areaPath, path)] = factor.Title
		}
		l.addSubFactors(areaPath, path, factor.Factors)
	}
}

func (l reportDisplayLabels) Area(path []string, fallback string) string {
	if title := l.Areas[areaLabelKey(path)]; title != "" {
		return title
	}
	return fallback
}

func (l reportDisplayLabels) Factor(areaPath []string, factorPath []string) string {
	for i := len(areaPath); i >= 0; i-- {
		if title := l.Factors[factorLabelKey(areaPath[:i], factorPath)]; title != "" {
			return title
		}
	}
	return FactorPath(factorPath).Display()
}

// AreaDisplayPath renders the absolute display path for an area, resolving each
// path prefix to its title. The root area (empty path) renders as
// `/ (<root title>)`; descendants render slash-prefixed, such as
// `/Services/Payments/Webhooks`. The path is built from the element array, never
// by parsing a joined string.
func (l reportDisplayLabels) AreaDisplayPath(path []string) string {
	if len(path) == 0 {
		return "/ (" + l.Area(nil, "root") + ")"
	}
	var out strings.Builder
	for i := range path {
		out.WriteString("/")
		out.WriteString(l.Area(path[:i+1], path[i]))
	}
	return out.String()
}

// FactorDisplayPath renders the display path for a factor, resolving each factor
// prefix to its title and joining with " / ". The path is built from the
// element array, never by parsing a joined string.
func (l reportDisplayLabels) FactorDisplayPath(areaPath []string, factorPath []string) string {
	if len(factorPath) == 0 {
		return FactorPath(factorPath).Display()
	}
	parts := make([]string, 0, len(factorPath))
	for i := range factorPath {
		parts = append(parts, l.Factor(areaPath, factorPath[:i+1]))
	}
	return strings.Join(parts, " / ")
}

func ratingDisplayLabels(scale []model.RatingLevel) map[string]string {
	labels := map[string]string{}
	for _, level := range scale {
		if level.Level == "" || level.Title == "" {
			continue
		}
		labels[level.Level] = level.Title
	}
	return labels
}

func displayAreaRatingState(state AreaRatingState, ratingLabels map[string]string) string {
	switch state.Kind {
	case AreaRatingStructural:
		return "(area group)"
	case AreaRatingNotAssessed, AreaRatingRated:
		if state.RatingResult == nil {
			return "not assessed"
		}
		return displayRatingResult(*state.RatingResult, ratingLabels)
	default:
		return state.Title
	}
}

func displayRatingResult(result RatingResult, ratingLabels map[string]string) string {
	if isNotAssessed(result) {
		return "not assessed"
	}
	if title := ratingLabels[result.Level]; title != "" {
		return title
	}
	if result.Level == "" {
		return "not assessed"
	}
	return result.Level
}

func isNotAssessed(result RatingResult) bool {
	return result.Kind.IsNotAssessed()
}

func displayList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(items, "; ")
}

func displayMissingMetadata(items []MissingMetadata) string {
	if len(items) == 0 {
		return "none"
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item.Title != "" {
			out = append(out, item.Title)
			continue
		}
		out = append(out, string(item.Field))
	}
	return strings.Join(out, "; ")
}

func displayScopeList(items []string, areas []AreaRatingSummary, labels reportDisplayLabels) string {
	if len(items) == 0 {
		return "none"
	}
	areaLabels := map[string]string{}
	for _, area := range areas {
		id := area.AreaPath.Display()
		if _, exists := areaLabels[id]; exists {
			areaLabels[id] = ""
			continue
		}
		areaLabels[id] = labels.Area(area.AreaPath.Elements(), id)
	}
	display := make([]string, 0, len(items))
	for _, item := range items {
		if label := areaLabels[item]; label != "" {
			display = append(display, label)
			continue
		}
		display = append(display, item)
	}
	return strings.Join(display, "; ")
}

func displayPath(path AreaPath) string {
	if len(path) == 0 {
		return "(root)"
	}
	return strings.Join(path, " / ")
}

func areaLabelKey(path []string) string {
	return strings.Join(path, "\x00")
}

func factorLabelKey(areaPath []string, factorPath []string) string {
	return areaLabelKey(areaPath) + "\x00" + strings.Join(factorPath, "\x00")
}

func appendPathElement(path []string, value string) []string {
	out := make([]string, 0, len(path)+1)
	out = append(out, path...)
	out = append(out, value)
	return out
}

func tableCell(value string) string {
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\n", " ")
	if value == "" {
		return " "
	}
	return value
}

func writeReportFile(path string, data []byte) error {
	if _, err := os.Stat(path); err == nil {
		return writeReplace(path, data)
	} else if !os.IsNotExist(err) {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Gate reports whether a built report rating passes a threshold against the
// given ordered rating scale.
func Gate(result *BuildReportReceipt, scale []string, threshold string) (bool, error) {
	if threshold == "" {
		return true, nil
	}
	thresholdIndex := -1
	ratingIndex := -1
	for i, level := range scale {
		if level == threshold {
			thresholdIndex = i
		}
		if result.RatingResult.Level != "" && level == result.RatingResult.Level {
			ratingIndex = i
		}
	}
	if thresholdIndex < 0 {
		return false, usagef("--at-or-below level %q is not in the run rating scale", threshold)
	}
	if isNotAssessed(result.RatingResult) || result.RatingResult.Level == "" {
		return false, nil
	}
	if ratingIndex < 0 {
		return false, fmt.Errorf("evaluation verdict %q is not in the run rating scale", result.RatingResult.Level)
	}
	return ratingIndex < thresholdIndex, nil
}

// GateReport loads a built report and evaluates it against a threshold.
func GateReport(path, threshold string) (*GateReceipt, error) {
	run, err := Inspect(path)
	if err != nil {
		return nil, err
	}
	if gaps := run.Renderable(); len(gaps) > 0 {
		return nil, nonReportableRunError(run.Path, gaps[0])
	}
	reportPath := filepath.Join(run.Path, "report.json")
	raw, err := os.ReadFile(reportPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("report.json has not been built for %s", run.Path)
		}
		return nil, fmt.Errorf("reading report.json: %w", err)
	}
	var report ReportDocument
	if err := DecodeSingleJSON(raw, &report); err != nil {
		return nil, fmt.Errorf("reading report.json: %w", err)
	}
	levels := make([]string, 0, len(run.Scale))
	for _, level := range run.Scale {
		levels = append(levels, level.Level)
	}
	build := &BuildReportReceipt{RatingResult: report.RatingResult}
	pass, err := Gate(build, levels, threshold)
	if err != nil {
		return nil, err
	}
	return &GateReceipt{
		SchemaVersion: SchemaVersion,
		Path:          filepath.ToSlash(reportPath),
		Threshold:     threshold,
		RatingResult:  report.RatingResult,
		Pass:          pass,
	}, nil
}

// ScaleLevels returns the ordered rating scale levels for a run.
func ScaleLevels(path string) ([]string, error) {
	run, err := Inspect(path)
	if err != nil {
		return nil, err
	}
	levels := make([]string, 0, len(run.Scale))
	for _, level := range run.Scale {
		levels = append(levels, level.Level)
	}
	return levels, nil
}
