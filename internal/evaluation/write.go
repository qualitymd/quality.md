package evaluation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

func AddRecord(kind WriteKind, runPath string, raw []byte) (*WriteResult, error) {
	runAbs, err := verifyRun(runPath)
	if err != nil {
		return nil, err
	}
	levels, err := ratingLevels(filepath.Join(runAbs, "model.md"))
	if err != nil {
		return nil, err
	}
	switch kind {
	case KindAssessment:
		return addAssessment(runAbs, raw, levels)
	case KindAnalysis:
		return addAnalysis(runAbs, raw, levels)
	case KindRecommendation:
		return addRecommendation(runAbs, raw)
	default:
		return nil, usagef("unknown record kind %q", kind)
	}
}

func DecodeSingleJSON(raw []byte, dst any) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		if errors.Is(err, io.EOF) {
			return usagef("input is empty")
		}
		return usagef("invalid JSON payload: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return usagef("input must contain exactly one JSON document")
	}
	return nil
}

func addAssessment(runAbs string, raw []byte, levels map[string]bool) (*WriteResult, error) {
	var payload AssessmentPayload
	if err := DecodeSingleJSON(raw, &payload); err != nil {
		return nil, err
	}
	if err := validateAssessment(payload, levels); err != nil {
		return nil, err
	}
	rec := AssessmentRecord{
		SchemaVersion:   SchemaVersion,
		Target:          payload.Target,
		TargetPath:      payload.TargetPath,
		Requirement:     payload.Requirement,
		Factors:         payload.Factors,
		Rating:          payload.Rating,
		NotAssessed:     payload.NotAssessed,
		CriterionSource: payload.CriterionSource,
		Findings:        payload.Findings,
		Rationale:       payload.Rationale,
		Recommendations: payload.Recommendations,
		Supersedes:      payload.Supersedes,
	}
	data, err := marshalJSON(rec)
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(runAbs, "assessments")
	path, err := writeNumbered(dir, Slug(payload.Target)+"-"+Slug(payload.Requirement)+".json", data)
	if err != nil {
		return nil, err
	}
	rel := filepath.ToSlash(path)
	return &WriteResult{
		SchemaVersion: SchemaVersion,
		Path:          rel,
		Kind:          KindAssessment,
		NextActions: []receipt.Action{{
			ID:      "show-status",
			Label:   "Inspect report readiness",
			Command: "qualitymd evaluation show-status " + filepath.ToSlash(runAbs),
		}},
	}, nil
}

func addAnalysis(runAbs string, raw []byte, levels map[string]bool) (*WriteResult, error) {
	var payload AnalysisPayload
	if err := DecodeSingleJSON(raw, &payload); err != nil {
		return nil, err
	}
	if err := validateAnalysis(payload, levels); err != nil {
		return nil, err
	}
	rec := AnalysisRecord{
		SchemaVersion:        SchemaVersion,
		Target:               payload.Target,
		TargetPath:           payload.TargetPath,
		LocalRating:          payload.LocalRating,
		FactorRatings:        payload.FactorRatings,
		AggregateRating:      payload.AggregateRating,
		AssessmentRecords:    payload.AssessmentRecords,
		ChildAnalysisRecords: payload.ChildAnalysisRecords,
		BindingConstraints:   payload.BindingConstraints,
	}
	data, err := marshalJSON(rec)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(runAbs, "analysis", Slug(payload.Target)+".json")
	_, statErr := os.Stat(path)
	created := os.IsNotExist(statErr)
	if err := writeReplace(path, data); err != nil {
		return nil, err
	}
	rel := filepath.ToSlash(path)
	return &WriteResult{SchemaVersion: SchemaVersion, Path: rel, Kind: KindAnalysis, Created: &created}, nil
}

func addRecommendation(runAbs string, raw []byte) (*WriteResult, error) {
	var payload RecommendationPayload
	if err := DecodeSingleJSON(raw, &payload); err != nil {
		return nil, err
	}
	if err := validateRecommendation(payload); err != nil {
		return nil, err
	}
	rec := RecommendationRecord{
		SchemaVersion:      SchemaVersion,
		Title:              payload.Title,
		Gap:                payload.Gap,
		EvidenceLocators:   payload.EvidenceLocators,
		AssessmentRecords:  payload.AssessmentRecords,
		RemediationOptions: payload.RemediationOptions,
		RecommendedOption:  payload.RecommendedOption,
		DoneCriterion:      payload.DoneCriterion,
		Supersedes:         payload.Supersedes,
	}
	data, err := renderRecommendation(rec)
	if err != nil {
		return nil, err
	}
	path, err := writeNumbered(filepath.Join(runAbs, "recommendations"), Slug(payload.Title)+".md", data)
	if err != nil {
		return nil, err
	}
	return &WriteResult{SchemaVersion: SchemaVersion, Path: filepath.ToSlash(path), Kind: KindRecommendation}, nil
}

func verifyRun(runPath string) (string, error) {
	abs, err := filepath.Abs(runPath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("reading run %s: %w", runPath, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not an evaluation run folder", runPath)
	}
	for _, name := range []string{"model.md", "design.md", "plan.md", "assessments", "analysis", "recommendations"} {
		if _, err := os.Stat(filepath.Join(abs, name)); err != nil {
			return "", fmt.Errorf("%s is not an evaluation run folder: missing %s", runPath, name)
		}
	}
	return abs, nil
}

func ratingLevels(path string) (map[string]bool, error) {
	doc, err := document.Parse(path)
	if err != nil {
		return nil, err
	}
	spec, err := model.Decode(doc)
	if err != nil {
		return nil, err
	}
	levels := map[string]bool{}
	for _, level := range spec.RatingScale {
		levels[level.Level] = true
	}
	return levels, nil
}

func validateAssessment(p AssessmentPayload, levels map[string]bool) error {
	for name, value := range map[string]string{
		"target":          p.Target,
		"requirement":     p.Requirement,
		"criterionSource": p.CriterionSource,
		"rationale":       p.Rationale,
	} {
		if err := requiredString(name, value); err != nil {
			return err
		}
	}
	if p.NotAssessed && p.Rating != nil {
		return usagef("rating must be null when notAssessed is true")
	}
	if !p.NotAssessed && p.Rating == nil {
		return usagef("rating is required unless notAssessed is true")
	}
	if p.Rating != nil && !levels[*p.Rating] {
		return usagef("rating %q is not defined by the run model", *p.Rating)
	}
	for i, finding := range p.Findings {
		if strings.TrimSpace(finding.Locator) == "" || strings.TrimSpace(finding.Observation) == "" || strings.TrimSpace(finding.Category) == "" {
			return usagef("findings[%d] must include locator, observation, and category", i)
		}
	}
	if p.Factors == nil {
		return usagef("factors is required")
	}
	if p.Findings == nil {
		return usagef("findings is required")
	}
	if p.Recommendations == nil {
		return usagef("recommendations is required")
	}
	for i, ref := range p.Supersedes {
		if strings.TrimSpace(ref) == "" {
			return usagef("supersedes[%d] is required", i)
		}
	}
	return nil
}

func validateAnalysis(p AnalysisPayload, levels map[string]bool) error {
	if err := requiredString("target", p.Target); err != nil {
		return err
	}
	if err := validateRatingResult("aggregateRating", &p.AggregateRating, levels, true); err != nil {
		return err
	}
	if p.LocalRating != nil {
		if err := validateRatingResult("localRating", p.LocalRating, levels, false); err != nil {
			return err
		}
	}
	for i := range p.FactorRatings {
		if strings.TrimSpace(p.FactorRatings[i].Factor) == "" {
			return usagef("factorRatings[%d].factor is required", i)
		}
		rating := RatingResult{
			Rating:      p.FactorRatings[i].Rating,
			NotAssessed: p.FactorRatings[i].NotAssessed,
			Rationale:   p.FactorRatings[i].Rationale,
		}
		if err := validateRatingResult(fmt.Sprintf("factorRatings[%d]", i), &rating, levels, true); err != nil {
			return err
		}
	}
	if p.FactorRatings == nil {
		return usagef("factorRatings is required")
	}
	if p.AssessmentRecords == nil {
		return usagef("assessmentRecords is required")
	}
	if p.ChildAnalysisRecords == nil {
		return usagef("childAnalysisRecords is required")
	}
	return nil
}

func validateRatingResult(name string, result *RatingResult, levels map[string]bool, requireRating bool) error {
	if strings.TrimSpace(result.Rationale) == "" {
		return usagef("%s.rationale is required", name)
	}
	if result.NotAssessed && result.Rating != nil {
		return usagef("%s.rating must be null when notAssessed is true", name)
	}
	if !result.NotAssessed && requireRating && result.Rating == nil {
		return usagef("%s.rating is required unless notAssessed is true", name)
	}
	if result.Rating != nil && !levels[*result.Rating] {
		return usagef("%s.rating %q is not defined by the run model", name, *result.Rating)
	}
	return nil
}

func validateRecommendation(p RecommendationPayload) error {
	for name, value := range map[string]string{
		"title":             p.Title,
		"gap":               p.Gap,
		"recommendedOption": p.RecommendedOption,
		"doneCriterion":     p.DoneCriterion,
	} {
		if err := requiredString(name, value); err != nil {
			return err
		}
	}
	if len(p.EvidenceLocators) == 0 {
		return usagef("evidenceLocators is required")
	}
	if len(p.RemediationOptions) == 0 {
		return usagef("remediationOptions is required")
	}
	for i, ref := range p.Supersedes {
		if strings.TrimSpace(ref) == "" {
			return usagef("supersedes[%d] is required", i)
		}
	}
	return nil
}

func renderRecommendation(rec RecommendationRecord) ([]byte, error) {
	var out bytes.Buffer
	out.WriteString("---\n")
	node := yaml.Node{Kind: yaml.MappingNode}
	addYAMLScalar(&node, "schemaVersion", fmt.Sprint(rec.SchemaVersion), "!!int")
	addYAMLScalar(&node, "title", rec.Title, "!!str")
	addYAMLScalar(&node, "gap", rec.Gap, "!!str")
	addYAMLSeq(&node, "evidenceLocators", rec.EvidenceLocators)
	addYAMLSeq(&node, "assessmentRecords", rec.AssessmentRecords)
	addYAMLSeq(&node, "remediationOptions", rec.RemediationOptions)
	addYAMLScalar(&node, "recommendedOption", rec.RecommendedOption, "!!str")
	addYAMLScalar(&node, "doneCriterion", rec.DoneCriterion, "!!str")
	if len(rec.Supersedes) > 0 {
		addYAMLSeq(&node, "supersedes", rec.Supersedes)
	}
	enc := yaml.NewEncoder(&out)
	enc.SetIndent(2)
	if err := enc.Encode(&node); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	out.WriteString("---\n\n")
	out.WriteString("# " + rec.Title + "\n\n")
	out.WriteString("## Gap\n\n" + rec.Gap + "\n\n")
	out.WriteString("## Evidence locators\n\n")
	for _, locator := range rec.EvidenceLocators {
		out.WriteString("- `" + locator + "`\n")
	}
	out.WriteString("\n## Remediation options\n\n")
	for _, option := range rec.RemediationOptions {
		out.WriteString("- " + option + "\n")
	}
	out.WriteString("\n## Recommended option\n\n" + rec.RecommendedOption + "\n\n")
	out.WriteString("## Done criterion\n\n" + rec.DoneCriterion + "\n")
	if len(rec.Supersedes) > 0 {
		out.WriteString("\n## Supersedes\n\n")
		for _, ref := range rec.Supersedes {
			out.WriteString("- `" + ref + "`\n")
		}
	}
	return out.Bytes(), nil
}

func addYAMLScalar(node *yaml.Node, key, value, tag string) {
	node.Content = append(node.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key},
		&yaml.Node{Kind: yaml.ScalarNode, Tag: tag, Value: value},
	)
}

func addYAMLSeq(node *yaml.Node, key string, values []string) {
	seq := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	for _, value := range values {
		seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value})
	}
	node.Content = append(node.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: key}, seq)
}
