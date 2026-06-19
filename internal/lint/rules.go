package lint

import (
	"strconv"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	qschema "github.com/qualitymd/quality.md/internal/schema"
	"gopkg.in/yaml.v3"
)

type schemaContext struct {
	targetName  string
	factorName  string
	requirement string
}

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
	s.checkSchemaProperties(qschema.Model, s.doc.Frontmatter, nil, schemaContext{})
	if _, value, _ := document.MapEntry(s.doc.Frontmatter, qschema.PropertyTitle); value == nil {
		s.add(RuleMissingTitle, "The model root declares no `title`; a title is recommended for readable reports.", s.locForMissing([]PathSegment{"title"}, "title"), nil)
	}
}

func (s *runState) checkSchemaProperties(node qschema.Node, owner *yaml.Node, base []PathSegment, context schemaContext) {
	if owner.Kind != yaml.MappingNode {
		return
	}
	for key, value := range document.MapEntries(owner) {
		path := appendPath(base, key.Value)
		locationLabel := label(path)
		property, ok := node.Property(key.Value)
		if !ok {
			if node.Kind == qschema.TargetKind && qschema.Model.HasProperty(key.Value) {
				s.add(RuleMisplacedRootKey, "The target `"+context.targetName+"` declares `"+key.Value+"`; `"+key.Value+"` is only valid on the model root.", s.loc(key, path, locationLabel), nil)
				continue
			}
			s.invalid(key, path, locationLabel, unknownPropertyMessage(node, key.Value, context))
			continue
		}
		s.checkSchemaProperty(node, property, owner, key, value, path, locationLabel, context)
	}
}

func (s *runState) checkSchemaProperty(node qschema.Node, property qschema.Property, owner, key, value *yaml.Node, path []PathSegment, locationLabel string, context schemaContext) {
	if node.Kind == qschema.RequirementKind && property.Name == qschema.PropertyAssessment {
		if value.Kind != yaml.ScalarNode || isEmpty(value) {
			s.add(RuleInvalidAssessment, "The requirement `"+context.requirement+"` has an invalid `assessment`; a requirement must declare exactly one non-empty scalar assessment.", s.loc(value, path, locationLabel), nil)
		}
		return
	}

	if property.Presence != qschema.RequiredPresence && isEmpty(value) {
		s.emptyProperty(owner, key, path, locationLabel)
		return
	}

	switch property.Shape {
	case qschema.ScalarShape:
		if property.Presence == qschema.RequiredPresence {
			if value.Kind != yaml.ScalarNode && !isEmpty(value) {
				s.invalid(key, path, locationLabel, "The `"+property.Name+"` property has the wrong YAML shape; it must be a scalar.")
			}
			return
		}
		s.checkOptionalScalar(owner, key, value, path, locationLabel)
	case qschema.MapShape:
		if isEmpty(value) {
			return
		}
		if value.Kind != yaml.MappingNode {
			s.invalid(key, path, locationLabel, wrongShapeMessage(property, context))
			return
		}
		if property.ValueShape == qschema.ScalarShape {
			s.checkScalarMapValues(value, path)
		}
	case qschema.SequenceShape:
		if isEmpty(value) {
			return
		}
		if value.Kind != yaml.SequenceNode {
			s.invalid(key, path, locationLabel, wrongShapeMessage(property, context))
			return
		}
		if property.ElementShape == qschema.ScalarShape {
			s.checkScalarSequenceItems(value, path, context)
		}
	}
}

func (s *runState) checkScalarMapValues(value *yaml.Node, path []PathSegment) {
	for ratingKey, ratingValue := range document.MapEntries(value) {
		if ratingValue.Kind != yaml.ScalarNode || isEmpty(ratingValue) {
			s.invalid(ratingValue, appendPath(path, ratingKey.Value), label(appendPath(path, ratingKey.Value)), "The rating override `"+ratingKey.Value+"` has the wrong YAML shape; override criteria must be non-empty scalars.")
		}
	}
}

func (s *runState) checkScalarSequenceItems(value *yaml.Node, path []PathSegment, context schemaContext) {
	for i, item := range value.Content {
		if item.Kind != yaml.ScalarNode || isEmpty(item) {
			s.invalid(item, appendPath(path, i), label(appendPath(path, i)), "The requirement `"+context.requirement+"` has a factor reference with the wrong YAML shape; each factor reference must be a non-empty scalar.")
		}
	}
}

func unknownPropertyMessage(node qschema.Node, key string, context schemaContext) string {
	switch node.Kind {
	case qschema.ModelKind:
		return "The frontmatter declares an unknown root key `" + key + "`; a QUALITY.md model may only use the specified model properties."
	case qschema.TargetKind:
		return "The target `" + context.targetName + "` declares an unknown key `" + key + "`; targets may only use target properties."
	case qschema.FactorKind:
		return "The factor `" + context.factorName + "` declares an unknown key `" + key + "`; factors may only use factor properties."
	case qschema.RequirementKind:
		return "The requirement `" + context.requirement + "` declares an unknown key `" + key + "`; requirements may only use requirement properties."
	case qschema.RatingLevelKind:
		return "A rating level declares an unknown key `" + key + "`; rating levels may only use " + formatSchemaKeys(qschema.RatingLevel.PropertyNames()) + "."
	default:
		return "Unknown property `" + key + "`."
	}
}

func formatSchemaKeys(keys []string) string {
	quoted := make([]string, 0, len(keys))
	for _, key := range keys {
		quoted = append(quoted, "`"+key+"`")
	}
	if len(quoted) == 0 {
		return "no properties"
	}
	if len(quoted) == 1 {
		return quoted[0]
	}
	return strings.Join(quoted[:len(quoted)-1], ", ") + ", and " + quoted[len(quoted)-1]
}

func wrongShapeMessage(property qschema.Property, context schemaContext) string {
	switch property.Name {
	case qschema.PropertyRatingScale:
		return "The `ratingScale` property has the wrong YAML shape; it must be a list of rating levels."
	case qschema.PropertyFactors:
		if context.requirement != "" {
			return "The requirement `" + context.requirement + "` has the wrong `factors` shape; factor references must be a sequence."
		}
	case qschema.PropertyRatings:
		return "The requirement `" + context.requirement + "` has the wrong `ratings` shape; ratings overrides must be a map."
	}
	switch property.Shape {
	case qschema.MapShape:
		return "The `" + property.Name + "` property has the wrong YAML shape; it must be a map."
	case qschema.SequenceShape:
		return "The `" + property.Name + "` property has the wrong YAML shape; it must be a list."
	case qschema.ScalarShape:
		return "The `" + property.Name + "` property has the wrong YAML shape; it must be a scalar."
	default:
		return "The `" + property.Name + "` property has the wrong YAML shape."
	}
}

func (s *runState) checkRatingScale() {
	key, scale, _ := document.MapEntry(s.doc.Frontmatter, qschema.PropertyRatingScale)
	property, _ := qschema.Model.Property(qschema.PropertyRatingScale)
	if scale == nil || isEmpty(scale) {
		s.add(RuleMissingRatingScale, "The model root declares no `ratingScale`; a QUALITY.md model requires one rating scale.", s.locForMissing([]PathSegment{"ratingScale"}, "ratingScale"), nil)
		return
	}
	if scale.Kind != yaml.SequenceNode {
		return
	}
	if len(scale.Content) < property.MinItems {
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
	s.checkSchemaProperties(qschema.RatingLevel, level, path, schemaContext{})
	if _, value, _ := document.MapEntry(level, qschema.PropertyLevel); value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
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
	if _, value, _ := document.MapEntry(level, qschema.PropertyCriterion); value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingCriterion, "A rating level declares no `criterion`; each rating level requires a non-empty criterion.", s.locForNodeOrMissing(value, appendPath(path, "criterion"), "ratingScale["+strconv.Itoa(index)+"].criterion"), nil)
	}
	if _, value, _ := document.MapEntry(level, qschema.PropertyDescription); value == nil {
		s.add(RuleMissingLevelDescription, "A rating level declares no `description`; a description is recommended for each level.", s.locForMissing(appendPath(path, "description"), "ratingScale["+strconv.Itoa(index)+"].description"), nil)
	}
}

func (s *runState) checkEmptyModel() {
	for _, group := range qschema.Model.RequiredAny {
		if group.Name == "model-content" && !s.rootHasAny(group.Properties) {
			s.add(RuleEmptyModel, "The model root supplies no entries under `factors`, `requirements`, or `targets`; a QUALITY.md model requires model content.", s.loc(s.doc.Frontmatter, []PathSegment{}, "frontmatter"), nil)
		}
	}
}

func (s *runState) rootHasAny(properties []string) bool {
	for _, property := range properties {
		switch property {
		case qschema.PropertyFactors:
			if len(s.root.factors) > 0 {
				return true
			}
		case qschema.PropertyRequirements:
			if len(s.root.requirements) > 0 {
				return true
			}
		case qschema.PropertyTargets:
			if len(s.root.targets) > 0 {
				return true
			}
		}
	}
	return false
}

func (s *runState) checkTargetShape(target *targetRef) {
	s.checkSchemaProperties(qschema.Target, target.node, target.path, schemaContext{targetName: target.name})
}

func (s *runState) checkFactorShape(factor *factorRef) {
	s.checkSchemaProperties(qschema.Factor, factor.node, factor.path, schemaContext{factorName: factor.name})
	if _, value, _ := document.MapEntry(factor.node, qschema.PropertyDescription); value == nil {
		s.add(RuleMissingFactorDescription, "The factor `"+factor.name+"` declares no `description`; a description is recommended for each factor.", s.locForMissing(appendPath(factor.path, "description"), label(appendPath(factor.path, "description"))), nil)
	}
}

func (s *runState) checkRequirementShape(req *requirementRef) {
	s.checkSchemaProperties(qschema.Requirement, req.node, req.path, schemaContext{requirement: req.statement})
	if _, value, _ := document.MapEntry(req.node, qschema.PropertyAssessment); value == nil {
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
		for _, resolved := range s.referencedFactors(req) {
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
	if req.factor == nil && !requirementReferencesFactor(req) {
		path := appendPath(req.path, "factors")
		s.add(RuleMissingFactorReference, "The requirement `"+req.statement+"` references no quality factor; place it under a factor or add one or more factor references under `factors`.", s.locForMissing(path, label(path)), nil)
	}
	if _, ratings, _ := document.MapEntry(req.node, qschema.PropertyRatings); ratings != nil && ratings.Kind == yaml.MappingNode {
		for key := range document.MapEntries(ratings) {
			if !s.levels[key.Value] {
				path := appendPath(req.path, "ratings", key.Value)
				s.add(RuleUnknownRatingKey, "The requirement `"+req.statement+"` has a `ratings` override for unknown level `"+key.Value+"`; override keys must name a rating-scale level.", s.loc(key, path, label(path)), nil)
			}
		}
	}
	if _, factors, _ := document.MapEntry(req.node, qschema.PropertyFactors); factors != nil && factors.Kind == yaml.SequenceNode {
		for i, item := range factors.Content {
			if item.Kind != yaml.ScalarNode || isEmpty(item) {
				continue
			}
			if s.resolveFactor(req.target, item.Value) == nil {
				path := appendPath(req.path, "factors", i)
				s.add(RuleUnknownFactor, "The requirement `"+req.statement+"` references unknown factor `"+item.Value+"`; factor references must resolve on the declaring target or an ancestor.", s.loc(item, path, label(path)), nil)
			}
		}
	}
}
