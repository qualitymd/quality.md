package runner

import (
	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/model"
)

// Deterministic frame generation. Frames are structural derivations of the
// model snapshot and planned scope, so the runner produces them itself; only
// judgment work goes to evaluators.

const (
	factorSynthesisGuidance = "protocol:factor-synthesis-default-v0"
	areaSynthesisGuidance   = "protocol:area-synthesis-default-v0"
)

func framePayload(kind evaluation.DataKind, body map[string]any) map[string]any {
	payload := map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(kind),
	}
	for key, value := range body {
		payload[key] = value
	}
	return payload
}

func ratingLevelRefs(spec *model.Spec) []any {
	out := make([]any, 0, len(spec.RatingScale))
	for _, level := range spec.RatingScale {
		out = append(out, evaluation.RatingReference(level.Level))
	}
	return out
}

func evaluationFramePayload(spec *model.Spec, manifest Manifest) map[string]any {
	return framePayload(evaluation.DataKindEvaluationFrame, map[string]any{
		"subject": map[string]any{"modelLocator": manifest.Model},
		"inputs":  map[string]any{"ratingLevelIds": ratingLevelRefs(spec)},
		"derivedContext": map[string]any{
			"rigor":              "standard",
			"evaluationPolicies": []any{"source-as-data", "secret-redaction"},
		},
	})
}

func areaEvaluationFramePayload(area *PlannedArea) map[string]any {
	inputs := map[string]any{
		"localRequirementIds": anyStrings(area.LocalRequirements),
		"rootFactorIds":       anyStrings(area.RootFactors),
		"childAreaIds":        anyStrings(area.ChildAreas),
	}
	if area.Source != "" {
		inputs["sourceRefs"] = []any{area.Source}
	}
	return framePayload(evaluation.DataKindAreaEvaluationFrame, map[string]any{
		"subject": map[string]any{"areaId": area.Ref},
		"inputs":  inputs,
	})
}

func requirementEvaluationFramePayload(spec *model.Spec, req *PlannedRequirement) map[string]any {
	subject := map[string]any{"requirementId": req.Ref}
	if len(req.Factors) > 0 {
		subject["factorIds"] = anyStrings(req.Factors)
	}
	inputs := map[string]any{"ratingLevelIds": ratingLevelRefs(spec)}
	if req.Assessment != "" {
		inputs["requirementAssessmentBasis"] = req.Assessment
	}
	criteria := make([]any, 0, len(spec.RatingScale))
	for _, level := range spec.RatingScale {
		criterion := level.Criterion
		source := "model_default"
		if override, ok := req.Ratings[level.Level]; ok && override != "" {
			criterion = override
			source = "requirement_override"
		}
		if criterion == "" {
			continue
		}
		criteria = append(criteria, map[string]any{
			"ratingLevelId": evaluation.RatingReference(level.Level),
			"criterion":     criterion,
			"source":        source,
		})
	}
	body := map[string]any{"subject": subject, "inputs": inputs}
	if len(criteria) > 0 {
		body["derivedContext"] = map[string]any{"appliedRatingCriteria": criteria}
	}
	return framePayload(evaluation.DataKindRequirementEvaluationFrame, body)
}

func factorAnalysisFramePayload(factor *PlannedFactor) map[string]any {
	directRefs := make([]any, 0, len(factor.Requirements))
	for _, reqRef := range factor.Requirements {
		directRefs = append(directRefs, routineRefPayload("RequirementRatingResult", map[string]any{"requirementId": reqRef}, ""))
	}
	childRefs := make([]any, 0, len(factor.ChildFactors))
	for _, childRef := range factor.ChildFactors {
		childRefs = append(childRefs, routineRefPayload("FactorAnalysisResult", map[string]any{"factorId": childRef}, "localAndDescendantAnalysis"))
	}
	return framePayload(evaluation.DataKindFactorAnalysisFrame, map[string]any{
		"subject": map[string]any{"areaId": factor.Area.Reference(), "factorId": factor.Ref},
		"inputs": map[string]any{
			"directRequirementRatingRefs": directRefs,
			"childFactorAnalysisRefs":     childRefs,
		},
		"derivedContext": map[string]any{
			"synthesisGuidanceRef": factorSynthesisGuidance,
			"emptySignalPolicy":    "ignore_empty",
		},
	})
}

func areaAnalysisFramePayload(area *PlannedArea) map[string]any {
	factorRefs := make([]any, 0, len(area.RootFactors))
	for _, factorRef := range area.RootFactors {
		factorRefs = append(factorRefs, routineRefPayload("FactorAnalysisResult", map[string]any{"factorId": factorRef}, "localAndDescendantAnalysis"))
	}
	childRefs := make([]any, 0, len(area.ChildAreas))
	for _, childRef := range area.ChildAreas {
		childRefs = append(childRefs, routineRefPayload("AreaAnalysisResult", map[string]any{"areaId": childRef}, "localAndDescendantAnalysis"))
	}
	return framePayload(evaluation.DataKindAreaAnalysisFrame, map[string]any{
		"subject": map[string]any{"areaId": area.Ref},
		"inputs": map[string]any{
			"factorAnalysisRefs":    factorRefs,
			"childAreaAnalysisRefs": childRefs,
		},
		"derivedContext": map[string]any{
			"synthesisGuidanceRef": areaSynthesisGuidance,
			"emptySignalPolicy":    "ignore_empty",
		},
	})
}

func routineRefPayload(kind string, subject map[string]any, selector string) map[string]any {
	ref := map[string]any{"kind": kind, "subject": subject}
	if selector != "" {
		ref["selector"] = selector
	}
	return ref
}
