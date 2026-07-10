package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const evaluationManifestPath = "data/evaluation-manifest.json"

// RunArtifactFile is the authoritative structured artifact runner-created
// runs keep at the run root instead of the multi-file data tree.
const RunArtifactFile = "evaluation.json"

func loadEvaluationManifest(runAbs string) (*EvaluationManifest, error) {
	raw, err := os.ReadFile(filepath.Join(runAbs, evaluationManifestPath))
	if os.IsNotExist(err) {
		return loadRunArtifactManifest(runAbs)
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", evaluationManifestPath, err)
	}
	payload, err := decodeDataPayload(raw)
	if err != nil {
		return nil, err
	}
	kind, err := payloadKind(payload)
	if err != nil {
		return nil, err
	}
	if kind != DataKindEvaluationManifest {
		return nil, usagef("%s kind = %s, want %s", evaluationManifestPath, kind, DataKindEvaluationManifest)
	}
	if err := validateDataPayload(kind, payload); err != nil {
		return nil, err
	}
	var manifest EvaluationManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", evaluationManifestPath, err)
	}
	return &manifest, nil
}

// loadRunArtifactManifest reads the manifest section of a runner-created
// run's evaluation.json.
func loadRunArtifactManifest(runAbs string) (*EvaluationManifest, error) {
	raw, err := os.ReadFile(filepath.Join(runAbs, RunArtifactFile))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", RunArtifactFile, err)
	}
	var doc struct {
		Manifest EvaluationManifest `json:"manifest"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", RunArtifactFile, err)
	}
	if doc.Manifest.EvaluationID == "" {
		return nil, fmt.Errorf("%s has no manifest", RunArtifactFile)
	}
	return &doc.Manifest, nil
}
