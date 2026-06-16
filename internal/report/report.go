// Package report renders placeholder assessment traversal results to the terminal.
package report

import (
	"fmt"
	"io"
	"sort"

	"github.com/charmbracelet/lipgloss"

	"github.com/qualitymd/quality.md/internal/eval"
)

var (
	notAssessedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	groupStyle       = lipgloss.NewStyle().Bold(true).Underline(true)
	detailStyle      = lipgloss.NewStyle().Faint(true)
)

// Print writes a grouped, styled summary of results to w.
func Print(w io.Writer, results eval.Results) {
	byGroup := map[string][]eval.Result{}
	var groups []string
	for _, it := range results.Items {
		group := it.Target
		if it.Factor != "" {
			group += " / " + it.Factor
		}
		if _, seen := byGroup[group]; !seen {
			groups = append(groups, group)
		}
		byGroup[group] = append(byGroup[group], it)
	}
	sort.Strings(groups)

	for _, group := range groups {
		fmt.Fprintln(w, groupStyle.Render(group))
		for _, it := range byGroup[group] {
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
	case eval.StatusRated:
		return "rated"
	default:
		return notAssessedStyle.Render("not assessed")
	}
}
