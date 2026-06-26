package evaluation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/workspace"
)

// CreateRun creates a numbered evaluation run folder and seeds its standard
// runtime files.
func CreateRun(opts Options) (*CreateRunReceipt, error) {
	if opts.Narrowing != "" && !IsPathSafeSlug(opts.Narrowing) {
		return nil, usagef("--narrowing must be a path-safe slug")
	}
	modelPath := opts.Model
	if modelPath == "" {
		modelPath = workspace.DefaultModelPath
	}
	if filepath.Clean(modelPath) == "." {
		return nil, usagef("--model %q must name a QUALITY.md file, not a directory", modelPath)
	}
	ws, err := workspace.Resolve(workspace.Options{
		RepoRoot:              opts.RepoRoot,
		Model:                 opts.Model,
		EvaluationDirOverride: opts.ResolveDir,
	})
	if err != nil {
		return nil, err
	}
	modelRaw, err := modelSnapshot(ws.Model.Abs, ws.Model.Rel)
	if err != nil {
		return nil, err
	}
	evalDirAbs := ws.Evaluations.Abs
	evalDirRel := ws.Evaluations.Rel
	if err := os.MkdirAll(evalDirAbs, 0o755); err != nil {
		return nil, fmt.Errorf("creating evaluation directory: %w", err)
	}
	number, name, err := nextRunName(evalDirAbs, opts.Narrowing)
	if err != nil {
		return nil, err
	}
	runAbs := filepath.Join(evalDirAbs, name)
	if err := os.Mkdir(runAbs, 0o755); err != nil {
		return nil, fmt.Errorf("creating run folder %s: %w", filepath.ToSlash(filepath.Join(evalDirRel, name)), err)
	}
	if err := createRunSkeleton(runAbs, modelRaw); err != nil {
		return nil, err
	}
	runRel := filepath.ToSlash(filepath.Join(evalDirRel, name))
	return &CreateRunReceipt{
		Path:   runRel,
		Number: number,
		NextActions: []receipt.Action{{
			ID:      "evaluation-data-set",
			Label:   "Record Evaluation v2 data",
			Command: "qualitymd evaluation data set " + runRel + " < payload.json",
		}},
	}, nil
}

func nextRunName(evalDirAbs, narrowing string) (int, string, error) {
	number, err := nextRunNumber(evalDirAbs)
	if err != nil {
		return 0, "", err
	}
	scope := "full"
	if narrowing != "" {
		scope = narrowing
	}
	return number, fmt.Sprintf("%04d-%s-eval", number, scope), nil
}

func createRunSkeleton(runAbs string, modelRaw []byte) error {
	if err := os.Mkdir(filepath.Join(runAbs, "data"), 0o755); err != nil {
		return fmt.Errorf("creating data: %w", err)
	}
	if err := os.WriteFile(filepath.Join(runAbs, ModelSnapshotFile), modelRaw, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", ModelSnapshotFile, err)
	}
	return nil
}

func modelSnapshot(modelAbs, modelRel string) ([]byte, error) {
	modelPath := modelRel
	info, err := os.Stat(modelAbs)
	if err != nil {
		return nil, fmt.Errorf("reading model %s: %w", modelPath, err)
	}
	if info.IsDir() {
		return nil, usagef("--model %q must name a QUALITY.md file, not a directory", modelPath)
	}
	raw, err := os.ReadFile(modelAbs)
	if err != nil {
		return nil, fmt.Errorf("reading model %s: %w", modelPath, err)
	}
	return raw, nil
}
