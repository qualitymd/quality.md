package runner

import (
	"encoding/json"
	"fmt"

	"github.com/qualitymd/quality.md/internal/evaluation"
	"github.com/qualitymd/quality.md/internal/evaluator"
)

// Work-request assembly. The runner owns every prompt: evaluators receive one
// bounded request with the frames, prior results, rating vocabulary, source
// bundle, and expected schema the judgment needs — nothing else.

func (e *engine) buildWorkRequest(unit *Unit) (evaluator.WorkRequest, error) {
	req := evaluator.WorkRequest{
		RunID:         e.artifact.Manifest.EvaluationID,
		WorkUnitID:    unit.ID,
		Kind:          string(unit.Kind),
		Subject:       unit.Subject,
		CorrelationID: fmt.Sprintf("%s#%s", e.artifact.Manifest.EvaluationID, unit.ID),
	}
	schema, err := expectedSchema(unit)
	if err != nil {
		return req, err
	}
	req.ExpectedSchema = schema

	switch unit.Kind {
	case KindResolveSource:
		e.fillResolveSourceRequest(unit, &req)
	case KindAssessRateRequirement:
		if err := e.fillAssessRateRequest(unit, &req); err != nil {
			return req, err
		}
	case KindAnalyzeFactor:
		e.fillFactorAnalysisRequest(unit, &req)
	case KindAnalyzeArea:
		e.fillAreaAnalysisRequest(unit, &req)
	case KindRankFindings:
		e.fillRankFindingsRequest(&req)
	case KindRecommend:
		e.fillRecommendRequest(&req)
	case KindRankRecommendations:
		e.fillRankRecommendationsRequest(&req)
	default:
		return req, fmt.Errorf("work unit %s is not evaluator-backed", unit.ID)
	}
	req.SourcePackageHash = bundleHash(req.Source)
	// The prefix hash covers the whole cacheable stable region — instructions,
	// schema, packaged source, and shared context — so the log records what a
	// provider prefix cache can actually reuse across an area's work units.
	req.PromptPrefixHash = hashJSON(map[string]any{
		"instructions":  req.Instructions,
		"schema":        string(req.ExpectedSchema),
		"sourcePackage": req.SourcePackageHash,
		"sharedContext": req.SharedContext,
	})
	return req, nil
}

func expectedSchema(unit *Unit) (json.RawMessage, error) {
	if unit.Kind == KindResolveSource {
		resolved := map[string]any{
			"type": "object",
			"properties": map[string]any{
				"files": map[string]any{
					"type":     "array",
					"minItems": 1,
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"path":    map[string]any{"type": "string", "minLength": 1},
							"content": map[string]any{"type": "string"},
						},
						"required":             []any{"path", "content"},
						"additionalProperties": false,
					},
				},
			},
			"required":             []any{"files"},
			"additionalProperties": false,
		}
		return json.MarshalIndent(resolved, "", "  ")
	}
	if unit.Kind == KindAssessRateRequirement {
		assessment, err := evaluation.EvaluationDataSchema(evaluation.DataKindRequirementAssessment)
		if err != nil {
			return nil, err
		}
		rating, err := evaluation.EvaluationDataSchema(evaluation.DataKindRequirementRating)
		if err != nil {
			return nil, err
		}
		composite := map[string]any{
			"type": "object",
			"properties": map[string]any{
				"assessment": json.RawMessage(assessment),
				"rating":     json.RawMessage(rating),
			},
			"required":             []any{"assessment", "rating"},
			"additionalProperties": false,
		}
		return json.MarshalIndent(composite, "", "  ")
	}
	if unit.Kind == KindRecommend {
		item, err := evaluation.EvaluationDataSchema(evaluation.DataKindRecommendation)
		if err != nil {
			return nil, err
		}
		composite := map[string]any{
			"type": "object",
			"properties": map[string]any{
				"recommendations": map[string]any{
					"type":     "array",
					"minItems": 1,
					"items":    json.RawMessage(item),
				},
			},
			"required":             []any{"recommendations"},
			"additionalProperties": false,
		}
		return json.MarshalIndent(composite, "", "  ")
	}
	raw, err := evaluation.EvaluationDataSchema(unit.DataKind)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

// fillResolveSourceRequest assembles a source resolution request: gather the
// material the area's selector describes and return it as files. The request
// carries the selector, its pinned kind, and the area frame — and an empty
// source bundle: the resolver is fed a description, never pre-gathered
// evidence, so the gatherer and the judge are never the same uncontrolled
// step.
func (e *engine) fillResolveSourceRequest(unit *Unit, req *evaluator.WorkRequest) {
	record := e.artifact.Sources[unit.Subject]
	areaFrame := e.payloadFor(unitID(KindFrameAreaEvaluation, unit.Subject))
	req.SharedContext = map[string]any{
		"areaEvaluationFrame": areaFrame,
	}
	req.Context = map[string]any{
		"sourceSelector": map[string]any{
			"selector": record.Selector,
			"kind":     record.Kind,
		},
	}
	req.Instructions = "Resolve this area's source selector: gather the material the selector describes and return " +
		"one JSON object of the form {\"files\": [{\"path\": string, \"content\": string}, ...]}.\n" +
		"- The selector describes a body of evidence; use your tools to locate exactly the material it names and " +
		"return it as text files.\n" +
		"- path is a stable label for each gathered item — a repo-relative path, ticket ID, or URL; paths must be " +
		"unique and non-empty.\n" +
		"- content is the gathered material itself, verbatim; do not summarize, assess, or rate it — a separate " +
		"judgment request evaluates the captured evidence.\n" +
		"- Gather only what the selector describes; do not widen to adjacent material.\n" +
		"- If the material the selector describes does not exist — including when the selector reads like a " +
		"filesystem path that names nothing — return the classified failure source_unavailable naming the selector " +
		"instead of improvising or substituting evidence."
}

// fillAssessRateRequest assembles the combined requirement judgment request:
// one evaluator call assesses the requirement and rates it from that
// assessment. Everything stable across the area's requirements — the packaged
// source and the area frame — goes to the shared context so the rendered
// prompt prefix repeats verbatim.
func (e *engine) fillAssessRateRequest(unit *Unit, req *evaluator.WorkRequest) error {
	planned := e.graph.Plan.Requirement(unit.Subject)
	frame := e.payloadFor(unitID(KindFrameRequirementEvaluation, unit.Subject))
	areaFrame := e.payloadFor(unitID(KindFrameAreaEvaluation, planned.Area.Reference()))
	bundle, err := e.areaSourceBundle(planned.Area.Reference())
	if err != nil {
		return err
	}
	req.Source = bundle.Files
	req.SharedContext = map[string]any{
		"areaEvaluationFrame": areaFrame,
	}
	if len(bundle.Missing) > 0 {
		req.SharedContext["unavailableSource"] = anyStrings(bundle.Missing)
	}
	req.Context = map[string]any{
		"requirement": map[string]any{
			"requirementId": planned.Ref,
			"title":         planned.Title,
			"description":   planned.Description,
			"assessment":    planned.Assessment,
		},
		"requirementEvaluationFrame": frame,
	}
	req.Instructions = "Assess this requirement against the packaged source evidence, then rate it from that " +
		"assessment, and return one JSON object of the form {\"assessment\": RequirementAssessmentResult, " +
		"\"rating\": RequirementRatingResult}.\n" +
		"- Set requirementId in both objects to the subject reference exactly.\n" +
		"Assessment:\n" +
		"- status is one of: assessed, partially_assessed, blocked, not_applicable.\n" +
		"- Record every finding with the full core shape (id, type, confidence, statement, condition, " +
		"criteria, basis, effect, evidence). Gap and risk findings carry severity; strength and note findings must not.\n" +
		"- criteria entries reference this requirement and a rating level from the frame's appliedRatingCriteria.\n" +
		"- Cite evidence sourceRef values from the packaged source paths.\n" +
		"- If required evidence is unavailable, say so via status, unknowns, and evaluationLimits instead of guessing.\n" +
		"- Use finding ids like gap-001, strength-001, risk-001, note-001, unique within this assessment.\n" +
		"Rating:\n" +
		"- Judge only from your assessment and the frame's appliedRatingCriteria; do not rate on evidence the " +
		"assessment does not record.\n" +
		"- status is one of: rated, not_rated, blocked, not_applicable. When rated, set ratingLevelId to the " +
		"highest rating level whose criterion the assessed evidence satisfies and explain the rationale.\n" +
		"- Record criteriaResults for each rating level considered and ratingDrivers referencing the assessment."
	return nil
}

func (e *engine) fillFactorAnalysisRequest(unit *Unit, req *evaluator.WorkRequest) {
	planned := e.graph.Plan.Factor(unit.Subject)
	frame := e.payloadFor(unitID(KindFrameFactorAnalysis, unit.Subject))
	ratings := map[string]any{}
	for _, reqRef := range planned.Requirements {
		ratings[reqRef] = e.requirementPayload(reqRef, evaluation.DataKindRequirementRating)
	}
	children := map[string]any{}
	for _, childRef := range planned.ChildFactors {
		children[childRef] = e.payloadFor(unitID(KindAnalyzeFactor, childRef))
	}
	req.Context = map[string]any{
		"factorAnalysisFrame":      frame,
		"directRequirementRatings": ratings,
		"childFactorAnalyses":      children,
	}
	req.Instructions = "Synthesize this factor's analysis from its direct requirement ratings and child factor " +
		"analyses, and return one FactorAnalysisResult JSON object.\n" +
		"- Set factorId to the subject reference exactly.\n" +
		"- Fill localAnalysis (direct inputs only) and localAndDescendantAnalysis (including child factors), each " +
		"with status, ratingLevelId when analyzed, rationale, inputRefs, and ratingDrivers.\n" +
		"- Follow the frame's synthesis guidance: the roll-up rating is bounded by the worst contributing input " +
		"(worst_bound); ignore empty inputs.\n" +
		"- Do not inspect new evidence; synthesize only the given inputs."
}

func (e *engine) fillAreaAnalysisRequest(unit *Unit, req *evaluator.WorkRequest) {
	planned := e.graph.Plan.Area(unit.Subject)
	frame := e.payloadFor(unitID(KindFrameAreaAnalysis, unit.Subject))
	factors := map[string]any{}
	for _, factorRef := range planned.RootFactors {
		factors[factorRef] = e.payloadFor(unitID(KindAnalyzeFactor, factorRef))
	}
	children := map[string]any{}
	for _, childRef := range planned.ChildAreas {
		children[childRef] = e.payloadFor(unitID(KindAnalyzeArea, childRef))
	}
	ratings := map[string]any{}
	for _, reqRef := range planned.LocalRequirements {
		ratings[reqRef] = e.requirementPayload(reqRef, evaluation.DataKindRequirementRating)
	}
	req.Context = map[string]any{
		"areaAnalysisFrame":       frame,
		"factorAnalyses":          factors,
		"childAreaAnalyses":       children,
		"localRequirementRatings": ratings,
	}
	req.Instructions = "Synthesize this area's analysis from its factor analyses, child area analyses, and local " +
		"requirement ratings, and return one AreaAnalysisResult JSON object.\n" +
		"- Set areaId to the subject reference exactly.\n" +
		"- Fill localAnalysis (local inputs only) and localAndDescendantAnalysis (including child areas), each with " +
		"status, ratingLevelId when analyzed, rationale, inputRefs, and ratingDrivers.\n" +
		"- Follow the frame's synthesis guidance: the roll-up rating is bounded by the worst contributing input " +
		"(worst_bound); ignore empty inputs.\n" +
		"- Do not inspect new evidence; synthesize only the given inputs."
}

// findingIndexEntry summarizes one persisted finding for advice work.
type findingIndexEntry struct {
	RequirementID string         `json:"requirementId"`
	FindingID     string         `json:"findingId"`
	FindingRef    map[string]any `json:"findingRef"`
	Type          string         `json:"type"`
	Severity      string         `json:"severity,omitempty"`
	Confidence    string         `json:"confidence,omitempty"`
	Statement     string         `json:"statement,omitempty"`
}

func (e *engine) findingIndex() []findingIndexEntry {
	var out []findingIndexEntry
	for _, planned := range e.graph.Plan.Requirements {
		assessment := e.requirementPayload(planned.Ref, evaluation.DataKindRequirementAssessment)
		if assessment == nil {
			continue
		}
		findings, _ := assessment["findings"].([]any)
		for _, item := range findings {
			finding, _ := item.(map[string]any)
			if finding == nil {
				continue
			}
			id, _ := finding["id"].(string)
			if id == "" {
				continue
			}
			entry := findingIndexEntry{
				RequirementID: planned.Ref,
				FindingID:     id,
				FindingRef: routineRefPayload("RequirementAssessmentResult",
					map[string]any{"requirementId": planned.Ref},
					evaluation.FindingSelector(id)),
			}
			entry.Type, _ = finding["type"].(string)
			entry.Severity, _ = finding["severity"].(string)
			entry.Confidence, _ = finding["confidence"].(string)
			entry.Statement, _ = finding["statement"].(string)
			out = append(out, entry)
		}
	}
	return out
}

func (e *engine) fillRankFindingsRequest(req *evaluator.WorkRequest) {
	findings := e.findingIndex()
	req.Context = map[string]any{"findings": findings}
	req.Instructions = "Rank every persisted finding across the evaluation scope and return one " +
		"FindingRankingResult JSON object.\n" +
		"- orderedFindings must contain exactly one entry per finding in the findings context, no more, no fewer.\n" +
		"- Copy each entry's findingRef object verbatim from the findings context.\n" +
		"- rank is 1-based and unique; tier is one of P1, P2, P3, P4 (P1 = act first).\n" +
		"- Rank by severity, confidence, and breadth of effect; give each entry a one-sentence rationale."
}

func (e *engine) fillRecommendRequest(req *evaluator.WorkRequest) {
	analyses := map[string]any{}
	for _, area := range e.graph.Plan.Areas {
		analyses[area.Ref] = e.payloadFor(unitID(KindAnalyzeArea, area.Ref))
	}
	req.Context = map[string]any{
		"findingRanking": e.payloadFor(string(KindRankFindings)),
		"findings":       e.findingIndex(),
		"areaAnalyses":   analyses,
	}
	req.Instructions = "Propose quality-management recommendations from the ranked findings and analyses, and " +
		"return one JSON object of the form {\"recommendations\": [RecommendationResult, ...]}.\n" +
		"- Each recommendation needs title, description, background, expectedValue, doneCriterion, impact " +
		"(high, medium, or low), confidence, and non-empty traceRefs pointing at the findings or analyses it follows from.\n" +
		"- Omit the id field; the runner assigns canonical recommendation IDs.\n" +
		"- Do not include planning fields (effort, roi, quickWin, priority, score).\n" +
		"- Cover the highest-ranked findings first; propose the smallest set of recommendations that addresses the " +
		"advice-driving findings (typically 2-6)."
}

func (e *engine) fillRankRecommendationsRequest(req *evaluator.WorkRequest) {
	recommendations := e.payloadsByWorkUnit(string(KindRecommend))
	req.Context = map[string]any{
		"recommendations": recommendations,
		"findings":        e.findingIndex(),
		"findingRanking":  e.payloadFor(string(KindRankFindings)),
	}
	req.Instructions = "Rank the recommendations and account for finding coverage, and return one " +
		"RecommendationRankingResult JSON object.\n" +
		"- orderedRecommendations must contain exactly one entry per recommendation in the recommendations context, " +
		"using each recommendation's id as recommendationRef, with 1-based unique ranks, impact, confidence, and rationale.\n" +
		"- findingCoverage must contain exactly one entry per finding in the findings context: copy findingRef " +
		"verbatim, set disposition to addressed_by_recommendation (with recommendationRefs listing covering " +
		"recommendation ids) or not_advice_driving (with a short rationale)."
}
