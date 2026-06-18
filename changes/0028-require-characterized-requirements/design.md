---
type: Design Doc
title: Require characterized requirements
description: How mandatory requirement characterization lands in lint, docs, and terminology.
tags: [specification, lint, terminology, design]
timestamp: 2026-06-18T00:00:00Z
---

# Require characterized requirements

Design behind the [Require characterized requirements](../0028-require-characterized-requirements.md)
change and its [functional spec](spec.md).

## Context

The lint traversal already records the two facts this change needs:

- every requirement has a declaring target; and
- every requirement either has a containing factor (`requirementRef.factor`) or
  is declared directly under a target (`factor == nil`).

That means the new conformance rule can be a narrow semantic check layered on
the existing parsed model. The structural schema should not make `factors`
globally required on a Requirement, because a nested requirement is already
characterized by its containing factor and should not have to repeat it.

## Approach

Add `RuleUncharacterizedRequirement` to the lint catalog:

```go
RuleUncharacterizedRequirement RuleID = "uncharacterized-requirement"
```

Catalog metadata:

- severity: `error`
- fixable: `false`
- description: "A requirement is not characterized by any factor."

Implement the check in the requirement-reference pass, next to `unknown-factor`
and `unknown-rating-key`. A helper should report whether a requirement has at
least one non-empty scalar factor entry:

```go
func (s *runState) requirementHasListedFactor(req *requirementRef) bool
```

The rule fires when both conditions hold:

1. `req.factor == nil`
2. `requirementHasListedFactor(req) == false`

Locate the finding at the missing `factors` key under the requirement:
`appendPath(req.path, "factors")`. That keeps the remediation location concrete
without pretending the key exists.

The message should use the new terminology and a deterministic repair hint:

> The requirement `<statement>` is not characterized by any factor; place it
> under a factor or list one or more `factors`.

Do not suppress the existing `empty-property` warning for `factors: []`. An
empty optional property is still mechanically wrong, and the same requirement is
also uncharacterized. The two findings point at different defects: omit empty
properties, and provide a factor association.

Keep `unknown-factor` as the resolver error for every listed factor that fails
to resolve. Update its static description and finding message so they no longer
always say "secondary factor"; the same code path validates direct
target-level `factors` and additional factors on nested requirements. Message
wording can be neutral:

> The requirement `<statement>` names unknown factor `<name>`; listed factors
> must resolve on the declaring target or an ancestor.

The helper currently named `secondaryFactors` should be renamed to something
neutral such as `listedFactors` or `resolvedListedFactors`, because direct
target-level lists are not secondary. `empty-factor` should use the renamed
helper without behavior changes.

## Durable Docs

`SPECIFICATION.md` needs the main semantic rewrite:

- Factor prose can still use "lens" as explanatory shorthand.
- Requirement prose should lead with "characterized by" for the mechanics.
- Direct target-level requirements must list factors.
- Listed factors are "secondary" only when a requirement is already nested under
  a factor.

`specs/cli/lint.md` adds the `uncharacterized-requirement` error row and updates
the `unknown-factor` row to refer to listed factors rather than only secondary
factors.

`README.md` and `internal/scaffold/skeleton.md` should nudge authored examples
toward characterized requirements. The scaffold should keep its placeholder
requirement under a factor, so it remains lint-valid without adding a direct
target-level `factors` example.

## Tests

Update `internal/lint/rules_test.go`:

- add at least two fixtures for `uncharacterized-requirement`, matching the
  rule-catalog fixture-count guard;
- update existing direct target-level requirement fixtures that currently omit
  `factors` so they remain valid under the stricter rule when they are not meant
  to test this case;
- keep a case proving a direct requirement characterized by an ancestor factor
  is valid;
- keep a case proving a nested requirement without `factors` is valid because
  the containing factor characterizes it; and
- update `unknown-factor` expectations to cover both direct listed factors and
  additional factors on nested requirements.

Run the full Go test suite after implementation. The schema snippet consistency
test should not require code changes, because the Requirement YAML snippet still
lists `factors` as optional: optional structurally, conditionally required only
for direct target-level requirements.

## Alternatives

- **Name the rule `missing-factor`.** Rejected because it is ambiguous in this
  catalog: it could mean no model `factors` property, no factor declaration, a
  missing `factors` key on a requirement, an unresolved listed factor, or an
  empty factor. `uncharacterized-requirement` names the semantic violation.
- **Make `factors` required in the structural schema.** Rejected because it would
  force nested requirements to repeat their containing factor and would blur the
  difference between structure and context-sensitive semantics.
- **Rename every "lens" occurrence.** Rejected because "lens" remains a useful
  shorthand for how a factor frames quality. This change only makes
  "characterized by" the primary term for the requirement-to-factor mechanics.

## Trade-offs & Risks

- Existing valid models with direct target-level requirements and no `factors`
  become invalid. That is intentional but breaking; docs and scaffold examples
  need to make the required shape obvious.
- The linter may emit both `empty-property` and
  `uncharacterized-requirement` for `factors: []`. That is noisier but accurate,
  and the remedies are different.
- If direct target-level requirements are common in tests, the implementation
  will touch more fixtures than the rule itself suggests. Keep those edits
  mechanical: add a factor association only where the test is unrelated to
  characterization.

## Open Questions

- None blocking.
