package lint

import "testing"

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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
    assessment: Inspect it.
---
`,
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
  "has an assessment":
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
  "has an assessment":
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
targets: {}
---
`,
	},
	{
		ruleID: RuleMisplacedRootKey,
		name:   "nested target title",
		model: validFrontmatter(`targets:
  api:
    title: API
    requirements:
      "has an assessment":
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleMisplacedRootKey,
		name:   "nested target rating scale",
		model: validFrontmatter(`targets:
  api:
    ratingScale:
      - level: target
        criterion: Meets it.
    requirements:
      "has an assessment":
        assessment: Inspect it.
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment absent",
		model: validFrontmatter(`requirements:
  "has an assessment": {}
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment empty",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: ""
`),
	},
	{
		ruleID: RuleInvalidAssessment,
		name:   "assessment list",
		model: validFrontmatter(`requirements:
  "has an assessment":
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
  "has an assessment":
    assessment: Inspect it.
    factors: [security]
`),
	},
	{
		ruleID: RuleUnknownFactor,
		name:   "child target names sibling factor",
		model: validFrontmatter(`targets:
  api:
    factors:
      reliability:
        description: Reliability.
    requirements:
      "api is reliable":
        assessment: Inspect it.
  worker:
    requirements:
      "worker is secure":
        assessment: Inspect it.
        factors: [reliability]
`),
	},
	{
		ruleID: RuleUnknownRatingKey,
		name:   "unknown override key",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: Inspect it.
    ratings:
      excellent: Exceeds it.
`),
	},
	{
		ruleID: RuleUnknownRatingKey,
		name:   "mixed known and unknown override keys",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: Inspect it.
    ratings:
      target: Meets it.
      poor: Does not meet it.
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
  "has an assessment":
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleMissingTitle,
		name:   "title absent with target content",
		model: `---
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
targets:
  api:
    requirements:
      "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
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
      "has an assessment":
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
          "has an assessment":
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
  "has an assessment":
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
  "has an assessment":
    assessment: Inspect it.
`),
	},
	{
		ruleID: RuleEmptyTarget,
		name:   "leaf target has no requirements",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: Inspect it.
targets:
  api:
    source: ./api
`),
	},
	{
		ruleID: RuleEmptyTarget,
		name:   "nested target subtree has no requirements",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: Inspect it.
targets:
  api:
    targets:
      handlers:
        source: ./handlers
`),
	},
	{
		ruleID:      RuleEmptyProperty,
		name:        "empty optional title",
		absentRules: []RuleID{RuleMissingTitle},
		model: `---
title: ""
ratingScale:
  - level: target
    description: Target.
    criterion: Meets it.
  - level: unacceptable
    description: Unacceptable.
    criterion: Does not meet it.
requirements:
  "has an assessment":
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
  "has an assessment":
    assessment: Inspect it.
---
`,
	},
	{
		ruleID: RuleEmptyProperty,
		name:   "empty optional secondary factors",
		model: validFrontmatter(`requirements:
  "has an assessment":
    assessment: Inspect it.
    factors: []
`),
	},
	{
		ruleID: RuleEmptyProperty,
		name:   "empty optional nested source",
		model: validFrontmatter(`targets:
  api:
    source:
    requirements:
      "has an assessment":
        assessment: Inspect it.
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

func TestRuleValidReferenceResolution(t *testing.T) {
	result, err := Check(writeModel(t, validFrontmatter(`factors:
  reliability:
    description: Reliability.
targets:
  api:
    requirements:
      "has an assessment":
        assessment: Inspect it.
        factors: [reliability]
`)))
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if hasRule(result, RuleUnknownFactor) {
		t.Fatalf("findings = %#v, ancestor factor should resolve", result.Findings)
	}
}
