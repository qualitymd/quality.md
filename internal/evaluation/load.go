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

type EvaluationRun struct {
	Path              string
	Design            string
	Plan              string
	PlannedCoverage   *PlannedCoverage
	PlanCoverageGaps  []EvaluationRunGap
	AssessmentResults []AssessmentResultRecord
	Analyses          []AnalysisRecord
	Recommendations   []RecommendationRecord
	Model             *model.Spec
	Scale             []model.RatingLevel
}

type EvaluationRecordCounts struct {
	AssessmentResults int `json:"assessmentResults"`
	Analyses          int `json:"analyses"`
	Recommendations   int `json:"recommendations"`
}

type EvaluationRunGap struct {
	Kind   string `json:"kind"`
	Ref    string `json:"ref"`
	Detail string `json:"detail"`
}

type EvaluationRunStatus struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Path          string                 `json:"path"`
	Reportable    bool                   `json:"reportable"`
	Counts        EvaluationRecordCounts `json:"counts"`
	Gaps          []EvaluationRunGap     `json:"gaps"`
	NextActions   []receipt.Action       `json:"nextActions"`
}

func Load(path string) (*EvaluationRun, error) {
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
	run := &EvaluationRun{Path: filepath.ToSlash(runAbs), Model: spec, Scale: spec.RatingScale}
	if raw, err := os.ReadFile(filepath.Join(runAbs, "design.md")); err != nil {
		return nil, fmt.Errorf("reading design.md: %w", err)
	} else {
		run.Design = string(raw)
	}
	if raw, err := os.ReadFile(filepath.Join(runAbs, "plan.md")); err != nil {
		return nil, fmt.Errorf("reading plan.md: %w", err)
	} else {
		run.Plan = string(raw)
		run.PlannedCoverage, run.PlanCoverageGaps = parsePlanCoverage(raw)
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "assessment-results"), "*.json", func(path string, raw []byte) error {
		var rec AssessmentResultRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("assessment-results", filepath.Base(path)))
		run.AssessmentResults = append(run.AssessmentResults, rec)
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

type planFrontmatter struct {
	Coverage *PlannedCoverage `yaml:"coverage"`
}

func parsePlanCoverage(raw []byte) (*PlannedCoverage, []EvaluationRunGap) {
	front, _, ok, err := splitOptionalMarkdownFrontmatter(raw)
	if err != nil {
		return nil, []EvaluationRunGap{invalidPlanCoverageGap(err)}
	}
	if !ok {
		return nil, nil
	}
	var fm planFrontmatter
	if err := yaml.Unmarshal(front, &fm); err != nil {
		return nil, []EvaluationRunGap{invalidPlanCoverageGap(fmt.Errorf("parsing frontmatter: %w", err))}
	}
	if fm.Coverage == nil {
		return nil, nil
	}
	if err := validatePlannedCoverage(*fm.Coverage); err != nil {
		return nil, []EvaluationRunGap{invalidPlanCoverageGap(err)}
	}
	sortPlannedCoverage(fm.Coverage)
	return fm.Coverage, nil
}

func splitOptionalMarkdownFrontmatter(raw []byte) ([]byte, []byte, bool, error) {
	if !bytes.HasPrefix(raw, []byte("---\n")) {
		return nil, raw, false, nil
	}
	rest := raw[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, raw, true, fmt.Errorf("unterminated frontmatter")
	}
	return rest[:end], rest[end+len("\n---\n"):], true, nil
}

func invalidPlanCoverageGap(err error) EvaluationRunGap {
	return EvaluationRunGap{Kind: "invalid-plan-coverage", Ref: "plan.md", Detail: err.Error()}
}

func (r *EvaluationRun) RecordCounts() EvaluationRecordCounts {
	return EvaluationRecordCounts{
		AssessmentResults: len(r.AssessmentResults),
		Analyses:          len(r.Analyses),
		Recommendations:   len(r.Recommendations),
	}
}

// ActiveRecommendationCount returns the number of recommendation records not
// superseded by a later recommendation record in this run.
func (r *EvaluationRun) ActiveRecommendationCount() int {
	superseded := map[string]bool{}
	known := map[string]string{}
	for _, rec := range r.Recommendations {
		for _, ref := range rec.Supersedes {
			if file, ok := resolveKnownRecommendation(known, ref); ok {
				superseded[file] = true
			}
		}
		id := strings.TrimSuffix(filepath.Base(rec.File), ".md")
		known[rec.File] = rec.File
		known[id] = rec.File
		known["recommendations/"+id+".md"] = rec.File
	}
	active := 0
	for _, rec := range r.Recommendations {
		if !superseded[rec.File] {
			active++
		}
	}
	return active
}

func resolveKnownRecommendation(known map[string]string, ref string) (string, bool) {
	ref = strings.TrimSpace(ref)
	if file, ok := known[ref]; ok {
		return file, true
	}
	if strings.HasPrefix(ref, "recommendations/") {
		id := strings.TrimSuffix(filepath.Base(ref), ".md")
		file, ok := known[id]
		return file, ok
	}
	file, ok := known["recommendations/"+ref+".md"]
	return file, ok
}

func (r *EvaluationRun) EvaluationRunStatus() EvaluationRunStatus {
	gaps := r.Renderable()
	status := EvaluationRunStatus{
		SchemaVersion: SchemaVersion,
		Path:          r.Path,
		Reportable:    len(gaps) == 0,
		Counts:        r.RecordCounts(),
		Gaps:          gaps,
	}
	if status.Reportable {
		status.NextActions = []receipt.Action{{
			ID:      "report-build",
			Label:   "Build the evaluation report",
			Command: "qualitymd evaluation report build " + r.Path,
		}}
	} else if hasReviewGap(gaps) {
		status.NextActions = []receipt.Action{{
			ID:      "review-gaps",
			Label:   "Resolve evaluation record gaps",
			Command: "qualitymd evaluation status " + r.Path,
		}}
	} else {
		status.NextActions = []receipt.Action{{
			ID:      "assessment-result-add",
			Label:   "Add the missing evaluation records",
			Command: "qualitymd evaluation assessment-result add " + r.Path,
		}}
	}
	return status
}

func (r *EvaluationRun) Renderable() []EvaluationRunGap {
	var gaps []EvaluationRunGap
	supersededAssessmentResults, supersedingGaps := r.assessmentResultSupersedingState()
	gaps = append(gaps, supersedingGaps...)
	assessmentResultState := r.renderableAssessmentResultState(supersededAssessmentResults)
	gaps = append(gaps, assessmentResultState.Gaps...)
	analysisState := r.renderableAnalysisState()
	gaps = append(gaps, r.plannedCoverageGaps(assessmentResultState.RecordsByIdentity, analysisState.RecordsByIdentity, supersededAssessmentResults)...)
	recommendations := r.renderableRecommendationRefs()
	gaps = append(gaps, r.recommendationReferenceGaps(recommendations)...)
	gaps = append(gaps, r.analysisReferenceGaps(assessmentResultState, analysisState, supersededAssessmentResults)...)
	gaps = append(gaps, r.assessmentResultRecommendationGaps(recommendations)...)
	return gaps
}

type renderableAssessmentResultState struct {
	Known             map[string]bool
	RecordsByIdentity map[string]string
	Gaps              []EvaluationRunGap
}

type renderableAnalysisState struct {
	Known             map[string]AnalysisRecord
	RecordsByIdentity map[string]string
}

func (r *EvaluationRun) renderableAssessmentResultState(superseded map[string]bool) renderableAssessmentResultState {
	state := renderableAssessmentResultState{
		Known:             map[string]bool{},
		RecordsByIdentity: map[string]string{},
	}
	identities := map[string]string{}
	for _, rec := range r.AssessmentResults {
		state.Known[rec.File] = true
		if superseded[rec.File] {
			continue
		}
		key := assessmentResultIdentity(rec)
		state.RecordsByIdentity[key] = rec.File
		if prior, ok := identities[key]; ok {
			state.Gaps = append(state.Gaps, EvaluationRunGap{Kind: "duplicate-assessment-result", Ref: rec.File, Detail: "duplicates " + prior})
			continue
		}
		identities[key] = rec.File
	}
	return state
}

func (r *EvaluationRun) renderableAnalysisState() renderableAnalysisState {
	state := renderableAnalysisState{
		Known:             map[string]AnalysisRecord{},
		RecordsByIdentity: map[string]string{},
	}
	for _, rec := range r.Analyses {
		state.Known[rec.File] = rec
		state.RecordsByIdentity[analysisIdentity(rec)] = rec.File
	}
	return state
}

func (r *EvaluationRun) renderableRecommendationRefs() map[string]bool {
	recs := map[string]bool{}
	for _, rec := range r.Recommendations {
		recs[rec.File] = true
		recs[strings.TrimSuffix(filepath.Base(rec.File), ".md")] = true
	}
	return recs
}

func (r *EvaluationRun) recommendationReferenceGaps(recommendations map[string]bool) []EvaluationRunGap {
	var gaps []EvaluationRunGap
	for _, rec := range r.Recommendations {
		for _, ref := range rec.Supersedes {
			if recommendationRefExists(recommendations, ref) {
				continue
			}
			gaps = append(gaps, EvaluationRunGap{Kind: "missing-superseded-recommendation", Ref: ref, Detail: "referenced by " + rec.File})
		}
	}
	return gaps
}

func (r *EvaluationRun) analysisReferenceGaps(assessmentResults renderableAssessmentResultState, analyses renderableAnalysisState, superseded map[string]bool) []EvaluationRunGap {
	if len(r.Analyses) == 0 {
		return []EvaluationRunGap{{Kind: "missing-analysis", Ref: "analysis/", Detail: "no analysis records are present"}}
	}
	var gaps []EvaluationRunGap
	rootAnalyses := 0
	for _, analysis := range r.Analyses {
		if len(analysis.TargetPath) == 0 {
			rootAnalyses++
		}
		gaps = append(gaps, assessmentResultReferenceGaps(analysis, assessmentResults.Known, superseded)...)
		gaps = append(gaps, childAnalysisReferenceGaps(analysis, analyses.Known)...)
	}
	return append(gaps, rootAnalysisGaps(rootAnalyses)...)
}

func assessmentResultReferenceGaps(analysis AnalysisRecord, assessmentResults, superseded map[string]bool) []EvaluationRunGap {
	var gaps []EvaluationRunGap
	for _, ref := range analysis.AssessmentResultRecords {
		if !assessmentResults[ref] {
			gaps = append(gaps, EvaluationRunGap{Kind: "missing-assessment-result", Ref: ref, Detail: "referenced by " + analysis.File})
			continue
		}
		if superseded[ref] {
			gaps = append(gaps, EvaluationRunGap{Kind: "superseded-assessment-result-reference", Ref: ref, Detail: "referenced by " + analysis.File})
		}
	}
	return gaps
}

func childAnalysisReferenceGaps(analysis AnalysisRecord, analyses map[string]AnalysisRecord) []EvaluationRunGap {
	var gaps []EvaluationRunGap
	for _, ref := range analysis.ChildAnalysisRecords {
		if _, ok := analyses[ref]; !ok {
			gaps = append(gaps, EvaluationRunGap{Kind: "missing-analysis", Ref: ref, Detail: "referenced by " + analysis.File})
		}
	}
	return gaps
}

func rootAnalysisGaps(rootAnalyses int) []EvaluationRunGap {
	switch {
	case rootAnalyses == 0:
		return []EvaluationRunGap{{Kind: "missing-root-analysis", Ref: "analysis/", Detail: "no analysis record has an empty targetPath for the in-scope root"}}
	case rootAnalyses > 1:
		return []EvaluationRunGap{{Kind: "duplicate-root-analysis", Ref: "analysis/", Detail: "multiple analysis records have an empty targetPath"}}
	default:
		return nil
	}
}

func (r *EvaluationRun) assessmentResultRecommendationGaps(recommendations map[string]bool) []EvaluationRunGap {
	var gaps []EvaluationRunGap
	for _, assessmentResult := range r.AssessmentResults {
		for _, ref := range assessmentResult.Recommendations {
			if recommendations[ref] || recommendations["recommendations/"+ref+".md"] {
				continue
			}
			gaps = append(gaps, EvaluationRunGap{Kind: "missing-recommendation", Ref: ref, Detail: "referenced by " + assessmentResult.File})
		}
	}
	return gaps
}

func assessmentResultIdentity(rec AssessmentResultRecord) string {
	return strings.Join(rec.TargetPath, "\x00") + "\x00" + rec.Requirement
}

func assessmentResultID(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".json")
}

type knownAssessmentResult struct {
	File   string
	Record AssessmentResultRecord
}

func (r *EvaluationRun) assessmentResultSupersedingState() (map[string]bool, []EvaluationRunGap) {
	known := map[string]knownAssessmentResult{}
	superseded := map[string]bool{}
	var gaps []EvaluationRunGap
	for _, rec := range r.AssessmentResults {
		for _, ref := range rec.Supersedes {
			ref = strings.TrimSpace(ref)
			prior, ok := resolveKnownAssessmentResult(known, ref)
			if !ok {
				gaps = append(gaps, EvaluationRunGap{Kind: "missing-superseded-assessment-result", Ref: ref, Detail: "referenced by " + rec.File})
				continue
			}
			if assessmentResultIdentity(prior.Record) != assessmentResultIdentity(rec) {
				gaps = append(gaps, EvaluationRunGap{Kind: "invalid-assessment-result-supersedes", Ref: ref, Detail: "referenced by " + rec.File + " with different targetPath or requirement"})
				continue
			}
			superseded[prior.File] = true
			superseded[assessmentResultID(prior.File)] = true
		}
		known[rec.File] = knownAssessmentResult{File: rec.File, Record: rec}
		known[assessmentResultID(rec.File)] = knownAssessmentResult{File: rec.File, Record: rec}
	}
	return superseded, gaps
}

func resolveKnownAssessmentResult(known map[string]knownAssessmentResult, ref string) (knownAssessmentResult, bool) {
	if rec, ok := known[ref]; ok {
		return rec, true
	}
	if strings.HasPrefix(ref, "assessment-results/") {
		rec, ok := known[strings.TrimSuffix(filepath.Base(ref), ".json")]
		return rec, ok
	}
	rec, ok := known["assessment-results/"+ref+".json"]
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

func (r *EvaluationRun) plannedCoverageGaps(assessmentResultRecordsByIdentity, analysisRecordsByIdentity map[string]string, supersededAssessmentResults map[string]bool) []EvaluationRunGap {
	gaps := append([]EvaluationRunGap(nil), r.PlanCoverageGaps...)
	if r.PlannedCoverage == nil {
		return gaps
	}
	plannedAssessmentResults := map[string]PlannedAssessmentResult{}
	for _, assessmentResult := range r.PlannedCoverage.AssessmentResults {
		plannedAssessmentResults[plannedAssessmentResultIdentity(assessmentResult)] = assessmentResult
	}
	plannedAnalyses := map[string]PlannedCoverageAnalysis{}
	for _, analysis := range r.PlannedCoverage.Analyses {
		plannedAnalyses[plannedAnalysisIdentity(analysis)] = analysis
	}
	for _, assessmentResult := range r.PlannedCoverage.AssessmentResults {
		key := plannedAssessmentResultIdentity(assessmentResult)
		if _, ok := assessmentResultRecordsByIdentity[key]; ok {
			continue
		}
		gaps = append(gaps, EvaluationRunGap{
			Kind:   "missing-planned-assessment-result",
			Ref:    "plan.md",
			Detail: describePlannedAssessmentResult(assessmentResult),
		})
	}
	for _, analysis := range r.PlannedCoverage.Analyses {
		key := plannedAnalysisIdentity(analysis)
		if _, ok := analysisRecordsByIdentity[key]; ok {
			continue
		}
		gaps = append(gaps, EvaluationRunGap{
			Kind:   "missing-planned-analysis",
			Ref:    "plan.md",
			Detail: describePlannedAnalysis(analysis),
		})
	}
	for _, assessmentResult := range r.AssessmentResults {
		if supersededAssessmentResults[assessmentResult.File] {
			continue
		}
		planned := PlannedAssessmentResult{TargetPath: assessmentResult.TargetPath, Requirement: assessmentResult.Requirement}
		if _, ok := plannedAssessmentResults[plannedAssessmentResultIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, EvaluationRunGap{
			Kind:   "unexpected-assessment-result",
			Ref:    assessmentResult.File,
			Detail: describePlannedAssessmentResult(planned),
		})
	}
	for _, analysis := range r.Analyses {
		planned := PlannedCoverageAnalysis{TargetPath: analysis.TargetPath}
		if _, ok := plannedAnalyses[plannedAnalysisIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, EvaluationRunGap{
			Kind:   "unexpected-analysis",
			Ref:    analysis.File,
			Detail: describePlannedAnalysis(planned),
		})
	}
	return gaps
}

func hasReviewGap(gaps []EvaluationRunGap) bool {
	for _, gap := range gaps {
		if gap.Kind == "duplicate-assessment-result" || gap.Kind == "duplicate-root-analysis" || gap.Kind == "missing-superseded-recommendation" || gap.Kind == "missing-superseded-assessment-result" || gap.Kind == "invalid-assessment-result-supersedes" || gap.Kind == "invalid-plan-coverage" || gap.Kind == "superseded-assessment-result-reference" || strings.HasPrefix(gap.Kind, "unexpected-") {
			return true
		}
	}
	return false
}
