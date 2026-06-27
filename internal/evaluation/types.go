// Package evaluation owns QUALITY.md Evaluation run folders and data.
package evaluation

import (
	"fmt"

	"github.com/qualitymd/quality.md/internal/receipt"
)

// SchemaVersion is the current evaluation data and receipt schema version.
const SchemaVersion = 3

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
