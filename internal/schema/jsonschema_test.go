package schema

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func generate(t *testing.T) map[string]any {
	t.Helper()
	data, err := GenerateJSON()
	if err != nil {
		t.Fatalf("GenerateJSON() error = %v", err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("GenerateJSON() produced invalid JSON: %v", err)
	}
	return doc
}

func TestGenerateJSONIsDeterministic(t *testing.T) {
	first, err := GenerateJSON()
	if err != nil {
		t.Fatalf("GenerateJSON() error = %v", err)
	}
	second, err := GenerateJSON()
	if err != nil {
		t.Fatalf("GenerateJSON() error = %v", err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("GenerateJSON() output is not byte-stable across runs")
	}
}

func TestGenerateJSONDeclaresDialectAndID(t *testing.T) {
	doc := generate(t)
	if got := doc["$schema"]; got != JSONSchemaDialect {
		t.Errorf("$schema = %v, want %v", got, JSONSchemaDialect)
	}
	if got := doc["$id"]; got != JSONSchemaID {
		t.Errorf("$id = %v, want %v", got, JSONSchemaID)
	}
}

func TestGenerateJSONRootRequiresTitleAndRatingScale(t *testing.T) {
	doc := generate(t)
	required := stringSet(t, doc["required"])
	for _, name := range []string{PropertyTitle, PropertyRatingScale} {
		if !required[name] {
			t.Errorf("root required missing %q", name)
		}
	}
	// description and source are optional and must not be required.
	for _, name := range []string{PropertyDescription, PropertySource} {
		if required[name] {
			t.Errorf("root required unexpectedly includes optional %q", name)
		}
	}
}

func TestGenerateJSONRootRequiresAnyContent(t *testing.T) {
	doc := generate(t)
	anyOf, ok := doc["anyOf"].([]any)
	if !ok {
		t.Fatalf("root anyOf = %T, want []any", doc["anyOf"])
	}
	got := map[string]bool{}
	for _, alt := range anyOf {
		for _, name := range stringSlice(t, alt.(map[string]any)["required"]) {
			got[name] = true
		}
	}
	for _, name := range []string{PropertyFactors, PropertyRequirements, PropertyAreas} {
		if !got[name] {
			t.Errorf("root anyOf missing alternative requiring %q", name)
		}
	}
}

func TestGenerateJSONRatingScaleMinItems(t *testing.T) {
	doc := generate(t)
	props := doc["properties"].(map[string]any)
	scale := props[PropertyRatingScale].(map[string]any)
	if got := scale["type"]; got != "array" {
		t.Errorf("ratingScale type = %v, want array", got)
	}
	if got := scale["minItems"]; got != float64(2) {
		t.Errorf("ratingScale minItems = %v, want 2", got)
	}
	items := scale["items"].(map[string]any)
	if got := items["$ref"]; got != "#/$defs/"+string(RatingLevelKind) {
		t.Errorf("ratingScale items $ref = %v", got)
	}
}

func TestGenerateJSONDefinesRecursiveNodes(t *testing.T) {
	doc := generate(t)
	defs, ok := doc["$defs"].(map[string]any)
	if !ok {
		t.Fatalf("$defs = %T, want map", doc["$defs"])
	}
	for _, kind := range []NodeKind{AreaKind, FactorKind, RequirementKind, RatingLevelKind} {
		if _, ok := defs[string(kind)]; !ok {
			t.Errorf("$defs missing %q", kind)
		}
	}
	// factors maps recurse into the factor definition.
	factor := defs[string(FactorKind)].(map[string]any)
	factorProps := factor["properties"].(map[string]any)
	nested := factorProps[PropertyFactors].(map[string]any)
	ref := nested["additionalProperties"].(map[string]any)
	if got := ref["$ref"]; got != "#/$defs/"+string(FactorKind) {
		t.Errorf("factor.factors $ref = %v, want self-reference", got)
	}
}

func TestGenerateJSONIncludesStrictModelNamePatterns(t *testing.T) {
	doc := generate(t)
	props := doc["properties"].(map[string]any)
	assertPropertyNamesPattern(t, props[PropertyAreas], ModelNamePattern)
	assertPropertyNamesPattern(t, props[PropertyFactors], ModelNamePattern)
	assertPropertyNamesPattern(t, props[PropertyRequirements], ModelNamePattern)

	scale := props[PropertyRatingScale].(map[string]any)
	items := scale["items"].(map[string]any)
	defs := doc["$defs"].(map[string]any)
	ratingLevel := defs[strings.TrimPrefix(items["$ref"].(string), "#/$defs/")].(map[string]any)
	level := ratingLevel["properties"].(map[string]any)[PropertyLevel].(map[string]any)
	if got := level["pattern"]; got != ModelNamePattern {
		t.Fatalf("rating level pattern = %v, want %s", got, ModelNamePattern)
	}

	requirement := defs[string(RequirementKind)].(map[string]any)
	requirementRequired := stringSet(t, requirement["required"])
	for _, name := range []string{PropertyTitle, PropertyAssessment} {
		if !requirementRequired[name] {
			t.Errorf("requirement required missing %q", name)
		}
	}
}

func TestGenerateJSONAllowsExtensions(t *testing.T) {
	doc := generate(t)
	// The format permits extension frontmatter, so no node object closes itself
	// off with additionalProperties. (A map-valued property such as
	// requirement.ratings legitimately uses additionalProperties to type its
	// entries, but that sits on the property schema, not on a node object.)
	if _, ok := doc["additionalProperties"]; ok {
		t.Error("root sets additionalProperties; extension frontmatter would be rejected")
	}
	defs := doc["$defs"].(map[string]any)
	for name, def := range defs {
		if _, ok := def.(map[string]any)["additionalProperties"]; ok {
			t.Errorf("$defs.%s sets additionalProperties on the node object", name)
		}
	}
}

func assertPropertyNamesPattern(t *testing.T, v any, want string) {
	t.Helper()
	propertyNames, ok := v.(map[string]any)["propertyNames"].(map[string]any)
	if !ok {
		t.Fatalf("propertyNames = %T, want map", v.(map[string]any)["propertyNames"])
	}
	if got := propertyNames["pattern"]; got != want {
		t.Fatalf("propertyNames.pattern = %v, want %s", got, want)
	}
}

func stringSet(t *testing.T, v any) map[string]bool {
	t.Helper()
	set := map[string]bool{}
	for _, s := range stringSlice(t, v) {
		set[s] = true
	}
	return set
}

func stringSlice(t *testing.T, v any) []string {
	t.Helper()
	raw, ok := v.([]any)
	if !ok {
		t.Fatalf("value = %T, want []any", v)
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok {
			t.Fatalf("element = %T, want string", item)
		}
		out = append(out, s)
	}
	return out
}
