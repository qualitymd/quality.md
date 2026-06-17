// Package schema defines the structural QUALITY.md frontmatter schema.
package schema

// NodeKind identifies a structural frontmatter node.
type NodeKind string

const (
	ModelKind       NodeKind = "model"
	TargetKind      NodeKind = "target"
	FactorKind      NodeKind = "factor"
	RequirementKind NodeKind = "requirement"
	RatingLevelKind NodeKind = "ratingLevel"
)

const (
	PropertyTitle        = "title"
	PropertyRatingScale  = "ratingScale"
	PropertyFactors      = "factors"
	PropertyRequirements = "requirements"
	PropertyTargets      = "targets"
	PropertySource       = "source"
	PropertyLevel        = "level"
	PropertyDescription  = "description"
	PropertyCriterion    = "criterion"
	PropertyAssessment   = "assessment"
	PropertyRatings      = "ratings"
)

// Shape is the YAML shape a property value must have.
type Shape string

const (
	ScalarShape   Shape = "scalar"
	MapShape      Shape = "map"
	SequenceShape Shape = "sequence"
)

// Presence describes whether a property is required, recommended, or optional.
type Presence string

const (
	RequiredPresence    Presence = "required"
	RecommendedPresence Presence = "recommended"
	OptionalPresence    Presence = "optional"
)

// Property is one valid key on a structural node.
type Property struct {
	Name         string
	Shape        Shape
	Presence     Presence
	ElementKind  NodeKind
	ElementShape Shape
	ValueShape   Shape
	MinItems     int
}

// RequiredAny requires at least one of its properties to contain entries.
type RequiredAny struct {
	Name       string
	Properties []string
}

// Node is the structural schema for a frontmatter node.
type Node struct {
	Kind        NodeKind
	Properties  []Property
	RequiredAny []RequiredAny
}

// Property returns the schema property named name.
func (n Node) Property(name string) (Property, bool) {
	for _, property := range n.Properties {
		if property.Name == name {
			return property, true
		}
	}
	return Property{}, false
}

// HasProperty reports whether name is valid on the node.
func (n Node) HasProperty(name string) bool {
	_, ok := n.Property(name)
	return ok
}

// PropertyNames returns the valid property names in schema order.
func (n Node) PropertyNames() []string {
	names := make([]string, 0, len(n.Properties))
	for _, property := range n.Properties {
		names = append(names, property.Name)
	}
	return names
}

// Model is the root QUALITY.md model schema.
var Model = Node{
	Kind: ModelKind,
	Properties: []Property{
		{Name: PropertyTitle, Shape: ScalarShape, Presence: RecommendedPresence},
		{Name: PropertyRatingScale, Shape: SequenceShape, Presence: RequiredPresence, ElementKind: RatingLevelKind, MinItems: 2},
		{Name: PropertyFactors, Shape: MapShape, Presence: OptionalPresence, ElementKind: FactorKind},
		{Name: PropertyRequirements, Shape: MapShape, Presence: OptionalPresence, ElementKind: RequirementKind},
		{Name: PropertyTargets, Shape: MapShape, Presence: OptionalPresence, ElementKind: TargetKind},
		{Name: PropertySource, Shape: ScalarShape, Presence: OptionalPresence},
	},
	RequiredAny: []RequiredAny{
		{Name: "model-content", Properties: []string{PropertyFactors, PropertyRequirements, PropertyTargets}},
	},
}

// Target is the recursive target-node schema.
var Target = Node{
	Kind: TargetKind,
	Properties: []Property{
		{Name: PropertyFactors, Shape: MapShape, Presence: OptionalPresence, ElementKind: FactorKind},
		{Name: PropertyRequirements, Shape: MapShape, Presence: OptionalPresence, ElementKind: RequirementKind},
		{Name: PropertyTargets, Shape: MapShape, Presence: OptionalPresence, ElementKind: TargetKind},
		{Name: PropertySource, Shape: ScalarShape, Presence: OptionalPresence},
	},
}

// Factor is the recursive factor-node schema.
var Factor = Node{
	Kind: FactorKind,
	Properties: []Property{
		{Name: PropertyDescription, Shape: ScalarShape, Presence: RecommendedPresence},
		{Name: PropertyFactors, Shape: MapShape, Presence: OptionalPresence, ElementKind: FactorKind},
		{Name: PropertyRequirements, Shape: MapShape, Presence: OptionalPresence, ElementKind: RequirementKind},
	},
}

// Requirement is the requirement-node schema.
var Requirement = Node{
	Kind: RequirementKind,
	Properties: []Property{
		{Name: PropertyAssessment, Shape: ScalarShape, Presence: RequiredPresence},
		{Name: PropertyFactors, Shape: SequenceShape, Presence: OptionalPresence, ElementShape: ScalarShape},
		{Name: PropertyRatings, Shape: MapShape, Presence: OptionalPresence, ValueShape: ScalarShape},
	},
}

// RatingLevel is one rating-scale level schema.
var RatingLevel = Node{
	Kind: RatingLevelKind,
	Properties: []Property{
		{Name: PropertyLevel, Shape: ScalarShape, Presence: RequiredPresence},
		{Name: PropertyTitle, Shape: ScalarShape, Presence: OptionalPresence},
		{Name: PropertyDescription, Shape: ScalarShape, Presence: RecommendedPresence},
		{Name: PropertyCriterion, Shape: ScalarShape, Presence: RequiredPresence},
	},
}

// Nodes lists every structural schema node.
var Nodes = []Node{Model, Target, Factor, Requirement, RatingLevel}
