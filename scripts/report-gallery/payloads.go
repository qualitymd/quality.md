package main

import (
	"fmt"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

var ratingLevelIDs = []any{"rating:outstanding", "rating:target", "rating:minimum", "rating:unacceptable"}

func galleryPayloads() []map[string]any {
	var payloads []map[string]any
	payloads = append(payloads, map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindEvaluationFrame),
		"subject":       map[string]any{"modelLocator": "QUALITY.md"},
		"inputs":        map[string]any{"ratingLevelIds": ratingLevelIDs},
		"derivedContext": map[string]any{
			"rigor": "full-evaluation",
			"evaluationPolicies": []any{
				"Read ratings worst-of within an area; balance-invariants is the named veto requirement.",
				"A not-assessed result stays visible as missing evidence and never counts as a low rating.",
			},
			"expectedEvaluationLimits": []any{map[string]any{
				"id":          "synthetic-source-omitted",
				"description": "This is an illustrative gallery evaluation: routine outputs are synthetic and the concrete source system behind the synthetic-source references is omitted.",
				"impact":      "Evidence references demonstrate report shape rather than checkable source claims.",
			}},
		},
	})
	payloads = append(payloads, areaFrames()...)
	payloads = append(payloads, requirementPayloads()...)
	payloads = append(payloads, factorPayloads()...)
	payloads = append(payloads, areaAnalyses()...)
	payloads = append(payloads, advicePayloads()...)
	return payloads
}

func areaFrames() []map[string]any {
	var payloads []map[string]any
	rootInputs := map[string]any{"sourceRefs": []any{"synthetic-source:ledgerlite-service"}}
	for _, req := range requirements {
		if req.Area == "" {
			inputsArrayAppend(rootInputs, "localRequirementIds", reqRef("", req.Name))
		}
	}
	inputsArrayAppend(rootInputs, "rootFactorIds", factorRef("", "agent-harnessability"))
	for _, area := range areas {
		inputsArrayAppend(rootInputs, "childAreaIds", areaRef(area.Name))
	}
	payloads = append(payloads, map[string]any{
		"schemaVersion":  evaluation.SchemaVersion,
		"kind":           string(evaluation.DataKindAreaEvaluationFrame),
		"subject":        map[string]any{"areaId": areaRef("")},
		"inputs":         rootInputs,
		"derivedContext": map[string]any{"scope": "The composite root: model-wide agent harnessability requirements plus six constituent areas."},
	})
	for _, area := range areas {
		inputs := map[string]any{"sourceRefs": []any{area.Source}}
		for _, req := range requirements {
			if req.Area == area.Name {
				inputsArrayAppend(inputs, "localRequirementIds", reqRef(area.Name, req.Name))
			}
		}
		for _, f := range factors {
			if f.Area == area.Name && !isSubFactorPath(f.Path) {
				inputsArrayAppend(inputs, "rootFactorIds", factorRef(area.Name, f.Path))
			}
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":  evaluation.SchemaVersion,
			"kind":           string(evaluation.DataKindAreaEvaluationFrame),
			"subject":        map[string]any{"areaId": areaRef(area.Name)},
			"inputs":         inputs,
			"derivedContext": map[string]any{"scope": area.Summary},
		})
	}
	return payloads
}

func requirementPayloads() []map[string]any {
	var payloads []map[string]any
	for _, req := range requirements {
		payloads = append(payloads,
			requirementFrame(req),
			requirementAssessment(req),
			requirementRating(req),
		)
	}
	return payloads
}

func requirementFrame(req requirementCase) map[string]any {
	var factorIDs []any
	for _, f := range req.Factors {
		factorIDs = append(factorIDs, factorRef(req.Area, f))
	}
	var criteria []any
	for _, c := range req.AppliedCriteria {
		criteria = append(criteria, map[string]any{"ratingLevelId": "rating:" + c.Level, "criterion": c.Text})
	}
	var sourceRefs []any
	for _, ref := range req.EvidenceRefs {
		sourceRefs = append(sourceRefs, ref)
	}
	return map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindRequirementEvaluationFrame),
		"subject": map[string]any{
			"requirementId": reqRef(req.Area, req.Name),
			"factorIds":     factorIDs,
		},
		"inputs": map[string]any{
			"ratingLevelIds":             ratingLevelIDs,
			"requirementAssessmentBasis": req.EvidenceQuestion,
		},
		"derivedContext": map[string]any{
			"evidenceTargets": []any{map[string]any{
				"id":         "primary-evidence",
				"question":   req.EvidenceQuestion,
				"sourceRefs": sourceRefs,
				"required":   true,
			}},
			"appliedRatingCriteria": criteria,
		},
	}
}

func buildFinding(req requirementCase, f findingCase) map[string]any {
	criterion := map[string]any{
		"requirementId": reqRef(req.Area, req.Name),
		"ratingLevelId": "rating:" + f.CriterionLevel,
		"criterion":     f.Criterion,
	}
	if f.CriterionRationale != "" {
		criterion["rationale"] = f.CriterionRationale
	}
	finding := map[string]any{
		"id":         f.ID,
		"type":       f.Type,
		"confidence": f.Confidence,
		"statement":  f.Statement,
		"condition":  f.Condition,
		"criteria":   []any{criterion},
		"basis": map[string]any{
			"status":    f.BasisStatus,
			"statement": f.Basis,
		},
		"effect": map[string]any{
			"statement":    f.Effect,
			"ratingEffect": f.RatingEffect,
		},
		"evidence": []any{map[string]any{
			"sourceRef": f.EvidenceRef,
			"statement": f.Evidence,
		}},
	}
	if f.Type == "gap" || f.Type == "risk" {
		finding["severity"] = f.Severity
	}
	return finding
}

func requirementAssessment(req requirementCase) map[string]any {
	var findings []any
	for _, f := range req.Findings {
		findings = append(findings, buildFinding(req, f))
	}
	if findings == nil {
		findings = []any{}
	}
	var factorIDs []any
	for _, f := range req.Factors {
		factorIDs = append(factorIDs, factorRef(req.Area, f))
	}
	payload := map[string]any{
		"schemaVersion":    evaluation.SchemaVersion,
		"kind":             string(evaluation.DataKindRequirementAssessment),
		"requirementId":    reqRef(req.Area, req.Name),
		"status":           req.AssessStatus,
		"evidenceSummary":  req.Summary,
		"summary":          req.Summary,
		"factors":          factorIDs,
		"findings":         findings,
		"confidence":       req.Confidence,
		"confidenceReason": req.ConfidenceReason,
	}
	if req.StatusReason != "" {
		payload["statusReason"] = req.StatusReason
	}
	if limits := limitList(req.Limits); limits != nil {
		payload["evaluationLimits"] = limits
	}
	return payload
}

func requirementRating(req requirementCase) map[string]any {
	payload := map[string]any{
		"schemaVersion":    evaluation.SchemaVersion,
		"kind":             string(evaluation.DataKindRequirementRating),
		"requirementId":    reqRef(req.Area, req.Name),
		"status":           req.RatingStatus,
		"rationale":        req.RatingRationale,
		"confidence":       req.Confidence,
		"confidenceReason": req.ConfidenceReason,
	}
	if req.RatingStatus != "rated" {
		payload["statusReason"] = req.StatusReason
		if missing := unknownList(req.MissingEvidence); missing != nil {
			payload["missingEvidence"] = missing
		}
		return payload
	}
	payload["ratingLevelId"] = "rating:" + req.Rating
	var drivers []any
	for _, f := range req.Findings {
		if !f.Driver {
			continue
		}
		drivers = append(drivers, map[string]any{
			"description":   f.Statement,
			"effect":        f.RatingEffect,
			"ratingLevelId": "rating:" + req.Rating,
			"inputRefs":     []any{findingRefByID(f.ID)},
		})
	}
	payload["ratingDrivers"] = drivers
	var results []any
	for _, r := range req.CriteriaResults {
		results = append(results, map[string]any{
			"ratingLevelId": "rating:" + r.Level,
			"matched":       r.Matched,
			"rationale":     r.Rationale,
		})
	}
	payload["criteriaResults"] = results
	return payload
}

func factorPayloads() []map[string]any {
	var payloads []map[string]any
	for _, f := range factors {
		factorID := factorRef(f.Area, f.Path)
		inputs := map[string]any{}
		for _, req := range requirements {
			if req.Area != f.Area {
				continue
			}
			for _, path := range req.Factors {
				if path == f.Path {
					inputsArrayAppend2(inputs, "directRequirementRatingRefs",
						routineRef(evaluation.DataKindRequirementRating, "requirementId", reqRef(req.Area, req.Name), ""))
				}
			}
		}
		for _, child := range f.Children {
			inputsArrayAppend2(inputs, "childFactorAnalysisRefs",
				routineRef(evaluation.DataKindFactorAnalysis, "factorId", factorRef(f.Area, child), "localAndDescendantAnalysis"))
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":  evaluation.SchemaVersion,
			"kind":           string(evaluation.DataKindFactorAnalysisFrame),
			"subject":        map[string]any{"factorId": factorID},
			"inputs":         inputs,
			"derivedContext": map[string]any{"synthesisGuidanceRef": "model-body/model-shape-and-how-to-read-ratings"},
		})
		payloads = append(payloads, map[string]any{
			"schemaVersion":              evaluation.SchemaVersion,
			"kind":                       string(evaluation.DataKindFactorAnalysis),
			"factorId":                   factorID,
			"localAnalysis":              factorLocalScope(f),
			"localAndDescendantAnalysis": factorAggregateScope(f),
		})
	}
	return payloads
}

func factorLocalScope(f factorCase) map[string]any {
	if len(f.Children) > 0 {
		return map[string]any{
			"status":       "empty",
			"statusReason": "The umbrella factor carries no direct requirements; its rating comes from sub-factor roll-up.",
			"confidence":   "none",
		}
	}
	return factorAggregateScope(f)
}

func factorAggregateScope(f factorCase) map[string]any {
	if f.Status == "blocked" {
		scope := map[string]any{
			"status":       "blocked",
			"statusReason": f.StatusReason,
			"rationale":    f.Summary,
			"confidence":   f.Confidence,
		}
		if limits := limitList(f.Limits); limits != nil {
			scope["evaluationLimits"] = limits
		}
		return scope
	}
	drivers := directRequirementDrivers(f)
	if len(f.Children) > 0 {
		drivers = childFactorDrivers(f)
	}
	return analyzedScope(f.Rating, f.Confidence, f.Summary, drivers)
}

func childFactorDrivers(f factorCase) []any {
	var drivers []any
	for _, child := range f.Children {
		childCase := factorByPath(f.Area, child)
		effect := "rates " + childCase.Rating
		if childCase.Status == "blocked" {
			effect = "blocked; contributes missing evidence"
		}
		drivers = append(drivers, driver(childCase.Title+": "+childCase.Summary, effect,
			routineRef(evaluation.DataKindFactorAnalysis, "factorId", factorRef(f.Area, child), "localAndDescendantAnalysis")))
	}
	return drivers
}

func directRequirementDrivers(f factorCase) []any {
	var drivers []any
	for _, req := range requirements {
		if req.Area != f.Area {
			continue
		}
		for _, path := range req.Factors {
			if path != f.Path {
				continue
			}
			headline := req.RatingRationale
			effect := "not rated"
			if len(req.Findings) > 0 {
				headline = req.Findings[0].Statement
				effect = req.Findings[0].RatingEffect
			}
			drivers = append(drivers, driver(headline, effect,
				routineRef(evaluation.DataKindRequirementRating, "requirementId", reqRef(req.Area, req.Name), "")))
		}
	}
	return drivers
}

func areaAnalyses() []map[string]any {
	var payloads []map[string]any
	for _, area := range areas {
		var refs []any
		var drivers []any
		for _, f := range factors {
			if f.Area != area.Name || isSubFactorPath(f.Path) {
				continue
			}
			ref := routineRef(evaluation.DataKindFactorAnalysis, "factorId", factorRef(f.Area, f.Path), "localAndDescendantAnalysis")
			refs = append(refs, ref)
			effect := "rates " + f.Rating
			if f.Status == "blocked" {
				effect = "blocked; contributes missing evidence"
			}
			drivers = append(drivers, driver(f.Title+": "+f.Summary, effect, ref))
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":  evaluation.SchemaVersion,
			"kind":           string(evaluation.DataKindAreaAnalysisFrame),
			"subject":        map[string]any{"areaId": areaRef(area.Name)},
			"inputs":         map[string]any{"factorAnalysisRefs": refs},
			"derivedContext": map[string]any{"synthesisGuidanceRef": "model-body/model-shape-and-how-to-read-ratings"},
		})
		scope := analyzedScope(area.Rating, area.Confidence, area.Summary, drivers)
		if limits := limitList(area.Limits); limits != nil {
			scope["evaluationLimits"] = limits
		}
		if missing := limitList(area.MissingEvidence); missing != nil {
			scope["incompleteInputs"] = missing
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion":              evaluation.SchemaVersion,
			"kind":                       string(evaluation.DataKindAreaAnalysis),
			"areaId":                     areaRef(area.Name),
			"localAnalysis":              scope,
			"localAndDescendantAnalysis": scope,
		})
	}
	payloads = append(payloads, rootAnalysisPayloads()...)
	return payloads
}

func rootAnalysisPayloads() []map[string]any {
	umbrellaRef := routineRef(evaluation.DataKindFactorAnalysis, "factorId", factorRef("", "agent-harnessability"), "localAndDescendantAnalysis")
	var childRefs []any
	var aggregateDrivers []any
	aggregateDrivers = append(aggregateDrivers, driver(
		"Agent Harnessability: "+factorByPath("", "agent-harnessability").Summary,
		"rates minimum", umbrellaRef))
	for _, area := range areas {
		ref := routineRef(evaluation.DataKindAreaAnalysis, "areaId", areaRef(area.Name), "localAndDescendantAnalysis")
		childRefs = append(childRefs, ref)
		aggregateDrivers = append(aggregateDrivers, driver(area.Summary, "rates "+area.Rating, ref))
	}
	frame := map[string]any{
		"schemaVersion": evaluation.SchemaVersion,
		"kind":          string(evaluation.DataKindAreaAnalysisFrame),
		"subject":       map[string]any{"areaId": areaRef("")},
		"inputs": map[string]any{
			"factorAnalysisRefs":    []any{umbrellaRef},
			"childAreaAnalysisRefs": childRefs,
		},
		"derivedContext": map[string]any{"synthesisGuidanceRef": "model-body/model-shape-and-how-to-read-ratings"},
	}
	localScope := analyzedScope(rootLocalAnalysis.Rating, rootLocalAnalysis.Confidence, rootLocalAnalysis.Summary,
		[]any{driver("Agent Harnessability sub-factor roll-up over the root's own requirements.", "rates minimum", umbrellaRef)})
	aggregateScope := analyzedScope(rootAggregateAnalysis.Rating, rootAggregateAnalysis.Confidence, rootAggregateAnalysis.Summary, aggregateDrivers)
	aggregateScope["evaluationLimits"] = []any{map[string]any{
		"id":          "below-required-margin",
		"description": "The body requires money-touching areas (api, service-contract, persistence) to land at target or better; all three sit at minimum.",
		"impact":      "Read the overall minimum as a stop-and-fix signal per the model's required margin, not as an acceptable steady state.",
	}}
	analysis := map[string]any{
		"schemaVersion":              evaluation.SchemaVersion,
		"kind":                       string(evaluation.DataKindAreaAnalysis),
		"areaId":                     areaRef(""),
		"localAnalysis":              localScope,
		"localAndDescendantAnalysis": aggregateScope,
	}
	return []map[string]any{frame, analysis}
}

func advicePayloads() []map[string]any {
	var ordered []any
	for i, ranked := range rankedFindings {
		ordered = append(ordered, map[string]any{
			"rank":       i + 1,
			"findingRef": findingRefByID(ranked.FindingID),
			"tier":       ranked.Tier,
			"rationale":  ranked.Rationale,
		})
	}
	payloads := []map[string]any{{
		"schemaVersion":   evaluation.SchemaVersion,
		"kind":            string(evaluation.DataKindFindingRanking),
		"orderedFindings": ordered,
		"rationale":       "Ranked by exposure of the money path first, then by what blocks judgment, then by agent-workflow drag; strengths anchor the tail.",
	}}
	addressedBy := map[string][]string{}
	for _, rec := range recommendations {
		var traces []any
		for _, findingID := range rec.Traces {
			traces = append(traces, findingRefByID(findingID))
			addressedBy[findingID] = append(addressedBy[findingID], rec.ID)
		}
		payloads = append(payloads, map[string]any{
			"schemaVersion": evaluation.SchemaVersion,
			"kind":          string(evaluation.DataKindRecommendation),
			"id":            rec.ID,
			"title":         rec.Title,
			"description":   rec.Description,
			"background":    rec.Background,
			"expectedValue": rec.ExpectedValue,
			"doneCriterion": rec.DoneCriterion,
			"impact":        rec.Impact,
			"confidence":    rec.Confidence,
			"traceRefs":     traces,
		})
	}
	var rankedRecs []any
	for i, rec := range recommendations {
		rankedRecs = append(rankedRecs, map[string]any{
			"rank":              i + 1,
			"recommendationRef": rec.ID,
			"impact":            rec.Impact,
			"confidence":        rec.Confidence,
			"rationale":         "Ordered by expected effect on the required margin: money-path gaps first, then enforcement, then assessability and workflow improvements.",
		})
	}
	var coverage []any
	for _, ranked := range rankedFindings {
		entry := map[string]any{"findingRef": findingRefByID(ranked.FindingID)}
		if recIDs, ok := addressedBy[ranked.FindingID]; ok {
			var refs []any
			for _, id := range recIDs {
				refs = append(refs, id)
			}
			entry["disposition"] = "addressed_by_recommendation"
			entry["recommendationRefs"] = refs
			entry["rationale"] = "Addressed by the traced recommendation."
		} else {
			entry["disposition"] = "not_advice_driving"
			entry["rationale"] = ranked.Rationale
		}
		coverage = append(coverage, entry)
	}
	payloads = append(payloads, map[string]any{
		"schemaVersion":          evaluation.SchemaVersion,
		"kind":                   string(evaluation.DataKindRecommendationRanking),
		"orderedRecommendations": rankedRecs,
		"findingCoverage":        coverage,
		"rationale":              "Recommendations target the findings that hold money-touching areas below the required margin, restore blocked judgment, and harden the agent loop.",
	})
	return payloads
}

func analyzedScope(rating, confidence, rationale string, drivers []any) map[string]any {
	return map[string]any{
		"status":        "analyzed",
		"ratingLevelId": "rating:" + rating,
		"rationale":     rationale,
		"ratingDrivers": drivers,
		"confidence":    confidence,
	}
}

func driver(description, effect string, ref map[string]any) map[string]any {
	return map[string]any{
		"description": description,
		"effect":      effect,
		"inputRefs":   []any{ref},
	}
}

func limitList(limits []limitCase) []any {
	var out []any
	for _, l := range limits {
		out = append(out, map[string]any{
			"id":          l.ID,
			"description": l.Description,
			"impact":      l.Impact,
		})
	}
	return out
}

func unknownList(limits []limitCase) []any {
	var out []any
	for _, l := range limits {
		out = append(out, map[string]any{
			"id":          l.ID,
			"description": l.Description,
			"impact":      l.Impact,
		})
	}
	return out
}

func findingRefByID(findingID string) map[string]any {
	for _, req := range requirements {
		for _, f := range req.Findings {
			if f.ID == findingID {
				return routineRef(evaluation.DataKindRequirementAssessment, "requirementId",
					reqRef(req.Area, req.Name), "findings["+findingID+"]")
			}
		}
	}
	panic(fmt.Sprintf("unknown finding id %q", findingID))
}

func factorByPath(area, path string) factorCase {
	for _, f := range factors {
		if f.Area == area && f.Path == path {
			return f
		}
	}
	panic(fmt.Sprintf("unknown factor %s::%s", area, path))
}

func isSubFactorPath(path string) bool {
	for _, r := range path {
		if r == '/' {
			return true
		}
	}
	return false
}

func routineRef(kind evaluation.DataKind, subjectKey, id, selector string) map[string]any {
	ref := map[string]any{
		"kind":    string(kind),
		"subject": map[string]any{subjectKey: id},
	}
	if selector != "" {
		ref["selector"] = selector
	}
	return ref
}

func areaRef(area string) string {
	if area == "" {
		return "area:root"
	}
	return "area:" + area
}

func factorRef(area, factorPath string) string {
	if area == "" {
		return "factor:root::" + factorPath
	}
	return "factor:" + area + "::" + factorPath
}

func reqRef(area, req string) string {
	if area == "" {
		return "requirement:root::" + req
	}
	return "requirement:" + area + "::" + req
}

func inputsArrayAppend(target map[string]any, key string, value string) {
	values, _ := target[key].([]any)
	target[key] = append(values, value)
}

func inputsArrayAppend2(target map[string]any, key string, value map[string]any) {
	values, _ := target[key].([]any)
	target[key] = append(values, value)
}
