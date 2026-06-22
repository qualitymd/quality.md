package evaluation

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// FieldProblem describes one caller-facing payload validation problem.
type FieldProblem struct {
	Field    string
	Problem  string
	Expected string
	Allowed  []string
}

// ValidationError aggregates payload validation problems.
type ValidationError struct {
	Problems []FieldProblem
}

func (e *ValidationError) Error() string {
	if len(e.Problems) == 0 {
		return "invalid payload"
	}
	lines := []string{"invalid payload:"}
	for _, problem := range e.Problems {
		detail := problem.Problem
		if problem.Expected != "" {
			detail += "; expected " + problem.Expected
		}
		if len(problem.Allowed) > 0 {
			detail += "; allowed values: " + strings.Join(problem.Allowed, ", ")
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", problem.Field, detail))
	}
	return strings.Join(lines, "\n")
}

type validationAccumulator struct {
	problems []FieldProblem
}

func (a *validationAccumulator) Add(field, problem string) {
	a.AddExpected(field, problem, "")
}

func (a *validationAccumulator) AddExpected(field, problem, expected string) {
	a.problems = append(a.problems, FieldProblem{Field: field, Problem: problem, Expected: expected})
}

func (a *validationAccumulator) AddAllowed(field, problem string, allowed []string) {
	a.problems = append(a.problems, FieldProblem{Field: field, Problem: problem, Allowed: allowed})
}

func (a *validationAccumulator) Merge(prefix string, err error) {
	if err == nil {
		return
	}
	var validationErr *ValidationError
	if !asValidationError(err, &validationErr) {
		a.Add(prefix, err.Error())
		return
	}
	for _, problem := range validationErr.Problems {
		problem.Field = prefixField(prefix, problem.Field)
		a.problems = append(a.problems, problem)
	}
}

func (a *validationAccumulator) Err() error {
	if len(a.problems) == 0 {
		return nil
	}
	return usagef("%w", &ValidationError{Problems: a.problems})
}

func asValidationError(err error, target **ValidationError) bool {
	return errors.As(err, target)
}

func prefixField(prefix, field string) string {
	if prefix == "" {
		return field
	}
	if field == "" {
		return prefix
	}
	return prefix + "." + field
}

type payloadField struct {
	Name     string
	Type     string
	Required bool
	Enum     []string
}

var payloadFieldAnnotations = map[RecordKind]map[string]payloadField{
	KindAssessmentResult: {
		"areaPath":        {Required: true},
		"requirement":     {Required: true},
		"factorPaths":     {Required: true},
		"ratingResult":    {Required: true},
		"criterionSource": {Required: true, Enum: []string{"rating-scale", "requirement"}},
		"findings":        {Required: true},
		"recommendations": {Required: true},
	},
	KindAnalysis: {
		"areaPath":                {Required: true},
		"factorRatingResults":     {Required: true},
		"aggregateRatingResult":   {Required: true},
		"assessmentResultRecords": {Required: true},
		"childAnalysisRecords":    {Required: true},
	},
	KindRecommendation: {
		"title":              {Required: true},
		"gap":                {Required: true},
		"evidenceLocators":   {Required: true},
		"remediationOptions": {Required: true},
		"recommendedOption":  {Required: true},
		"doneCriterion":      {Required: true},
	},
}

// PayloadHelp returns the JSON payload field table for a record kind.
func PayloadHelp(kind RecordKind) string {
	fields := payloadFields(kind)
	if len(fields) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Payload fields:\n\n")
	b.WriteString("| Field | Type | Required | Allowed values |\n")
	b.WriteString("| --- | --- | --- | --- |\n")
	for _, field := range fields {
		required := "no"
		if field.Required {
			required = "yes"
		}
		allowed := "-"
		if len(field.Enum) > 0 {
			allowed = strings.Join(field.Enum, ", ")
		}
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n", field.Name, field.Type, required, allowed)
	}
	return b.String()
}

func payloadFields(kind RecordKind) []payloadField {
	var typ reflect.Type
	switch kind {
	case KindAssessmentResult:
		typ = reflect.TypeOf(AssessmentResultInput{})
	case KindAnalysis:
		typ = reflect.TypeOf(AnalysisInput{})
	case KindRecommendation:
		typ = reflect.TypeOf(RecommendationInput{})
	default:
		return nil
	}
	annotations := payloadFieldAnnotations[kind]
	fields := make([]payloadField, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := jsonFieldName(field)
		if name == "" {
			continue
		}
		info := annotations[name]
		info.Name = name
		info.Type = jsonTypeDescription(field.Type)
		fields = append(fields, info)
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})
	return fields
}

// CanonicalPayloadExample returns a complete valid JSON payload for a write kind.
func CanonicalPayloadExample(kind RecordKind) string {
	switch kind {
	case KindAssessmentResult:
		return assessmentResultPayloadExample
	case KindAnalysis:
		return analysisPayloadExample
	case KindRecommendation:
		return recommendationPayloadExample
	default:
		return "{}\n"
	}
}

const assessmentResultPayloadExample = `[
  {
    "areaPath": [],
    "requirement": "Has tests",
    "factorPaths": [],
    "ratingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "Evidence supports the target level."
    },
    "criterionSource": "rating-scale",
    "findings": [
      {
        "locator": "tests/example_test.go:1",
        "observation": "The requirement is covered by a focused test.",
        "category": "coverage",
        "severity": "low",
        "evidence": [
          {
            "kind": "source",
            "ref": "tests/example_test.go:1"
          }
        ]
      }
    ],
    "recommendations": []
  }
]
`

const analysisPayloadExample = `[
  {
    "areaPath": [],
    "localRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The local assessment result reaches target."
    },
    "factorRatingResults": [],
    "aggregateRatingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "The root local rating binds the aggregate rating."
    },
    "assessmentResultRecords": [
      "assessments/001-root-has-tests.json"
    ],
    "childAnalysisRecords": []
  }
]
`

const recommendationPayloadExample = `[
  {
    "title": "Improve test coverage",
    "gap": "The evaluation found a requirement with thin test evidence.",
    "evidenceLocators": [
      "assessments/001-root-has-tests.json"
    ],
    "assessmentResultRecords": [
      "assessments/001-root-has-tests.json"
    ],
    "remediationOptions": [
      "Add focused tests for the requirement"
    ],
    "recommendedOption": "Add focused tests for the requirement",
    "doneCriterion": "The affected requirement reaches target with current test evidence."
  }
]
`

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return ""
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return field.Name
	}
	return name
}

func jsonTypeDescription(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}
	switch typ.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice, reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return "string"
		}
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	default:
		return "value"
	}
}
