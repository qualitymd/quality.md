package cli

import (
	"bytes"
	"encoding/json"
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
  has-assessment:
    title: Has an assessment
    factors: [reliability]
    assessment: Inspect it.
factors:
  reliability:
    description: Reliability.
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
	if got := codeFor(err); got != ExitProblems {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitProblems)
	}
	if !strings.Contains(out.String(), `"valid": false`) || !strings.Contains(out.String(), `"ruleId": "too-few-levels"`) {
		t.Fatalf("stdout = %s, want JSON lint result", out.String())
	}
	if !strings.Contains(out.String(), `"command": "qualitymd lint `+path+`"`) {
		t.Fatalf("stdout = %s, want rerun next action", out.String())
	}
	if hasTerminalEscape(out.String()) {
		t.Fatalf("stdout = %q, want JSON without terminal escapes", out.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestLintHumanInvalidWritesNextActionToStderr(t *testing.T) {
	path := writeLintModel(t, `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
requirements:
  has-assessment:
    title: Has an assessment
    factors: [reliability]
    assessment: Inspect it.
factors:
  reliability:
    description: Reliability.
---
`)

	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"lint", path})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want lint error")
	}
	if !strings.Contains(out.String(), "too-few-levels") {
		t.Fatalf("stdout = %q, want lint findings", out.String())
	}
	if !strings.Contains(stderr.String(), "Next: qualitymd lint "+path) {
		t.Fatalf("stderr = %q, want next-action footer", stderr.String())
	}
}

func TestLintJSONFixableFindingPrefersFixAction(t *testing.T) {
	path := writeLintModel(t, `---
title: Example
description: ""
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
  has-assessment:
    title: Has an assessment
    factors: [reliability]
    assessment: Inspect it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
---
`)

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"lint", "--json", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var result struct {
		NextActions []struct {
			ID      string `json:"id"`
			Command string `json:"command"`
		} `json:"nextActions"`
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stdout = %s", err, out.String())
	}
	if len(result.NextActions) != 1 || result.NextActions[0].ID != "fix" || result.NextActions[0].Command != "qualitymd lint --fix "+path {
		t.Fatalf("nextActions = %#v, want fix action", result.NextActions)
	}
}

func TestLintHumanValidFile(t *testing.T) {
	path := writeLintModel(t, `---
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
requirements:
  has-assessment:
    title: Has an assessment
    factors: [reliability]
    assessment: Inspect it.
factors:
  reliability:
    title: Reliability
    description: Reliability.
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
	if got := codeFor(nil); got != ExitOK {
		t.Fatalf("codeFor(nil) = %d, want %d", got, ExitOK)
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
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitUsage)
	}
}

func TestLintMissingFileMapsToInternal(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"lint", filepath.Join(t.TempDir(), "missing.md")})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want missing-file error")
	}
	if got := codeFor(err); got != ExitInternal {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitInternal)
	}
}

func TestMalformedInvocationsMapToUsage(t *testing.T) {
	for name, args := range map[string][]string{
		"unknown flag":    {"lint", "--bogus"},
		"too many args":   {"lint", "one.md", "two.md"},
		"unknown command": {"bogus"},
	} {
		t.Run(name, func(t *testing.T) {
			cmd := newRootCmd()
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})
			cmd.SetArgs(args)
			err := cmd.Execute()
			if err == nil {
				t.Fatal("Execute() error = nil, want usage error")
			}
			if got := codeFor(err); got != ExitUsage {
				t.Fatalf("codeFor(error) = %d, want %d; err = %v", got, ExitUsage, err)
			}
		})
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
