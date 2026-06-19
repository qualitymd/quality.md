package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatusJSONEmitsSnapshot(t *testing.T) {
	repo := newStatusRepo(t)
	path := writeStatusModel(t, repo)

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"status", "--json", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	var result struct {
		SchemaVersion int    `json:"schemaVersion"`
		Readiness     string `json:"readiness"`
		Model         struct {
			Valid bool `json:"valid"`
		} `json:"model"`
	}
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json.Unmarshal() error = %v; stdout = %s", err, out.String())
	}
	if result.SchemaVersion != 1 || result.Readiness != "ready-to-evaluate" || !result.Model.Valid {
		t.Fatalf("result = %#v, want ready valid snapshot", result)
	}
}

func TestStatusHumanIsCompact(t *testing.T) {
	repo := newStatusRepo(t)
	path := writeStatusModel(t, repo)

	var out bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"status", path})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"Status", "QUALITY.md: present, valid", "Readiness: ready-to-evaluate"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, want %q", out.String(), want)
		}
	}
	if strings.Contains(out.String(), "ruleId") {
		t.Fatalf("stdout = %q, want compact human output", out.String())
	}
}

func TestStatusRejectsStdinSentinelAsUsage(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"status", "-"})
	err := cmd.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want usage error")
	}
	if got := codeFor(err); got != ExitUsage {
		t.Fatalf("codeFor(error) = %d, want %d", got, ExitUsage)
	}
}

func TestStatusInvalidModelStillExitsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte("---\ntitle: Invalid\n---\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"status", path})
	if got := execute(context.Background(), cmd); got != ExitOK {
		t.Fatalf("execute() = %d, want %d", got, ExitOK)
	}
}

func newStatusRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	return repo
}

func writeStatusModel(t *testing.T, repo string) string {
	t.Helper()
	path := filepath.Join(repo, "QUALITY.md")
	if err := os.WriteFile(path, []byte(`---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  "starts":
    factors: [reliability]
    assessment: Run it.
factors:
  reliability:
    description: Reliability.
---
`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}
