package lint

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/qualitymd/quality.md/internal/receipt"
)

const schemaVersion = 1

// Severity is a lint finding severity.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// RuleID identifies a lint rule. It is a named string type so it marshals to a
// plain JSON string while keeping emission sites type-checked against the catalog.
type RuleID string

const (
	RuleInvalidFrontmatter       RuleID = "invalid-frontmatter"
	RuleMissingTitle             RuleID = "missing-title"
	RuleMissingRatingScale       RuleID = "missing-rating-scale"
	RuleTooFewLevels             RuleID = "too-few-levels"
	RuleMissingLevelName         RuleID = "missing-level-name"
	RuleDuplicateLevel           RuleID = "duplicate-level"
	RuleMissingCriterion         RuleID = "missing-criterion"
	RuleMissingLevelDescription  RuleID = "missing-level-description"
	RuleEmptyModel               RuleID = "empty-model"
	RuleMisplacedRootKey         RuleID = "misplaced-root-key"
	RuleEmptyTarget              RuleID = "empty-target"
	RuleEmptyFactor              RuleID = "empty-factor"
	RuleMissingFactorDescription RuleID = "missing-factor-description"
	RuleInvalidAssessment        RuleID = "invalid-assessment"
	RuleUnknownRatingKey         RuleID = "unknown-rating-key"
	RuleUnknownFactor            RuleID = "unknown-factor"
	RuleEmptyProperty            RuleID = "empty-property"
)

// Rule describes one lint rule in the catalog.
type Rule struct {
	ID          RuleID
	Severity    Severity
	Fixable     bool
	Description string
}

// Rules is the catalog of every lint rule and the single source of truth for
// each rule's severity and fixability.
var Rules = []Rule{
	{RuleInvalidFrontmatter, SeverityError, false, "The frontmatter is missing or has the wrong shape for a QUALITY.md model."},
	{RuleMissingTitle, SeverityWarning, false, "The model root declares no title."},
	{RuleMissingRatingScale, SeverityError, false, "The model root declares no rating scale."},
	{RuleTooFewLevels, SeverityError, false, "The rating scale has fewer than two levels."},
	{RuleMissingLevelName, SeverityError, false, "A rating level declares no level name."},
	{RuleDuplicateLevel, SeverityError, false, "A rating level name is duplicated within the rating scale."},
	{RuleMissingCriterion, SeverityError, false, "A rating level declares no criterion."},
	{RuleMissingLevelDescription, SeverityWarning, false, "A rating level declares no description."},
	{RuleEmptyModel, SeverityError, false, "The model root supplies no factors, requirements, or targets."},
	{RuleMisplacedRootKey, SeverityError, false, "A root-only key appears on a nested target."},
	{RuleEmptyTarget, SeverityWarning, false, "A target reaches no requirements in its subtree."},
	{RuleEmptyFactor, SeverityWarning, false, "A factor leads to no requirements."},
	{RuleMissingFactorDescription, SeverityWarning, false, "A factor declares no description."},
	{RuleInvalidAssessment, SeverityError, false, "A requirement has no single non-empty scalar assessment."},
	{RuleUnknownRatingKey, SeverityError, false, "A ratings override names a level outside the rating scale."},
	{RuleUnknownFactor, SeverityError, false, "A requirement names a secondary factor that does not resolve."},
	{RuleEmptyProperty, SeverityWarning, true, "An optional property is present but empty and should be omitted."},
}

// rulesByID indexes the catalog so the recording path can derive a finding's
// severity (and fixability) from its rule ID.
var rulesByID = func() map[RuleID]Rule {
	m := make(map[RuleID]Rule, len(Rules))
	for _, rule := range Rules {
		m[rule.ID] = rule
	}
	return m
}()

// PathSegment is one element of a model path: a string key or an integer index.
type PathSegment = any

// Result is the JSON contract emitted by qualitymd lint.
type Result struct {
	SchemaVersion int              `json:"schemaVersion"`
	Path          string           `json:"path"`
	Valid         bool             `json:"valid"`
	Summary       Summary          `json:"summary"`
	Findings      []Finding        `json:"findings"`
	Repairs       []RepairRecord   `json:"repairs"`
	NextActions   []receipt.Action `json:"nextActions"`
}

// Summary counts findings and repairs in a lint result.
type Summary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
	Fixable  int `json:"fixable"`
	Fixed    int `json:"fixed"`
}

// Finding is one rule violation.
type Finding struct {
	RuleID   RuleID   `json:"ruleId"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	Location Location `json:"location"`
	Fixable  bool     `json:"fixable"`
}

// Location names the smallest stable model location for a finding.
type Location struct {
	Path      string        `json:"path"`
	ModelPath []PathSegment `json:"modelPath"`
	Label     string        `json:"label"`
	Line      int           `json:"line,omitempty"`
	Column    int           `json:"column,omitempty"`
}

// RepairRecord reports one repair applied by --fix.
type RepairRecord struct {
	RuleID   RuleID   `json:"ruleId"`
	Message  string   `json:"message"`
	Location Location `json:"location"`
}

// JSON formats the result as the lint JSON document.
func (r Result) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}

// Err returns a non-nil error when the result contains error findings.
func (r Result) Err() error {
	if r.Summary.Errors == 0 {
		return nil
	}
	return lintError{result: r}
}

func compareFindings(a, b Finding) int {
	if cmp := compareLocations(a.Location, b.Location); cmp != 0 {
		return cmp
	}
	if cmp := severityRank(a.Severity) - severityRank(b.Severity); cmp != 0 {
		return cmp
	}
	return strings.Compare(string(a.RuleID), string(b.RuleID))
}

func compareLocations(a, b Location) int {
	if a.Line > 0 && b.Line > 0 {
		if a.Line != b.Line {
			return a.Line - b.Line
		}
		if a.Column != b.Column {
			return a.Column - b.Column
		}
	}
	if cmp := comparePath(a.ModelPath, b.ModelPath); cmp != 0 {
		return cmp
	}
	return strings.Compare(a.Label, b.Label)
}

func comparePath(a, b []PathSegment) int {
	minLen := min(len(a), len(b))
	for i := 0; i < minLen; i++ {
		switch av := a[i].(type) {
		case int:
			if bv, ok := b[i].(int); ok {
				if av != bv {
					return av - bv
				}
				continue
			}
		case string:
			if bv, ok := b[i].(string); ok {
				if cmp := strings.Compare(av, bv); cmp != 0 {
					return cmp
				}
				continue
			}
		}
		if cmp := strings.Compare(fmt.Sprint(a[i]), fmt.Sprint(b[i])); cmp != 0 {
			return cmp
		}
	}
	return len(a) - len(b)
}

func severityRank(severity Severity) int {
	switch severity {
	case SeverityError:
		return 0
	case SeverityWarning:
		return 1
	case SeverityInfo:
		return 2
	default:
		return 3
	}
}
