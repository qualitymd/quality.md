package evaluation

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

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
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	ReportMD      string           `json:"reportMd"`
	ReportJSON    string           `json:"reportJson"`
	Rating        *string          `json:"rating"`
	NotAssessed   bool             `json:"notAssessed"`
	NextActions   []receipt.Action `json:"nextActions,omitempty"`
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
	md := renderReportMarkdown(report)
	js, err := marshalJSON(report)
	if err != nil {
		return nil, err
	}
	reportMD := filepath.Join(run.Path, "report.md")
	reportJSON := filepath.Join(run.Path, "report.json")
	if err := writeReportFile(reportMD, md); err != nil {
		return nil, err
	}
	if err := writeReportFile(reportJSON, js); err != nil {
		return nil, err
	}
	return &BuildResult{
		SchemaVersion: SchemaVersion,
		Path:          run.Path,
		ReportMD:      filepath.ToSlash(reportMD),
		ReportJSON:    filepath.ToSlash(reportJSON),
		Rating:        report.Rating.Rating,
		NotAssessed:   report.Rating.NotAssessed,
	}, nil
}

func (r *Run) Report() (ReportJSON, error) {
	analyses := append([]AnalysisRecord(nil), r.Analyses...)
	if len(analyses) == 0 {
		return ReportJSON{}, fmt.Errorf("run has no analysis records")
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
				return ReportJSON{}, fmt.Errorf("run has multiple root analysis records")
			}
			candidate := analysis
			root = &candidate
		}
	}
	if root == nil {
		return ReportJSON{}, fmt.Errorf("run has no root analysis record")
	}
	rootAnalysis := *root
	context := r.reportContext()
	report := ReportJSON{
		SchemaVersion: SchemaVersion,
		Summary: ReportSummary{
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

	assessmentsByFile := map[string]AssessmentRecord{}
	evidenceSeen := map[string]bool{}
	limitationSeen := map[string]bool{}
	addLimitation := func(limitation string) {
		limitation = strings.TrimSpace(limitation)
		key := limitationKey(limitation)
		if key == "" || limitationSeen[key] {
			return
		}
		limitationSeen[key] = true
		report.Limitations = append(report.Limitations, limitation)
	}
	addLimitations := func(text string) {
		for _, limitation := range limitationSentences(text) {
			addLimitation(limitation)
		}
	}
	for _, limitation := range context.Limitations {
		addLimitations(limitation)
		addLimitation(limitation)
	}
	addLimitations(rootAnalysis.AggregateRating.Rationale)
	supersededAssessments, _ := r.assessmentSupersedingState()
	for _, assessment := range r.Assessments {
		active := !supersededAssessments[assessment.File]
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
			addLimitations(assessment.Requirement + ": " + assessment.Rationale)
		}
		addLimitations(assessment.Rationale)
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
				report.EvidenceBasis = append(report.EvidenceBasis, ReportEvidence{
					Kind: evidence.Kind,
					Ref:  evidence.Ref,
				})
			}
		}
	}
	for _, analysis := range analyses {
		structural := len(analysis.AssessmentRecords) == 0 && len(analysis.ChildAnalysisRecords) > 0
		if analysis.LocalRating != nil {
			addLimitations(analysis.LocalRating.Rationale)
		}
		addLimitations(analysis.AggregateRating.Rationale)
		for _, factor := range analysis.FactorRatings {
			addLimitations(factor.Rationale)
		}
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
		for _, ref := range analysis.AssessmentRecords {
			if assessment, ok := assessmentsByFile[ref]; ok && assessment.NotAssessed {
				target.NotAssessed = append(target.NotAssessed, assessment.Requirement)
			}
		}
		note := ""
		if structural {
			note = "structural grouping target"
		} else if len(target.NotAssessed) > 0 {
			note = "has not-assessed requirements"
		}
		report.TargetSummary = append(report.TargetSummary, ReportTargetSummary{
			Target:              target.Target,
			TargetPath:          target.TargetPath,
			LocalRating:         target.LocalRating,
			AggregateRating:     target.AggregateRating,
			CoveredRequirements: len(analysis.AssessmentRecords),
			Note:                note,
		})
		if len(context.InScope) == 0 {
			report.Scope.InScope = append(report.Scope.InScope, target.Target)
		}
		report.Targets = append(report.Targets, target)
	}
	supersededRecommendations := r.supersededRecommendations()
	for _, rec := range r.Recommendations {
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		report.Recommendations = append(report.Recommendations, ReportRecommendation{
			ID:            strings.TrimSuffix(filepath.Base(rec.File), ".md"),
			Path:          rec.File,
			DoneCriterion: rec.DoneCriterion,
			Active:        !supersededRecommendations[rec.File] && !supersededRecommendations[id],
			Supersedes:    append([]string(nil), rec.Supersedes...),
		})
	}
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
	return report, nil
}

func renderReportMarkdown(report ReportJSON) []byte {
	var out bytes.Buffer
	out.WriteString("# Evaluation Report\n\n")
	out.WriteString("## Summary\n\n")
	out.WriteString("- **Subject:** " + report.Summary.Subject + "\n")
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
	out.WriteString("- **Rating:** " + displayRating(report.Summary.Rating, report.Summary.NotAssessed) + "\n")
	if report.Summary.Rationale != "" {
		out.WriteString("- **Rationale:** " + report.Summary.Rationale + "\n")
	}
	out.WriteString("\n## Scope\n\n")
	out.WriteString(report.Scope.Description + "\n\n")
	if report.Scope.Narrowing != "" {
		out.WriteString("- **Narrowing:** " + report.Scope.Narrowing + "\n")
	} else {
		out.WriteString("- **Narrowing:** whole recorded run\n")
	}
	out.WriteString("- **In scope:** " + displayList(report.Scope.InScope) + "\n")
	if len(report.Scope.OutOfScope) > 0 {
		out.WriteString("- **Out of scope:** " + displayList(report.Scope.OutOfScope) + "\n")
	} else {
		out.WriteString("- **Out of scope:** not recorded\n")
	}
	if len(report.Scope.NotRecorded) > 0 {
		out.WriteString("- **Metadata not recorded:** " + displayList(report.Scope.NotRecorded) + "\n")
	}
	out.WriteString("\n## Top Risks and Limitations\n\n")
	risks := riskFindings(report.FindingSummaries)
	summaryRisks := firstFindingSummaries(risks, 8)
	summaryLimitations := firstStrings(report.Limitations, 8)
	if len(summaryRisks) == 0 && len(summaryLimitations) == 0 {
		out.WriteString("No top risks or limitations were recorded in the summary data.\n")
	} else {
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
	out.WriteString("\n## Evidence Basis\n\n")
	if len(report.EvidenceBasis) == 0 {
		out.WriteString("No command or source evidence basis was recorded in findings.\n")
	} else {
		for _, evidence := range report.EvidenceBasis {
			out.WriteString("- **" + evidence.Kind + ":** `" + evidence.Ref + "`\n")
		}
	}
	out.WriteString("\n## Next Action\n\n")
	if report.NextAction.RecommendationID != "" {
		out.WriteString("- [" + report.NextAction.RecommendationID + "](" + report.NextAction.RecommendationPath + ") - " + report.NextAction.Summary + "\n")
	} else {
		out.WriteString(report.NextAction.Summary + "\n")
	}
	out.WriteString("\n## Target Summary\n\n")
	out.WriteString("| Target | Local | Aggregate | Covered Requirements | Note |\n")
	out.WriteString("| --- | --- | --- | --- | --- |\n")
	for _, target := range report.TargetSummary {
		out.WriteString("| " + tableCell(target.Target) + " | " + tableCell(displayRatingResult(target.LocalRating)) + " | " + tableCell(displayRatingResult(target.AggregateRating)) + " | " + fmt.Sprintf("%d", target.CoveredRequirements) + " | " + tableCell(target.Note) + " |\n")
	}
	out.WriteString("\n## Target Details\n\n")
	for _, target := range report.Targets {
		out.WriteString("### " + target.Target + "\n\n")
		out.WriteString("- **Path:** " + displayPath(target.TargetPath) + "\n")
		out.WriteString("- **Local rating:** " + displayRatingResult(target.LocalRating) + "\n")
		if target.LocalRating.Rationale != "" {
			out.WriteString("  - " + target.LocalRating.Rationale + "\n")
		}
		out.WriteString("- **Aggregate rating:** " + displayRatingResult(target.AggregateRating) + "\n")
		if target.AggregateRating.Rationale != "" {
			out.WriteString("  - " + target.AggregateRating.Rationale + "\n")
		}
		for _, factor := range target.FactorRatings {
			out.WriteString("- **Factor " + factor.Factor + ":** " + displayRating(factor.Rating, factor.NotAssessed) + "\n")
			if factor.Rationale != "" {
				out.WriteString("  - " + factor.Rationale + "\n")
			}
		}
		out.WriteString("- **Analysis record:** `" + target.AnalysisRecord + "`\n")
		if len(target.NotAssessed) > 0 {
			out.WriteString("- **Not assessed:** " + strings.Join(target.NotAssessed, "; ") + "\n")
		}
		out.WriteString("\n")
	}
	out.WriteString("## Requirements\n\n")
	for _, assessment := range report.Assessments {
		out.WriteString("### " + assessment.Requirement + "\n\n")
		state := "active"
		if !assessment.Active {
			state = "superseded"
		}
		out.WriteString("- **State:** " + state + "\n")
		out.WriteString("- **Target:** " + assessment.Target + "\n")
		out.WriteString("- **Rating:** " + displayRating(assessment.Rating, assessment.NotAssessed) + "\n")
		out.WriteString("- **Assessment record:** `" + assessment.AssessmentRecord + "`\n")
		if assessment.Rationale != "" {
			out.WriteString("- **Rationale:** " + assessment.Rationale + "\n")
		}
		out.WriteString("\n")
	}
	out.WriteString("## Findings\n\n")
	for _, finding := range report.FindingSummaries {
		out.WriteString("- `" + finding.AssessmentRecord + "`")
		if finding.Locator != "" {
			out.WriteString(" at `" + finding.Locator + "`")
		}
		out.WriteString(": " + finding.Summary + "\n")
	}
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
	return out.Bytes()
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

func displayRating(rating *string, notAssessed bool) string {
	if notAssessed || rating == nil {
		return "not assessed"
	}
	return *rating
}

func displayRatingResult(result ReportRatingResult) string {
	switch result.Kind {
	case "structural":
		return "n/a (structural)"
	default:
		return displayRating(result.Rating, result.NotAssessed)
	}
}

func displayList(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	return strings.Join(items, "; ")
}

func displayPath(path []string) string {
	if len(path) == 0 {
		return "(root)"
	}
	return strings.Join(path, " / ")
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
		return false, usagef("--fail-at-or-below level %q is not in the run rating scale", threshold)
	}
	if result.NotAssessed || result.Rating == nil {
		return false, nil
	}
	if ratingIndex < 0 {
		return false, fmt.Errorf("root rating %q is not in the run rating scale", *result.Rating)
	}
	return ratingIndex < thresholdIndex, nil
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
