package cli

// Human-output styling. Every command's human surface draws from the one brand
// palette here so help, errors, and command output read as a single program.
// Styling is a terminal convenience only: it is applied solely when the target
// writer is a terminal and NO_COLOR is unset (see colorEnabled), so piped,
// redirected, and agent-driven output stays plain and byte-stable per the
// CLI spec's agent-accessibility baseline.

import (
	"io"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/exp/charmtone"
	"github.com/charmbracelet/x/term"
)

// brandColorScheme is qualitymd's Fang colorscheme: the stack default with the
// title and program name in the brand green, so styled help and version share
// the palette the command renderers below use.
func brandColorScheme(c lipgloss.LightDarkFunc) fang.ColorScheme {
	scheme := fang.DefaultColorScheme(c)
	scheme.Title = charmtone.Guac
	scheme.Program = c(charmtone.Guac, charmtone.Julep)
	return scheme
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
	styleSuccess = lipgloss.NewStyle().Foreground(charmtone.Guac).Bold(true)
	styleError   = lipgloss.NewStyle().Foreground(charmtone.Cherry).Bold(true)
	styleWarning = lipgloss.NewStyle().Foreground(charmtone.Mustard).Bold(true)
	styleInfo    = lipgloss.NewStyle().Foreground(charmtone.Malibu).Bold(true)
	styleRule    = lipgloss.NewStyle().Foreground(charmtone.Squid)
	styleDim     = lipgloss.NewStyle().Foreground(charmtone.Squid)
	styleCommand = lipgloss.NewStyle().Foreground(charmtone.Malibu)
	styleHeader  = lipgloss.NewStyle().Foreground(charmtone.Squid).Bold(true)
)

const (
	glyphSuccess = "✓"
	glyphError   = "✗"
	glyphWarning = "⚠"
	glyphInfo    = "•"
)
