package lint

import (
	"strconv"

	"github.com/qualitymd/quality.md/internal/spec"
	"gopkg.in/yaml.v3"
)

func (s *runState) run() {
	if s.doc.Frontmatter.Kind != yaml.MappingNode {
		s.add(RuleInvalidFrontmatter, "The frontmatter is not a model mapping; a QUALITY.md frontmatter block must be a map of model properties.", s.loc(s.doc.Frontmatter, nil, "frontmatter"), nil)
		return
	}
	s.root = &targetRef{node: s.doc.Frontmatter, path: []PathSegment{}}
	s.checkRoot()
	s.checkRatingScale()
	s.walkModel()
	s.checkEmptyModel()
	s.checkTargets(s.root)
	s.checkFactors(s.root)
	s.checkRequirements(s.root)
	s.sort()
}

func (s *runState) checkRoot() {
	allowed := map[string]bool{
		"title": true, "ratingScale": true, "factors": true,
		"requirements": true, "targets": true, "source": true,
	}
	for key, value := range spec.MapEntries(s.doc.Frontmatter) {
		if !allowed[key.Value] {
			s.invalid(key, appendPath(nil, key.Value), key.Value, "The frontmatter declares an unknown root key `"+key.Value+"`; a QUALITY.md model may only use the specified model properties.")
			continue
		}
		switch key.Value {
		case "title", "source":
			s.checkOptionalScalar(s.doc.Frontmatter, key, value, []PathSegment{key.Value}, key.Value)
		case "factors", "requirements", "targets":
			if isEmpty(value) {
				s.emptyProperty(s.doc.Frontmatter, key, []PathSegment{key.Value}, key.Value)
			} else if value.Kind != yaml.MappingNode {
				s.invalid(key, []PathSegment{key.Value}, key.Value, "The `"+key.Value+"` property has the wrong YAML shape; it must be a map.")
			}
		case "ratingScale":
			if !isEmpty(value) && value.Kind != yaml.SequenceNode {
				s.invalid(key, []PathSegment{"ratingScale"}, "ratingScale", "The `ratingScale` property has the wrong YAML shape; it must be a list of rating levels.")
			}
		}
	}
	if _, value, _ := spec.MapEntry(s.doc.Frontmatter, "title"); value == nil {
		s.add(RuleMissingTitle, "The model root declares no `title`; a title is recommended for readable reports.", s.locForMissing([]PathSegment{"title"}, "title"), nil)
	}
}

func (s *runState) checkRatingScale() {
	key, scale, _ := spec.MapEntry(s.doc.Frontmatter, "ratingScale")
	if scale == nil || isEmpty(scale) {
		s.add(RuleMissingRatingScale, "The model root declares no `ratingScale`; a QUALITY.md model requires one rating scale.", s.locForMissing([]PathSegment{"ratingScale"}, "ratingScale"), nil)
		return
	}
	if scale.Kind != yaml.SequenceNode {
		return
	}
	if len(scale.Content) < 2 {
		s.add(RuleTooFewLevels, "The `ratingScale` has fewer than two levels; a QUALITY.md rating scale requires at least two levels.", s.loc(key, []PathSegment{"ratingScale"}, "ratingScale"), nil)
	}
	seen := map[string]Location{}
	for i, level := range scale.Content {
		levelPath := []PathSegment{"ratingScale", i}
		if level.Kind != yaml.MappingNode {
			s.invalid(level, levelPath, "ratingScale["+strconv.Itoa(i)+"]", "A rating level has the wrong YAML shape; each `ratingScale` item must be a map.")
			continue
		}
		s.checkRatingLevel(level, levelPath, i, seen)
	}
}

func (s *runState) checkRatingLevel(level *yaml.Node, path []PathSegment, index int, seen map[string]Location) {
	allowed := map[string]bool{"level": true, "title": true, "description": true, "criterion": true}
	for key, value := range spec.MapEntries(level) {
		if !allowed[key.Value] {
			s.invalid(key, appendPath(path, key.Value), label(appendPath(path, key.Value)), "A rating level declares an unknown key `"+key.Value+"`; rating levels may only use `level`, `title`, `description`, and `criterion`.")
			continue
		}
		switch key.Value {
		case "title", "description":
			s.checkOptionalScalar(level, key, value, appendPath(path, key.Value), label(appendPath(path, key.Value)))
		case "level", "criterion":
			if value.Kind != yaml.ScalarNode && !isEmpty(value) {
				s.invalid(key, appendPath(path, key.Value), label(appendPath(path, key.Value)), "The `"+key.Value+"` property has the wrong YAML shape; it must be a scalar.")
			}
		}
	}
	if _, value, _ := spec.MapEntry(level, "level"); value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingLevelName, "A rating level declares no `level` name; each rating level requires a non-empty `level`.", s.locForNodeOrMissing(value, appendPath(path, "level"), "ratingScale["+strconv.Itoa(index)+"].level"), nil)
	} else {
		name := value.Value
		if prior, ok := seen[name]; ok {
			s.add(RuleDuplicateLevel, "The rating level `"+name+"` is duplicated; each `level` name must be unique within `ratingScale`.", s.loc(value, appendPath(path, "level"), "ratingScale["+strconv.Itoa(index)+"].level"), nil)
			s.add(RuleDuplicateLevel, "The rating level `"+name+"` is duplicated; each `level` name must be unique within `ratingScale`.", prior, nil)
		} else {
			seen[name] = s.loc(value, appendPath(path, "level"), "ratingScale["+strconv.Itoa(index)+"].level")
			s.levels[name] = true
		}
	}
	if _, value, _ := spec.MapEntry(level, "criterion"); value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingCriterion, "A rating level declares no `criterion`; each rating level requires a non-empty criterion.", s.locForNodeOrMissing(value, appendPath(path, "criterion"), "ratingScale["+strconv.Itoa(index)+"].criterion"), nil)
	}
	if _, value, _ := spec.MapEntry(level, "description"); value == nil {
		s.add(RuleMissingLevelDescription, "A rating level declares no `description`; a description is recommended for each level.", s.locForMissing(appendPath(path, "description"), "ratingScale["+strconv.Itoa(index)+"].description"), nil)
	}
}

func (s *runState) checkEmptyModel() {
	if len(s.root.factors) == 0 && len(s.root.requirements) == 0 && len(s.root.targets) == 0 {
		s.add(RuleEmptyModel, "The model root supplies no entries under `factors`, `requirements`, or `targets`; a QUALITY.md model requires model content.", s.loc(s.doc.Frontmatter, []PathSegment{}, "frontmatter"), nil)
	}
}

func (s *runState) checkTargetShape(target *targetRef) {
	allowed := map[string]bool{"factors": true, "requirements": true, "targets": true, "source": true, "title": true, "ratingScale": true}
	for key, value := range spec.MapEntries(target.node) {
		path := appendPath(target.path, key.Value)
		switch {
		case key.Value == "title" || key.Value == "ratingScale":
			s.add(RuleMisplacedRootKey, "The target `"+target.name+"` declares `"+key.Value+"`; `"+key.Value+"` is only valid on the model root.", s.loc(key, path, label(path)), nil)
		case !allowed[key.Value]:
			s.invalid(key, path, label(path), "The target `"+target.name+"` declares an unknown key `"+key.Value+"`; targets may only use target properties.")
		case key.Value == "source":
			s.checkOptionalScalar(target.node, key, value, path, label(path))
		case key.Value == "factors" || key.Value == "requirements" || key.Value == "targets":
			if isEmpty(value) {
				s.emptyProperty(target.node, key, path, label(path))
			} else if value.Kind != yaml.MappingNode {
				s.invalid(key, path, label(path), "The `"+key.Value+"` property has the wrong YAML shape; it must be a map.")
			}
		}
	}
}

func (s *runState) checkFactorShape(factor *factorRef) {
	allowed := map[string]bool{"description": true, "factors": true, "requirements": true}
	for key, value := range spec.MapEntries(factor.node) {
		path := appendPath(factor.path, key.Value)
		switch {
		case !allowed[key.Value]:
			s.invalid(key, path, label(path), "The factor `"+factor.name+"` declares an unknown key `"+key.Value+"`; factors may only use factor properties.")
		case key.Value == "description":
			s.checkOptionalScalar(factor.node, key, value, path, label(path))
		case key.Value == "factors" || key.Value == "requirements":
			if isEmpty(value) {
				s.emptyProperty(factor.node, key, path, label(path))
			} else if value.Kind != yaml.MappingNode {
				s.invalid(key, path, label(path), "The `"+key.Value+"` property has the wrong YAML shape; it must be a map.")
			}
		}
	}
	if _, value, _ := spec.MapEntry(factor.node, "description"); value == nil {
		s.add(RuleMissingFactorDescription, "The factor `"+factor.name+"` declares no `description`; a description is recommended for each factor.", s.locForMissing(appendPath(factor.path, "description"), label(appendPath(factor.path, "description"))), nil)
	}
}

func (s *runState) checkRequirementShape(req *requirementRef) {
	allowed := map[string]bool{"assessment": true, "factors": true, "ratings": true}
	for key, value := range spec.MapEntries(req.node) {
		path := appendPath(req.path, key.Value)
		switch {
		case !allowed[key.Value]:
			s.invalid(key, path, label(path), "The requirement `"+req.statement+"` declares an unknown key `"+key.Value+"`; requirements may only use requirement properties.")
		case key.Value == "assessment":
			if value.Kind != yaml.ScalarNode || isEmpty(value) {
				s.add(RuleInvalidAssessment, "The requirement `"+req.statement+"` has an invalid `assessment`; a requirement must declare exactly one non-empty scalar assessment.", s.loc(value, path, label(path)), nil)
			}
		case key.Value == "factors":
			if isEmpty(value) {
				s.emptyProperty(req.node, key, path, label(path))
			} else if value.Kind != yaml.SequenceNode {
				s.invalid(key, path, label(path), "The requirement `"+req.statement+"` has the wrong `factors` shape; secondary factors must be a list.")
			} else {
				for i, item := range value.Content {
					if item.Kind != yaml.ScalarNode || isEmpty(item) {
						s.invalid(item, appendPath(path, i), label(appendPath(path, i)), "The requirement `"+req.statement+"` has a secondary factor with the wrong YAML shape; each secondary factor must be a non-empty scalar.")
					}
				}
			}
		case key.Value == "ratings":
			if isEmpty(value) {
				s.emptyProperty(req.node, key, path, label(path))
			} else if value.Kind != yaml.MappingNode {
				s.invalid(key, path, label(path), "The requirement `"+req.statement+"` has the wrong `ratings` shape; ratings overrides must be a map.")
			} else {
				for ratingKey, ratingValue := range spec.MapEntries(value) {
					if ratingValue.Kind != yaml.ScalarNode || isEmpty(ratingValue) {
						s.invalid(ratingValue, appendPath(path, ratingKey.Value), label(appendPath(path, ratingKey.Value)), "The rating override `"+ratingKey.Value+"` has the wrong YAML shape; override criteria must be non-empty scalars.")
					}
				}
			}
		}
	}
	if _, value, _ := spec.MapEntry(req.node, "assessment"); value == nil {
		s.add(RuleInvalidAssessment, "The requirement `"+req.statement+"` has no `assessment`; a requirement must declare exactly one non-empty scalar assessment.", s.locForMissing(appendPath(req.path, "assessment"), label(appendPath(req.path, "assessment"))), nil)
	}
}

func (s *runState) checkTargets(target *targetRef) bool {
	hasRequirements := len(target.requirements) > 0
	for _, factor := range target.factors {
		if factorHasRequirements(factor) {
			hasRequirements = true
		}
	}
	for _, child := range target.targets {
		if s.checkTargets(child) {
			hasRequirements = true
		}
	}
	if target.parent != nil && !hasRequirements {
		s.add(RuleEmptyTarget, "The target `"+target.name+"` reaches no requirements in its subtree; each target should lead to at least one requirement.", s.loc(target.node, target.path, label(target.path)), nil)
	}
	return hasRequirements
}

func (s *runState) checkFactors(target *targetRef) {
	for _, factor := range target.factors {
		s.checkFactor(factor)
	}
	for _, child := range target.targets {
		s.checkFactors(child)
	}
}

func (s *runState) checkFactor(factor *factorRef) bool {
	has := len(factor.requirements) > 0
	for _, child := range factor.factors {
		if s.checkFactor(child) {
			has = true
		}
	}
	for _, req := range allRequirements(factor.target) {
		for _, resolved := range s.secondaryFactors(req) {
			if resolved == factor {
				has = true
			}
		}
	}
	if !has {
		s.add(RuleEmptyFactor, "The factor `"+factor.name+"` leads to no requirements; each factor should be tied to at least one requirement.", s.loc(factor.node, factor.path, label(factor.path)), nil)
	}
	return has
}

func (s *runState) checkRequirements(target *targetRef) {
	for _, req := range target.requirements {
		s.checkRequirementRefs(req)
	}
	for _, factor := range target.factors {
		s.checkFactorRequirements(factor)
	}
	for _, child := range target.targets {
		s.checkRequirements(child)
	}
}

func (s *runState) checkFactorRequirements(factor *factorRef) {
	for _, req := range factor.requirements {
		s.checkRequirementRefs(req)
	}
	for _, child := range factor.factors {
		s.checkFactorRequirements(child)
	}
}

func (s *runState) checkRequirementRefs(req *requirementRef) {
	if _, ratings, _ := spec.MapEntry(req.node, "ratings"); ratings != nil && ratings.Kind == yaml.MappingNode {
		for key := range spec.MapEntries(ratings) {
			if !s.levels[key.Value] {
				path := appendPath(req.path, "ratings", key.Value)
				s.add(RuleUnknownRatingKey, "The requirement `"+req.statement+"` has a `ratings` override for unknown level `"+key.Value+"`; override keys must name a rating-scale level.", s.loc(key, path, label(path)), nil)
			}
		}
	}
	if _, factors, _ := spec.MapEntry(req.node, "factors"); factors != nil && factors.Kind == yaml.SequenceNode {
		for i, item := range factors.Content {
			if item.Kind != yaml.ScalarNode || isEmpty(item) {
				continue
			}
			if s.resolveFactor(req.target, item.Value) == nil {
				path := appendPath(req.path, "factors", i)
				s.add(RuleUnknownFactor, "The requirement `"+req.statement+"` names unknown secondary factor `"+item.Value+"`; secondary factors must resolve on the declaring target or an ancestor.", s.loc(item, path, label(path)), nil)
			}
		}
	}
}
