package cli

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/exp/charmtone"

	"github.com/qualitymd/quality.md/internal/lint"
	"github.com/qualitymd/quality.md/internal/receipt"
)

func TestColorEnabledOffForNonTerminalWriter(t *testing.T) {
	if colorEnabled(&bytes.Buffer{}) {
		t.Fatal("colorEnabled(buffer) = true, want false for a non-terminal writer")
	}
}

func TestColorEnabledOffWhenNoColorSet(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if colorEnabled(panicFDWriter{}) {
		t.Fatal("colorEnabled() = true, want false when NO_COLOR is set")
	}
}

func TestBrandColorSchemeAvoidsDefaultAccentColors(t *testing.T) {
	scheme := brandColorScheme(lipgloss.LightDark(false))
	for name, got := range map[string]any{
		"title":         scheme.Title,
		"program":       scheme.Program,
		"command":       scheme.Command,
		"flag":          scheme.Flag,
		"quoted string": scheme.QuotedString,
	} {
		for _, unwanted := range []any{charmtone.Charple, charmtone.Pony, charmtone.Cheeky, charmtone.Guac} {
			if reflect.DeepEqual(got, unwanted) {
				t.Fatalf("%s color = %#v, want neutral color", name, got)
			}
		}
	}
}

func TestRenderLintPlainIsStable(t *testing.T) {
	result := lint.Result{
		Path: "QUALITY.md",
		Summary: lint.Summary{
			Errors:   1,
			Warnings: 1,
			Fixed:    2,
		},
		Findings: []lint.Finding{
			{
				RuleID:   lint.RuleTooFewLevels,
				Severity: lint.SeverityError,
				Message:  "The rating scale has fewer than two levels.",
				Location: lint.Location{Label: "ratingScale", Line: 4},
			},
			{
				RuleID:   lint.RuleMissingTitle,
				Severity: lint.SeverityWarning,
				Message:  "The model root declares no title.",
				Location: lint.Location{Label: "frontmatter"},
			},
		},
	}

	var out bytes.Buffer
	if err := renderLintPlain(&out, result); err != nil {
		t.Fatalf("renderLintPlain() error = %v", err)
	}
	want := "Applied 2 repair(s).\n" +
		"error too-few-levels: The rating scale has fewer than two levels. (ratingScale)\n" +
		"warning missing-title: The model root declares no title. (frontmatter)\n" +
		"\n1 error(s), 1 warning(s).\n"
	if out.String() != want {
		t.Fatalf("renderLintPlain() =\n%q\nwant\n%q", out.String(), want)
	}
}

func TestRenderLintStyledCarriesGlyphsAndCounts(t *testing.T) {
	result := lint.Result{
		Path: "QUALITY.md",
		Summary: lint.Summary{
			Errors:   2,
			Warnings: 1,
		},
		Findings: []lint.Finding{
			{
				RuleID:   lint.RuleTooFewLevels,
				Severity: lint.SeverityError,
				Message:  "fewer than two levels",
				Location: lint.Location{Path: "QUALITY.md", Label: "ratingScale", Line: 4},
			},
		},
	}

	var out bytes.Buffer
	if err := renderLintStyled(&out, result); err != nil {
		t.Fatalf("renderLintStyled() error = %v", err)
	}
	got := out.String()
	for _, want := range []string{glyphError, "too-few-levels", "QUALITY.md:4", "2 errors", "1 warning"} {
		if !strings.Contains(got, want) {
			t.Fatalf("renderLintStyled() = %q, want substring %q", got, want)
		}
	}
}

func TestRenderLintStyledValidShowsSuccessGlyph(t *testing.T) {
	var out bytes.Buffer
	if err := renderLintStyled(&out, lint.Result{Path: "QUALITY.md"}); err != nil {
		t.Fatalf("renderLintStyled() error = %v", err)
	}
	got := out.String()
	if !strings.Contains(got, glyphSuccess) || !strings.Contains(got, "is valid.") {
		t.Fatalf("renderLintStyled() = %q, want success glyph and valid message", got)
	}
}

func TestRenderInitHumanPlainIsStable(t *testing.T) {
	var out bytes.Buffer
	if err := renderInitHuman(&out, "QUALITY.md", nil, "qualitymd lint QUALITY.md"); err != nil {
		t.Fatalf("renderInitHuman() error = %v", err)
	}
	want := "Created QUALITY.md\n\nNext: qualitymd lint QUALITY.md\n"
	if out.String() != want {
		t.Fatalf("renderInitHuman() = %q, want %q", out.String(), want)
	}
}

func TestRenderInitHumanPlainIncludesAgentInstructions(t *testing.T) {
	var out bytes.Buffer
	if err := renderInitHuman(&out, "QUALITY.md", []string{"AGENTS.md"}, "qualitymd lint QUALITY.md"); err != nil {
		t.Fatalf("renderInitHuman() error = %v", err)
	}
	want := "Created QUALITY.md\nAgent instructions: AGENTS.md\n\nNext: qualitymd lint QUALITY.md\n"
	if out.String() != want {
		t.Fatalf("renderInitHuman() = %q, want %q", out.String(), want)
	}
}

func TestRenderNextActionsPlainIsStable(t *testing.T) {
	var out bytes.Buffer
	actions := []receipt.Action{{ID: "rerun-lint", Label: "Re-run validation", Command: "qualitymd lint QUALITY.md"}}
	if err := renderNextActions(&out, actions); err != nil {
		t.Fatalf("renderNextActions() error = %v", err)
	}
	want := "\nNext: qualitymd lint QUALITY.md\n"
	if out.String() != want {
		t.Fatalf("renderNextActions() = %q, want %q", out.String(), want)
	}
}

func TestLintSummaryPluralizesAndOmitsZeroNotes(t *testing.T) {
	got := lintSummary(lint.Summary{Errors: 1, Warnings: 2})
	for _, want := range []string{"1 error", "2 warnings"} {
		if !strings.Contains(got, want) {
			t.Fatalf("lintSummary() = %q, want substring %q", got, want)
		}
	}
	if strings.Contains(got, "note") {
		t.Fatalf("lintSummary() = %q, want no note segment when Info is zero", got)
	}
	if withNotes := lintSummary(lint.Summary{Info: 3}); !strings.Contains(withNotes, "3 notes") {
		t.Fatalf("lintSummary() = %q, want 3 notes", withNotes)
	}
}

func hasTerminalEscape(s string) bool {
	return strings.Contains(s, "\x1b[") || strings.Contains(s, "\x1b]")
}
