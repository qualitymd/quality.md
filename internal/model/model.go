// Package model defines the typed QUALITY.md frontmatter model.
package model

import (
	"fmt"

	"github.com/qualitymd/quality.md/internal/document"
)

// Spec is the typed QUALITY.md frontmatter model used by callers that need a
// valid model after lint has accepted the document.
type Spec struct {
	Title        string                 `yaml:"title,omitempty" json:"title,omitempty"`
	RatingScale  []RatingLevel          `yaml:"ratingScale" json:"ratingScale"`
	Factors      map[string]Factor      `yaml:"factors,omitempty" json:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty" json:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty" json:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty" json:"source,omitempty"`

	// Path is the source file; not part of the YAML.
	Path string `yaml:"-" json:"-"`
}

// RatingLevel is one level in a model's rating scale. Description states what
// the level means across the whole model and is never overridden; Criterion is
// the default rule for rating a requirement's findings and MAY be overridden per
// requirement via Requirement.Ratings.
type RatingLevel struct {
	Level       string `yaml:"level" json:"level"`
	Title       string `yaml:"title,omitempty" json:"title,omitempty"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Criterion   string `yaml:"criterion" json:"criterion"`
}

// Target is a recursive target node in the quality model.
type Target struct {
	Factors      map[string]Factor      `yaml:"factors,omitempty" json:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty" json:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty" json:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty" json:"source,omitempty"`
}

// Factor is a quality lens scoped to the target where it is declared.
type Factor struct {
	Description  string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Factors      map[string]Factor      `yaml:"factors,omitempty" json:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty" json:"requirements,omitempty"`
}

// Requirement is one assessable expectation.
type Requirement struct {
	Assessment string            `yaml:"assessment" json:"assessment"`
	Factors    []string          `yaml:"factors,omitempty" json:"factors,omitempty"`
	Ratings    map[string]string `yaml:"ratings,omitempty" json:"ratings,omitempty"`
}

// Decode unmarshals a parsed document into the typed model. Callers should only
// use this after lint has accepted the document.
func Decode(doc *document.Document) (*Spec, error) {
	var out Spec
	if err := doc.Frontmatter.Decode(&out); err != nil {
		return nil, fmt.Errorf("%s: decoding model: %w", doc.Path, err)
	}
	out.Path = doc.Path
	return &out, nil
}
