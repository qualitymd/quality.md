package evaluation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/qualitymd/quality.md/internal/models"
	"github.com/qualitymd/quality.md/internal/receipt"
	"gopkg.in/yaml.v3"
)

type config struct {
	EvaluationDir string `yaml:"evaluationDir"`
}

func CreateRun(opts Options) (*CreateRunResult, error) {
	if opts.Altitude != "subject" && opts.Altitude != "model" {
		return nil, usagef("--altitude must be subject or model")
	}
	if opts.Narrowing != "" && !IsPathSafeSlug(opts.Narrowing) {
		return nil, usagef("--narrowing must be a path-safe slug")
	}
	repoRoot := opts.RepoRoot
	var err error
	if repoRoot == "" {
		repoRoot, err = FindRepoRoot("")
		if err != nil {
			return nil, err
		}
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
	number, err := nextRunNumber(evalDirAbs)
	if err != nil {
		return nil, err
	}
	name := fmt.Sprintf("%04d-%s-quality-eval", number, opts.Altitude)
	if opts.Narrowing != "" {
		name = fmt.Sprintf("%04d-%s-%s-quality-eval", number, opts.Altitude, opts.Narrowing)
	}
	runAbs := filepath.Join(evalDirAbs, name)
	if err := os.Mkdir(runAbs, 0o755); err != nil {
		return nil, fmt.Errorf("creating run folder %s: %w", filepath.ToSlash(filepath.Join(evalDirRel, name)), err)
	}
	for _, subdir := range []string{"assessments", "analysis", "recommendations"} {
		if err := os.Mkdir(filepath.Join(runAbs, subdir), 0o755); err != nil {
			return nil, fmt.Errorf("creating %s: %w", subdir, err)
		}
	}
	if err := os.WriteFile(filepath.Join(runAbs, "model.md"), modelRaw, 0o644); err != nil {
		return nil, fmt.Errorf("writing model.md: %w", err)
	}
	for file, content := range map[string]string{
		"design.md": "# Evaluation design\n",
		"plan.md":   "# Evaluation plan\n",
	} {
		if err := os.WriteFile(filepath.Join(runAbs, file), []byte(content), 0o644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", file, err)
		}
	}
	runRel := filepath.ToSlash(filepath.Join(evalDirRel, name))
	return &CreateRunResult{
		Path:     runRel,
		Number:   number,
		Altitude: opts.Altitude,
		NextActions: []receipt.Action{{
			ID:      "add-record",
			Label:   "Record evaluation judgments",
			Command: "qualitymd evaluation add-record assessment " + runRel,
		}},
	}, nil
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
	subject := opts.Subject
	if subject == "" {
		subject = "QUALITY.md"
	}
	if filepath.Clean(subject) == "." {
		return nil, usagef("--subject %q must name a QUALITY.md file, not a directory", subject)
	}
	subjectAbs, subjectRel, err := ResolveRepoPath(repoRoot, subject)
	if err != nil {
		return nil, usagef("--subject %s", err)
	}
	info, err := os.Stat(subjectAbs)
	if err != nil {
		return nil, fmt.Errorf("reading subject %s: %w", subject, err)
	}
	if info.IsDir() {
		return nil, usagef("--subject %q must name a QUALITY.md file, not a directory", subject)
	}
	if opts.Altitude == "model" {
		raw, err := models.Markdown("quality-meta-model", subjectRel)
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
	raw, err := os.ReadFile(subjectAbs)
	if err != nil {
		return nil, fmt.Errorf("reading subject %s: %w", subject, err)
	}
	return raw, nil
}
