package qualitymd

import _ "embed"

//go:embed SPECIFICATION.md
var specification []byte

// Specification returns the bundled QUALITY.md format specification.
func Specification() []byte {
	return append([]byte(nil), specification...)
}
