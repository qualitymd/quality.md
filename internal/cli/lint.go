package cli

import (
	"fmt"
	"io"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/lint"
)

func newLintCmd() *cobra.Command {
	var jsonOutput bool
	var fix bool
	cmd := &cobra.Command{
		Use:   "lint [path]",
		Short: "Validate a QUALITY.md file",
		Example: "  qualitymd lint\n" +
			"  qualitymd lint docs/QUALITY.md\n" +
			"  qualitymd lint --fix\n" +
			"  qualitymd lint --json",
		Args: usage(cobra.MaximumNArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "QUALITY.md"
			if len(args) == 1 {
				path = args[0]
			}
			if path == "-" {
				return fmt.Errorf("lint does not read from stdin yet; pass a file path")
			}

			var (
				result lint.Result
				err    error
			)
			if fix {
				result, err = lint.Fix(path)
			} else {
				result, err = lint.Check(path)
			}
			if err != nil {
				return err
			}

			if jsonOutput {
				data, err := result.JSON()
				if err != nil {
					return err
				}
				if _, err := cmd.OutOrStdout().Write(append(data, '\n')); err != nil {
					return err
				}
			} else if err := renderLintHuman(cmd, result); err != nil {
				return err
			}
			if err := result.Err(); err != nil {
				return silentProblems(err)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a machine-readable JSON lint result")
	cmd.Flags().BoolVar(&fix, "fix", false, "apply deterministic in-place repairs")
	return cmd
}

func renderLintHuman(cmd *cobra.Command, result lint.Result) error {
	out := cmd.OutOrStdout()
	if !colorEnabled(out) {
		return renderLintPlain(out, result)
	}
	return renderLintStyled(out, result)
}

// renderLintPlain is the non-terminal rendering: no color or glyphs, stable
// bytes. It is the canonical form; the styled renderer only dresses it up.
func renderLintPlain(out io.Writer, result lint.Result) error {
	if result.Summary.Fixed > 0 {
		if _, err := fmt.Fprintf(out, "Applied %d repair(s).\n", result.Summary.Fixed); err != nil {
			return err
		}
	}
	if len(result.Findings) == 0 {
		_, err := fmt.Fprintf(out, "%s is valid.\n", result.Path)
		return err
	}
	for _, finding := range result.Findings {
		if _, err := fmt.Fprintf(out, "%s %s: %s (%s)\n", finding.Severity, finding.RuleID, finding.Message, finding.Location.Label); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(out, "\n%d error(s), %d warning(s).\n", result.Summary.Errors, result.Summary.Warnings)
	return err
}

// renderLintStyled is the terminal rendering: a severity glyph and color per
// finding, a clickable file:line location where the finding has a source
// position, and a colored count summary. It carries the same facts as the
// plain form.
func renderLintStyled(out io.Writer, result lint.Result) error {
	if result.Summary.Fixed > 0 {
		if _, err := fmt.Fprintf(out, "%s Applied %d repair(s).\n", styleSuccess.Render(glyphSuccess), result.Summary.Fixed); err != nil {
			return err
		}
	}
	if len(result.Findings) == 0 {
		_, err := fmt.Fprintf(out, "%s %s is valid.\n", styleSuccess.Render(glyphSuccess), result.Path)
		return err
	}
	for _, finding := range result.Findings {
		glyph, sev := severityStyle(finding.Severity)
		line := fmt.Sprintf("%s %s %s  %s",
			glyph,
			sev.Render(string(finding.Severity)),
			styleRule.Render(string(finding.RuleID)),
			finding.Message,
		)
		if loc := findingLocation(finding.Location); loc != "" {
			line += " " + styleDim.Render("("+loc+")")
		}
		if _, err := fmt.Fprintln(out, line); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(out, "\n%s\n", lintSummary(result.Summary))
	return err
}

// severityStyle maps a finding severity to its glyph and color.
func severityStyle(severity lint.Severity) (string, lipgloss.Style) {
	switch severity {
	case lint.SeverityError:
		return styleError.Render(glyphError), styleError
	case lint.SeverityWarning:
		return styleWarning.Render(glyphWarning), styleWarning
	default:
		return styleInfo.Render(glyphInfo), styleInfo
	}
}

// findingLocation renders the model-path label, appending a clickable
// file:line when the finding carries a source position. The label is always
// present so styled output conveys the same location as the plain form.
func findingLocation(loc lint.Location) string {
	if loc.Line > 0 {
		return fmt.Sprintf("%s, %s:%d", loc.Label, loc.Path, loc.Line)
	}
	return loc.Label
}

// lintSummary renders the colored count footer, naming only the non-zero
// severities and pluralizing each.
func lintSummary(summary lint.Summary) string {
	parts := []string{
		styleError.Render(count(summary.Errors, "error")),
		styleWarning.Render(count(summary.Warnings, "warning")),
	}
	if summary.Info > 0 {
		parts = append(parts, styleInfo.Render(count(summary.Info, "note")))
	}
	return strings.Join(parts, styleDim.Render(" · "))
}

func count(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
