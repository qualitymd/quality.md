package cli

import (
	"bytes"
	"io"
	"os"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/spf13/cobra"

	qualitymd "github.com/qualitymd/quality.md"
)

func newSchemaCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "schema",
		Short: "Emit the companion JSON Schema for QUALITY.md frontmatter",
		Example: "  qualitymd schema\n" +
			"  qualitymd schema > quality.schema.json",
		Args: usage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return writeSchema(cmd.OutOrStdout(), qualitymd.Schema())
		},
	}
}

// writeSchema mirrors writeMarkdown's two-branch shape. When output must be
// plain — w is not a terminal or NO_COLOR is set — it writes the verbatim JSON
// bytes and nothing else, so `qualitymd schema > quality.schema.json` reproduces
// the artifact byte-for-byte. On a terminal it syntax-highlights for readability
// and pages; highlighting only adds ANSI escapes around tokens and never runs on
// the bytes a redirect captures, so it cannot change the artifact.
func writeSchema(w io.Writer, jsonBytes []byte) error {
	if !colorEnabled(w) {
		_, err := w.Write(jsonBytes)
		return err
	}
	highlighted, err := highlightJSON(jsonBytes)
	if err != nil {
		// Highlighting is a convenience; fall back to the plain bytes rather
		// than failing the command over a rendering error.
		_, werr := w.Write(jsonBytes)
		return werr
	}
	return page(w, highlighted)
}

func highlightJSON(jsonBytes []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := quick.Highlight(&buf, string(jsonBytes), "json", terminalFormatterName(), schemaStyle); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// schemaStyle is the chroma style used to highlight emitted JSON. It is isolated
// here as the single place to swap the palette; a brand-exact custom style is a
// possible future refinement.
const schemaStyle = "monokai"

// terminalFormatterName selects the chroma terminal formatter, preferring
// truecolor when the terminal advertises it and falling back to the 256-color
// formatter otherwise.
func terminalFormatterName() string {
	switch os.Getenv("COLORTERM") {
	case "truecolor", "24bit":
		return "terminal16m"
	default:
		return "terminal256"
	}
}
