package evaluation

import (
	"fmt"
	"path/filepath"
	"strings"
)

type EvaluationRunListEntry struct {
	Path       string                 `json:"path"`
	Subject    string                 `json:"subject"`
	Narrowing  string                 `json:"narrowing,omitempty"`
	Counts     EvaluationRecordCounts `json:"counts"`
	Reportable bool                   `json:"reportable"`
	Gaps       int                    `json:"gaps"`
}

type EvaluationRunList struct {
	SchemaVersion int                      `json:"schemaVersion"`
	Runs          []EvaluationRunListEntry `json:"runs"`
}

type EvaluationRecordList struct {
	SchemaVersion int                  `json:"schemaVersion"`
	Path          string               `json:"path"`
	Kind          EvaluationRecordKind `json:"kind"`
	Records       []string             `json:"records"`
}

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
	evalDirAbs, evalDirRel, err := EvaluationDir(repoRoot, evaluationDir)
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

func ListRuns(repoRoot, evaluationDir, state string) (*EvaluationRunList, error) {
	if repoRoot == "" {
		var err error
		repoRoot, err = FindRepoRoot("")
		if err != nil {
			return nil, err
		}
	}
	evalDirAbs, evalDirRel, err := EvaluationDir(repoRoot, evaluationDir)
	if err != nil {
		return nil, err
	}
	runs, err := ListRunDirs(evalDirAbs, evalDirRel)
	if err != nil {
		return nil, fmt.Errorf("listing evaluation runs: %w", err)
	}
	result := &EvaluationRunList{SchemaVersion: SchemaVersion}
	for _, dir := range runs {
		run, err := Inspect(dir.Abs)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", dir.Rel, err)
		}
		status := run.EvaluationRunStatus()
		if !includeRunState(status.Reportable, state) {
			continue
		}
		result.Runs = append(result.Runs, EvaluationRunListEntry{
			Path:       dir.Rel,
			Subject:    run.Model.Title,
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
	if match == nil || match[3] == "" {
		return ""
	}
	return strings.TrimPrefix(strings.TrimSuffix(match[3], "-quality-eval"), "-")
}

func ListRecords(kind EvaluationRecordKind, runPath string) (*EvaluationRecordList, error) {
	run, err := Inspect(runPath)
	if err != nil {
		return nil, err
	}
	result := &EvaluationRecordList{SchemaVersion: SchemaVersion, Path: filepath.ToSlash(run.Path), Kind: kind}
	records, err := recordFilesForKind(run.AbsPath, kind)
	if err != nil {
		return nil, usagef("unknown record kind %q", kind)
	}
	result.Records = records
	return result, nil
}

func recordFilesForKind(runAbs string, kind EvaluationRecordKind) ([]string, error) {
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
