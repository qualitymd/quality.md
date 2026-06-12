// Package report renders evaluation results to the terminal.
package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/lipgloss"

	"github.com/qualitymd/quality.md/internal/eval"
)

var (
	passStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	failStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	skipStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	factorStyle = lipgloss.NewStyle().Bold(true).Underline(true)
	detailStyle = lipgloss.NewStyle().Faint(true)
)

// Print writes a grouped, styled summary of results to w.
func Print(w io.Writer, results eval.Results) {
	byFactor := map[string][]eval.Result{}
	var factors []string
	for _, it := range results.Items {
		if _, seen := byFactor[it.Factor]; !seen {
			factors = append(factors, it.Factor)
		}
		byFactor[it.Factor] = append(byFactor[it.Factor], it)
	}
	sort.Strings(factors)

	for _, f := range factors {
		fmt.Fprintln(w, factorStyle.Render(f))
		for _, it := range byFactor[f] {
			fmt.Fprintf(w, "  %s  %s\n", icon(it.Status), it.Requirement)
			if it.Detail != "" {
				fmt.Fprintf(w, "       %s\n", detailStyle.Render(it.Detail))
			}
		}
		fmt.Fprintln(w)
	}
}

func icon(s eval.Status) string {
	switch s {
	case eval.StatusPass:
		return passStyle.Render("✓ pass")
	case eval.StatusFail:
		return failStyle.Render("✗ fail")
	default:
		return skipStyle.Render("• skip")
	}
}
