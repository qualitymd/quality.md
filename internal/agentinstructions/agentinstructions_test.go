package agentinstructions

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUpdateCreatesAgentsPointer(t *testing.T) {
	dir := t.TempDir()

	results, err := Update(UpdateOptions{Dir: dir, ModelPath: "QUALITY.md"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1: %#v", len(results), results)
	}
	if results[0].Path != "AGENTS.md" || !results[0].Created || results[0].Updated {
		t.Fatalf("result = %#v, want created AGENTS.md", results[0])
	}
	got := readFile(t, filepath.Join(dir, "AGENTS.md"))
	want := "<!-- Added by qualitymd init. -->\nSee [QUALITY.md](QUALITY.md) for this project's quality model.\n"
	if got != want {
		t.Fatalf("AGENTS.md = %q, want %q", got, want)
	}
}

func TestUpdateUpdatesDetectedInstructionFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "CLAUDE.md"), "# Claude\n")
	writeFile(t, filepath.Join(dir, "GEMINI.md"), "# Gemini")

	results, err := Update(UpdateOptions{Dir: dir, ModelPath: "docs/QUALITY.md"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("len(results) = %d, want 3: %#v", len(results), results)
	}
	for _, name := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"} {
		got := readFile(t, filepath.Join(dir, name))
		if !strings.Contains(got, "See [QUALITY.md](docs/QUALITY.md) for this project's quality model.") {
			t.Fatalf("%s missing relative pointer: %q", name, got)
		}
	}
}

func TestUpdateIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	if _, err := Update(UpdateOptions{Dir: dir, ModelPath: "QUALITY.md"}); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	results, err := Update(UpdateOptions{Dir: dir, ModelPath: "QUALITY.md"})
	if err != nil {
		t.Fatalf("second Update() error = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("second results = %#v, want none", results)
	}
	got := readFile(t, filepath.Join(dir, "AGENTS.md"))
	if count := strings.Count(got, "<!-- Added by qualitymd init. -->"); count != 1 {
		t.Fatalf("marker count = %d, want 1 in %q", count, got)
	}
}

func TestUpdateDeduplicatesSymlinkedInstructionFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "AGENTS.md"), "# Agents\n")
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "CLAUDE.md")); err != nil {
		t.Skipf("os.Symlink() error = %v", err)
	}
	if err := os.Symlink("AGENTS.md", filepath.Join(dir, "GEMINI.md")); err != nil {
		t.Skipf("os.Symlink() error = %v", err)
	}

	results, err := Update(UpdateOptions{Dir: dir, ModelPath: "QUALITY.md"})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1: %#v", len(results), results)
	}
	got := readFile(t, filepath.Join(dir, "AGENTS.md"))
	if count := strings.Count(got, "<!-- Added by qualitymd init. -->"); count != 1 {
		t.Fatalf("marker count = %d, want 1 in %q", count, got)
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile(%q) error = %v", path, err)
	}
	return string(data)
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v", path, err)
	}
}
