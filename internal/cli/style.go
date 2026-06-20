package cli

// Human-output styling. Every command's human surface draws from the one brand
// palette here so help, errors, and command output read as a single program.
// Styling is a terminal convenience only: it is applied solely when the target
// writer is a terminal and NO_COLOR is unset (see colorEnabled), so piped,
// redirected, and agent-driven output stays plain and byte-stable per the
// CLI spec's agent-accessibility baseline.

import (
	"fmt"
	"image/color"
	"io"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/charmbracelet/x/term"
)

// brandColorScheme keeps Fang-rendered help mostly monochrome. Color is
// reserved for status output below, not the command reference itself.
func brandColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	base := c(charmtone.Charcoal, charmtone.Ash)
	dim := c(charmtone.Squid, charmtone.Oyster)
	codeblock := c(charmtone.Salt, lipgloss.Color("#2F2E36"))
	return fang.ColorScheme{
		Base:           base,
		Title:          base,
		Description:    base,
		Codeblock:      codeblock,
		Program:        base,
		DimmedArgument: dim,
		Comment:        dim,
		Flag:           base,
		FlagDefault:    dim,
		Command:        base,
		QuotedString:   base,
		Argument:       base,
		Help:           base,
		Dash:           dim,
		ErrorHeader: [2]color.Color{
			charmtone.Butter,
			charmtone.Cherry,
		},
		ErrorDetails: charmtone.Cherry,
	}
}

// colorEnabled reports whether human styling should be applied to w. It mirrors
// the Fang / Lip Gloss non-TTY behavior: never style when NO_COLOR is set or
// when w is not a terminal, so redirected and piped output is plain.
func colorEnabled(w io.Writer) bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	fdWriter, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	return term.IsTerminal(fdWriter.Fd())
}

// Brand styles, used only on the color-enabled path. Colors are charmtone keys
// chosen to stay legible on both light and dark terminals.
var (
	styleBrand   = lipgloss.NewStyle().Bold(true)
	styleSuccess = lipgloss.NewStyle().Foreground(charmtone.Guac)
	styleError   = lipgloss.NewStyle().Foreground(charmtone.Cherry).Bold(true)
	styleWarning = lipgloss.NewStyle().Foreground(charmtone.Mustard).Bold(true)
	styleInfo    = lipgloss.NewStyle().Foreground(charmtone.Squid)
	styleRule    = lipgloss.NewStyle().Foreground(charmtone.Squid)
	styleDim     = lipgloss.NewStyle().Foreground(charmtone.Squid)
	styleCommand = lipgloss.NewStyle().Bold(true)
)

const (
	glyphSuccess = "✓"
	glyphError   = "✗"
	glyphWarning = "⚠"
	glyphInfo    = "•"
)

func renderRootWelcome(w io.Writer) error {
	if !colorEnabled(w) {
		_, err := fmt.Fprint(w, rootWelcomePlain())
		return err
	}
	_, err := fmt.Fprint(w, rootWelcomeStyled())
	return err
}

func rootWelcomePlain() string {
	return "QUALITY.md\n\n" +
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
}

func rootWelcomeStyled() string {
	return styleBrand.Render("QUALITY.md") + "\n\n" +
		"The companion CLI for the QUALITY.md file format for evaluating and improving\n" +
		"the quality of AI assistant projects and harnesses.\n\n" +
		"Designed to be used with the companion agent skill:\n" +
		"  " + styleCommand.Render("npx skills add qualitymd/quality.md") + "\n\n" +
		styleDim.Render("Start:") + "\n" +
		"  " + styleCommand.Render("qualitymd init") + "\n" +
		"  " + styleCommand.Render("qualitymd lint QUALITY.md") + "\n\n" +
		styleDim.Render("Continue:") + "\n" +
		"  " + styleCommand.Render("qualitymd status") + "\n" +
		"  " + styleCommand.Render("qualitymd evaluation create") + "\n\n" +
		styleDim.Render("More:") + "\n" +
		"  " + styleCommand.Render("qualitymd --help") + "\n" +
		"  " + styleDim.Render("docs:") + " https://getquality.md\n" +
		"  " + styleDim.Render("report issues:") + " https://github.com/qualitymd/quality.md/issues\n"
}
