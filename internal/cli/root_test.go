package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteMapsLintProblemsWithoutExtraStderr(t *testing.T) {
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

	if got := execute(context.Background(), cmd); got != ExitProblems {
		t.Fatalf("execute() = %d, want %d", got, ExitProblems)
	}
	if !strings.Contains(out.String(), `"valid": false`) {
		t.Fatalf("stdout = %s, want JSON lint result", out.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecuteMapsUsageAndInternalErrors(t *testing.T) {
	for name, tc := range map[string]struct {
		args []string
		code int
	}{
		"usage":    {args: []string{"lint", "--bogus"}, code: ExitUsage},
		"internal": {args: []string{"lint", filepath.Join(t.TempDir(), "missing.md")}, code: ExitInternal},
	} {
		t.Run(name, func(t *testing.T) {
			var stderr bytes.Buffer
			cmd := newRootCmd()
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&stderr)
			cmd.SetArgs(tc.args)

			if got := execute(context.Background(), cmd); got != tc.code {
				t.Fatalf("execute() = %d, want %d", got, tc.code)
			}
			if stderr.Len() == 0 {
				t.Fatal("stderr is empty, want rendered error")
			}
		})
	}
}

func TestExecuteMapsSuccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "QUALITY.md")
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
  "has an assessment":
    assessment: Inspect it.
---
`), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	cmd := newRootCmd()
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetArgs([]string{"lint", path})

	if got := execute(context.Background(), cmd); got != ExitOK {
		t.Fatalf("execute() = %d, want %d", got, ExitOK)
	}
}
