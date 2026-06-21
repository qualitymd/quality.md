package evaluation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

type config struct {
	EvaluationDir string `yaml:"evaluationDir"`
}

const debugLogSeed = `# Evaluation debug log

This log records notable events involving the evaluation process itself. It is
not an assessment record, rating rationale, report, or evidence store;
evaluation evidence belongs in assessment, analysis, and recommendation records.

## Events
`

func CreateRun(opts Options) (*CreateRunReceipt, error) {
	if opts.Narrowing != "" && !IsPathSafeSlug(opts.Narrowing) {
		return nil, usagef("--narrowing must be a path-safe slug")
	}
	repoRoot, err := repoRootForCreate(opts)
	if err != nil {
		return nil, err
	}
	modelRaw, err := modelSnapshot(repoRoot, opts)
	if err != nil {
		return nil, err
	}
	evalDirValue, err := evaluationDirValue(repoRoot, opts.EvaluationDir)
	if err != nil {
		return nil, err
	}
	evalDirAbs, evalDirRel, err := ResolveRepoPath(repoRoot, evalDirValue)
	if err != nil {
		return nil, err
	}
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
			ID:      "assessment-add",
			Label:   "Record evaluation judgments",
			Command: "qualitymd evaluation assessment add " + runRel,
		}},
	}, nil
}

func repoRootForCreate(opts Options) (string, error) {
	if opts.RepoRoot != "" {
		return opts.RepoRoot, nil
	}
	return FindRepoRoot("")
}

func nextRunName(evalDirAbs, narrowing string) (int, string, error) {
	number, err := nextRunNumber(evalDirAbs)
	if err != nil {
		return 0, "", err
	}
	name := fmt.Sprintf("%04d-quality-eval", number)
	if narrowing != "" {
		name = fmt.Sprintf("%04d-%s-quality-eval", number, narrowing)
	}
	return number, name, nil
}

func createRunSkeleton(runAbs string, modelRaw []byte) error {
	for _, subdir := range []string{"assessments", "analysis", "recommendations"} {
		if err := os.Mkdir(filepath.Join(runAbs, subdir), 0o755); err != nil {
			return fmt.Errorf("creating %s: %w", subdir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(runAbs, "model.md"), modelRaw, 0o644); err != nil {
		return fmt.Errorf("writing model.md: %w", err)
	}
	for file, content := range map[string]string{
		"debug-log.md": debugLogSeed,
		"design.md":    "# Evaluation design\n",
		"plan.md":      "# Evaluation plan\n",
	} {
		if err := os.WriteFile(filepath.Join(runAbs, file), []byte(content), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", file, err)
		}
	}
	return nil
}

func evaluationDirValue(repoRoot, override string) (string, error) {
	if override != "" {
		return override, nil
	}
	raw, err := os.ReadFile(filepath.Join(repoRoot, ".quality", "config.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return "quality/evaluations", nil
		}
		return "", fmt.Errorf("reading .quality/config.yaml: %w", err)
	}
	var cfg config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return "", fmt.Errorf("parsing .quality/config.yaml: %w", err)
	}
	if cfg.EvaluationDir == "" {
		return "quality/evaluations", nil
	}
	return cfg.EvaluationDir, nil
}

func modelSnapshot(repoRoot string, opts Options) ([]byte, error) {
	modelPath := opts.Model
	if modelPath == "" {
		modelPath = "QUALITY.md"
	}
	if filepath.Clean(modelPath) == "." {
		return nil, usagef("--model %q must name a QUALITY.md file, not a directory", modelPath)
	}
	modelAbs, _, err := ResolveRepoPath(repoRoot, modelPath)
	if err != nil {
		return nil, usagef("--model %s", err)
	}
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
