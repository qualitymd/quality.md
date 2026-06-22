package evaluation

import (
	"slices"
	"testing"

	"github.com/qualitymd/quality.md/internal/model"
)

func modelReferenceSpec() *model.Spec {
	return &model.Spec{
		Title: "Example",
		RatingScale: []model.RatingLevel{
			{Level: "target", Title: "Target"},
			{Level: "minimum", Title: "Minimum"},
		},
		Factors: map[string]model.Factor{
			"security": {
				Title: "Security",
				Factors: map[string]model.Factor{
					"secrets": {Title: "Secrets"},
				},
			},
		},
		Areas: map[string]model.Area{
			"webhooks": {
				Title: "Webhooks",
				Areas: map[string]model.Area{
					"delivery": {
						Title: "Delivery",
						Factors: map[string]model.Factor{
							"reliability": {
								Title: "Reliability",
								Factors: map[string]model.Factor{
									"retry-behavior": {Title: "Retry behavior"},
								},
							},
						},
					},
				},
			},
		},
	}
}

func TestModelReferenceRendering(t *testing.T) {
	if got := (AreaPath{}).Reference(); got != "area:root" {
		t.Fatalf("root area reference = %q", got)
	}
	if got := (AreaPath{"webhooks", "delivery"}).Reference(); got != "area:webhooks/delivery" {
		t.Fatalf("nested area reference = %q", got)
	}
	if got := FactorReference(AreaPath{"webhooks", "delivery"}, FactorPath{"reliability", "retry-behavior"}); got != "factor:webhooks/delivery::reliability/retry-behavior" {
		t.Fatalf("factor reference = %q", got)
	}
	if got := RatingReference("target"); got != "rating:target" {
		t.Fatalf("rating reference = %q", got)
	}
}

func TestParseModelReferences(t *testing.T) {
	spec := modelReferenceSpec()

	area, err := ParseAreaReference(spec, "area:webhooks/delivery")
	if err != nil {
		t.Fatalf("ParseAreaReference() error = %v", err)
	}
	if !slices.Equal(area.Elements(), []string{"webhooks", "delivery"}) {
		t.Fatalf("area path = %#v", area)
	}

	area, factor, err := ParseFactorReference(spec, "factor:webhooks/delivery::reliability/retry-behavior")
	if err != nil {
		t.Fatalf("ParseFactorReference() error = %v", err)
	}
	if !slices.Equal(area.Elements(), []string{"webhooks", "delivery"}) {
		t.Fatalf("factor declaring area = %#v", area)
	}
	if !slices.Equal(factor.Elements(), []string{"reliability", "retry-behavior"}) {
		t.Fatalf("factor path = %#v", factor)
	}

	level, err := ParseRatingReference(spec, "rating:target")
	if err != nil {
		t.Fatalf("ParseRatingReference() error = %v", err)
	}
	if level != "target" {
		t.Fatalf("level = %q", level)
	}
}

func TestParseModelReferencesRejectInvalidOrUnresolvedInput(t *testing.T) {
	spec := modelReferenceSpec()
	for _, tc := range []struct {
		name string
		call func() error
	}{
		{name: "area shorthand", call: func() error {
			_, err := ParseAreaReference(spec, "webhooks/delivery")
			return err
		}},
		{name: "area bad segment", call: func() error {
			_, err := ParseAreaReference(spec, "area:webhooks.delivery")
			return err
		}},
		{name: "area missing", call: func() error {
			_, err := ParseAreaReference(spec, "area:webhooks/missing")
			return err
		}},
		{name: "factor missing separator", call: func() error {
			_, _, err := ParseFactorReference(spec, "factor:webhooks/delivery/reliability")
			return err
		}},
		{name: "factor missing", call: func() error {
			_, _, err := ParseFactorReference(spec, "factor:webhooks/delivery::security")
			return err
		}},
		{name: "rating bad id", call: func() error {
			_, err := ParseRatingReference(spec, "rating:not.acceptable")
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
