// Package evaluation owns QUALITY.md evaluation run folders and records.
package evaluation

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qualitymd/quality.md/internal/receipt"
)

// SchemaVersion is the current evaluation record and receipt schema version.
const SchemaVersion = 1

// UsageError marks an evaluation error as invalid user input.
type UsageError struct {
	Err error
}

func (e *UsageError) Error() string { return e.Err.Error() }
func (e *UsageError) Unwrap() error { return e.Err }

func usagef(format string, args ...any) error {
	return &UsageError{Err: fmt.Errorf(format, args...)}
}

// Options configures evaluation run creation.
type Options struct {
	RepoRoot   string
	ResolveDir string
	Narrowing  string
	Model      string
}

// CreateRunReceipt is the JSON contract emitted after creating a run.
type CreateRunReceipt struct {
	Path        string           `json:"path"`
	Number      int              `json:"-"`
	NextActions []receipt.Action `json:"nextActions"`
}

// RecordKind identifies an evaluation record family.
type RecordKind string

const (
	KindAssessmentResult RecordKind = "assessment"
	KindAnalysis         RecordKind = "analysis"
	KindRecommendation   RecordKind = "recommendation"
)

// WriteRecordReceipt is the JSON contract emitted after writing records.
type WriteRecordReceipt struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path,omitempty"`
	Paths         []string         `json:"paths"`
	Kind          RecordKind       `json:"kind"`
	DryRun        bool             `json:"dryRun,omitempty"`
	Created       *bool            `json:"created,omitempty"`
	NextActions   []receipt.Action `json:"nextActions,omitempty"`
}

// WriteOptions configures evaluation record writes.
type WriteOptions struct {
	DryRun bool
}

// PlannedCoverage records the assessment and analysis records expected for a
// run.
type PlannedCoverage struct {
	AssessmentResults []PlannedAssessmentResult `json:"assessmentResults" yaml:"assessmentResults"`
	Analyses          []PlannedCoverageAnalysis `json:"analyses" yaml:"analyses"`
}

// PlannedAssessmentResult identifies one expected assessment result.
type PlannedAssessmentResult struct {
	AreaPath    AreaPath `json:"areaPath" yaml:"areaPath"`
	Requirement string   `json:"requirement" yaml:"requirement"`
}

// PlannedCoverageAnalysis identifies one expected analysis record.
type PlannedCoverageAnalysis struct {
	AreaPath AreaPath `json:"areaPath" yaml:"areaPath"`
}

// AreaPath is a stable path from the root area to a nested area.
type AreaPath []string

// Clone returns a copy of the area path.
func (p AreaPath) Clone() AreaPath {
	if len(p) == 0 {
		return AreaPath{}
	}
	return append(AreaPath(nil), p...)
}

// Elements returns the path elements as a string slice.
func (p AreaPath) Elements() []string {
	return []string(p)
}

// IdentityKey returns a stable string key for equality comparisons.
func (p AreaPath) IdentityKey() string {
	return strings.Join(p, "\x00")
}

// Display returns the human-facing path label.
func (p AreaPath) Display() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

// FactorPath is a stable path from an area's factor set to a nested factor.
type FactorPath []string

// Clone returns a copy of the factor path.
func (p FactorPath) Clone() FactorPath {
	if len(p) == 0 {
		return FactorPath{}
	}
	return append(FactorPath(nil), p...)
}

// Elements returns the path elements as a string slice.
func (p FactorPath) Elements() []string {
	return []string(p)
}

// IdentityKey returns a stable string key for equality comparisons.
func (p FactorPath) IdentityKey() string {
	return strings.Join(p, "\x00")
}

// Display returns the human-facing path label.
func (p FactorPath) Display() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

// Evidence identifies supporting evidence for a finding.
type Evidence struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

// Finding records one observation in an assessment result.
type Finding struct {
	Locator     string               `json:"locator"`
	Observation string               `json:"observation"`
	Category    string               `json:"category"`
	Severity    FindingSeverityLevel `json:"severity"`
	Evidence    []Evidence           `json:"evidence,omitempty"`
	Attributes  map[string]any       `json:"attributes,omitempty"`
}

// FindingSeverityLevel identifies a finding severity.
type FindingSeverityLevel string

const (
	FindingSeverityCritical FindingSeverityLevel = "critical"
	FindingSeverityHigh     FindingSeverityLevel = "high"
	FindingSeverityMedium   FindingSeverityLevel = "medium"
	FindingSeverityLow      FindingSeverityLevel = "low"
	FindingSeverityInfo     FindingSeverityLevel = "info"
)

// FindingSeverity is the report-display form of a severity level.
type FindingSeverity struct {
	Level FindingSeverityLevel `json:"level"`
	Title string               `json:"title"`
}

// Valid reports whether s is a canonical finding severity.
func (s FindingSeverityLevel) Valid() bool {
	switch s {
	case FindingSeverityCritical, FindingSeverityHigh, FindingSeverityMedium, FindingSeverityLow, FindingSeverityInfo:
		return true
	default:
		return false
	}
}

// Title returns the display title for a severity level.
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

// IsRisk reports whether a severity represents a risk finding.
func (s FindingSeverityLevel) IsRisk() bool {
	return s.Valid() && s != FindingSeverityInfo
}

// Display returns the report-display form of a severity level.
func (s FindingSeverityLevel) Display() FindingSeverity {
	return FindingSeverity{Level: s, Title: s.Title()}
}

func findingSeverityLevels() string {
	return "critical, high, medium, low, or info"
}

// AssessmentResultInput is the user-authored payload for an assessment result.
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

// AssessmentResultRecord is a persisted assessment result record.
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

// RatingResult records a rating verdict or not-assessed state and rationale.
type RatingResult struct {
	Kind      RatingResultKind `json:"kind"`
	Level     string           `json:"level,omitempty"`
	Rationale string           `json:"rationale"`
}

// RatingResultKind identifies the shape of a rating result.
type RatingResultKind string

const (
	RatingResultRated       RatingResultKind = "rated"
	RatingResultNotAssessed RatingResultKind = "not-assessed"
)

// Valid reports whether k is a canonical rating result kind.
func (k RatingResultKind) Valid() bool {
	switch k {
	case RatingResultRated, RatingResultNotAssessed:
		return true
	default:
		return false
	}
}

// IsRated reports whether k represents a rated result.
func (k RatingResultKind) IsRated() bool {
	return k == RatingResultRated
}

// IsNotAssessed reports whether k represents a not-assessed result.
func (k RatingResultKind) IsNotAssessed() bool {
	return k == RatingResultNotAssessed
}

// AreaRatingKind identifies how an area's Area-only rating should be
// interpreted.
type AreaRatingKind string

const (
	AreaRatingRated       AreaRatingKind = "rated"
	AreaRatingNotAssessed AreaRatingKind = "not-assessed"
	AreaRatingStructural  AreaRatingKind = "structural"
)

// AreaRatingState is the report-display form of an area's Area-only rating
// state. Rated and not-assessed states carry the underlying rating result;
// structural area-group states carry no nested result.
type AreaRatingState struct {
	Kind         AreaRatingKind `json:"kind"`
	RatingResult *RatingResult  `json:"ratingResult,omitempty"`
	Title        string         `json:"title"`
}

func areaRatingStateFromResult(result *RatingResult) AreaRatingState {
	if result == nil {
		return AreaRatingState{Kind: AreaRatingStructural, Title: "Structural"}
	}
	clone := *result
	if clone.Kind.IsNotAssessed() {
		return AreaRatingState{Kind: AreaRatingNotAssessed, RatingResult: &clone, Title: "Not assessed"}
	}
	return AreaRatingState{Kind: AreaRatingRated, RatingResult: &clone, Title: "Rated"}
}

// RecordLifecycleState identifies whether a record is active or superseded.
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

// Active reports whether s represents an active record.
func (s RecordLifecycleState) Active() bool {
	return s == RecordLifecycleActive
}

// ReportNextStepKind identifies the kind of next action a report recommends.
type ReportNextStepKind string

const (
	ReportNextStepRecommendation ReportNextStepKind = "recommendation"
	ReportNextStepNone           ReportNextStepKind = "none"
)

// Rigor identifies the rigor recorded for an evaluation run.
type Rigor string

const (
	RigorQuick    Rigor = "quick"
	RigorStandard Rigor = "standard"
	RigorDeep     Rigor = "deep"
)

// Valid reports whether r is a canonical evaluation rigor.
func (r Rigor) Valid() bool {
	switch r {
	case RigorQuick, RigorStandard, RigorDeep:
		return true
	default:
		return false
	}
}

// Display returns the human-facing rigor label.
func (r Rigor) Display() string {
	if r == "" {
		return "not recorded"
	}
	return strings.ToUpper(string(r[:1])) + string(r[1:])
}

// Level identifies the level at which an evaluation was performed.
type Level string

const (
	LevelModel Level = "model"
)

// Valid reports whether l is a canonical evaluation level.
func (l Level) Valid() bool {
	switch l {
	case LevelModel:
		return true
	default:
		return false
	}
}

// MissingMetadataKind identifies scope metadata absent from a report.
type MissingMetadataKind string

const (
	MissingMetadataRigor      MissingMetadataKind = "rigor"
	MissingMetadataOutOfScope MissingMetadataKind = "out-of-scope-areas"
	MissingMetadataUnknown    MissingMetadataKind = "unknown"
)

// MissingMetadata describes report metadata that was not recorded.
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

// FactorRatingResult records the rating result for one factor path.
type FactorRatingResult struct {
	FactorPath   FactorPath   `json:"factorPath"`
	RatingResult RatingResult `json:"ratingResult"`
}

// RatingConstraint records the assessment-level reason constraining a rating.
type RatingConstraint struct {
	AssessmentResultRecord string  `json:"assessmentResultRecord,omitempty"`
	Requirement            string  `json:"requirement,omitempty"`
	Level                  *string `json:"level,omitempty"`
}

// AnalysisInput is the user-authored payload for an analysis record.
type AnalysisInput struct {
	AreaPath                AreaPath             `json:"areaPath"`
	LocalRatingResult       *RatingResult        `json:"localRatingResult"`
	FactorRatingResults     []FactorRatingResult `json:"factorRatingResults"`
	AggregateRatingResult   RatingResult         `json:"aggregateRatingResult"`
	AssessmentResultRecords []string             `json:"assessmentResultRecords"`
	ChildAnalysisRecords    []string             `json:"childAnalysisRecords"`
	RatingConstraints       []RatingConstraint   `json:"ratingConstraints,omitempty"`
}

// AnalysisRecord is a persisted analysis record.
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

// RecommendationInput is the user-authored payload for a recommendation record.
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

// RecommendationRecord is a persisted recommendation record.
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
