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
				Requirements: map[string]model.Requirement{
					"no-committed-secrets": {Title: "No credentials are committed"},
				},
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
								Requirements: map[string]model.Requirement{
									"retry-window": {Title: "Retry window is bounded"},
								},
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
	if got := (AreaPath{}).Display(); got != "/" {
		t.Fatalf("root area display = %q", got)
	}
	if got := (AreaPath{}).Reference(); got != "area:root" {
		t.Fatalf("root area reference = %q", got)
	}
	if got := (AreaPath{}).UnqualifiedReference(); got != "root" {
		t.Fatalf("root area unqualified reference = %q", got)
	}
	if got := (AreaPath{"webhooks", "delivery"}).Display(); got != "webhooks/delivery" {
		t.Fatalf("nested area display = %q", got)
	}
	if got := (AreaPath{"webhooks", "delivery"}).Reference(); got != "area:webhooks/delivery" {
		t.Fatalf("nested area reference = %q", got)
	}
	if got := (AreaPath{"webhooks", "delivery"}).UnqualifiedReference(); got != "webhooks/delivery" {
		t.Fatalf("nested area unqualified reference = %q", got)
	}
	if got := (FactorPath{"reliability", "retry-behavior"}).Display(); got != "reliability/retry-behavior" {
		t.Fatalf("factor display = %q", got)
	}
	if got := FactorReference(AreaPath{"webhooks", "delivery"}, FactorPath{"reliability", "retry-behavior"}); got != "factor:webhooks/delivery::reliability/retry-behavior" {
		t.Fatalf("factor reference = %q", got)
	}
	if got := UnqualifiedFactorReference(AreaPath{"webhooks", "delivery"}, FactorPath{"reliability", "retry-behavior"}); got != "webhooks/delivery::reliability/retry-behavior" {
		t.Fatalf("factor unqualified reference = %q", got)
	}
	if got := UnqualifiedFactorReference(AreaPath{}, FactorPath{"security"}); got != "root::security" {
		t.Fatalf("root factor unqualified reference = %q", got)
	}
	if got := RequirementReference(AreaPath{"webhooks", "delivery"}, "retry-window"); got != "requirement:webhooks/delivery::retry-window" {
		t.Fatalf("requirement reference = %q", got)
	}
	if got := UnqualifiedRequirementReference(AreaPath{}, "no-committed-secrets"); got != "root::no-committed-secrets" {
		t.Fatalf("root requirement unqualified reference = %q", got)
	}
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

	area, requirement, err := ParseRequirementReference(spec, "requirement:webhooks/delivery::retry-window")
	if err != nil {
		t.Fatalf("ParseRequirementReference() error = %v", err)
	}
	if !slices.Equal(area.Elements(), []string{"webhooks", "delivery"}) {
		t.Fatalf("requirement declaring area = %#v", area)
	}
	if requirement != "retry-window" {
		t.Fatalf("requirement = %q", requirement)
	}

	level, err := ParseRatingReference(spec, "rating:target")
	if err != nil {
		t.Fatalf("ParseRatingReference() error = %v", err)
	}
	if level != "target" {
		t.Fatalf("level = %q", level)
	}
}

func TestParseUnqualifiedModelReferences(t *testing.T) {
	spec := modelReferenceSpec()

	area, err := ParseUnqualifiedAreaReference(spec, "webhooks/delivery")
	if err != nil {
		t.Fatalf("ParseUnqualifiedAreaReference() error = %v", err)
	}
	if !slices.Equal(area.Elements(), []string{"webhooks", "delivery"}) {
		t.Fatalf("area path = %#v", area)
	}

	root, err := ParseUnqualifiedAreaReference(spec, "root")
	if err != nil {
		t.Fatalf("ParseUnqualifiedAreaReference(root) error = %v", err)
	}
	if len(root.Elements()) != 0 {
		t.Fatalf("root area path = %#v", root)
	}

	area, factor, err := ParseUnqualifiedFactorReference(spec, "webhooks/delivery::reliability/retry-behavior")
	if err != nil {
		t.Fatalf("ParseUnqualifiedFactorReference() error = %v", err)
	}
	if !slices.Equal(area.Elements(), []string{"webhooks", "delivery"}) {
		t.Fatalf("factor declaring area = %#v", area)
	}
	if !slices.Equal(factor.Elements(), []string{"reliability", "retry-behavior"}) {
		t.Fatalf("factor path = %#v", factor)
	}

	area, factor, err = ParseUnqualifiedFactorReference(spec, "root::security/secrets")
	if err != nil {
		t.Fatalf("ParseUnqualifiedFactorReference(root) error = %v", err)
	}
	if len(area.Elements()) != 0 {
		t.Fatalf("root factor declaring area = %#v", area)
	}
	if !slices.Equal(factor.Elements(), []string{"security", "secrets"}) {
		t.Fatalf("root factor path = %#v", factor)
	}

	area, requirement, err := ParseUnqualifiedRequirementReference(spec, "root::no-committed-secrets")
	if err != nil {
		t.Fatalf("ParseUnqualifiedRequirementReference(root) error = %v", err)
	}
	if len(area.Elements()) != 0 {
		t.Fatalf("root requirement declaring area = %#v", area)
	}
	if requirement != "no-committed-secrets" {
		t.Fatalf("root requirement = %q", requirement)
	}

	level, err := ParseUnqualifiedRatingReference(spec, "target")
	if err != nil {
		t.Fatalf("ParseUnqualifiedRatingReference() error = %v", err)
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
		{name: "unqualified area rejects typed reference", call: func() error {
			_, err := ParseUnqualifiedAreaReference(spec, "area:webhooks/delivery")
			return err
		}},
		{name: "unqualified area rejects display root", call: func() error {
			_, err := ParseUnqualifiedAreaReference(spec, "/")
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
		{name: "unqualified factor rejects typed reference", call: func() error {
			_, _, err := ParseUnqualifiedFactorReference(spec, "factor:webhooks/delivery::reliability")
			return err
		}},
		{name: "factor missing", call: func() error {
			_, _, err := ParseFactorReference(spec, "factor:webhooks/delivery::security")
			return err
		}},
		{name: "requirement missing separator", call: func() error {
			_, _, err := ParseRequirementReference(spec, "requirement:webhooks/delivery/retry-window")
			return err
		}},
		{name: "unqualified requirement rejects typed reference", call: func() error {
			_, _, err := ParseUnqualifiedRequirementReference(spec, "requirement:webhooks/delivery::retry-window")
			return err
		}},
		{name: "requirement bad name", call: func() error {
			_, _, err := ParseRequirementReference(spec, "requirement:webhooks/delivery::retry.window")
			return err
		}},
		{name: "requirement missing", call: func() error {
			_, _, err := ParseRequirementReference(spec, "requirement:webhooks/delivery::unknown-requirement")
			return err
		}},
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
