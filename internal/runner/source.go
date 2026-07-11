package runner

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
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

// globMetacharacters marks a selector (or selector segment) as a glob pattern
// rather than a literal path.
const globMetacharacters = "*?["

// SourceKind classifies an effective source selector for resolver dispatch.
// Kinds are detected from the bare selector string at run creation and pinned
// in the run artifact, so filesystem changes mid-run cannot re-dispatch a
// selector to a different resolver.
type SourceKind string

const (
	// SourceKindPath is a literal filesystem path selector, resolved by the
	// deterministic workspace walk.
	SourceKindPath SourceKind = "path"
	// SourceKindGlob is a glob-pattern selector, expanded by the deterministic
	// workspace walk.
	SourceKindGlob SourceKind = "glob"
	// SourceKindProse is a selector describing a body of evidence that is not
	// workspace filesystem material; resolving it needs tools, so it is
	// dispatched to a resolution-capable evaluator.
	SourceKindProse SourceKind = "prose"
)

// Resolver names recorded in the run artifact's per-area sources provenance.
const (
	// ResolverWalk is the deterministic workspace walk for path and glob
	// selectors.
	ResolverWalk = "walk"
	// ResolverHarness is checkpoint-dispatched resolution served by the
	// invoking harness.
	ResolverHarness = "harness"
)

// detectSourceKind classifies one effective source selector. The filesystem
// interpretation always wins: glob metacharacters make a glob; a selector that
// is absolute or escapes the workspace stays a path (workspace containment
// keeps it unresolvable — it never falls back to prose); a selector naming an
// existing filesystem entry is a path; only a selector that cannot be
// filesystem material is prose.
func detectSourceKind(workspaceRoot, selector string) SourceKind {
	ref := strings.TrimSpace(selector)
	if strings.ContainsAny(ref, globMetacharacters) {
		return SourceKindGlob
	}
	if ref == "" || filepath.IsAbs(ref) {
		return SourceKindPath
	}
	cleaned := path.Clean(filepath.ToSlash(ref))
	if path.IsAbs(cleaned) || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return SourceKindPath
	}
	if _, err := os.Stat(filepath.Join(workspaceRoot, filepath.FromSlash(cleaned))); err == nil {
		return SourceKindPath
	}
	return SourceKindProse
}

// resolverForKind maps a pinned selector kind to the resolver that serves it.
func resolverForKind(kind SourceKind) string {
	if kind == SourceKindProse {
		return ResolverHarness
	}
	return ResolverWalk
}

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
	// Missing lists source selectors that did not resolve to readable files.
	Missing []string
	// Truncated reports whether the bundle hit its size cap.
	Truncated bool
}

// sourceUnavailableError reports an effective source selector that resolved
// to zero readable source files, so the affected work fails loudly as
// source_unavailable instead of being judged against empty evidence.
type sourceUnavailableError struct {
	areaRef  string
	selector string
}

func (e *sourceUnavailableError) Error() string {
	return fmt.Sprintf("%s source %q resolved to no readable source files; fix the source selector or add the material it names, then run again",
		e.areaRef, e.selector)
}

// sourceUnavailableFailure classifies err as a source_unavailable run failure
// when it reports an unresolved source selector; nil otherwise.
func sourceUnavailableFailure(err error) *Failure {
	var unavailable *sourceUnavailableError
	if !errors.As(err, &unavailable) {
		return nil
	}
	return &Failure{Category: evaluator.FailureSourceUnavailable, Detail: unavailable.Error()}
}

// selectorUnsupportedDetail words the selector_unsupported failure: the run
// cannot resolve this kind of selector, which is a different remedy — change
// the selector or the evaluator, not the material — than source_unavailable.
func selectorUnsupportedDetail(areaRef, selector string, kind SourceKind, evaluatorName string) string {
	return fmt.Sprintf("%s source %q is a %s selector, and evaluator %q cannot serve source resolution requests; "+
		"evaluate through harness dispatch (--evaluator harness) or change the selector to a workspace path or glob",
		areaRef, selector, kind, evaluatorName)
}

// areaSourceBundle returns the area's source bundle, dispatching on the run's
// pinned selector kind: path and glob selectors package through the unchanged
// deterministic walk on first use; a prose selector reads the bundle its
// resolution work unit captured into the run artifact. Either way the bundle
// is memoized for every later work unit in the area.
func (e *engine) areaSourceBundle(areaRef string) (*SourceBundle, error) {
	if bundle, ok := e.sourceBundles[areaRef]; ok {
		return bundle, nil
	}
	record := e.artifact.Sources[areaRef]
	if record != nil && SourceKind(record.Kind) == SourceKindProse {
		bundle, err := record.capturedBundle(areaRef)
		if err != nil {
			return nil, err
		}
		e.sourceBundles[areaRef] = bundle
		return bundle, nil
	}
	bundle, err := e.packageAreaSource(areaRef, record)
	if err != nil {
		return nil, err
	}
	e.sourceBundles[areaRef] = bundle
	return bundle, nil
}

// packageAreaSource runs the deterministic walk for a path or glob selector
// and completes the area's provenance record with the packaged file hashes.
// The plan's area source is the effective selector from the shared model
// resolver, so it is never empty; a selector that packages nothing — matching
// no files, or only unreadable ones — is a sourceUnavailableError. The
// document default (the file's directory) always packages at least the
// QUALITY.md file itself.
func (e *engine) packageAreaSource(areaRef string, record *SourceRecord) (*SourceBundle, error) {
	selector := e.graph.Plan.Area(areaRef).Source
	bundle, err := packageSource(e.workspaceRoot, selector)
	if err != nil {
		return nil, err
	}
	if len(bundle.Files) == 0 {
		return nil, &sourceUnavailableError{areaRef: areaRef, selector: selector}
	}
	if record != nil && record.BundleHash == "" {
		record.completeFromBundle(bundle, e.timestamp(), false)
	}
	return bundle, nil
}

// capturedBundle rebuilds the prose area's source bundle from the captured
// provenance record. The area's resolution work unit is a completed
// dependency of every consumer, so an uncaptured record is an internal error,
// never a silent empty bundle.
func (r *SourceRecord) capturedBundle(areaRef string) (*SourceBundle, error) {
	if r.BundleHash == "" || len(r.Files) == 0 {
		return nil, fmt.Errorf("internal error: %s has no captured source bundle for its resolved selector", areaRef)
	}
	bundle := &SourceBundle{Hash: r.BundleHash, Truncated: r.Truncated}
	for _, file := range r.Files {
		bundle.Files = append(bundle.Files, evaluator.SourceFile{
			Path:      file.Path,
			Content:   file.Content,
			SHA256:    file.SHA256,
			Truncated: file.Truncated,
		})
	}
	return bundle, nil
}

// captureResolvedSource validates a resolution result's returned material and
// captures it as the area's bounded, hashed source bundle: non-empty unique
// file paths, the walk's per-file and per-bundle caps with truncation marks,
// SHA-256 per file, and the same bundle hash function as walked source. Prose
// file paths are labels for gathered material (a ticket ID, a URL, a
// repo-relative path), recorded and hashed verbatim.
func captureResolvedSource(payload map[string]any) (*SourceBundle, error) {
	items, ok := payload["files"].([]any)
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("source resolution result must carry a non-empty files array")
	}
	bundle := &SourceBundle{}
	seen := map[string]struct{}{}
	total := 0
	for i, item := range items {
		entry, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("files[%d] must be an object", i)
		}
		filePath, _ := entry["path"].(string)
		if strings.TrimSpace(filePath) == "" {
			return nil, fmt.Errorf("files[%d] must carry a non-empty path", i)
		}
		if _, dup := seen[filePath]; dup {
			return nil, fmt.Errorf("files[%d] duplicates path %q", i, filePath)
		}
		seen[filePath] = struct{}{}
		content, ok := entry["content"].(string)
		if !ok {
			return nil, fmt.Errorf("files[%d] must carry string content", i)
		}
		sum := sha256.Sum256([]byte(content))
		file := evaluator.SourceFile{Path: filePath, SHA256: hex.EncodeToString(sum[:])}
		if len(content) > maxSourceFileBytes {
			content = content[:maxSourceFileBytes]
			file.Truncated = true
		}
		if total >= maxSourceBundleBytes {
			bundle.Truncated = true
			break
		}
		if total+len(content) > maxSourceBundleBytes {
			content = content[:maxSourceBundleBytes-total]
			file.Truncated = true
			bundle.Truncated = true
		}
		file.Content = content
		total += len(content)
		bundle.Files = append(bundle.Files, file)
	}
	bundle.Hash = bundleHash(bundle.Files)
	return bundle, nil
}

// packageSource resolves an area source selector — a literal path or a glob —
// relative to the workspace root and builds the bounded bundle. Selectors are
// workspace-contained: an escaping selector is reported missing rather than
// followed.
func packageSource(workspaceRoot, sourceRef string) (*SourceBundle, error) {
	bundle := &SourceBundle{}
	paths, err := resolveSourceSelector(workspaceRoot, sourceRef)
	if err != nil {
		return nil, fmt.Errorf("resolving source %s: %w", sourceRef, err)
	}
	sort.Strings(paths)
	total := 0
	for _, filePath := range paths {
		if total >= maxSourceBundleBytes {
			bundle.Truncated = true
			break
		}
		file, ok, err := packageSourceFile(workspaceRoot, filePath)
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
	if len(bundle.Files) == 0 && strings.TrimSpace(sourceRef) != "" {
		bundle.Missing = append(bundle.Missing, sourceRef)
	}
	bundle.Hash = bundleHash(bundle.Files)
	return bundle, nil
}

// resolveSourceSelector expands one selector to the absolute file paths it
// selects. A selector that is empty, escapes the workspace, or names nothing
// resolves to no paths without error.
func resolveSourceSelector(workspaceRoot, sourceRef string) ([]string, error) {
	ref := strings.TrimSpace(sourceRef)
	if ref == "" || filepath.IsAbs(ref) {
		return nil, nil
	}
	selector := path.Clean(filepath.ToSlash(ref))
	if path.IsAbs(selector) || selector == ".." || strings.HasPrefix(selector, "../") {
		return nil, nil
	}
	if strings.ContainsAny(selector, globMetacharacters) {
		return expandSourceGlob(workspaceRoot, selector)
	}
	abs := filepath.Join(workspaceRoot, filepath.FromSlash(selector))
	info, err := os.Stat(abs)
	if err != nil {
		return nil, nil
	}
	if info.IsDir() {
		return walkSourceDir(abs)
	}
	return []string{abs}, nil
}

// expandSourceGlob walks from the pattern's longest non-glob parent directory
// and keeps each file whose workspace-relative path matches the pattern.
// Because the walk starts at that literal prefix, a skipped directory named
// literally in the pattern (vendor/**) is walked as an explicit selection,
// while a recursive glob (**/*.md) still skips dependency and artifact
// directories.
func expandSourceGlob(workspaceRoot, pattern string) ([]string, error) {
	root := workspaceRoot
	if prefix := literalGlobPrefix(pattern); prefix != "" {
		root = filepath.Join(workspaceRoot, filepath.FromSlash(prefix))
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	walked, err := walkSourceDir(root)
	if err != nil {
		return nil, err
	}
	var matches []string
	for _, abs := range walked {
		rel, err := filepath.Rel(workspaceRoot, abs)
		if err != nil {
			continue
		}
		if matchSourceGlob(pattern, filepath.ToSlash(rel)) {
			matches = append(matches, abs)
		}
	}
	return matches, nil
}

// literalGlobPrefix returns the pattern's leading path segments before the
// first segment containing a glob metacharacter. A glob pattern always has
// such a segment, so the prefix never consumes the whole pattern.
func literalGlobPrefix(pattern string) string {
	segments := strings.Split(pattern, "/")
	literal := 0
	for _, segment := range segments {
		if strings.ContainsAny(segment, globMetacharacters) {
			break
		}
		literal++
	}
	return strings.Join(segments[:literal], "/")
}

// matchSourceGlob reports whether the slash-form workspace-relative path
// matches the selector pattern. Segments match with path.Match semantics; a
// `**` segment matches any number of path segments, including none.
func matchSourceGlob(pattern, name string) bool {
	return matchGlobSegments(strings.Split(pattern, "/"), strings.Split(name, "/"))
}

func matchGlobSegments(pattern, name []string) bool {
	for len(pattern) > 0 {
		if pattern[0] == "**" {
			for skip := 0; skip <= len(name); skip++ {
				if matchGlobSegments(pattern[1:], name[skip:]) {
					return true
				}
			}
			return false
		}
		if len(name) == 0 {
			return false
		}
		if ok, err := path.Match(pattern[0], name[0]); err != nil || !ok {
			return false
		}
		pattern, name = pattern[1:], name[1:]
	}
	return len(name) == 0
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
		// Only regular files are readable evidence: symlinked directories and
		// files, sockets, and devices are skipped, never walked or read.
		if !entry.Type().IsRegular() {
			return nil
		}
		paths = append(paths, path)
		return nil
	})
	return paths, err
}

func packageSourceFile(workspaceRoot, abs string) (evaluator.SourceFile, bool, error) {
	// The defensive Lstat keeps non-regular entries — however selected — from
	// reaching os.ReadFile, mirroring the walk's regular-file rule.
	info, err := os.Lstat(abs)
	if err != nil || !info.Mode().IsRegular() {
		return evaluator.SourceFile{}, false, nil
	}
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
