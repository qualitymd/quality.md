---
type: Functional Specification
title: /quality improve
description: Behavioral component spec for applying confirmed QUALITY.md evaluation recommendations and verifying the result.
tags: [skill, quality, mode, improve]
timestamp: 2026-06-22T00:00:00Z
---

# /quality improve

`improve` is the `/quality` skill mode that evaluates, presents
recommendations, applies one confirmed option, and verifies the result with a
fresh evaluation. It implements the shared contracts in the parent
[/quality skill](../quality-skill.md) spec and owns only the improve-specific
behavior below.

The runtime procedure lives at
[`skills/quality/modes/improve.md`](../../../../skills/quality/modes/improve.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`improve` is selected when the user asks the skill to improve evaluated quality,
apply or act on an evaluation recommendation, or close a specific quality gap.

The mode's purpose is to turn one confirmed recommendation option into a change
and prove the before/after delta against the same QUALITY.md model criteria.
`improve` adds a confirmed apply step to the same evaluation process used by
[`evaluate`](evaluate.md).

## Mutation surface and artifacts

Before explicit confirmation, `improve` may mutate only evaluation artifacts.
After confirmation, it may mutate the artifact class named in the decision
brief: evaluated source, `QUALITY.md`, the quality log, or a combination of
those. It **MUST** make that mutation surface unambiguous before editing.

`improve` **MUST NOT** edit evaluated source files or `QUALITY.md` until the user
explicitly confirms a recommendation and option to apply. It **MUST NOT** treat
an obvious or recommended fix as consent.

When the confirmed change alters the QUALITY.md model, `improve` **MUST** write
a quality-log entry for the coherent model change, cross-linking the source
evaluation run and recommendation when applicable.

## Required flow

`improve` **MUST** first behave as [`evaluate`](evaluate.md): resolve scope and
rigor, inspect history, create a run through the CLI, record assessment,
analysis, recommendation, and report artifacts, and summarize current gaps.

Before applying a recommendation, `improve` **MUST** present a decision brief
that names the proposed action, artifact class, evidence or reason, recommended
option, at least one non-mutating alternative, and done criterion or
verification.

After confirmation, `improve` **MUST** apply the chosen option, create a new
numbered evaluation folder, re-evaluate the affected scope, and compare the new
evidence to the recommendation's done criterion.

## Stop conditions

`improve` inherits all [`evaluate`](evaluate.md) stop conditions before the
recommendation stage.

`improve` **MUST** stop before applying when there is no active recommendation
for the selected scope, the recommendation lacks a verifiable done criterion,
the user declines confirmation, or the proposed mutation surface cannot be made
clear.

If verification after apply is incomplete or the rating does not move, the mode
**MUST** report the result as limited rather than fully confirmed.

## Completion criteria

`improve` is complete when it reports the recommendation, applied option,
changed artifacts, before evidence, after evidence, verification performed,
rating movement when any, remaining gaps or limits, and any quality-log entry
written for model changes.
