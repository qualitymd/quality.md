package cli

import (
	"encoding/json"
	"fmt"
	"text/tabwriter"

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
	table := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
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
