package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/evaluator"
)

func writeSourceFile(t *testing.T, root, rel, content string) {
	t.Helper()
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatalf("MkdirAll(%s) error = %v", rel, err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile(%s) error = %v", rel, err)
	}
}

func bundlePaths(bundle *SourceBundle) []string {
	paths := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		paths = append(paths, file.Path)
	}
	return paths
}

func TestPackageSourceDefaultSelectorWalksWorkspace(t *testing.T) {
	root := t.TempDir()
	writeSourceFile(t, root, "QUALITY.md", "---\ntitle: T\n---\n")
	writeSourceFile(t, root, "docs/guide.md", "guide")
	writeSourceFile(t, root, "node_modules/pkg/index.js", "skipped")

	bundle, err := packageSource(root, ".")
	if err != nil {
		t.Fatalf("packageSource(.) error = %v", err)
	}
	want := []string{"QUALITY.md", "docs/guide.md"}
	got := bundlePaths(bundle)
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("packaged = %v, want %v", got, want)
	}
}

func TestPackageSourceGlobSelectors(t *testing.T) {
	root := t.TempDir()
	writeSourceFile(t, root, "QUALITY.md", "---\ntitle: T\n---\n")
	writeSourceFile(t, root, "docs/guide.md", "guide")
	writeSourceFile(t, root, "docs/nested/deep.md", "deep")
	writeSourceFile(t, root, "docs/notes.txt", "notes")
	writeSourceFile(t, root, "node_modules/pkg/readme.md", "skipped")
	writeSourceFile(t, root, "vendor/lib/lib.go", "vendored")

	for _, tc := range []struct {
		name     string
		selector string
		want     []string
	}{
		{"segment star", "docs/*.md", []string{"docs/guide.md"}},
		{"recursive doublestar", "docs/**/*.md", []string{"docs/guide.md", "docs/nested/deep.md"}},
		{"root doublestar skips dependency dirs", "**/*.md", []string{"QUALITY.md", "docs/guide.md", "docs/nested/deep.md"}},
		{"literal prefix opts skipped dir back in", "vendor/**", []string{"vendor/lib/lib.go"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bundle, err := packageSource(root, tc.selector)
			if err != nil {
				t.Fatalf("packageSource(%q) error = %v", tc.selector, err)
			}
			got := bundlePaths(bundle)
			if len(got) != len(tc.want) {
				t.Fatalf("packaged = %v, want %v", got, tc.want)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("packaged = %v, want %v", got, tc.want)
				}
			}
			if len(bundle.Missing) != 0 {
				t.Fatalf("missing = %v, want none", bundle.Missing)
			}
		})
	}
}

func TestPackageSourceSkipsSymlinkedEntries(t *testing.T) {
	root := t.TempDir()
	writeSourceFile(t, root, "QUALITY.md", "---\ntitle: T\n---\n")
	writeSourceFile(t, root, "skills/quality/SKILL.md", "skill")
	if err := os.MkdirAll(filepath.Join(root, ".claude", "skills"), 0o755); err != nil {
		t.Fatalf("MkdirAll(.claude/skills) error = %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "skills", "quality"), filepath.Join(root, ".claude", "skills", "quality")); err != nil {
		t.Fatalf("Symlink(dir) error = %v", err)
	}
	if err := os.Symlink(filepath.Join(root, "QUALITY.md"), filepath.Join(root, "LINK.md")); err != nil {
		t.Fatalf("Symlink(file) error = %v", err)
	}

	bundle, err := packageSource(root, ".")
	if err != nil {
		t.Fatalf("packageSource(.) error = %v, want symlinks skipped without error", err)
	}
	for _, file := range bundle.Files {
		if file.Path == "LINK.md" || file.Path == ".claude/skills/quality/SKILL.md" {
			t.Fatalf("packaged = %v, want symlinked entries skipped", bundlePaths(bundle))
		}
	}
	if len(bundle.Files) != 2 {
		t.Fatalf("packaged = %v, want QUALITY.md and skills/quality/SKILL.md", bundlePaths(bundle))
	}
}

func TestPackageSourceUnresolvedSelectorIsMissing(t *testing.T) {
	root := t.TempDir()
	writeSourceFile(t, root, "QUALITY.md", "---\ntitle: T\n---\n")
	for _, selector := range []string{"absent", "absent/**/*.md", "../escape", "docs/*.nope"} {
		bundle, err := packageSource(root, selector)
		if err != nil {
			t.Fatalf("packageSource(%q) error = %v", selector, err)
		}
		if len(bundle.Files) != 0 || len(bundle.Missing) != 1 || bundle.Missing[0] != selector {
			t.Fatalf("packageSource(%q) = files %v missing %v, want no files and the selector missing",
				selector, bundlePaths(bundle), bundle.Missing)
		}
	}
}

func TestMatchSourceGlob(t *testing.T) {
	for _, tc := range []struct {
		pattern string
		name    string
		want    bool
	}{
		{"docs/*.md", "docs/guide.md", true},
		{"docs/*.md", "docs/nested/deep.md", false},
		{"docs/**/*.md", "docs/guide.md", true},
		{"docs/**/*.md", "docs/nested/deep.md", true},
		{"**/*.md", "QUALITY.md", true},
		{"**", "docs/nested/deep.md", true},
		{"docs/gu?de.[lm]d", "docs/guide.md", true},
		{"docs/*.md", "docs/notes.txt", false},
	} {
		if got := matchSourceGlob(tc.pattern, tc.name); got != tc.want {
			t.Errorf("matchSourceGlob(%q, %q) = %v, want %v", tc.pattern, tc.name, got, tc.want)
		}
	}
}

const inheritedSourceModel = `---
title: Inherited source model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: Minimum
    description: Minimum.
    criterion: Barely meets it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
requirements:
  has-tests:
    title: Has tests
    assessment: Inspect tests.
    factors: [reliability]
areas:
  api:
    title: API
    factors:
      durability:
        title: Durability
        description: Durability.
    requirements:
      has-docs:
        title: Has docs
        assessment: Inspect docs.
        factors: [durability]
---
`

// TestRunResolvesDefaultAndInheritedSource covers R1/R2 end to end: a
// source-less root resolves to the document's directory and a source-less
// child area inherits it, so judgment runs against real evidence instead of
// an empty bundle.
func TestRunResolvesDefaultAndInheritedSource(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeSourceFile(t, repo, "QUALITY.md", inheritedSourceModel)
	writeSourceFile(t, repo, "src/main.txt", "package main\n")

	scripted := newScriptedEvaluator()
	captured := &sourceCapturingEvaluator{inner: scripted}
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(captured),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	if !captured.sawPath("QUALITY.md") || !captured.sawPath("src/main.txt") {
		t.Fatalf("packaged source paths = %v, want the document directory walked for the inherited default", captured.paths())
	}
}

// TestRunFailsSourceUnavailable: a filesystem-kind selector that resolves to
// nothing fails the run with the classified source_unavailable category
// naming the selector, instead of judging against empty evidence. The
// selector is a glob so its kind stays filesystem under detection — a bare
// path naming nothing detects as prose and classifies differently (see
// TestRunFailsSelectorUnsupported).
func TestRunFailsSourceUnavailable(t *testing.T) {
	repo := testRunnerRepo(t)
	model := strings.ReplaceAll(testModel, "source: src", "source: does-not-exist/**/*.md")
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(model), 0o644); err != nil {
		t.Fatalf("WriteFile(QUALITY.md) error = %v", err)
	}

	scripted := newScriptedEvaluator()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusFailed {
		t.Fatalf("status = %q, want failed", result.Status)
	}
	if result.Failure == nil || result.Failure.Category != evaluator.FailureSourceUnavailable {
		t.Fatalf("failure = %+v, want %s", result.Failure, evaluator.FailureSourceUnavailable)
	}
	if !strings.Contains(result.Failure.Detail, "does-not-exist") {
		t.Fatalf("failure detail = %q, want the unresolved selector named", result.Failure.Detail)
	}
	if len(scripted.calls) != 0 {
		t.Fatalf("evaluator calls = %v, want none for an unresolved source", scripted.calls)
	}
}

// TestDetectSourceKind covers R1's detection order: glob metacharacters win,
// escaping and absolute selectors stay filesystem selectors (containment
// keeps them unresolvable — they never fall back to prose), an existing
// filesystem entry is a path, and only a selector that cannot be filesystem
// material is prose.
func TestDetectSourceKind(t *testing.T) {
	root := t.TempDir()
	writeSourceFile(t, root, "QUALITY.md", "---\ntitle: T\n---\n")
	writeSourceFile(t, root, "docs/guide.md", "guide")
	for _, tc := range []struct {
		selector string
		want     SourceKind
	}{
		{".", SourceKindPath},
		{"docs", SourceKindPath},
		{"docs/guide.md", SourceKindPath},
		{"docs/**/*.md", SourceKindGlob},
		{"specs/*.md", SourceKindGlob},
		{"../shared", SourceKindPath},
		{"/etc/passwd", SourceKindPath},
		{"docs/guids", SourceKindProse},
		{"all specs and open tickets", SourceKindProse},
		{"the deployed API", SourceKindProse},
	} {
		if got := detectSourceKind(root, tc.selector); got != tc.want {
			t.Errorf("detectSourceKind(%q) = %q, want %q", tc.selector, got, tc.want)
		}
	}
}

// TestCaptureResolvedSource covers the capture contract for
// resolver-returned material: unique non-empty paths, per-file and per-bundle
// caps with truncation marks, SHA-256 over the full returned content, and
// the shared bundle hash function.
func TestCaptureResolvedSource(t *testing.T) {
	t.Run("valid files capture with hashes", func(t *testing.T) {
		bundle, err := captureResolvedSource(map[string]any{
			"files": []any{
				map[string]any{"path": "tickets/T-1", "content": "first"},
				map[string]any{"path": "https://example.com/T-2", "content": "second"},
			},
		})
		if err != nil {
			t.Fatalf("captureResolvedSource() error = %v", err)
		}
		if len(bundle.Files) != 2 || bundle.Files[0].Path != "tickets/T-1" {
			t.Fatalf("files = %v", bundlePaths(bundle))
		}
		if bundle.Files[0].SHA256 == "" || bundle.Hash == "" {
			t.Errorf("captured bundle lacks hashes: %+v", bundle)
		}
		if bundle.Hash != bundleHash(bundle.Files) {
			t.Errorf("bundle hash differs from the shared bundle hash function")
		}
	})
	t.Run("rejects empty duplicate and malformed files", func(t *testing.T) {
		for name, payload := range map[string]map[string]any{
			"no files":       {"files": []any{}},
			"missing files":  {},
			"empty path":     {"files": []any{map[string]any{"path": " ", "content": "x"}}},
			"duplicate path": {"files": []any{map[string]any{"path": "a", "content": "x"}, map[string]any{"path": "a", "content": "y"}}},
			"non-string":     {"files": []any{map[string]any{"path": "a", "content": 42}}},
		} {
			if _, err := captureResolvedSource(payload); err == nil {
				t.Errorf("captureResolvedSource(%s) accepted invalid material", name)
			}
		}
	})
	t.Run("caps apply with truncation marks", func(t *testing.T) {
		big := strings.Repeat("x", maxSourceFileBytes+1)
		bundle, err := captureResolvedSource(map[string]any{
			"files": []any{map[string]any{"path": "big", "content": big}},
		})
		if err != nil {
			t.Fatalf("captureResolvedSource() error = %v", err)
		}
		if len(bundle.Files[0].Content) != maxSourceFileBytes || !bundle.Files[0].Truncated {
			t.Errorf("file cap not applied: len %d truncated %v", len(bundle.Files[0].Content), bundle.Files[0].Truncated)
		}
	})
}

const proseSourceModel = `---
title: Prose source model
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: minimum
    title: Minimum
    description: Minimum.
    criterion: Barely meets it.
factors:
  reliability:
    title: Reliability
requirements:
  has-tests:
    title: Open tickets are triaged
    assessment: Inspect the open tickets.
    factors: [reliability]
source: open tickets in the support queue
---
`

// TestRunFailsSelectorUnsupported: a prose selector under an evaluator that
// cannot serve source resolution requests fails the run at plan time with
// selector_unsupported — distinct from source_unavailable — naming the
// selector, its kind, and the remedy, before any judgment is dispatched.
func TestRunFailsSelectorUnsupported(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeSourceFile(t, repo, "QUALITY.md", proseSourceModel)

	scripted := newScriptedEvaluator()
	result, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		SelectEvaluator: scriptedSelect(scripted),
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Status != StatusFailed {
		t.Fatalf("status = %q, want failed", result.Status)
	}
	if result.Failure == nil || result.Failure.Category != evaluator.FailureSelectorUnsupported {
		t.Fatalf("failure = %+v, want %s", result.Failure, evaluator.FailureSelectorUnsupported)
	}
	for _, want := range []string{"open tickets in the support queue", "prose", "harness"} {
		if !strings.Contains(result.Failure.Detail, want) {
			t.Errorf("failure detail %q does not name %q", result.Failure.Detail, want)
		}
	}
	if len(scripted.calls) != 0 {
		t.Errorf("evaluator calls = %v, want none before the plan-time check", scripted.calls)
	}
	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	record := artifact.Sources["area:root"]
	if record == nil || record.Kind != string(SourceKindProse) || record.Resolver != ResolverHarness {
		t.Errorf("sources record = %+v, want pinned prose/harness", record)
	}
}

// TestHarnessResolvesProseSelector covers R2/R3/R5/R6 end to end: the run's
// first checkpoint is the resolution work request (empty source, selector in
// context), the returned material is captured — bounded, hashed, persisted
// with provenance — before dependent judgment, judgment requests carry the
// captured bundle as data, and a workspace write after capture does not
// invalidate the frozen prose bundle.
func TestHarnessResolvesProseSelector(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeSourceFile(t, repo, "QUALITY.md", proseSourceModel)

	scripted := newScriptedEvaluator()
	result := startHarnessRun(t, repo)
	if result.Status != StatusAwaitingEvaluator {
		t.Fatalf("status = %q, want awaiting_evaluator", result.Status)
	}
	first := soleRequest(t, result)
	if first.Kind != string(KindResolveSource) || first.Subject != "area:root" {
		t.Fatalf("first checkpoint = %s %s, want the resolution request", first.Kind, first.Subject)
	}
	if len(first.Source) != 0 {
		t.Fatalf("resolution request carries pre-gathered source: %v", first.Source)
	}
	selectorCtx, _ := first.Context["sourceSelector"].(map[string]any)
	if selectorCtx == nil || selectorCtx["selector"] != "open tickets in the support queue" || selectorCtx["kind"] != string(SourceKindProse) {
		t.Fatalf("sourceSelector context = %v, want the selector and its kind", first.Context["sourceSelector"])
	}

	sawCapturedSource := false
	checkpoints := 0
	for result.Status == StatusAwaitingEvaluator {
		checkpoints++
		if checkpoints > 20 {
			t.Fatalf("run did not complete after %d checkpoints", checkpoints)
		}
		envelopes := make([]evaluator.HarnessResultEnvelope, 0, len(result.EvaluatorRequests))
		for _, req := range result.EvaluatorRequests {
			if UnitKind(req.Kind) == KindAssessRateRequirement {
				if len(req.Source) != 1 || req.Source[0].Path != "tickets/T-1" {
					t.Fatalf("judgment request source = %v, want the captured bundle", req.Source)
				}
				sawCapturedSource = true
				// A workspace write after capture must not invalidate the
				// frozen prose bundle (reproducibility of record, no re-stat).
				writeSourceFile(t, repo, "unrelated.md", "written mid-run")
			}
			envelopes = append(envelopes, envelopeFor(t, scripted, req))
		}
		var err error
		result, err = submitHarnessResults(t, repo, result.Path, envelopes)
		if err != nil {
			t.Fatalf("submit batch error = %v", err)
		}
	}
	if result.Status != StatusCompleted {
		t.Fatalf("status = %q (failure: %+v), want completed", result.Status, result.Failure)
	}
	if !sawCapturedSource {
		t.Errorf("no judgment request carried the captured bundle")
	}

	artifact, err := NewStore(filepath.Join(repo, result.Path)).Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	record := artifact.Sources["area:root"]
	if record == nil {
		t.Fatalf("no sources record for area:root")
	}
	if record.Kind != string(SourceKindProse) || record.Resolver != ResolverHarness || record.HarnessRuntime != "claude-code" {
		t.Errorf("provenance = %+v, want prose/harness/claude-code", record)
	}
	if record.BundleHash == "" || record.CapturedAt == "" {
		t.Errorf("captured record lacks bundle hash or timestamp: %+v", record)
	}
	if len(record.Files) != 1 || record.Files[0].Path != "tickets/T-1" ||
		record.Files[0].Content == "" || record.Files[0].SHA256 == "" {
		t.Errorf("captured files = %+v, want the gathered material with content and hash", record.Files)
	}
}

// TestHarnessProsePinsKindAcrossResume: kind is detected once, at run
// creation, and honored on resume — a matching filesystem entry appearing
// mid-run cannot re-dispatch the selector to the deterministic resolver or
// change the graph shape.
func TestHarnessProsePinsKindAcrossResume(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeSourceFile(t, repo, "QUALITY.md", proseSourceModel)

	first := startHarnessRun(t, repo)
	firstRequest := soleRequest(t, first)
	if firstRequest.Kind != string(KindResolveSource) {
		t.Fatalf("first checkpoint = %s, want resolveSource", firstRequest.Kind)
	}
	// The selector now names an existing filesystem entry.
	writeSourceFile(t, repo, "open tickets in the support queue", "now a file")

	again, err := Run(context.Background(), Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Resume:          first.Path,
		SelectEvaluator: harnessSelect(),
	})
	if err != nil {
		t.Fatalf("Run(resume) error = %v", err)
	}
	againRequest := soleRequest(t, again)
	if againRequest.Kind != string(KindResolveSource) || againRequest.RequestID != firstRequest.RequestID {
		t.Fatalf("resumed checkpoint = %+v, want the identical pinned resolution request", againRequest)
	}
}

// TestDryRunSurfacesSourceDispatchPlan: dry-run names each area's detected
// kind and serving resolver, and the work-unit counts include the resolution
// unit, so the dispatch plan is visible before anything runs.
func TestDryRunSurfacesSourceDispatchPlan(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	writeSourceFile(t, repo, "QUALITY.md", proseSourceModel)

	preview, err := DryRun(Options{
		RepoRoot:        repo,
		Model:           filepath.Join(repo, "QUALITY.md"),
		Evaluator:       "harness",
		SelectEvaluator: harnessSelect(),
	})
	if err != nil {
		t.Fatalf("DryRun() error = %v", err)
	}
	if len(preview.Sources) != 1 {
		t.Fatalf("sources = %+v, want one area entry", preview.Sources)
	}
	entry := preview.Sources[0]
	if entry.Area != "area:root" || entry.Selector != "open tickets in the support queue" ||
		entry.Kind != string(SourceKindProse) || entry.Resolver != ResolverHarness {
		t.Errorf("source plan = %+v, want the prose detection named", entry)
	}
	// 1 evaluation frame, 1 area frame, 1 resolution, 1 requirement frame,
	// 1 combined judgment, 1 factor x2, 1 area analysis x2, 3 advice,
	// 1 report build.
	if preview.WorkUnits.EvaluatorUnits != 7 {
		t.Errorf("evaluator units = %d, want 7 including the resolution unit", preview.WorkUnits.EvaluatorUnits)
	}
}

// sourceCapturingEvaluator records packaged source paths seen in requests.
type sourceCapturingEvaluator struct {
	inner evaluator.Evaluator
	seen  []string
}

func (s *sourceCapturingEvaluator) Kind() string { return s.inner.Kind() }

func (s *sourceCapturingEvaluator) Capabilities() evaluator.Capabilities {
	return evaluator.Capabilities{}
}

func (s *sourceCapturingEvaluator) Evaluate(ctx context.Context, req evaluator.WorkRequest) (evaluator.WorkResult, error) {
	for _, file := range req.Source {
		s.seen = append(s.seen, file.Path)
	}
	return s.inner.Evaluate(ctx, req)
}

func (s *sourceCapturingEvaluator) sawPath(path string) bool {
	for _, seen := range s.seen {
		if seen == path {
			return true
		}
	}
	return false
}

func (s *sourceCapturingEvaluator) paths() []string { return s.seen }
