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
	Rating           ReportRating           `json:"rating"`
	Targets          []ReportTarget         `json:"targets"`
	Assessments      []ReportAssessment     `json:"assessments"`
	FindingSummaries []FindingSummary       `json:"findingSummaries"`
	Recommendations  []ReportRecommendation `json:"recommendations"`
}

type ReportRating struct {
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Subject     string  `json:"subject"`
	Rationale   string  `json:"rationale"`
}

type ReportTarget struct {
	Target             string         `json:"target"`
	TargetPath         []string       `json:"targetPath"`
	LocalRating        *string        `json:"localRating,omitempty"`
	LocalRationale     string         `json:"localRationale,omitempty"`
	AggregateRating    *string        `json:"aggregateRating,omitempty"`
	AggregateRationale string         `json:"aggregateRationale,omitempty"`
	FactorRatings      []FactorRating `json:"factorRatings,omitempty"`
	AnalysisRecord     string         `json:"analysisRecord"`
	NotAssessed        []string       `json:"notAssessed,omitempty"`
}

type ReportAssessment struct {
	AssessmentRecord string   `json:"assessmentRecord"`
	Target           string   `json:"target"`
	TargetPath       []string `json:"targetPath"`
	Requirement      string   `json:"requirement"`
	Rating           *string  `json:"rating"`
	NotAssessed      bool     `json:"notAssessed"`
	Rationale        string   `json:"rationale"`
}

type FindingSummary struct {
	AssessmentRecord string `json:"assessmentRecord"`
	Locator          string `json:"locator"`
	Category         string `json:"category"`
	Severity         string `json:"severity,omitempty"`
	Summary          string `json:"summary"`
}

type ReportRecommendation struct {
	ID            string `json:"id"`
	Path          string `json:"path"`
	DoneCriterion string `json:"doneCriterion"`
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
	slices.SortFunc(analyses, func(a, b AnalysisRecord) int {
		if cmp := len(a.TargetPath) - len(b.TargetPath); cmp != 0 {
			return cmp
		}
		return strings.Compare(strings.Join(a.TargetPath, "/"), strings.Join(b.TargetPath, "/"))
	})
	root := analyses[0]
	for _, analysis := range analyses {
		if len(analysis.TargetPath) == 0 {
			root = analysis
			break
		}
	}
	report := ReportJSON{
		SchemaVersion: SchemaVersion,
		Rating: ReportRating{
			Rating:      root.AggregateRating.Rating,
			NotAssessed: root.AggregateRating.NotAssessed,
			Subject:     root.Target,
			Rationale:   root.AggregateRating.Rationale,
		},
	}
	assessmentsByFile := map[string]AssessmentRecord{}
	for _, assessment := range r.Assessments {
		assessmentsByFile[assessment.File] = assessment
		report.Assessments = append(report.Assessments, ReportAssessment{
			AssessmentRecord: assessment.File,
			Target:           assessment.Target,
			TargetPath:       assessment.TargetPath,
			Requirement:      assessment.Requirement,
			Rating:           assessment.Rating,
			NotAssessed:      assessment.NotAssessed,
			Rationale:        assessment.Rationale,
		})
		for _, finding := range assessment.Findings {
			report.FindingSummaries = append(report.FindingSummaries, FindingSummary{
				AssessmentRecord: assessment.File,
				Locator:          finding.Locator,
				Category:         finding.Category,
				Severity:         finding.Severity,
				Summary:          finding.Observation,
			})
		}
	}
	for _, analysis := range analyses {
		target := ReportTarget{
			Target:             analysis.Target,
			TargetPath:         analysis.TargetPath,
			AnalysisRecord:     analysis.File,
			AggregateRating:    analysis.AggregateRating.Rating,
			AggregateRationale: analysis.AggregateRating.Rationale,
			FactorRatings:      analysis.FactorRatings,
		}
		if analysis.LocalRating != nil {
			target.LocalRating = analysis.LocalRating.Rating
			target.LocalRationale = analysis.LocalRating.Rationale
		}
		for _, ref := range analysis.AssessmentRecords {
			if assessment, ok := assessmentsByFile[ref]; ok && assessment.NotAssessed {
				target.NotAssessed = append(target.NotAssessed, assessment.Requirement)
			}
		}
		report.Targets = append(report.Targets, target)
	}
	for _, rec := range r.Recommendations {
		report.Recommendations = append(report.Recommendations, ReportRecommendation{
			ID:            strings.TrimSuffix(filepath.Base(rec.File), ".md"),
			Path:          rec.File,
			DoneCriterion: rec.DoneCriterion,
		})
	}
	return report, nil
}

func renderReportMarkdown(report ReportJSON) []byte {
	var out bytes.Buffer
	out.WriteString("# Evaluation Report\n\n")
	out.WriteString("## Rating\n\n")
	out.WriteString("- **Subject:** " + report.Rating.Subject + "\n")
	out.WriteString("- **Rating:** " + displayRating(report.Rating.Rating, report.Rating.NotAssessed) + "\n")
	if report.Rating.Rationale != "" {
		out.WriteString("- **Rationale:** " + report.Rating.Rationale + "\n")
	}
	out.WriteString("\n## Targets\n\n")
	for _, target := range report.Targets {
		out.WriteString("### " + target.Target + "\n\n")
		out.WriteString("- **Path:** " + displayPath(target.TargetPath) + "\n")
		out.WriteString("- **Local rating:** " + displayRating(target.LocalRating, false) + "\n")
		if target.LocalRationale != "" {
			out.WriteString("  - " + target.LocalRationale + "\n")
		}
		out.WriteString("- **Aggregate rating:** " + displayRating(target.AggregateRating, false) + "\n")
		if target.AggregateRationale != "" {
			out.WriteString("  - " + target.AggregateRationale + "\n")
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
	for _, rec := range report.Recommendations {
		out.WriteString("- [" + rec.ID + "](" + rec.Path + ") — " + rec.DoneCriterion + "\n")
	}
	return out.Bytes()
}

func displayRating(rating *string, notAssessed bool) string {
	if notAssessed || rating == nil {
		return "not assessed"
	}
	return *rating
}

func displayPath(path []string) string {
	if len(path) == 0 {
		return "(root)"
	}
	return strings.Join(path, " / ")
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
