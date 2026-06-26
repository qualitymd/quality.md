package evaluation

import (
	"testing"

	"github.com/qualitymd/quality.md/internal/model"
)

func ratingReferenceSpec() *model.Spec {
	return &model.Spec{
		Title: "Example",
		RatingScale: []model.RatingLevel{
			{Level: "target", Title: "Target"},
			{Level: "minimum", Title: "Minimum"},
		},
	}
}

func TestRatingReferenceRendering(t *testing.T) {
	if got := RatingReference("target"); got != "rating:target" {
		t.Fatalf("rating reference = %q", got)
	}
	if got := RatingDisplay("target"); got != "target" {
		t.Fatalf("rating display = %q", got)
	}
	if got := UnqualifiedRatingReference("target"); got != "target" {
		t.Fatalf("rating unqualified reference = %q", got)
	}
}

func TestParseRatingReference(t *testing.T) {
	spec := ratingReferenceSpec()

	level, err := ParseRatingReference(spec, "rating:target")
	if err != nil {
		t.Fatalf("ParseRatingReference() error = %v", err)
	}
	if level != "target" {
		t.Fatalf("level = %q", level)
	}

	level, err = ParseUnqualifiedRatingReference(spec, "target")
	if err != nil {
		t.Fatalf("ParseUnqualifiedRatingReference() error = %v", err)
	}
	if level != "target" {
		t.Fatalf("level = %q", level)
	}
}

func TestParseRatingReferenceRejectInvalidOrUnresolvedInput(t *testing.T) {
	spec := ratingReferenceSpec()
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{name: "rating bad id", call: func() error {
			_, err := ParseRatingReference(spec, "rating:not.acceptable")
			return err
		}},
		{name: "unqualified rating rejects typed reference", call: func() error {
			_, err := ParseUnqualifiedRatingReference(spec, "rating:target")
			return err
		}},
		{name: "rating missing", call: func() error {
			_, err := ParseRatingReference(spec, "rating:outstanding")
			return err
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.call(); err == nil {
				t.Fatal("parse error = nil, want rejection")
			}
		})
	}
}
