package model

import "testing"

func TestEffectiveSource(t *testing.T) {
	spec := &Spec{
		Source: "src",
		Areas: map[string]Area{
			"api": {
				Areas: map[string]Area{
					"handlers": {},
				},
			},
			"cli": {
				Source: "cmd",
				Areas: map[string]Area{
					"flags": {},
				},
			},
		},
	}
	for _, tc := range []struct {
		name     string
		path     AreaPath
		selector string
		state    SourceState
	}{
		{"root declared", nil, "src", SourceStateDeclared},
		{"child inherits root", AreaPath{"api"}, "src", SourceStateInherited},
		{"grandchild inherits root", AreaPath{"api", "handlers"}, "src", SourceStateInherited},
		{"child declared", AreaPath{"cli"}, "cmd", SourceStateDeclared},
		{"grandchild inherits nearest ancestor", AreaPath{"cli", "flags"}, "cmd", SourceStateInherited},
	} {
		t.Run(tc.name, func(t *testing.T) {
			selector, state := EffectiveSource(spec, tc.path)
			if selector != tc.selector || state != tc.state {
				t.Fatalf("EffectiveSource(%v) = %q/%q, want %q/%q", tc.path, selector, state, tc.selector, tc.state)
			}
		})
	}
}

func TestEffectiveSourceDefaultsToDocumentDirectory(t *testing.T) {
	spec := &Spec{
		Areas: map[string]Area{
			"api": {Areas: map[string]Area{"handlers": {}}},
		},
	}
	for _, path := range []AreaPath{nil, {"api"}, {"api", "handlers"}} {
		selector, state := EffectiveSource(spec, path)
		if selector != DefaultSource || state != SourceStateDefault {
			t.Fatalf("EffectiveSource(%v) = %q/%q, want %q/%q", path, selector, state, DefaultSource, SourceStateDefault)
		}
	}
}
