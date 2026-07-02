package model

import "slices"

// Kind identifies the element type a projected node represents.
type Kind string

const (
	KindArea        Kind = "area"
	KindFactor      Kind = "factor"
	KindRequirement Kind = "requirement"
)

// Element is one node in the read-only model projection: its canonical
// reference ID, kind, label, the ID of its structural parent (empty for the
// root area), and its immediate children in the documented projection order.
type Element struct {
	ID       string     `json:"id"`
	Kind     Kind       `json:"kind"`
	Label    string     `json:"label"`
	ParentID string     `json:"parentId,omitempty"`
	Children []*Element `json:"children,omitempty"`
}

// Project walks spec and returns the rooted Element tree. One recursive walk
// visits the rooted area, its factors (sub-factors nested), its requirements,
// then child areas — the documented SB4 order — with siblings of each kind in
// lexicographic key order. The result is byte-stable for a given model.
func Project(spec *Spec) *Element {
	root := &Element{
		ID:    AreaPath{}.Reference(),
		Kind:  KindArea,
		Label: areaLabel(spec.Title, "root"),
	}
	buildArea(root, nil, spec.Factors, spec.Requirements, spec.Areas)
	return root
}

// buildArea populates an area element's children from the factor, requirement,
// and child-area maps declared at areaPath.
func buildArea(parent *Element, areaPath AreaPath, factors map[string]Factor, requirements map[string]Requirement, areas map[string]Area) {
	for _, name := range sortedKeys(factors) {
		factor := factors[name]
		parent.Children = append(parent.Children, buildFactor(parent.ID, areaPath, FactorPath{name}, factor))
	}
	for _, name := range sortedKeys(requirements) {
		req := requirements[name]
		parent.Children = append(parent.Children, &Element{
			ID:       RequirementReference(areaPath, name),
			Kind:     KindRequirement,
			Label:    labelOrKey(req.Title, name),
			ParentID: parent.ID,
		})
	}
	for _, name := range sortedKeys(areas) {
		area := areas[name]
		childPath := appendPath(areaPath, name)
		child := &Element{
			ID:       childPath.Reference(),
			Kind:     KindArea,
			Label:    labelOrKey(area.Title, name),
			ParentID: parent.ID,
		}
		buildArea(child, childPath, area.Factors, area.Requirements, area.Areas)
		parent.Children = append(parent.Children, child)
	}
}

// buildFactor populates a factor element's children: its sub-factors (nested)
// then its directly declared requirements.
func buildFactor(parentID string, areaPath AreaPath, factorPath FactorPath, factor Factor) *Element {
	node := &Element{
		ID:       FactorReference(areaPath, factorPath),
		Kind:     KindFactor,
		Label:    labelOrKey(factor.Title, factorPath[len(factorPath)-1]),
		ParentID: parentID,
	}
	for _, name := range sortedKeys(factor.Factors) {
		sub := factor.Factors[name]
		node.Children = append(node.Children, buildFactor(node.ID, areaPath, appendFactorPath(factorPath, name), sub))
	}
	for _, name := range sortedKeys(factor.Requirements) {
		req := factor.Requirements[name]
		node.Children = append(node.Children, &Element{
			ID:       RequirementReference(areaPath, name),
			Kind:     KindRequirement,
			Label:    labelOrKey(req.Title, name),
			ParentID: node.ID,
		})
	}
	return node
}

// Flatten returns the projection as a flat pre-order slice (parent before
// children), the documented enumeration order for list output.
func Flatten(root *Element) []*Element {
	var out []*Element
	var walk func(*Element)
	walk = func(e *Element) {
		out = append(out, e)
		for _, child := range e.Children {
			walk(child)
		}
	}
	walk(root)
	return out
}

// Find returns the element whose canonical ID equals id, or nil.
func Find(root *Element, id string) *Element {
	for _, e := range Flatten(root) {
		if e.ID == id {
			return e
		}
	}
	return nil
}

func areaLabel(title, fallback string) string {
	return labelOrKey(title, fallback)
}

func labelOrKey(title, key string) string {
	if title != "" {
		return title
	}
	return key
}

func appendPath(path AreaPath, name string) AreaPath {
	out := make(AreaPath, 0, len(path)+1)
	out = append(out, path...)
	return append(out, name)
}

func appendFactorPath(path FactorPath, name string) FactorPath {
	out := make(FactorPath, 0, len(path)+1)
	out = append(out, path...)
	return append(out, name)
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
