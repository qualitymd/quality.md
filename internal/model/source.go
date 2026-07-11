package model

// Source resolution is model semantics (SPECIFICATION.md §Source resolution):
// an area evaluates the source it declares, else the nearest declaring
// ancestor's, else the document default — the directory containing the
// QUALITY.md file. One resolver serves every consumer (the evaluation runner,
// status) so the surfaces cannot silently disagree.

// SourceState describes how an area's effective source selector is resolved.
type SourceState string

const (
	// SourceStateDeclared means the area declares its own source.
	SourceStateDeclared SourceState = "declared"
	// SourceStateInherited means the area inherits a source declared by an
	// ancestor area.
	SourceStateInherited SourceState = "inherited"
	// SourceStateDefault means no area in the chain declares a source, so the
	// area resolves to the document's default source: the directory containing
	// the QUALITY.md file. This is a deliberate, valid choice, not a defect.
	SourceStateDefault SourceState = "default"
)

// DefaultSource is the effective selector when no area in the chain declares
// a source: the directory containing the QUALITY.md document
// (SPECIFICATION.md §Document structure).
const DefaultSource = "."

// EffectiveSource resolves the area at path to its effective source selector
// and the state that produced it. The nearest declaring area on the chain
// from the root to the target wins; a chain with no declaration resolves to
// DefaultSource.
func EffectiveSource(spec *Spec, path AreaPath) (string, SourceState) {
	selector := spec.Source
	// declaredAt indexes the declaring node on the chain: 0 is the root,
	// i+1 is path[i], and the target is len(path).
	declaredAt := -1
	if selector != "" {
		declaredAt = 0
	}
	areas := spec.Areas
	for i, name := range path {
		area, ok := areas[name]
		if !ok {
			break
		}
		if area.Source != "" {
			selector = area.Source
			declaredAt = i + 1
		}
		areas = area.Areas
	}
	switch {
	case declaredAt == len(path):
		return selector, SourceStateDeclared
	case declaredAt >= 0:
		return selector, SourceStateInherited
	default:
		return DefaultSource, SourceStateDefault
	}
}
