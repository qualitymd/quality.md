package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/qualitymd/quality.md/internal/evaluation"
)

func TestEvaluationAssessmentAddCommandAcceptsBatch(t *testing.T) {
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
	cmd.SetIn(strings.NewReader(`[
  {
    "targetPath": [],
    "requirement": "Has tests",
    "criterionSource": "rating-scale",
    "findings": [],
    "recommendations": [],
    "factorPaths": [],
    "ratingResult": {
      "kind": "rated",
      "level": "target",
      "rationale": "Tests cover the requirement."
    }
  }
]`))
	cmd.SetArgs([]string{"evaluation", "assessment-result", "add", "--file", "-", "--json", runPath})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"paths": [`) || !strings.Contains(out.String(), "assessment-results/001-root-has-tests.json") {
		t.Fatalf("stdout = %s, want batched write receipt", out.String())
	}
	raw, err := os.ReadFile(filepath.Join(runPath, "assessment-results", "001-root-has-tests.json"))
	if err != nil {
		t.Fatalf("reading assessment result record: %v", err)
	}
	if !strings.Contains(string(raw), `"schemaVersion": 1`) || !strings.Contains(string(raw), `"requirement": "Has tests"`) {
		t.Fatalf("assessment result record = %s, want canonical JSON", raw)
	}
}

func TestEvaluationListAndLatestStatusCommands(t *testing.T) {
	repo := testEvaluationRepo(t)
	if _, err := evaluation.CreateRun(evaluation.Options{RepoRoot: repo, Subject: "QUALITY.md"}); err != nil {
		t.Fatalf("CreateRun() error = %v", err)
	}
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Fatalf("restore cwd: %v", err)
		}
	})
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("chdir repo: %v", err)
	}

	cmd := newRootCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "list", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"runs": [`) || !strings.Contains(out.String(), `"reportable": false`) {
		t.Fatalf("list stdout = %s, want run list", out.String())
	}

	cmd = newRootCmd()
	out.Reset()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "status", "--latest", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("status Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), `"path": `) || !strings.Contains(out.String(), `"reportable": false`) {
		t.Fatalf("status stdout = %s, want latest status", out.String())
	}
}

func TestEvaluationOldCommandNamesAreRejected(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"evaluation", "create-run"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if !strings.Contains(err.Error(), `unknown command "create-run"`) {
		t.Fatalf("Execute() error = %v, want unknown old command", err)
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
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
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
