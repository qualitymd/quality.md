// Package evaluation owns QUALITY.md Evaluation run folders and data.
package evaluation

import (
	"fmt"
	"strings"

	"github.com/qualitymd/quality.md/internal/receipt"
)

// SchemaVersion is the current evaluation data and receipt schema version.
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

// AreaPath is a stable path from the root area to a nested area.
type AreaPath []string

// Clone returns a copy of the Area path.
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

// Display returns the human-facing Area path label.
func (p AreaPath) Display() string {
	if len(p) == 0 {
		return "/"
	}
	return strings.Join(p, "/")
}

func (p AreaPath) referencePath() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

// Reference returns the canonical typed model reference for an Area path.
func (p AreaPath) Reference() string {
	return "area:" + p.referencePath()
}

// UnqualifiedReference returns the fixed-type Area reference.
func (p AreaPath) UnqualifiedReference() string {
	return p.referencePath()
}

// FactorPath is a stable path from an Area's Factor set to a nested Factor.
type FactorPath []string

// Clone returns a copy of the Factor path.
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

// Display returns the human-facing Factor path label.
func (p FactorPath) Display() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

func (p FactorPath) referencePath() string {
	if len(p) == 0 {
		return "root"
	}
	return strings.Join(p, "/")
}

// FactorReference returns the canonical typed model reference for a Factor path
// declared by areaPath.
func FactorReference(areaPath AreaPath, factorPath FactorPath) string {
	return "factor:" + areaPath.referencePath() + "::" + factorPath.referencePath()
}

// UnqualifiedFactorReference returns the fixed-type Factor reference.
func UnqualifiedFactorReference(areaPath AreaPath, factorPath FactorPath) string {
	return areaPath.referencePath() + "::" + factorPath.referencePath()
}

// RequirementReference returns the canonical typed model reference for a
// Requirement name declared by areaPath.
func RequirementReference(areaPath AreaPath, requirementName string) string {
	return "requirement:" + areaPath.referencePath() + "::" + requirementName
}

// UnqualifiedRequirementReference returns the fixed-type Requirement reference.
func UnqualifiedRequirementReference(areaPath AreaPath, requirementName string) string {
	return areaPath.referencePath() + "::" + requirementName
}

// RatingReference returns the canonical typed model reference for a Rating
// Level ID.
func RatingReference(level string) string {
	return "rating:" + level
}

// RatingDisplay returns the human-facing label for a Rating Level ID.
func RatingDisplay(level string) string {
	return level
}

// UnqualifiedRatingReference returns the fixed-type Rating Level reference.
func UnqualifiedRatingReference(level string) string {
	return level
}

// RatingResult records an Evaluation report receipt rating verdict or
// not-assessed state.
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
