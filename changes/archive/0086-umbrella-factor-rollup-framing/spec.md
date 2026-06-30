---
type: Functional Specification
title: Umbrella factor roll-up framing — functional spec
description: What the /quality skill guidance must say about how the Agent Harnessability umbrella factor is rated — by rolling up its sub-factors, not by requirements on the parent.
tags: [skill, authoring, factors, harnessability, rollup]
timestamp: 2026-06-24T00:00:00Z
---

# Umbrella factor roll-up framing — functional spec

Companion to the
[Umbrella factor roll-up framing](../0086-umbrella-factor-rollup-framing.md)
change case. This spec states _what_ the `/quality` skill guidance must say about
how the model-wide **Agent Harnessability** umbrella factor (0081) is rated. There
is no design doc: the correction has no design choice to settle.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

0081 introduced Agent Harnessability as a deliberate umbrella factor decomposed
into six independently assessable sub-factors. To prevent silent double-counting,
authors are told not to attach requirements to the umbrella itself. The guidance
compressed that rule into "do not rate the parent directly" (bundled guide) and
"a deliberate umbrella that does not roll up directly" (spec mirror).

Both overstate it. The QUALITY.md format defines no special umbrella semantics. A
factor that has sub-factors and carries no requirements of its own is rated by
rolling up its children — the same mechanism as a grouping area with no local
rating, already specified in
[`SPECIFICATION.md`](../../../SPECIFICATION.md). The umbrella _is_ rated; its
rating just comes from its sub-factors rather than from parent-level requirements.
The overstated phrasing reads as an exemption from rating, and a setup agent
paraphrased it into a generated model as "never rated directly" — a claim that is
simply wrong. The guidance must state the narrow, accurate rule.

## Scope

Covered: the prose in the bundled authoring guide and its durable spec mirror that
describes how the Agent Harnessability umbrella factor is rated.

Deferred / non-goals: no change to the 0081 umbrella structure, its six
sub-factors, the double-count boundary notes, or the QUALITY.md format; no rewrite
of archived Change Cases or append-only logs.

## Requirements

### State the umbrella's rating accurately

The authoring guide and its spec mirror **MUST** describe the Agent Harnessability
umbrella factor as carrying **no requirements of its own**, with its assessment
living in the sub-factors and its overall rating coming from rolling those
sub-factors up.

The guidance **MUST NOT** state or imply that the umbrella factor is "never rated",
"not rated directly", or "does not roll up". Those phrasings contradict the
format's roll-up semantics, under which a factor with sub-factors and no
own-requirements is rated from its children.

The guidance **SHOULD** anchor the rule to the existing grouping-area behavior — a
node with no local rating whose aggregate reflects only its children — so authors
see it as an instance of a general format rule, not a harnessability special case.

> Rationale: the real discipline is "attach requirements to the sub-factors, not
> the parent, to avoid double-counting," not "the umbrella has no rating." Naming
> the rule accurately keeps it from leaking into generated models as a false
> exemption. — 0086

### Preserve the 0081 decomposition and boundaries

This change **MUST NOT** alter the umbrella's six sub-factors, their definitions,
or the double-count boundary notes from 0081 and 0085. It corrects only the
description of how the umbrella is rated.

> Rationale: the umbrella structure is correct; only the roll-up phrasing was
> wrong. — 0086

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/guides/authoring-md.md` - replace the requirement
  that the guide present Agent Harnessability as "a deliberate umbrella that does
  not roll up directly" with the accurate framing (carries no requirements of its
  own; rated by rolling up its sub-factors), per the rating-accuracy requirement
  above.
- `specs/skills/quality-skill/guides/log.md` - record the contract revision.

### To rename

None

### To delete

None
