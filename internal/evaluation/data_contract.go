package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qualitymd/quality.md/internal/model"
	qschema "github.com/qualitymd/quality.md/internal/schema"
)

const evaluationDataSchemaID = "https://getquality.md/evaluation-data.schema.json"

type dataFieldType string

const (
	dataAny           dataFieldType = "any"
	dataBool          dataFieldType = "bool"
	dataNumber        dataFieldType = "number"
	dataString        dataFieldType = "string"
	dataStringArray   dataFieldType = "stringArray"
	dataObject        dataFieldType = "object"
	dataArray         dataFieldType = "array"
	dataAreaID        dataFieldType = "areaID"
	dataRequirementID dataFieldType = "requirementID"
	dataFactorID      dataFieldType = "factorID"
	dataRatingLevelID dataFieldType = "ratingLevelID"
)

type dataField struct {
	Name     string
	Type     dataFieldType
	Required bool
	Enum     []string
	Object   *dataObjectContract
	Element  *dataField
	MinItems int
}

type dataObjectContract struct {
	Fields []dataField
	Open   bool
}

type dataKindContract struct {
	Kind        DataKind
	Description string
	Object      dataObjectContract
	Example     func() map[string]any
}

type DataVerifyReceipt struct {
	SchemaVersion int                 `json:"schemaVersion"`
	Path          string              `json:"path"`
	Valid         bool                `json:"valid"`
	Checked       int                 `json:"checked"`
	Failures      []DataVerifyFailure `json:"failures,omitempty"`
}

type DataVerifyFailure struct {
	Path   string   `json:"path"`
	Kind   DataKind `json:"kind,omitempty"`
	Reason string   `json:"reason"`
}

var dataContracts = map[DataKind]dataKindContract{}

var dataContractOrder = []DataKind{
	DataKindEvaluationFrame,
	DataKindAreaEvaluationFrame,
	DataKindRequirementEvaluationFrame,
	DataKindRequirementAssessment,
	DataKindRequirementRating,
	DataKindFactorAnalysisFrame,
	DataKindFactorAnalysis,
	DataKindAreaAnalysisFrame,
	DataKindAreaAnalysis,
	DataKindEvaluationOutput,
}

//nolint:funlen // The registry is intentionally declared in one place.
func init() {
	for _, contract := range []dataKindContract{
		{
			Kind:        DataKindEvaluationFrame,
			Description: "Run-level frame for one evaluation.",
			Object: topContract(DataKindEvaluationFrame,
				field("subject", dataObject, false, object(
					field("modelLocator", dataString, false),
				)),
				field("inputs", dataObject, false, object(
					field("requestedScope", dataString, false),
					field("ratingLevelIds", dataArray, false, arrayOf(dataRatingLevelID)),
					field("areaIds", dataArray, false, arrayOf(dataAreaID)),
					field("factorIds", dataArray, false, arrayOf(dataFactorID)),
				)),
				field("derivedContext", dataObject, false, object(
					field("resolvedScope", dataString, false),
					field("rigor", dataString, false),
					field("evaluationPolicies", dataStringArray, false),
					field("expectedEvaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				)),
			),
			Example: func() map[string]any { return evaluationFrameExample(DataKindEvaluationFrame) },
		},
		{
			Kind:        DataKindAreaEvaluationFrame,
			Description: "Area-local frame for source boundaries and local model structure.",
			Object: topContract(DataKindAreaEvaluationFrame,
				field("subject", dataObject, true, object(field("areaId", dataAreaID, true))),
				field("inputs", dataObject, false, object(
					field("sourceRefs", dataArray, false, arrayOfAny()),
					field("localRequirementIds", dataArray, false, arrayOf(dataRequirementID)),
					field("rootFactorIds", dataArray, false, arrayOf(dataFactorID)),
					field("childAreaIds", dataArray, false, arrayOf(dataAreaID)),
				)),
				field("derivedContext", dataObject, false, object(
					field("scope", dataString, false),
					field("expectedEvaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				)),
			),
			Example: func() map[string]any { return areaEvaluationFrameExample(DataKindAreaEvaluationFrame) },
		},
		{
			Kind:        DataKindRequirementEvaluationFrame,
			Description: "Requirement-local frame for evidence targets and applied rating criteria.",
			Object: topContract(DataKindRequirementEvaluationFrame,
				field("subject", dataObject, true, object(
					field("requirementId", dataRequirementID, true),
					field("factorIds", dataArray, false, arrayOf(dataFactorID)),
				)),
				field("inputs", dataObject, false, object(
					field("ratingLevelIds", dataArray, false, arrayOf(dataRatingLevelID)),
					field("requirementAssessmentBasis", dataString, false),
					field("ratingOverrides", dataObject, false, openObject()),
				)),
				field("derivedContext", dataObject, false, object(
					field("evidenceTargets", dataArray, false, arrayOfObject(evidenceTargetContract())),
					field("appliedRatingCriteria", dataArray, false, arrayOfObject(criteriaContract())),
					field("stopConditions", dataArray, false, arrayOfObject(limitContract())),
					field("expectedEvaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				)),
			),
			Example: func() map[string]any { return requirementEvaluationFrameExample(DataKindRequirementEvaluationFrame) },
		},
		{
			Kind:        DataKindRequirementAssessment,
			Description: "Evidence assessment result for one Requirement.",
			Object: topContract(DataKindRequirementAssessment,
				field("requirementId", dataRequirementID, true),
				field("status", dataString, true, enum("assessed", "partially_assessed", "not_assessed", "blocked")),
				field("statusReason", dataString, false),
				field("evidenceSummary", dataString, false),
				field("summary", dataString, false),
				field("factors", dataArray, false, arrayOf(dataFactorID)),
				field("evidenceTargetCoverage", dataArray, false, arrayOfObject(openObject())),
				field("findings", dataArray, true, arrayOfObject(findingContract())),
				field("unknowns", dataArray, false, arrayOfObject(unknownContract())),
				field("evaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				field("confidence", dataString, false, enum("high", "medium", "low", "none")),
				field("confidenceReason", dataString, false),
			),
			Example: func() map[string]any { return requirementAssessmentExample(DataKindRequirementAssessment) },
		},
		{
			Kind:        DataKindRequirementRating,
			Description: "Rating result for one Requirement Assessment.",
			Object: topContract(DataKindRequirementRating,
				field("requirementId", dataRequirementID, true),
				field("status", dataString, true, enum("rated", "not_rated", "blocked")),
				field("statusReason", dataString, false),
				field("ratingLevelId", dataRatingLevelID, false),
				field("rationale", dataString, false),
				field("ratingDrivers", dataArray, false, arrayOfObject(ratingDriverContract())),
				field("criteriaResults", dataArray, false, arrayOfObject(criteriaResultContract())),
				field("missingEvidence", dataArray, false, arrayOfObject(unknownContract())),
				field("evaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				field("confidence", dataString, false, enum("high", "medium", "low", "none")),
				field("confidenceReason", dataString, false),
			),
			Example: func() map[string]any { return requirementRatingExample(DataKindRequirementRating) },
		},
		{
			Kind:        DataKindFactorAnalysisFrame,
			Description: "Factor-local frame for synthesis inputs.",
			Object: topContract(DataKindFactorAnalysisFrame,
				field("subject", dataObject, true, object(
					field("areaId", dataAreaID, false),
					field("factorId", dataFactorID, true),
				)),
				field("inputs", dataObject, false, object(
					field("directRequirementRatingRefs", dataArray, false, arrayOfObject(routineRefContract())),
					field("childFactorAnalysisRefs", dataArray, false, arrayOfObject(routineRefContract())),
				)),
				field("derivedContext", dataObject, false, object(
					field("synthesisGuidanceRef", dataString, false),
					field("emptySignalPolicy", dataString, false),
					field("stopConditions", dataArray, false, arrayOfObject(limitContract())),
					field("expectedEvaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				)),
			),
			Example: func() map[string]any { return factorAnalysisFrameExample(DataKindFactorAnalysisFrame) },
		},
		{
			Kind:        DataKindFactorAnalysis,
			Description: "Analysis result for one Factor node.",
			Object:      analysisResultContract(DataKindFactorAnalysis, "factorId", dataFactorID),
			Example: func() map[string]any {
				return scopedAnalysisExample(DataKindFactorAnalysis, "factorId", exampleFactorID())
			},
		},
		{
			Kind:        DataKindAreaAnalysisFrame,
			Description: "Area-local frame for synthesis inputs.",
			Object: topContract(DataKindAreaAnalysisFrame,
				field("subject", dataObject, true, object(field("areaId", dataAreaID, true))),
				field("inputs", dataObject, false, object(
					field("factorAnalysisRefs", dataArray, false, arrayOfObject(routineRefContract())),
					field("childAreaAnalysisRefs", dataArray, false, arrayOfObject(routineRefContract())),
				)),
				field("derivedContext", dataObject, false, object(
					field("synthesisGuidanceRef", dataString, false),
					field("emptySignalPolicy", dataString, false),
					field("stopConditions", dataArray, false, arrayOfObject(limitContract())),
					field("expectedEvaluationLimits", dataArray, false, arrayOfObject(limitContract())),
				)),
			),
			Example: func() map[string]any { return areaAnalysisFrameExample(DataKindAreaAnalysisFrame) },
		},
		{
			Kind:        DataKindAreaAnalysis,
			Description: "Analysis result for one Area.",
			Object:      areaAnalysisResultContract(),
			Example: func() map[string]any {
				return scopedAnalysisExample(DataKindAreaAnalysis, "areaId", exampleAreaID())
			},
		},
		{
			Kind:        DataKindEvaluationOutput,
			Description: "CLI-owned completed evaluation output index generated by report build.",
			Object: topContract(DataKindEvaluationOutput,
				field("runReportRef", dataObject, true, reportRefContract()),
				field("headlineResultRef", dataObject, true, routineRefContract()),
				field("headlineReportRef", dataObject, true, reportRefContract()),
				field("rootAreaAnalysisRef", dataObject, false, routineRefContract()),
				field("areaOutputs", dataArray, true, arrayOfObject(areaOutputContract())),
				field("reportOutputs", dataArray, true, arrayOfObject(reportRefContract())),
			),
			Example: evaluationOutputExample,
		},
	} {
		dataContracts[contract.Kind] = contract
	}
}

func topContract(kind DataKind, fields ...dataField) dataObjectContract {
	base := []dataField{
		field("schemaVersion", dataNumber, true),
		field("kind", dataString, true, enum(string(kind))),
	}
	return object(append(base, fields...)...)
}

func analysisResultContract(kind DataKind, idName string, idType dataFieldType) dataObjectContract {
	return topContract(kind,
		field(idName, idType, true),
		field("localAnalysis", dataObject, true, analysisScopeContract()),
		field("localAndDescendantAnalysis", dataObject, true, analysisScopeContract()),
	)
}

func areaAnalysisResultContract() dataObjectContract {
	return topContract(DataKindAreaAnalysis,
		field("areaId", dataAreaID, true),
		field("localAnalysis", dataObject, true, analysisScopeContract()),
		field("localAndDescendantAnalysis", dataObject, true, analysisScopeContract()),
	)
}

func analysisScopeContract() dataObjectContract {
	return object(
		field("status", dataString, true, enum("analyzed", "empty", "not_analyzed", "blocked")),
		field("statusReason", dataString, false),
		field("ratingLevelId", dataRatingLevelID, false),
		field("rationale", dataString, false),
		field("inputRefs", dataArray, false, arrayOfObject(routineRefContract())),
		field("ratingDrivers", dataArray, false, arrayOfObject(ratingDriverContract())),
		field("incompleteInputs", dataArray, false, arrayOfObject(limitContract())),
		field("evaluationLimits", dataArray, false, arrayOfObject(limitContract())),
		field("confidence", dataString, false, enum("high", "medium", "low", "none")),
		field("confidenceReason", dataString, false),
	)
}

func findingContract() dataObjectContract {
	fields := append(findingCoreFields(),
		field("candidateActions", dataArray, false, arrayOfObject(candidateActionContract())),
	)
	return object(fields...)
}

func findingCoreFields() []dataField {
	return []dataField{
		field("id", dataString, true),
		field("type", dataString, true, enum("strength", "gap", "risk", "unknown", "note")),
		field("severity", dataString, true, enum("critical", "high", "medium", "low")),
		field("confidence", dataString, true, enum("high", "medium", "low", "none")),
		field("statement", dataString, true),
		field("condition", dataString, true),
		field("criteria", dataArray, true, arrayOfObject(findingCriterionContract()), minItems(1)),
		field("cause", dataObject, true, findingCauseContract()),
		field("effect", dataObject, true, findingEffectContract()),
		field("evidence", dataArray, true, arrayOfObject(findingEvidenceContract()), minItems(1)),
	}
}

func findingCriterionContract() dataObjectContract {
	return object(
		field("requirementId", dataRequirementID, true),
		field("ratingLevelId", dataRatingLevelID, true),
		field("criterion", dataString, true),
		field("rationale", dataString, false),
	)
}

func findingCauseContract() dataObjectContract {
	return object(
		field("status", dataString, true, enum("verified", "plausible", "not_assessed", "not_applicable")),
		field("statement", dataString, true),
		field("rationale", dataString, false),
		field("evidence", dataArray, false, arrayOfObject(findingEvidenceContract())),
	)
}

func findingEffectContract() dataObjectContract {
	return object(
		field("statement", dataString, true),
		field("rationale", dataString, false),
		field("ratingEffect", dataString, false),
	)
}

func findingEvidenceContract() dataObjectContract {
	return object(
		field("sourceRef", dataString, true),
		field("statement", dataString, true),
		field("rationale", dataString, false),
	)
}

func candidateActionContract() dataObjectContract {
	return object(
		field("id", dataString, true),
		field("description", dataString, true),
		field("rationale", dataString, false),
	)
}

func ratingDriverContract() dataObjectContract {
	return object(
		field("id", dataString, false),
		field("description", dataString, false),
		field("summary", dataString, false),
		field("requirementRatingDriver", dataString, false),
		field("effect", dataString, false),
		field("impact", dataString, false),
		field("ratingLevelId", dataRatingLevelID, false),
		field("inputRefs", dataArray, false, arrayOfObject(routineRefContract())),
	)
}

func criteriaResultContract() dataObjectContract {
	return object(
		field("ratingLevelId", dataRatingLevelID, true),
		field("matched", dataBool, true),
		field("rationale", dataString, false),
	)
}

func criteriaContract() dataObjectContract {
	return object(
		field("ratingLevelId", dataRatingLevelID, true),
		field("criterion", dataString, true),
		field("source", dataString, false),
		field("adaptationRationale", dataString, false),
	)
}

func evidenceTargetContract() dataObjectContract {
	return object(
		field("id", dataString, true),
		field("question", dataString, false),
		field("purpose", dataString, false),
		field("sourceRefs", dataArray, false, arrayOfAny()),
		field("required", dataBool, false),
	)
}

func unknownContract() dataObjectContract {
	return object(
		field("id", dataString, false),
		field("description", dataString, false),
		field("reason", dataString, false),
		field("ref", dataString, false),
		field("impact", dataString, false),
	)
}

func limitContract() dataObjectContract {
	return object(
		field("id", dataString, false),
		field("type", dataString, false),
		field("scope", dataString, false),
		field("ref", dataString, false),
		field("description", dataString, false),
		field("reason", dataString, false),
		field("impact", dataString, false),
	)
}

func routineRefContract() dataObjectContract {
	return object(
		field("kind", dataString, true, enum(kindStrings(supportedDataKinds)...)),
		field("subject", dataObject, true, object(
			field("areaId", dataAreaID, false),
			field("factorId", dataFactorID, false),
			field("requirementId", dataRequirementID, false),
		)),
		field("selector", dataString, false),
	)
}

func reportRefContract() dataObjectContract {
	return object(
		field("kind", dataString, true, enum(kindStrings(reportKinds)...)),
		field("areaId", dataAreaID, false),
		field("factorId", dataFactorID, false),
		field("requirementId", dataRequirementID, false),
		field("path", dataString, true),
	)
}

func areaOutputContract() dataObjectContract {
	return object(
		field("areaId", dataAreaID, true),
		field("areaEvaluationFrameRef", dataObject, true, routineRefContract()),
		field("areaAnalysisResultRef", dataObject, true, routineRefContract()),
		field("factorAnalysisRefs", dataArray, true, arrayOfObject(routineRefContract())),
		field("requirementAssessmentRefs", dataArray, true, arrayOfObject(routineRefContract())),
		field("requirementRatingRefs", dataArray, true, arrayOfObject(routineRefContract())),
		field("reportRefs", dataArray, true, arrayOfObject(reportRefContract())),
	)
}

func object(fields ...dataField) dataObjectContract {
	return dataObjectContract{Fields: fields}
}

func openObject() dataObjectContract {
	return dataObjectContract{Open: true}
}

func field(name string, typ dataFieldType, required bool, opts ...any) dataField {
	f := dataField{Name: name, Type: typ, Required: required}
	for _, opt := range opts {
		switch v := opt.(type) {
		case []string:
			f.Enum = v
		case dataObjectContract:
			f.Object = &v
		case dataField:
			f.Element = &v
		case dataMinItems:
			f.MinItems = int(v)
		}
	}
	return f
}

type dataMinItems int

func minItems(n int) dataMinItems {
	return dataMinItems(n)
}

func enum(values ...string) []string {
	return values
}

func arrayOf(typ dataFieldType) dataField {
	return field("", typ, false)
}

func arrayOfObject(obj dataObjectContract) dataField {
	return field("", dataObject, false, obj)
}

func arrayOfAny() dataField {
	return field("", dataAny, false)
}

func resolveDataKindArg(raw string) (DataKind, error) {
	candidate := DataKind(raw)
	if _, ok := dataContracts[candidate]; ok {
		return candidate, nil
	}
	for _, kind := range dataContractOrder {
		if raw == kebabDataKind(kind) {
			return kind, nil
		}
	}
	return "", usagef("unknown evaluation data kind %q", raw)
}

// ResolveDataKind accepts a canonical or kebab-case Evaluation data kind.
func ResolveDataKind(raw string) (DataKind, error) {
	return resolveDataKindArg(raw)
}

func kebabDataKind(kind DataKind) string {
	var b strings.Builder
	for i, r := range string(kind) {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('-')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

func validateDataPayload(kind DataKind, payload map[string]any) error {
	contract, ok := dataContracts[kind]
	if !ok {
		return usagef("unsupported evaluation data kind %q", kind)
	}
	return validateDataObject(contract.Object, payload, string(kind))
}

func validateDataPayloadForModel(kind DataKind, payload map[string]any, spec *model.Spec) error {
	if err := validateDataPayload(kind, payload); err != nil {
		return err
	}
	if err := validatePayloadModelBindings(kind, payload, spec); err != nil {
		return err
	}
	return validatePayloadSemantics(kind, payload)
}

func validateDataObject(contract dataObjectContract, payload map[string]any, path string) error {
	fields := map[string]dataField{}
	for _, f := range contract.Fields {
		fields[f.Name] = f
		if f.Required {
			if _, ok := payload[f.Name]; !ok {
				return usagef("%s is missing required field %s", path, f.Name)
			}
		}
	}
	if !contract.Open {
		for name := range payload {
			if _, ok := fields[name]; !ok {
				return usagef("%s contains unknown field %s", path, name)
			}
		}
	}
	for _, f := range contract.Fields {
		value, ok := payload[f.Name]
		if !ok {
			continue
		}
		if err := validateDataValue(f, value, path+"."+f.Name); err != nil {
			return err
		}
	}
	return nil
}

//nolint:cyclop,gocognit // The recursive validator mirrors the field type enum.
func validateDataValue(field dataField, value any, path string) error {
	switch field.Type {
	case dataAny:
		return nil
	case dataBool:
		if _, ok := value.(bool); !ok {
			return usagef("%s must be a boolean", path)
		}
	case dataNumber:
		if _, ok := numericSchemaVersion(value); !ok {
			return usagef("%s must be a number", path)
		}
	case dataString, dataRatingLevelID:
		s, ok := value.(string)
		if !ok || s == "" {
			return usagef("%s must be a non-empty string", path)
		}
		if len(field.Enum) > 0 && !slices.Contains(field.Enum, s) {
			return usagef("%s = %q, want one of %s", path, s, strings.Join(field.Enum, ", "))
		}
	case dataStringArray:
		return validateDataStringArray(field, value, path)
	case dataAreaID:
		_, err := areaIDFrom(value)
		if err != nil {
			return usagef("%s: %w", path, err)
		}
	case dataRequirementID:
		_, err := requirementIDFrom(value)
		if err != nil {
			return usagef("%s: %w", path, err)
		}
	case dataFactorID:
		_, err := factorIDFrom(value)
		if err != nil {
			return usagef("%s: %w", path, err)
		}
	case dataObject:
		obj, ok := value.(map[string]any)
		if !ok {
			return usagef("%s must be an object", path)
		}
		if field.Object != nil {
			if err := validateDataObject(*field.Object, obj, path); err != nil {
				return err
			}
		}
	case dataArray:
		return validateDataArray(field, value, path)
	}
	return nil
}

func validateDataStringArray(field dataField, value any, path string) error {
	items, ok := value.([]any)
	if !ok {
		return usagef("%s must be an array", path)
	}
	if err := validateDataMinItems(field, len(items), path); err != nil {
		return err
	}
	for i, item := range items {
		if s, ok := item.(string); !ok || s == "" {
			return usagef("%s[%d] must be a non-empty string", path, i)
		}
	}
	return nil
}

func validateDataArray(field dataField, value any, path string) error {
	items, ok := value.([]any)
	if !ok {
		return usagef("%s must be an array", path)
	}
	if err := validateDataMinItems(field, len(items), path); err != nil {
		return err
	}
	if field.Element == nil {
		return nil
	}
	for i, item := range items {
		if err := validateDataValue(*field.Element, item, fmt.Sprintf("%s[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}

func validateDataMinItems(field dataField, itemCount int, path string) error {
	if field.MinItems > 0 && itemCount < field.MinItems {
		return usagef("%s must contain at least %d item(s)", path, field.MinItems)
	}
	return nil
}

func validatePayloadSemantics(kind DataKind, payload map[string]any) error {
	switch kind {
	case DataKindRequirementAssessment:
		return validateRequirementAssessmentResultSemantics(payload)
	default:
		return nil
	}
}

func validateRequirementAssessmentResultSemantics(payload map[string]any) error {
	for i, finding := range objectSlice(payload["findings"]) {
		path := fmt.Sprintf("RequirementAssessmentResult.findings[%d]", i)
		seen := map[string]struct{}{}
		for j, action := range objectSlice(finding["candidateActions"]) {
			actionPath := fmt.Sprintf("%s.candidateActions[%d]", path, j)
			id := firstString(action, "id")
			if _, ok := seen[id]; ok {
				return usagef("%s.id %q is duplicated within %s.candidateActions", actionPath, id, path)
			}
			seen[id] = struct{}{}
		}
	}
	return nil
}

//nolint:gocognit,nestif // Model binding is centralized to keep write/verify parity.
func validatePayloadModelBindings(kind DataKind, payload map[string]any, spec *model.Spec) error {
	if spec == nil {
		return usagef("%s cannot be model-bound without a model snapshot", kind)
	}
	if err := walkModelReferences(payload, spec, string(kind), nil); err != nil {
		return err
	}
	return nil
}

//nolint:cyclop,gocognit,nestif // Recursive descent keeps model-reference validation uniform.
func walkModelReferences(value any, spec *model.Spec, path string, key *string) error {
	switch v := value.(type) {
	case map[string]any:
		for childKey, child := range v {
			k := childKey
			if err := walkModelReferences(child, spec, path+"."+childKey, &k); err != nil {
				return err
			}
		}
	case []any:
		if key != nil {
			if singular, ok := repeatedModelReferenceKey(*key); ok {
				for i, item := range v {
					if err := validateModelReferenceString(spec, singular, item, fmt.Sprintf("%s[%d]", path, i)); err != nil {
						return err
					}
				}
				return nil
			}
		}
		for i, item := range v {
			if err := walkModelReferences(item, spec, fmt.Sprintf("%s[%d]", path, i), key); err != nil {
				return err
			}
		}
	case string:
		if key != nil {
			if err := validateModelReferenceString(spec, *key, v, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func repeatedModelReferenceKey(key string) (string, bool) {
	switch key {
	case "areaIds", "childAreaIds":
		return "areaId", true
	case "factorIds", "rootFactorIds", "factors":
		return "factorId", true
	case "localRequirementIds":
		return "requirementId", true
	case "ratingLevelIds":
		return "ratingLevelId", true
	default:
		return "", false
	}
}

func validateModelReferenceString(spec *model.Spec, key string, value any, path string) error {
	switch key {
	case "areaId":
		return validateAreaReferenceString(spec, value, path)
	case "requirementId":
		return validateRequirementReferenceString(spec, value, path)
	case "factorId":
		return validateFactorReferenceString(spec, value, path)
	case "ratingLevelId":
		ref, ok := value.(string)
		if !ok {
			return usagef("%s must be a qualified rating reference string", path)
		}
		if _, err := ParseRatingReference(spec, ref); err != nil {
			return err
		}
	}
	return nil
}

func validateAreaReferenceString(spec *model.Spec, value any, path string) error {
	area, err := areaIDFrom(value)
	if err != nil {
		return usagef("%s: %w", path, err)
	}
	if !model.AreaExists(spec, area) {
		return usagef("%s does not resolve in the model", path)
	}
	return nil
}

func validateRequirementReferenceString(spec *model.Spec, value any, path string) error {
	req, err := requirementIDFrom(value)
	if err != nil {
		return usagef("%s: %w", path, err)
	}
	if !model.AreaExists(spec, req.DeclaringArea) {
		return usagef("%s declares an Area absent from the model", path)
	}
	if !model.RequirementExists(spec, req.DeclaringArea, req.Name) {
		return usagef("%s does not resolve in the model", path)
	}
	return nil
}

func validateFactorReferenceString(spec *model.Spec, value any, path string) error {
	factor, err := factorIDFrom(value)
	if err != nil {
		return usagef("%s: %w", path, err)
	}
	if !model.AreaExists(spec, factor.DeclaringArea) {
		return usagef("%s declares an Area absent from the model", path)
	}
	if !model.FactorExists(spec, factor.DeclaringArea, factor.Path) {
		return usagef("%s does not resolve in the model", path)
	}
	return nil
}

func prefixedModelReferencePattern(prefix string) string {
	body := strings.TrimPrefix(strings.TrimSuffix(qschema.ModelNamePattern, "$"), "^")
	return "^" + prefix + "(?:root|" + body + "(?:/" + body + ")*)$"
}

func ratingReferencePattern() string {
	body := strings.TrimPrefix(strings.TrimSuffix(qschema.ModelNamePattern, "$"), "^")
	return "^rating:" + body + "$"
}

func factorReferencePattern() string {
	body := strings.TrimPrefix(strings.TrimSuffix(qschema.ModelNamePattern, "$"), "^")
	path := "(?:root|" + body + "(?:/" + body + ")*)"
	return "^factor:" + path + "::" + body + "(?:/" + body + ")*$"
}

func requirementReferencePattern() string {
	body := strings.TrimPrefix(strings.TrimSuffix(qschema.ModelNamePattern, "$"), "^")
	path := "(?:root|" + body + "(?:/" + body + ")*)"
	return "^requirement:" + path + "::" + body + "$"
}

func EvaluationDataSchema(kind DataKind) ([]byte, error) {
	doc, err := evaluationDataSchemaDoc(kind)
	if err != nil {
		return nil, err
	}
	return canonicalJSON(doc)
}

func evaluationDataSchemaDoc(kind DataKind) (map[string]any, error) {
	if kind != "" {
		contract, ok := dataContracts[kind]
		if !ok {
			return nil, usagef("unknown evaluation data kind %q", kind)
		}
		doc := schemaForObject(contract.Object)
		doc["$schema"] = qschema.JSONSchemaDialect
		doc["$id"] = evaluationDataSchemaID + "/" + string(kind)
		return doc, nil
	}

	defs := map[string]any{}
	for _, k := range dataContractOrder {
		contract := dataContracts[k]
		defs[string(k)] = schemaForObject(contract.Object)
	}
	doc := map[string]any{
		"$schema": qschema.JSONSchemaDialect,
		"$id":     evaluationDataSchemaID,
		"$defs":   defs,
	}
	refs := make([]any, 0, len(dataContractOrder))
	for _, k := range dataContractOrder {
		refs = append(refs, map[string]any{"$ref": "#/$defs/" + string(k)})
	}
	doc["oneOf"] = refs
	return doc, nil
}

func schemaForObject(contract dataObjectContract) map[string]any {
	props := map[string]any{}
	var required []any
	for _, f := range contract.Fields {
		props[f.Name] = schemaForField(f)
		if f.Required {
			required = append(required, f.Name)
		}
	}
	out := map[string]any{"type": "object", "properties": props}
	if len(required) > 0 {
		out["required"] = required
	}
	if !contract.Open {
		out["additionalProperties"] = false
	}
	return out
}

//nolint:cyclop // Schema generation follows the same field type enum as validation.
func schemaForField(field dataField) map[string]any {
	switch field.Type {
	case dataAny:
		return map[string]any{}
	case dataBool:
		return map[string]any{"type": "boolean"}
	case dataNumber:
		return map[string]any{"type": "integer"}
	case dataString, dataRatingLevelID:
		out := map[string]any{"type": "string"}
		if len(field.Enum) > 0 {
			out["enum"] = anyStrings(field.Enum)
		}
		if field.Type == dataRatingLevelID {
			out["pattern"] = ratingReferencePattern()
		}
		return out
	case dataStringArray:
		out := map[string]any{"type": "array", "items": map[string]any{"type": "string"}}
		if field.MinItems > 0 {
			out["minItems"] = field.MinItems
		}
		return out
	case dataAreaID:
		return map[string]any{"type": "string", "pattern": prefixedModelReferencePattern("area:")}
	case dataRequirementID:
		return map[string]any{"type": "string", "pattern": requirementReferencePattern()}
	case dataFactorID:
		return map[string]any{"type": "string", "pattern": factorReferencePattern()}
	case dataObject:
		if field.Object == nil {
			return map[string]any{"type": "object"}
		}
		return schemaForObject(*field.Object)
	case dataArray:
		items := map[string]any{}
		if field.Element != nil {
			items = schemaForField(*field.Element)
		}
		out := map[string]any{"type": "array", "items": items}
		if field.MinItems > 0 {
			out["minItems"] = field.MinItems
		}
		return out
	default:
		return map[string]any{}
	}
}

func VerifyData(runPath string) (*DataVerifyReceipt, error) {
	spec, err := loadRunModel(runPath)
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", ModelSnapshotFile, err)
	}
	result := &DataVerifyReceipt{SchemaVersion: SchemaVersion, Path: filepath.ToSlash(runPath), Valid: true}
	root := filepath.Join(runPath, "data")
	payloads := map[string]map[string]any{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		rel, err := filepath.Rel(runPath, path)
		if err != nil {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			result.addFailure(rel, "", err.Error())
			return nil
		}
		payload, err := decodeDataPayload(raw)
		if err != nil {
			result.addFailure(rel, "", err.Error())
			return nil
		}
		kind, err := payloadKind(payload)
		if err != nil {
			result.addFailure(rel, "", err.Error())
			return nil
		}
		if _, ok := dataContracts[kind]; !ok {
			result.addFailure(rel, kind, fmt.Sprintf("unsupported evaluation data kind %q", kind))
			return nil
		}
		result.Checked++
		if err := validateDataPayloadForModel(kind, payload, spec); err != nil {
			result.addFailure(rel, kind, err.Error())
		}
		payloads[filepath.ToSlash(rel)] = payload
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("verifying evaluation data: %w", err)
	}
	for _, failure := range validateEffectivePayloads(payloads) {
		result.addFailure(failure.Path, failure.Kind, failure.Reason)
	}
	result.Valid = len(result.Failures) == 0
	return result, nil
}

func (r *DataVerifyReceipt) addFailure(rel string, kind DataKind, reason string) {
	r.Failures = append(r.Failures, DataVerifyFailure{Path: filepath.ToSlash(rel), Kind: kind, Reason: reason})
}

func (kind DataKind) String() string {
	return string(kind)
}
