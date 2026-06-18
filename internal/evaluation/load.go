package evaluation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

type Run struct {
	Path            string
	Design          string
	Plan            string
	PlannedCoverage *PlannedCoverage
	Assessments     []AssessmentRecord
	Analyses        []AnalysisRecord
	Recommendations []RecommendationRecord
	Scale           []model.RatingLevel
}

type Counts struct {
	Assessments     int `json:"assessments"`
	Analyses        int `json:"analyses"`
	Recommendations int `json:"recommendations"`
}

type Gap struct {
	Kind   string `json:"kind"`
	Ref    string `json:"ref"`
	Detail string `json:"detail"`
}

type Status struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Reportable    bool             `json:"reportable"`
	Counts        Counts           `json:"counts"`
	Gaps          []Gap            `json:"gaps"`
	NextActions   []receipt.Action `json:"nextActions"`
}

func Load(path string) (*Run, error) {
	runAbs, err := verifyRun(path)
	if err != nil {
		return nil, err
	}
	doc, err := document.Parse(filepath.Join(runAbs, "model.md"))
	if err != nil {
		return nil, err
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, err
	}
	run := &Run{Path: filepath.ToSlash(runAbs), Scale: spec.RatingScale}
	if raw, err := os.ReadFile(filepath.Join(runAbs, "design.md")); err != nil {
		return nil, fmt.Errorf("reading design.md: %w", err)
	} else {
		run.Design = string(raw)
	}
	if raw, err := os.ReadFile(filepath.Join(runAbs, "plan.md")); err != nil {
		return nil, fmt.Errorf("reading plan.md: %w", err)
	} else {
		run.Plan = string(raw)
	}
	if raw, err := os.ReadFile(filepath.Join(runAbs, plannedCoverageFile)); err == nil {
		var coverage PlannedCoverage
		if err := DecodeSingleJSON(raw, &coverage); err != nil {
			return nil, fmt.Errorf("reading %s: %w", plannedCoverageFile, err)
		}
		if err := validatePlannedCoverage(coverage); err != nil {
			return nil, fmt.Errorf("reading %s: %w", plannedCoverageFile, err)
		}
		sortPlannedCoverage(&coverage)
		run.PlannedCoverage = &coverage
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading %s: %w", plannedCoverageFile, err)
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "assessments"), "*.json", func(path string, raw []byte) error {
		var rec AssessmentRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("assessments", filepath.Base(path)))
		run.Assessments = append(run.Assessments, rec)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "analysis"), "*.json", func(path string, raw []byte) error {
		var rec AnalysisRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("analysis", filepath.Base(path)))
		run.Analyses = append(run.Analyses, rec)
		return nil
	}); err != nil {
		return nil, err
	}
	if err := loadRecommendations(filepath.Join(runAbs, "recommendations"), &run.Recommendations); err != nil {
		return nil, err
	}
	return run, nil
}

func loadJSONRecords(dir, pattern string, decode func(string, []byte) error) error {
	paths, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return err
	}
	slices.Sort(paths)
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		if err := decode(path, raw); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	return nil
}

func loadRecommendations(dir string, out *[]RecommendationRecord) error {
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		return err
	}
	slices.Sort(paths)
	for _, path := range paths {
		rec, err := parseRecommendation(path)
		if err != nil {
			return err
		}
		rec.File = filepath.ToSlash(filepath.Join("recommendations", filepath.Base(path)))
		*out = append(*out, rec)
	}
	return nil
}

func parseRecommendation(path string) (RecommendationRecord, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return RecommendationRecord{}, fmt.Errorf("reading %s: %w", path, err)
	}
	front, body, err := splitMarkdownFrontmatter(raw)
	if err != nil {
		return RecommendationRecord{}, fmt.Errorf("%s: %w", path, err)
	}
	var rec RecommendationRecord
	if err := yaml.Unmarshal(front, &rec); err != nil {
		return RecommendationRecord{}, fmt.Errorf("%s: parsing frontmatter: %w", path, err)
	}
	if rec.SchemaVersion != SchemaVersion {
		return RecommendationRecord{}, fmt.Errorf("%s: schemaVersion = %d, want %d", path, rec.SchemaVersion, SchemaVersion)
	}
	rec.Body = string(body)
	return rec, nil
}

func splitMarkdownFrontmatter(raw []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(raw, []byte("---\n")) {
		return nil, nil, fmt.Errorf("missing runtime frontmatter")
	}
	rest := raw[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, nil, fmt.Errorf("unterminated runtime frontmatter")
	}
	return rest[:end], rest[end+len("\n---\n"):], nil
}

func (r *Run) Counts() Counts {
	return Counts{
		Assessments:     len(r.Assessments),
		Analyses:        len(r.Analyses),
		Recommendations: len(r.Recommendations),
	}
}

func (r *Run) Status() Status {
	gaps := r.Renderable()
	status := Status{
		SchemaVersion: SchemaVersion,
		Path:          r.Path,
		Reportable:    len(gaps) == 0,
		Counts:        r.Counts(),
		Gaps:          gaps,
	}
	if status.Reportable {
		status.NextActions = []receipt.Action{{
			ID:      "build-report",
			Label:   "Build the evaluation report",
			Command: "qualitymd evaluation build-report " + r.Path,
		}}
	} else if hasReviewGap(gaps) {
		status.NextActions = []receipt.Action{{
			ID:      "review-gaps",
			Label:   "Resolve evaluation record gaps",
			Command: "qualitymd evaluation show-status " + r.Path,
		}}
	} else {
		status.NextActions = []receipt.Action{{
			ID:      "add-record",
			Label:   "Add the missing evaluation records",
			Command: "qualitymd evaluation add-record <kind> " + r.Path,
		}}
	}
	return status
}

func (r *Run) Renderable() []Gap {
	var gaps []Gap
	assessments := map[string]bool{}
	supersededAssessments, supersedingGaps := r.assessmentSupersedingState()
	gaps = append(gaps, supersedingGaps...)
	assessmentRecordsByIdentity := map[string]string{}
	assessmentIdentities := map[string]string{}
	for _, rec := range r.Assessments {
		assessments[rec.File] = true
		if supersededAssessments[rec.File] {
			continue
		}
		key := assessmentIdentity(rec)
		assessmentRecordsByIdentity[key] = rec.File
		if prior, ok := assessmentIdentities[key]; ok {
			gaps = append(gaps, Gap{Kind: "duplicate-assessment", Ref: rec.File, Detail: "duplicates " + prior})
			continue
		}
		assessmentIdentities[key] = rec.File
	}
	analyses := map[string]AnalysisRecord{}
	analysisRecordsByIdentity := map[string]string{}
	for _, rec := range r.Analyses {
		analyses[rec.File] = rec
		analysisRecordsByIdentity[analysisIdentity(rec)] = rec.File
	}
	gaps = append(gaps, r.plannedCoverageGaps(assessmentRecordsByIdentity, analysisRecordsByIdentity, supersededAssessments)...)
	recs := map[string]bool{}
	for _, rec := range r.Recommendations {
		recs[rec.File] = true
		recs[strings.TrimSuffix(filepath.Base(rec.File), ".md")] = true
	}
	for _, rec := range r.Recommendations {
		for _, ref := range rec.Supersedes {
			if recommendationRefExists(recs, ref) {
				continue
			}
			gaps = append(gaps, Gap{Kind: "missing-superseded-recommendation", Ref: ref, Detail: "referenced by " + rec.File})
		}
	}
	if len(r.Analyses) == 0 {
		gaps = append(gaps, Gap{Kind: "missing-analysis", Ref: "analysis/", Detail: "no analysis records are present"})
	}
	rootAnalyses := 0
	for _, analysis := range r.Analyses {
		if len(analysis.TargetPath) == 0 {
			rootAnalyses++
		}
		for _, ref := range analysis.AssessmentRecords {
			if !assessments[ref] {
				gaps = append(gaps, Gap{Kind: "missing-assessment", Ref: ref, Detail: "referenced by " + analysis.File})
				continue
			}
			if supersededAssessments[ref] {
				gaps = append(gaps, Gap{Kind: "superseded-assessment-reference", Ref: ref, Detail: "referenced by " + analysis.File})
			}
		}
		for _, ref := range analysis.ChildAnalysisRecords {
			if _, ok := analyses[ref]; !ok {
				gaps = append(gaps, Gap{Kind: "missing-analysis", Ref: ref, Detail: "referenced by " + analysis.File})
			}
		}
	}
	if len(r.Analyses) > 0 && rootAnalyses == 0 {
		gaps = append(gaps, Gap{Kind: "missing-root-analysis", Ref: "analysis/", Detail: "no analysis record has an empty targetPath for the in-scope root"})
	}
	if rootAnalyses > 1 {
		gaps = append(gaps, Gap{Kind: "duplicate-root-analysis", Ref: "analysis/", Detail: "multiple analysis records have an empty targetPath"})
	}
	for _, assessment := range r.Assessments {
		for _, ref := range assessment.Recommendations {
			if recs[ref] || recs["recommendations/"+ref+".md"] {
				continue
			}
			gaps = append(gaps, Gap{Kind: "missing-recommendation", Ref: ref, Detail: "referenced by " + assessment.File})
		}
	}
	return gaps
}

func assessmentIdentity(rec AssessmentRecord) string {
	return strings.Join(rec.TargetPath, "\x00") + "\x00" + rec.Requirement
}

func assessmentID(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".json")
}

type knownAssessment struct {
	File   string
	Record AssessmentRecord
}

func (r *Run) assessmentSupersedingState() (map[string]bool, []Gap) {
	known := map[string]knownAssessment{}
	superseded := map[string]bool{}
	var gaps []Gap
	for _, rec := range r.Assessments {
		for _, ref := range rec.Supersedes {
			ref = strings.TrimSpace(ref)
			prior, ok := resolveKnownAssessment(known, ref)
			if !ok {
				gaps = append(gaps, Gap{Kind: "missing-superseded-assessment", Ref: ref, Detail: "referenced by " + rec.File})
				continue
			}
			if assessmentIdentity(prior.Record) != assessmentIdentity(rec) {
				gaps = append(gaps, Gap{Kind: "invalid-assessment-supersedes", Ref: ref, Detail: "referenced by " + rec.File + " with different targetPath or requirement"})
				continue
			}
			superseded[prior.File] = true
			superseded[assessmentID(prior.File)] = true
		}
		known[rec.File] = knownAssessment{File: rec.File, Record: rec}
		known[assessmentID(rec.File)] = knownAssessment{File: rec.File, Record: rec}
	}
	return superseded, gaps
}

func resolveKnownAssessment(known map[string]knownAssessment, ref string) (knownAssessment, bool) {
	if rec, ok := known[ref]; ok {
		return rec, true
	}
	if strings.HasPrefix(ref, "assessments/") {
		rec, ok := known[strings.TrimSuffix(filepath.Base(ref), ".json")]
		return rec, ok
	}
	rec, ok := known["assessments/"+ref+".json"]
	return rec, ok
}

func analysisIdentity(rec AnalysisRecord) string {
	return strings.Join(rec.TargetPath, "\x00")
}

func recommendationRefExists(recs map[string]bool, ref string) bool {
	ref = strings.TrimSpace(ref)
	if recs[ref] {
		return true
	}
	if strings.HasPrefix(ref, "recommendations/") {
		return recs[strings.TrimSuffix(filepath.Base(ref), ".md")]
	}
	return recs["recommendations/"+ref+".md"]
}

func (r *Run) plannedCoverageGaps(assessmentRecordsByIdentity, analysisRecordsByIdentity map[string]string, supersededAssessments map[string]bool) []Gap {
	if r.PlannedCoverage == nil {
		return nil
	}
	plannedAssessments := map[string]PlannedCoverageAssessment{}
	for _, assessment := range r.PlannedCoverage.Assessments {
		plannedAssessments[plannedAssessmentIdentity(assessment)] = assessment
	}
	plannedAnalyses := map[string]PlannedCoverageAnalysis{}
	for _, analysis := range r.PlannedCoverage.Analyses {
		plannedAnalyses[plannedAnalysisIdentity(analysis)] = analysis
	}
	var gaps []Gap
	for _, assessment := range r.PlannedCoverage.Assessments {
		key := plannedAssessmentIdentity(assessment)
		if _, ok := assessmentRecordsByIdentity[key]; ok {
			continue
		}
		gaps = append(gaps, Gap{
			Kind:   "missing-planned-assessment",
			Ref:    plannedCoverageFile,
			Detail: describePlannedAssessment(assessment),
		})
	}
	for _, analysis := range r.PlannedCoverage.Analyses {
		key := plannedAnalysisIdentity(analysis)
		if _, ok := analysisRecordsByIdentity[key]; ok {
			continue
		}
		gaps = append(gaps, Gap{
			Kind:   "missing-planned-analysis",
			Ref:    plannedCoverageFile,
			Detail: describePlannedAnalysis(analysis),
		})
	}
	for _, assessment := range r.Assessments {
		if supersededAssessments[assessment.File] {
			continue
		}
		planned := PlannedCoverageAssessment{TargetPath: assessment.TargetPath, Requirement: assessment.Requirement}
		if _, ok := plannedAssessments[plannedAssessmentIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, Gap{
			Kind:   "unexpected-assessment",
			Ref:    assessment.File,
			Detail: describePlannedAssessment(planned),
		})
	}
	for _, analysis := range r.Analyses {
		planned := PlannedCoverageAnalysis{TargetPath: analysis.TargetPath}
		if _, ok := plannedAnalyses[plannedAnalysisIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, Gap{
			Kind:   "unexpected-analysis",
			Ref:    analysis.File,
			Detail: describePlannedAnalysis(planned),
		})
	}
	return gaps
}

func hasReviewGap(gaps []Gap) bool {
	for _, gap := range gaps {
		if gap.Kind == "duplicate-assessment" || gap.Kind == "duplicate-root-analysis" || gap.Kind == "missing-superseded-recommendation" || gap.Kind == "missing-superseded-assessment" || gap.Kind == "invalid-assessment-supersedes" || gap.Kind == "superseded-assessment-reference" || strings.HasPrefix(gap.Kind, "unexpected-") {
			return true
		}
	}
	return false
}
