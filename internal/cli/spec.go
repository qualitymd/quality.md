package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"

	qualitymd "github.com/qualitymd/quality.md"
)

func newSpecCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "spec",
		Short: "Emit the QUALITY.md format specification",
		Example: "  qualitymd spec\n" +
			"  qualitymd spec | glow\n" +
			"  qualitymd spec > SPECIFICATION.md",
		Args: usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeSpec(cmd.OutOrStdout(), qualitymd.Specification())
		},
	}
}

func writeSpec(w io.Writer, markdown []byte) error {
	if !colorEnabled(w) {
		_, err := w.Write(markdown)
		return err
	}
	rendered, err := renderSpec(markdown)
	if err != nil {
		return err
	}
	return page(w, rendered)
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

// page sends content through the user's pager so a long rendered spec scrolls
// like `glow` or `git log`. It only applies when w is the real terminal file;
// if no pager is configured or it fails to start, the content is written
// directly, so paging is never load-bearing.
func page(w io.Writer, content []byte) error {
	tty, ok := w.(*os.File)
	if !ok {
		_, err := w.Write(content)
		return err
	}

	var pager *exec.Cmd
	if custom := os.Getenv("PAGER"); custom != "" {
		pager = exec.Command("sh", "-c", custom)
	} else if less, err := exec.LookPath("less"); err == nil {
		// -R passes the rendered ANSI through; -F quits if it fits one screen.
		pager = exec.Command(less, "-R", "-F")
	} else {
		_, err := w.Write(content)
		return err
	}

	pager.Stdin = bytes.NewReader(content)
	pager.Stdout = tty
	pager.Stderr = os.Stderr
	if err := pager.Start(); err != nil {
		_, werr := w.Write(content)
		return werr
	}
	return pager.Wait()
}
