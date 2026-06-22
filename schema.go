package qualitymd

import _ "embed"

//go:generate go run ./internal/schema/gen

//go:embed quality.schema.json
var schemaJSON []byte

// Schema returns the bundled companion JSON Schema for QUALITY.md frontmatter.
// The schema is structural-only and non-normative; SPECIFICATION.md remains the
// normative source of truth.
func Schema() []byte {
	return append([]byte(nil), schemaJSON...)
}
