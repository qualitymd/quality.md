---
type: Functional Specification
title: Require factor references
description: Requirements for making requirements without factor references invalid and aligning requirement-to-factor terminology.
tags: [specification, lint, terminology]
timestamp: 2026-06-18T00:00:00Z
---

# Require factor references

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

The format currently permits a requirement declared directly under a target to
inform only the target's local rating, with no factor reference. That
"unlensed" special case weakens factor roll-ups and complicates the model:
requirements can be assessed without naming the quality characteristic they
measure. Requiring every requirement to be connected to at least one factor
makes the model simpler and makes each assessment's quality concern explicit.

The terminology should also distinguish connection by placement from explicit
factor references. A requirement nested under a factor is connected to that
factor by placement. A `factors` property on a nested requirement adds explicit
references to additional factors; those factor references are secondary. A
`factors` property on a requirement declared directly under a target is not
secondary to anything; it is the set of explicit factor references for that
requirement.

## Requirements

Every requirement **MUST** be connected to at least one factor.

> Rationale: A requirement that affects only a target's local rating creates a
> special case with no quality-characteristic roll-up. Mandatory
> required factor connections keep every assessment tied to a quality factor. —
> 0028

A requirement nested under a factor or sub-factor **MUST** be connected to that
containing factor by placement.

A requirement nested under a factor or sub-factor **MAY** declare additional
factor references under `factors`; each referenced factor **MUST** resolve to a
factor declared on the requirement's target or an ancestor target. In this case,
the `factors` references are secondary factors.

A requirement declared directly under a target **MUST** declare one or more
factor references under `factors`; each referenced factor **MUST** resolve to a
factor declared on that target or an ancestor target. In this case, the
`factors` references are direct factor references, not secondary factors.

`qualitymd lint` **MUST** emit a `missing-factor-reference` error for each
requirement that is declared directly under a target and has no non-empty scalar
entry under `factors`. Missing `factors`, `factors: null`, `factors: []`, and
sequences with only null or empty entries all leave the requirement without a
factor reference.

> Rationale: The repair requires choosing model intent, so the linter should
> reject the incomplete structure but not invent a factor reference. — 0028

`qualitymd lint` **MUST** continue to emit an error when a referenced factor
does not resolve to an in-scope factor.

The `missing-factor-reference` rule **MUST NOT** be fixable.

The durable format, lint, public README, and scaffold guidance **SHOULD** use
"reference" or "referenced" when describing the mechanical relationship between
requirements and factors.

"Lens" and related wording **MAY** remain where they are useful shorthand for how
a factor frames a quality concern, but they **SHOULD NOT** be the primary term
for the requirement-to-factor reference mechanics.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` — require every requirement to be connected to at least one
  factor, retain "lens" only as helpful shorthand, clarify that `factors`
  references are secondary only when a requirement is already nested under a
  factor, and state that null or empty factor-reference values do not satisfy a
  direct target-level requirement (per the reference and terminology
  requirements above).
- `specs/cli/lint.md` — add the `missing-factor-reference` error rule,
  record that it is not fixable, state its interaction with null or empty
  optional `factors`, and update `unknown-factor` wording to describe factor
  references rather than only secondary factors (per the lint requirements
  above).

### To delete

None
