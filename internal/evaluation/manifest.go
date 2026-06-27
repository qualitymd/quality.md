package evaluation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const runManifestPath = "data/run-manifest.json"

func loadRunManifest(runAbs string) (*RunManifest, error) {
	raw, err := os.ReadFile(filepath.Join(runAbs, runManifestPath))
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", runManifestPath, err)
	}
	payload, err := decodeDataPayload(raw)
	if err != nil {
		return nil, err
	}
	kind, err := payloadKind(payload)
	if err != nil {
		return nil, err
	}
	if kind != DataKindRunManifest {
		return nil, usagef("%s kind = %s, want %s", runManifestPath, kind, DataKindRunManifest)
	}
	if err := validateDataPayload(kind, payload); err != nil {
		return nil, err
	}
	var manifest RunManifest
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", runManifestPath, err)
	}
	return &manifest, nil
}
