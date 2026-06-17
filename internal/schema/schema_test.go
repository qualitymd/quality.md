package schema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpecificationSchemaSnippetsMatchDeclaration(t *testing.T) {
	spec := readSpecification(t)

	for _, tc := range []struct {
		heading string
		node    Node
	}{
		{"Model", Model},
		{"Target", Target},
		{"Factor", Factor},
		{"Requirement", Requirement},
	} {
		t.Run(tc.heading, func(t *testing.T) {
			got := parseTopLevelProperties(firstYAMLBlock(t, spec, "#### "+tc.heading))
			assertPropertiesMatch(t, got, tc.node)
		})
	}

	t.Run("RatingLevel", func(t *testing.T) {
		modelSnippet := firstYAMLBlock(t, spec, "#### Model")
		got := parseTopLevelProperties(ratingLevelSnippet(t, modelSnippet))
		assertPropertiesMatch(t, got, RatingLevel)
	})

	ratingScale, ok := Model.Property(PropertyRatingScale)
	if !ok {
		t.Fatal("model schema has no ratingScale property")
	}
	if ratingScale.MinItems > 0 && !strings.Contains(spec, "At least two rating levels MUST be supplied.") {
		t.Fatalf("SPECIFICATION.md does not document ratingScale MinItems = %d", ratingScale.MinItems)
	}
	for _, group := range Model.RequiredAny {
		if group.Name != "model-content" {
			continue
		}
		for _, property := range group.Properties {
			if !strings.Contains(spec, property) {
				t.Fatalf("SPECIFICATION.md does not mention required-any property %q", property)
			}
		}
		if !strings.Contains(spec, "An entry on either factors, requirements, or targets MUST be supplied.") {
			t.Fatal("SPECIFICATION.md does not document the model-content required-any group")
		}
	}
}

func readSpecification(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "SPECIFICATION.md"))
	if err != nil {
		t.Fatalf("os.ReadFile(SPECIFICATION.md) error = %v", err)
	}
	return string(data)
}

func firstYAMLBlock(t *testing.T, markdown, heading string) string {
	t.Helper()
	sectionStart := strings.Index(markdown, heading)
	if sectionStart == -1 {
		t.Fatalf("heading %q not found", heading)
	}
	fenceStart := strings.Index(markdown[sectionStart:], "```yaml")
	if fenceStart == -1 {
		t.Fatalf("yaml block after %q not found", heading)
	}
	blockStart := sectionStart + fenceStart + len("```yaml")
	blockEnd := strings.Index(markdown[blockStart:], "```")
	if blockEnd == -1 {
		t.Fatalf("yaml block after %q is not closed", heading)
	}
	return markdown[blockStart : blockStart+blockEnd]
}

func ratingLevelSnippet(t *testing.T, modelSnippet string) string {
	t.Helper()
	lines := strings.Split(modelSnippet, "\n")
	var out []string
	inRatingScale := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, PropertyRatingScale+":") {
			inRatingScale = true
			continue
		}
		if !inRatingScale {
			continue
		}
		if line != "" && line[0] != ' ' {
			break
		}
		out = append(out, strings.TrimSpace(line))
	}
	if len(out) == 0 {
		t.Fatal("model snippet contains no ratingScale item")
	}
	return strings.Join(out, "\n")
}

func parseTopLevelProperties(snippet string) map[string]Presence {
	out := map[string]Presence{}
	for _, line := range strings.Split(snippet, "\n") {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "#") {
			continue
		}
		if line[0] == ' ' {
			continue
		}
		key, presence, ok := parsePropertyLine(line)
		if ok {
			out[key] = presence
		}
	}
	return out
}

func parsePropertyLine(line string) (string, Presence, bool) {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "- ")
	if strings.HasPrefix(line, "<") {
		return "", "", false
	}
	parts := strings.SplitN(line, "#", 2)
	field := strings.TrimSpace(parts[0])
	colon := strings.Index(field, ":")
	if colon == -1 {
		return "", "", false
	}
	key := strings.TrimSpace(field[:colon])
	if key == "" {
		return "", "", false
	}
	if len(parts) == 1 {
		return key, "", true
	}
	comment := strings.TrimSpace(parts[1])
	presenceText := comment
	if before, _, ok := strings.Cut(presenceText, ";"); ok {
		presenceText = before
	}
	presenceText = strings.TrimSuffix(strings.TrimSpace(presenceText), "*")
	switch presenceText {
	case "Required":
		return key, RequiredPresence, true
	case "Recommended":
		return key, RecommendedPresence, true
	case "Optional":
		return key, OptionalPresence, true
	default:
		return key, "", true
	}
}

func assertPropertiesMatch(t *testing.T, got map[string]Presence, node Node) {
	t.Helper()
	want := map[string]Presence{}
	for _, property := range node.Properties {
		want[property.Name] = property.Presence
	}
	for name, presence := range want {
		gotPresence, ok := got[name]
		if !ok {
			t.Fatalf("%s missing from SPECIFICATION.md snippet; got keys %#v", name, got)
		}
		if gotPresence != presence {
			t.Fatalf("%s presence = %q, want %q", name, gotPresence, presence)
		}
	}
	for name := range got {
		if _, ok := want[name]; !ok {
			t.Fatalf("SPECIFICATION.md snippet has extra property %q; schema keys %#v", name, want)
		}
	}
}
