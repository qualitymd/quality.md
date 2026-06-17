package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLintJSONExitsNonZeroOnErrors(t *testing.T) {
	path := writeLintModel(t, `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
requirements:
  "has an assessment":
    assessment: Inspect it.
---
`)

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"lint", "--json", path})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want lint error")
	}
	if !strings.Contains(out.String(), `"valid": false`) || !strings.Contains(out.String(), `"ruleId": "too-few-levels"`) {
		t.Fatalf("stdout = %s, want JSON lint result", out.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestLintHumanValidFile(t *testing.T) {
	path := writeLintModel(t, `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  "has an assessment":
    assessment: Inspect it.
---
`)

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"lint", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), path+" is valid.") {
		t.Fatalf("stdout = %q, want valid message", out.String())
	}
}

func TestLintRejectsStdinSentinel(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"lint", "-"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want stdin sentinel error")
	}
	if !strings.Contains(err.Error(), "does not read from stdin") {
		t.Fatalf("Execute() error = %v, want stdin message", err)
	}
}

func writeLintModel(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return path
}
