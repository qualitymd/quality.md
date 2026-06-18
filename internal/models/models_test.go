package models

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestCatalogIncludesQualityMetaModel(t *testing.T) {
	got := Catalog()
	if len(got) != 1 {
		t.Fatalf("len(Catalog()) = %d, want 1", len(got))
	}
	if got[0].Name != "quality-meta-model" {
		t.Fatalf("Catalog()[0].Name = %q, want quality-meta-model", got[0].Name)
	}
	if got[0].Title == "" || got[0].Description == "" {
		t.Fatalf("Catalog()[0] = %#v, want title and description", got[0])
	}
}

func TestMarkdownWithoutSourceRewriteIsAuthoredBytes(t *testing.T) {
	want, err := os.ReadFile("quality-meta-model.md")
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	got, err := Markdown("quality-meta-model", "")
	if err != nil {
		t.Fatalf("Markdown() error = %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("Markdown() did not return authored bytes")
	}
}

func TestMarkdownRewritesSource(t *testing.T) {
	got, err := Markdown("quality-meta-model", "docs/QUALITY.md")
	if err != nil {
		t.Fatalf("Markdown() error = %v", err)
	}
	if !bytes.Contains(got, []byte("source: docs/QUALITY.md\n")) {
		t.Fatalf("Markdown() = %s, want rewritten source", got[:120])
	}
	if bytes.Contains(got, []byte("source: ./QUALITY.md\n")) {
		t.Fatal("Markdown() still contains authored source")
	}
}

func TestStructuredViewRewritesSourceAndOmitsPathFromJSON(t *testing.T) {
	view, err := Structured("quality-meta-model", "QUALITY.md")
	if err != nil {
		t.Fatalf("Structured() error = %v", err)
	}
	if view.SchemaVersion != 1 {
		t.Fatalf("SchemaVersion = %d, want 1", view.SchemaVersion)
	}
	if view.Model.Source != "QUALITY.md" {
		t.Fatalf("Model.Source = %q, want QUALITY.md", view.Model.Source)
	}
	if view.BodyMarkdown == "" || view.BodyMarkdown[0] != '\n' {
		t.Fatalf("BodyMarkdown = %q, want original body with leading newline", view.BodyMarkdown[:min(len(view.BodyMarkdown), 20)])
	}
	data, err := json.Marshal(view)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if bytes.Contains(data, []byte("Path")) {
		t.Fatalf("JSON contains internal Path field: %s", data)
	}
}

func TestUnknownModelErrors(t *testing.T) {
	if _, err := Markdown("missing", ""); err == nil {
		t.Fatal("Markdown() error = nil, want unknown model error")
	}
	if _, err := Structured("missing", ""); err == nil {
		t.Fatal("Structured() error = nil, want unknown model error")
	}
}
