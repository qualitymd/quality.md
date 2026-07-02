package evaluation

import (
	"fmt"
	"path/filepath"
)

// RunListEntry summarizes one discovered evaluation run.
type RunListEntry struct {
	Path           string          `json:"path"`
	RootArea       string          `json:"rootArea"`
	RequestedScope RunScope        `json:"requestedScope"`
	PlannedScope   PlannedRunScope `json:"plannedScope"`
	DataArtifacts  int             `json:"dataArtifacts"`
	Reportable     bool            `json:"reportable"`
	Gaps           int             `json:"gaps"`
}

// RunList is the JSON contract emitted by evaluation run listings.
type RunList struct {
	SchemaVersion int            `json:"schemaVersion"`
	Runs          []RunListEntry `json:"runs"`
}

// RunSelection describes how a command chooses an evaluation run.
type RunSelection struct {
	Model         string
	EvaluationDir string
	RunArg        string
	Latest        bool
}

// ResolvedRun is an evaluation run path resolved for command execution.
type ResolvedRun struct {
	Path        string
	DisplayPath string
	Model       string
}

// ResolveRun resolves an explicit run path or the latest discovered run.
func ResolveRun(repoRoot, evaluationDir, runArg string, latest bool) (string, error) {
	result, err := ResolveRunSelection(RunSelection{
		Model:         modelFromRepoRoot(repoRoot),
		EvaluationDir: evaluationDir,
		RunArg:        runArg,
		Latest:        latest,
	})
	if err != nil {
		return "", err
	}
	return result.Path, nil
}

// ResolveRunSelection resolves an explicit run path or the latest discovered
// run for a selected model.
func ResolveRunSelection(selection RunSelection) (*ResolvedRun, error) {
	runArg := selection.RunArg
	latest := selection.Latest
	if runArg != "" && latest {
		return nil, usagef("pass a run path or --latest, not both")
	}
	if runArg == "" && !latest {
		return nil, usagef("pass a run path or --latest")
	}
	if runArg != "" {
		if selection.Model == "" {
			return &ResolvedRun{Path: runArg}, nil
		}
		path, displayPath, err := resolveModelRelativeRun(selection.Model, runArg)
		if err != nil {
			return nil, err
		}
		return &ResolvedRun{Path: path, DisplayPath: displayPath, Model: selection.Model}, nil
	}
	evalDirAbs, evalDirRel, err := ResolveDirForModel(selection.Model, selection.EvaluationDir)
	if err != nil {
		return nil, err
	}
	runs, err := ListRunDirs(evalDirAbs, evalDirRel)
	if err != nil {
		return nil, fmt.Errorf("listing evaluation runs: %w", err)
	}
	if len(runs) == 0 {
		return nil, usagef("no evaluation runs found")
	}
	run := runs[len(runs)-1]
	return &ResolvedRun{Path: run.Abs, DisplayPath: run.Rel, Model: selection.Model}, nil
}

// ListRuns lists discovered evaluation runs, optionally filtered by state.
func ListRuns(repoRoot, evaluationDir, state string) (*RunList, error) {
	return ListRunsForModel(modelFromRepoRoot(repoRoot), evaluationDir, state)
}

// ListRunsForModel lists discovered evaluation runs for a selected model,
// optionally filtered by state.
func ListRunsForModel(model, evaluationDir, state string) (*RunList, error) {
	evalDirAbs, evalDirRel, err := ResolveDirForModel(model, evaluationDir)
	if err != nil {
		return nil, err
	}
	runs, err := ListRunDirs(evalDirAbs, evalDirRel)
	if err != nil {
		return nil, fmt.Errorf("listing evaluation runs: %w", err)
	}
	result := &RunList{SchemaVersion: SchemaVersion}
	for _, dir := range runs {
		run, err := Inspect(dir.Abs)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", dir.Rel, err)
		}
		status := run.Status()
		if !includeRunState(status.Reportable, state) {
			continue
		}
		manifest, err := loadEvaluationManifest(dir.Abs)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", filepath.ToSlash(filepath.Join(dir.Rel, evaluationManifestPath)), err)
		}
		result.Runs = append(result.Runs, RunListEntry{
			Path:           dir.Rel,
			RootArea:       run.Model.Title,
			RequestedScope: manifest.RequestedScope,
			PlannedScope:   manifest.PlannedScope,
			DataArtifacts:  status.Data.Artifacts,
			Reportable:     status.Reportable,
			Gaps:           len(status.Gaps),
		})
	}
	return result, nil
}

func modelFromRepoRoot(repoRoot string) string {
	if repoRoot == "" {
		return ""
	}
	return filepath.Join(repoRoot, "QUALITY.md")
}

func includeRunState(reportable bool, state string) bool {
	switch state {
	case "", "all":
		return true
	case "reportable", "complete":
		return reportable
	case "incomplete":
		return !reportable
	default:
		return false
	}
}

// ValidateRunState checks an evaluation run listing state filter.
func ValidateRunState(state string) error {
	switch state {
	case "", "all", "reportable", "complete", "incomplete":
		return nil
	default:
		return usagef("--state must be one of: all, complete, reportable, incomplete")
	}
}
