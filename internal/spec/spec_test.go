package spec

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseAndDecodeCurrentSchema(t *testing.T) {
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

	doc, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	model, err := Decode(doc)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if len(model.RatingScale) != 2 {
		t.Fatalf("len(RatingScale) = %d, want 2", len(model.RatingScale))
	}
	if got := model.Targets["api"].Factors["reliability"].Requirements["writes are durable"].Assessment; got == "" {
		t.Fatal("assessment was not loaded")
	}
}

func TestParseRequiresFrontmatterFence(t *testing.T) {
	path := writeModel(t, `title: Example
ratingScale: []
`)

	_, err := Parse(path)
	if err == nil {
		t.Fatal("Parse() error = nil, want error")
	}
	var parseErr *ParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("Parse() error = %T, want ParseError", err)
	}
	if !strings.Contains(err.Error(), "frontmatter") {
		t.Fatalf("Parse() error = %v, want frontmatter message", err)
	}
}

func TestRenderPreservesMarkdownBody(t *testing.T) {
	path := writeModel(t, `---
title: Example
ratingScale:
  - level: target
    criterion: Meets the requirement.
  - level: unacceptable
    criterion: Does not meet the requirement.
requirements:
  "has an assessment":
    assessment: Inspect it.
---

# Example

Keep this body exactly.
`)

	doc, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	rendered, err := Render(doc)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(string(rendered), "---\n\n# Example\n\nKeep this body exactly.\n") {
		t.Fatalf("rendered body changed:\n%s", rendered)
	}
}

func TestMapEntries(t *testing.T) {
	mapping := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			{Kind: yaml.ScalarNode, Value: "first"},
			{Kind: yaml.ScalarNode, Value: "1"},
			{Kind: yaml.ScalarNode, Value: "second"},
			{Kind: yaml.ScalarNode, Value: "2"},
		},
	}

	var keys []string
	for key, value := range MapEntries(mapping) {
		keys = append(keys, key.Value+"="+value.Value)
	}
	if want := []string{"first=1", "second=2"}; !reflect.DeepEqual(keys, want) {
		t.Fatalf("MapEntries() = %v, want %v", keys, want)
	}

	count := 0
	for range MapEntries(nil) {
		count++
	}
	if count != 0 {
		t.Fatalf("MapEntries(nil) yielded %d pairs, want 0", count)
	}

	scalar := &yaml.Node{Kind: yaml.ScalarNode, Value: "not a mapping"}
	count = 0
	for range MapEntries(scalar) {
		count++
	}
	if count != 0 {
		t.Fatalf("MapEntries(scalar) yielded %d pairs, want 0", count)
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
