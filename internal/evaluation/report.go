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

type ReportJSON struct {
	SchemaVersion    int                    `json:"schemaVersion"`
	Summary          ReportSummary          `json:"summary"`
	Rating           ReportRating           `json:"rating"`
	Scope            ReportScope            `json:"scope"`
	EvidenceBasis    []ReportEvidence       `json:"evidenceBasis"`
	Limitations      []string               `json:"limitations"`
	NextAction       ReportNextAction       `json:"nextAction"`
	TargetSummary    []ReportTargetSummary  `json:"targetSummary"`
	Targets          []ReportTarget         `json:"targets"`
	Assessments      []ReportAssessment     `json:"assessments"`
	FindingSummaries []FindingSummary       `json:"findingSummaries"`
	Recommendations  []ReportRecommendation `json:"recommendations"`
}

type ReportSummary struct {
	Run         string  `json:"run,omitempty"`
	Subject     string  `json:"subject"`
	Altitude    string  `json:"altitude"`
	Effort      string  `json:"effort"`
	Narrowing   string  `json:"narrowing"`
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Rationale   string  `json:"rationale"`
}

type ReportRating struct {
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Subject     string  `json:"subject"`
	Rationale   string  `json:"rationale"`
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

type ReportTargetSummary struct {
	Target              string             `json:"target"`
	TargetPath          []string           `json:"targetPath"`
	LocalRating         ReportRatingResult `json:"localRating"`
	AggregateRating     ReportRatingResult `json:"aggregateRating"`
	CoveredRequirements int                `json:"coveredRequirements"`
	Note                string             `json:"note,omitempty"`
}

type ReportTarget struct {
	Target          string             `json:"target"`
	TargetPath      []string           `json:"targetPath"`
	LocalRating     ReportRatingResult `json:"localRating"`
	AggregateRating ReportRatingResult `json:"aggregateRating"`
	FactorRatings   []FactorRating     `json:"factorRatings"`
	AnalysisRecord  string             `json:"analysisRecord"`
	NotAssessed     []string           `json:"notAssessed"`
	Structural      bool               `json:"structural"`
}

type ReportRatingResult struct {
	Kind        string  `json:"kind"`
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Rationale   string  `json:"rationale"`
}

type ReportAssessment struct {
	AssessmentRecord string   `json:"assessmentRecord"`
	Target           string   `json:"target"`
	TargetPath       []string `json:"targetPath"`
	Requirement      string   `json:"requirement"`
	Rating           *string  `json:"rating"`
	NotAssessed      bool     `json:"notAssessed"`
	Rationale        string   `json:"rationale"`
	Active           bool     `json:"active"`
	Supersedes       []string `json:"supersedes,omitempty"`
}

type FindingSummary struct {
	AssessmentRecord string `json:"assessmentRecord"`
	Locator          string `json:"locator"`
	Category         string `json:"category"`
	Severity         string `json:"severity,omitempty"`
	Summary          string `json:"summary"`
}

type ReportRecommendation struct {
	ID            string   `json:"id"`
	Path          string   `json:"path"`
	DoneCriterion string   `json:"doneCriterion"`
	Active        bool     `json:"active"`
	Supersedes    []string `json:"supersedes,omitempty"`
}

type BuildResult struct {
	SchemaVersion   int              `json:"schemaVersion"`
	Path            string           `json:"path"`
	ReportSummaryMD string           `json:"reportSummaryMd"`
	ReportMD        string           `json:"reportMd"`
	ReportJSON      string           `json:"reportJson"`
	Rating          *string          `json:"rating"`
	NotAssessed     bool             `json:"notAssessed"`
	NextActions     []receipt.Action `json:"nextActions,omitempty"`
}

type GateResult struct {
	SchemaVersion int     `json:"schemaVersion"`
	Path          string  `json:"path"`
	Threshold     string  `json:"threshold"`
	Rating        *string `json:"rating"`
	NotAssessed   bool    `json:"notAssessed"`
	Pass          bool    `json:"pass"`
}

func BuildReport(path string) (*BuildResult, error) {
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
	reportSummaryMD := filepath.Join(run.Path, "report-summary.md")
	reportMD := filepath.Join(run.Path, "report.md")
	reportJSON := filepath.Join(run.Path, "report.json")
	if err := writeReportFile(reportSummaryMD, summaryMD); err != nil {
		return nil, err
	}
	if err := writeReportFile(reportMD, md); err != nil {
		return nil, err
	}
	if err := writeReportFile(reportJSON, js); err != nil {
		return nil, err
	}
	return &BuildResult{
		SchemaVersion:   SchemaVersion,
		Path:            run.Path,
		ReportSummaryMD: filepath.ToSlash(reportSummaryMD),
		ReportMD:        filepath.ToSlash(reportMD),
		ReportJSON:      filepath.ToSlash(reportJSON),
		Rating:          report.Rating.Rating,
		NotAssessed:     report.Rating.NotAssessed,
	}, nil
}

func (r *Run) Report() (ReportJSON, error) {
	analyses, rootAnalysis, err := sortedReportAnalyses(r.Analyses)
	if err != nil {
		return ReportJSON{}, err
	}
	context := r.reportContext()
	report := newReportJSON(r.Path, rootAnalysis, context)
	collector := limitationCollector{seen: map[string]bool{}, out: &report.Limitations}
	collector.addContext(context, rootAnalysis)

	supersededAssessments, _ := r.assessmentSupersedingState()
	assessmentsByFile := addReportAssessments(&report, r.Assessments, supersededAssessments, &collector)
	addReportTargets(&report, analyses, assessmentsByFile, context, &collector)
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

func newReportJSON(runPath string, rootAnalysis AnalysisRecord, context runContext) ReportJSON {
	report := ReportJSON{
		SchemaVersion: SchemaVersion,
		Summary: ReportSummary{
			Run:         filepath.Base(runPath),
			Subject:     rootAnalysis.Target,
			Altitude:    context.Altitude,
			Effort:      context.Effort,
			Narrowing:   context.Narrowing,
			Rating:      rootAnalysis.AggregateRating.Rating,
			NotAssessed: rootAnalysis.AggregateRating.NotAssessed,
			Rationale:   rootAnalysis.AggregateRating.Rationale,
		},
		Rating: ReportRating{
			Rating:      rootAnalysis.AggregateRating.Rating,
			NotAssessed: rootAnalysis.AggregateRating.NotAssessed,
			Subject:     rootAnalysis.Target,
			Rationale:   rootAnalysis.AggregateRating.Rationale,
		},
		Scope: ReportScope{
			Recorded:    context.Recorded,
			Description: context.ScopeDescription,
			Narrowing:   context.Narrowing,
			InScope:     append([]string(nil), context.InScope...),
			OutOfScope:  append([]string(nil), context.OutOfScope...),
			NotRecorded: []string{},
		},
		EvidenceBasis:    []ReportEvidence{},
		Limitations:      []string{},
		TargetSummary:    []ReportTargetSummary{},
		Targets:          []ReportTarget{},
		Assessments:      []ReportAssessment{},
		FindingSummaries: []FindingSummary{},
		Recommendations:  []ReportRecommendation{},
	}
	normalizeReportScope(&report)
	return report
}

func normalizeReportScope(report *ReportJSON) {
	if report.Scope.Description == "" {
		report.Scope.Description = "Evaluation scope reconstructed from the run's analysis and assessment records."
	}
	if report.Scope.InScope == nil {
		report.Scope.InScope = []string{}
	}
	if report.Scope.OutOfScope == nil {
		report.Scope.OutOfScope = []string{}
	}
	if report.Summary.Effort == "" {
		report.Scope.NotRecorded = append(report.Scope.NotRecorded, "effort")
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
	c.addText(rootAnalysis.AggregateRating.Rationale)
}

func addReportAssessments(report *ReportJSON, assessments []AssessmentRecord, superseded map[string]bool, collector *limitationCollector) map[string]AssessmentRecord {
	assessmentsByFile := map[string]AssessmentRecord{}
	evidenceSeen := map[string]bool{}
	for _, assessment := range assessments {
		active := !superseded[assessment.File]
		assessmentsByFile[assessment.File] = assessment
		report.Assessments = append(report.Assessments, ReportAssessment{
			AssessmentRecord: assessment.File,
			Target:           assessment.Target,
			TargetPath:       assessment.TargetPath,
			Requirement:      assessment.Requirement,
			Rating:           assessment.Rating,
			NotAssessed:      assessment.NotAssessed,
			Rationale:        assessment.Rationale,
			Active:           active,
			Supersedes:       append([]string(nil), assessment.Supersedes...),
		})
		if !active {
			continue
		}
		if assessment.NotAssessed {
			collector.addText(assessment.Requirement + ": " + assessment.Rationale)
		}
		collector.addText(assessment.Rationale)
		addReportFindings(report, assessment, evidenceSeen)
	}
	return assessmentsByFile
}

func addReportFindings(report *ReportJSON, assessment AssessmentRecord, evidenceSeen map[string]bool) {
	for _, finding := range assessment.Findings {
		report.FindingSummaries = append(report.FindingSummaries, FindingSummary{
			AssessmentRecord: assessment.File,
			Locator:          finding.Locator,
			Category:         finding.Category,
			Severity:         finding.Severity,
			Summary:          finding.Observation,
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

func addReportTargets(report *ReportJSON, analyses []AnalysisRecord, assessmentsByFile map[string]AssessmentRecord, context runContext, collector *limitationCollector) {
	for _, analysis := range analyses {
		structural := len(analysis.AssessmentRecords) == 0 && len(analysis.ChildAnalysisRecords) > 0
		if analysis.LocalRating != nil {
			collector.addText(analysis.LocalRating.Rationale)
		}
		collector.addText(analysis.AggregateRating.Rationale)
		for _, factor := range analysis.FactorRatings {
			collector.addText(factor.Rationale)
		}
		target := reportTargetFromAnalysis(analysis, structural)
		for _, ref := range analysis.AssessmentRecords {
			if assessment, ok := assessmentsByFile[ref]; ok && assessment.NotAssessed {
				target.NotAssessed = append(target.NotAssessed, assessment.Requirement)
			}
		}
		report.TargetSummary = append(report.TargetSummary, reportTargetSummary(target, len(analysis.AssessmentRecords)))
		if len(context.InScope) == 0 {
			report.Scope.InScope = append(report.Scope.InScope, target.Target)
		}
		report.Targets = append(report.Targets, target)
	}
}

func reportTargetFromAnalysis(analysis AnalysisRecord, structural bool) ReportTarget {
	target := ReportTarget{
		Target:          analysis.Target,
		TargetPath:      analysis.TargetPath,
		AnalysisRecord:  analysis.File,
		LocalRating:     reportRatingFromPointer(analysis.LocalRating),
		AggregateRating: reportRatingFromValue(analysis.AggregateRating),
		FactorRatings:   nonNilFactorRatings(analysis.FactorRatings),
		NotAssessed:     []string{},
		Structural:      structural,
	}
	if structural {
		target.LocalRating = ReportRatingResult{
			Kind:      "structural",
			Rationale: "Structural grouping target; local rating does not apply.",
		}
	}
	return target
}

func reportTargetSummary(target ReportTarget, coveredRequirements int) ReportTargetSummary {
	note := ""
	if target.Structural {
		note = "structural grouping target"
	} else if len(target.NotAssessed) > 0 {
		note = "has not-assessed requirements"
	}
	return ReportTargetSummary{
		Target:              target.Target,
		TargetPath:          target.TargetPath,
		LocalRating:         target.LocalRating,
		AggregateRating:     target.AggregateRating,
		CoveredRequirements: coveredRequirements,
		Note:                note,
	}
}

func addReportRecommendations(report *ReportJSON, recommendations []RecommendationRecord, superseded map[string]bool) {
	for _, rec := range recommendations {
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		report.Recommendations = append(report.Recommendations, ReportRecommendation{
			ID:            strings.TrimSuffix(filepath.Base(rec.File), ".md"),
			Path:          rec.File,
			DoneCriterion: rec.DoneCriterion,
			Active:        !superseded[rec.File] && !superseded[id],
			Supersedes:    append([]string(nil), rec.Supersedes...),
		})
	}
}

func setReportNextAction(report *ReportJSON) {
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

func renderReportMarkdown(report ReportJSON, labels reportDisplayLabels) []byte {
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

func writeReportSummarySection(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("## Summary\n\n")
	out.WriteString("- **Subject:** " + labels.Target(nil, report.Summary.Subject) + "\n")
	if report.Summary.Altitude != "" {
		out.WriteString("- **Altitude:** " + report.Summary.Altitude + "\n")
	} else {
		out.WriteString("- **Altitude:** not recorded\n")
	}
	if report.Summary.Effort != "" {
		out.WriteString("- **Effort:** " + report.Summary.Effort + "\n")
	} else {
		out.WriteString("- **Effort:** not recorded\n")
	}
	if report.Summary.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Summary.Narrowing + "\n")
	}
	out.WriteString("- **Rating:** " + displayRating(report.Summary.Rating, report.Summary.NotAssessed, labels.Ratings) + "\n")
	if report.Summary.Rationale != "" {
		out.WriteString("- **Rationale:** " + report.Summary.Rationale + "\n")
	}
}

func writeReportScopeSection(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
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

func writeReportRisksAndLimitationsSection(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Top Risks and Limitations\n\n")
	risks := riskFindings(report.FindingSummaries)
	summaryRisks := firstFindingSummaries(risks, 8)
	summaryLimitations := firstStrings(report.Limitations, 8)
	if len(summaryRisks) == 0 && len(summaryLimitations) == 0 {
		out.WriteString("No top risks or limitations were recorded in the summary data.\n")
		return
	}
	for _, finding := range summaryRisks {
		out.WriteString("- `" + finding.AssessmentRecord + "`")
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

func writeReportEvidenceSection(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Evidence Basis\n\n")
	if len(report.EvidenceBasis) == 0 {
		out.WriteString("No command or source evidence basis was recorded in findings.\n")
	} else {
		for _, evidence := range report.EvidenceBasis {
			out.WriteString("- **" + evidence.Kind + ":** `" + evidence.Ref + "`\n")
		}
	}
}

func writeReportNextActionSection(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Next Action\n\n")
	if report.NextAction.RecommendationID != "" {
		out.WriteString("- [" + report.NextAction.RecommendationID + "](" + report.NextAction.RecommendationPath + ") - " + report.NextAction.Summary + "\n")
	} else {
		out.WriteString(report.NextAction.Summary + "\n")
	}
}

func writeReportTargetSummarySection(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("\n## Target Summary\n\n")
	out.WriteString("| Target | Local | Aggregate | Covered Requirements | Note |\n")
	out.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, target := range report.TargetSummary {
		out.WriteString("| " + tableCell(labels.Target(target.TargetPath, target.Target)) + " | " + tableCell(displayRatingResult(target.LocalRating, labels.Ratings)) + " | " + tableCell(displayRatingResult(target.AggregateRating, labels.Ratings)) + " | " + fmt.Sprintf("%d", target.CoveredRequirements) + " | " + tableCell(target.Note) + " |\n")
	}
}

func writeReportTargetDetailsSection(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("\n## Target Details\n\n")
	for _, target := range report.Targets {
		writeReportTargetDetail(out, target, labels)
	}
}

func writeReportTargetDetail(out *bytes.Buffer, target ReportTarget, labels reportDisplayLabels) {
	out.WriteString("### " + labels.Target(target.TargetPath, target.Target) + "\n\n")
	out.WriteString("- **Path:** " + displayPath(target.TargetPath) + "\n")
	out.WriteString("- **Local rating:** " + displayRatingResult(target.LocalRating, labels.Ratings) + "\n")
	writeOptionalRationale(out, target.LocalRating.Rationale)
	out.WriteString("- **Aggregate rating:** " + displayRatingResult(target.AggregateRating, labels.Ratings) + "\n")
	writeOptionalRationale(out, target.AggregateRating.Rationale)
	for _, factor := range target.FactorRatings {
		out.WriteString("- **Factor " + labels.Factor(target.TargetPath, factor.Factor) + ":** " + displayRating(factor.Rating, factor.NotAssessed, labels.Ratings) + "\n")
		writeOptionalRationale(out, factor.Rationale)
	}
	out.WriteString("- **Analysis record:** `" + target.AnalysisRecord + "`\n")
	if len(target.NotAssessed) > 0 {
		out.WriteString("- **Not assessed:** " + strings.Join(target.NotAssessed, "; ") + "\n")
	}
	out.WriteString("\n")
}

func writeOptionalRationale(out *bytes.Buffer, rationale string) {
	if rationale != "" {
		out.WriteString("  - " + rationale + "\n")
	}
}

func writeReportRequirementsSection(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("## Requirements\n\n")
	for _, assessment := range report.Assessments {
		out.WriteString("### " + assessment.Requirement + "\n\n")
		state := "active"
		if !assessment.Active {
			state = "superseded"
		}
		out.WriteString("- **State:** " + state + "\n")
		out.WriteString("- **Target:** " + labels.Target(assessment.TargetPath, assessment.Target) + "\n")
		out.WriteString("- **Rating:** " + displayRating(assessment.Rating, assessment.NotAssessed, labels.Ratings) + "\n")
		out.WriteString("- **Assessment record:** `" + assessment.AssessmentRecord + "`\n")
		if assessment.Rationale != "" {
			out.WriteString("- **Rationale:** " + assessment.Rationale + "\n")
		}
		out.WriteString("\n")
	}
}

func writeReportFindingsSection(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("## Findings\n\n")
	for _, finding := range report.FindingSummaries {
		out.WriteString("- `" + finding.AssessmentRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		out.WriteString(": " + finding.Summary + "\n")
	}
}

func writeReportAdviceSection(out *bytes.Buffer, report ReportJSON) {
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

func renderReportSummaryMarkdown(report ReportJSON, labels reportDisplayLabels) []byte {
	var out bytes.Buffer
	out.WriteString("# Quality Evaluation Summary\n\n")
	writeSummaryHeader(&out, report, labels)
	writeSummaryHeadline(&out, report, labels)
	writeSummaryTopRisks(&out, report)
	writeSummaryRatingSummary(&out, report, labels)
	writeSummaryLimitations(&out, report)
	writeSummaryNextAction(&out, report)
	return out.Bytes()
}

func writeSummaryHeader(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	if report.Summary.Run != "" {
		out.WriteString("**Run:** `" + report.Summary.Run + "`\n")
	}
	if report.Summary.Subject != "" {
		out.WriteString("**Subject:** `" + labels.Target(nil, report.Summary.Subject) + "`\n")
	}
	out.WriteString("**Scope:** " + summaryScope(report) + "\n")
	out.WriteString("**Effort:** " + summaryValue(report.Summary.Effort) + "\n")
	out.WriteString("**Root rating:** " + displayRating(report.Summary.Rating, report.Summary.NotAssessed, labels.Ratings) + "\n")
	out.WriteString("**Full report:** [report.md](report.md)\n")
	out.WriteString("**Machine report:** [report.json](report.json)\n\n")
}

func writeSummaryHeadline(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("## Headline\n\n")
	if report.Summary.Rationale != "" {
		out.WriteString(report.Summary.Rationale + "\n")
	} else if report.Summary.NotAssessed {
		out.WriteString("The root target was not assessed in this run.\n")
	} else {
		out.WriteString("The run completed with root rating " + displayRating(report.Summary.Rating, report.Summary.NotAssessed, labels.Ratings) + ".\n")
	}
}

func writeSummaryTopRisks(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Top Risks\n\n")
	risks := firstFindingSummaries(riskFindings(report.FindingSummaries), 5)
	if len(risks) == 0 {
		out.WriteString("None recorded.\n")
		return
	}
	for i, finding := range risks {
		fmt.Fprintf(out, "%d. ", i+1)
		if finding.Severity != "" {
			out.WriteString("**" + finding.Severity + "** - ")
		}
		out.WriteString(finding.Summary)
		writeFindingRecordSuffix(out, finding)
		out.WriteString("\n")
	}
}

func writeFindingRecordSuffix(out *bytes.Buffer, finding FindingSummary) {
	if finding.AssessmentRecord == "" {
		return
	}
	out.WriteString(" (`" + finding.AssessmentRecord + "`")
	if finding.Locator != "" {
		out.WriteString(" at `" + finding.Locator + "`")
	}
	out.WriteString(")")
}

func writeSummaryRatingSummary(out *bytes.Buffer, report ReportJSON, labels reportDisplayLabels) {
	out.WriteString("\n## Rating Summary\n\n")
	if len(report.TargetSummary) == 0 {
		out.WriteString("No target ratings were recorded.\n")
	} else {
		out.WriteString("| Target | Aggregate rating | Reason |\n")
		out.WriteString("| --- | --- | --- |\n")
		for _, target := range report.TargetSummary {
			out.WriteString("| " + tableCell(labels.Target(target.TargetPath, target.Target)) + " | " + tableCell(displayRatingResult(target.AggregateRating, labels.Ratings)) + " | " + tableCell(target.AggregateRating.Rationale) + " |\n")
		}
	}
}

func writeSummaryLimitations(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Limitations\n\n")
	limitations := firstStrings(report.Limitations, 5)
	if len(limitations) == 0 {
		out.WriteString("None recorded.\n")
	} else {
		for _, limitation := range limitations {
			out.WriteString("- " + limitation + "\n")
		}
		if len(report.Limitations) > len(limitations) {
			out.WriteString("- Additional limitations are available in [report.md](report.md).\n")
		}
	}
}

func writeSummaryNextAction(out *bytes.Buffer, report ReportJSON) {
	out.WriteString("\n## Next Action\n\n")
	active := activeRecommendations(report.Recommendations)
	if report.NextAction.RecommendationID != "" {
		out.WriteString(report.NextAction.Summary + "\n\n")
	} else {
		out.WriteString(report.NextAction.Summary + "\n")
	}
	if len(active) > 0 {
		out.WriteString("See active recommendations:\n\n")
		for _, rec := range active {
			out.WriteString("- [" + rec.ID + "](" + rec.Path + ") - " + rec.DoneCriterion + "\n")
		}
	}
}

func summaryScope(report ReportJSON) string {
	if report.Summary.Narrowing != "" {
		return report.Summary.Narrowing
	}
	if report.Scope.Narrowing != "" {
		return report.Scope.Narrowing
	}
	if report.Scope.Description != "" {
		return report.Scope.Description
	}
	return "whole recorded run"
}

func summaryValue(value string) string {
	if value == "" {
		return "not recorded"
	}
	return value
}

func activeRecommendations(recommendations []ReportRecommendation) []ReportRecommendation {
	active := []ReportRecommendation{}
	for _, rec := range recommendations {
		if rec.Active {
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

func firstActiveRecommendation(recommendations []ReportRecommendation) (ReportRecommendation, bool) {
	for _, rec := range recommendations {
		if rec.Active {
			return rec, true
		}
	}
	return ReportRecommendation{}, false
}

type runContext struct {
	Recorded         bool
	Altitude         string
	Narrowing        string
	Effort           string
	ScopeDescription string
	InScope          []string
	OutOfScope       []string
	Limitations      []string
}

func (r *Run) reportContext() runContext {
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
	if value := designParams["Effort"]; value != "" {
		context.Effort = strings.ToLower(cleanInlineCode(value))
		context.Recorded = true
	}
	if value := designParams["Scope"]; value != "" {
		context.ScopeDescription = cleanInlineCode(value)
		context.Recorded = true
	} else if value := designParams["Scope description"]; value != "" {
		context.ScopeDescription = cleanInlineCode(value)
		context.Recorded = true
	}
	if effort := firstParagraph(markdownSection(r.Plan, "Effort")); effort != "" {
		context.Recorded = true
		if context.Effort == "" {
			context.Effort = inferEffort(strings.ToLower(effort))
		}
		if context.ScopeDescription == "" {
			context.ScopeDescription = effort
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
		context.Effort = inferEffort(parts[2])
	}
	return context
}

func inferEffort(narrowing string) string {
	for _, effort := range []string{"quick", "standard", "deep"} {
		if strings.Contains(narrowing, effort) {
			return effort
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

func reportRatingFromPointer(result *RatingResult) ReportRatingResult {
	if result == nil {
		return ReportRatingResult{
			Kind:        "not-assessed",
			NotAssessed: true,
			Rationale:   "No local rating was recorded.",
		}
	}
	return reportRatingFromValue(*result)
}

func reportRatingFromValue(result RatingResult) ReportRatingResult {
	kind := "rated"
	if result.NotAssessed || result.Rating == nil {
		kind = "not-assessed"
	}
	return ReportRatingResult{
		Kind:        kind,
		Rating:      result.Rating,
		NotAssessed: result.NotAssessed,
		Rationale:   result.Rationale,
	}
}

func nonNilFactorRatings(factors []FactorRating) []FactorRating {
	if factors == nil {
		return []FactorRating{}
	}
	return factors
}

func riskFindings(findings []FindingSummary) []FindingSummary {
	risks := []FindingSummary{}
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

func firstFindingSummaries(items []FindingSummary, limit int) []FindingSummary {
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
		if strings.TrimSpace(factor.Title) != "" {
			l.Factors[factorLabelKey(targetPath, key)] = factor.Title
		}
		l.addFactors(targetPath, factor.Factors)
	}
}

func (l reportDisplayLabels) Target(path []string, fallback string) string {
	if title := l.Targets[targetLabelKey(path)]; title != "" {
		return title
	}
	return fallback
}

func (l reportDisplayLabels) Factor(targetPath []string, fallback string) string {
	for i := len(targetPath); i >= 0; i-- {
		if title := l.Factors[factorLabelKey(targetPath[:i], fallback)]; title != "" {
			return title
		}
	}
	return fallback
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

func displayRating(rating *string, notAssessed bool, ratingLabels map[string]string) string {
	if notAssessed || rating == nil {
		return "not assessed"
	}
	if title := ratingLabels[*rating]; title != "" {
		return title
	}
	return *rating
}

func displayRatingResult(result ReportRatingResult, ratingLabels map[string]string) string {
	switch result.Kind {
	case "structural":
		return "n/a (structural)"
	default:
		return displayRating(result.Rating, result.NotAssessed, ratingLabels)
	}
}

func displayList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(items, "; ")
}

func displayScopeList(items []string, targets []ReportTargetSummary, labels reportDisplayLabels) string {
	if len(items) == 0 {
		return "none"
	}
	targetLabels := map[string]string{}
	for _, target := range targets {
		if _, exists := targetLabels[target.Target]; exists {
			targetLabels[target.Target] = ""
			continue
		}
		targetLabels[target.Target] = labels.Target(target.TargetPath, target.Target)
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

func targetLabelKey(path []string) string {
	return strings.Join(path, "\x00")
}

func factorLabelKey(targetPath []string, factor string) string {
	return targetLabelKey(targetPath) + "\x00" + factor
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

func Gate(result *BuildResult, scale []string, threshold string) (bool, error) {
	if threshold == "" {
		return true, nil
	}
	thresholdIndex := -1
	ratingIndex := -1
	for i, level := range scale {
		if level == threshold {
			thresholdIndex = i
		}
		if result.Rating != nil && level == *result.Rating {
			ratingIndex = i
		}
	}
	if thresholdIndex < 0 {
		return false, usagef("--at-or-below level %q is not in the run rating scale", threshold)
	}
	if result.NotAssessed || result.Rating == nil {
		return false, nil
	}
	if ratingIndex < 0 {
		return false, fmt.Errorf("root rating %q is not in the run rating scale", *result.Rating)
	}
	return ratingIndex < thresholdIndex, nil
}

func GateReport(path, threshold string) (*GateResult, error) {
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
	var report ReportJSON
	if err := DecodeSingleJSON(raw, &report); err != nil {
		return nil, fmt.Errorf("reading report.json: %w", err)
	}
	levels := make([]string, 0, len(run.Scale))
	for _, level := range run.Scale {
		levels = append(levels, level.Level)
	}
	build := &BuildResult{Rating: report.Rating.Rating, NotAssessed: report.Rating.NotAssessed}
	pass, err := Gate(build, levels, threshold)
	if err != nil {
		return nil, err
	}
	return &GateResult{
		SchemaVersion: SchemaVersion,
		Path:          filepath.ToSlash(reportPath),
		Threshold:     threshold,
		Rating:        report.Rating.Rating,
		NotAssessed:   report.Rating.NotAssessed,
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
