package evaluation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	runNameRE    = regexp.MustCompile(`^(\d{4})-(subject|model)(-[a-z0-9-]+)?-quality-eval$`)
	recordNameRE = regexp.MustCompile(`^(\d+)-`)
)

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
