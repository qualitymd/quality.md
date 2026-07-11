package lint

import (
	"os"
	"path/filepath"
	"testing"
)

// ruleCases is the exhaustive per-rule fixture table. Every rule in the catalog
// must appear at least twice (asserted below), and each case may assert that
// certain other rules are absent.
var ruleCases = []struct {
	ruleID      RuleID
	name        string
	model       string
	absentRules []RuleID
}{
	{
		ruleID: RuleInvalidFrontmatter,
		name:   "no frontmatter fence",
		model:  "title: Example\n",
	},
	{
		ruleID: RuleInvalidFrontmatter,
		name:   "invalid YAML",
		model:  "---\ntitle: [\n---\n",
	},
	{
		ruleID: RuleInvalidFrontmatter,
		name:   "frontmatter is a list",
		model:  "---\n- title\n---\n",
	},
	{
		ruleID: RuleMissingRatingScale,
		name:   "rating scale absent",
		model: `---
title: Example
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID:      RuleMissingRatingScale,
		name:        "rating scale null",
		absentRules: []RuleID{RuleInvalidFrontmatter},
		model: `---
title: Example
ratingScale:
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID:      RuleMissingRatingScale,
		name:        "rating scale empty string",
		absentRules: []RuleID{RuleInvalidFrontmatter},
		model: `---
title: Example
ratingScale: ""
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleTooFewLevels,
		name:   "one level",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleTooFewLevels,
		name:   "one level with missing criterion",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingLevelName,
		name:   "level name absent",
		model: `---
title: Example
ratingScale:
  - description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingLevelName,
		name:   "level name null",
		model: `---
title: Example
ratingScale:
  - level:
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleInvalidRatingLevelID,
		name:   "level id contains dot",
		model: `---
title: Example
ratingScale:
  - level: tar.get
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleInvalidRatingLevelID,
		name:   "level id has trailing separator",
		model: `---
title: Example
ratingScale:
  - level: target_
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleDuplicateLevel,
		name:   "two duplicate levels",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: target
    description: Duplicate.
    criterion: Also meets it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleDuplicateLevel,
		name:   "duplicate after distinct level",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
  - level: target
    description: Duplicate.
    criterion: Also meets it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleInvalidConfig,
		name:   "root config is empty",
		model: validFrontmatter(`config:
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidConfig,
		name:   "root config is absolute",
		model: validFrontmatter(`config: /tmp/outside.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleMissingCriterion,
		name:   "criterion absent",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingCriterion,
		name:   "criterion empty",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: ""
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleEmptyModel,
		name:   "no model content",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
---
`,
	},
	{
		ruleID: RuleEmptyModel,
		name:   "empty maps supply no entries",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
factors: {}
requirements: {}
areas: {}
---
`,
	},
	{
		ruleID: RuleMisplacedRootKey,
		name:   "nested area rating scale",
		model: validFrontmatter(`areas:
  api:
    ratingScale:
      - level: target
        criterion: Meets it.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleMisplacedRootKey,
		name:   "deep nested area rating scale",
		model: validFrontmatter(`areas:
  api:
    areas:
      handlers:
        ratingScale:
          - level: target
            criterion: Meets it.
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidAreaName,
		name:   "area name contains slash",
		model: validFrontmatter(`areas:
  api/service:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidAreaName,
		name:   "nested area name starts with separator",
		model: validFrontmatter(`areas:
  api:
    areas:
      _handlers:
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleReservedAreaName,
		name:   "root area name at top level",
		model: validFrontmatter(`areas:
  root:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleReservedAreaName,
		name:   "root area name nested",
		model: validFrontmatter(`areas:
  api:
    areas:
      root:
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment absent",
		model: validFrontmatter(`requirements:
  missing-assessment: {}
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment empty",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: ""
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment list",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment:
      - Inspect it.
`),
	},
	{
		ruleID: RuleUnknownFactor,
		name:   "root requirement names missing factor",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: [security]
`),
	},
	{
		ruleID: RuleUnknownFactor,
		name:   "child area names sibling factor",
		model: validFrontmatter(`areas:
  api:
    factors:
      reliability:
        description: Reliability.
    requirements:
      api-reliable:
        title: API is reliable
        assessment: Inspect it.
  worker:
    requirements:
      worker-secure:
        title: Worker is secure
        assessment: Inspect it.
        factors: [reliability]
`),
	},
	{
		ruleID: RuleUnknownFactor,
		name:   "child area names ancestor factor",
		model: validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
areas:
  api:
    requirements:
      api-reliable:
        title: API is reliable
        assessment: Inspect it.
        factors: [reliability]
`),
	},
	{
		ruleID: RuleInvalidRequirementName,
		name:   "requirement name contains spaces",
		model: validFrontmatter(`requirements:
  "has an assessment":
    title: Has an assessment
    assessment: Inspect it.
    factors: [reliability]
factors:
  reliability:
    title: Reliability
    description: Reliability.
`),
	},
	{
		ruleID: RuleInvalidRequirementName,
		name:   "nested requirement name starts with separator",
		model: validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
    requirements:
      _has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleDuplicateRequirement,
		name:   "area and factor requirements share name",
		model: validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: [reliability]
`),
	},
	{
		ruleID: RuleDuplicateRequirement,
		name:   "sibling factors share requirement name",
		model: validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
  security:
    title: Security
    description: Security.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleUnknownRatingKey,
		name:   "unknown override key",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    ratings:
      excellent: Exceeds it.
`),
	},
	{
		ruleID: RuleUnknownRatingKey,
		name:   "mixed known and unknown override keys",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    ratings:
      target: Meets it.
      poor: Does not meet it.
`),
	},
	{
		ruleID: RuleMissingFactorReference,
		name:   "direct requirement missing factors",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleMissingFactorReference,
		name:   "direct requirement empty factors",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: []
`),
	},
	{
		ruleID: RuleMissingTitle,
		name:   "title absent",
		model: `---
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingTitle,
		name:   "title absent with area content",
		model: `---
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
areas:
  api:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingLevelDescription,
		name:   "first level description absent",
		model: `---
title: Example
ratingScale:
  - level: target
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingLevelDescription,
		name:   "second level description absent",
		model: `---
title: Example
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingFactorDescription,
		name:   "root factor description absent",
		model: validFrontmatter(`factors:
  reliability:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleMissingFactorDescription,
		name:   "nested factor description absent",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
    factors:
      availability:
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidFactorName,
		name:   "factor name contains space",
		model: validFrontmatter(`factors:
  service health:
    description: Service health.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidFactorName,
		name:   "nested factor name ends with separator",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
    factors:
      availability-:
        description: Availability.
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleDuplicateFactorName,
		name:   "root factor and subfactor share name",
		model: validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
    factors:
      reliability:
        title: Reliability detail
        description: Reliability detail.
        requirements:
          has-assessment:
            title: Has an assessment
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleDuplicateFactorName,
		name:   "sibling branches share nested factor name",
		model: validFrontmatter(`factors:
  operations:
    title: Operations
    description: Operations.
    factors:
      security:
        title: Security
        description: Security.
        requirements:
          operations-secure:
            title: Operations are secure
            assessment: Inspect it.
  delivery:
    title: Delivery
    description: Delivery.
    factors:
      security:
        title: Security
        description: Security.
        requirements:
          delivery-secure:
            title: Delivery is secure
            assessment: Inspect it.
`),
	},
	{
		ruleID: RuleEmptyFactor,
		name:   "root factor has no requirements",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleEmptyFactor,
		name:   "nested factor has no requirements",
		model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
    factors:
      availability:
        description: Availability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleEmptyArea,
		name:   "leaf area has no requirements",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
areas:
  api:
    source: ./api
`),
	},
	{
		ruleID: RuleEmptyArea,
		name:   "nested area subtree has no requirements",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
areas:
  api:
    areas:
      handlers:
        source: ./handlers
`),
	},
	{
		ruleID:      RuleEmptyProperty,
		name:        "empty optional description",
		absentRules: []RuleID{RuleMissingTitle},
		model: `---
title: Example
description: ""
ratingScale:
  - level: target
    title: Target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    title: Unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID:      RuleEmptyProperty,
		name:        "empty optional level description",
		absentRules: []RuleID{RuleMissingLevelDescription},
		model: `---
title: Example
ratingScale:
  - level: target
    description: ""
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleEmptyProperty,
		name:   "empty optional factor references",
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: []
`),
	},
	{
		ruleID: RuleEmptyProperty,
		name:   "empty optional nested source",
		model: validFrontmatter(`areas:
  api:
    source:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
	},
	{
		ruleID:      RuleUnknownKey,
		name:        "root extension key",
		absentRules: []RuleID{RuleInvalidFrontmatter},
		model: validFrontmatter(`x-owner: platform-team
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
	},
	{
		ruleID:      RuleUnknownKey,
		name:        "requirement extension key",
		absentRules: []RuleID{RuleInvalidFrontmatter},
		model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    assessmnt: A typo of assessment.
`),
	},
}

func TestRules(t *testing.T) {
	// Every rule in the catalog must have at least two distinct fixtures.
	counts := map[RuleID]int{}
	for _, tc := range ruleCases {
		counts[tc.ruleID]++
	}
	for _, rule := range Rules {
		if counts[rule.ID] < 2 {
			t.Errorf("rule %s has %d fixtures, want at least 2", rule.ID, counts[rule.ID])
		}
	}

	for _, tc := range ruleCases {
		t.Run(string(tc.ruleID)+"/"+tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, tc.model))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if !hasRule(result, tc.ruleID) {
				t.Fatalf("findings = %#v, want rule %s", result.Findings, tc.ruleID)
			}
			for _, absent := range tc.absentRules {
				if hasRule(result, absent) {
					t.Fatalf("findings = %#v, did not want rule %s", result.Findings, absent)
				}
			}
		})
	}
}

func TestRuleAreaLocalReferenceResolution(t *testing.T) {
	result, err := Check(writeModel(t, validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
areas:
  api:
    factors:
      reliability:
        title: Reliability
        description: API reliability.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
        factors: [reliability]
`)))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if hasRule(result, RuleUnknownFactor) {
		t.Fatalf("findings = %#v, local factor should resolve", result.Findings)
	}
	if hasRule(result, RuleMissingFactorReference) {
		t.Fatalf("findings = %#v, direct requirement references a local factor", result.Findings)
	}
	if hasRule(result, RuleDuplicateFactorName) {
		t.Fatalf("findings = %#v, same factor name on different areas should not warn", result.Findings)
	}
}

func TestRuleAncestorFactorReferenceDoesNotResolve(t *testing.T) {
	result, err := Check(writeModel(t, validFrontmatter(`factors:
  reliability:
    title: Reliability
    description: Reliability.
areas:
  api:
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
        factors: [reliability]
`)))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !hasRule(result, RuleUnknownFactor) {
		t.Fatalf("findings = %#v, ancestor factor should not resolve", result.Findings)
	}
	if hasRule(result, RuleMissingFactorReference) {
		t.Fatalf("findings = %#v, non-empty factor reference should satisfy reference presence", result.Findings)
	}
}

func TestRuleRequirementFactorReferences(t *testing.T) {
	for _, tc := range []struct {
		name              string
		body              string
		wantEmptyProperty bool
	}{
		{
			name: "missing factors",
			body: `factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`,
		},
		{
			name:              "null factors",
			wantEmptyProperty: true,
			body: `factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors:
`,
		},
		{
			name:              "empty factors",
			wantEmptyProperty: true,
			body: `factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: []
`,
		},
		{
			name: "only empty factor entries",
			body: `factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: ["", null]
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, validFrontmatter(tc.body)))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if !hasRule(result, RuleMissingFactorReference) {
				t.Fatalf("findings = %#v, want %s", result.Findings, RuleMissingFactorReference)
			}
			if tc.wantEmptyProperty && !hasRule(result, RuleEmptyProperty) {
				t.Fatalf("findings = %#v, want %s", result.Findings, RuleEmptyProperty)
			}
		})
	}

	for _, tc := range []struct {
		name string
		body string
	}{
		{
			name: "direct requirement references local factor",
			body: `factors:
  reliability:
    description: Reliability.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    factors: [reliability]
`,
		},
		{
			name: "nested requirement declared under containing factor",
			body: `factors:
  reliability:
    description: Reliability.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, validFrontmatter(tc.body)))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if hasRule(result, RuleMissingFactorReference) {
				t.Fatalf("findings = %#v, requirement should be connected to a factor", result.Findings)
			}
		})
	}
}

func TestAreaDisplayFieldsAreValid(t *testing.T) {
	for _, tc := range []struct {
		name  string
		model string
	}{
		{
			name: "title",
			model: validFrontmatter(`areas:
  api:
    title: API
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
		{
			name: "description",
			model: validFrontmatter(`areas:
  api:
    description: Functional specifications for the API.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
		{
			name: "title and description",
			model: validFrontmatter(`areas:
  api:
    title: API
    description: Functional specifications for the API.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, tc.model))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if hasRule(result, RuleMisplacedRootKey) || hasRule(result, RuleInvalidFrontmatter) {
				t.Fatalf("findings = %#v, area display fields should be valid", result.Findings)
			}
		})
	}
}

func TestAreaDisplayFieldShapesAreValidated(t *testing.T) {
	for _, tc := range []struct {
		name  string
		model string
	}{
		{
			name: "title list",
			model: validFrontmatter(`areas:
  api:
    title: [API]
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
		{
			name: "description map",
			model: validFrontmatter(`areas:
  api:
    description:
      text: Functional specifications for the API.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, tc.model))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if !hasRule(result, RuleInvalidFrontmatter) {
				t.Fatalf("findings = %#v, want %s", result.Findings, RuleInvalidFrontmatter)
			}
			if hasRule(result, RuleMisplacedRootKey) {
				t.Fatalf("findings = %#v, area display fields should not be root-key findings", result.Findings)
			}
		})
	}
}

func TestSchemaDrivenUnknownKeys(t *testing.T) {
	for _, tc := range []struct {
		name  string
		model string
	}{
		{
			name: "root",
			model: validFrontmatter(`unexpected: true
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
		},
		{
			name: "rating level",
			model: `---
title: Example
ratingScale:
  - level: target
    criterion: Meets it.
    unexpected: true
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
---
`,
		},
		{
			name: "target",
			model: validFrontmatter(`areas:
  api:
    unexpected: true
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
		{
			name: "factor",
			model: validFrontmatter(`factors:
  reliability:
    description: Reliability.
    unexpected: true
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
`),
		},
		{
			name: "requirement",
			model: validFrontmatter(`requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
    unexpected: true
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, tc.model))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if !hasRule(result, RuleUnknownKey) {
				t.Fatalf("findings = %#v, want %s", result.Findings, RuleUnknownKey)
			}
			if hasRule(result, RuleInvalidFrontmatter) {
				t.Fatalf("findings = %#v, extension frontmatter must not be %s", result.Findings, RuleInvalidFrontmatter)
			}
			for _, finding := range result.Findings {
				if finding.RuleID == RuleUnknownKey && finding.Severity != SeverityWarning {
					t.Fatalf("unknown-key severity = %s, want %s", finding.Severity, SeverityWarning)
				}
			}
		})
	}
}

// TestExtensionFrontmatterLintsValid pins the spec's Extensions rule: a
// conforming document carrying extension frontmatter stays valid and loads.
func TestExtensionFrontmatterLintsValid(t *testing.T) {
	path := writeModel(t, validFrontmatter(`x-owner: platform-team
factors:
  reliability:
    title: Reliability
    description: Reliability.
    requirements:
      has-assessment:
        title: Has an assessment
        assessment: Inspect it.
        x-ticket: QM-42
`))
	result, err := Check(path)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if !result.Valid || result.Summary.Errors != 0 {
		t.Fatalf("result = %+v, want extension frontmatter to lint valid", result.Summary)
	}
	if result.Summary.Warnings == 0 || !hasRule(result, RuleUnknownKey) {
		t.Fatalf("findings = %#v, want unknown-key advisories for the extension keys", result.Findings)
	}
	if _, err := Load(path); err != nil {
		t.Fatalf("Load() error = %v, want extension frontmatter to load", err)
	}
}

func TestRootConfigToolingKey(t *testing.T) {
	valid, err := Check(writeModel(t, validFrontmatter(`config: .quality/config.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`)))
	if err != nil {
		t.Fatalf("Check(valid config) error = %v", err)
	}
	if hasRule(valid, RuleInvalidFrontmatter) || hasRule(valid, RuleInvalidConfig) {
		t.Fatalf("findings = %#v, want root config accepted", valid.Findings)
	}

	for _, tc := range []struct {
		name  string
		model string
	}{
		{
			name: "map",
			model: validFrontmatter(`config:
  path: .quality/config.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
		},
		{
			name: "absolute",
			model: validFrontmatter(`config: /tmp/config.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Check(writeModel(t, tc.model))
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			if !hasRule(result, RuleInvalidConfig) {
				t.Fatalf("findings = %#v, want %s", result.Findings, RuleInvalidConfig)
			}
		})
	}
}

func TestRootConfigToolingKeyIsModelRelativeInsideRepository(t *testing.T) {
	repo := t.TempDir()
	if err := os.Mkdir(filepath.Join(repo, ".git"), 0o755); err != nil {
		t.Fatalf("Mkdir(.git) error = %v", err)
	}
	modelDir := filepath.Join(repo, "packages", "api")
	if err := os.MkdirAll(modelDir, 0o755); err != nil {
		t.Fatalf("MkdirAll(model dir) error = %v", err)
	}
	inside := filepath.Join(modelDir, "QUALITY.md")
	if err := os.WriteFile(inside, []byte(validFrontmatter(`config: ../shared/quality-config.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`)), 0o600); err != nil {
		t.Fatalf("WriteFile(inside) error = %v", err)
	}
	result, err := Check(inside)
	if err != nil {
		t.Fatalf("Check(inside) error = %v", err)
	}
	if hasRule(result, RuleInvalidConfig) {
		t.Fatalf("findings = %#v, want in-repository model-relative config accepted", result.Findings)
	}

	outside := filepath.Join(modelDir, "QUALITY-outside.md")
	if err := os.WriteFile(outside, []byte(validFrontmatter(`config: ../../../outside.yaml
requirements:
  has-assessment:
    title: Has an assessment
    assessment: Inspect it.
`)), 0o600); err != nil {
		t.Fatalf("WriteFile(outside) error = %v", err)
	}
	result, err = Check(outside)
	if err != nil {
		t.Fatalf("Check(outside) error = %v", err)
	}
	if !hasRule(result, RuleInvalidConfig) {
		t.Fatalf("findings = %#v, want repository-escaping config rejected", result.Findings)
	}
}
