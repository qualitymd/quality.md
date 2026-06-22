package evaluation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
func WriteRecords(kind RecordKind, runPath string, raw []byte, options ...WriteOptions) (*WriteRecordReceipt, error) {
	opts := WriteOptions{}
	if len(options) > 0 {
		opts = options[0]
	}
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
		return addAssessmentResults(runAbs, runDisplay, raw, levels, opts)
	case KindAnalysis:
		spec, err := decodeRunModel(filepath.Join(runAbs, "model.md"))
		if err != nil {
			return nil, err
		}
		return setAnalyses(runAbs, runDisplay, raw, ratingLevelSetFromSpec(spec), newFactorVocabulary(spec), opts)
	case KindRecommendation:
		return addRecommendations(runAbs, runDisplay, raw, opts)
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
	values, _, err := decodeJSONList[T](raw)
	return values, err
}

func decodeJSONList[T any](raw []byte) ([]T, bool, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, false, usagef("input is empty")
	}
	var values []T
	switch trimmed[0] {
	case '[':
		rawItems, err := decodeRawJSONArray(trimmed)
		if err != nil {
			return nil, true, err
		}
		values = make([]T, 0, len(rawItems))
		var acc validationAccumulator
		for i, rawItem := range rawItems {
			value, err := decodePayloadObject[T](rawItem)
			if err != nil {
				acc.Merge(fmt.Sprintf("[%d]", i), err)
				values = append(values, value)
				continue
			}
			values = append(values, value)
		}
		if err := acc.Err(); err != nil {
			return values, true, err
		}
	case '{':
		value, err := decodePayloadObject[T](trimmed)
		if err != nil {
			return []T{value}, false, err
		}
		values = []T{value}
	default:
		return nil, false, usagef("input must be a JSON object or array")
	}
	return values, trimmed[0] == '[', nil
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

func decodeRawJSONArray(raw []byte) ([]json.RawMessage, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	var values []json.RawMessage
	if err := dec.Decode(&values); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, usagef("input is empty")
		}
		return nil, usagef("invalid JSON payload: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return nil, usagef("input must contain exactly one JSON document")
	}
	return values, nil
}

func decodePayloadObject[T any](raw []byte) (T, error) {
	var zero T
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return zero, usagef("input is empty")
	}
	if trimmed[0] != '{' {
		return zero, usagef("input must be a JSON object or array")
	}
	if !json.Valid(trimmed) {
		return zero, usagef("invalid JSON payload")
	}
	var acc validationAccumulator
	collectJSONShapeProblems(trimmed, reflect.TypeOf(zero), "", &acc)
	if err := acc.Err(); err != nil {
		return zero, err
	}
	var value T
	if err := decodeJSONDocument(trimmed, &value); err != nil {
		return zero, err
	}
	return value, nil
}

func collectJSONShapeProblems(raw json.RawMessage, typ reflect.Type, field string, acc *validationAccumulator) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return
	}
	unwrapped, ok := unwrapJSONType(raw, typ)
	if !ok {
		return
	}
	if !jsonShapeMatches(raw, unwrapped) {
		acc.AddExpected(field, "has the wrong JSON type", jsonTypeDescription(unwrapped))
		return
	}
	switch unwrapped.Kind() {
	case reflect.Struct:
		collectJSONObjectProblems(raw, unwrapped, field, acc)
	case reflect.Slice, reflect.Array:
		collectJSONArrayProblems(raw, unwrapped, field, acc)
	}
}

func unwrapJSONType(raw []byte, typ reflect.Type) (reflect.Type, bool) {
	for typ.Kind() == reflect.Pointer {
		if bytes.Equal(raw, []byte("null")) {
			return typ, false
		}
		typ = typ.Elem()
	}
	return typ, true
}

func collectJSONObjectProblems(raw json.RawMessage, typ reflect.Type, field string, acc *validationAccumulator) {
	var object map[string]json.RawMessage
	if err := json.Unmarshal(raw, &object); err != nil {
		acc.Add(field, "must be a JSON object")
		return
	}
	fields := jsonStructFields(typ)
	allowed := sortedMapKeys(fields)
	for name, value := range object {
		fieldInfo, ok := fields[name]
		if !ok {
			acc.AddAllowed(prefixField(field, name), "is not a supported field", allowed)
			continue
		}
		collectJSONShapeProblems(value, fieldInfo.Type, prefixField(field, name), acc)
	}
}

func collectJSONArrayProblems(raw json.RawMessage, typ reflect.Type, field string, acc *validationAccumulator) {
	if typ.Elem().Kind() == reflect.Uint8 {
		return
	}
	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		acc.Add(field, "must be a JSON array")
		return
	}
	for i, item := range items {
		collectJSONShapeProblems(item, typ.Elem(), fmt.Sprintf("%s[%d]", field, i), acc)
	}
}

func sortedMapKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for name := range values {
		keys = append(keys, name)
	}
	sort.Strings(keys)
	return keys
}

func jsonStructFields(typ reflect.Type) map[string]reflect.StructField {
	fields := map[string]reflect.StructField{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name := jsonFieldName(field)
		if name == "" {
			continue
		}
		fields[name] = field
	}
	return fields
}

func jsonShapeMatches(raw []byte, typ reflect.Type) bool {
	if bytes.Equal(raw, []byte("null")) {
		return nullableJSONKind(typ.Kind())
	}
	switch typ.Kind() {
	case reflect.String:
		return raw[0] == '"'
	case reflect.Bool:
		return bytes.Equal(raw, []byte("true")) || bytes.Equal(raw, []byte("false"))
	case reflect.Slice, reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return raw[0] == '"'
		}
		return raw[0] == '['
	case reflect.Map, reflect.Struct:
		return raw[0] == '{'
	case reflect.Interface:
		return true
	default:
		return !numericJSONKind(typ.Kind()) || (raw[0] >= '0' && raw[0] <= '9') || raw[0] == '-'
	}
}

func nullableJSONKind(kind reflect.Kind) bool {
	return kind == reflect.Pointer || kind == reflect.Interface || kind == reflect.Map || kind == reflect.Slice
}

func numericJSONKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

type plannedRecordWrite struct {
	Path    string
	Data    []byte
	Created *bool
}

func addAssessmentResults(runAbs, runDisplay string, raw []byte, levels map[string]bool, opts WriteOptions) (*WriteRecordReceipt, error) {
	payloads, fromArray, err := decodeJSONList[AssessmentResultInput](raw)
	if err != nil && len(payloads) == 0 {
		return nil, err
	}
	var acc validationAccumulator
	acc.Merge("", err)
	for i, payload := range payloads {
		prefix := validationPrefix(fromArray, i)
		acc.Merge(prefix, validateAssessmentResult(payload, levels))
	}
	if err := acc.Err(); err != nil {
		return nil, err
	}
	var writes []plannedRecordWrite
	nextNumber, err := nextRecordNumber(filepath.Join(runAbs, "assessments"))
	if err != nil {
		return nil, err
	}
	for _, payload := range payloads {
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
		path := filepath.Join(runAbs, "assessments", fmt.Sprintf("%03d-%s", nextNumber, areaPathSlug(payload.AreaPath)+"-"+Slug(payload.Requirement)+".json"))
		nextNumber++
		writes = append(writes, plannedRecordWrite{Path: path, Data: data})
	}
	return commitPlannedWrites(runAbs, runDisplay, KindAssessmentResult, writes, opts, []receipt.Action{{
		ID:      "evaluation-status",
		Label:   "Inspect report readiness",
		Command: "qualitymd evaluation status " + runDisplay,
	}})
}

func setAnalyses(runAbs, runDisplay string, raw []byte, levels map[string]bool, vocab factorVocabulary, opts WriteOptions) (*WriteRecordReceipt, error) {
	payloads, fromArray, err := decodeJSONList[AnalysisInput](raw)
	if err != nil && len(payloads) == 0 {
		return nil, err
	}
	var acc validationAccumulator
	acc.Merge("", err)
	for i, payload := range payloads {
		prefix := validationPrefix(fromArray, i)
		acc.Merge(prefix, validateAnalysis(payload, levels, vocab))
	}
	if err := acc.Err(); err != nil {
		return nil, err
	}
	var writes []plannedRecordWrite
	for _, payload := range payloads {
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
		write := plannedRecordWrite{Path: path, Data: data}
		if len(payloads) == 1 {
			write.Created = &created
		}
		writes = append(writes, write)
	}
	return commitPlannedWrites(runAbs, runDisplay, KindAnalysis, writes, opts, nil)
}

func addRecommendations(runAbs, runDisplay string, raw []byte, opts WriteOptions) (*WriteRecordReceipt, error) {
	payloads, fromArray, err := decodeJSONList[RecommendationInput](raw)
	if err != nil && len(payloads) == 0 {
		return nil, err
	}
	var acc validationAccumulator
	acc.Merge("", err)
	for i, payload := range payloads {
		prefix := validationPrefix(fromArray, i)
		acc.Merge(prefix, validateRecommendation(payload))
	}
	if err := acc.Err(); err != nil {
		return nil, err
	}
	var writes []plannedRecordWrite
	nextNumber, err := nextRecordNumber(filepath.Join(runAbs, "recommendations"))
	if err != nil {
		return nil, err
	}
	for _, payload := range payloads {
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
		path := filepath.Join(runAbs, "recommendations", fmt.Sprintf("%03d-%s", nextNumber, Slug(payload.Title)+".md"))
		nextNumber++
		writes = append(writes, plannedRecordWrite{Path: path, Data: data})
	}
	return commitPlannedWrites(runAbs, runDisplay, KindRecommendation, writes, opts, nil)
}

func validationPrefix(fromArray bool, index int) string {
	if !fromArray {
		return ""
	}
	return fmt.Sprintf("[%d]", index)
}

func commitPlannedWrites(runAbs, runDisplay string, kind RecordKind, writes []plannedRecordWrite, opts WriteOptions, nextActions []receipt.Action) (*WriteRecordReceipt, error) {
	paths := make([]string, 0, len(writes))
	var createdPtr *bool
	for _, write := range writes {
		paths = append(paths, displayRecordPath(runAbs, runDisplay, write.Path))
		if write.Created != nil {
			createdPtr = write.Created
		}
		if opts.DryRun {
			continue
		}
		var err error
		switch kind {
		case KindAnalysis:
			err = writeReplace(write.Path, write.Data)
		default:
			err = writeCreate(write.Path, write.Data)
		}
		if err != nil {
			return nil, err
		}
	}
	return &WriteRecordReceipt{
		SchemaVersion: SchemaVersion,
		Path:          singlePath(paths),
		Paths:         paths,
		Kind:          kind,
		DryRun:        opts.DryRun,
		Created:       createdPtr,
		NextActions:   nextActions,
	}, nil
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
	var acc validationAccumulator
	acc.Merge("", validateAssessmentResultRequiredStrings(p))
	acc.Merge("", validateRatingResult("ratingResult", &p.RatingResult, levels))
	if p.AreaPath == nil {
		acc.Add("areaPath", "is required")
	}
	if p.FactorPaths == nil {
		acc.Add("factorPaths", "is required")
	}
	if p.Findings == nil {
		acc.Add("findings", "is required")
	}
	if p.Recommendations == nil {
		acc.Add("recommendations", "is required")
	}
	acc.Merge("", validateAssessmentResultFindings(p.Findings))
	acc.Merge("", validateRequiredStrings("supersedes", p.Supersedes))
	return acc.Err()
}

func validateAssessmentResultRequiredStrings(p AssessmentResultInput) error {
	var acc validationAccumulator
	for name, value := range map[string]string{
		"requirement":     p.Requirement,
		"criterionSource": p.CriterionSource,
	} {
		if strings.TrimSpace(value) == "" {
			acc.Add(name, "is required")
		}
	}
	return acc.Err()
}

func validateAssessmentResultFindings(findings []Finding) error {
	var acc validationAccumulator
	for i, finding := range findings {
		prefix := fmt.Sprintf("findings[%d]", i)
		if strings.TrimSpace(finding.Locator) == "" {
			acc.Add(prefix+".locator", "is required")
		}
		if strings.TrimSpace(finding.Observation) == "" {
			acc.Add(prefix+".observation", "is required")
		}
		if strings.TrimSpace(finding.Category) == "" {
			acc.Add(prefix+".category", "is required")
		}
		if strings.TrimSpace(string(finding.Severity)) == "" {
			acc.AddAllowed(prefix+".severity", "is required", []string{"critical", "high", "medium", "low", "info"})
		} else if !finding.Severity.Valid() {
			acc.AddAllowed(prefix+".severity", "is not a supported value", []string{"critical", "high", "medium", "low", "info"})
		}
	}
	return acc.Err()
}

func validateRequiredStrings(name string, values []string) error {
	var acc validationAccumulator
	for i, value := range values {
		if strings.TrimSpace(value) == "" {
			acc.Add(fmt.Sprintf("%s[%d]", name, i), "is required")
		}
	}
	return acc.Err()
}

func validateAnalysis(p AnalysisInput, levels map[string]bool, vocab factorVocabulary) error {
	var acc validationAccumulator
	if p.AreaPath == nil {
		acc.Add("areaPath", "is required")
	}
	acc.Merge("", validateRatingResult("aggregateRatingResult", &p.AggregateRatingResult, levels))
	if p.LocalRatingResult != nil {
		acc.Merge("", validateRatingResult("localRatingResult", p.LocalRatingResult, levels))
	}
	seenFactorPaths := map[string]bool{}
	for i := range p.FactorRatingResults {
		factorPath := p.FactorRatingResults[i].FactorPath
		if len(factorPath) == 0 {
			acc.Add(fmt.Sprintf("factorRatingResults[%d].factorPath", i), "is required")
		}
		key := factorPath.IdentityKey()
		if seenFactorPaths[key] {
			acc.Add(fmt.Sprintf("factorRatingResults[%d].factorPath", i), fmt.Sprintf("%q is a duplicate within the analysis record", factorPath.Display()))
		}
		seenFactorPaths[key] = true
		if !vocab.Resolves(p.AreaPath, factorPath) {
			acc.Add(fmt.Sprintf("factorRatingResults[%d].factorPath", i), fmt.Sprintf("%q does not resolve against the area's declared or inherited factor vocabulary", factorPath.Display()))
		}
		acc.Merge("", validateRatingResult(fmt.Sprintf("factorRatingResults[%d].ratingResult", i), &p.FactorRatingResults[i].RatingResult, levels))
	}
	if p.FactorRatingResults == nil {
		acc.Add("factorRatingResults", "is required")
	}
	if p.AssessmentResultRecords == nil {
		acc.Add("assessmentResultRecords", "is required")
	}
	if p.ChildAnalysisRecords == nil {
		acc.Add("childAnalysisRecords", "is required")
	}
	return acc.Err()
}

func validateRatingResult(name string, result *RatingResult, levels map[string]bool) error {
	var acc validationAccumulator
	if strings.TrimSpace(result.Rationale) == "" {
		acc.Add(name+".rationale", "is required")
	}
	switch result.Kind {
	case RatingResultRated:
		if strings.TrimSpace(result.Level) == "" {
			acc.Add(name+".level", "is required when kind is rated")
		} else if !levels[result.Level] {
			acc.Add(name+".level", fmt.Sprintf("%q is not defined by the run model", result.Level))
		}
	case RatingResultNotAssessed:
		if strings.TrimSpace(result.Level) != "" {
			acc.Add(name+".level", "must be empty when kind is not-assessed")
		}
	default:
		acc.AddAllowed(name+".kind", "is not a supported value", []string{"rated", "not-assessed"})
	}
	return acc.Err()
}

func validateRecommendation(p RecommendationInput) error {
	var acc validationAccumulator
	for name, value := range map[string]string{
		"title":             p.Title,
		"gap":               p.Gap,
		"recommendedOption": p.RecommendedOption,
		"doneCriterion":     p.DoneCriterion,
	} {
		if strings.TrimSpace(value) == "" {
			acc.Add(name, "is required")
		}
	}
	if len(p.EvidenceLocators) == 0 {
		acc.Add("evidenceLocators", "is required")
	}
	if len(p.RemediationOptions) == 0 {
		acc.Add("remediationOptions", "is required")
	}
	acc.Merge("", validateRequiredStrings("supersedes", p.Supersedes))
	return acc.Err()
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
