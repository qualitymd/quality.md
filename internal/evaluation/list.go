package evaluation

import (
	"fmt"
	"path/filepath"
	"strings"
)

// RunListEntry summarizes one discovered evaluation run.
type RunListEntry struct {
	Path       string       `json:"path"`
	RootArea   string       `json:"rootArea"`
	Narrowing  string       `json:"narrowing,omitempty"`
	Counts     RecordCounts `json:"counts"`
	Reportable bool         `json:"reportable"`
	Gaps       int          `json:"gaps"`
}

// RunList is the JSON contract emitted by evaluation run listings.
type RunList struct {
	SchemaVersion int            `json:"schemaVersion"`
	Runs          []RunListEntry `json:"runs"`
}

// RecordList is the JSON contract emitted by record listings.
type RecordList struct {
	SchemaVersion int        `json:"schemaVersion"`
	Path          string     `json:"path"`
	Kind          RecordKind `json:"kind"`
	Records       []string   `json:"records"`
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
			Path:       dir.Rel,
			RootArea:   run.Model.Title,
			Narrowing:  narrowingFromRunName(dir.Name),
			Counts:     status.Counts,
			Reportable: status.Reportable,
			Gaps:       len(status.Gaps),
		})
	}
	return result, nil
}

func includeRunState(reportable bool, state string) bool {
	switch state {
	case "", "all":
		return true
	case "reportable":
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
	case "", "all", "reportable", "incomplete":
		return nil
	default:
		return usagef("--state must be one of: all, reportable, incomplete")
	}
}

func narrowingFromRunName(name string) string {
	match := runNameRE.FindStringSubmatch(name)
	if match == nil || match[2] == "" {
		return ""
	}
	suffix := match[2]
	if suffix == "subject" || suffix == "model" {
		return ""
	}
	for _, prefix := range []string{"subject-", "model-"} {
		if strings.HasPrefix(suffix, prefix) {
			return strings.TrimPrefix(suffix, prefix)
		}
	}
	return suffix
}

// ListRecords lists record files of the requested kind for a run.
func ListRecords(kind RecordKind, runPath string) (*RecordList, error) {
	run, err := Inspect(runPath)
	if err != nil {
		return nil, err
	}
	result := &RecordList{SchemaVersion: SchemaVersion, Path: filepath.ToSlash(run.Path), Kind: kind}
	records, err := recordFilesForKind(run.AbsPath, kind)
	if err != nil {
		return nil, usagef("unknown record kind %q", kind)
	}
	result.Records = records
	return result, nil
}

func recordFilesForKind(runAbs string, kind RecordKind) ([]string, error) {
	var dir, pattern string
	switch kind {
	case KindAssessmentResult:
		dir, pattern = "assessments", "*.json"
	case KindAnalysis:
		dir, pattern = "analysis", "*.json"
	case KindRecommendation:
		dir, pattern = "recommendations", "*.md"
	default:
		return nil, usagef("unknown record kind %q", kind)
	}
	var records []string
	for _, path := range globRecordFiles(filepath.Join(runAbs, dir), pattern) {
		records = append(records, filepath.ToSlash(filepath.Join(dir, filepath.Base(path))))
	}
	return records, nil
}
