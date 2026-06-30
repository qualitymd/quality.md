---
type: Change Case
title: Umbrella factor roll-up framing
description: Correct guidance that overstated the Agent Harnessability umbrella factor as "not rated directly" / "does not roll up" to the accurate "carries no requirements of its own; its rating rolls up from sub-factors."
status: Done
tags: [skill, authoring, factors, harnessability, rollup]
timestamp: 2026-06-24T00:00:00Z
---

# Umbrella factor roll-up framing

A **Change Case** to correct a misleading bit of guidance about the model-wide
**Agent Harnessability** umbrella factor introduced in 0081. The guidance told
authors the umbrella is "not rated directly" and the spec mirror said it "does
not roll up directly." Both overstate the rule: an umbrella factor with
sub-factors _does_ receive a rating — it rolls up from its children, the same way
a grouping area with no local rating does. What is actually true is narrower: the
umbrella carries **no requirements of its own**; its assessment lives in the
sub-factors, and its rating comes from rolling those up.

This wording leaked into a generated `QUALITY.md`, where the setup agent
paraphrased the guide as `# Agent Harnessability is a deliberate UMBRELLA factor
— never rated directly.` — a stronger, inaccurate claim. Fixing the source
guidance keeps future generated models from repeating it.

Detail lives in:

- [Functional spec](0086-umbrella-factor-rollup-framing/spec.md) - what the
  guidance must say. No design doc: this is a wording correction with no design
  choice to settle.

## Motivation

"Not rated directly" / "does not roll up directly" reads as a special-case rule
that exempts the umbrella from rating. It does not. The format defines no special
umbrella semantics: a factor with sub-factors and no own-requirements simply rolls
up from its children, exactly like the grouping-area rule already in
[`SPECIFICATION.md`](../../SPECIFICATION.md). Describing it as "never rated" invites
two errors — a reader thinking the umbrella's overall quality is invisible in a
report, or an author concluding the evaluator special-cases the factor. The
accurate framing ("carries no requirements of its own; rolls up from its
sub-factors") teaches the real discipline: attach requirements to the
sub-factors, not the parent, to avoid silent double-counting.

## Scope

Covered: correct the umbrella roll-up framing in the bundled authoring guide and
its durable spec mirror so both say the umbrella carries no requirements of its
own and is rated by rolling up its sub-factors. Verify by a full-repo sweep that
no other live file repeats the overstated framing.

Deferred / non-goals: no QUALITY.md format or schema change (the grouping-area /
roll-up semantics in `SPECIFICATION.md` are already correct); no CLI or Go change;
no change to the 0081 umbrella structure or its six sub-factors; no rewrite of
archived Change Cases or append-only logs (0081's `spec.md`/`design.md` keep the
original wording as frozen history).

## Affected artifacts

Derived by sweeping the whole repo (an `Explore` agent plus `grep`) for the
phrases "rated directly", "rate the parent", "roll up directly", "does not roll
up", "never rated", and "no local rating", and for "umbrella" near
harnessability. The sweep found the overstated framing in exactly two live files
(below) and otherwise only in frozen `changes/archive/0081-*` records. Grouped by
kind; empty kinds are deliberate.

### Code

None - this is skill authoring guidance, not CLI or Go behavior.

### Format spec (`SPECIFICATION.md`)

None - the roll-up and grouping-area semantics are already correct; this change
aligns the skill's prose to them, it does not change the format.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0086-umbrella-factor-rollup-framing/spec.md#durable-spec-changes)
section is authoritative. In summary:

- `specs/skills/quality-skill/guides/authoring-md.md` - replace "a deliberate
  umbrella that does not roll up directly" with the accurate framing (carries no
  requirements of its own; rating comes from rolling up the sub-factors).
- `specs/skills/quality-skill/guides/log.md` - record the contract revision.

### Durable docs

None beyond the bundled skill below.

### Bundled skill (`skills/quality/`)

- `skills/quality/guides/authoring.md` - runtime counterpart: replace "do not
  rate the parent directly" with the accurate framing.

### Install / scaffold

None - `internal/scaffold/skeleton.md` carries no umbrella roll-up claim.

### Changelog

None - an internal guidance wording correction with no user-facing format, CLI, or
naming change.

## Children

- [Functional spec](0086-umbrella-factor-rollup-framing/spec.md) - what the
  corrected guidance must say.

## Status

`Done`. Landed and archived. The two durable artifacts (bundled authoring guide
and its spec mirror) were corrected and reconciled with the Affected artifacts
list. A full-repo sweep confirmed no other live file repeats the overstated
framing; the only remaining occurrences are the frozen
`changes/archive/0081-harnessability-factor/` `spec.md` and `design.md`, which
stay unchanged as historical record. No code, format-spec, or CLI change was
needed.
