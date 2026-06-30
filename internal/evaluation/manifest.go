package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const evaluationManifestPath = "data/evaluation-manifest.json"

func loadEvaluationManifest(runAbs string) (*EvaluationManifest, error) {
	raw, err := os.ReadFile(filepath.Join(runAbs, evaluationManifestPath))
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
