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

// Run is a loaded evaluation run folder and its decoded records.
type Run struct {
	Path              string
	AbsPath           string
	Design            string
	Plan              string
	PlannedCoverage   *PlannedCoverage
	PlanCoverageGaps  []RunGap
	CompatibilityGaps []RunGap
	AssessmentResults []AssessmentResultRecord
	Analyses          []AnalysisRecord
	Recommendations   []RecommendationRecord
	Model             *model.Spec
	Scale             []model.RatingLevel
	recordCounts      RecordCounts
}

// RecordCounts counts the record files discovered in a run.
type RecordCounts struct {
	AssessmentResults int `json:"assessmentResults"`
	Analyses          int `json:"analyses"`
	Recommendations   int `json:"recommendations"`
}

// RunGap describes a compatibility or reportability gap in a run.
type RunGap struct {
	Kind   RunGapKind `json:"kind"`
	Ref    string     `json:"ref"`
	Detail string     `json:"detail"`
}

// RunGapKind identifies a class of evaluation run gap.
type RunGapKind string

const (
	GapInvalidPlanCoverage                 RunGapKind = "invalid-plan-coverage"
	GapDuplicateAssessmentResult           RunGapKind = "duplicate-assessment-result"
	GapMissingFindingSeverity              RunGapKind = "missing-finding-severity"
	GapInvalidFindingSeverity              RunGapKind = "invalid-finding-severity"
	GapInvalidRatingResult                 RunGapKind = "invalid-rating-result"
	GapMissingSupersededRecommendation     RunGapKind = "missing-superseded-recommendation"
	GapMissingAssessmentResult             RunGapKind = "missing-assessment-result"
	GapSupersededAssessmentResultReference RunGapKind = "superseded-assessment-result-reference"
	GapMissingAnalysis                     RunGapKind = "missing-analysis"
	GapMissingRootAnalysis                 RunGapKind = "missing-root-analysis"
	GapDuplicateRootAnalysis               RunGapKind = "duplicate-root-analysis"
	GapMissingRecommendation               RunGapKind = "missing-recommendation"
	GapMissingSupersededAssessmentResult   RunGapKind = "missing-superseded-assessment-result"
	GapInvalidAssessmentResultSupersedes   RunGapKind = "invalid-assessment-result-supersedes"
	GapMissingPlannedAssessmentResult      RunGapKind = "missing-planned-assessment-result"
	GapMissingPlannedAnalysis              RunGapKind = "missing-planned-analysis"
	GapUnexpectedAssessmentResult          RunGapKind = "unexpected-assessment-result"
	GapUnexpectedAnalysis                  RunGapKind = "unexpected-analysis"
	GapMalformedEvaluationRecord           RunGapKind = "malformed-evaluation-record"
	GapUnreadableEvaluationRecord          RunGapKind = "unreadable-evaluation-record"
	GapMissingRecordSchemaVersion          RunGapKind = "missing-record-schema-version"
	GapUnsupportedRecordSchemaVersion      RunGapKind = "unsupported-record-schema-version"
	GapIncompleteEvaluationRecord          RunGapKind = "incomplete-evaluation-record"
	GapMissingEvaluationData               RunGapKind = "missing-evaluation-data"
	GapMalformedEvaluationData             RunGapKind = "malformed-evaluation-data"
	GapIncompleteEvaluationData            RunGapKind = "incomplete-evaluation-data"
)

// RequiresReview reports whether a gap requires human or agent reconciliation.
func (k RunGapKind) RequiresReview() bool {
	switch k {
	case GapDuplicateAssessmentResult,
		GapDuplicateRootAnalysis,
		GapMissingSupersededRecommendation,
		GapMissingSupersededAssessmentResult,
		GapInvalidAssessmentResultSupersedes,
		GapInvalidPlanCoverage,
		GapMissingFindingSeverity,
		GapInvalidFindingSeverity,
		GapInvalidRatingResult,
		GapSupersededAssessmentResultReference,
		GapUnexpectedAssessmentResult,
		GapUnexpectedAnalysis,
		GapMalformedEvaluationRecord,
		GapUnreadableEvaluationRecord,
		GapMissingRecordSchemaVersion,
		GapUnsupportedRecordSchemaVersion,
		GapIncompleteEvaluationRecord:
		return true
	default:
		return false
	}
}

// RunStatus is the JSON contract emitted by evaluation status.
type RunStatus struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Reportable    bool             `json:"reportable"`
	Counts        RecordCounts     `json:"counts"`
	Gaps          []RunGap         `json:"gaps"`
	NextActions   []receipt.Action `json:"nextActions"`
}

// Load reads a current-contract evaluation run and rejects incompatible records.
func Load(path string) (*Run, error) {
	return load(path, false)
}

// Inspect reads an evaluation run in tolerant mode for history/status views.
func Inspect(path string) (*Run, error) {
	return load(path, true)
}

func load(path string, tolerant bool) (*Run, error) {
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
	run := &Run{
		Path:    displayRunPath(runAbs),
		AbsPath: filepath.ToSlash(runAbs),
		Model:   spec,
		Scale:   spec.RatingScale,
	}
	run.recordCounts = countRecordFiles(runAbs)
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
	loadAssessment := func(path string, raw []byte) error {
		var rec AssessmentResultRecord
		if err := json.Unmarshal(raw, &rec); err != nil {
			return err
		}
		if rec.SchemaVersion != SchemaVersion {
			return fmt.Errorf("schemaVersion = %d, want %d", rec.SchemaVersion, SchemaVersion)
		}
		rec.File = filepath.ToSlash(filepath.Join("assessments", filepath.Base(path)))
		run.AssessmentResults = append(run.AssessmentResults, rec)
		return nil
	}
	loadAnalysis := func(path string, raw []byte) error {
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
	}
	if tolerant {
		run.CompatibilityGaps = append(run.CompatibilityGaps, inspectJSONRecords(filepath.Join(runAbs, "assessments"), "assessments", assessmentRequiredFields(), loadAssessment)...)
		run.CompatibilityGaps = append(run.CompatibilityGaps, inspectJSONRecords(filepath.Join(runAbs, "analysis"), "analysis", analysisRequiredFields(), loadAnalysis)...)
		run.CompatibilityGaps = append(run.CompatibilityGaps, inspectRecommendations(filepath.Join(runAbs, "recommendations"), &run.Recommendations)...)
		return run, nil
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "assessments"), "*.json", loadAssessment); err != nil {
		return nil, err
	}
	if err := loadJSONRecords(filepath.Join(runAbs, "analysis"), "*.json", loadAnalysis); err != nil {
		return nil, err
	}
	if err := loadRecommendations(filepath.Join(runAbs, "recommendations"), &run.Recommendations); err != nil {
		return nil, err
	}
	return run, nil
}

func countRecordFiles(runAbs string) RecordCounts {
	return RecordCounts{
		AssessmentResults: len(globRecordFiles(filepath.Join(runAbs, "assessments"), "*.json")),
		Analyses:          len(globRecordFiles(filepath.Join(runAbs, "analysis"), "*.json")),
		Recommendations:   len(globRecordFiles(filepath.Join(runAbs, "recommendations"), "*.md")),
	}
}

func globRecordFiles(dir, pattern string) []string {
	paths, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil
	}
	slices.Sort(paths)
	return paths
}

func inspectJSONRecords(dir, recordDir string, required []string, decode func(string, []byte) error) []RunGap {
	var gaps []RunGap
	for _, path := range globRecordFiles(dir, "*.json") {
		file := filepath.ToSlash(filepath.Join(recordDir, filepath.Base(path)))
		raw, err := os.ReadFile(path)
		if err != nil {
			gaps = append(gaps, unreadableRecordGap(file, err))
			continue
		}
		if gap := jsonRecordCompatibilityGap(file, raw, required); gap != nil {
			gaps = append(gaps, *gap)
			continue
		}
		if err := decode(path, raw); err != nil {
			gaps = append(gaps, malformedRecordGap(file, err))
		}
	}
	return gaps
}

func jsonRecordCompatibilityGap(file string, raw []byte, required []string) *RunGap {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return gapPtr(malformedRecordGap(file, err))
	}
	versionRaw, ok := fields["schemaVersion"]
	if !ok || isJSONNull(versionRaw) {
		return gapPtr(RunGap{Kind: GapMissingRecordSchemaVersion, Ref: file, Detail: "schemaVersion is required"})
	}
	var version int
	if err := json.Unmarshal(versionRaw, &version); err != nil {
		return gapPtr(RunGap{Kind: GapUnsupportedRecordSchemaVersion, Ref: file, Detail: "schemaVersion must be an integer"})
	}
	if version != SchemaVersion {
		return gapPtr(RunGap{Kind: GapUnsupportedRecordSchemaVersion, Ref: file, Detail: fmt.Sprintf("schemaVersion = %d, want %d", version, SchemaVersion)})
	}
	for _, field := range required {
		raw, ok := fields[field]
		if !ok || (isJSONNull(raw) && field != "localRatingResult") {
			return gapPtr(RunGap{Kind: GapIncompleteEvaluationRecord, Ref: file, Detail: field + " is required"})
		}
	}
	return nil
}

func isJSONNull(raw json.RawMessage) bool {
	return string(bytes.TrimSpace(raw)) == "null"
}

func inspectRecommendations(dir string, out *[]RecommendationRecord) []RunGap {
	var gaps []RunGap
	for _, path := range globRecordFiles(dir, "*.md") {
		file := filepath.ToSlash(filepath.Join("recommendations", filepath.Base(path)))
		raw, err := os.ReadFile(path)
		if err != nil {
			gaps = append(gaps, unreadableRecordGap(file, err))
			continue
		}
		front, body, err := splitMarkdownFrontmatter(raw)
		if err != nil {
			gaps = append(gaps, malformedRecordGap(file, err))
			continue
		}
		if gap := recommendationCompatibilityGap(file, front); gap != nil {
			gaps = append(gaps, *gap)
			continue
		}
		var rec RecommendationRecord
		if err := yaml.Unmarshal(front, &rec); err != nil {
			gaps = append(gaps, malformedRecordGap(file, err))
			continue
		}
		rec.Body = string(body)
		rec.File = file
		*out = append(*out, rec)
	}
	return gaps
}

func recommendationCompatibilityGap(file string, front []byte) *RunGap {
	var fields map[string]any
	if err := yaml.Unmarshal(front, &fields); err != nil {
		return gapPtr(malformedRecordGap(file, err))
	}
	rawVersion, ok := fields["schemaVersion"]
	if !ok || rawVersion == nil {
		return gapPtr(RunGap{Kind: GapMissingRecordSchemaVersion, Ref: file, Detail: "schemaVersion is required"})
	}
	version, ok := rawVersion.(int)
	if !ok {
		return gapPtr(RunGap{Kind: GapUnsupportedRecordSchemaVersion, Ref: file, Detail: "schemaVersion must be an integer"})
	}
	if version != SchemaVersion {
		return gapPtr(RunGap{Kind: GapUnsupportedRecordSchemaVersion, Ref: file, Detail: fmt.Sprintf("schemaVersion = %d, want %d", version, SchemaVersion)})
	}
	for _, field := range recommendationRequiredFields() {
		if value, ok := fields[field]; !ok || value == nil {
			return gapPtr(RunGap{Kind: GapIncompleteEvaluationRecord, Ref: file, Detail: field + " is required"})
		}
	}
	return nil
}

func malformedRecordGap(file string, err error) RunGap {
	return RunGap{Kind: GapMalformedEvaluationRecord, Ref: file, Detail: err.Error()}
}

func unreadableRecordGap(file string, err error) RunGap {
	return RunGap{Kind: GapUnreadableEvaluationRecord, Ref: file, Detail: err.Error()}
}

func gapPtr(gap RunGap) *RunGap {
	return &gap
}

func assessmentRequiredFields() []string {
	return []string{"areaPath", "requirement", "factorPaths", "ratingResult", "criterionSource", "findings", "recommendations"}
}

func analysisRequiredFields() []string {
	return []string{"areaPath", "localRatingResult", "factorRatingResults", "aggregateRatingResult", "assessmentResultRecords", "childAnalysisRecords"}
}

func recommendationRequiredFields() []string {
	return []string{"title", "gap", "evidenceLocators", "assessmentResultRecords", "remediationOptions", "recommendedOption", "doneCriterion"}
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

func parsePlanCoverage(raw []byte) (*PlannedCoverage, []RunGap) {
	front, _, ok, err := splitOptionalMarkdownFrontmatter(raw)
	if err != nil {
		return nil, []RunGap{invalidPlanCoverageGap(err)}
	}
	if !ok {
		return nil, nil
	}
	var fm planFrontmatter
	if err := yaml.Unmarshal(front, &fm); err != nil {
		return nil, []RunGap{invalidPlanCoverageGap(fmt.Errorf("parsing frontmatter: %w", err))}
	}
	if fm.Coverage == nil {
		return nil, nil
	}
	if err := validatePlannedCoverage(*fm.Coverage); err != nil {
		return nil, []RunGap{invalidPlanCoverageGap(err)}
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

func invalidPlanCoverageGap(err error) RunGap {
	return RunGap{Kind: GapInvalidPlanCoverage, Ref: "plan.md", Detail: err.Error()}
}

// RecordCounts returns the counts of record files discovered in the run,
// captured at load time. These are file-glob counts and are distinct from the
// lengths of the decoded result slices after tolerant filtering.
func (r *Run) RecordCounts() RecordCounts {
	return r.recordCounts
}

// ActiveRecommendationCount returns the number of recommendation records not
// superseded by a later recommendation record in this run.
func (r *Run) ActiveRecommendationCount() int {
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

// Status returns the run's report-readiness status. Reportable is true only
// when Renderable returns no gaps. NextActions is always populated with the
// single next step: build the report when reportable, resolve record gaps when
// a gap requires review, or add the missing records otherwise.
func (r *Run) Status() RunStatus {
	if isV2DataRun(r.AbsPath) {
		gaps := v2RenderableGaps(r.AbsPath)
		status := RunStatus{
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
		} else {
			status.NextActions = []receipt.Action{{
				ID:      "evaluation-data-set",
				Label:   "Add missing Evaluation v2 data",
				Command: "qualitymd evaluation data set " + r.Path + " --file <payload.json>",
			}}
		}
		return status
	}
	gaps := r.Renderable()
	status := RunStatus{
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
			ID:      "assessment-add",
			Label:   "Add the missing evaluation records",
			Command: "qualitymd evaluation assessment add " + r.Path,
		}}
	}
	return status
}

// Renderable returns the gaps that block building a report for the run. An
// empty (or nil) slice means the run is reportable; BuildReport and GateReport
// gate on len(Renderable()) == 0.
func (r *Run) Renderable() []RunGap {
	gaps := append([]RunGap(nil), r.CompatibilityGaps...)
	supersededAssessmentResults, supersedingGaps := r.assessmentResultSupersedingState()
	gaps = append(gaps, supersedingGaps...)
	assessmentResultState := r.renderableAssessmentResultState(supersededAssessmentResults)
	gaps = append(gaps, assessmentResultState.Gaps...)
	analysisState := r.renderableAnalysisState()
	gaps = append(gaps, analysisState.Gaps...)
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
	Gaps              []RunGap
}

type renderableAnalysisState struct {
	Known             map[string]AnalysisRecord
	RecordsByIdentity map[string]string
	Gaps              []RunGap
}

func (r *Run) renderableAssessmentResultState(superseded map[string]bool) renderableAssessmentResultState {
	state := renderableAssessmentResultState{
		Known:             map[string]bool{},
		RecordsByIdentity: map[string]string{},
	}
	identities := map[string]string{}
	for _, rec := range r.AssessmentResults {
		state.Known[rec.File] = true
		state.Gaps = append(state.Gaps, assessmentResultRatingGaps(rec, r.ratingLevelSet())...)
		state.Gaps = append(state.Gaps, assessmentResultFindingSeverityGaps(rec)...)
		if superseded[rec.File] {
			continue
		}
		key := assessmentResultIdentity(rec)
		state.RecordsByIdentity[key] = rec.File
		if prior, ok := identities[key]; ok {
			state.Gaps = append(state.Gaps, RunGap{Kind: GapDuplicateAssessmentResult, Ref: rec.File, Detail: "duplicates " + prior})
			continue
		}
		identities[key] = rec.File
	}
	return state
}

func (r *Run) ratingLevelSet() map[string]bool {
	levels := map[string]bool{}
	for _, level := range r.Scale {
		levels[level.Level] = true
	}
	return levels
}

func assessmentResultRatingGaps(rec AssessmentResultRecord, levels map[string]bool) []RunGap {
	if detail := ratingResultGapDetail("ratingResult", rec.RatingResult, levels); detail != "" {
		return []RunGap{{Kind: GapInvalidRatingResult, Ref: rec.File, Detail: detail}}
	}
	return nil
}

func assessmentResultFindingSeverityGaps(rec AssessmentResultRecord) []RunGap {
	var gaps []RunGap
	for i, finding := range rec.Findings {
		switch {
		case strings.TrimSpace(string(finding.Severity)) == "":
			gaps = append(gaps, RunGap{Kind: GapMissingFindingSeverity, Ref: rec.File, Detail: fmt.Sprintf("findings[%d] must include severity", i)})
		case !finding.Severity.Valid():
			gaps = append(gaps, RunGap{Kind: GapInvalidFindingSeverity, Ref: rec.File, Detail: fmt.Sprintf("findings[%d].severity must be one of %s", i, findingSeverityLevels())})
		}
	}
	return gaps
}

func (r *Run) renderableAnalysisState() renderableAnalysisState {
	state := renderableAnalysisState{
		Known:             map[string]AnalysisRecord{},
		RecordsByIdentity: map[string]string{},
	}
	for _, rec := range r.Analyses {
		state.Known[rec.File] = rec
		state.RecordsByIdentity[analysisIdentity(rec)] = rec.File
		state.Gaps = append(state.Gaps, analysisRatingGaps(rec, r.ratingLevelSet())...)
	}
	return state
}

func analysisRatingGaps(rec AnalysisRecord, levels map[string]bool) []RunGap {
	var gaps []RunGap
	if rec.LocalRatingResult != nil {
		if detail := ratingResultGapDetail("localRatingResult", *rec.LocalRatingResult, levels); detail != "" {
			gaps = append(gaps, RunGap{Kind: GapInvalidRatingResult, Ref: rec.File, Detail: detail})
		}
	}
	if detail := ratingResultGapDetail("aggregateRatingResult", rec.AggregateRatingResult, levels); detail != "" {
		gaps = append(gaps, RunGap{Kind: GapInvalidRatingResult, Ref: rec.File, Detail: detail})
	}
	for i, factor := range rec.FactorRatingResults {
		if detail := ratingResultGapDetail(fmt.Sprintf("factorRatingResults[%d].ratingResult", i), factor.RatingResult, levels); detail != "" {
			gaps = append(gaps, RunGap{Kind: GapInvalidRatingResult, Ref: rec.File, Detail: detail})
		}
	}
	return gaps
}

func ratingResultGapDetail(name string, result RatingResult, levels map[string]bool) string {
	if strings.TrimSpace(result.Rationale) == "" {
		return name + ".rationale is required"
	}
	switch result.Kind {
	case RatingResultRated:
		if strings.TrimSpace(result.Level) == "" {
			return name + ".level is required when kind is rated"
		}
		if !levels[result.Level] {
			return fmt.Sprintf("%s.level %q is not defined by the run model", name, result.Level)
		}
	case RatingResultNotAssessed:
		if strings.TrimSpace(result.Level) != "" {
			return name + ".level must be empty when kind is not-assessed"
		}
	default:
		return name + ".kind must be rated or not-assessed"
	}
	return ""
}

func (r *Run) renderableRecommendationRefs() map[string]bool {
	recs := map[string]bool{}
	for _, rec := range r.Recommendations {
		recs[rec.File] = true
		recs[strings.TrimSuffix(filepath.Base(rec.File), ".md")] = true
	}
	return recs
}

func (r *Run) recommendationReferenceGaps(recommendations map[string]bool) []RunGap {
	var gaps []RunGap
	for _, rec := range r.Recommendations {
		for _, ref := range rec.Supersedes {
			if recommendationRefExists(recommendations, ref) {
				continue
			}
			gaps = append(gaps, RunGap{Kind: GapMissingSupersededRecommendation, Ref: ref, Detail: "referenced by " + rec.File})
		}
	}
	return gaps
}

func (r *Run) analysisReferenceGaps(assessmentResults renderableAssessmentResultState, analyses renderableAnalysisState, superseded map[string]bool) []RunGap {
	var gaps []RunGap
	rootAnalyses := 0
	for _, analysis := range r.Analyses {
		if len(analysis.AreaPath) == 0 {
			rootAnalyses++
		}
		gaps = append(gaps, assessmentResultReferenceGaps(analysis, assessmentResults.Known, superseded)...)
		gaps = append(gaps, childAnalysisReferenceGaps(analysis, analyses.Known)...)
	}
	return append(gaps, rootAnalysisGaps(rootAnalyses)...)
}

func assessmentResultReferenceGaps(analysis AnalysisRecord, assessmentResults, superseded map[string]bool) []RunGap {
	var gaps []RunGap
	for _, ref := range analysis.AssessmentResultRecords {
		if !assessmentResults[ref] {
			gaps = append(gaps, RunGap{Kind: GapMissingAssessmentResult, Ref: ref, Detail: "referenced by " + analysis.File})
			continue
		}
		if superseded[ref] {
			gaps = append(gaps, RunGap{Kind: GapSupersededAssessmentResultReference, Ref: ref, Detail: "referenced by " + analysis.File})
		}
	}
	return gaps
}

func childAnalysisReferenceGaps(analysis AnalysisRecord, analyses map[string]AnalysisRecord) []RunGap {
	var gaps []RunGap
	for _, ref := range analysis.ChildAnalysisRecords {
		if _, ok := analyses[ref]; !ok {
			gaps = append(gaps, RunGap{Kind: GapMissingAnalysis, Ref: ref, Detail: "referenced by " + analysis.File})
		}
	}
	return gaps
}

func rootAnalysisGaps(rootAnalyses int) []RunGap {
	switch {
	case rootAnalyses == 0:
		return []RunGap{{Kind: GapMissingRootAnalysis, Ref: "analysis/", Detail: "no analysis record has an empty areaPath for the in-scope root"}}
	case rootAnalyses > 1:
		return []RunGap{{Kind: GapDuplicateRootAnalysis, Ref: "analysis/", Detail: "multiple analysis records have an empty areaPath"}}
	default:
		return nil
	}
}

func (r *Run) assessmentResultRecommendationGaps(recommendations map[string]bool) []RunGap {
	var gaps []RunGap
	for _, assessmentResult := range r.AssessmentResults {
		for _, ref := range assessmentResult.Recommendations {
			if recommendations[ref] || recommendations["recommendations/"+ref+".md"] {
				continue
			}
			gaps = append(gaps, RunGap{Kind: GapMissingRecommendation, Ref: ref, Detail: "referenced by " + assessmentResult.File})
		}
	}
	return gaps
}

func assessmentResultIdentity(rec AssessmentResultRecord) string {
	return strings.Join(rec.AreaPath, "\x00") + "\x00" + rec.Requirement
}

func assessmentResultID(file string) string {
	return strings.TrimSuffix(filepath.Base(file), ".json")
}

type knownAssessmentResult struct {
	File   string
	Record AssessmentResultRecord
}

func (r *Run) assessmentResultSupersedingState() (map[string]bool, []RunGap) {
	known := map[string]knownAssessmentResult{}
	superseded := map[string]bool{}
	var gaps []RunGap
	for _, rec := range r.AssessmentResults {
		for _, ref := range rec.Supersedes {
			ref = strings.TrimSpace(ref)
			prior, ok := resolveKnownAssessmentResult(known, ref)
			if !ok {
				gaps = append(gaps, RunGap{Kind: GapMissingSupersededAssessmentResult, Ref: ref, Detail: "referenced by " + rec.File})
				continue
			}
			if assessmentResultIdentity(prior.Record) != assessmentResultIdentity(rec) {
				gaps = append(gaps, RunGap{Kind: GapInvalidAssessmentResultSupersedes, Ref: ref, Detail: "referenced by " + rec.File + " with different areaPath or requirement"})
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
	if strings.HasPrefix(ref, "assessments/") {
		rec, ok := known[strings.TrimSuffix(filepath.Base(ref), ".json")]
		return rec, ok
	}
	rec, ok := known["assessments/"+ref+".json"]
	return rec, ok
}

func analysisIdentity(rec AnalysisRecord) string {
	return strings.Join(rec.AreaPath, "\x00")
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

func (r *Run) plannedCoverageGaps(assessmentResultRecordsByIdentity, analysisRecordsByIdentity map[string]string, supersededAssessmentResults map[string]bool) []RunGap {
	gaps := append([]RunGap(nil), r.PlanCoverageGaps...)
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
		gaps = append(gaps, RunGap{
			Kind:   GapMissingPlannedAssessmentResult,
			Ref:    "plan.md",
			Detail: describePlannedAssessmentResult(assessmentResult),
		})
	}
	for _, analysis := range r.PlannedCoverage.Analyses {
		key := plannedAnalysisIdentity(analysis)
		if _, ok := analysisRecordsByIdentity[key]; ok {
			continue
		}
		gaps = append(gaps, RunGap{
			Kind:   GapMissingPlannedAnalysis,
			Ref:    "plan.md",
			Detail: describePlannedAnalysis(analysis),
		})
	}
	for _, assessmentResult := range r.AssessmentResults {
		if supersededAssessmentResults[assessmentResult.File] {
			continue
		}
		planned := PlannedAssessmentResult{AreaPath: assessmentResult.AreaPath, Requirement: assessmentResult.Requirement}
		if _, ok := plannedAssessmentResults[plannedAssessmentResultIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, RunGap{
			Kind:   GapUnexpectedAssessmentResult,
			Ref:    assessmentResult.File,
			Detail: describePlannedAssessmentResult(planned),
		})
	}
	for _, analysis := range r.Analyses {
		planned := PlannedCoverageAnalysis{AreaPath: analysis.AreaPath}
		if _, ok := plannedAnalyses[plannedAnalysisIdentity(planned)]; ok {
			continue
		}
		gaps = append(gaps, RunGap{
			Kind:   GapUnexpectedAnalysis,
			Ref:    analysis.File,
			Detail: describePlannedAnalysis(planned),
		})
	}
	return gaps
}

func hasReviewGap(gaps []RunGap) bool {
	for _, gap := range gaps {
		if gap.Kind.RequiresReview() {
			return true
		}
	}
	return false
}
