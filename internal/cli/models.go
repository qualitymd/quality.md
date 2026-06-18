package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/models"
)

func newModelsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Work with bundled QUALITY.md models",
	}
	cmd.AddCommand(newModelsListCmd())
	cmd.AddCommand(newModelsViewCmd())
	return cmd
}

func newModelsListCmd() *cobra.Command {
	var jsonOutput bool
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bundled QUALITY.md models",
		Example: "  qualitymd models list\n" +
			"  qualitymd models list --json",
		Args: usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			entries := models.Catalog()
			if jsonOutput {
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(entries)
			}
			return renderModelsList(cmd, entries)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit the bundled model catalog as JSON")
	return cmd
}

func newModelsViewCmd() *cobra.Command {
	var jsonOutput bool
	var source string
	cmd := &cobra.Command{
		Use:   "view <name>",
		Short: "Emit a bundled QUALITY.md model",
		Example: "  qualitymd models view quality-meta-model\n" +
			"  qualitymd models view quality-meta-model --source QUALITY.md > model.md\n" +
			"  qualitymd models view quality-meta-model --source QUALITY.md --json",
		Args: usage(cobra.ExactArgs(1)),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if jsonOutput {
				view, err := models.Structured(name, source)
				if err != nil {
					return usageError(err)
				}
				encoder := json.NewEncoder(cmd.OutOrStdout())
				encoder.SetIndent("", "  ")
				return encoder.Encode(view)
			}
			markdown, err := models.Markdown(name, source)
			if err != nil {
				return usageError(err)
			}
			return writeMarkdown(cmd.OutOrStdout(), markdown)
		},
	}
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit a structured JSON representation")
	cmd.Flags().StringVar(&source, "source", "", "rewrite the model's apex source before emitting it")
	return cmd
}

func renderModelsList(cmd *cobra.Command, entries []models.Entry) error {
	out := cmd.OutOrStdout()
	if colorEnabled(out) {
		return renderModelsListStyled(out, entries)
	}
	return renderModelsListPlain(out, entries)
}

func renderModelsListPlain(out io.Writer, entries []models.Entry) error {
	table := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	if _, err := fmt.Fprintln(table, "NAME\tTITLE\tDESCRIPTION"); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := fmt.Fprintf(table, "%s\t%s\t%s\n", entry.Name, entry.Title, entry.Description); err != nil {
			return err
		}
	}
	return table.Flush()
}

func renderModelsListStyled(out io.Writer, entries []models.Entry) error {
	nameWidth, titleWidth := len("NAME"), len("TITLE")
	for _, entry := range entries {
		nameWidth = max(nameWidth, lipgloss.Width(entry.Name))
		titleWidth = max(titleWidth, lipgloss.Width(entry.Title))
	}

	if _, err := fmt.Fprintf(out, "%s  %s  %s\n",
		styleHeader.Render(padRight("NAME", nameWidth)),
		styleHeader.Render(padRight("TITLE", titleWidth)),
		styleHeader.Render("DESCRIPTION"),
	); err != nil {
		return err
	}
	for _, entry := range entries {
		if _, err := fmt.Fprintf(out, "%s  %s  %s\n",
			styleCommand.Render(padRight(entry.Name, nameWidth)),
			padRight(entry.Title, titleWidth),
			entry.Description,
		); err != nil {
			return err
		}
	}
	return nil
}

func padRight(s string, width int) string {
	padding := width - lipgloss.Width(s)
	if padding <= 0 {
		return s
	}
	return s + strings.Repeat(" ", padding)
}
