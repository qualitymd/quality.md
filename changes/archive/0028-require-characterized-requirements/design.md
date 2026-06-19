---
type: Design Doc
title: Require factor references
description: How mandatory factor references land in lint, docs, and terminology.
tags: [specification, lint, terminology, design]
timestamp: 2026-06-18T00:00:00Z
---

# Require factor references

Design behind the [Require factor references](../0028-require-characterized-requirements.md)
change and its [functional spec](spec.md).

## Context

The lint traversal already records the two facts this change needs:

- every requirement has a declaring target; and
- every requirement either has a containing factor (`requirementRef.factor`) or
  is declared directly under a target (`factor == nil`).

That means the new conformance rule can be a narrow semantic check layered on
the existing parsed model. The structural schema should not make `factors`
globally required on a Requirement, because a nested requirement is already
connected to its containing factor by placement.

## Approach

Add `RuleMissingFactorReference` to the lint catalog:

```go
RuleMissingFactorReference RuleID = "missing-factor-reference"
```

Catalog metadata:

- severity: `error`
- fixable: `false`
- description: "A direct target-level requirement references no quality factor."

Implement the check in the requirement-reference pass, next to `unknown-factor`
and `unknown-rating-key`. A helper should report whether a requirement has at
least one non-empty scalar factor entry:

```go
func requirementReferencesFactor(req *requirementRef) bool
```

The rule fires when both conditions hold:

1. `req.factor == nil`
2. `requirementReferencesFactor(req) == false`

Locate the finding at the missing `factors` key under the requirement:
`appendPath(req.path, "factors")`. That keeps the remediation location concrete
without pretending the key exists.

The message should use the new terminology and a deterministic repair hint:

> The requirement `<statement>` references no quality factor; place it under a
> factor or add one or more factor references under `factors`.

Do not suppress the existing `empty-property` warning for `factors: null` or
`factors: []`. A null or empty optional property is still mechanically wrong,
and the same requirement also has no factor reference. The two findings point at
different defects: omit empty properties, and provide a factor reference.

Keep `unknown-factor` as the resolver error for every factor reference that fails
to resolve. Update its static description and finding message so they no longer
always say "secondary factor"; the same code path validates direct
target-level `factors` and additional factors on nested requirements. Message
wording can be neutral:

> The requirement `<statement>` references unknown factor `<name>`; factor
> references must resolve on the declaring target or an ancestor.

The helper currently named `secondaryFactors` should be renamed to something
neutral such as `referencedFactors` or `resolvedFactorReferences`, because
direct target-level references are not secondary. `empty-factor` should use the renamed
helper without behavior changes.

## Documentation Approach

Keep durable spec accounting in the functional spec and durable doc accounting
in the parent change case. Implementation should then make the recorded
terminology pass mechanically: use "reference" for requirement-to-factor
mechanics, reserve "characterize" for factors describing targets, keep "lens"
only as explanatory shorthand, and ensure examples remain lint-valid under
mandatory factor references.

## Tests

Update `internal/lint/rules_test.go`:

- add at least two fixtures for `missing-factor-reference`, matching the
  rule-catalog fixture-count guard;
- cover null and empty `factors` values so both `empty-property` and
  `missing-factor-reference` are expected for direct target-level
  requirements;
- update existing direct target-level requirement fixtures that currently omit
  `factors` so they remain valid under the stricter rule when they are not meant
  to test this case;
- keep a case proving a direct requirement references an ancestor factor
  is valid;
- keep a case proving a nested requirement without `factors` is valid because
  it is connected to its containing factor by placement; and
- update `unknown-factor` expectations to cover both direct factor references
  and additional factors on nested requirements.

Run the full Go test suite after implementation. The schema snippet consistency
test should not require code changes, because the Requirement YAML snippet still
shows `factors` as optional: optional structurally, conditionally required only
for direct target-level requirements.

## Alternatives

- **Name the rule `missing-factor`.** Rejected because it is ambiguous in this
  catalog: it could mean no model `factors` property, no factor declaration, a
  missing `factors` key on a requirement, an unresolved factor reference, or an
  empty factor. `missing-factor-reference` names the semantic violation.
- **Make `factors` required in the structural schema.** Rejected because it would
  force nested requirements to repeat their containing factor and would blur the
  difference between structure and context-sensitive semantics.
- **Rename every "lens" occurrence.** Rejected because "lens" remains a useful
  shorthand for how a factor frames quality. This change only makes
  "reference" the primary term for the requirement-to-factor mechanics.

## Trade-offs & Risks

- Existing valid models with direct target-level requirements and no `factors`
  become invalid. That is intentional but breaking; docs and scaffold examples
  need to make the required shape obvious.
- The linter may emit both `empty-property` and
  `missing-factor-reference` for `factors: null` or `factors: []`. That is
  noisier but accurate, and the remedies are different.
- If direct target-level requirements are common in tests, the implementation
  will touch more fixtures than the rule itself suggests. Keep those edits
  mechanical: add a factor reference only where the test is unrelated to
  missing factor references.

## Open Questions

- None blocking.
