package lint

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	qschema "github.com/qualitymd/quality.md/internal/schema"
	"github.com/qualitymd/quality.md/internal/workspace"
	"gopkg.in/yaml.v3"
)

var modelNamePattern = regexp.MustCompile(qschema.ModelNamePattern)

type schemaContext struct {
	areaName    string
	factorName  string
	requirement string
}

func (s *runState) run() {
	if s.doc.Frontmatter.Kind != yaml.MappingNode {
		s.add(RuleInvalidFrontmatter, "The frontmatter is not a model mapping; a QUALITY.md frontmatter block must be a map of model properties.", s.loc(s.doc.Frontmatter, nil, "frontmatter"), nil)
		return
	}
	s.root = &areaRef{node: s.doc.Frontmatter, path: []PathSegment{}}
	s.checkRoot()
	s.checkRatingScale()
	s.walkModel()
	s.checkEmptyModel()
	s.checkAreas(s.root)
	s.checkFactors(s.root)
	s.checkRequirements(s.root)
	s.sort()
}

func (s *runState) checkRoot() {
	s.checkSchemaProperties(qschema.Model, s.doc.Frontmatter, nil, schemaContext{})
	s.checkRootConfig()
	s.checkRequiredTitle(s.doc.Frontmatter, []PathSegment{}, "model root", "The model root declares no `title`; a model title is required for human-facing display.")
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
			if node.Kind == qschema.ModelKind && s.options.Rules.UnknownKey.AllowedRootKeys[key.Value] {
				continue
			}
			if node.Kind == qschema.AreaKind && qschema.Model.HasProperty(key.Value) {
				s.add(RuleMisplacedRootKey, "The area `"+context.areaName+"` declares `"+key.Value+"`; `"+key.Value+"` is only valid on the model root.", s.loc(key, path, locationLabel), nil)
				continue
			}
			s.invalid(key, path, locationLabel, unknownPropertyMessage(node, key.Value, context))
			continue
		}
		s.checkSchemaProperty(node, property, owner, key, value, path, locationLabel, context)
	}
}

func (s *runState) checkRootConfig() {
	key, value, _ := document.MapEntry(s.doc.Frontmatter, workspace.FrontmatterConfigField)
	if key == nil {
		return
	}
	path := []PathSegment{workspace.FrontmatterConfigField}
	if value.Kind != yaml.ScalarNode || strings.TrimSpace(value.Value) == "" {
		s.add(RuleInvalidConfig, "The root `config` value must be a non-empty repository-relative scalar path.", s.locForNodeOrMissing(value, path, workspace.FrontmatterConfigField), nil)
		return
	}
	if _, err := workspace.CleanRepoRelative(value.Value); err != nil {
		s.add(RuleInvalidConfig, "The root `config` value is invalid: "+err.Error()+".", s.loc(value, path, workspace.FrontmatterConfigField), nil)
	}
}

func (s *runState) checkSchemaProperty(node qschema.Node, property qschema.Property, owner, key, value *yaml.Node, path []PathSegment, locationLabel string, context schemaContext) {
	if node.Kind == qschema.RequirementKind && property.Name == qschema.PropertyAssessment {
		s.checkRequiredAssessment(value, path, locationLabel, context)
		return
	}

	if property.Presence != qschema.RequiredPresence && isEmpty(value) {
		s.emptyProperty(owner, key, path, locationLabel)
		return
	}

	switch property.Shape {
	case qschema.ScalarShape:
		if !s.checkRequiredScalar(property, key, value, path, locationLabel) {
			s.checkOptionalScalar(owner, key, value, path, locationLabel)
		}
	case qschema.MapShape:
		if s.validContainerShape(property, key, value, path, locationLabel, context, yaml.MappingNode) && property.ValueShape == qschema.ScalarShape {
			s.checkScalarMapValues(value, path)
		}
	case qschema.SequenceShape:
		if s.validContainerShape(property, key, value, path, locationLabel, context, yaml.SequenceNode) && property.ElementShape == qschema.ScalarShape {
			s.checkScalarSequenceItems(value, path, context)
		}
	}
}

func (s *runState) checkRequiredAssessment(value *yaml.Node, path []PathSegment, locationLabel string, context schemaContext) {
	if value.Kind != yaml.ScalarNode || isEmpty(value) {
		s.add(RuleInvalidAssessment, "The requirement `"+context.requirement+"` has an invalid `assessment`; a requirement must declare exactly one non-empty scalar assessment.", s.loc(value, path, locationLabel), nil)
	}
}

func (s *runState) checkRequiredScalar(property qschema.Property, key, value *yaml.Node, path []PathSegment, locationLabel string) bool {
	if property.Presence != qschema.RequiredPresence {
		return false
	}
	if value.Kind != yaml.ScalarNode && !isEmpty(value) {
		s.invalid(key, path, locationLabel, "The `"+property.Name+"` property has the wrong YAML shape; it must be a scalar.")
	}
	return true
}

func (s *runState) validContainerShape(property qschema.Property, key, value *yaml.Node, path []PathSegment, locationLabel string, context schemaContext, want yaml.Kind) bool {
	if isEmpty(value) {
		return false
	}
	if value.Kind != want {
		s.invalid(key, path, locationLabel, wrongShapeMessage(property, context))
		return false
	}
	return true
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
	case qschema.AreaKind:
		return "The area `" + context.areaName + "` declares an unknown key `" + key + "`; areas may only use area properties."
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
	s.checkRatingLevelID(level, path, index, seen)
	if _, value, _ := document.MapEntry(level, qschema.PropertyCriterion); value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingCriterion, "A rating level declares no `criterion`; each rating level requires a non-empty criterion.", s.locForNodeOrMissing(value, appendPath(path, "criterion"), "ratingScale["+strconv.Itoa(index)+"].criterion"), nil)
	}
	s.checkRequiredTitle(level, path, "rating level", "The rating level at `ratingScale["+strconv.Itoa(index)+"]` declares no `title`; each rating level requires a human-facing title.")
	if _, value, _ := document.MapEntry(level, qschema.PropertyDescription); value == nil {
		s.add(RuleMissingLevelDescription, "A rating level declares no `description`; a description is recommended for each level.", s.locForMissing(appendPath(path, "description"), "ratingScale["+strconv.Itoa(index)+"].description"), nil)
	}
}

func (s *runState) checkRatingLevelID(level *yaml.Node, path []PathSegment, index int, seen map[string]Location) {
	_, value, _ := document.MapEntry(level, qschema.PropertyLevel)
	label := "ratingScale[" + strconv.Itoa(index) + "].level"
	levelPath := appendPath(path, "level")
	if value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingLevelName, "A rating level declares no `level` name; each rating level requires a non-empty `level`.", s.locForNodeOrMissing(value, levelPath, label), nil)
		return
	}

	name := value.Value
	location := s.loc(value, levelPath, label)
	if !validModelName(name) {
		s.add(RuleInvalidRatingLevelID, "The rating level ID `"+name+"` is invalid; Area names, Factor names, and Rating Level IDs must match "+qschema.ModelNamePattern+".", location, nil)
	}
	if prior, ok := seen[name]; ok {
		s.add(RuleDuplicateLevel, "The rating level `"+name+"` is duplicated; each `level` name must be unique within `ratingScale`.", location, nil)
		s.add(RuleDuplicateLevel, "The rating level `"+name+"` is duplicated; each `level` name must be unique within `ratingScale`.", prior, nil)
		return
	}
	seen[name] = location
	s.levels[name] = true
}

func (s *runState) checkEmptyModel() {
	for _, group := range qschema.Model.RequiredAny {
		if group.Name == "model-content" && !s.rootHasAny(group.Properties) {
			s.add(RuleEmptyModel, "The model root supplies no entries under `factors`, `requirements`, or `areas`; a QUALITY.md model requires model content.", s.loc(s.doc.Frontmatter, []PathSegment{}, "frontmatter"), nil)
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
		case qschema.PropertyAreas:
			if len(s.root.areas) > 0 {
				return true
			}
		}
	}
	return false
}

func (s *runState) checkAreaShape(area *areaRef) {
	s.checkSchemaProperties(qschema.Area, area.node, area.path, schemaContext{areaName: area.name})
	s.checkRequiredTitle(area.node, area.path, "area", "The area `"+area.name+"` declares no `title`; each area requires a human-facing title.")
}

func (s *runState) checkFactorShape(factor *factorRef) {
	s.checkSchemaProperties(qschema.Factor, factor.node, factor.path, schemaContext{factorName: factor.name})
	s.checkRequiredTitle(factor.node, factor.path, "factor", "The factor `"+factor.name+"` declares no `title`; each factor requires a human-facing title.")
	if _, value, _ := document.MapEntry(factor.node, qschema.PropertyDescription); value == nil {
		s.add(RuleMissingFactorDescription, "The factor `"+factor.name+"` declares no `description`; a description is recommended for each factor.", s.locForMissing(appendPath(factor.path, "description"), label(appendPath(factor.path, "description"))), nil)
	}
}

func validModelName(name string) bool {
	return modelNamePattern.MatchString(name)
}

func (s *runState) checkRequiredTitle(owner *yaml.Node, base []PathSegment, kind, message string) {
	_, value, _ := document.MapEntry(owner, qschema.PropertyTitle)
	path := appendPath(base, qschema.PropertyTitle)
	locationLabel := label(path)
	if kind == "model root" {
		locationLabel = "title"
	}
	if value == nil || isEmpty(value) || value.Kind != yaml.ScalarNode {
		s.add(RuleMissingTitle, message, s.locForNodeOrMissing(value, path, locationLabel), nil)
	}
}

func (s *runState) checkRequirementShape(req *requirementRef) {
	s.checkSchemaProperties(qschema.Requirement, req.node, req.path, schemaContext{requirement: req.statement})
	if _, value, _ := document.MapEntry(req.node, qschema.PropertyAssessment); value == nil {
		s.add(RuleInvalidAssessment, "The requirement `"+req.statement+"` has no `assessment`; a requirement must declare exactly one non-empty scalar assessment.", s.locForMissing(appendPath(req.path, "assessment"), label(appendPath(req.path, "assessment"))), nil)
	}
}

func (s *runState) checkAreas(area *areaRef) bool {
	hasRequirements := len(area.requirements) > 0
	for _, factor := range area.factors {
		if factorHasRequirements(factor) {
			hasRequirements = true
		}
	}
	for _, child := range area.areas {
		if s.checkAreas(child) {
			hasRequirements = true
		}
	}
	if area.parent != nil && !hasRequirements {
		s.add(RuleEmptyArea, "The area `"+area.name+"` reaches no requirements in its subtree; each area should lead to at least one requirement.", s.loc(area.node, area.path, label(area.path)), nil)
	}
	return hasRequirements
}

func (s *runState) checkFactors(area *areaRef) {
	for _, factor := range area.factors {
		s.checkFactor(factor)
	}
	for _, child := range area.areas {
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
	for _, req := range allRequirements(factor.area) {
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

func (s *runState) checkRequirements(area *areaRef) {
	for _, req := range area.requirements {
		s.checkRequirementRefs(req)
	}
	for _, factor := range area.factors {
		s.checkFactorRequirements(factor)
	}
	for _, child := range area.areas {
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
			if s.resolveFactor(req.area, item.Value) == nil {
				path := appendPath(req.path, "factors", i)
				s.add(RuleUnknownFactor, "The requirement `"+req.statement+"` references unknown factor `"+item.Value+"`; factor references must resolve on the declaring area or an ancestor.", s.loc(item, path, label(path)), nil)
			}
		}
	}
}
