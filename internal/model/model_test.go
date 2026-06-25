package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/qualitymd/quality.md/internal/document"
)

func TestDecodeCurrentSchema(t *testing.T) {
	path := writeModel(t, `---
title: Example
ratingScale:
  - level: target
    title: Target
    criterion: Meets the requirement.
  - level: unacceptable
    title: Unacceptable
    criterion: Does not meet the requirement.
areas:
  api:
    title: API
    source: ./internal/api
    factors:
      reliability:
        title: Reliability
        description: The API continues to behave under expected failure modes.
        requirements:
          durable-writes:
            title: Writes are durable
            assessment: Inspect the write path and tests for durable commits.
---

# Example quality model
`)

	doc, err := document.Parse(path)
	if err != nil {
		t.Fatalf("document.Parse() error = %v", err)
	}
	model, err := Decode(doc)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if len(model.RatingScale) != 2 {
		t.Fatalf("len(RatingScale) = %d, want 2", len(model.RatingScale))
	}
	if got := model.Areas["api"].Title; got != "API" {
		t.Fatalf("area title = %q, want API", got)
	}
	if got := model.Areas["api"].Factors["reliability"].Title; got != "Reliability" {
		t.Fatalf("factor title = %q, want Reliability", got)
	}
	requirement := model.Areas["api"].Factors["reliability"].Requirements["durable-writes"]
	if got := requirement.Title; got != "Writes are durable" {
		t.Fatalf("requirement title = %q, want Writes are durable", got)
	}
	if got := requirement.Assessment; got == "" {
		t.Fatal("assessment was not loaded")
	}
}

func writeModel(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "QUALITY.md")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}
	return path
}
