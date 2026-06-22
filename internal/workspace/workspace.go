// Package workspace resolves the project-local paths used by qualitymd
// tooling for one selected QUALITY.md file.
package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/qualitymd/quality.md/internal/document"
	"gopkg.in/yaml.v3"
)

const (
	DefaultModelPath       = "QUALITY.md"
	DefaultConfigPath      = ".quality/config.yaml"
	DefaultDataDir         = ".quality"
	DefaultEvaluationDir   = ".quality/evaluations"
	DefaultQualityLogDir   = ".quality/log"
	FrontmatterConfigField = "config"
)

// PathRef carries both absolute and repository-relative path forms.
type PathRef struct {
	Abs string
	Rel string
}

// Workspace is the resolved operating context for one QUALITY.md file.
type Workspace struct {
	RepoRoot PathRef
	Model    PathRef

	Config        PathRef
	ConfigPresent bool

	DataDir     PathRef
	Evaluations PathRef
	Log         PathRef
}

// Options controls workspace resolution.
type Options struct {
	RepoRoot              string
	Model                 string
	EvaluationDirOverride string
}

type config struct {
	EvaluationDir string `yaml:"evaluationDir"`
}

// Resolve returns the workspace paths for the selected model file.
func Resolve(opts Options) (*Workspace, error) {
	modelInput := opts.Model
	if modelInput == "" {
		modelInput = DefaultModelPath
	}

	repoRoot, modelAbs, err := resolveRoots(opts.RepoRoot, modelInput)
	if err != nil {
		return nil, err
	}
	modelRel, err := relToRepo(repoRoot, modelAbs)
	if err != nil {
		return nil, fmt.Errorf("resolving model path: %w", err)
	}

	configValue, err := configPathFromModel(modelAbs)
	if err != nil {
		return nil, err
	}
	if configValue == "" {
		configValue = DefaultConfigPath
	}
	configAbs, configRel, err := ResolveRepoPath(repoRoot, configValue)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg, configPresent, err := readConfig(configAbs, configRel)
	if err != nil {
		return nil, err
	}
	evaluationDir := DefaultEvaluationDir
	if cfg.EvaluationDir != "" {
		evaluationDir = cfg.EvaluationDir
	}
	if opts.EvaluationDirOverride != "" {
		evaluationDir = opts.EvaluationDirOverride
	}

	dataAbs, dataRel, err := ResolveRepoPath(repoRoot, DefaultDataDir)
	if err != nil {
		return nil, err
	}
	evalAbs, evalRel, err := ResolveRepoPath(repoRoot, evaluationDir)
	if err != nil {
		return nil, fmt.Errorf("evaluationDir: %w", err)
	}
	logAbs, logRel, err := ResolveRepoPath(repoRoot, DefaultQualityLogDir)
	if err != nil {
		return nil, err
	}

	return &Workspace{
		RepoRoot: PathRef{Abs: repoRoot, Rel: "."},
		Model:    PathRef{Abs: modelAbs, Rel: modelRel},
		Config: PathRef{
			Abs: configAbs,
			Rel: configRel,
		},
		ConfigPresent: configPresent,
		DataDir:       PathRef{Abs: dataAbs, Rel: dataRel},
		Evaluations:   PathRef{Abs: evalAbs, Rel: evalRel},
		Log:           PathRef{Abs: logAbs, Rel: logRel},
	}, nil
}

func resolveRoots(repoRootInput, modelInput string) (string, string, error) {
	if repoRootInput != "" {
		repoRoot, err := filepath.Abs(repoRootInput)
		if err != nil {
			return "", "", err
		}
		modelAbs := modelInput
		if !filepath.IsAbs(modelAbs) {
			modelAbs = filepath.Join(repoRoot, modelInput)
		}
		modelAbs, err = filepath.Abs(modelAbs)
		if err != nil {
			return "", "", err
		}
		return repoRoot, modelAbs, nil
	}

	modelAbs := modelInput
	if !filepath.IsAbs(modelAbs) {
		var err error
		modelAbs, err = filepath.Abs(modelInput)
		if err != nil {
			return "", "", err
		}
	}
	repoRoot, err := FindRepoRoot(modelAbs)
	if err != nil {
		return "", "", err
	}
	return repoRoot, modelAbs, nil
}

func configPathFromModel(modelAbs string) (string, error) {
	doc, err := document.Parse(modelAbs)
	if err != nil {
		return "", err
	}
	_, value, _ := document.MapEntry(doc.Frontmatter, FrontmatterConfigField)
	if value == nil {
		return "", nil
	}
	if value.Kind != yaml.ScalarNode || strings.TrimSpace(value.Value) == "" {
		return "", fmt.Errorf("root config must be a non-empty repository-relative scalar path")
	}
	_, err = CleanRepoRelative(value.Value)
	if err != nil {
		return "", err
	}
	return value.Value, nil
}

func readConfig(configAbs, configRel string) (config, bool, error) {
	raw, err := os.ReadFile(configAbs)
	if err != nil {
		if os.IsNotExist(err) {
			return config{}, false, nil
		}
		return config{}, false, fmt.Errorf("reading %s: %w", configRel, err)
	}
	var cfg config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return config{}, true, fmt.Errorf("parsing %s: %w", configRel, err)
	}
	return cfg, true, nil
}

// FindRepoRoot walks upward from start until it finds a Git repository root.
func FindRepoRoot(start string) (string, error) {
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err == nil && !info.IsDir() {
		abs = filepath.Dir(abs)
	} else if os.IsNotExist(err) {
		abs = filepath.Dir(abs)
	} else if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(abs, ".git")); err == nil {
			return abs, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", fmt.Errorf("could not find repository root from %s", start)
		}
		abs = parent
	}
}

// ResolveRepoPath validates a repository-relative path and returns absolute and
// slash-normalized relative forms.
func ResolveRepoPath(repoRoot, value string) (string, string, error) {
	rel, err := CleanRepoRelative(value)
	if err != nil {
		return "", "", err
	}
	abs := filepath.Join(repoRoot, filepath.FromSlash(rel))
	checkedRel, err := relToRepo(repoRoot, abs)
	if err != nil {
		return "", "", err
	}
	return abs, checkedRel, nil
}

// CleanRepoRelative validates and normalizes a repository-relative path without
// requiring the repository root.
func CleanRepoRelative(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("path must be non-empty")
	}
	if filepath.IsAbs(value) {
		return "", fmt.Errorf("path %q must be repository-relative", value)
	}
	clean := filepath.Clean(filepath.FromSlash(value))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes the repository", value)
	}
	return filepath.ToSlash(clean), nil
}

func relToRepo(repoRoot, abs string) (string, error) {
	rel, err := filepath.Rel(repoRoot, abs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path %q escapes the repository", abs)
	}
	return filepath.ToSlash(rel), nil
}
