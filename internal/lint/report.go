package lint

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

type repairOp struct {
	record RepairRecord
	apply  func() error
}

type runState struct {
	doc      *document.Document
	options  Options
	root     *areaRef
	findings []Finding
	repairs  []repairOp
	levels   map[string]bool
}

func newRunState(doc *document.Document, options Options) *runState {
	return &runState{
		doc:     doc,
		options: options,
		levels:  map[string]bool{},
	}
}

func (s *runState) result(repairs []RepairRecord) Result {
	summary := Summary{Fixed: len(repairs)}
	for _, finding := range s.findings {
		switch finding.Severity {
		case SeverityError:
			summary.Errors++
		case SeverityWarning:
			summary.Warnings++
		case SeverityInfo:
			summary.Info++
		}
		if finding.Fixable {
			summary.Fixable++
		}
	}
	if repairs == nil {
		repairs = []RepairRecord{}
	}
	findings := s.findings
	if findings == nil {
		findings = []Finding{}
	}
	return Result{
		SchemaVersion: schemaVersion,
		Path:          s.doc.Path,
		Valid:         summary.Errors == 0,
		Summary:       summary,
		Findings:      findings,
		Repairs:       repairs,
		NextActions:   nextActions(s.doc.Path, summary),
	}
}

func nextActions(path string, summary Summary) []receipt.Action {
	if summary.Fixable > 0 {
		return []receipt.Action{
			{
				ID:      "fix",
				Label:   "Apply deterministic lint repairs",
				Command: "qualitymd lint --fix " + path,
			},
		}
	}
	if summary.Errors > 0 {
		return []receipt.Action{
			{
				ID:      "rerun-lint",
				Label:   "Re-run validation",
				Command: "qualitymd lint " + path,
			},
		}
	}
	return []receipt.Action{}
}

func (s *runState) checkOptionalScalar(parent, key, value *yaml.Node, path []PathSegment, locationLabel string) {
	if isEmpty(value) {
		s.emptyProperty(parent, key, path, locationLabel)
		return
	}
	if value.Kind != yaml.ScalarNode {
		s.invalid(key, path, locationLabel, "The `"+key.Value+"` property has the wrong YAML shape; it must be a scalar.")
	}
}

func (s *runState) emptyProperty(parent, key *yaml.Node, path []PathSegment, locationLabel string) {
	message := "The optional property `" + key.Value + "` is empty; empty optional properties should be omitted."
	location := s.loc(key, path, locationLabel)
	s.add(RuleEmptyProperty, message, location, &repairOp{
		record: RepairRecord{
			RuleID:   RuleEmptyProperty,
			Message:  "Removed empty optional property `" + key.Value + "`.",
			Location: location,
		},
		apply: func() error {
			if !document.RemoveMapEntry(parent, key.Value) {
				return fmt.Errorf("property no longer exists")
			}
			return nil
		},
	})
}

func (s *runState) invalid(node *yaml.Node, path []PathSegment, locationLabel, message string) {
	s.add(RuleInvalidFrontmatter, message, s.loc(node, path, locationLabel), nil)
}

// add records a finding for ruleID, deriving its severity and fixability from
// the rule catalog so the catalog is the single source of truth.
func (s *runState) add(ruleID RuleID, message string, location Location, repair *repairOp) {
	rule := rulesByID[ruleID]
	s.findings = append(s.findings, Finding{
		RuleID:   ruleID,
		Severity: rule.Severity,
		Message:  message,
		Location: location,
		Fixable:  rule.Fixable,
	})
	if repair != nil {
		s.repairs = append(s.repairs, *repair)
	}
}

func (s *runState) loc(node *yaml.Node, path []PathSegment, locationLabel string) Location {
	location := Location{
		Path:      s.doc.Path,
		ModelPath: clonePath(path),
		Label:     locationLabel,
	}
	if node != nil && node.Line > 0 {
		// The +1 compensates for the frontmatter YAML being parsed from the block
		// that begins on the line after the opening "---" fence.
		location.Line = node.Line + 1
		location.Column = node.Column
	}
	return location
}

func (s *runState) locForMissing(path []PathSegment, locationLabel string) Location {
	return Location{
		Path:      s.doc.Path,
		ModelPath: clonePath(path),
		Label:     locationLabel,
	}
}

func (s *runState) locForNodeOrMissing(node *yaml.Node, path []PathSegment, locationLabel string) Location {
	if node == nil {
		return s.locForMissing(path, locationLabel)
	}
	return s.loc(node, path, locationLabel)
}

func (s *runState) sort() {
	slices.SortFunc(s.findings, compareFindings)
	slices.SortFunc(s.repairs, func(a, b repairOp) int {
		return compareLocations(a.record.Location, b.record.Location)
	})
}

func isEmpty(node *yaml.Node) bool {
	if node == nil {
		return true
	}
	switch node.Kind {
	case yaml.ScalarNode:
		return node.Tag == "!!null" || strings.TrimSpace(node.Value) == ""
	case yaml.MappingNode, yaml.SequenceNode:
		return len(node.Content) == 0
	default:
		return false
	}
}

func appendPath(path []PathSegment, parts ...PathSegment) []PathSegment {
	out := clonePath(path)
	out = append(out, parts...)
	return out
}

func clonePath(path []PathSegment) []PathSegment {
	if path == nil {
		return []PathSegment{}
	}
	out := make([]PathSegment, len(path))
	copy(out, path)
	return out
}

func label(path []PathSegment) string {
	if len(path) == 0 {
		return "frontmatter"
	}
	parts := make([]string, 0, len(path))
	for _, part := range path {
		switch v := part.(type) {
		case string:
			parts = append(parts, v)
		case int:
			parts = append(parts, strconv.Itoa(v))
		default:
			parts = append(parts, fmt.Sprint(v))
		}
	}
	return strings.Join(parts, ".")
}
