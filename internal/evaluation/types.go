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
	Altitude      string
	Narrowing     string
	Subject       string
}

type CreateRunResult struct {
	Path        string           `json:"path"`
	Number      int              `json:"-"`
	Altitude    string           `json:"altitude"`
	NextActions []receipt.Action `json:"nextActions"`
}

type WriteKind string

const (
	KindAssessment     WriteKind = "assessment"
	KindAnalysis       WriteKind = "analysis"
	KindRecommendation WriteKind = "recommendation"
)

type WriteResult struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Kind          WriteKind        `json:"kind"`
	Created       *bool            `json:"created,omitempty"`
	NextActions   []receipt.Action `json:"nextActions,omitempty"`
}

type PlannedCoverageResult struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	NextActions   []receipt.Action `json:"nextActions,omitempty"`
}

type PlannedCoverage struct {
	SchemaVersion int                         `json:"schemaVersion"`
	Assessments   []PlannedCoverageAssessment `json:"assessments"`
	Analyses      []PlannedCoverageAnalysis   `json:"analyses"`
}

type PlannedCoverageAssessment struct {
	TargetPath  []string `json:"targetPath"`
	Requirement string   `json:"requirement"`
}

type PlannedCoverageAnalysis struct {
	TargetPath []string `json:"targetPath"`
}

type Evidence struct {
	Kind string `json:"kind"`
	Ref  string `json:"ref"`
}

type Finding struct {
	Locator     string         `json:"locator"`
	Observation string         `json:"observation"`
	Category    string         `json:"category"`
	Severity    string         `json:"severity,omitempty"`
	Evidence    []Evidence     `json:"evidence,omitempty"`
	Attributes  map[string]any `json:"attributes,omitempty"`
}

type AssessmentPayload struct {
	Target          string    `json:"target"`
	TargetPath      []string  `json:"targetPath"`
	Requirement     string    `json:"requirement"`
	Factors         []string  `json:"factors"`
	Rating          *string   `json:"rating"`
	NotAssessed     bool      `json:"notAssessed"`
	CriterionSource string    `json:"criterionSource"`
	Findings        []Finding `json:"findings"`
	Rationale       string    `json:"rationale"`
	Recommendations []string  `json:"recommendations"`
	Supersedes      []string  `json:"supersedes,omitempty"`
}

type AssessmentRecord struct {
	SchemaVersion   int       `json:"schemaVersion"`
	Target          string    `json:"target"`
	TargetPath      []string  `json:"targetPath"`
	Requirement     string    `json:"requirement"`
	Factors         []string  `json:"factors"`
	Rating          *string   `json:"rating"`
	NotAssessed     bool      `json:"notAssessed"`
	CriterionSource string    `json:"criterionSource"`
	Findings        []Finding `json:"findings"`
	Rationale       string    `json:"rationale"`
	Recommendations []string  `json:"recommendations"`
	Supersedes      []string  `json:"supersedes,omitempty"`
	File            string    `json:"-"`
}

type RatingResult struct {
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Rationale   string  `json:"rationale"`
}

type FactorRating struct {
	Factor      string  `json:"factor"`
	Rating      *string `json:"rating"`
	NotAssessed bool    `json:"notAssessed"`
	Rationale   string  `json:"rationale"`
}

type BindingConstraint struct {
	Record      string  `json:"record,omitempty"`
	Requirement string  `json:"requirement,omitempty"`
	Rating      *string `json:"rating,omitempty"`
}

type AnalysisPayload struct {
	Target               string              `json:"target"`
	TargetPath           []string            `json:"targetPath"`
	LocalRating          *RatingResult       `json:"localRating"`
	FactorRatings        []FactorRating      `json:"factorRatings"`
	AggregateRating      RatingResult        `json:"aggregateRating"`
	AssessmentRecords    []string            `json:"assessmentRecords"`
	ChildAnalysisRecords []string            `json:"childAnalysisRecords"`
	BindingConstraints   []BindingConstraint `json:"bindingConstraints,omitempty"`
}

type AnalysisRecord struct {
	SchemaVersion        int                 `json:"schemaVersion"`
	Target               string              `json:"target"`
	TargetPath           []string            `json:"targetPath"`
	LocalRating          *RatingResult       `json:"localRating"`
	FactorRatings        []FactorRating      `json:"factorRatings"`
	AggregateRating      RatingResult        `json:"aggregateRating"`
	AssessmentRecords    []string            `json:"assessmentRecords"`
	ChildAnalysisRecords []string            `json:"childAnalysisRecords"`
	BindingConstraints   []BindingConstraint `json:"bindingConstraints,omitempty"`
	File                 string              `json:"-"`
}

type RecommendationPayload struct {
	Title              string   `json:"title"`
	Gap                string   `json:"gap"`
	EvidenceLocators   []string `json:"evidenceLocators"`
	AssessmentRecords  []string `json:"assessmentRecords"`
	RemediationOptions []string `json:"remediationOptions"`
	RecommendedOption  string   `json:"recommendedOption"`
	DoneCriterion      string   `json:"doneCriterion"`
	Supersedes         []string `json:"supersedes,omitempty"`
}

type RecommendationRecord struct {
	SchemaVersion      int      `json:"schemaVersion" yaml:"schemaVersion"`
	Title              string   `json:"title" yaml:"title"`
	Gap                string   `json:"gap" yaml:"gap"`
	EvidenceLocators   []string `json:"evidenceLocators" yaml:"evidenceLocators"`
	AssessmentRecords  []string `json:"assessmentRecords" yaml:"assessmentRecords"`
	RemediationOptions []string `json:"remediationOptions" yaml:"remediationOptions"`
	RecommendedOption  string   `json:"recommendedOption" yaml:"recommendedOption"`
	DoneCriterion      string   `json:"doneCriterion" yaml:"doneCriterion"`
	Supersedes         []string `json:"supersedes,omitempty" yaml:"supersedes,omitempty"`
	Body               string   `json:"-" yaml:"-"`
	File               string   `json:"-" yaml:"-"`
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
