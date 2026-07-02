package model

import (
	"fmt"
	"regexp"
	"strings"

	qschema "github.com/qualitymd/quality.md/internal/schema"
)

var modelReferenceNamePattern = regexp.MustCompile(qschema.ModelNamePattern)

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

// Display returns the human-facing area path label.
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

// Reference returns the canonical typed model reference for an area path.
func (p AreaPath) Reference() string {
	return "area:" + p.referencePath()
}

// UnqualifiedReference returns the fixed-type area reference.
func (p AreaPath) UnqualifiedReference() string {
	return p.referencePath()
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

// Display returns the human-facing factor path label.
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

// FactorReference returns the canonical typed model reference for a factor path
// declared by areaPath.
func FactorReference(areaPath AreaPath, factorPath FactorPath) string {
	return "factor:" + areaPath.referencePath() + "::" + factorPath.referencePath()
}

// UnqualifiedFactorReference returns the fixed-type factor reference.
func UnqualifiedFactorReference(areaPath AreaPath, factorPath FactorPath) string {
	return areaPath.referencePath() + "::" + factorPath.referencePath()
}

// RequirementReference returns the canonical typed model reference for a
// requirement name declared by areaPath.
func RequirementReference(areaPath AreaPath, requirementName string) string {
	return "requirement:" + areaPath.referencePath() + "::" + requirementName
}

// UnqualifiedRequirementReference returns the fixed-type requirement reference.
func UnqualifiedRequirementReference(areaPath AreaPath, requirementName string) string {
	return areaPath.referencePath() + "::" + requirementName
}

// ParseAreaReference resolves a canonical area model reference against spec.
func ParseAreaReference(spec *Spec, ref string) (AreaPath, error) {
	if spec == nil {
		return nil, fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "area:")
	if !ok {
		return nil, fmt.Errorf("area model reference %q must start with area: prefix", ref)
	}
	return parseAreaReferenceBody(spec, ref, body, "area model reference")
}

// ParseUnqualifiedAreaReference resolves an area reference in a context where
// the expected reference type is fixed.
func ParseUnqualifiedAreaReference(spec *Spec, ref string) (AreaPath, error) {
	if spec == nil {
		return nil, fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	return parseAreaReferenceBody(spec, ref, ref, "unqualified area reference")
}

func parseAreaReferenceBody(spec *Spec, ref, body, label string) (AreaPath, error) {
	path, err := parseReferencePath(body)
	if err != nil {
		return nil, fmt.Errorf("%s %q is invalid: %w", label, ref, err)
	}
	if !AreaExists(spec, path) {
		return nil, fmt.Errorf("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(path), nil
}

// ParseFactorReference resolves a canonical factor model reference against spec.
func ParseFactorReference(spec *Spec, ref string) (AreaPath, FactorPath, error) {
	if spec == nil {
		return nil, nil, fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "factor:")
	if !ok {
		return nil, nil, fmt.Errorf("factor model reference %q must start with factor: prefix", ref)
	}
	return parseFactorReferenceBody(spec, ref, body, "factor model reference")
}

// ParseUnqualifiedFactorReference resolves a factor reference in a context where
// the expected reference type is fixed.
func ParseUnqualifiedFactorReference(spec *Spec, ref string) (AreaPath, FactorPath, error) {
	if spec == nil {
		return nil, nil, fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	return parseFactorReferenceBody(spec, ref, ref, "unqualified factor reference")
}

func parseFactorReferenceBody(spec *Spec, ref, body, label string) (AreaPath, FactorPath, error) {
	areaPart, factorPart, ok := strings.Cut(body, "::")
	if !ok {
		return nil, nil, fmt.Errorf("%s %q must contain :: between area and factor paths", label, ref)
	}
	areaPath, err := parseReferencePath(areaPart)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %q has invalid area path: %w", label, ref, err)
	}
	factorPath, err := parseReferencePath(factorPart)
	if err != nil {
		return nil, nil, fmt.Errorf("%s %q has invalid factor path: %w", label, ref, err)
	}
	if len(factorPath) == 0 {
		return nil, nil, fmt.Errorf("%s %q must name a factor path", label, ref)
	}
	if !AreaExists(spec, areaPath) {
		return nil, nil, fmt.Errorf("%s %q declares an area that does not resolve in the model", label, ref)
	}
	if !FactorExists(spec, areaPath, factorPath) {
		return nil, nil, fmt.Errorf("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(areaPath), FactorPath(factorPath), nil
}

// ParseRequirementReference resolves a canonical requirement model reference
// against spec.
func ParseRequirementReference(spec *Spec, ref string) (AreaPath, string, error) {
	if spec == nil {
		return nil, "", fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "requirement:")
	if !ok {
		return nil, "", fmt.Errorf("requirement model reference %q must start with requirement: prefix", ref)
	}
	return parseRequirementReferenceBody(spec, ref, body, "requirement model reference")
}

// ParseUnqualifiedRequirementReference resolves a requirement reference in a
// context where the expected reference type is fixed.
func ParseUnqualifiedRequirementReference(spec *Spec, ref string) (AreaPath, string, error) {
	if spec == nil {
		return nil, "", fmt.Errorf("model reference %q cannot resolve without a model", ref)
	}
	return parseRequirementReferenceBody(spec, ref, ref, "unqualified requirement reference")
}

func parseRequirementReferenceBody(spec *Spec, ref, body, label string) (AreaPath, string, error) {
	areaPart, requirementName, ok := strings.Cut(body, "::")
	if !ok {
		return nil, "", fmt.Errorf("%s %q must contain :: between area path and requirement name", label, ref)
	}
	areaPath, err := parseReferencePath(areaPart)
	if err != nil {
		return nil, "", fmt.Errorf("%s %q has invalid area path: %w", label, ref, err)
	}
	if !ValidReferenceName(requirementName) {
		return nil, "", fmt.Errorf("%s %q has an invalid requirement name", label, ref)
	}
	if !AreaExists(spec, areaPath) {
		return nil, "", fmt.Errorf("%s %q declares an area that does not resolve in the model", label, ref)
	}
	if !RequirementExists(spec, areaPath, requirementName) {
		return nil, "", fmt.Errorf("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(areaPath), requirementName, nil
}

func parseReferencePath(raw string) ([]string, error) {
	if raw == "root" {
		return nil, nil
	}
	if raw == "" {
		return nil, fmt.Errorf("path is empty")
	}
	parts := strings.Split(raw, "/")
	for _, part := range parts {
		if !ValidReferenceName(part) {
			return nil, fmt.Errorf("segment %q does not match %s", part, qschema.ModelNamePattern)
		}
	}
	return parts, nil
}

// ValidReferenceName reports whether name matches the strict model-name grammar
// shared by area, factor, requirement, and rating level identifiers.
func ValidReferenceName(name string) bool {
	return modelReferenceNamePattern.MatchString(name)
}

// AreaExists reports whether path resolves to an area declared in spec.
func AreaExists(spec *Spec, path []string) bool {
	areas := spec.Areas
	for _, element := range path {
		area, ok := areas[element]
		if !ok {
			return false
		}
		areas = area.Areas
	}
	return true
}

// FactorExists reports whether factorPath resolves to a factor declared by the
// area at areaPath in spec.
func FactorExists(spec *Spec, areaPath []string, factorPath []string) bool {
	factors, ok := factorsForArea(spec, areaPath)
	if !ok {
		return false
	}
	for i, element := range factorPath {
		factor, ok := factors[element]
		if !ok {
			return false
		}
		if i == len(factorPath)-1 {
			return true
		}
		factors = factor.Factors
	}
	return false
}

// RequirementExists reports whether requirementName resolves to a requirement
// declared by the area at areaPath in spec, directly or under one of its factors.
func RequirementExists(spec *Spec, areaPath []string, requirementName string) bool {
	areaRequirements, areaFactors, ok := requirementsForArea(spec, areaPath)
	if !ok {
		return false
	}
	if _, ok := areaRequirements[requirementName]; ok {
		return true
	}
	return factorRequirementExists(areaFactors, requirementName)
}

func factorRequirementExists(factors map[string]Factor, requirementName string) bool {
	for _, factor := range factors {
		if _, ok := factor.Requirements[requirementName]; ok {
			return true
		}
		if factorRequirementExists(factor.Factors, requirementName) {
			return true
		}
	}
	return false
}

func requirementsForArea(spec *Spec, areaPath []string) (map[string]Requirement, map[string]Factor, bool) {
	if len(areaPath) == 0 {
		return spec.Requirements, spec.Factors, true
	}
	areas := spec.Areas
	for i, element := range areaPath {
		area, ok := areas[element]
		if !ok {
			return nil, nil, false
		}
		if i == len(areaPath)-1 {
			return area.Requirements, area.Factors, true
		}
		areas = area.Areas
	}
	return nil, nil, false
}

func factorsForArea(spec *Spec, areaPath []string) (map[string]Factor, bool) {
	if len(areaPath) == 0 {
		return spec.Factors, true
	}
	areas := spec.Areas
	for i, element := range areaPath {
		area, ok := areas[element]
		if !ok {
			return nil, false
		}
		if i == len(areaPath)-1 {
			return area.Factors, true
		}
		areas = area.Areas
	}
	return nil, false
}
