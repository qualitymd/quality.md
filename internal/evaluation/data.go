package evaluation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
)

// DataKind identifies an Evaluation structured payload kind.
type DataKind string

const (
	DataKindRunManifest                DataKind = "RunManifest"
	DataKindEvaluationFrame            DataKind = "EvaluationFrame"
	DataKindAreaEvaluationFrame        DataKind = "AreaEvaluationFrame"
	DataKindRequirementEvaluationFrame DataKind = "RequirementEvaluationFrame"
	DataKindRequirementAssessment      DataKind = "RequirementAssessmentResult"
	DataKindRequirementRating          DataKind = "RequirementRatingResult"
	DataKindFactorAnalysisFrame        DataKind = "FactorAnalysisFrame"
	DataKindFactorAnalysis             DataKind = "FactorAnalysisResult"
	DataKindAreaAnalysisFrame          DataKind = "AreaAnalysisFrame"
	DataKindAreaAnalysis               DataKind = "AreaAnalysisResult"
	DataKindFindingRanking             DataKind = "FindingRankingResult"
	DataKindRecommendation             DataKind = "RecommendationResult"
	DataKindRecommendationRanking      DataKind = "RecommendationRankingResult"
	DataKindEvaluationOutput           DataKind = "EvaluationOutputResult"
)

// supportedDataKinds lists every Evaluation payload kind the CLI can persist,
// including the CLI-owned EvaluationOutputResult. It is the single typed source
// for the reference-kind vocabulary: any reference may name any of these kinds.
var supportedDataKinds = []DataKind{
	DataKindRunManifest,
	DataKindEvaluationFrame,
	DataKindAreaEvaluationFrame,
	DataKindRequirementEvaluationFrame,
	DataKindRequirementAssessment,
	DataKindRequirementRating,
	DataKindFactorAnalysisFrame,
	DataKindFactorAnalysis,
	DataKindAreaAnalysisFrame,
	DataKindAreaAnalysis,
	DataKindFindingRanking,
	DataKindRecommendation,
	DataKindRecommendationRanking,
	DataKindEvaluationOutput,
}

// acceptedDataKinds is the agent-writable subset of supportedDataKinds: every
// supported kind except the CLI-owned EvaluationOutputResult, which only
// evaluation report build generates. It gates data set payloads.
var acceptedDataKinds = supportedDataKindsExcept(DataKindRunManifest, DataKindEvaluationOutput)

func supportedDataKindsExcept(excluded ...DataKind) []DataKind {
	kinds := make([]DataKind, 0, len(supportedDataKinds))
	for _, kind := range supportedDataKinds {
		if !slices.Contains(excluded, kind) {
			kinds = append(kinds, kind)
		}
	}
	return kinds
}

// kindStrings renders a typed kind slice as plain strings for enum constraints.
func kindStrings[T ~string](kinds []T) []string {
	out := make([]string, len(kinds))
	for i, kind := range kinds {
		out[i] = string(kind)
	}
	return out
}

// DataSetOptions configures Evaluation data writes.
type DataSetOptions struct {
	DryRun bool
}

// DataSetReceipt is emitted after validating or writing Evaluation data.
type DataSetReceipt struct {
	SchemaVersion int              `json:"schemaVersion"`
	Count         int              `json:"count"`
	Writes        []DataSetWrite   `json:"writes"`
	DryRun        bool             `json:"dryRun,omitempty"`
	NextActions   []receipt.Action `json:"nextActions,omitempty"`
}

// DataSetWrite identifies one payload write in a batch receipt.
type DataSetWrite struct {
	Index int      `json:"index"`
	Kind  DataKind `json:"kind"`
	Path  string   `json:"path"`
}

type dataWriteCandidate struct {
	Index     int
	Kind      DataKind
	Path      string
	Canonical []byte
}

// DataKindList lists Evaluation data kinds.
type DataKindList struct {
	SchemaVersion int            `json:"schemaVersion"`
	Kinds         []DataKindInfo `json:"kinds"`
}

// DataKindInfo describes one Evaluation data kind.
type DataKindInfo struct {
	Kind          DataKind `json:"kind"`
	AgentWritable bool     `json:"agentWritable"`
	Description   string   `json:"description"`
}

// DataList lists stored Evaluation data artifacts.
type DataList struct {
	SchemaVersion int             `json:"schemaVersion"`
	Path          string          `json:"path"`
	Artifacts     []DataListEntry `json:"artifacts"`
}

// DataListEntry identifies one stored Evaluation data artifact.
type DataListEntry struct {
	Kind DataKind `json:"kind"`
	Path string   `json:"path"`
}

// DataQuery identifies an Evaluation data artifact by kind and model ID.
type DataQuery struct {
	Kind            DataKind
	AreaRef         string
	FactorRef       string
	RequirementRef  string
	Selector        string
	AllowCLIOwned   bool
	RequireArtifact bool
}

// SetData validates and writes a batch of Evaluation data payloads.
func SetData(runPath string, raw []byte, opts DataSetOptions) (*DataSetReceipt, error) {
	payloads, err := decodeDataPayloadBatch(raw)
	if err != nil {
		return nil, err
	}
	spec, err := loadRunModel(runPath)
	if err != nil {
		return nil, fmt.Errorf("loading %s: %w", ModelSnapshotFile, err)
	}
	candidates, err := dataWriteCandidates(payloads, spec)
	if err != nil {
		return nil, err
	}
	if err := rejectDuplicateCandidatePaths(candidates); err != nil {
		return nil, err
	}
	if err := validateEffectiveDataSet(runPath, candidates); err != nil {
		return nil, err
	}
	if !opts.DryRun {
		if err := writeDataCandidates(runPath, candidates); err != nil {
			return nil, err
		}
	}
	return &DataSetReceipt{
		SchemaVersion: SchemaVersion,
		Count:         len(candidates),
		Writes:        dataSetWrites(candidates),
		DryRun:        opts.DryRun,
		NextActions: []receipt.Action{{
			ID:      "evaluation-status",
			Label:   "Inspect evaluation data status",
			Command: "qualitymd evaluation status " + filepath.ToSlash(runPath),
		}},
	}, nil
}

// ListData lists stored Evaluation data artifacts.
func ListData(runPath string, kind DataKind) (*DataList, error) {
	root := filepath.Join(runPath, "data")
	result := &DataList{SchemaVersion: SchemaVersion, Path: filepath.ToSlash(runPath)}
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return result, nil
	}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil
		}
		payloadKind, err := payloadKind(payload)
		if err != nil {
			return nil
		}
		if kind != "" && payloadKind != kind {
			return nil
		}
		rel, err := filepath.Rel(runPath, path)
		if err != nil {
			return err
		}
		result.Artifacts = append(result.Artifacts, DataListEntry{
			Kind: payloadKind,
			Path: filepath.ToSlash(rel),
		})
		return nil
	}); err != nil {
		return nil, fmt.Errorf("listing evaluation data: %w", err)
	}
	slices.SortFunc(result.Artifacts, func(a, b DataListEntry) int {
		return strings.Compare(a.Path, b.Path)
	})
	return result, nil
}

// GetData reads one stored Evaluation data artifact.
func GetData(runPath string, query DataQuery) ([]byte, string, error) {
	rel, err := dataPathForQuery(query)
	if err != nil {
		return nil, "", err
	}
	raw, err := os.ReadFile(filepath.Join(runPath, rel))
	if err != nil {
		return nil, "", fmt.Errorf("reading %s: %w", filepath.ToSlash(rel), err)
	}
	return raw, filepath.ToSlash(rel), nil
}

// EvaluationDataKinds lists Evaluation data kinds.
func EvaluationDataKinds() *DataKindList {
	var infos []DataKindInfo
	for _, kind := range acceptedDataKinds {
		contract := dataContracts[kind]
		infos = append(infos, DataKindInfo{
			Kind:          kind,
			AgentWritable: true,
			Description:   contract.Description,
		})
	}
	infos = append(infos, DataKindInfo{
		Kind:          DataKindRunManifest,
		AgentWritable: false,
		Description:   dataContracts[DataKindRunManifest].Description,
	})
	infos = append(infos, DataKindInfo{
		Kind:          DataKindEvaluationOutput,
		AgentWritable: false,
		Description:   dataContracts[DataKindEvaluationOutput].Description,
	})
	return &DataKindList{SchemaVersion: SchemaVersion, Kinds: infos}
}

// DataExample returns a complete example JSON artifact for kind.
func DataExample(kind DataKind) ([]byte, error) {
	contract, ok := dataContracts[kind]
	if !ok {
		return nil, usagef("unknown evaluation data kind %q", kind)
	}
	return canonicalJSON(contract.Example())
}

func decodeDataPayload(raw []byte) (map[string]any, error) {
	var payload map[string]any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&payload); err != nil {
		return nil, usagef("invalid JSON payload: %w", err)
	}
	if len(payload) == 0 {
		return nil, usagef("payload must be a JSON object")
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, usagef("payload must contain one JSON object")
	}
	if v, ok := payload["schemaVersion"]; !ok {
		return nil, usagef("payload is missing schemaVersion")
	} else if version, ok := numericSchemaVersion(v); !ok || version != SchemaVersion {
		return nil, usagef("unsupported schemaVersion %v", v)
	}
	return payload, nil
}

func decodeDataPayloadBatch(raw []byte) ([]map[string]any, error) {
	var values []any
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	if err := dec.Decode(&values); err != nil {
		return nil, usagef("invalid JSON payload array: %w", err)
	}
	if len(values) == 0 {
		return nil, usagef("payload batch must contain at least one JSON object")
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, usagef("payload batch must contain one JSON array")
	}
	payloads := make([]map[string]any, 0, len(values))
	var failures []string
	for i, value := range values {
		payload, ok := value.(map[string]any)
		if !ok || len(payload) == 0 {
			failures = append(failures, fmt.Sprintf("payload[%d]: payload must be a JSON object", i))
			continue
		}
		payloads = append(payloads, payload)
	}
	if len(failures) > 0 {
		return nil, batchUsageError(failures)
	}
	return payloads, nil
}

func dataWriteCandidates(payloads []map[string]any, spec *model.Spec) ([]dataWriteCandidate, error) {
	candidates := make([]dataWriteCandidate, 0, len(payloads))
	var failures []string
	for i, payload := range payloads {
		candidate, err := dataWriteCandidateForPayload(i, payload, spec)
		if err != nil {
			failures = append(failures, fmt.Sprintf("payload[%d]: %s", i, err))
			continue
		}
		candidates = append(candidates, candidate)
	}
	if len(failures) > 0 {
		return nil, batchUsageError(failures)
	}
	return candidates, nil
}

func dataWriteCandidateForPayload(index int, payload map[string]any, spec *model.Spec) (dataWriteCandidate, error) {
	if v, ok := payload["schemaVersion"]; !ok {
		return dataWriteCandidate{}, usagef("payload is missing schemaVersion")
	} else if version, ok := numericSchemaVersion(v); !ok || version != SchemaVersion {
		return dataWriteCandidate{}, usagef("unsupported schemaVersion %v", v)
	}
	kind, err := payloadKind(payload)
	if err != nil {
		return dataWriteCandidate{}, err
	}
	if kind == DataKindRunManifest || kind == DataKindEvaluationOutput {
		return dataWriteCandidate{}, usagef("%s is CLI-owned and cannot be written with evaluation data set", kind)
	}
	if !slices.Contains(acceptedDataKinds, kind) {
		return dataWriteCandidate{}, usagef("unsupported evaluation data kind %q", kind)
	}
	if err := validateDataPayloadForModel(kind, payload, spec); err != nil {
		return dataWriteCandidate{}, err
	}
	rel, err := dataPathForPayload(kind, payload)
	if err != nil {
		return dataWriteCandidate{}, err
	}
	canonical, err := canonicalJSON(payload)
	if err != nil {
		return dataWriteCandidate{}, err
	}
	return dataWriteCandidate{Index: index, Kind: kind, Path: filepath.ToSlash(rel), Canonical: canonical}, nil
}

func rejectDuplicateCandidatePaths(candidates []dataWriteCandidate) error {
	firstByPath := map[string]int{}
	var failures []string
	for _, candidate := range candidates {
		if first, ok := firstByPath[candidate.Path]; ok {
			failures = append(failures, fmt.Sprintf("payload[%d] and payload[%d] derive the same path %s", first, candidate.Index, candidate.Path))
			continue
		}
		firstByPath[candidate.Path] = candidate.Index
	}
	if len(failures) > 0 {
		return batchUsageError(failures)
	}
	return nil
}

func validateEffectiveDataSet(runPath string, candidates []dataWriteCandidate) error {
	payloads, err := effectiveDataPayloads(runPath, candidates)
	if err != nil {
		return err
	}
	failures := validateEffectivePayloads(payloads)
	if len(failures) == 0 {
		return nil
	}
	out := make([]string, len(failures))
	for i, failure := range failures {
		out[i] = failure.Path + ": " + failure.Reason
	}
	return batchUsageError(out)
}

type effectiveDataFailure struct {
	Path   string
	Kind   DataKind
	Reason string
}

func effectiveDataPayloads(runPath string, candidates []dataWriteCandidate) (map[string]map[string]any, error) {
	payloads := map[string]map[string]any{}
	root := filepath.Join(runPath, "data")
	if _, err := os.Stat(root); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
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
			return usagef("%s: %w", filepath.ToSlash(rel), err)
		}
		payload, err := decodeDataPayload(raw)
		if err != nil {
			return usagef("%s: %w", filepath.ToSlash(rel), err)
		}
		payloads[filepath.ToSlash(rel)] = payload
		return nil
	}); err != nil {
		return nil, err
	}
	for _, candidate := range candidates {
		payload, err := decodeDataPayload(candidate.Canonical)
		if err != nil {
			return nil, usagef("%s: %w", candidate.Path, err)
		}
		payloads[candidate.Path] = payload
	}
	return payloads, nil
}

func validateEffectivePayloads(payloads map[string]map[string]any) []effectiveDataFailure {
	var failures []effectiveDataFailure
	for path, payload := range payloads {
		kind, err := payloadKind(payload)
		if err != nil {
			failures = append(failures, effectiveDataFailure{Path: path, Reason: err.Error()})
			continue
		}
		switch kind {
		case DataKindRequirementRating:
			failures = append(failures, validateEffectiveRequirementRating(path, payload, payloads)...)
		case DataKindFactorAnalysis, DataKindAreaAnalysis:
			failures = append(failures, validateEffectiveAnalysisResult(path, kind, payload, payloads)...)
		case DataKindFindingRanking:
			failures = append(failures, validateEffectiveFindingRanking(path, payload, payloads)...)
		case DataKindRecommendation:
			failures = append(failures, validateEffectiveRecommendation(path, payload, payloads)...)
		case DataKindRecommendationRanking:
			failures = append(failures, validateEffectiveRecommendationRanking(path, payload, payloads)...)
		}
	}
	return failures
}

func validateEffectiveRequirementRating(path string, payload map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	if firstString(payload, "status") != "rated" {
		return nil
	}
	var failures []effectiveDataFailure
	req, err := topRequirementID(payload)
	if err != nil {
		return []effectiveDataFailure{{Path: path, Kind: DataKindRequirementRating, Reason: err.Error()}}
	}
	assessmentPath := requirementDataPath(req, "requirement-assessment-result.json")
	assessment, ok := payloads[assessmentPath]
	if !ok {
		failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRequirementRating, Reason: "rated Requirement requires paired RequirementAssessmentResult with at least one finding"})
	} else {
		status := firstString(assessment, "status")
		if status != "assessed" && status != "partially_assessed" {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRequirementRating, Reason: "rated Requirement requires paired RequirementAssessmentResult status assessed or partially_assessed"})
		}
		if len(objectSlice(assessment["findings"])) == 0 {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRequirementRating, Reason: "rated Requirement requires paired RequirementAssessmentResult with at least one finding"})
		}
	}
	if len(objectSlice(payload["ratingDrivers"])) == 0 {
		failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRequirementRating, Reason: "rated Requirement requires at least one ratingDrivers entry"})
	}
	failures = append(failures, validateRatingDriverInputRefs(path, DataKindRequirementRating, payload, payloads)...)
	return failures
}

func validateEffectiveAnalysisResult(path string, kind DataKind, payload map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	var failures []effectiveDataFailure
	for _, scopeName := range []string{"localAnalysis", "localAndDescendantAnalysis"} {
		scope := objectMap(payload[scopeName])
		if firstString(scope, "status") != "analyzed" || firstString(scope, "ratingLevelId") == "" {
			continue
		}
		if len(objectSlice(scope["ratingDrivers"])) == 0 {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: scopeName + " with a ratingLevelId requires at least one ratingDrivers entry"})
			continue
		}
		failures = append(failures, validateRatingDriverInputRefs(path+"#"+scopeName, kind, scope, payloads)...)
	}
	return failures
}

func validateRatingDriverInputRefs(path string, kind DataKind, owner map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	var failures []effectiveDataFailure
	for i, driver := range objectSlice(owner["ratingDrivers"]) {
		inputRefs := objectSlice(driver["inputRefs"])
		if len(inputRefs) == 0 {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: fmt.Sprintf("ratingDrivers[%d] requires at least one inputRefs entry", i)})
			continue
		}
		for j, ref := range inputRefs {
			refPath, err := dataPathForRoutineRef(ref)
			if err != nil {
				failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: fmt.Sprintf("ratingDrivers[%d].inputRefs[%d]: %s", i, j, err)})
				continue
			}
			if _, ok := payloads[refPath]; !ok {
				failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: fmt.Sprintf("ratingDrivers[%d].inputRefs[%d] does not resolve to %s", i, j, refPath)})
			}
		}
	}
	return failures
}

func validateEffectiveFindingRanking(path string, payload map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	expected := effectiveFindingRefs(payloads)
	seen := map[string]struct{}{}
	var failures []effectiveDataFailure
	for i, entry := range objectSlice(payload["orderedFindings"]) {
		ref := objectMap(entry["findingRef"])
		key, err := findingRefKey(ref, payloads)
		if err != nil {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindFindingRanking, Reason: fmt.Sprintf("orderedFindings[%d].findingRef: %s", i, err)})
			continue
		}
		if _, ok := expected[key]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindFindingRanking, Reason: fmt.Sprintf("orderedFindings[%d].findingRef does not resolve to an in-scope Finding", i)})
			continue
		}
		if _, ok := seen[key]; ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindFindingRanking, Reason: fmt.Sprintf("orderedFindings[%d].findingRef duplicates %s", i, key)})
			continue
		}
		seen[key] = struct{}{}
	}
	for key := range expected {
		if _, ok := seen[key]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindFindingRanking, Reason: fmt.Sprintf("orderedFindings missing %s", key)})
		}
	}
	failures = append(failures, validateUniquePositiveRanks(path, DataKindFindingRanking, "orderedFindings", objectSlice(payload["orderedFindings"]))...)
	return failures
}

func validateEffectiveRecommendation(path string, payload map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	var failures []effectiveDataFailure
	for i, ref := range objectSlice(payload["traceRefs"]) {
		if err := validateAdviceTraceRef(ref, payloads); err != nil {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendation, Reason: fmt.Sprintf("traceRefs[%d]: %s", i, err)})
		}
	}
	return failures
}

func validateEffectiveRecommendationRanking(path string, payload map[string]any, payloads map[string]map[string]any) []effectiveDataFailure {
	recommendations := effectiveRecommendationIDs(payloads)
	var failures []effectiveDataFailure
	failures = append(failures, validateRankedRecommendations(path, payload, recommendations)...)
	failures = append(failures, validateFindingCoverage(path, payload, payloads, recommendations, effectiveFindingRefs(payloads))...)
	return failures
}

func validateRankedRecommendations(path string, payload map[string]any, recommendations map[string]string) []effectiveDataFailure {
	ranked := map[string]struct{}{}
	var failures []effectiveDataFailure
	for i, entry := range objectSlice(payload["orderedRecommendations"]) {
		ref := firstString(entry, "recommendationRef")
		if _, ok := recommendations[ref]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("orderedRecommendations[%d].recommendationRef does not resolve to a RecommendationResult", i)})
			continue
		}
		if _, ok := ranked[ref]; ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("orderedRecommendations[%d].recommendationRef duplicates %s", i, ref)})
			continue
		}
		ranked[ref] = struct{}{}
	}
	for id := range recommendations {
		if _, ok := ranked[id]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("orderedRecommendations missing %s", id)})
		}
	}
	failures = append(failures, validateUniquePositiveRanks(path, DataKindRecommendationRanking, "orderedRecommendations", objectSlice(payload["orderedRecommendations"]))...)
	return failures
}

func validateFindingCoverage(path string, payload map[string]any, payloads map[string]map[string]any, recommendations map[string]string, expectedFindings map[string]struct{}) []effectiveDataFailure {
	covered := map[string]struct{}{}
	var failures []effectiveDataFailure
	for i, entry := range objectSlice(payload["findingCoverage"]) {
		ref := objectMap(entry["findingRef"])
		key, err := findingRefKey(ref, payloads)
		if err != nil {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].findingRef: %s", i, err)})
			continue
		}
		if _, ok := expectedFindings[key]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].findingRef does not resolve to an in-scope Finding", i)})
			continue
		}
		if _, ok := covered[key]; ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].findingRef duplicates %s", i, key)})
			continue
		}
		covered[key] = struct{}{}
		failures = append(failures, validateCoverageDisposition(path, i, entry, recommendations)...)
	}
	for key := range expectedFindings {
		if _, ok := covered[key]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage missing %s", key)})
		}
	}
	return failures
}

func validateCoverageDisposition(path string, index int, entry map[string]any, recommendations map[string]string) []effectiveDataFailure {
	switch firstString(entry, "disposition") {
	case "addressed_by_recommendation":
		return validateAddressedCoverageRefs(path, index, entry, recommendations)
	case "not_advice_driving":
		if firstString(entry, "rationale") == "" {
			return []effectiveDataFailure{{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].rationale is required for not_advice_driving", index)}}
		}
	}
	return nil
}

func validateAddressedCoverageRefs(path string, index int, entry map[string]any, recommendations map[string]string) []effectiveDataFailure {
	refs := stringValues(entry["recommendationRefs"])
	if len(refs) == 0 {
		return []effectiveDataFailure{{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].recommendationRefs requires at least one RecommendationResult ID", index)}}
	}
	var failures []effectiveDataFailure
	for _, rec := range refs {
		if _, ok := recommendations[rec]; !ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: DataKindRecommendationRanking, Reason: fmt.Sprintf("findingCoverage[%d].recommendationRefs includes unknown RecommendationResult %s", index, rec)})
		}
	}
	return failures
}

func validateUniquePositiveRanks(path string, kind DataKind, field string, entries []map[string]any) []effectiveDataFailure {
	seen := map[int]struct{}{}
	var failures []effectiveDataFailure
	for i, entry := range entries {
		rank, ok := rankField(entry)
		if !ok || rank < 1 {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: fmt.Sprintf("%s[%d].rank must be a positive integer", field, i)})
			continue
		}
		if _, ok := seen[rank]; ok {
			failures = append(failures, effectiveDataFailure{Path: path, Kind: kind, Reason: fmt.Sprintf("%s[%d].rank duplicates %d", field, i, rank)})
			continue
		}
		seen[rank] = struct{}{}
	}
	return failures
}

func effectiveFindingRefs(payloads map[string]map[string]any) map[string]struct{} {
	out := map[string]struct{}{}
	for path, payload := range payloads {
		kind, err := payloadKind(payload)
		if err != nil || kind != DataKindRequirementAssessment {
			continue
		}
		for _, finding := range objectSlice(payload["findings"]) {
			id := firstString(finding, "id")
			if id != "" {
				out[path+"#findings["+id+"]"] = struct{}{}
			}
		}
	}
	return out
}

func effectiveRecommendationIDs(payloads map[string]map[string]any) map[string]string {
	out := map[string]string{}
	for path, payload := range payloads {
		kind, err := payloadKind(payload)
		if err != nil || kind != DataKindRecommendation {
			continue
		}
		id := firstString(payload, "id")
		if id != "" {
			out[id] = path
		}
	}
	return out
}

func findingRefKey(ref map[string]any, payloads map[string]map[string]any) (string, error) {
	path, err := dataPathForRoutineRef(ref)
	if err != nil {
		return "", err
	}
	if DataKind(firstString(ref, "kind")) != DataKindRequirementAssessment {
		return "", usagef("kind must be %s", DataKindRequirementAssessment)
	}
	selector := firstString(ref, "selector")
	if !strings.HasPrefix(selector, "findings[") || !strings.HasSuffix(selector, "]") {
		return "", usagef("selector must look like findings[<id>]")
	}
	payload, ok := payloads[path]
	if !ok {
		return "", usagef("referenced assessment %s is missing", path)
	}
	want := strings.TrimSuffix(strings.TrimPrefix(selector, "findings["), "]")
	for _, finding := range objectSlice(payload["findings"]) {
		if firstString(finding, "id") == want {
			return path + "#" + selector, nil
		}
	}
	return "", usagef("selector %s does not resolve in %s", selector, path)
}

func validateAdviceTraceRef(ref map[string]any, payloads map[string]map[string]any) error {
	kind := DataKind(firstString(ref, "kind"))
	if kind == DataKindRequirementAssessment && firstString(ref, "selector") != "" {
		_, err := findingRefKey(ref, payloads)
		return err
	}
	path, err := dataPathForRoutineRef(ref)
	if err != nil {
		return err
	}
	if _, ok := payloads[path]; !ok {
		return usagef("does not resolve to %s", path)
	}
	return nil
}

//nolint:cyclop // Centralized kind-to-path dispatch keeps data path ownership explicit.
func dataPathForRoutineRef(ref map[string]any) (string, error) {
	kindValue := firstString(ref, "kind")
	if kindValue == "" {
		return "", usagef("missing kind")
	}
	kind := DataKind(kindValue)
	subject := objectMap(ref["subject"])
	payload := map[string]any{}
	switch kind {
	case DataKindRunManifest:
		return "data/run-manifest.json", nil
	case DataKindEvaluationFrame:
	case DataKindAreaEvaluationFrame, DataKindAreaAnalysisFrame:
		payload["subject"] = subject
	case DataKindAreaAnalysis:
		payload["areaId"] = subject["areaId"]
	case DataKindRequirementEvaluationFrame:
		payload["subject"] = subject
	case DataKindRequirementAssessment, DataKindRequirementRating:
		payload["requirementId"] = subject["requirementId"]
	case DataKindFactorAnalysisFrame:
		payload["subject"] = subject
	case DataKindFactorAnalysis:
		payload["factorId"] = subject["factorId"]
	case DataKindFindingRanking:
		return "data/advice/finding-ranking-result.json", nil
	case DataKindRecommendation:
		id := firstString(ref, "id", "recommendationId")
		if id == "" {
			id = firstString(subject, "recommendationId")
		}
		if id == "" {
			return "", usagef("missing recommendation id")
		}
		return recommendationDataPath(id), nil
	case DataKindRecommendationRanking:
		return "data/advice/recommendation-ranking-result.json", nil
	case DataKindEvaluationOutput:
		return "data/evaluation-output-result.json", nil
	default:
		return "", usagef("unsupported evaluation data kind %q", kind)
	}
	return dataPathForPayload(kind, payload)
}

func dataSetWrites(candidates []dataWriteCandidate) []DataSetWrite {
	writes := make([]DataSetWrite, len(candidates))
	for i, candidate := range candidates {
		writes[i] = DataSetWrite{Index: candidate.Index, Kind: candidate.Kind, Path: candidate.Path}
	}
	return writes
}

type committedDataWrite struct {
	target      string
	backup      string
	backupMoved bool
	targetWrote bool
}

func writeDataCandidates(runPath string, candidates []dataWriteCandidate) error {
	tempDir, err := os.MkdirTemp(runPath, ".data-set-*")
	if err != nil {
		return fmt.Errorf("creating data staging directory: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	stageRoot := filepath.Join(tempDir, "stage")
	backupRoot := filepath.Join(tempDir, "backup")
	for _, candidate := range candidates {
		stagePath := filepath.Join(stageRoot, filepath.FromSlash(candidate.Path))
		if err := os.MkdirAll(filepath.Dir(stagePath), 0o755); err != nil {
			return fmt.Errorf("creating staged data directory: %w", err)
		}
		if err := os.WriteFile(stagePath, candidate.Canonical, 0o644); err != nil {
			return fmt.Errorf("staging %s: %w", candidate.Path, err)
		}
	}

	var committed []committedDataWrite
	for _, candidate := range candidates {
		target := filepath.Join(runPath, filepath.FromSlash(candidate.Path))
		stagePath := filepath.Join(stageRoot, filepath.FromSlash(candidate.Path))
		state := committedDataWrite{
			target: target,
			backup: filepath.Join(backupRoot, filepath.FromSlash(candidate.Path)),
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			rollbackCommittedDataWrites(committed)
			return fmt.Errorf("creating data directory: %w", err)
		}
		if _, err := os.Stat(target); err == nil {
			if err := os.MkdirAll(filepath.Dir(state.backup), 0o755); err != nil {
				rollbackCommittedDataWrites(committed)
				return fmt.Errorf("creating data backup directory: %w", err)
			}
			if err := os.Rename(target, state.backup); err != nil {
				rollbackCommittedDataWrites(committed)
				return fmt.Errorf("backing up %s: %w", candidate.Path, err)
			}
			state.backupMoved = true
		} else if !os.IsNotExist(err) {
			rollbackCommittedDataWrites(committed)
			return fmt.Errorf("checking %s: %w", candidate.Path, err)
		}
		if err := os.Rename(stagePath, target); err != nil {
			rollbackCommittedDataWrites(append(committed, state))
			return fmt.Errorf("writing %s: %w", candidate.Path, err)
		}
		state.targetWrote = true
		committed = append(committed, state)
	}
	return nil
}

func rollbackCommittedDataWrites(committed []committedDataWrite) {
	for i := len(committed) - 1; i >= 0; i-- {
		state := committed[i]
		if state.targetWrote {
			_ = os.Remove(state.target)
		}
		if state.backupMoved {
			_ = os.MkdirAll(filepath.Dir(state.target), 0o755)
			_ = os.Rename(state.backup, state.target)
		}
	}
}

func batchUsageError(failures []string) error {
	return usagef("invalid evaluation data batch:\n- %s", strings.Join(failures, "\n- "))
}

func canonicalJSON(v any) ([]byte, error) {
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(raw, '\n'), nil
}

func numericSchemaVersion(v any) (int, bool) {
	switch n := v.(type) {
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	case float64:
		return int(n), n == float64(int(n))
	default:
		return 0, false
	}
}

func payloadKind(payload map[string]any) (DataKind, error) {
	raw, ok := payload["kind"].(string)
	if !ok || raw == "" {
		return "", usagef("payload is missing kind")
	}
	return DataKind(raw), nil
}

//nolint:cyclop // Centralized kind-to-path dispatch keeps payload routing explicit.
func dataPathForPayload(kind DataKind, payload map[string]any) (string, error) {
	switch kind {
	case DataKindRunManifest:
		return "data/run-manifest.json", nil
	case DataKindEvaluationFrame:
		return "data/frame/evaluation-frame.json", nil
	case DataKindAreaEvaluationFrame:
		areaID, err := subjectAreaID(payload)
		return areaDataPath(areaID, "area-evaluation-frame.json"), err
	case DataKindAreaAnalysisFrame:
		areaID, err := subjectAreaID(payload)
		return areaDataPath(areaID, "area-analysis-frame.json"), err
	case DataKindAreaAnalysis:
		areaID, err := topAreaID(payload)
		return areaDataPath(areaID, "area-analysis-result.json"), err
	case DataKindRequirementEvaluationFrame:
		req, err := subjectRequirementID(payload)
		return requirementDataPath(req, "requirement-evaluation-frame.json"), err
	case DataKindRequirementAssessment:
		req, err := topRequirementID(payload)
		return requirementDataPath(req, "requirement-assessment-result.json"), err
	case DataKindRequirementRating:
		req, err := topRequirementID(payload)
		return requirementDataPath(req, "requirement-rating-result.json"), err
	case DataKindFactorAnalysisFrame:
		factor, err := subjectFactorID(payload)
		return factorDataPath(factor, "factor-analysis-frame.json"), err
	case DataKindFactorAnalysis:
		factor, err := topFactorID(payload)
		return factorDataPath(factor, "factor-analysis-result.json"), err
	case DataKindFindingRanking:
		return "data/advice/finding-ranking-result.json", nil
	case DataKindRecommendation:
		id := firstString(payload, "id")
		if !safeModelName(id) {
			return "", usagef("RecommendationResult.id must be a path-safe non-empty string")
		}
		return recommendationDataPath(id), nil
	case DataKindRecommendationRanking:
		return "data/advice/recommendation-ranking-result.json", nil
	default:
		return "", usagef("unsupported evaluation data kind %q", kind)
	}
}

//nolint:cyclop // Centralized kind-to-path dispatch keeps query routing explicit.
func dataPathForQuery(query DataQuery) (string, error) {
	switch query.Kind {
	case DataKindRunManifest:
		return "data/run-manifest.json", nil
	case DataKindEvaluationFrame:
		return "data/frame/evaluation-frame.json", nil
	case DataKindAreaEvaluationFrame:
		areaID, err := parseAreaRef(query.AreaRef)
		return areaDataPath(areaID, "area-evaluation-frame.json"), err
	case DataKindAreaAnalysisFrame:
		areaID, err := parseAreaRef(query.AreaRef)
		return areaDataPath(areaID, "area-analysis-frame.json"), err
	case DataKindAreaAnalysis:
		areaID, err := parseAreaRef(query.AreaRef)
		return areaDataPath(areaID, "area-analysis-result.json"), err
	case DataKindRequirementEvaluationFrame:
		req, err := parseRequirementRef(query.RequirementRef)
		return requirementDataPath(req, "requirement-evaluation-frame.json"), err
	case DataKindRequirementAssessment:
		req, err := parseRequirementRef(query.RequirementRef)
		return requirementDataPath(req, "requirement-assessment-result.json"), err
	case DataKindRequirementRating:
		req, err := parseRequirementRef(query.RequirementRef)
		return requirementDataPath(req, "requirement-rating-result.json"), err
	case DataKindFactorAnalysisFrame:
		factor, err := parseFactorRef(query.FactorRef)
		return factorDataPath(factor, "factor-analysis-frame.json"), err
	case DataKindFactorAnalysis:
		factor, err := parseFactorRef(query.FactorRef)
		return factorDataPath(factor, "factor-analysis-result.json"), err
	case DataKindFindingRanking:
		return "data/advice/finding-ranking-result.json", nil
	case DataKindRecommendation:
		if query.Selector == "" {
			return "", usagef("--selector must name the RecommendationResult id")
		}
		return recommendationDataPath(query.Selector), nil
	case DataKindRecommendationRanking:
		return "data/advice/recommendation-ranking-result.json", nil
	case DataKindEvaluationOutput:
		return "data/evaluation-output-result.json", nil
	default:
		return "", usagef("--kind must be a known evaluation data kind")
	}
}

type requirementID struct {
	DeclaringArea []string
	Name          string
}

type factorID struct {
	DeclaringArea []string
	Path          []string
}

func subjectMap(payload map[string]any) (map[string]any, error) {
	subject, ok := payload["subject"].(map[string]any)
	if !ok {
		return nil, usagef("payload subject must be an object")
	}
	return subject, nil
}

func subjectAreaID(payload map[string]any) ([]string, error) {
	subject, err := subjectMap(payload)
	if err != nil {
		return nil, err
	}
	return areaIDFrom(subject["areaId"])
}

func topAreaID(payload map[string]any) ([]string, error) {
	return areaIDFrom(payload["areaId"])
}

func subjectRequirementID(payload map[string]any) (requirementID, error) {
	subject, err := subjectMap(payload)
	if err != nil {
		return requirementID{}, err
	}
	return requirementIDFrom(subject["requirementId"])
}

func topRequirementID(payload map[string]any) (requirementID, error) {
	return requirementIDFrom(payload["requirementId"])
}

func subjectFactorID(payload map[string]any) (factorID, error) {
	subject, err := subjectMap(payload)
	if err != nil {
		return factorID{}, err
	}
	return factorIDFrom(subject["factorId"])
}

func topFactorID(payload map[string]any) (factorID, error) {
	return factorIDFrom(payload["factorId"])
}

func areaIDFrom(v any) ([]string, error) {
	ref, ok := v.(string)
	if !ok {
		return nil, usagef("areaId must be a qualified area reference string")
	}
	if !strings.HasPrefix(ref, "area:") {
		return nil, usagef("areaId must start with area:")
	}
	return parseAreaRef(ref)
}

func requirementIDFrom(v any) (requirementID, error) {
	ref, ok := v.(string)
	if !ok {
		return requirementID{}, usagef("requirementId must be a qualified requirement reference string")
	}
	if !strings.HasPrefix(ref, "requirement:") {
		return requirementID{}, usagef("requirementId must start with requirement:")
	}
	return parseRequirementRef(ref)
}

func factorIDFrom(v any) (factorID, error) {
	ref, ok := v.(string)
	if !ok {
		return factorID{}, usagef("factorId must be a qualified factor reference string")
	}
	if !strings.HasPrefix(ref, "factor:") {
		return factorID{}, usagef("factorId must start with factor:")
	}
	return parseFactorRef(ref)
}

func safeModelName(name string) bool {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return false
	}
	for _, r := range name {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func areaDataPath(areaID []string, file string) string {
	parts := []string{"data", "areas"}
	if len(areaID) == 0 {
		parts = append(parts, "root")
	} else {
		parts = append(parts, areaID...)
	}
	parts = append(parts, file)
	return filepath.ToSlash(filepath.Join(parts...))
}

func requirementDataPath(req requirementID, file string) string {
	parts := areaBaseParts(req.DeclaringArea)
	parts = append(parts, "requirements", req.Name, file)
	return filepath.ToSlash(filepath.Join(parts...))
}

func factorDataPath(factor factorID, file string) string {
	parts := areaBaseParts(factor.DeclaringArea)
	for _, name := range factor.Path {
		parts = append(parts, "factors", name)
	}
	parts = append(parts, file)
	return filepath.ToSlash(filepath.Join(parts...))
}

func recommendationDataPath(id string) string {
	return filepath.ToSlash(filepath.Join("data", "advice", "recommendations", id, "recommendation-result.json"))
}

func rankField(m map[string]any) (int, bool) {
	value, ok := m["rank"]
	if !ok {
		return 0, false
	}
	switch n := value.(type) {
	case json.Number:
		i, err := n.Int64()
		return int(i), err == nil
	case float64:
		i := int(n)
		return i, n == float64(i)
	default:
		return 0, false
	}
}

func areaBaseParts(areaID []string) []string {
	parts := []string{"data", "areas"}
	if len(areaID) == 0 {
		return append(parts, "root")
	}
	return append(parts, areaID...)
}

func parseAreaRef(ref string) ([]string, error) {
	ref = strings.TrimPrefix(ref, "area:")
	if ref == "" || ref == "root" || ref == "/" {
		return []string{}, nil
	}
	parts := strings.Split(ref, "/")
	for _, part := range parts {
		if !safeModelName(part) {
			return nil, usagef("invalid area ref %q", ref)
		}
	}
	return parts, nil
}

func parseRequirementRef(ref string) (requirementID, error) {
	ref = strings.TrimPrefix(ref, "requirement:")
	areaPart, name, ok := strings.Cut(ref, "::")
	if !ok || name == "" {
		return requirementID{}, usagef("requirement ref must look like requirement:<area-path>::<requirement-name>")
	}
	area, err := parseAreaRef(areaPart)
	if err != nil {
		return requirementID{}, err
	}
	if !safeModelName(name) {
		return requirementID{}, usagef("invalid requirement name %q", name)
	}
	return requirementID{DeclaringArea: area, Name: name}, nil
}

func parseFactorRef(ref string) (factorID, error) {
	ref = strings.TrimPrefix(ref, "factor:")
	areaPart, factorPart, ok := strings.Cut(ref, "::")
	if !ok || factorPart == "" {
		return factorID{}, usagef("factor ref must look like factor:<area-path>::<factor-path>")
	}
	area, err := parseAreaRef(areaPart)
	if err != nil {
		return factorID{}, err
	}
	path := strings.Split(factorPart, "/")
	for _, part := range path {
		if !safeModelName(part) {
			return factorID{}, usagef("invalid factor path %q", factorPart)
		}
	}
	return factorID{DeclaringArea: area, Path: path}, nil
}

func runManifestExample() map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(DataKindRunManifest),
		"number":        1,
		"model":         "QUALITY.md",
		"requestedScope": map[string]any{
			"areaId":       exampleAreaID(),
			"factorFilter": []any{exampleFactorID()},
		},
		"plannedScope": map[string]any{
			"areaId":       exampleAreaID(),
			"factorFilter": []any{exampleFactorID()},
		},
	}
}

func evaluationFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject":       map[string]any{"modelLocator": "QUALITY.md"},
		"inputs": map[string]any{
			"ratingLevelIds": []any{RatingReference("target"), RatingReference("unacceptable")},
		},
		"derivedContext": map[string]any{
			"rigor":                    "standard",
			"evaluationPolicies":       []any{"source-as-data", "secret-redaction"},
			"expectedEvaluationLimits": []any{exampleLimit("source-gap", "One source was unavailable.", "Confidence is lower.")},
		},
	}
}

func areaEvaluationFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject":       map[string]any{"areaId": exampleAreaID()},
		"inputs": map[string]any{
			"sourceRefs":          []any{"QUALITY.md", "tests/example_test.go"},
			"localRequirementIds": []any{exampleRequirementID()},
			"rootFactorIds":       []any{exampleFactorID()},
			"childAreaIds":        []any{exampleChildAreaID()},
		},
		"derivedContext": map[string]any{
			"scope":                    "root Area",
			"expectedEvaluationLimits": []any{exampleLimit("area-source-gap", "A child Area source was not inspected.", "Child roll-up confidence is lower.")},
		},
	}
}

func requirementEvaluationFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject": map[string]any{
			"requirementId": exampleRequirementID(),
			"factorIds":     []any{exampleFactorID()},
		},
		"inputs": map[string]any{
			"ratingLevelIds":             []any{RatingReference("target"), RatingReference("unacceptable")},
			"requirementAssessmentBasis": "Inspect tests.",
			"ratingOverrides":            map[string]any{"target": "Focused tests cover the stated behavior."},
		},
		"derivedContext": map[string]any{
			"evidenceTargets":          []any{map[string]any{"id": "tests", "question": "Do tests exist?", "purpose": "Assessment basis", "sourceRefs": []any{"tests/example_test.go"}, "required": true}},
			"appliedRatingCriteria":    []any{map[string]any{"ratingLevelId": RatingReference("target"), "criterion": "Tests cover the requirement.", "source": "model_default", "adaptationRationale": "Applies the model criterion to this Requirement."}},
			"stopConditions":           []any{exampleLimit("missing-required-evidence", "Stop if required evidence cannot be inspected.", "Requirement status becomes blocked or partial.")},
			"expectedEvaluationLimits": []any{exampleLimit("narrow-test-review", "Only test files are in scope for this example.", "Runtime behavior may remain uninspected.")},
		},
	}
}

func requirementAssessmentExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion":          SchemaVersion,
		"kind":                   string(kind),
		"requirementId":          exampleRequirementID(),
		"status":                 "assessed",
		"statusReason":           "Evidence target was addressed.",
		"evidenceSummary":        "A focused test was found.",
		"evidenceTargetCoverage": []any{map[string]any{"id": "tests", "description": "Test evidence was inspected."}},
		"findings": []any{
			map[string]any{
				"id":         "strength-001",
				"type":       "strength",
				"severity":   "medium",
				"confidence": "medium",
				"statement":  "Focused test coverage is present.",
				"condition":  "A focused test covers the requirement's primary path.",
				"criteria": []any{map[string]any{
					"requirementId": exampleRequirementID(),
					"ratingLevelId": RatingReference("target"),
					"criterion":     "Tests cover the requirement.",
					"rationale":     "The target criterion is the relevant bar for the observed test coverage.",
				}},
				"basis": map[string]any{
					"status":    "not_applicable",
					"statement": "No separate basis beyond the cited evidence is claimed for this strength.",
				},
				"effect": map[string]any{
					"statement":    "The evidence supports the target rating for the primary path.",
					"ratingEffect": "supports target",
				},
				"evidence": []any{map[string]any{
					"sourceRef": "tests/example_test.go",
					"statement": "A focused test file exists for the requirement.",
					"rationale": "The evidence directly addresses the target.",
				}},
			},
			map[string]any{
				"id":         "gap-001",
				"type":       "gap",
				"severity":   "medium",
				"confidence": "medium",
				"statement":  "Edge-case behavior is not covered by tests.",
				"condition":  "Edge-case paths around the requirement lack visible tests.",
				"criteria": []any{map[string]any{
					"requirementId": exampleRequirementID(),
					"ratingLevelId": RatingReference("target"),
					"criterion":     "Tests cover the requirement, including meaningful boundary behavior.",
					"rationale":     "Boundary behavior is needed to satisfy the target criterion.",
				}},
				"basis": map[string]any{
					"status":    "plausible",
					"statement": "The visible tests focus on the primary path.",
					"rationale": "No broader test inventory was inspected in this example.",
					"evidence": []any{map[string]any{
						"sourceRef": "tests/example_test.go",
						"statement": "The example test evidence is narrow.",
					}},
				},
				"effect": map[string]any{
					"statement":    "The requirement cannot be rated above acceptable without boundary evidence.",
					"ratingEffect": "holds below target",
					"rationale":    "The observed gap prevents distinguishing target from acceptable.",
				},
				"evidence": []any{map[string]any{
					"sourceRef": "tests/example_test.go",
					"statement": "Only narrow test evidence was reviewed.",
				}},
				"candidateActions": []any{map[string]any{
					"id":          "action-001",
					"description": "Add focused tests for the boundary and error paths.",
					"rationale":   "Closing the untested edge cases would lift coverage of the requirement.",
				}},
			},
		},
		"unknowns":         []any{map[string]any{"id": "coverage-depth", "description": "Broader edge-case coverage was not inspected."}},
		"evaluationLimits": []any{map[string]any{"id": "test-only", "description": "Only test evidence was reviewed.", "impact": "Confidence remains medium."}},
		"confidence":       "medium",
		"confidenceReason": "Evidence is relevant but narrow.",
	}
}

func requirementRatingExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion":    SchemaVersion,
		"kind":             string(kind),
		"requirementId":    exampleRequirementID(),
		"status":           "rated",
		"statusReason":     "Assessment maps to a Rating Level.",
		"ratingLevelId":    RatingReference("target"),
		"rationale":        "Evidence satisfies the target criterion.",
		"ratingDrivers":    []any{map[string]any{"description": "Focused evidence satisfies the target criterion.", "effect": "supports target", "ratingLevelId": RatingReference("target"), "inputRefs": []any{routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "")}}},
		"criteriaResults":  []any{map[string]any{"ratingLevelId": RatingReference("target"), "matched": true, "rationale": "Criterion is satisfied."}},
		"missingEvidence":  []any{map[string]any{"id": "edge-cases", "description": "Edge-case evidence was not reviewed."}},
		"evaluationLimits": []any{map[string]any{"id": "narrow-review", "description": "Review was narrow.", "impact": "Confidence remains medium."}},
		"confidence":       "medium",
		"confidenceReason": "Rating follows a narrow assessment.",
	}
}

func factorAnalysisFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject":       map[string]any{"areaId": exampleAreaID(), "factorId": exampleFactorID()},
		"inputs": map[string]any{
			"directRequirementRatingRefs": []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"childFactorAnalysisRefs":     []any{routineRef(DataKindFactorAnalysis, map[string]any{"factorId": exampleChildFactorID()}, "localAndDescendantAnalysis")},
		},
		"derivedContext": map[string]any{
			"synthesisGuidanceRef":     "protocol:factor-synthesis-default-v0",
			"emptySignalPolicy":        "ignore_empty",
			"stopConditions":           []any{exampleLimit("missing-ratings", "Stop if direct Requirement ratings are unavailable.", "Factor analysis becomes incomplete.")},
			"expectedEvaluationLimits": []any{exampleLimit("factor-source-limit", "Only direct Requirement ratings and child Factor analyses are synthesized.", "Other evidence is out of scope.")},
		},
	}
}

func areaAnalysisFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject":       map[string]any{"areaId": exampleAreaID()},
		"inputs": map[string]any{
			"factorAnalysisRefs":    []any{routineRef(DataKindFactorAnalysis, map[string]any{"factorId": exampleFactorID()}, "localAndDescendantAnalysis")},
			"childAreaAnalysisRefs": []any{routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleChildAreaID()}, "localAndDescendantAnalysis")},
		},
		"derivedContext": map[string]any{
			"synthesisGuidanceRef":     "protocol:area-synthesis-default-v0",
			"emptySignalPolicy":        "ignore_empty",
			"stopConditions":           []any{exampleLimit("missing-factor-analysis", "Stop if in-scope Factor analyses are unavailable.", "Area analysis becomes incomplete.")},
			"expectedEvaluationLimits": []any{exampleLimit("area-rollup-limit", "Only available Factor and child Area analyses are synthesized.", "Roll-up confidence is limited.")},
		},
	}
}

func scopedAnalysisExample(kind DataKind, idField string, id any) map[string]any {
	example := map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		idField:         id,
		"localAnalysis": map[string]any{
			"status":           "analyzed",
			"statusReason":     "Local inputs were synthesized.",
			"ratingLevelId":    RatingReference("target"),
			"rationale":        "Local signal satisfies the target criterion.",
			"inputRefs":        []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"ratingDrivers":    []any{map[string]any{"description": "Representative driver.", "effect": "supports target", "ratingLevelId": RatingReference("target"), "inputRefs": []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")}}},
			"incompleteInputs": []any{map[string]any{"id": "missing-descendant", "description": "One descendant input was unavailable.", "impact": "Roll-up confidence is lower."}},
			"evaluationLimits": []any{map[string]any{"id": "source-limit", "description": "Only local evidence was inspected.", "impact": "Rating may change with more sources."}},
			"confidence":       "medium",
			"confidenceReason": "Evidence is representative but partial.",
		},
		"localAndDescendantAnalysis": map[string]any{
			"status":           "analyzed",
			"statusReason":     "Local and descendant inputs were synthesized.",
			"ratingLevelId":    RatingReference("target"),
			"rationale":        "Available local and descendant signal satisfies the target criterion.",
			"inputRefs":        []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"ratingDrivers":    []any{map[string]any{"description": "Representative roll-up driver.", "effect": "supports target", "ratingLevelId": RatingReference("target"), "inputRefs": []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")}}},
			"incompleteInputs": []any{map[string]any{"id": "missing-child", "description": "One child input was unavailable.", "impact": "Roll-up is provisional."}},
			"evaluationLimits": []any{map[string]any{"id": "rollup-limit", "description": "The roll-up used available inputs only.", "impact": "Confidence is limited."}},
			"confidence":       "medium",
			"confidenceReason": "Roll-up uses available inputs only.",
		},
	}
	return example
}

func findingRankingExample() map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(DataKindFindingRanking),
		"subject":       map[string]any{"scopeRef": "run"},
		"orderedFindings": []any{
			map[string]any{
				"rank":       1,
				"findingRef": routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[gap-001]"),
				"tier":       "P1",
				"rationale":  "The gap most directly shapes the next quality-management move.",
			},
			map[string]any{
				"rank":       2,
				"findingRef": routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[strength-001]"),
				"tier":       "P3",
				"rationale":  "The strength is useful context but less urgent than the gap.",
			},
		},
		"rationale": "Findings are ordered by advice relevance, severity, rating influence, confidence, and model order.",
	}
}

func recommendationExample() map[string]any {
	return map[string]any{
		"schemaVersion":       SchemaVersion,
		"kind":                string(DataKindRecommendation),
		"id":                  "rec-001",
		"title":               "Review the next quality bar",
		"whyItMatters":        "The current evidence suggests the evaluated entity may already meet the present bar.",
		"recommendedNextMove": "Review whether the next evaluation should use sharper criteria.",
		"expectedBenefit":     "The Model stays useful as the evaluated entity improves.",
		"howToKnowItWorked":   "The review records either an updated bar or a rationale for keeping the current bar.",
		"impact":              "high",
		"confidence":          "medium",
		"traceRefs":           []any{routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[strength-001]")},
	}
}

func recommendationRankingExample() map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(DataKindRecommendationRanking),
		"orderedRecommendations": []any{map[string]any{
			"rank":              1,
			"recommendationRef": "rec-001",
			"impact":            "high",
			"confidence":        "medium",
			"rationale":         "This recommendation has the strongest expected quality-management impact.",
		}},
		"findingCoverage": []any{
			map[string]any{
				"findingRef":         routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[gap-001]"),
				"disposition":        "addressed_by_recommendation",
				"recommendationRefs": []any{"rec-001"},
			},
			map[string]any{
				"findingRef":  routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[strength-001]"),
				"disposition": "not_advice_driving",
				"rationale":   "The strength is context only and does not change the next quality-management move.",
			},
		},
		"rationale": "Recommendation ranking closes Advice after Finding coverage is accounted for.",
	}
}

func exampleAreaID() string {
	return model.AreaPath{}.Reference()
}

func exampleChildAreaID() string {
	return model.AreaPath{"operations"}.Reference()
}

func exampleRequirementID() string {
	return model.RequirementReference(nil, "has-tests")
}

func exampleFactorID() string {
	return model.FactorReference(nil, model.FactorPath{"verification"})
}

func exampleChildFactorID() string {
	return model.FactorReference(nil, model.FactorPath{"verification", "coverage"})
}

func exampleLimit(id, description, impact string) map[string]any {
	return map[string]any{"id": id, "description": description, "impact": impact}
}

func routineRef(kind DataKind, subject any, selector string) map[string]any {
	ref := map[string]any{"kind": string(kind), "subject": subject}
	if selector != "" {
		ref["selector"] = selector
	}
	return ref
}

func areaReportRef(areaID any, path string) map[string]any {
	return map[string]any{"kind": string(ReportKindArea), "areaId": areaID, "path": path}
}

func factorReportRef(areaID, factorID any, path string) map[string]any {
	return map[string]any{"kind": string(ReportKindFactor), "areaId": areaID, "factorId": factorID, "path": path}
}

func requirementReportRef(areaID, requirementID any, path string) map[string]any {
	return map[string]any{"kind": string(ReportKindRequirement), "areaId": areaID, "requirementId": requirementID, "path": path}
}

func runReportRef(path string) map[string]any {
	return map[string]any{"kind": string(ReportKindRun), "path": path}
}

func evaluationOutputExample() map[string]any {
	runReport := runReportRef("report.md")
	areaReport := areaReportRef(exampleAreaID(), "root-area.md")
	factorReport := factorReportRef(exampleAreaID(), exampleFactorID(), "factors/verification/verification-factor.md")
	requirementReport := requirementReportRef(exampleAreaID(), exampleRequirementID(), "requirements/has-tests/has-tests-requirement.md")
	return map[string]any{
		"schemaVersion":         SchemaVersion,
		"kind":                  string(DataKindEvaluationOutput),
		"runReportRef":          runReport,
		"scopedAreaAnalysisRef": routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleAreaID()}, "localAndDescendantAnalysis"),
		"rootAreaAnalysisRef":   routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleAreaID()}, "localAndDescendantAnalysis"),
		"areaOutputs": []any{map[string]any{
			"areaId":                    exampleAreaID(),
			"areaEvaluationFrameRef":    routineRef(DataKindAreaEvaluationFrame, map[string]any{"areaId": exampleAreaID()}, ""),
			"areaAnalysisResultRef":     routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleAreaID()}, ""),
			"factorAnalysisRefs":        []any{routineRef(DataKindFactorAnalysis, map[string]any{"factorId": exampleFactorID()}, "localAndDescendantAnalysis")},
			"requirementAssessmentRefs": []any{routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"requirementRatingRefs":     []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"reportRefs":                []any{areaReport, factorReport, requirementReport},
		}},
		"reportOutputs": []any{runReport, areaReport, factorReport, requirementReport},
	}
}
