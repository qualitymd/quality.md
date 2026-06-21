// Package evaluation owns QUALITY.md evaluation run folders and records.
package evaluation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qualitymd/quality.md/internal/receipt"
)

const SchemaVersion = 1

type UsageError struct {
	Err error
}

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

func usagef(format string, args ...any) error {
	return &UsageError{Err: fmt.Errorf(format, args...)}
}

type Options struct {
	RepoRoot      string
	EvaluationDir string
	Narrowing     string
	Model         string
}

type CreateRunReceipt struct {
	Path        string           `json:"path"`
	Number      int              `json:"-"`
	NextActions []receipt.Action `json:"nextActions"`
}

type EvaluationRecordKind string

const (
	KindAssessmentResult EvaluationRecordKind = "assessment"
	KindAnalysis         EvaluationRecordKind = "analysis"
	KindRecommendation   EvaluationRecordKind = "recommendation"
)

type WriteRecordReceipt struct {
	SchemaVersion int                  `json:"schemaVersion"`
	Path          string               `json:"path,omitempty"`
	Paths         []string             `json:"paths"`
	Kind          EvaluationRecordKind `json:"kind"`
	Created       *bool                `json:"created,omitempty"`
	NextActions   []receipt.Action     `json:"nextActions,omitempty"`
}

type PlannedCoverage struct {
	AssessmentResults []PlannedAssessmentResult `json:"assessmentResults" yaml:"assessmentResults"`
	Analyses          []PlannedCoverageAnalysis `json:"analyses" yaml:"analyses"`
}

type PlannedAssessmentResult struct {
	AreaPath    AreaPath `json:"areaPath" yaml:"areaPath"`
	Requirement string   `json:"requirement" yaml:"requirement"`
}

type PlannedCoverageAnalysis struct {
	AreaPath AreaPath `json:"areaPath" yaml:"areaPath"`
}

type AreaPath []string

func (p AreaPath) Clone() AreaPath {
	if len(p) == 0 {
		return AreaPath{}
	}
	return append(AreaPath(nil), p...)
}

func (p AreaPath) Elements() []string {
	return []string(p)
}

func (p AreaPath) IdentityKey() string {
	return strings.Join(p, "\x00")
}

func (p AreaPath) Display() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

type FactorPath []string

func (p FactorPath) Clone() FactorPath {
	if len(p) == 0 {
		return FactorPath{}
	}
	return append(FactorPath(nil), p...)
}

func (p FactorPath) Elements() []string {
	return []string(p)
}

func (p FactorPath) Display() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

type Evidence struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

type Finding struct {
	Locator     string               `json:"locator"`
	Observation string               `json:"observation"`
	Category    string               `json:"category"`
	Severity    FindingSeverityLevel `json:"severity"`
	Evidence    []Evidence           `json:"evidence,omitempty"`
	Attributes  map[string]any       `json:"attributes,omitempty"`
}

type FindingSeverityLevel string

const (
	FindingSeverityCritical FindingSeverityLevel = "critical"
	FindingSeverityHigh     FindingSeverityLevel = "high"
	FindingSeverityMedium   FindingSeverityLevel = "medium"
	FindingSeverityLow      FindingSeverityLevel = "low"
	FindingSeverityInfo     FindingSeverityLevel = "info"
)

type FindingSeverity struct {
	Level FindingSeverityLevel `json:"level"`
	Title string               `json:"title"`
}

func (s FindingSeverityLevel) Valid() bool {
	switch s {
	case FindingSeverityCritical, FindingSeverityHigh, FindingSeverityMedium, FindingSeverityLow, FindingSeverityInfo:
		return true
	default:
		return false
	}
}

func (s FindingSeverityLevel) Title() string {
	switch s {
	case FindingSeverityCritical:
		return "Critical"
	case FindingSeverityHigh:
		return "High"
	case FindingSeverityMedium:
		return "Medium"
	case FindingSeverityLow:
		return "Low"
	case FindingSeverityInfo:
		return "Info"
	default:
		return string(s)
	}
}

func (s FindingSeverityLevel) IsRisk() bool {
	return s.Valid() && s != FindingSeverityInfo
}

func (s FindingSeverityLevel) Display() FindingSeverity {
	return FindingSeverity{Level: s, Title: s.Title()}
}

func findingSeverityLevels() string {
	return "critical, high, medium, low, or info"
}

type AssessmentResultInput struct {
	AreaPath        AreaPath     `json:"areaPath"`
	Requirement     string       `json:"requirement"`
	FactorPaths     []FactorPath `json:"factorPaths"`
	RatingResult    RatingResult `json:"ratingResult"`
	CriterionSource string       `json:"criterionSource"`
	Findings        []Finding    `json:"findings"`
	Recommendations []string     `json:"recommendations"`
	Supersedes      []string     `json:"supersedes,omitempty"`
}

type AssessmentResultRecord struct {
	SchemaVersion   int          `json:"schemaVersion"`
	AreaPath        AreaPath     `json:"areaPath"`
	Requirement     string       `json:"requirement"`
	FactorPaths     []FactorPath `json:"factorPaths"`
	RatingResult    RatingResult `json:"ratingResult"`
	CriterionSource string       `json:"criterionSource"`
	Findings        []Finding    `json:"findings"`
	Recommendations []string     `json:"recommendations"`
	Supersedes      []string     `json:"supersedes,omitempty"`
	File            string       `json:"-"`
}

type RatingResult struct {
	Kind      RatingResultKind `json:"kind"`
	Level     string           `json:"level,omitempty"`
	Rationale string           `json:"rationale"`
}

type RatingResultKind string

const (
	RatingResultRated       RatingResultKind = "rated"
	RatingResultNotAssessed RatingResultKind = "not-assessed"
)

func (k RatingResultKind) Valid() bool {
	switch k {
	case RatingResultRated, RatingResultNotAssessed:
		return true
	default:
		return false
	}
}

func (k RatingResultKind) IsRated() bool {
	return k == RatingResultRated
}

func (k RatingResultKind) IsNotAssessed() bool {
	return k == RatingResultNotAssessed
}

type LocalRatingKind string

const (
	LocalRatingRated       LocalRatingKind = "rated"
	LocalRatingNotAssessed LocalRatingKind = "not-assessed"
	LocalRatingStructural  LocalRatingKind = "structural"
)

type LocalRatingState struct {
	Kind         LocalRatingKind `json:"kind"`
	RatingResult *RatingResult   `json:"ratingResult,omitempty"`
	Title        string          `json:"title"`
}

func localRatingStateFromResult(result *RatingResult) LocalRatingState {
	if result == nil {
		return LocalRatingState{Kind: LocalRatingStructural, Title: "Structural"}
	}
	clone := *result
	if clone.Kind.IsNotAssessed() {
		return LocalRatingState{Kind: LocalRatingNotAssessed, RatingResult: &clone, Title: "Not assessed"}
	}
	return LocalRatingState{Kind: LocalRatingRated, RatingResult: &clone, Title: "Rated"}
}

type RecordLifecycleState string

const (
	RecordLifecycleActive     RecordLifecycleState = "active"
	RecordLifecycleSuperseded RecordLifecycleState = "superseded"
)

func lifecycleState(active bool) RecordLifecycleState {
	if active {
		return RecordLifecycleActive
	}
	return RecordLifecycleSuperseded
}

func (s RecordLifecycleState) Active() bool {
	return s == RecordLifecycleActive
}

type ReportNextStepKind string

const (
	ReportNextStepRecommendation ReportNextStepKind = "recommendation"
	ReportNextStepNone           ReportNextStepKind = "none"
)

type EvaluationRigor string

const (
	EvaluationRigorQuick    EvaluationRigor = "quick"
	EvaluationRigorStandard EvaluationRigor = "standard"
	EvaluationRigorDeep     EvaluationRigor = "deep"
)

func (r EvaluationRigor) Valid() bool {
	switch r {
	case EvaluationRigorQuick, EvaluationRigorStandard, EvaluationRigorDeep:
		return true
	default:
		return false
	}
}

func (r EvaluationRigor) Display() string {
	if r == "" {
		return "not recorded"
	}
	return strings.ToUpper(string(r[:1])) + string(r[1:])
}

type EvaluationLevel string

const (
	EvaluationLevelModel EvaluationLevel = "model"
)

func (l EvaluationLevel) Valid() bool {
	switch l {
	case EvaluationLevelModel:
		return true
	default:
		return false
	}
}

type MissingMetadataKind string

const (
	MissingMetadataRigor      MissingMetadataKind = "rigor"
	MissingMetadataOutOfScope MissingMetadataKind = "out-of-scope-areas"
	MissingMetadataUnknown    MissingMetadataKind = "unknown"
)

type MissingMetadata struct {
	Field MissingMetadataKind `json:"field"`
	Title string              `json:"title"`
}

func missingMetadata(kind MissingMetadataKind) MissingMetadata {
	switch kind {
	case MissingMetadataRigor:
		return MissingMetadata{Field: kind, Title: "Rigor"}
	case MissingMetadataOutOfScope:
		return MissingMetadata{Field: kind, Title: "Out-of-scope areas"}
	default:
		return MissingMetadata{Field: MissingMetadataUnknown, Title: "Unknown metadata"}
	}
}

type FactorRatingResult struct {
	FactorPath   FactorPath   `json:"factorPath"`
	RatingResult RatingResult `json:"ratingResult"`
}

type RatingConstraint struct {
	AssessmentResultRecord string  `json:"assessmentResultRecord,omitempty"`
	Requirement            string  `json:"requirement,omitempty"`
	Level                  *string `json:"level,omitempty"`
}

type AnalysisInput struct {
	AreaPath                AreaPath             `json:"areaPath"`
	LocalRatingResult       *RatingResult        `json:"localRatingResult"`
	FactorRatingResults     []FactorRatingResult `json:"factorRatingResults"`
	AggregateRatingResult   RatingResult         `json:"aggregateRatingResult"`
	AssessmentResultRecords []string             `json:"assessmentResultRecords"`
	ChildAnalysisRecords    []string             `json:"childAnalysisRecords"`
	RatingConstraints       []RatingConstraint   `json:"ratingConstraints,omitempty"`
}

type AnalysisRecord struct {
	SchemaVersion           int                  `json:"schemaVersion"`
	AreaPath                AreaPath             `json:"areaPath"`
	LocalRatingResult       *RatingResult        `json:"localRatingResult"`
	FactorRatingResults     []FactorRatingResult `json:"factorRatingResults"`
	AggregateRatingResult   RatingResult         `json:"aggregateRatingResult"`
	AssessmentResultRecords []string             `json:"assessmentResultRecords"`
	ChildAnalysisRecords    []string             `json:"childAnalysisRecords"`
	RatingConstraints       []RatingConstraint   `json:"ratingConstraints,omitempty"`
	File                    string               `json:"-"`
}

type RecommendationInput struct {
	Title                   string   `json:"title"`
	Gap                     string   `json:"gap"`
	EvidenceLocators        []string `json:"evidenceLocators"`
	AssessmentResultRecords []string `json:"assessmentResultRecords"`
	RemediationOptions      []string `json:"remediationOptions"`
	RecommendedOption       string   `json:"recommendedOption"`
	DoneCriterion           string   `json:"doneCriterion"`
	Supersedes              []string `json:"supersedes,omitempty"`
}

type RecommendationRecord struct {
	SchemaVersion           int      `json:"schemaVersion" yaml:"schemaVersion"`
	Title                   string   `json:"title" yaml:"title"`
	Gap                     string   `json:"gap" yaml:"gap"`
	EvidenceLocators        []string `json:"evidenceLocators" yaml:"evidenceLocators"`
	AssessmentResultRecords []string `json:"assessmentResultRecords" yaml:"assessmentResultRecords"`
	RemediationOptions      []string `json:"remediationOptions" yaml:"remediationOptions"`
	RecommendedOption       string   `json:"recommendedOption" yaml:"recommendedOption"`
	DoneCriterion           string   `json:"doneCriterion" yaml:"doneCriterion"`
	Supersedes              []string `json:"supersedes,omitempty" yaml:"supersedes,omitempty"`
	Body                    string   `json:"-" yaml:"-"`
	File                    string   `json:"-" yaml:"-"`
}

func marshalJSON(v any) ([]byte, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func requiredString(name, value string) error {
	if strings.TrimSpace(value) == "" {
		return usagef("%s is required", name)
	}
	return nil
}
