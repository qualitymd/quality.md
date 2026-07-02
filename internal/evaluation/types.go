// Package evaluation owns QUALITY.md evaluation run folders and data.
package evaluation

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"
	"time"

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
	Area       string
	Factors    []string
	Model      string
}

// RunScope records the faithful user-requested evaluation scope.
type RunScope struct {
	AreaID       string   `json:"areaId,omitempty"`
	FactorFilter []string `json:"factorFilter,omitempty"`
}

// PlannedRunScope records the normalized evaluation scope used by deterministic
// run consumers.
type PlannedRunScope struct {
	AreaID       string   `json:"areaId"`
	FactorFilter []string `json:"factorFilter"`
}

// RunMetadata records repository-local run numbering and folder identity.
type RunMetadata struct {
	Number int    `json:"number"`
	Label  string `json:"label"`
}

// EvaluationManifest is the CLI-owned Evaluation metadata written when a run is created.
type EvaluationManifest struct {
	SchemaVersion  int             `json:"schemaVersion"`
	Kind           DataKind        `json:"kind"`
	EvaluationID   string          `json:"evaluationId"`
	CreatedAt      string          `json:"createdAt"`
	Model          string          `json:"model"`
	RequestedScope RunScope        `json:"requestedScope"`
	PlannedScope   PlannedRunScope `json:"plannedScope"`
	Run            RunMetadata     `json:"run"`
}

const (
	runIDTailLength = 12
	runIDAlphabet   = "0123456789abcdefghjkmnpqrstvwxyz"
)

func newEvaluationIdentity() (string, string, error) {
	createdAt := time.Now().UTC()
	tail, err := randomEvaluationIDTail(runIDTailLength)
	if err != nil {
		return "", "", err
	}
	id := createdAt.Format("20060102T150405Z") + "-" + tail
	return id, createdAt.Format(time.RFC3339), nil
}

func randomEvaluationIDTail(length int) (string, error) {
	raw := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return "", fmt.Errorf("generating evaluation ID: %w", err)
	}
	var b strings.Builder
	b.Grow(length)
	for _, value := range raw {
		b.WriteByte(runIDAlphabet[int(value)&31])
	}
	return b.String(), nil
}

// CreateRunReceipt is the JSON contract emitted after creating a run.
type CreateRunReceipt struct {
	Path        string           `json:"path"`
	Number      int              `json:"-"`
	NextActions []receipt.Action `json:"nextActions"`
}

// RatingReference returns the canonical typed model reference for a rating
// level ID.
func RatingReference(level string) string {
	return "rating:" + level
}

// RatingDisplay returns the human-facing label for a rating level ID.
func RatingDisplay(level string) string {
	return level
}

// UnqualifiedRatingReference returns the fixed-type rating level reference.
func UnqualifiedRatingReference(level string) string {
	return level
}

// RatingResult records an evaluation report receipt rating verdict or
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
