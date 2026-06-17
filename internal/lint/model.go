package lint

import (
	"github.com/qualitymd/quality.md/internal/document"
	qschema "github.com/qualitymd/quality.md/internal/schema"
	"gopkg.in/yaml.v3"
)

type targetRef struct {
	name         string
	node         *yaml.Node
	path         []PathSegment
	parent       *targetRef
	factors      []*factorRef
	requirements []*requirementRef
	targets      []*targetRef
}

type factorRef struct {
	name         string
	node         *yaml.Node
	path         []PathSegment
	target       *targetRef
	parent       *factorRef
	factors      []*factorRef
	requirements []*requirementRef
}

type requirementRef struct {
	statement string
	node      *yaml.Node
	path      []PathSegment
	target    *targetRef
	factor    *factorRef
}

func (s *runState) walkModel() {
	s.root.factors = s.walkFactors(s.root, nil, s.doc.Frontmatter, []PathSegment{})
	s.root.requirements = s.walkRequirements(s.root, nil, s.doc.Frontmatter, []PathSegment{})
	s.root.targets = s.walkTargets(s.root, s.doc.Frontmatter, []PathSegment{})
}

func (s *runState) walkTargets(parent *targetRef, node *yaml.Node, base []PathSegment) []*targetRef {
	_, targets, _ := document.MapEntry(node, qschema.PropertyTargets)
	if targets == nil || targets.Kind != yaml.MappingNode {
		return nil
	}
	var out []*targetRef
	for key, value := range document.MapEntries(targets) {
		path := appendPath(base, qschema.PropertyTargets, key.Value)
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, label(path), "The target `"+key.Value+"` has the wrong YAML shape; each target must be a map.")
			continue
		}
		target := &targetRef{name: key.Value, node: value, path: path, parent: parent}
		s.checkTargetShape(target)
		target.factors = s.walkFactors(target, nil, value, path)
		target.requirements = s.walkRequirements(target, nil, value, path)
		target.targets = s.walkTargets(target, value, path)
		out = append(out, target)
	}
	return out
}

func (s *runState) walkFactors(target *targetRef, parent *factorRef, node *yaml.Node, base []PathSegment) []*factorRef {
	_, factors, _ := document.MapEntry(node, qschema.PropertyFactors)
	if factors == nil || factors.Kind != yaml.MappingNode {
		return nil
	}
	var out []*factorRef
	for key, value := range document.MapEntries(factors) {
		path := appendPath(base, qschema.PropertyFactors, key.Value)
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, label(path), "The factor `"+key.Value+"` has the wrong YAML shape; each factor must be a map.")
			continue
		}
		factor := &factorRef{name: key.Value, node: value, path: path, target: target, parent: parent}
		s.checkFactorShape(factor)
		factor.factors = s.walkFactors(target, factor, value, path)
		factor.requirements = s.walkRequirements(target, factor, value, path)
		out = append(out, factor)
	}
	return out
}

func (s *runState) walkRequirements(target *targetRef, factor *factorRef, node *yaml.Node, base []PathSegment) []*requirementRef {
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
		req := &requirementRef{statement: key.Value, node: value, path: path, target: target, factor: factor}
		s.checkRequirementShape(req)
		out = append(out, req)
	}
	return out
}

func (s *runState) secondaryFactors(req *requirementRef) []*factorRef {
	_, factors, _ := document.MapEntry(req.node, qschema.PropertyFactors)
	if factors == nil || factors.Kind != yaml.SequenceNode {
		return nil
	}
	var out []*factorRef
	for _, item := range factors.Content {
		if item.Kind != yaml.ScalarNode || isEmpty(item) {
			continue
		}
		if factor := s.resolveFactor(req.target, item.Value); factor != nil {
			out = append(out, factor)
		}
	}
	return out
}

func (s *runState) resolveFactor(target *targetRef, name string) *factorRef {
	for current := target; current != nil; current = current.parent {
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

func allRequirements(target *targetRef) []*requirementRef {
	var out []*requirementRef
	var walkTarget func(*targetRef)
	var walkFactor func(*factorRef)
	walkFactor = func(factor *factorRef) {
		out = append(out, factor.requirements...)
		for _, child := range factor.factors {
			walkFactor(child)
		}
	}
	walkTarget = func(target *targetRef) {
		out = append(out, target.requirements...)
		for _, factor := range target.factors {
			walkFactor(factor)
		}
		for _, child := range target.targets {
			walkTarget(child)
		}
	}
	walkTarget(target)
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
