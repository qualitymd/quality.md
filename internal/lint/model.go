package lint

import (
	"github.com/qualitymd/quality.md/internal/document"
	qschema "github.com/qualitymd/quality.md/internal/schema"
	"gopkg.in/yaml.v3"
)

type areaRef struct {
	name         string
	node         *yaml.Node
	path         []PathSegment
	parent       *areaRef
	factors      []*factorRef
	requirements []*requirementRef
	areas        []*areaRef
}

type factorRef struct {
	name         string
	node         *yaml.Node
	path         []PathSegment
	area         *areaRef
	parent       *factorRef
	factors      []*factorRef
	requirements []*requirementRef
}

type requirementRef struct {
	statement string
	node      *yaml.Node
	path      []PathSegment
	area      *areaRef
	factor    *factorRef
}

func (s *runState) walkModel() {
	s.root.factors = s.walkFactors(s.root, nil, s.doc.Frontmatter, []PathSegment{})
	s.root.requirements = s.walkRequirements(s.root, nil, s.doc.Frontmatter, []PathSegment{})
	s.root.areas = s.walkAreas(s.root, s.doc.Frontmatter, []PathSegment{})
}

func (s *runState) walkAreas(parent *areaRef, node *yaml.Node, base []PathSegment) []*areaRef {
	_, areas, _ := document.MapEntry(node, qschema.PropertyAreas)
	if areas == nil || areas.Kind != yaml.MappingNode {
		return nil
	}
	var out []*areaRef
	for key, value := range document.MapEntries(areas) {
		path := appendPath(base, qschema.PropertyAreas, key.Value)
		if key.Kind != yaml.ScalarNode || !validModelName(key.Value) {
			s.add(RuleInvalidAreaName, "The area name `"+key.Value+"` is invalid; Area names must match "+qschema.ModelNamePattern+".", s.loc(key, path, label(path)), nil)
		}
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, label(path), "The area `"+key.Value+"` has the wrong YAML shape; each area must be a map.")
			continue
		}
		area := &areaRef{name: key.Value, node: value, path: path, parent: parent}
		s.checkAreaShape(area)
		area.factors = s.walkFactors(area, nil, value, path)
		area.requirements = s.walkRequirements(area, nil, value, path)
		area.areas = s.walkAreas(area, value, path)
		out = append(out, area)
	}
	return out
}

func (s *runState) walkFactors(area *areaRef, parent *factorRef, node *yaml.Node, base []PathSegment) []*factorRef {
	_, factors, _ := document.MapEntry(node, qschema.PropertyFactors)
	if factors == nil || factors.Kind != yaml.MappingNode {
		return nil
	}
	var out []*factorRef
	for key, value := range document.MapEntries(factors) {
		path := appendPath(base, qschema.PropertyFactors, key.Value)
		if key.Kind != yaml.ScalarNode || !validModelName(key.Value) {
			s.add(RuleInvalidFactorName, "The factor name `"+key.Value+"` is invalid; Factor names must match "+qschema.ModelNamePattern+".", s.loc(key, path, label(path)), nil)
		}
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, label(path), "The factor `"+key.Value+"` has the wrong YAML shape; each factor must be a map.")
			continue
		}
		factor := &factorRef{name: key.Value, node: value, path: path, area: area, parent: parent}
		s.checkFactorShape(factor)
		factor.factors = s.walkFactors(area, factor, value, path)
		factor.requirements = s.walkRequirements(area, factor, value, path)
		out = append(out, factor)
	}
	return out
}

func (s *runState) walkRequirements(area *areaRef, factor *factorRef, node *yaml.Node, base []PathSegment) []*requirementRef {
	_, requirements, _ := document.MapEntry(node, qschema.PropertyRequirements)
	if requirements == nil || requirements.Kind != yaml.MappingNode {
		return nil
	}
	var out []*requirementRef
	for key, value := range document.MapEntries(requirements) {
		path := appendPath(base, qschema.PropertyRequirements, key.Value)
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, label(path), "The requirement `"+key.Value+"` has the wrong YAML shape; each requirement must be a map.")
			continue
		}
		req := &requirementRef{statement: key.Value, node: value, path: path, area: area, factor: factor}
		s.checkRequirementShape(req)
		out = append(out, req)
	}
	return out
}

func (s *runState) referencedFactors(req *requirementRef) []*factorRef {
	_, factors, _ := document.MapEntry(req.node, qschema.PropertyFactors)
	if factors == nil || factors.Kind != yaml.SequenceNode {
		return nil
	}
	var out []*factorRef
	for _, item := range factors.Content {
		if item.Kind != yaml.ScalarNode || isEmpty(item) {
			continue
		}
		if factor := s.resolveFactor(req.area, item.Value); factor != nil {
			out = append(out, factor)
		}
	}
	return out
}

func requirementReferencesFactor(req *requirementRef) bool {
	_, factors, _ := document.MapEntry(req.node, qschema.PropertyFactors)
	if factors == nil || factors.Kind != yaml.SequenceNode {
		return false
	}
	for _, item := range factors.Content {
		if item.Kind == yaml.ScalarNode && !isEmpty(item) {
			return true
		}
	}
	return false
}

func (s *runState) resolveFactor(area *areaRef, name string) *factorRef {
	for current := area; current != nil; current = current.parent {
		if found := findFactor(current.factors, name); found != nil {
			return found
		}
	}
	return nil
}

func findFactor(factors []*factorRef, name string) *factorRef {
	for _, factor := range factors {
		if factor.name == name {
			return factor
		}
		if found := findFactor(factor.factors, name); found != nil {
			return found
		}
	}
	return nil
}

func allRequirements(area *areaRef) []*requirementRef {
	var out []*requirementRef
	var walkArea func(*areaRef)
	var walkFactor func(*factorRef)
	walkFactor = func(factor *factorRef) {
		out = append(out, factor.requirements...)
		for _, child := range factor.factors {
			walkFactor(child)
		}
	}
	walkArea = func(area *areaRef) {
		out = append(out, area.requirements...)
		for _, factor := range area.factors {
			walkFactor(factor)
		}
		for _, child := range area.areas {
			walkArea(child)
		}
	}
	walkArea(area)
	return out
}

func factorHasRequirements(factor *factorRef) bool {
	if len(factor.requirements) > 0 {
		return true
	}
	for _, child := range factor.factors {
		if factorHasRequirements(child) {
			return true
		}
	}
	return false
}
