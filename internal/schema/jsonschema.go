package schema

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// JSON Schema identity for the companion artifact.
const (
	// JSONSchemaDialect is the JSON Schema dialect the companion artifact
	// declares. Draft 2020-12 supports the $defs/$ref, additionalProperties,
	// anyOf, and minItems constructs the structural schema needs, and is
	// understood by the YAML language servers that back editor validation.
	JSONSchemaDialect = "https://json-schema.org/draft/2020-12/schema"

	// JSONSchemaID is the stable, unversioned identifier for the companion
	// schema. A $id is an identifier, not a fetch URL, so reserving the
	// project-domain identity is independent of whether the schema is ever
	// hosted there.
	JSONSchemaID = "https://getquality.md/quality.schema.json"
)

const (
	jsonSchemaTitle       = "QUALITY.md frontmatter"
	jsonSchemaDescription = "Structural JSON Schema for QUALITY.md frontmatter, derived from the qualitymd linter's schema."
	// jsonSchemaComment marks the artifact non-normative and subordinate so a
	// consumer cannot mistake passing structural validation for full
	// conformance, and warns hand-editors that the file is generated.
	jsonSchemaComment = "Non-normative and subordinate to SPECIFICATION.md (https://getquality.md). " +
		"Structural-only: passing this schema does not imply full conformance. Semantic rules " +
		"(factor-reference resolution, rating-override keys, the placement-dependent factor-connection rule, " +
		"and rating-level ordering and uniqueness) are enforced by `qualitymd lint`, not here. " +
		"Generated from internal/schema and guarded against drift by a consistency test; " +
		"do not edit by hand — run `go generate ./...`."
)

// GenerateJSON renders the structural frontmatter schema as a JSON Schema
// (draft 2020-12) document. It is derived entirely from the Node definitions in
// this package, so the companion artifact cannot encode a structural rule the
// linter does not enforce, or omit one it does.
//
// The output is deterministic: encoding/json emits object keys in sorted order,
// so re-running GenerateJSON over an unchanged schema yields byte-identical
// output. A consistency test relies on that to detect drift between the
// committed artifact and this generator.
func GenerateJSON() ([]byte, error) {
	root, err := nodeSchema(Model)
	if err != nil {
		return nil, err
	}
	root["$schema"] = JSONSchemaDialect
	root["$id"] = JSONSchemaID
	root["title"] = jsonSchemaTitle
	root["description"] = jsonSchemaDescription
	root["$comment"] = jsonSchemaComment

	defs := map[string]any{}
	for _, node := range []Node{Area, Factor, Requirement, RatingLevel} {
		def, err := nodeSchema(node)
		if err != nil {
			return nil, err
		}
		defs[string(node.Kind)] = def
	}
	root["$defs"] = defs

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	if err := enc.Encode(root); err != nil {
		return nil, fmt.Errorf("encoding JSON schema: %w", err)
	}
	return buf.Bytes(), nil
}

// nodeSchema renders one structural node as an object schema. Nodes never set
// additionalProperties:false: the format permits extension frontmatter, so a
// closed schema would reject conforming documents that use it.
func nodeSchema(n Node) (map[string]any, error) {
	properties := map[string]any{}
	var required []string
	for _, property := range n.Properties {
		shape, err := propertySchema(property)
		if err != nil {
			return nil, err
		}
		properties[property.Name] = shape
		if property.Presence == RequiredPresence {
			required = append(required, property.Name)
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if len(required) > 0 {
		schema["required"] = required
	}

	switch len(n.RequiredAny) {
	case 0:
	case 1:
		schema["anyOf"] = anyOfRequired(n.RequiredAny[0])
	default:
		// Each RequiredAny is an independent "at least one of" constraint, so
		// they compose under allOf rather than collapsing into one anyOf.
		allOf := make([]any, 0, len(n.RequiredAny))
		for _, requiredAny := range n.RequiredAny {
			allOf = append(allOf, map[string]any{"anyOf": anyOfRequired(requiredAny)})
		}
		schema["allOf"] = allOf
	}

	return schema, nil
}

// anyOfRequired expresses "at least one of these properties is present" as an
// anyOf of single-property required clauses.
func anyOfRequired(r RequiredAny) []any {
	alternatives := make([]any, 0, len(r.Properties))
	for _, name := range r.Properties {
		alternatives = append(alternatives, map[string]any{"required": []string{name}})
	}
	return alternatives
}

// propertySchema renders the JSON Schema for one property's value shape.
func propertySchema(p Property) (map[string]any, error) {
	switch p.Shape {
	case ScalarShape:
		schema := map[string]any{"type": "string"}
		if p.Pattern != "" {
			schema["pattern"] = p.Pattern
		}
		return schema, nil
	case MapShape:
		element, err := elementSchema(p)
		if err != nil {
			return nil, err
		}
		schema := map[string]any{
			"type":                 "object",
			"additionalProperties": element,
		}
		if p.KeyPattern != "" {
			propertyNames := map[string]any{"pattern": p.KeyPattern}
			if p.ElementKind == AreaKind {
				propertyNames = map[string]any{
					"allOf": []any{
						map[string]any{"pattern": p.KeyPattern},
						map[string]any{"not": map[string]any{"enum": []any{"root"}}},
					},
				}
			}
			schema["propertyNames"] = propertyNames
		}
		return schema, nil
	case SequenceShape:
		element, err := elementSchema(p)
		if err != nil {
			return nil, err
		}
		schema := map[string]any{
			"type":  "array",
			"items": element,
		}
		if p.MinItems > 0 {
			schema["minItems"] = p.MinItems
		}
		return schema, nil
	default:
		return nil, fmt.Errorf("schema: property %q has unknown shape %q", p.Name, p.Shape)
	}
}

// elementSchema resolves the element type of a map value or sequence item: a
// $ref to a node definition when the element is a structural node, or a scalar
// type when it is a plain string.
func elementSchema(p Property) (map[string]any, error) {
	if p.ElementKind != "" {
		return map[string]any{"$ref": "#/$defs/" + string(p.ElementKind)}, nil
	}
	if p.ElementShape == ScalarShape || p.ValueShape == ScalarShape {
		return map[string]any{"type": "string"}, nil
	}
	return nil, fmt.Errorf("schema: property %q has an unresolved element type", p.Name)
}
