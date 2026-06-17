package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/qualitymd/quality.md/internal/scaffold"
)

func newInitCmd() *cobra.Command {
	var force bool
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Scaffold a starter QUALITY.md",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "QUALITY.md"
			if len(args) == 1 {
				path = args[0]
			}

			if path == "-" {
				_, err := cmd.OutOrStdout().Write(scaffold.Bytes())
				return err
			}

			if err := scaffold.Create(path, force); err != nil {
				return err
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Created %s\n\nNext: qualitymd lint %s\n", path, path)
			return nil
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "overwrite an existing file")
	return cmd
}
