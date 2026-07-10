package runner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/qualitymd/quality.md/internal/evaluator"
)

// Source packaging is runner-owned and deterministic: the same model, scope,
// and file state always produce the same bounded bundle with the same hashes.
// Evaluated source content is packaged as data, never as instructions.

const (
	// maxSourceFileBytes truncates one packaged file.
	maxSourceFileBytes = 64 * 1024
	// maxSourceBundleBytes caps one work unit's whole bundle.
	maxSourceBundleBytes = 512 * 1024
	// binarySniffBytes is how much of a file is inspected for binary content.
	binarySniffBytes = 8 * 1024
)

// skippedSourceDirs are dependency, VCS, and artifact directories excluded
// from source walking unless directly selected.
var skippedSourceDirs = map[string]struct{}{
	".git":         {},
	".quality":     {},
	"node_modules": {},
	"vendor":       {},
	"dist":         {},
}

// SourceBundle is the packaged, hashed source for one work unit.
type SourceBundle struct {
	Files []evaluator.SourceFile
	// Hash is the stable hash of the whole bundle.
	Hash string
	// Missing lists source refs that did not resolve to readable files.
	Missing []string
	// Truncated reports whether the bundle hit its size cap.
	Truncated bool
}

// areaSourceBundle returns the area's packaged source bundle, packaging it on
// first use and reusing the memoized bundle for every later work unit in the
// area.
func (e *engine) areaSourceBundle(areaRef string) (*SourceBundle, error) {
	if bundle, ok := e.sourceBundles[areaRef]; ok {
		return bundle, nil
	}
	bundle, err := packageSource(e.workspaceRoot, e.graph.Plan.Area(areaRef).Source)
	if err != nil {
		return nil, err
	}
	e.sourceBundles[areaRef] = bundle
	return bundle, nil
}

// packageSource resolves an area source ref relative to the workspace root
// and builds the bounded bundle. Paths are workspace-contained: an escaping
// ref is reported missing rather than followed.
func packageSource(workspaceRoot, sourceRef string) (*SourceBundle, error) {
	bundle := &SourceBundle{}
	if strings.TrimSpace(sourceRef) == "" {
		bundle.Hash = bundleHash(nil)
		return bundle, nil
	}
	clean := filepath.Clean(filepath.FromSlash(sourceRef))
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		bundle.Missing = append(bundle.Missing, sourceRef)
		bundle.Hash = bundleHash(nil)
		return bundle, nil
	}
	abs := filepath.Join(workspaceRoot, clean)
	info, err := os.Stat(abs)
	if err != nil {
		bundle.Missing = append(bundle.Missing, sourceRef)
		bundle.Hash = bundleHash(nil)
		return bundle, nil
	}
	var paths []string
	if info.IsDir() {
		paths, err = walkSourceDir(abs)
		if err != nil {
			return nil, fmt.Errorf("walking source %s: %w", sourceRef, err)
		}
	} else {
		paths = []string{abs}
	}
	sort.Strings(paths)
	total := 0
	for _, path := range paths {
		if total >= maxSourceBundleBytes {
			bundle.Truncated = true
			break
		}
		file, ok, err := packageSourceFile(workspaceRoot, path)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if total+len(file.Content) > maxSourceBundleBytes {
			file.Content = file.Content[:maxSourceBundleBytes-total]
			file.Truncated = true
			bundle.Truncated = true
		}
		total += len(file.Content)
		bundle.Files = append(bundle.Files, file)
	}
	bundle.Hash = bundleHash(bundle.Files)
	return bundle, nil
}

func walkSourceDir(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if _, skip := skippedSourceDirs[entry.Name()]; skip && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func packageSourceFile(workspaceRoot, abs string) (evaluator.SourceFile, bool, error) {
	raw, err := os.ReadFile(abs)
	if err != nil {
		return evaluator.SourceFile{}, false, fmt.Errorf("reading source %s: %w", abs, err)
	}
	sniff := raw
	if len(sniff) > binarySniffBytes {
		sniff = sniff[:binarySniffBytes]
	}
	if bytes.IndexByte(sniff, 0) >= 0 {
		return evaluator.SourceFile{}, false, nil
	}
	sum := sha256.Sum256(raw)
	truncated := false
	if len(raw) > maxSourceFileBytes {
		raw = raw[:maxSourceFileBytes]
		truncated = true
	}
	rel, err := filepath.Rel(workspaceRoot, abs)
	if err != nil {
		rel = abs
	}
	return evaluator.SourceFile{
		Path:      filepath.ToSlash(rel),
		Content:   string(raw),
		SHA256:    hex.EncodeToString(sum[:]),
		Truncated: truncated,
	}, true, nil
}

func bundleHash(files []evaluator.SourceFile) string {
	h := sha256.New()
	for _, file := range files {
		_, _ = fmt.Fprintf(h, "%s\x00%s\x00", file.Path, file.SHA256)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// hashJSON fingerprints a value through its canonical JSON encoding
// (encoding/json sorts map keys, so the hash is stable).
func hashJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}
