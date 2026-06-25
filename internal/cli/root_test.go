package cli

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecuteNoArgsShowsConciseWelcome(t *testing.T) {
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)

	if got := execute(context.Background(), cmd); got != ExitOK {
		t.Fatalf("execute() = %d, want %d", got, ExitOK)
	}

	want := "QUALITY.md\n\n" +
		"The companion CLI for the QUALITY.md file format for evaluating and improving\n" +
		"the quality of AI assistant projects and harnesses.\n\n" +
		"Designed to be used with the companion agent skill:\n" +
		"  npx skills add qualitymd/quality.md\n\n" +
		"Start:\n" +
		"  qualitymd init\n" +
		"  qualitymd lint QUALITY.md\n\n" +
		"Continue:\n" +
		"  qualitymd status\n" +
		"  qualitymd evaluation create\n\n" +
		"More:\n" +
		"  qualitymd --help\n" +
		"  docs: https://getquality.md\n" +
		"  report issues: https://github.com/qualitymd/quality.md/issues\n"
	if out.String() != want {
		t.Fatalf("stdout =\n%q\nwant\n%q", out.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecuteHelpShowsFullReference(t *testing.T) {
	var out, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"--help"})

	if got := execute(context.Background(), cmd); got != ExitOK {
		t.Fatalf("execute() = %d, want %d", got, ExitOK)
	}

	got := out.String()
	for _, want := range []string{
		"USAGE",
		"COMMON TASKS",
		"MANAGE",
		"qualitymd [command] [--flags]",
		"npx skills add qualitymd/quality.md",
		"https://getquality.md",
		"https://github.com/qualitymd/quality.md/issues",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("stdout = %q, want substring %q", got, want)
		}
	}
	if strings.Contains(got, "Start:\n  qualitymd init") {
		t.Fatalf("stdout = %q, want full help instead of root welcome", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecuteMapsLintProblemsWithoutExtraStderr(t *testing.T) {
	path := writeLintModel(t, `---
title: Example
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
requirements:
  has-assessment:
    title: Has an assessment
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

func TestBuildInfoPrefersLdflagsVersion(t *testing.T) {
	savedVersion, savedCommit := version, commit
	t.Cleanup(func() { version, commit = savedVersion, savedCommit })

	version, commit = "v1.2.3", "abc1234"
	if v, c := buildInfo(); v != "v1.2.3" || c != "abc1234" {
		t.Fatalf("buildInfo() = (%q, %q), want stamped (v1.2.3, abc1234)", v, c)
	}
}

func TestBuildInfoFallsBackForDevVersion(t *testing.T) {
	savedVersion, savedCommit := version, commit
	t.Cleanup(func() { version, commit = savedVersion, savedCommit })

	version, commit = "dev", "none"
	if v, _ := buildInfo(); v == "none" || v == "unknown (built from source)" {
		t.Fatalf("buildInfo() version = %q, want a dev or recovered version", v)
	}
}

func TestFallbackBuildInfoVersionUsesModuleVersion(t *testing.T) {
	v, c := fallbackBuildInfoVersion("v1.2.3", "abcdef123456")
	if v != "v1.2.3" || c != "abcdef1" {
		t.Fatalf("fallbackBuildInfoVersion() = (%q, %q), want module version and short commit", v, c)
	}
}

func TestFallbackBuildInfoVersionIncludesDevRevision(t *testing.T) {
	v, c := fallbackBuildInfoVersion("(devel)", "abcdef123456")
	if v != "dev (abcdef1)" || c != "" {
		t.Fatalf("fallbackBuildInfoVersion() = (%q, %q), want dev short revision and empty commit", v, c)
	}
}

func TestExecuteMapsSuccess(t *testing.T) {
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
