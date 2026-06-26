package model

import "testing"

// projectionSpec returns a model exercising areas, nested factors, and
// requirements declared both directly under an area and under a factor.
func projectionSpec() *Spec {
	return &Spec{
		Title: "Example",
		Factors: map[string]Factor{
			"security": {
				Title: "Security",
				Requirements: map[string]Requirement{
					"no-committed-secrets": {Title: "No credentials are committed"},
				},
				Factors: map[string]Factor{
					"secrets": {Title: "Secrets"},
				},
			},
		},
		Requirements: map[string]Requirement{
			"has-readme": {Title: "Has a README"},
		},
		Areas: map[string]Area{
			"webhooks": {
				Title: "Webhooks",
				Factors: map[string]Factor{
					"reliability": {Title: "Reliability"},
				},
			},
		},
	}
}

func TestProjectStructureAndIDs(t *testing.T) {
	root := Project(projectionSpec())
	if root.ID != "area:root" || root.Kind != KindArea || root.Label != "Example" {
		t.Fatalf("root = %+v", root)
	}
	if root.ParentID != "" {
		t.Fatalf("root parent = %q, want empty", root.ParentID)
	}

	byID := map[string]*Element{}
	for _, e := range Flatten(root) {
		byID[e.ID] = e
	}
	cases := []struct {
		id       string
		kind     Kind
		parentID string
	}{
		{"area:root", KindArea, ""},
		{"factor:root::security", KindFactor, "area:root"},
		{"factor:root::security/secrets", KindFactor, "factor:root::security"},
		{"requirement:root::no-committed-secrets", KindRequirement, "factor:root::security"},
		{"requirement:root::has-readme", KindRequirement, "area:root"},
		{"area:webhooks", KindArea, "area:root"},
		{"factor:webhooks::reliability", KindFactor, "area:webhooks"},
	}
	for _, tc := range cases {
		got, ok := byID[tc.id]
		if !ok {
			t.Fatalf("missing element %q", tc.id)
		}
		if got.Kind != tc.kind {
			t.Errorf("%s kind = %q, want %q", tc.id, got.Kind, tc.kind)
		}
		if got.ParentID != tc.parentID {
			t.Errorf("%s parentId = %q, want %q", tc.id, got.ParentID, tc.parentID)
		}
	}
}

func TestProjectDeterministicOrder(t *testing.T) {
	spec := projectionSpec()
	first := flatIDs(Project(spec))
	for i := 0; i < 5; i++ {
		if got := flatIDs(Project(spec)); !equalStrings(got, first) {
			t.Fatalf("projection order not stable:\n%v\n%v", first, got)
		}
	}
	// Factors precede requirements precede child areas at the root.
	want := []string{
		"area:root",
		"factor:root::security",
		"factor:root::security/secrets",
		"requirement:root::no-committed-secrets",
		"requirement:root::has-readme",
		"area:webhooks",
		"factor:webhooks::reliability",
	}
	if !equalStrings(first, want) {
		t.Fatalf("projection order = %v, want %v", first, want)
	}
}

func TestFind(t *testing.T) {
	root := Project(projectionSpec())
	if got := Find(root, "factor:root::security"); got == nil || got.Label != "Security" {
		t.Fatalf("Find(factor:root::security) = %+v", got)
	}
	if got := Find(root, "requirement:nope::missing"); got != nil {
		t.Fatalf("Find(missing) = %+v, want nil", got)
	}
}

func flatIDs(root *Element) []string {
	var ids []string
	for _, e := range Flatten(root) {
		ids = append(ids, e.ID)
	}
	return ids
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
