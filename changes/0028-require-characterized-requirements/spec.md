---
type: Functional Specification
title: Require characterized requirements
description: Requirements for making uncharacterized requirements invalid and aligning factor-association terminology.
tags: [specification, lint, terminology]
timestamp: 2026-06-18T00:00:00Z
---

# Require characterized requirements

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in IETF RFC 2119.

## Background / Motivation

The format currently permits a requirement declared directly under a target to
inform only the target's local rating, with no factor association. That
"unlensed" special case weakens factor roll-ups and complicates the model:
requirements can be assessed without naming the quality characteristic they
measure. Requiring every requirement to be characterized by at least one factor
makes the model simpler and makes each assessment's quality concern explicit.

The terminology should also distinguish two cases that the current "secondary
factor" wording blurs. A `factors` list on a requirement nested under a factor
adds other factors beyond its containing factor; those listed factors are
secondary. A `factors` list on a requirement declared directly under a target is
not secondary to anything; it is the set of factors that characterize the
requirement.

## Requirements

Every requirement **MUST** be characterized by at least one factor.

> Rationale: A requirement that affects only a target's local rating creates a
> special case with no quality-characteristic roll-up. Mandatory
> characterization keeps every assessment tied to an explicit quality factor. —
> 0028

A requirement nested under a factor or sub-factor **MUST** be characterized by
that containing factor.

A requirement nested under a factor or sub-factor **MAY** list additional
`factors`; each listed factor **MUST** resolve to a factor declared on the
requirement's target or an ancestor target. In this case, the listed factors are
secondary factors.

A requirement declared directly under a target **MUST** list one or more
`factors`; each listed factor **MUST** resolve to a factor declared on that
target or an ancestor target. In this case, the listed factors are not secondary
factors; they are the factors that characterize the requirement.

`qualitymd lint` **MUST** emit an `uncharacterized-requirement` error for each
requirement that is declared directly under a target and has no non-empty
`factors` list.

> Rationale: The repair requires choosing model intent, so the linter should
> reject the incomplete structure but not invent a factor association. — 0028

`qualitymd lint` **MUST** continue to emit an error when a listed factor does not
resolve to an in-scope factor.

The `uncharacterized-requirement` rule **MUST NOT** be fixable.

The durable format, lint, public README, and scaffold guidance **SHOULD** use
"characterize" or "characterized by" when describing the mechanical relationship
between requirements and factors.

"Lens" and related wording **MAY** remain where they are useful shorthand for how
a factor frames a quality concern, but they **SHOULD NOT** be the primary term
for the requirement-to-factor association mechanics.

## Durable spec changes

### To add

None

### To modify

- `SPECIFICATION.md` — require every requirement to be characterized by at least
  one factor, retain "lens" only as helpful shorthand, and clarify that listed
  factors are secondary only when a requirement is already nested under a factor
  (per the characterization and terminology requirements above).
- `specs/cli/lint.md` — add the `uncharacterized-requirement` error rule,
  record that it is not fixable, and update `unknown-factor` wording to describe
  listed factors rather than only secondary factors (per the lint requirements
  above).

### To delete

None
