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

// AddRecord writes one evaluation record and preserves the legacy single-path
// receipt field when exactly one record is created.
func AddRecord(kind RecordKind, runPath string, raw []byte) (*WriteRecordReceipt, error) {
	result, err := WriteRecords(kind, runPath, raw)
	if err != nil {
		return nil, err
	}
	if len(result.Paths) == 1 && result.Path == "" {
		result.Path = result.Paths[0]
	}
	return result, nil
}

// WriteRecords writes one or more evaluation records from a JSON payload.
func WriteRecords(kind RecordKind, runPath string, raw []byte) (*WriteRecordReceipt, error) {
	runAbs, err := verifyRun(runPath)
	if err != nil {
		return nil, err
	}
	runDisplay := displayRunPath(runAbs)
	switch kind {
	case KindAssessmentResult:
		levels, err := ratingLevels(filepath.Join(runAbs, "model.md"))
		if err != nil {
			return nil, err
		}
		return addAssessmentResults(runAbs, runDisplay, raw, levels)
	case KindAnalysis:
		spec, err := decodeRunModel(filepath.Join(runAbs, "model.md"))
		if err != nil {
			return nil, err
		}
		return setAnalyses(runAbs, runDisplay, raw, ratingLevelSetFromSpec(spec), newFactorVocabulary(spec))
	case KindRecommendation:
		return addRecommendations(runAbs, runDisplay, raw)
	default:
		return nil, usagef("unknown record kind %q", kind)
	}
}

// DecodeSingleJSON decodes one strict JSON document into dst.
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

// DecodeJSONList decodes the input into a slice of T. It accepts either a
// single JSON object (returned as a one-element slice) or a JSON array. Decoding
// is strict: unknown fields and trailing documents are rejected. Empty or
// otherwise invalid input returns a usage error.
func DecodeJSONList[T any](raw []byte) ([]T, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, usagef("input is empty")
	}
	var values []T
	switch trimmed[0] {
	case '[':
		if err := decodeJSONDocument(trimmed, &values); err != nil {
			return nil, err
		}
	case '{':
		var value T
		if err := decodeJSONDocument(trimmed, &value); err != nil {
			return nil, err
		}
		values = []T{value}
	default:
		return nil, usagef("input must be a JSON object or array")
	}
	return values, nil
}

func decodeJSONDocument(raw []byte, dst any) error {
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

func addAssessmentResults(runAbs, runDisplay string, raw []byte, levels map[string]bool) (*WriteRecordReceipt, error) {
	payloads, err := DecodeJSONList[AssessmentResultInput](raw)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, payload := range payloads {
		if err := validateAssessmentResult(payload, levels); err != nil {
			return nil, err
		}
		rec := AssessmentResultRecord{
			SchemaVersion:   SchemaVersion,
			AreaPath:        payload.AreaPath,
			Requirement:     payload.Requirement,
			FactorPaths:     payload.FactorPaths,
			RatingResult:    payload.RatingResult,
			CriterionSource: payload.CriterionSource,
			Findings:        payload.Findings,
			Recommendations: payload.Recommendations,
			Supersedes:      payload.Supersedes,
		}
		data, err := marshalJSON(rec)
		if err != nil {
			return nil, err
		}
		path, err := writeNumbered(filepath.Join(runAbs, "assessments"), areaPathSlug(payload.AreaPath)+"-"+Slug(payload.Requirement)+".json", data)
		if err != nil {
			return nil, err
		}
		paths = append(paths, displayRecordPath(runAbs, runDisplay, path))
	}
	return &WriteRecordReceipt{
		SchemaVersion: SchemaVersion,
		Path:          singlePath(paths),
		Paths:         paths,
		Kind:          KindAssessmentResult,
		NextActions: []receipt.Action{{
			ID:      "evaluation-status",
			Label:   "Inspect report readiness",
			Command: "qualitymd evaluation status " + runDisplay,
		}},
	}, nil
}

func setAnalyses(runAbs, runDisplay string, raw []byte, levels map[string]bool, vocab factorVocabulary) (*WriteRecordReceipt, error) {
	payloads, err := DecodeJSONList[AnalysisInput](raw)
	if err != nil {
		return nil, err
	}
	var paths []string
	var createdPtr *bool
	for _, payload := range payloads {
		if err := validateAnalysis(payload, levels, vocab); err != nil {
			return nil, err
		}
		rec := AnalysisRecord{
			SchemaVersion:           SchemaVersion,
			AreaPath:                payload.AreaPath,
			LocalRatingResult:       payload.LocalRatingResult,
			FactorRatingResults:     payload.FactorRatingResults,
			AggregateRatingResult:   payload.AggregateRatingResult,
			AssessmentResultRecords: payload.AssessmentResultRecords,
			ChildAnalysisRecords:    payload.ChildAnalysisRecords,
			RatingConstraints:       payload.RatingConstraints,
		}
		data, err := marshalJSON(rec)
		if err != nil {
			return nil, err
		}
		path := filepath.Join(runAbs, "analysis", areaPathSlug(payload.AreaPath)+".json")
		_, statErr := os.Stat(path)
		created := os.IsNotExist(statErr)
		if len(payloads) == 1 {
			createdPtr = &created
		}
		if err := writeReplace(path, data); err != nil {
			return nil, err
		}
		paths = append(paths, displayRecordPath(runAbs, runDisplay, path))
	}
	return &WriteRecordReceipt{SchemaVersion: SchemaVersion, Path: singlePath(paths), Paths: paths, Kind: KindAnalysis, Created: createdPtr}, nil
}

func addRecommendations(runAbs, runDisplay string, raw []byte) (*WriteRecordReceipt, error) {
	payloads, err := DecodeJSONList[RecommendationInput](raw)
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, payload := range payloads {
		if err := validateRecommendation(payload); err != nil {
			return nil, err
		}
		rec := RecommendationRecord{
			SchemaVersion:           SchemaVersion,
			Title:                   payload.Title,
			Gap:                     payload.Gap,
			EvidenceLocators:        payload.EvidenceLocators,
			AssessmentResultRecords: payload.AssessmentResultRecords,
			RemediationOptions:      payload.RemediationOptions,
			RecommendedOption:       payload.RecommendedOption,
			DoneCriterion:           payload.DoneCriterion,
			Supersedes:              payload.Supersedes,
		}
		data, err := renderRecommendation(rec)
		if err != nil {
			return nil, err
		}
		path, err := writeNumbered(filepath.Join(runAbs, "recommendations"), Slug(payload.Title)+".md", data)
		if err != nil {
			return nil, err
		}
		paths = append(paths, displayRecordPath(runAbs, runDisplay, path))
	}
	return &WriteRecordReceipt{SchemaVersion: SchemaVersion, Path: singlePath(paths), Paths: paths, Kind: KindRecommendation}, nil
}

func singlePath(paths []string) string {
	if len(paths) == 1 {
		return paths[0]
	}
	return ""
}

func areaPathSlug(areaPath []string) string {
	if len(areaPath) == 0 {
		return "root"
	}
	return Slug(strings.Join(areaPath, "-"))
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

func displayRunPath(runAbs string) string {
	runAbs = filepath.Clean(runAbs)
	repoRoot, err := FindRepoRoot(runAbs)
	if err == nil {
		if rel, relErr := filepath.Rel(repoRoot, runAbs); relErr == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".." {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(runAbs)
}

func displayRecordPath(runAbs, runDisplay, recordAbs string) string {
	rel, err := filepath.Rel(runAbs, recordAbs)
	if err != nil {
		return filepath.ToSlash(recordAbs)
	}
	return filepath.ToSlash(filepath.Join(runDisplay, rel))
}

func decodeRunModel(path string) (*model.Spec, error) {
	doc, err := document.Parse(path)
	if err != nil {
		return nil, err
	}
	return model.Decode(doc)
}

func ratingLevelSetFromSpec(spec *model.Spec) map[string]bool {
	levels := map[string]bool{}
	if spec == nil {
		return levels
	}
	for _, level := range spec.RatingScale {
		levels[level.Level] = true
	}
	return levels
}

func ratingLevels(path string) (map[string]bool, error) {
	spec, err := decodeRunModel(path)
	if err != nil {
		return nil, err
	}
	return ratingLevelSetFromSpec(spec), nil
}

func validateAssessmentResult(p AssessmentResultInput, levels map[string]bool) error {
	if err := validateAssessmentResultRequiredStrings(p); err != nil {
		return err
	}
	if err := validateRatingResult("ratingResult", &p.RatingResult, levels); err != nil {
		return err
	}
	if p.AreaPath == nil {
		return usagef("areaPath is required")
	}
	if p.FactorPaths == nil {
		return usagef("factorPaths is required")
	}
	if p.Findings == nil {
		return usagef("findings is required")
	}
	if p.Recommendations == nil {
		return usagef("recommendations is required")
	}
	if err := validateAssessmentResultFindings(p.Findings); err != nil {
		return err
	}
	if err := validateRequiredStrings("supersedes", p.Supersedes); err != nil {
		return err
	}
	return nil
}

func validateAssessmentResultRequiredStrings(p AssessmentResultInput) error {
	for name, value := range map[string]string{
		"requirement":     p.Requirement,
		"criterionSource": p.CriterionSource,
	} {
		if err := requiredString(name, value); err != nil {
			return err
		}
	}
	return nil
}

func validateAssessmentResultFindings(findings []Finding) error {
	for i, finding := range findings {
		if strings.TrimSpace(finding.Locator) == "" || strings.TrimSpace(finding.Observation) == "" || strings.TrimSpace(finding.Category) == "" || strings.TrimSpace(string(finding.Severity)) == "" {
			return usagef("findings[%d] must include locator, observation, category, and severity", i)
		}
		if !finding.Severity.Valid() {
			return usagef("findings[%d].severity must be one of %s", i, findingSeverityLevels())
		}
	}
	return nil
}

func validateRequiredStrings(name string, values []string) error {
	for i, value := range values {
		if strings.TrimSpace(value) == "" {
			return usagef("%s[%d] is required", name, i)
		}
	}
	return nil
}

func validateAnalysis(p AnalysisInput, levels map[string]bool, vocab factorVocabulary) error {
	if p.AreaPath == nil {
		return usagef("areaPath is required")
	}
	if err := validateRatingResult("aggregateRatingResult", &p.AggregateRatingResult, levels); err != nil {
		return err
	}
	if p.LocalRatingResult != nil {
		if err := validateRatingResult("localRatingResult", p.LocalRatingResult, levels); err != nil {
			return err
		}
	}
	seenFactorPaths := map[string]bool{}
	for i := range p.FactorRatingResults {
		factorPath := p.FactorRatingResults[i].FactorPath
		if len(factorPath) == 0 {
			return usagef("factorRatingResults[%d].factorPath is required", i)
		}
		key := factorPath.IdentityKey()
		if seenFactorPaths[key] {
			return usagef("factorRatingResults[%d].factorPath %q is a duplicate within the analysis record", i, factorPath.Display())
		}
		seenFactorPaths[key] = true
		if !vocab.Resolves(p.AreaPath, factorPath) {
			return usagef("factorRatingResults[%d].factorPath %q does not resolve against the area's declared or inherited factor vocabulary", i, factorPath.Display())
		}
		if err := validateRatingResult(fmt.Sprintf("factorRatingResults[%d].ratingResult", i), &p.FactorRatingResults[i].RatingResult, levels); err != nil {
			return err
		}
	}
	if p.FactorRatingResults == nil {
		return usagef("factorRatingResults is required")
	}
	if p.AssessmentResultRecords == nil {
		return usagef("assessmentResultRecords is required")
	}
	if p.ChildAnalysisRecords == nil {
		return usagef("childAnalysisRecords is required")
	}
	return nil
}

func validateRatingResult(name string, result *RatingResult, levels map[string]bool) error {
	if strings.TrimSpace(result.Rationale) == "" {
		return usagef("%s.rationale is required", name)
	}
	switch result.Kind {
	case RatingResultRated:
		if strings.TrimSpace(result.Level) == "" {
			return usagef("%s.level is required when kind is rated", name)
		}
		if !levels[result.Level] {
			return usagef("%s.level %q is not defined by the run model", name, result.Level)
		}
	case RatingResultNotAssessed:
		if strings.TrimSpace(result.Level) != "" {
			return usagef("%s.level must be empty when kind is not-assessed", name)
		}
	default:
		return usagef("%s.kind must be rated or not-assessed", name)
	}
	return nil
}

func validateRecommendation(p RecommendationInput) error {
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
	addYAMLSeq(&node, "assessmentResultRecords", rec.AssessmentResultRecords)
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

// factorVocabulary resolves factor paths against a run model's declared and
// inherited factor vocabulary. It is built once from the parsed model and used
// by the analysis write path to reject factor paths that cannot be resolved.
type factorVocabulary struct {
	spec *model.Spec
}

func newFactorVocabulary(spec *model.Spec) factorVocabulary {
	return factorVocabulary{spec: spec}
}

// Resolves reports whether factorPath resolves against the vocabulary visible to
// the area at areaPath: the union of the root model factors, every ancestor
// area's factors, and the area's own factors. Path elements are matched against
// factor-map keys or factor titles, case-insensitively. When the model is
// absent, resolution is permissive so model-free runs are not blocked.
func (v factorVocabulary) Resolves(areaPath AreaPath, factorPath FactorPath) bool {
	if v.spec == nil {
		return true
	}
	if len(factorPath) == 0 {
		return false
	}
	for _, scope := range v.factorScopes(areaPath) {
		if resolveFactorPath(scope, factorPath) {
			return true
		}
	}
	return false
}

// factorScopes returns the factor maps visible to the area at areaPath, from the
// root model down through each ancestor area to the area itself.
func (v factorVocabulary) factorScopes(areaPath AreaPath) []map[string]model.Factor {
	scopes := []map[string]model.Factor{v.spec.Factors}
	areas := v.spec.Areas
	for _, element := range areaPath {
		area, ok := areas[element]
		if !ok {
			break
		}
		scopes = append(scopes, area.Factors)
		areas = area.Areas
	}
	return scopes
}

// resolveFactorPath reports whether factorPath resolves within a single factor
// map by walking each element to a matching factor or sub-factor.
func resolveFactorPath(factors map[string]model.Factor, factorPath FactorPath) bool {
	current := factors
	for i, element := range factorPath {
		factor, ok := matchFactor(current, element)
		if !ok {
			return false
		}
		if i == len(factorPath)-1 {
			return true
		}
		current = factor.Factors
	}
	return false
}

// matchFactor finds the factor in factors whose map key or title matches element
// case-insensitively.
func matchFactor(factors map[string]model.Factor, element string) (model.Factor, bool) {
	target := strings.ToLower(strings.TrimSpace(element))
	if factor, ok := factors[element]; ok {
		return factor, true
	}
	for key, factor := range factors {
		if strings.ToLower(strings.TrimSpace(key)) == target {
			return factor, true
		}
		if strings.ToLower(strings.TrimSpace(factor.Title)) == target {
			return factor, true
		}
	}
	return model.Factor{}, false
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
