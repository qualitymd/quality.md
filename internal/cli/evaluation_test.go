package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

func TestEvaluationSetPlannedCoverageCommand(t *testing.T) {
	repo := testEvaluationRepo(t)
	run, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Subject: "QUALITY.md"})
	if err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	runPath := filepath.Join(repo, run.Path)

	cmd := newRootCmd()
	var out, stderr bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetIn(strings.NewReader(`{
  "schemaVersion": 1,
  "assessments": [
    { "targetPath": [], "requirement": "Has tests" }
  ],
  "analyses": [
    { "targetPath": [] }
  ]
}`))
	cmd.SetArgs([]string{"evaluation", "set-planned-coverage", "--file", "-", "--json", runPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"path": `) || !strings.Contains(out.String(), "planned-coverage.json") {
		t.Fatalf("stdout = %s, want planned coverage write receipt", out.String())
	}
	raw, err := os.ReadFile(filepath.Join(runPath, "planned-coverage.json"))
	if err != nil {
		t.Fatalf("reading planned-coverage.json: %v", err)
	}
	if !strings.Contains(string(raw), `"schemaVersion": 1`) || !strings.Contains(string(raw), `"assessments": [`) {
		t.Fatalf("planned-coverage.json = %s, want canonical JSON", raw)
	}
}

func testEvaluationRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repo, "QUALITY.md"), []byte(`---
title: Test model
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  Has tests:
    assessment: Inspect tests.
---
`), 0o644); err != nil {
		t.Fatalf("write QUALITY.md: %v", err)
	}
	return repo
}
