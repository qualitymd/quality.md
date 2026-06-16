package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadCurrentSchema(t *testing.T) {
	path := writeModel(t, `---
title: Example
ratingScale:
  - level: target
    criterion: Meets the requirement.
  - level: unacceptable
    criterion: Does not meet the requirement.
targets:
  api:
    source: ./internal/api
    factors:
      reliability:
        description: The API continues to behave under expected failure modes.
        requirements:
          "writes are durable":
            assessment: Inspect the write path and tests for durable commits.
---

# Example quality model
`)

	model, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(model.RatingScale) != 2 {
		t.Fatalf("len(RatingScale) = %d, want 2", len(model.RatingScale))
	}
	if got := model.Targets["api"].Factors["reliability"].Requirements["writes are durable"].Assessment; got == "" {
		t.Fatal("assessment was not loaded")
	}
}

func TestLoadRejectsRootRatings(t *testing.T) {
	path := writeModel(t, `---
ratings:
  - level: target
    criterion: Meets the requirement.
  - level: unacceptable
    criterion: Does not meet the requirement.
requirements:
  "has an assessment":
    assessment: Inspect it.
---
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "field ratings not found") {
		t.Fatalf("Load() error = %v, want unknown ratings field", err)
	}
}

func TestLoadRejectsOldPromptRequirement(t *testing.T) {
	path := writeModel(t, `---
ratingScale:
  - level: target
    criterion: Meets the requirement.
  - level: unacceptable
    criterion: Does not meet the requirement.
requirements:
  "has an assessment":
    prompt: Inspect it.
---
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "field prompt not found") {
		t.Fatalf("Load() error = %v, want unknown prompt field", err)
	}
}

func TestLoadRejectsRatingScaleBelowRoot(t *testing.T) {
	path := writeModel(t, `---
ratingScale:
  - level: target
    criterion: Meets the requirement.
  - level: unacceptable
    criterion: Does not meet the requirement.
targets:
  api:
    ratingScale:
      - level: target
        criterion: Meets the requirement.
    requirements:
      "has an assessment":
        assessment: Inspect it.
---
`)

	_, err := Load(path)
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "field ratingScale not found") {
		t.Fatalf("Load() error = %v, want unknown nested ratingScale field", err)
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
