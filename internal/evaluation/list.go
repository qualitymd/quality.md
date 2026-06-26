package evaluation

import "fmt"

// RunListEntry summarizes one discovered evaluation run.
type RunListEntry struct {
	Path          string `json:"path"`
	RootArea      string `json:"rootArea"`
	Narrowing     string `json:"narrowing,omitempty"`
	DataArtifacts int    `json:"dataArtifacts"`
	Reportable    bool   `json:"reportable"`
	Gaps          int    `json:"gaps"`
}

// RunList is the JSON contract emitted by evaluation run listings.
type RunList struct {
	SchemaVersion int            `json:"schemaVersion"`
	Runs          []RunListEntry `json:"runs"`
}

// ResolveRun resolves an explicit run path or the latest discovered run.
func ResolveRun(repoRoot, evaluationDir, runArg string, latest bool) (string, error) {
	if runArg != "" && latest {
		return "", usagef("pass a run path or --latest, not both")
	}
	if runArg == "" && !latest {
		return "", usagef("pass a run path or --latest")
	}
	if runArg != "" {
		return runArg, nil
	}
	if repoRoot == "" {
		var err error
		repoRoot, err = FindRepoRoot("")
		if err != nil {
			return "", err
		}
	}
	evalDirAbs, evalDirRel, err := ResolveDir(repoRoot, evaluationDir)
	if err != nil {
		return "", err
	}
	runs, err := ListRunDirs(evalDirAbs, evalDirRel)
	if err != nil {
		return "", fmt.Errorf("listing evaluation runs: %w", err)
	}
	if len(runs) == 0 {
		return "", usagef("no evaluation runs found")
	}
	return runs[len(runs)-1].Rel, nil
}

// ListRuns lists discovered evaluation runs, optionally filtered by state.
func ListRuns(repoRoot, evaluationDir, state string) (*RunList, error) {
	if repoRoot == "" {
		var err error
		repoRoot, err = FindRepoRoot("")
		if err != nil {
			return nil, err
		}
	}
	evalDirAbs, evalDirRel, err := ResolveDir(repoRoot, evaluationDir)
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
		result.Runs = append(result.Runs, RunListEntry{
			Path:          dir.Rel,
			RootArea:      run.Model.Title,
			Narrowing:     narrowingFromRunName(dir.Name),
			DataArtifacts: status.Data.Artifacts,
			Reportable:    status.Reportable,
			Gaps:          len(status.Gaps),
		})
	}
	return result, nil
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

func narrowingFromRunName(name string) string {
	parsed, ok := parseRunName(name)
	if !ok {
		return ""
	}
	return parsed.narrowing
}
