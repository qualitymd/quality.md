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
	DataKindEvaluationFrame            DataKind = "EvaluationFrame"
	DataKindAreaEvaluationFrame        DataKind = "AreaEvaluationFrame"
	DataKindRequirementEvaluationFrame DataKind = "RequirementEvaluationFrame"
	DataKindRequirementAssessment      DataKind = "RequirementAssessmentResult"
	DataKindRequirementRating          DataKind = "RequirementRatingResult"
	DataKindFactorAnalysisFrame        DataKind = "FactorAnalysisFrame"
	DataKindFactorAnalysis             DataKind = "FactorAnalysisResult"
	DataKindAreaAnalysisFrame          DataKind = "AreaAnalysisFrame"
	DataKindAreaAnalysis               DataKind = "AreaAnalysisResult"
	DataKindEvaluationOutput           DataKind = "EvaluationOutputResult"
)

// supportedDataKinds lists every Evaluation payload kind the CLI can persist,
// including the CLI-owned EvaluationOutputResult. It is the single typed source
// for the reference-kind vocabulary: any reference may name any of these kinds.
var supportedDataKinds = []DataKind{
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

// acceptedDataKinds is the agent-writable subset of supportedDataKinds: every
// supported kind except the CLI-owned EvaluationOutputResult, which only
// evaluation report build generates. It gates data set payloads.
var acceptedDataKinds = supportedDataKindsExcept(DataKindEvaluationOutput)

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
	if kind == DataKindEvaluationOutput {
		return dataWriteCandidate{}, usagef("%s is generated by evaluation report build and cannot be written with evaluation data set", kind)
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

func dataPathForPayload(kind DataKind, payload map[string]any) (string, error) {
	switch kind {
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
	default:
		return "", usagef("unsupported evaluation data kind %q", kind)
	}
}

func dataPathForQuery(query DataQuery) (string, error) {
	switch query.Kind {
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

func evaluationFrameExample(kind DataKind) map[string]any {
	return map[string]any{
		"schemaVersion": SchemaVersion,
		"kind":          string(kind),
		"subject":       map[string]any{"modelLocator": "QUALITY.md"},
		"inputs": map[string]any{
			"requestedScope": "full evaluation",
			"ratingLevelIds": []any{RatingReference("target"), RatingReference("unacceptable")},
			"areaIds":        []any{exampleAreaID()},
			"factorIds":      []any{exampleFactorID()},
		},
		"derivedContext": map[string]any{
			"resolvedScope":            "full evaluation",
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
				"id":          "tests-present",
				"type":        "strength",
				"severity":    "medium",
				"description": "A focused test covers the requirement.",
				"evidence":    map[string]any{"sourceRef": "tests/example_test.go"},
				"rationale":   "The evidence directly addresses the target.",
			},
			map[string]any{
				"id":          "edge-cases-untested",
				"type":        "gap",
				"severity":    "medium",
				"description": "Edge-case paths around the requirement lack tests.",
				"evidence":    map[string]any{"sourceRef": "tests/example_test.go"},
				"rationale":   "Only the happy path is exercised.",
				"actions": []any{map[string]any{
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
	if kind == DataKindAreaAnalysis {
		example["findings"] = []any{map[string]any{
			"id":         "area-finding-1",
			"type":       "gap",
			"severity":   "medium",
			"confidence": "medium",
			"summary":    "Area-level synthesis found a representative gap.",
			"rationale":  "The gap is synthesized from Requirement and Factor results.",
			"inputRefs":  []any{routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "findings[gap-1]")},
			"factorRelationships": []any{map[string]any{
				"factorId":     exampleFactorID(),
				"relationship": "primary-driver",
				"rationale":    "The finding directly affects the Factor analysis.",
			}},
		}}
	}
	return example
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

func evaluationOutputExample() map[string]any {
	areaReport := areaReportRef(exampleAreaID(), "report.md")
	factorReport := factorReportRef(exampleAreaID(), exampleFactorID(), "factors/verification/verification-factor.md")
	requirementReport := requirementReportRef(exampleAreaID(), exampleRequirementID(), "requirements/has-tests/has-tests-requirement.md")
	return map[string]any{
		"schemaVersion":       SchemaVersion,
		"kind":                string(DataKindEvaluationOutput),
		"rootAreaAnalysisRef": routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleAreaID()}, "localAndDescendantAnalysis"),
		"areaOutputs": []any{map[string]any{
			"areaId":                    exampleAreaID(),
			"areaEvaluationFrameRef":    routineRef(DataKindAreaEvaluationFrame, map[string]any{"areaId": exampleAreaID()}, ""),
			"areaAnalysisResultRef":     routineRef(DataKindAreaAnalysis, map[string]any{"areaId": exampleAreaID()}, ""),
			"factorAnalysisRefs":        []any{routineRef(DataKindFactorAnalysis, map[string]any{"factorId": exampleFactorID()}, "localAndDescendantAnalysis")},
			"requirementAssessmentRefs": []any{routineRef(DataKindRequirementAssessment, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"requirementRatingRefs":     []any{routineRef(DataKindRequirementRating, map[string]any{"requirementId": exampleRequirementID()}, "")},
			"reportRefs":                []any{areaReport, factorReport, requirementReport},
		}},
		"reportOutputs": []any{areaReport, factorReport, requirementReport},
	}
}
