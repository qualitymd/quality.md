package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeModelFixture writes a model exercising nested factors, child areas, and
// requirements declared both directly under an area and under a factor.
func writeModelFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte(`---
title: Example
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
factors:
  performance:
    title: Performance
    factors:
      latency:
        title: Latency
requirements:
  has-readme:
    title: Has a README
    factors: [performance]
    assessment: Check.
areas:
  client-app:
    title: Client App
    factors:
      performance:
        title: Performance
        factors:
          latency:
            title: Latency
        requirements:
          fast-startup:
            title: Fast startup
            assessment: Measure.
    requirements:
      has-tests:
        title: Has tests
        factors: [performance]
        assessment: Run.
---
`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

// runModel executes the model command tree, returning stdout and the mapped
// exit code. stderr is discarded.
func runModel(t *testing.T, args ...string) (string, int) {
	t.Helper()
	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs(args)
	code := codeFor(cmd.Execute())
	return out.String(), code
}

func TestModelTreeJSONNestsWithIDs(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "tree", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; out = %s", code, ExitOK, out)
	}
	var root struct {
		ID       string `json:"id"`
		Kind     string `json:"kind"`
		Label    string `json:"label"`
		Children []struct {
			ID   string `json:"id"`
			Kind string `json:"kind"`
		} `json:"children"`
	}
	if err := json.Unmarshal([]byte(out), &root); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; out = %s", err, out)
	}
	if root.ID != "area:root" || root.Kind != "area" || root.Label != "Example" {
		t.Fatalf("root = %+v", root)
	}
	// Root children begin with factors, then requirements, then child areas.
	wantFirst := []string{"factor:root::performance", "requirement:root::has-readme", "area:client-app"}
	if len(root.Children) != len(wantFirst) {
		t.Fatalf("root children = %+v, want %v", root.Children, wantFirst)
	}
	for i, want := range wantFirst {
		if root.Children[i].ID != want {
			t.Fatalf("child %d id = %q, want %q", i, root.Children[i].ID, want)
		}
	}
}

func TestModelTreeAreaRootsSubtree(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "tree", "--area", "area:client-app", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; out = %s", code, ExitOK, out)
	}
	if !strings.HasPrefix(out, "area:client-app") {
		t.Fatalf("tree did not root at subtree; out = %q", out)
	}
	if strings.Contains(out, "area:root") {
		t.Fatalf("subtree leaked root; out = %q", out)
	}

	if _, code := runModel(t, "model", "tree", "--area", "client-app", path); code != ExitUsage {
		t.Fatalf("bare-path --area exit = %d, want %d", code, ExitUsage)
	}
	if _, code := runModel(t, "model", "tree", "--area", "area:nope", path); code != ExitUsage {
		t.Fatalf("unknown --area exit = %d, want %d", code, ExitUsage)
	}
}

func TestModelTreeDepthZeroEmitsRootedNodeOnly(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "tree", "--depth", "0", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d", code, ExitOK)
	}
	if strings.Count(strings.TrimSpace(out), "\n") != 0 {
		t.Fatalf("--depth 0 emitted more than the rooted node; out = %q", out)
	}
	if !strings.HasPrefix(out, "area:root") {
		t.Fatalf("out = %q, want only the rooted node", out)
	}
}

func TestModelListEnumeratesEveryElement(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "list", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; out = %s", code, ExitOK, out)
	}
	var rows []struct {
		ID       string `json:"id"`
		Kind     string `json:"kind"`
		Label    string `json:"label"`
		ParentID string `json:"parentId"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; out = %s", err, out)
	}
	got := map[string]string{}
	for _, r := range rows {
		got[r.ID] = r.ParentID
	}
	want := map[string]string{
		"area:root":                            "",
		"factor:root::performance":             "area:root",
		"factor:root::performance/latency":     "factor:root::performance",
		"requirement:root::has-readme":         "area:root",
		"area:client-app":                      "area:root",
		"factor:client-app::performance":       "area:client-app",
		"requirement:client-app::fast-startup": "factor:client-app::performance",
		"requirement:client-app::has-tests":    "area:client-app",
	}
	for id, parent := range want {
		gotParent, ok := got[id]
		if !ok {
			t.Errorf("missing element %q", id)
			continue
		}
		if gotParent != parent {
			t.Errorf("%s parentId = %q, want %q", id, gotParent, parent)
		}
	}
}

func TestModelListTypeFilter(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "list", "--type", "factor", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d", code, ExitOK)
	}
	var rows []struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; out = %s", err, out)
	}
	if len(rows) == 0 {
		t.Fatal("no factors returned")
	}
	for _, r := range rows {
		if r.Kind != "factor" {
			t.Fatalf("got kind %q, want only factor", r.Kind)
		}
	}
	if _, code := runModel(t, "model", "list", "--type", "bogus", path); code != ExitUsage {
		t.Fatalf("--type bogus exit = %d, want %d", code, ExitUsage)
	}
}

func TestModelListAreaAndTypeCombine(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "list", "--area", "area:client-app", "--type", "requirement", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d", code, ExitOK)
	}
	var rows []struct {
		ID   string `json:"id"`
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; out = %s", err, out)
	}
	for _, r := range rows {
		if r.Kind != "requirement" {
			t.Fatalf("got kind %q, want only requirement", r.Kind)
		}
		if !strings.HasPrefix(r.ID, "requirement:client-app::") {
			t.Fatalf("got id %q outside client-app subtree", r.ID)
		}
	}
	if len(rows) != 2 {
		t.Fatalf("got %d requirements, want 2 (has-tests, fast-startup)", len(rows))
	}
}

func TestModelGetFactorDetail(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "get", "factor:client-app::performance", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; out = %s", code, ExitOK, out)
	}
	var detail struct {
		ID           string   `json:"id"`
		Kind         string   `json:"kind"`
		Factors      []string `json:"factors"`
		Requirements []string `json:"requirements"`
	}
	if err := json.Unmarshal([]byte(out), &detail); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; out = %s", err, out)
	}
	if detail.Kind != "factor" {
		t.Fatalf("kind = %q, want factor", detail.Kind)
	}
	if len(detail.Factors) != 1 || detail.Factors[0] != "factor:client-app::performance/latency" {
		t.Fatalf("sub-factor ids = %v", detail.Factors)
	}
	if len(detail.Requirements) != 1 || detail.Requirements[0] != "requirement:client-app::fast-startup" {
		t.Fatalf("requirement ids = %v", detail.Requirements)
	}
}

func TestModelGetUnknownExitsUsageWithSuggestions(t *testing.T) {
	path := writeModelFixture(t)
	var stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"model", "get", "requirement:nope::missing", path})
	err := cmd.Execute()
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("exit = %d, want %d", got, ExitUsage)
	}
	if err == nil || !strings.Contains(err.Error(), "requirement:nope::missing") {
		t.Fatalf("error = %v, want it to name the unresolved id", err)
	}
	if !strings.Contains(err.Error(), "did you mean") {
		t.Fatalf("error = %v, want near-match suggestions", err)
	}
}

func TestModelCanonicalIDsRoundTrip(t *testing.T) {
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "list", "--json", path)
	if code != ExitOK {
		t.Fatalf("list exit = %d, want %d", code, ExitOK)
	}
	var rows []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	for _, r := range rows {
		got, code := runModel(t, "model", "get", r.ID, "--json", path)
		if code != ExitOK {
			t.Fatalf("get %s exit = %d, want %d", r.ID, code, ExitOK)
		}
		var detail struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal([]byte(got), &detail); err != nil {
			t.Fatalf("get %s unmarshal error = %v", r.ID, err)
		}
		if detail.ID != r.ID {
			t.Fatalf("round-trip id = %q, want %q", detail.ID, r.ID)
		}
	}
}

func TestModelJSONIsByteStableAndPlain(t *testing.T) {
	path := writeModelFixture(t)
	first, code := runModel(t, "model", "list", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d", code, ExitOK)
	}
	for i := 0; i < 3; i++ {
		got, _ := runModel(t, "model", "list", "--json", path)
		if got != first {
			t.Fatalf("list --json not byte-stable across runs")
		}
	}
	if strings.Contains(first, "\x1b[") {
		t.Fatalf("--json output carried styling bytes: %q", first)
	}
}

func TestModelUnreadableFileExitsInternal(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.md")
	var stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"model", "tree", missing})
	err := cmd.Execute()
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("exit = %d, want %d", got, ExitInternal)
	}
	if err == nil || !strings.Contains(err.Error(), "lint") {
		t.Fatalf("error = %v, want it to point at lint", err)
	}
}

func TestModelGetReadsModelByPath(t *testing.T) {
	// A model file is read by path with no run awareness; a snapshot is just a
	// model file, so the same path-based read yields that file's IDs.
	path := writeModelFixture(t)
	out, code := runModel(t, "model", "get", "area:root", "--json", path)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d; out = %s", code, ExitOK, out)
	}
	if !strings.Contains(out, `"id": "area:root"`) {
		t.Fatalf("out = %q, want area:root detail", out)
	}
}
