package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"github.com/qualitymd/quality.md/internal/model"
	"github.com/qualitymd/quality.md/internal/receipt"
	"github.com/qualitymd/quality.md/internal/workspace"
)

// CreateRun creates a numbered evaluation run folder and seeds its standard
// runtime files.
func CreateRun(opts Options) (*CreateRunReceipt, error) {
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
	spec, err := loadModel(ws.Model.Abs)
	if err != nil {
		return nil, err
	}
	scope, err := resolveCreateScope(spec, opts)
	if err != nil {
		return nil, err
	}
	scope, err = completeRunIdentity(scope)
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
	number, name, err := nextRunName(evalDirAbs, scope.PlannedScope)
	if err != nil {
		return nil, err
	}
	scope.Number = number
	scope.Model = ws.Model.Rel
	runAbs := filepath.Join(evalDirAbs, name)
	if err := os.Mkdir(runAbs, 0o755); err != nil {
		return nil, fmt.Errorf("creating run folder %s: %w", filepath.ToSlash(filepath.Join(evalDirRel, name)), err)
	}
	if err := createRunSkeleton(runAbs, modelRaw, scope); err != nil {
		return nil, err
	}
	runRel := filepath.ToSlash(filepath.Join(evalDirRel, name))
	modelArgs := modelFlagArgs(opts.Model)
	return &CreateRunReceipt{
		Path:   runRel,
		Number: number,
		NextActions: []receipt.Action{{
			ID:      "evaluation-data-set",
			Label:   "Record Evaluation data",
			Command: "qualitymd evaluation data set" + modelArgs + " " + runRel + " < payloads.json",
		}},
	}, nil
}

func loadModel(path string) (*model.Spec, error) {
	doc, err := document.Parse(path)
	if err != nil {
		return nil, err
	}
	return model.Decode(doc)
}

func resolveCreateScope(spec *model.Spec, opts Options) (RunManifest, error) {
	requested := RunScope{}
	planned := PlannedRunScope{AreaID: model.AreaPath{}.Reference(), FactorFilter: []string{}}

	if opts.Area != "" {
		area, err := model.ParseAreaReference(spec, opts.Area)
		if err != nil {
			return RunManifest{}, usagef("--area: %v", err)
		}
		requested.AreaID = area.Reference()
		planned.AreaID = area.Reference()
	}

	for _, value := range opts.Factors {
		area, factor, err := model.ParseFactorReference(spec, value)
		if err != nil {
			return RunManifest{}, usagef("--factor: %v", err)
		}
		factorRef := model.FactorReference(area, factor)
		if requested.AreaID == "" {
			requested.AreaID = area.Reference()
		}
		if planned.AreaID == (model.AreaPath{}).Reference() && opts.Area == "" && len(planned.FactorFilter) == 0 {
			planned.AreaID = area.Reference()
		}
		if requested.AreaID != area.Reference() || planned.AreaID != area.Reference() {
			return RunManifest{}, usagef("--factor %s does not belong to --area %s", value, planned.AreaID)
		}
		requested.FactorFilter = append(requested.FactorFilter, factorRef)
		planned.FactorFilter = append(planned.FactorFilter, factorRef)
	}

	return RunManifest{
		SchemaVersion:  SchemaVersion,
		Kind:           DataKindRunManifest,
		RequestedScope: requested,
		PlannedScope:   planned,
	}, nil
}

func completeRunIdentity(manifest RunManifest) (RunManifest, error) {
	id, createdAt, err := newRunIdentity()
	if err != nil {
		return RunManifest{}, err
	}
	manifest.ID = id
	manifest.CreatedAt = createdAt
	return manifest, nil
}

func modelFlagArgs(model string) string {
	if model == "" || model == workspace.DefaultModelPath {
		return ""
	}
	return " --model " + model
}

func nextRunName(evalDirAbs string, scope PlannedRunScope) (int, string, error) {
	number, err := nextRunNumber(evalDirAbs)
	if err != nil {
		return 0, "", err
	}
	return number, fmt.Sprintf("%04d-%s-eval", number, runScopeSlug(scope)), nil
}

func runScopeSlug(scope PlannedRunScope) string {
	if scope.AreaID == (model.AreaPath{}).Reference() && len(scope.FactorFilter) == 0 {
		return "full"
	}
	var parts []string
	area, err := areaIDFrom(scope.AreaID)
	if err == nil && len(area) > 0 {
		parts = append(parts, area...)
	} else {
		parts = append(parts, "root")
	}
	for _, ref := range scope.FactorFilter {
		id, err := factorIDFrom(ref)
		if err != nil {
			continue
		}
		parts = append(parts, id.Path...)
	}
	return strings.Join(parts, "-")
}

func createRunSkeleton(runAbs string, modelRaw []byte, manifest RunManifest) error {
	if err := os.Mkdir(filepath.Join(runAbs, "data"), 0o755); err != nil {
		return fmt.Errorf("creating data: %w", err)
	}
	if err := os.WriteFile(filepath.Join(runAbs, ModelSnapshotFile), modelRaw, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", ModelSnapshotFile, err)
	}
	raw, err := canonicalJSON(manifest)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(runAbs, "data", "run-manifest.json"), raw, 0o644); err != nil {
		return fmt.Errorf("writing data/run-manifest.json: %w", err)
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
