// Package spec loads and models QUALITY.md frontmatter.
package spec

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Spec is a parsed QUALITY.md document.
type Spec struct {
	Title        string                 `yaml:"title,omitempty"`
	RatingScale  []RatingLevel          `yaml:"ratingScale"`
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty"`

	// Path is the source file; not part of the YAML.
	Path string `yaml:"-"`
}

// RatingLevel is one level in a model's rating scale. Description states what
// the level means across the whole model and is never overridden; Criterion is
// the default rule for rating a requirement's findings and MAY be overridden per
// requirement via Requirement.Ratings.
type RatingLevel struct {
	Level       string `yaml:"level"`
	Title       string `yaml:"title,omitempty"`
	Description string `yaml:"description,omitempty"`
	Criterion   string `yaml:"criterion"`
}

// Target is a recursive target node in the quality model.
type Target struct {
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
	Targets      map[string]Target      `yaml:"targets,omitempty"`
	Source       string                 `yaml:"source,omitempty"`
}

// Factor is a quality lens scoped to the target where it is declared.
type Factor struct {
	Description  string                 `yaml:"description,omitempty"`
	Factors      map[string]Factor      `yaml:"factors,omitempty"`
	Requirements map[string]Requirement `yaml:"requirements,omitempty"`
}

// Requirement is one assessable expectation.
type Requirement struct {
	Assessment string            `yaml:"assessment"`
	Factors    []string          `yaml:"factors,omitempty"`
	Ratings    map[string]string `yaml:"ratings,omitempty"`
}

// Load reads the spec at path (defaulting to QUALITY.md) and parses its
// frontmatter.
func Load(path string) (*Spec, error) {
	if path == "" {
		path = "QUALITY.md"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading spec: %w", err)
	}

	fm, err := frontmatter(raw)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}

	var s Spec
	decoder := yaml.NewDecoder(bytes.NewReader(fm))
	decoder.KnownFields(true)
	if err := decoder.Decode(&s); err != nil {
		return nil, fmt.Errorf("%s: parsing spec: %w", path, err)
	}
	s.Path = path
	if err := validate(s); err != nil {
		return nil, fmt.Errorf("%s: invalid QUALITY.md: %w", path, err)
	}
	return &s, nil
}

func validate(s Spec) error {
	if len(s.RatingScale) < 2 {
		return fmt.Errorf("ratingScale must declare at least two levels")
	}
	levels := map[string]bool{}
	for i, level := range s.RatingScale {
		if level.Level == "" {
			return fmt.Errorf("ratingScale[%d].level is required", i)
		}
		if levels[level.Level] {
			return fmt.Errorf("ratingScale level %q is duplicated", level.Level)
		}
		levels[level.Level] = true
		if level.Criterion == "" {
			return fmt.Errorf("ratingScale[%d].criterion is required", i)
		}
	}
	if len(s.Factors) == 0 && len(s.Requirements) == 0 && len(s.Targets) == 0 {
		return fmt.Errorf("one of factors, requirements, or targets is required")
	}
	for name, req := range s.Requirements {
		if err := validateRequirement(name, req, levels); err != nil {
			return err
		}
	}
	for name, factor := range s.Factors {
		if err := validateFactor("factor "+name, factor, levels); err != nil {
			return err
		}
	}
	for name, target := range s.Targets {
		if err := validateTarget("target "+name, target, levels); err != nil {
			return err
		}
	}
	return nil
}

func validateTarget(path string, target Target, levels map[string]bool) error {
	for name, req := range target.Requirements {
		if err := validateRequirement(path+" requirement "+name, req, levels); err != nil {
			return err
		}
	}
	for name, factor := range target.Factors {
		if err := validateFactor(path+" factor "+name, factor, levels); err != nil {
			return err
		}
	}
	for name, child := range target.Targets {
		if err := validateTarget(path+" target "+name, child, levels); err != nil {
			return err
		}
	}
	return nil
}

func validateFactor(path string, factor Factor, levels map[string]bool) error {
	for name, req := range factor.Requirements {
		if err := validateRequirement(path+" requirement "+name, req, levels); err != nil {
			return err
		}
	}
	for name, child := range factor.Factors {
		if err := validateFactor(path+" factor "+name, child, levels); err != nil {
			return err
		}
	}
	return nil
}

func validateRequirement(path string, req Requirement, levels map[string]bool) error {
	if req.Assessment == "" {
		return fmt.Errorf("%s: assessment is required", path)
	}
	for level := range req.Ratings {
		if !levels[level] {
			return fmt.Errorf("%s: ratings override names unknown level %q", path, level)
		}
	}
	return nil
}

// frontmatter extracts the YAML block delimited by a leading and trailing
// `---` fence. A document without a fence is treated as YAML in its entirety.
func frontmatter(raw []byte) ([]byte, error) {
	const fence = "---"
	trimmed := bytes.TrimLeft(raw, " \t\r\n")
	if !bytes.HasPrefix(trimmed, []byte(fence)) {
		return raw, nil
	}
	rest := trimmed[len(fence):]
	if i := bytes.Index(rest, []byte("\n"+fence)); i >= 0 {
		return rest[:i], nil
	}
	return nil, fmt.Errorf("unterminated frontmatter: missing closing %q", fence)
}
