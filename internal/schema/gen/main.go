// Command gen writes the committed companion JSON Schema (quality.schema.json)
// from the structural schema in internal/schema. It is the entrypoint for the
// `go generate ./...` directive in the repo-root schema.go; the committed file
// it writes is embedded into the binary and emitted verbatim by
// `qualitymd schema`.
package main

import (
	"log"
	"os"

	"github.com/qualitymd/quality.md/internal/schema"
)

const outputPath = "quality.schema.json"

func main() {
	data, err := schema.GenerateJSON()
	if err != nil {
		log.Fatalf("generating %s: %v", outputPath, err)
	}
	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		log.Fatalf("writing %s: %v", outputPath, err)
	}
}
