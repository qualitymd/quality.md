package cli

import (
	"io"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"

	qualitymd "github.com/qualitymd/quality.md"
)

func newSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec",
		Short: "Emit the QUALITY.md format specification",
		Args:  usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeSpec(cmd.OutOrStdout(), qualitymd.Specification())
		},
	}
}

func writeSpec(w io.Writer, markdown []byte) error {
	if shouldRenderSpec(w) {
		rendered, err := renderSpec(markdown)
		if err != nil {
			return err
		}
		_, err = w.Write(rendered)
		return err
	}
	_, err := w.Write(markdown)
	return err
}

func shouldRenderSpec(w io.Writer) bool {
	if _, ok := os.LookupEnv("NO_COLOR"); ok {
		return false
	}
	fdWriter, ok := w.(interface{ Fd() uintptr })
	if !ok {
		return false
	}
	return term.IsTerminal(fdWriter.Fd())
}

func renderSpec(markdown []byte) ([]byte, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return nil, err
	}
	return renderer.RenderBytes(markdown)
}
