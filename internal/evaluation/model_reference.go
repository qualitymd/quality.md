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
	path, err := parseReferencePath(body)
	if err != nil {
		return nil, usagef("area model reference %q is invalid: %w", ref, err)
	}
	if !areaExists(spec, path) {
		return nil, usagef("area model reference %q does not resolve in the model", ref)
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
	areaPart, factorPart, ok := strings.Cut(body, "::")
	if !ok {
		return nil, nil, usagef("factor model reference %q must separate area and factor paths with ::", ref)
	}
	areaPath, err := parseReferencePath(areaPart)
	if err != nil {
		return nil, nil, usagef("factor model reference %q has invalid area path: %w", ref, err)
	}
	factorPath, err := parseReferencePath(factorPart)
	if err != nil {
		return nil, nil, usagef("factor model reference %q has invalid factor path: %w", ref, err)
	}
	if len(factorPath) == 0 {
		return nil, nil, usagef("factor model reference %q must name a factor path", ref)
	}
	if !areaExists(spec, areaPath) {
		return nil, nil, usagef("factor model reference %q declares an area that does not resolve in the model", ref)
	}
	if !factorExists(spec, areaPath, factorPath) {
		return nil, nil, usagef("factor model reference %q does not resolve in the model", ref)
	}
	return AreaPath(areaPath), FactorPath(factorPath), nil
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
	if !validReferenceName(level) {
		return "", usagef("rating model reference %q has an invalid Rating Level ID", ref)
	}
	for _, candidate := range spec.RatingScale {
		if candidate.Level == level {
			return level, nil
		}
	}
	return "", usagef("rating model reference %q does not resolve in the model", ref)
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
