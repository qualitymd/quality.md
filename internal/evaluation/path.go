package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/qualitymd/quality.md/internal/workspace"
)

var (
	currentRunNameRE = regexp.MustCompile(`^(\d{4})-([a-z0-9-]+)-eval$`)
)

// ModelSnapshotFile is the run-folder filename for the frozen copy of the
// resolved working-tree model captured when a run is created. The name marks it
// as a point-in-time snapshot, distinct from the live working-tree model.
const ModelSnapshotFile = "model-snapshot.md"

// FindRepoRoot walks upward from start until it finds a Git repository root.
func FindRepoRoot(start string) (string, error) {
	return workspace.FindRepoRoot(start)
}

// ResolveRepoPath validates a repository-relative path and returns absolute and
// slash-normalized relative forms.
func ResolveRepoPath(repoRoot, value string) (string, string, error) {
	return workspace.ResolveRepoPath(repoRoot, value)
}

// ResolveDir resolves the configured evaluation directory from a repository
// root, returning both absolute and model-relative paths for the root model.
func ResolveDir(repoRoot, override string) (string, string, error) {
	ws, err := workspace.Resolve(workspace.Options{
		RepoRoot:              repoRoot,
		EvaluationDirOverride: override,
	})
	if err != nil {
		return "", "", err
	}
	return ws.Evaluations.Abs, ws.Evaluations.Rel, nil
}

// ResolveDirForModel resolves the configured evaluation directory for a
// selected model, returning absolute and model-relative paths.
func ResolveDirForModel(model, override string) (string, string, error) {
	ws, err := workspace.Resolve(workspace.Options{
		Model:                 model,
		EvaluationDirOverride: override,
	})
	if err != nil {
		return "", "", err
	}
	return ws.Evaluations.Abs, ws.Evaluations.Rel, nil
}

// RunDir is one recognized evaluation run folder.
type RunDir struct {
	Number int
	Name   string
	Abs    string
	Rel    string
}

type runName struct {
	number int
}

// ListRunDirs returns recognized evaluation run folders in deterministic order.
func ListRunDirs(evalDirAbs, evalDirRel string) ([]RunDir, error) {
	entries, err := os.ReadDir(evalDirAbs)
	if err != nil {
		return nil, err
	}
	var runs []RunDir
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		runAbs := filepath.Join(evalDirAbs, entry.Name())
		manifest, err := loadEvaluationManifest(runAbs)
		if err != nil {
			parsed, ok := parseRunName(entry.Name())
			if !ok {
				continue
			}
			runs = append(runs, RunDir{
				Number: parsed.number,
				Name:   entry.Name(),
				Abs:    runAbs,
				Rel:    filepath.ToSlash(filepath.Join(evalDirRel, entry.Name())),
			})
			continue
		}
		runs = append(runs, RunDir{
			Number: manifest.Run.Number,
			Name:   entry.Name(),
			Abs:    runAbs,
			Rel:    filepath.ToSlash(filepath.Join(evalDirRel, entry.Name())),
		})
	}
	slices.SortFunc(runs, func(a, b RunDir) int {
		if a.Number != b.Number {
			return a.Number - b.Number
		}
		return strings.Compare(a.Name, b.Name)
	})
	return runs, nil
}

// Slug normalizes a string into the path-safe slug form used for run names.
func Slug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	lastHyphen := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastHyphen = false
			continue
		}
		if !lastHyphen {
			b.WriteByte('-')
			lastHyphen = true
		}
	}
	return strings.Trim(b.String(), "-")
}

// IsPathSafeSlug reports whether s is already in path-safe slug form.
func IsPathSafeSlug(s string) bool {
	return s != "" && Slug(s) == s
}

func isCurrentRunScope(scope string) bool {
	if !IsPathSafeSlug(scope) {
		return false
	}
	for _, segment := range strings.Split(scope, "-") {
		if segment == "quality" {
			return false
		}
	}
	return true
}

func parseRunName(name string) (runName, bool) {
	if match := currentRunNameRE.FindStringSubmatch(name); match != nil {
		n, err := strconv.Atoi(match[1])
		if err != nil {
			return runName{}, false
		}
		scope := match[2]
		if !isCurrentRunScope(scope) {
			return runName{}, false
		}
		return runName{number: n}, true
	}
	return runName{}, false
}

func nextRunNumber(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	maxN := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		manifest, err := loadEvaluationManifest(filepath.Join(dir, entry.Name()))
		if err == nil && manifest.Run.Number > maxN {
			maxN = manifest.Run.Number
			continue
		}
		parsed, ok := parseRunName(entry.Name())
		if ok && parsed.number > maxN {
			maxN = parsed.number
		}
	}
	return maxN + 1, nil
}

func verifyRun(runPath string) (string, error) {
	abs, err := filepath.Abs(runPath)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("reading run %s: %w", runPath, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not an evaluation run folder", runPath)
	}
	if _, err := os.Stat(filepath.Join(abs, ModelSnapshotFile)); err != nil {
		return "", fmt.Errorf("%s is not an evaluation run folder: missing %s", runPath, ModelSnapshotFile)
	}
	// Runner-created runs carry one authoritative evaluation.json instead of
	// the historical multi-file data tree.
	if _, err := os.Stat(filepath.Join(abs, RunArtifactFile)); err == nil {
		return abs, nil
	}
	dataInfo, err := os.Stat(filepath.Join(abs, "data"))
	if err != nil {
		return "", fmt.Errorf("%s is not an evaluation run folder: missing data or %s", runPath, RunArtifactFile)
	}
	if !dataInfo.IsDir() {
		return "", fmt.Errorf("%s is not an evaluation run folder: data is not a directory", runPath)
	}
	return abs, nil
}

func displayRunPath(runAbs string) string {
	runAbs = filepath.Clean(runAbs)
	repoRoot, err := FindRepoRoot(runAbs)
	if err == nil {
		if rel, relErr := filepath.Rel(repoRoot, runAbs); relErr == nil && rel != "." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".." {
			return filepath.ToSlash(rel)
		}
	}
	return filepath.ToSlash(runAbs)
}

func resolveModelRelativeRun(model, runPath string) (string, string, error) {
	ws, err := workspace.Resolve(workspace.Options{Model: model})
	if err != nil {
		return "", "", err
	}
	if filepath.IsAbs(runPath) {
		return runPath, displayRunPath(runPath), nil
	}
	abs, rel, _, err := workspace.ResolveWorkspacePath(ws.RepoRoot.Abs, ws.WorkspaceRoot.Abs, runPath)
	if err != nil {
		return "", "", err
	}
	return abs, rel, nil
}
