package qualitymd_test

import (
	"bytes"
	"testing"

	qualitymd "github.com/qualitymd/quality.md"
	"github.com/qualitymd/quality.md/internal/schema"
)

// TestSchemaMatchesGenerator guards the no-drift property: the committed,
// embedded quality.schema.json must equal a fresh render of the structural
// schema. If this fails, the artifact is stale — run `go generate ./...`.
func TestSchemaMatchesGenerator(t *testing.T) {
	want, err := schema.GenerateJSON()
	if err != nil {
		t.Fatalf("schema.GenerateJSON() error = %v", err)
	}
	if !bytes.Equal(qualitymd.Schema(), want) {
		t.Fatal("bundled quality.schema.json is stale; run `go generate ./...`")
	}
}
