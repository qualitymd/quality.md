package evaluation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/qualitymd/quality.md/internal/model"
	qschema "github.com/qualitymd/quality.md/internal/schema"
)

var modelReferenceNamePattern = regexp.MustCompile(qschema.ModelNamePattern)

// ParseAreaReference resolves a canonical Area model reference against spec.
func ParseAreaReference(spec *model.Spec, ref string) (AreaPath, error) {
	if spec == nil {
		return nil, usagef("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "area:")
	if !ok {
		return nil, usagef("area model reference %q must start with area:", ref)
	}
	return parseAreaReferenceBody(spec, ref, body, "area model reference")
}

// ParseUnqualifiedAreaReference resolves an Area reference in a context where
// the expected reference type is fixed.
func ParseUnqualifiedAreaReference(spec *model.Spec, ref string) (AreaPath, error) {
	if spec == nil {
		return nil, usagef("model reference %q cannot resolve without a model", ref)
	}
	return parseAreaReferenceBody(spec, ref, ref, "unqualified area reference")
}

func parseAreaReferenceBody(spec *model.Spec, ref, body, label string) (AreaPath, error) {
	path, err := parseReferencePath(body)
	if err != nil {
		return nil, usagef("%s %q is invalid: %w", label, ref, err)
	}
	if !areaExists(spec, path) {
		return nil, usagef("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(path), nil
}

// ParseFactorReference resolves a canonical Factor model reference against spec.
func ParseFactorReference(spec *model.Spec, ref string) (AreaPath, FactorPath, error) {
	if spec == nil {
		return nil, nil, usagef("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "factor:")
	if !ok {
		return nil, nil, usagef("factor model reference %q must start with factor:", ref)
	}
	return parseFactorReferenceBody(spec, ref, body, "factor model reference")
}

// ParseUnqualifiedFactorReference resolves a Factor reference in a context where
// the expected reference type is fixed.
func ParseUnqualifiedFactorReference(spec *model.Spec, ref string) (AreaPath, FactorPath, error) {
	if spec == nil {
		return nil, nil, usagef("model reference %q cannot resolve without a model", ref)
	}
	return parseFactorReferenceBody(spec, ref, ref, "unqualified factor reference")
}

func parseFactorReferenceBody(spec *model.Spec, ref, body, label string) (AreaPath, FactorPath, error) {
	areaPart, factorPart, ok := strings.Cut(body, "::")
	if !ok {
		return nil, nil, usagef("%s %q must separate area and factor paths with ::", label, ref)
	}
	areaPath, err := parseReferencePath(areaPart)
	if err != nil {
		return nil, nil, usagef("%s %q has invalid area path: %w", label, ref, err)
	}
	factorPath, err := parseReferencePath(factorPart)
	if err != nil {
		return nil, nil, usagef("%s %q has invalid factor path: %w", label, ref, err)
	}
	if len(factorPath) == 0 {
		return nil, nil, usagef("%s %q must name a factor path", label, ref)
	}
	if !areaExists(spec, areaPath) {
		return nil, nil, usagef("%s %q declares an area that does not resolve in the model", label, ref)
	}
	if !factorExists(spec, areaPath, factorPath) {
		return nil, nil, usagef("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(areaPath), FactorPath(factorPath), nil
}

// ParseRequirementReference resolves a canonical Requirement model reference
// against spec.
func ParseRequirementReference(spec *model.Spec, ref string) (AreaPath, string, error) {
	if spec == nil {
		return nil, "", usagef("model reference %q cannot resolve without a model", ref)
	}
	body, ok := strings.CutPrefix(ref, "requirement:")
	if !ok {
		return nil, "", usagef("requirement model reference %q must start with requirement:", ref)
	}
	return parseRequirementReferenceBody(spec, ref, body, "requirement model reference")
}

// ParseUnqualifiedRequirementReference resolves a Requirement reference in a
// context where the expected reference type is fixed.
func ParseUnqualifiedRequirementReference(spec *model.Spec, ref string) (AreaPath, string, error) {
	if spec == nil {
		return nil, "", usagef("model reference %q cannot resolve without a model", ref)
	}
	return parseRequirementReferenceBody(spec, ref, ref, "unqualified requirement reference")
}

func parseRequirementReferenceBody(spec *model.Spec, ref, body, label string) (AreaPath, string, error) {
	areaPart, requirementName, ok := strings.Cut(body, "::")
	if !ok {
		return nil, "", usagef("%s %q must separate area path and Requirement name with ::", label, ref)
	}
	areaPath, err := parseReferencePath(areaPart)
	if err != nil {
		return nil, "", usagef("%s %q has invalid area path: %w", label, ref, err)
	}
	if !validReferenceName(requirementName) {
		return nil, "", usagef("%s %q has an invalid Requirement name", label, ref)
	}
	if !areaExists(spec, areaPath) {
		return nil, "", usagef("%s %q declares an area that does not resolve in the model", label, ref)
	}
	if !requirementExists(spec, areaPath, requirementName) {
		return nil, "", usagef("%s %q does not resolve in the model", label, ref)
	}
	return AreaPath(areaPath), requirementName, nil
}

// ParseRatingReference resolves a canonical Rating Level model reference against
// spec.
func ParseRatingReference(spec *model.Spec, ref string) (string, error) {
	if spec == nil {
		return "", usagef("model reference %q cannot resolve without a model", ref)
	}
	level, ok := strings.CutPrefix(ref, "rating:")
	if !ok {
		return "", usagef("rating model reference %q must start with rating:", ref)
	}
	return parseRatingReferenceBody(spec, ref, level, "rating model reference")
}

// ParseUnqualifiedRatingReference resolves a Rating Level reference in a context
// where the expected reference type is fixed.
func ParseUnqualifiedRatingReference(spec *model.Spec, ref string) (string, error) {
	if spec == nil {
		return "", usagef("model reference %q cannot resolve without a model", ref)
	}
	return parseRatingReferenceBody(spec, ref, ref, "unqualified rating reference")
}

func parseRatingReferenceBody(spec *model.Spec, ref, level, label string) (string, error) {
	if !validReferenceName(level) {
		return "", usagef("%s %q has an invalid Rating Level ID", label, ref)
	}
	for _, candidate := range spec.RatingScale {
		if candidate.Level == level {
			return level, nil
		}
	}
	return "", usagef("%s %q does not resolve in the model", label, ref)
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
		if !validReferenceName(part) {
			return nil, fmt.Errorf("segment %q does not match %s", part, qschema.ModelNamePattern)
		}
	}
	return parts, nil
}

func validReferenceName(name string) bool {
	return modelReferenceNamePattern.MatchString(name)
}

func areaExists(spec *model.Spec, path []string) bool {
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

func factorExists(spec *model.Spec, areaPath []string, factorPath []string) bool {
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

func requirementExists(spec *model.Spec, areaPath []string, requirementName string) bool {
	areaRequirements, areaFactors, ok := requirementsForArea(spec, areaPath)
	if !ok {
		return false
	}
	if _, ok := areaRequirements[requirementName]; ok {
		return true
	}
	return factorRequirementExists(areaFactors, requirementName)
}

func factorRequirementExists(factors map[string]model.Factor, requirementName string) bool {
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

func requirementsForArea(spec *model.Spec, areaPath []string) (map[string]model.Requirement, map[string]model.Factor, bool) {
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

func factorsForArea(spec *model.Spec, areaPath []string) (map[string]model.Factor, bool) {
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
