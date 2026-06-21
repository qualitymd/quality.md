package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	runNameRE    = regexp.MustCompile(`^(\d{4})(?:-((?:subject|model)(?:-[a-z0-9-]+)?|[a-z0-9-]+))?-quality-eval$`)
	recordNameRE = regexp.MustCompile(`^(\d+)-`)
)

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
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		abs = filepath.Dir(abs)
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
	if filepath.IsAbs(value) {
		return "", "", fmt.Errorf("path %q must be repository-relative", value)
	}
	clean := filepath.Clean(value)
	if clean == "." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || clean == ".." {
		return "", "", fmt.Errorf("path %q escapes the repository", value)
	}
	abs := filepath.Join(repoRoot, clean)
	rel, err := filepath.Rel(repoRoot, abs)
	if err != nil {
		return "", "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("path %q escapes the repository", value)
	}
	return abs, filepath.ToSlash(rel), nil
}

// EvaluationDir resolves the configured evaluation directory from a repository
// root, returning both absolute and repository-relative paths.
func EvaluationDir(repoRoot, override string) (string, string, error) {
	value, err := evaluationDirValue(repoRoot, override)
	if err != nil {
		return "", "", err
	}
	return ResolveRepoPath(repoRoot, value)
}

// RunDir is one recognized evaluation run folder.
type RunDir struct {
	Number int
	Name   string
	Abs    string
	Rel    string
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
		match := runNameRE.FindStringSubmatch(entry.Name())
		if match == nil {
			continue
		}
		n, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		runs = append(runs, RunDir{
			Number: n,
			Name:   entry.Name(),
			Abs:    filepath.Join(evalDirAbs, entry.Name()),
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

func nextRunNumber(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	maxN := 0
	for _, entry := range entries {
		m := runNameRE.FindStringSubmatch(entry.Name())
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxN {
			maxN = n
		}
	}
	return maxN + 1, nil
}

func nextRecordNumber(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	maxN := 0
	for _, entry := range entries {
		m := recordNameRE.FindStringSubmatch(entry.Name())
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err == nil && n > maxN {
			maxN = n
		}
	}
	return maxN + 1, nil
}
