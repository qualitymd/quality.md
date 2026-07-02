package evaluation

import (
	"strings"

	"github.com/qualitymd/quality.md/internal/model"
)

// ParseRatingReference resolves a canonical rating level model reference against
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

// ParseUnqualifiedRatingReference resolves a rating level reference in a context
// where the expected reference type is fixed.
func ParseUnqualifiedRatingReference(spec *model.Spec, ref string) (string, error) {
	if spec == nil {
		return "", usagef("model reference %q cannot resolve without a model", ref)
	}
	return parseRatingReferenceBody(spec, ref, ref, "unqualified rating reference")
}

func parseRatingReferenceBody(spec *model.Spec, ref, level, label string) (string, error) {
	if !model.ValidReferenceName(level) {
		return "", usagef("%s %q has an invalid rating level ID", label, ref)
	}
	for _, candidate := range spec.RatingScale {
		if candidate.Level == level {
			return level, nil
		}
	}
	return "", usagef("%s %q does not resolve in the model", label, ref)
}
