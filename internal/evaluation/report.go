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

type EvaluationReportDocument struct {
	SchemaVersion     int                       `json:"schemaVersion"`
	Summary           ReportSummary             `json:"summary"`
	RatingResult      RatingResult              `json:"ratingResult"`
	Scope             ReportScope               `json:"scope"`
	EvidenceBasis     []ReportEvidence          `json:"evidenceBasis"`
	Limitations       []string                  `json:"limitations"`
	NextAction        ReportNextAction          `json:"nextAction"`
	TargetSummary     []TargetRatingSummary     `json:"targetSummary"`
	Targets           []TargetEvaluationDetail  `json:"targets"`
	AssessmentResults []AssessmentResultDigest  `json:"assessmentResults"`
	FindingSummaries  []FindingDigest           `json:"findingSummaries"`
	Recommendations   []RecommendationReference `json:"recommendations"`
}

type ReportJSON = EvaluationReportDocument

type ReportSummary struct {
	Run          string       `json:"run,omitempty"`
	Subject      string       `json:"subject"`
	Altitude     string       `json:"altitude"`
	Rigor        string       `json:"rigor"`
	Narrowing    string       `json:"narrowing"`
	RatingResult RatingResult `json:"ratingResult"`
}

type ReportScope struct {
	Recorded    bool     `json:"recorded"`
	Description string   `json:"description"`
	Narrowing   string   `json:"narrowing"`
	InScope     []string `json:"inScope"`
	OutOfScope  []string `json:"outOfScope"`
	NotRecorded []string `json:"notRecorded"`
}

type ReportEvidence struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

type ReportNextAction struct {
	Kind               string `json:"kind"`
	Summary            string `json:"summary"`
	RecommendationID   string `json:"recommendationId,omitempty"`
	RecommendationPath string `json:"recommendationPath,omitempty"`
}

type TargetRatingSummary struct {
	TargetPath            []string      `json:"targetPath"`
	LocalRatingResult     *RatingResult `json:"localRatingResult"`
	AggregateRatingResult RatingResult  `json:"aggregateRatingResult"`
	CoveredRequirements   int           `json:"coveredRequirements"`
	Note                  string        `json:"note,omitempty"`
}

type TargetEvaluationDetail struct {
	TargetPath              []string             `json:"targetPath"`
	LocalRatingResult       *RatingResult        `json:"localRatingResult"`
	AggregateRatingResult   RatingResult         `json:"aggregateRatingResult"`
	FactorRatingResults     []FactorRatingResult `json:"factorRatingResults"`
	AnalysisRecord          string               `json:"analysisRecord"`
	NotAssessedRequirements []string             `json:"notAssessedRequirements"`
	Structural              bool                 `json:"structural"`
}

type AssessmentResultDigest struct {
	AssessmentResultRecord string       `json:"assessmentResultRecord"`
	TargetPath             []string     `json:"targetPath"`
	Requirement            string       `json:"requirement"`
	RatingResult           RatingResult `json:"ratingResult"`
	Active                 bool         `json:"active"`
	Supersedes             []string     `json:"supersedes,omitempty"`
}

type FindingDigest struct {
	AssessmentResultRecord string `json:"assessmentResultRecord"`
	Locator                string `json:"locator"`
	Category               string `json:"category"`
	Severity               string `json:"severity,omitempty"`
	Summary                string `json:"summary"`
}

type RecommendationReference struct {
	ID            string   `json:"id"`
	Path          string   `json:"path"`
	DoneCriterion string   `json:"doneCriterion"`
	Active        bool     `json:"active"`
	Supersedes    []string `json:"supersedes,omitempty"`
}

type BuildReportReceipt struct {
	SchemaVersion   int              `json:"schemaVersion"`
	Path            string           `json:"path"`
	ReportSummaryMD string           `json:"reportSummaryMd"`
	ReportMD        string           `json:"reportMd"`
	ReportJSON      string           `json:"reportJson"`
	RatingResult    RatingResult     `json:"ratingResult"`
	NextActions     []receipt.Action `json:"nextActions,omitempty"`
}

type GateReceipt struct {
	SchemaVersion int          `json:"schemaVersion"`
	Path          string       `json:"path"`
	Threshold     string       `json:"threshold"`
	RatingResult  RatingResult `json:"ratingResult"`
	Pass          bool         `json:"pass"`
}

func BuildReport(path string) (*BuildReportReceipt, error) {
	run, err := Load(path)
	if err != nil {
		return nil, err
	}
	if gaps := run.Renderable(); len(gaps) > 0 {
		return nil, fmt.Errorf("run is not reportable: %s %s", gaps[0].Kind, gaps[0].Ref)
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

func (r *EvaluationRun) Report() (EvaluationReportDocument, error) {
	analyses, rootAnalysis, err := sortedReportAnalyses(r.Analyses)
	if err != nil {
		return EvaluationReportDocument{}, err
	}
	context := r.reportContext()
	report := newReportDocument(r.Path, rootAnalysis, context)
	collector := limitationCollector{seen: map[string]bool{}, out: &report.Limitations}
	collector.addContext(context, rootAnalysis)

	supersededAssessmentResults, _ := r.assessmentResultSupersedingState()
	assessmentResultsByFile := addAssessmentResultSummaries(&report, r.AssessmentResults, supersededAssessmentResults, &collector)
	addReportTargets(&report, analyses, assessmentResultsByFile, context, &collector)
	addReportRecommendations(&report, r.Recommendations, r.supersededRecommendations())
	setReportNextAction(&report)
	return report, nil
}

func sortedReportAnalyses(records []AnalysisRecord) ([]AnalysisRecord, AnalysisRecord, error) {
	analyses := append([]AnalysisRecord(nil), records...)
	if len(analyses) == 0 {
		return nil, AnalysisRecord{}, fmt.Errorf("run has no analysis records")
	}
	slices.SortFunc(analyses, func(a, b AnalysisRecord) int {
		if cmp := len(a.TargetPath) - len(b.TargetPath); cmp != 0 {
			return cmp
		}
		return strings.Compare(strings.Join(a.TargetPath, "/"), strings.Join(b.TargetPath, "/"))
	})
	var root *AnalysisRecord
	for _, analysis := range analyses {
		if len(analysis.TargetPath) == 0 {
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

func newReportDocument(runPath string, rootAnalysis AnalysisRecord, context runContext) EvaluationReportDocument {
	rootRating := rootAnalysis.AggregateRatingResult
	report := EvaluationReportDocument{
		SchemaVersion: SchemaVersion,
		Summary: ReportSummary{
			Run:          filepath.Base(runPath),
			Subject:      targetPathDisplay(rootAnalysis.TargetPath),
			Altitude:     context.Altitude,
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
			NotRecorded: []string{},
		},
		EvidenceBasis:     []ReportEvidence{},
		Limitations:       []string{},
		TargetSummary:     []TargetRatingSummary{},
		Targets:           []TargetEvaluationDetail{},
		AssessmentResults: []AssessmentResultDigest{},
		FindingSummaries:  []FindingDigest{},
		Recommendations:   []RecommendationReference{},
	}
	normalizeReportScope(&report)
	return report
}

func normalizeReportScope(report *EvaluationReportDocument) {
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
		report.Scope.NotRecorded = append(report.Scope.NotRecorded, "rigor")
	}
	if len(report.Scope.OutOfScope) == 0 {
		report.Scope.NotRecorded = append(report.Scope.NotRecorded, "out-of-scope areas")
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

func addAssessmentResultSummaries(report *EvaluationReportDocument, records []AssessmentResultRecord, superseded map[string]bool, collector *limitationCollector) map[string]AssessmentResultRecord {
	recordsByFile := map[string]AssessmentResultRecord{}
	evidenceSeen := map[string]bool{}
	for _, record := range records {
		active := !superseded[record.File]
		recordsByFile[record.File] = record
		report.AssessmentResults = append(report.AssessmentResults, AssessmentResultDigest{
			AssessmentResultRecord: record.File,
			TargetPath:             record.TargetPath,
			Requirement:            record.Requirement,
			RatingResult:           record.RatingResult,
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

func addReportFindings(report *EvaluationReportDocument, record AssessmentResultRecord, evidenceSeen map[string]bool) {
	for _, finding := range record.Findings {
		report.FindingSummaries = append(report.FindingSummaries, FindingDigest{
			AssessmentResultRecord: record.File,
			Locator:                finding.Locator,
			Category:               finding.Category,
			Severity:               finding.Severity,
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

func addReportTargets(report *EvaluationReportDocument, analyses []AnalysisRecord, assessmentResultsByFile map[string]AssessmentResultRecord, context runContext, collector *limitationCollector) {
	for _, analysis := range analyses {
		structural := len(analysis.AssessmentResultRecords) == 0 && len(analysis.ChildAnalysisRecords) > 0
		if analysis.LocalRatingResult != nil {
			collector.addText(analysis.LocalRatingResult.Rationale)
		}
		collector.addText(analysis.AggregateRatingResult.Rationale)
		for _, factor := range analysis.FactorRatingResults {
			collector.addText(factor.RatingResult.Rationale)
		}
		target := targetEvaluationFromAnalysis(analysis, structural)
		for _, ref := range analysis.AssessmentResultRecords {
			if record, ok := assessmentResultsByFile[ref]; ok && isNotAssessed(record.RatingResult) {
				target.NotAssessedRequirements = append(target.NotAssessedRequirements, record.Requirement)
			}
		}
		report.TargetSummary = append(report.TargetSummary, targetEvaluationSummary(target, len(analysis.AssessmentResultRecords)))
		if len(context.InScope) == 0 {
			report.Scope.InScope = append(report.Scope.InScope, targetPathDisplay(target.TargetPath))
		}
		report.Targets = append(report.Targets, target)
	}
}

func targetEvaluationFromAnalysis(analysis AnalysisRecord, structural bool) TargetEvaluationDetail {
	target := TargetEvaluationDetail{
		TargetPath:              append([]string(nil), analysis.TargetPath...),
		AnalysisRecord:          analysis.File,
		LocalRatingResult:       cloneRatingResult(analysis.LocalRatingResult),
		AggregateRatingResult:   analysis.AggregateRatingResult,
		FactorRatingResults:     nonNilFactorRatingResults(analysis.FactorRatingResults),
		NotAssessedRequirements: []string{},
		Structural:              structural,
	}
	if structural {
		target.LocalRatingResult = nil
	}
	return target
}

func targetEvaluationSummary(target TargetEvaluationDetail, coveredRequirements int) TargetRatingSummary {
	note := ""
	if target.Structural {
		note = "structural grouping target"
	} else if len(target.NotAssessedRequirements) > 0 {
		note = "has not-assessed requirements"
	}
	return TargetRatingSummary{
		TargetPath:            append([]string(nil), target.TargetPath...),
		LocalRatingResult:     cloneRatingResult(target.LocalRatingResult),
		AggregateRatingResult: target.AggregateRatingResult,
		CoveredRequirements:   coveredRequirements,
		Note:                  note,
	}
}

func addReportRecommendations(report *EvaluationReportDocument, recommendations []RecommendationRecord, superseded map[string]bool) {
	for _, rec := range recommendations {
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		report.Recommendations = append(report.Recommendations, RecommendationReference{
			ID:            id,
			Path:          rec.File,
			DoneCriterion: rec.DoneCriterion,
			Active:        !superseded[rec.File] && !superseded[id],
			Supersedes:    append([]string(nil), rec.Supersedes...),
		})
	}
}

func setReportNextAction(report *EvaluationReportDocument) {
	if first, ok := firstActiveRecommendation(report.Recommendations); ok {
		report.NextAction = ReportNextAction{
			Kind:               "recommendation",
			Summary:            first.DoneCriterion,
			RecommendationID:   first.ID,
			RecommendationPath: first.Path,
		}
	} else {
		report.NextAction = ReportNextAction{
			Kind:    "none",
			Summary: "No recommendation records exist for this run.",
		}
	}
}

func renderReportMarkdown(report EvaluationReportDocument, labels reportDisplayLabels) []byte {
	var out bytes.Buffer
	out.WriteString("# Evaluation Report\n\n")
	writeReportSummarySection(&out, report, labels)
	writeReportScopeSection(&out, report, labels)
	writeReportRisksAndLimitationsSection(&out, report)
	writeReportEvidenceSection(&out, report)
	writeReportNextActionSection(&out, report)
	writeReportTargetSummarySection(&out, report, labels)
	writeReportTargetDetailsSection(&out, report, labels)
	writeReportRequirementsSection(&out, report, labels)
	writeReportFindingsSection(&out, report)
	writeReportAdviceSection(&out, report)
	return out.Bytes()
}

func writeReportSummarySection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Summary\n\n")
	out.WriteString("- **Subject:** " + labels.Target(nil, report.Summary.Subject) + "\n")
	if report.Summary.Altitude != "" {
		out.WriteString("- **Altitude:** " + report.Summary.Altitude + "\n")
	} else {
		out.WriteString("- **Altitude:** not recorded\n")
	}
	if report.Summary.Rigor != "" {
		out.WriteString("- **Rigor:** " + report.Summary.Rigor + "\n")
	} else {
		out.WriteString("- **Rigor:** not recorded\n")
	}
	if report.Summary.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Summary.Narrowing + "\n")
	}
	out.WriteString("- **Rating:** " + displayRatingResult(report.Summary.RatingResult, labels.Ratings) + "\n")
	if report.Summary.RatingResult.Rationale != "" {
		out.WriteString("- **Rationale:** " + report.Summary.RatingResult.Rationale + "\n")
	}
}

func writeReportScopeSection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Scope\n\n")
	out.WriteString(report.Scope.Description + "\n\n")
	if report.Scope.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Scope.Narrowing + "\n")
	} else {
		out.WriteString("- **Narrowing:** whole recorded run\n")
	}
	out.WriteString("- **In scope:** " + displayScopeList(report.Scope.InScope, report.TargetSummary, labels) + "\n")
	if len(report.Scope.OutOfScope) > 0 {
		out.WriteString("- **Out of scope:** " + displayList(report.Scope.OutOfScope) + "\n")
	} else {
		out.WriteString("- **Out of scope:** not recorded\n")
	}
	if len(report.Scope.NotRecorded) > 0 {
		out.WriteString("- **Metadata not recorded:** " + displayList(report.Scope.NotRecorded) + "\n")
	}
}

func writeReportRisksAndLimitationsSection(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Top Risks and Limitations\n\n")
	risks := riskFindings(report.FindingSummaries)
	summaryRisks := firstFindingSummaries(risks, 8)
	summaryLimitations := firstStrings(report.Limitations, 8)
	if len(summaryRisks) == 0 && len(summaryLimitations) == 0 {
		out.WriteString("No top risks or limitations were recorded in the summary data.\n")
		return
	}
	for _, finding := range summaryRisks {
		out.WriteString("- `" + finding.AssessmentResultRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		if finding.Severity != "" {
			out.WriteString(" [" + finding.Severity + "]")
		}
		out.WriteString(": " + finding.Summary + "\n")
	}
	for _, limitation := range summaryLimitations {
		out.WriteString("- Limitation: " + limitation + "\n")
	}
	if len(risks) > len(summaryRisks) || len(report.Limitations) > len(summaryLimitations) {
		out.WriteString("- Additional risks or limitations are available in `report.json`.\n")
	}
}

func writeReportEvidenceSection(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Evidence Basis\n\n")
	if len(report.EvidenceBasis) == 0 {
		out.WriteString("No command or source evidence basis was recorded in findings.\n")
	} else {
		for _, evidence := range report.EvidenceBasis {
			out.WriteString("- **" + evidence.Kind + ":** `" + evidence.Ref + "`\n")
		}
	}
}

func writeReportNextActionSection(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Next Action\n\n")
	if report.NextAction.RecommendationID != "" {
		out.WriteString("- [" + report.NextAction.RecommendationID + "](" + report.NextAction.RecommendationPath + ") - " + report.NextAction.Summary + "\n")
	} else {
		out.WriteString(report.NextAction.Summary + "\n")
	}
}

func writeReportTargetSummarySection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Target Summary\n\n")
	out.WriteString("| Target | Local | Aggregate | Covered Requirements | Note |\n")
	out.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, target := range report.TargetSummary {
		out.WriteString("| " + tableCell(labels.Target(target.TargetPath, targetPathDisplay(target.TargetPath))) + " | " + tableCell(displayOptionalRatingResult(target.LocalRatingResult, labels.Ratings)) + " | " + tableCell(displayRatingResult(target.AggregateRatingResult, labels.Ratings)) + " | " + fmt.Sprintf("%d", target.CoveredRequirements) + " | " + tableCell(target.Note) + " |\n")
	}
}

func writeReportTargetDetailsSection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Target Details\n\n")
	for _, target := range report.Targets {
		writeReportTargetDetail(out, target, labels)
	}
}

func writeReportTargetDetail(out *bytes.Buffer, target TargetEvaluationDetail, labels reportDisplayLabels) {
	out.WriteString("### " + labels.Target(target.TargetPath, targetPathDisplay(target.TargetPath)) + "\n\n")
	out.WriteString("- **Path:** " + displayPath(target.TargetPath) + "\n")
	out.WriteString("- **Local rating:** " + displayOptionalRatingResult(target.LocalRatingResult, labels.Ratings) + "\n")
	if target.LocalRatingResult != nil {
		writeOptionalRationale(out, target.LocalRatingResult.Rationale)
	}
	out.WriteString("- **Aggregate rating:** " + displayRatingResult(target.AggregateRatingResult, labels.Ratings) + "\n")
	writeOptionalRationale(out, target.AggregateRatingResult.Rationale)
	for _, factor := range target.FactorRatingResults {
		out.WriteString("- **Factor " + labels.Factor(target.TargetPath, factor.FactorPath) + ":** " + displayRatingResult(factor.RatingResult, labels.Ratings) + "\n")
		writeOptionalRationale(out, factor.RatingResult.Rationale)
	}
	out.WriteString("- **Analysis record:** `" + target.AnalysisRecord + "`\n")
	if len(target.NotAssessedRequirements) > 0 {
		out.WriteString("- **Not assessed:** " + strings.Join(target.NotAssessedRequirements, "; ") + "\n")
	}
	out.WriteString("\n")
}

func writeOptionalRationale(out *bytes.Buffer, rationale string) {
	if rationale != "" {
		out.WriteString("  - " + rationale + "\n")
	}
}

func writeReportRequirementsSection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Requirements\n\n")
	for _, result := range report.AssessmentResults {
		out.WriteString("### " + result.Requirement + "\n\n")
		state := "active"
		if !result.Active {
			state = "superseded"
		}
		out.WriteString("- **State:** " + state + "\n")
		out.WriteString("- **Target:** " + labels.Target(result.TargetPath, targetPathDisplay(result.TargetPath)) + "\n")
		out.WriteString("- **Rating:** " + displayRatingResult(result.RatingResult, labels.Ratings) + "\n")
		out.WriteString("- **Assessment result record:** `" + result.AssessmentResultRecord + "`\n")
		if result.RatingResult.Rationale != "" {
			out.WriteString("- **Rationale:** " + result.RatingResult.Rationale + "\n")
		}
		out.WriteString("\n")
	}
}

func writeReportFindingsSection(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("## Findings\n\n")
	for _, finding := range report.FindingSummaries {
		out.WriteString("- `" + finding.AssessmentResultRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		out.WriteString(": " + finding.Summary + "\n")
	}
}

func writeReportAdviceSection(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Advice\n\n")
	if len(report.Recommendations) == 0 {
		out.WriteString("No recommendation records exist for this run.\n")
	} else {
		for _, rec := range report.Recommendations {
			state := "active"
			if !rec.Active {
				state = "superseded"
			}
			out.WriteString("- [" + rec.ID + "](" + rec.Path + ") [" + state + "] - " + rec.DoneCriterion + "\n")
		}
	}
}

func renderReportSummaryMarkdown(report EvaluationReportDocument, labels reportDisplayLabels) []byte {
	var out bytes.Buffer
	out.WriteString("# Quality Evaluation Summary\n\n")
	writeSummaryKeyDetails(&out, report, labels)
	writeSummarySection(&out, report, labels)
	writeSummaryTopIssues(&out, report)
	writeSummaryRecommendations(&out, report)
	writeSummaryScopeAndLimitations(&out, report, labels)
	return out.Bytes()
}

func writeSummaryKeyDetails(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("| Field | Value |\n")
	out.WriteString("| --- | --- |\n")
	out.WriteString("| Subject | " + tableCell(labels.Target(nil, report.Summary.Subject)) + " |\n")
	if report.Summary.Run != "" {
		out.WriteString("| Run | `" + tableCell(report.Summary.Run) + "` |\n")
	}
	out.WriteString("| Scope | " + tableCell(summaryScope(report)) + " |\n")
	out.WriteString("| Rigor | " + tableCell(summaryValueTitle(report.Summary.Rigor)) + " |\n")
	out.WriteString("| Overall rating | " + tableCell(displayRatingResult(report.Summary.RatingResult, labels.Ratings)) + " |\n")
	out.WriteString("| Full report | [report.md](report.md) |\n")
	out.WriteString("| Machine report | [report.json](report.json) |\n\n")
}

func writeSummarySection(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("## Summary\n\n")
	if report.Summary.RatingResult.Rationale != "" {
		out.WriteString(report.Summary.RatingResult.Rationale + "\n")
	} else if isNotAssessed(report.Summary.RatingResult) {
		out.WriteString("The evaluation did not produce an overall rating.\n")
	} else {
		out.WriteString("The evaluation completed with overall rating " + displayRatingResult(report.Summary.RatingResult, labels.Ratings) + ".\n")
	}
	writeSummaryRatingTable(out, report, labels)
}

func writeSummaryRatingTable(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n\n")
	if len(report.TargetSummary) == 0 {
		out.WriteString("No target ratings were recorded.\n")
		return
	}
	out.WriteString("| Target | Local rating | Overall rating | Driver |\n")
	out.WriteString("| --- | --- | --- | --- |\n")
	for _, target := range report.TargetSummary {
		out.WriteString("| " + tableCell(labels.Target(target.TargetPath, targetPathDisplay(target.TargetPath))) + " | " + tableCell(displayOptionalRatingResult(target.LocalRatingResult, labels.Ratings)) + " | " + tableCell(displayRatingResult(target.AggregateRatingResult, labels.Ratings)) + " | " + tableCell(summaryDriver(target)) + " |\n")
	}
}

func summaryDriver(target TargetRatingSummary) string {
	if target.AggregateRatingResult.Rationale != "" {
		return target.AggregateRatingResult.Rationale
	}
	if target.Note != "" {
		return target.Note
	}
	return ""
}

func writeSummaryTopIssues(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Top Issues\n\n")
	issues := firstFindingSummaries(riskFindings(report.FindingSummaries), 5)
	if len(issues) == 0 {
		out.WriteString("None recorded.\n")
		return
	}
	for i, finding := range issues {
		fmt.Fprintf(out, "%d. ", i+1)
		if finding.Severity != "" {
			out.WriteString("**" + finding.Severity + "**  \n   ")
		} else if finding.Category != "" {
			out.WriteString("**" + finding.Category + "**  \n   ")
		}
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

func writeSummaryRecommendations(out *bytes.Buffer, report EvaluationReportDocument) {
	out.WriteString("\n## Recommendations\n\n")
	active := activeRecommendations(report.Recommendations)
	if report.NextAction.RecommendationID != "" {
		out.WriteString("Primary next action: use `" + report.NextAction.RecommendationID + "`.\n\n")
	} else {
		out.WriteString(report.NextAction.Summary + "\n")
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

func writeSummaryScopeAndLimitations(out *bytes.Buffer, report EvaluationReportDocument, labels reportDisplayLabels) {
	out.WriteString("\n## Scope & Limitations\n\n")
	out.WriteString("Scope: **" + summaryScope(report) + "**\n")
	inScope := displayScopeList(report.Scope.InScope, report.TargetSummary, labels)
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

func summaryScope(report EvaluationReportDocument) string {
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
	case "", "whole model", "whole subject", "whole recorded run", "full evaluation", "evaluation scope reconstructed from the run's analysis and assessment result records.":
		return true
	default:
		return false
	}
}

func summaryValue(value string) string {
	if value == "" {
		return "not recorded"
	}
	return value
}

func summaryValueTitle(value string) string {
	value = summaryValue(value)
	if value == "not recorded" {
		return value
	}
	return strings.ToUpper(value[:1]) + value[1:]
}

func activeRecommendations(recommendations []RecommendationReference) []RecommendationReference {
	active := []RecommendationReference{}
	for _, rec := range recommendations {
		if rec.Active {
			active = append(active, rec)
		}
	}
	return active
}

func (r *EvaluationRun) supersededRecommendations() map[string]bool {
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
		if rec.Active {
			return rec, true
		}
	}
	return RecommendationReference{}, false
}

type runContext struct {
	Recorded         bool
	Altitude         string
	Narrowing        string
	Rigor            string
	ScopeDescription string
	InScope          []string
	OutOfScope       []string
	Limitations      []string
}

func (r *EvaluationRun) reportContext() runContext {
	context := inferRunContext(r.Path)
	designParams := markdownBulletMap(markdownSection(r.Design, "Resolved parameters"))
	if value := designParams["Altitude"]; value != "" {
		context.Altitude = cleanInlineCode(value)
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
		context.Rigor = strings.ToLower(cleanInlineCode(value))
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
	parts := strings.SplitN(name, "-", 3)
	if len(parts) < 2 {
		return runContext{}
	}
	context := runContext{Altitude: parts[1]}
	if len(parts) == 3 {
		context.Narrowing = parts[2]
		context.Rigor = inferRigor(parts[2])
	}
	return context
}

func inferRigor(narrowing string) string {
	for _, rigor := range []string{"quick", "standard", "deep"} {
		if strings.Contains(narrowing, rigor) {
			return rigor
		}
	}
	return ""
}

func markdownSection(doc, heading string) string {
	lines := strings.Split(doc, "\n")
	target := "## " + heading
	var out []string
	inSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.EqualFold(trimmed, target) {
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

func riskFindings(findings []FindingDigest) []FindingDigest {
	risks := []FindingDigest{}
	for _, finding := range findings {
		if finding.Severity == "" {
			continue
		}
		switch strings.ToLower(finding.Severity) {
		case "info", "informational", "note":
			continue
		default:
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
	Targets map[string]string
	Factors map[string]string
}

func reportDisplayLabelsFromModel(spec *model.Spec) reportDisplayLabels {
	labels := reportDisplayLabels{
		Ratings: map[string]string{},
		Targets: map[string]string{},
		Factors: map[string]string{},
	}
	if spec == nil {
		return labels
	}
	labels.Ratings = ratingDisplayLabels(spec.RatingScale)
	if strings.TrimSpace(spec.Title) != "" {
		labels.Targets[targetLabelKey(nil)] = spec.Title
	}
	labels.addFactors(nil, spec.Factors)
	labels.addTargets(nil, spec.Targets)
	return labels
}

func (l reportDisplayLabels) addTargets(parentPath []string, targets map[string]model.Target) {
	for key, target := range targets {
		path := appendPathElement(parentPath, key)
		if strings.TrimSpace(target.Title) != "" {
			l.Targets[targetLabelKey(path)] = target.Title
		}
		l.addFactors(path, target.Factors)
		l.addTargets(path, target.Targets)
	}
}

func (l reportDisplayLabels) addFactors(targetPath []string, factors map[string]model.Factor) {
	for key, factor := range factors {
		path := appendPathElement(nil, key)
		if strings.TrimSpace(factor.Title) != "" {
			l.Factors[factorLabelKey(targetPath, path)] = factor.Title
		}
		l.addSubFactors(targetPath, path, factor.Factors)
	}
}

func (l reportDisplayLabels) addSubFactors(targetPath, parentFactorPath []string, factors map[string]model.Factor) {
	for key, factor := range factors {
		path := appendPathElement(parentFactorPath, key)
		if strings.TrimSpace(factor.Title) != "" {
			l.Factors[factorLabelKey(targetPath, path)] = factor.Title
		}
		l.addSubFactors(targetPath, path, factor.Factors)
	}
}

func (l reportDisplayLabels) Target(path []string, fallback string) string {
	if title := l.Targets[targetLabelKey(path)]; title != "" {
		return title
	}
	return fallback
}

func (l reportDisplayLabels) Factor(targetPath []string, factorPath []string) string {
	for i := len(targetPath); i >= 0; i-- {
		if title := l.Factors[factorLabelKey(targetPath[:i], factorPath)]; title != "" {
			return title
		}
	}
	return factorPathDisplay(factorPath)
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

func displayOptionalRatingResult(result *RatingResult, ratingLabels map[string]string) string {
	if result == nil {
		return "n/a (structural)"
	}
	return displayRatingResult(*result, ratingLabels)
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
	return result.Kind == "not-assessed"
}

func displayList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(items, "; ")
}

func displayScopeList(items []string, targets []TargetRatingSummary, labels reportDisplayLabels) string {
	if len(items) == 0 {
		return "none"
	}
	targetLabels := map[string]string{}
	for _, target := range targets {
		id := targetPathDisplay(target.TargetPath)
		if _, exists := targetLabels[id]; exists {
			targetLabels[id] = ""
			continue
		}
		targetLabels[id] = labels.Target(target.TargetPath, id)
	}
	display := make([]string, 0, len(items))
	for _, item := range items {
		if label := targetLabels[item]; label != "" {
			display = append(display, label)
			continue
		}
		display = append(display, item)
	}
	return strings.Join(display, "; ")
}

func displayPath(path []string) string {
	if len(path) == 0 {
		return "(root)"
	}
	return strings.Join(path, " / ")
}

func targetPathDisplay(path []string) string {
	if len(path) == 0 {
		return "root"
	}
	return strings.Join(path, "/")
}

func factorPathDisplay(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return strings.Join(path, "/")
}

func targetLabelKey(path []string) string {
	return strings.Join(path, "\x00")
}

func factorLabelKey(targetPath []string, factorPath []string) string {
	return targetLabelKey(targetPath) + "\x00" + strings.Join(factorPath, "\x00")
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
		return false, fmt.Errorf("overall rating %q is not in the run rating scale", result.RatingResult.Level)
	}
	return ratingIndex < thresholdIndex, nil
}

func GateReport(path, threshold string) (*GateReceipt, error) {
	run, err := Load(path)
	if err != nil {
		return nil, err
	}
	reportPath := filepath.Join(run.Path, "report.json")
	raw, err := os.ReadFile(reportPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("report.json has not been built for %s", run.Path)
		}
		return nil, fmt.Errorf("reading report.json: %w", err)
	}
	var report EvaluationReportDocument
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

func ScaleLevels(path string) ([]string, error) {
	run, err := Load(path)
	if err != nil {
		return nil, err
	}
	levels := make([]string, 0, len(run.Scale))
	for _, level := range run.Scale {
		levels = append(levels, level.Level)
	}
	return levels, nil
}
