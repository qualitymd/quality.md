---
type: Change Case
title: Skill rigor and efficiency
description: Upgrade the /quality skill prompt to make evaluations more rigorous, reliable, and efficient, independent of the CLI work.
status: Done
tags: [skill, evaluation, rigor]
timestamp: 2026-06-17T00:00:00Z
---

# Skill rigor and efficiency

Sharpen how the [`/quality` skill](../../specs/skills/quality-skill/quality-skill.md)
runs evaluations: operationalize effort, demand verified evidence, re-check the
findings that bind the headline rating, and emit artifacts efficiently.

- [Functional spec](0017-skill-rigor-efficiency/spec.md) - what the change must
  do.
- [Design doc](0017-skill-rigor-efficiency/design.md) - how the skill prompt
  makes it so.

This change touches the skill prompt only. It needs no CLI changes and stands on
its own, independent of the CLI-side work (0012–0016); it remains valuable even
if that work never lands.

## Motivation

A real subject-evaluation run with the skill was slow and under-specified on
rigor. The judgment was cheap; the cost and risk were in process. The named
effort levels (`quick`/`standard`/`deep`) are not operationalized, evidence can
be asserted from memory rather than verified, the finding that drives the
headline rating is not independently re-checked, and records are emitted serially
when they could be batched.

## Scope

Covered — requirements on the skill's evaluation behavior:

- Operationalize the effort levels in observable, decidable terms.
- Evidence rigor: claims about code/CLI/tool behavior verified by a cited command
  or search; finding locators pinned to verifiable positions.
- Re-check the one or two findings that bind the headline rating before they
  drive the report.
- Recommendation actionability: include route hints when the owning path,
  package, workflow, or maintainer surface is inferable from evidence.
- Execution efficiency: batch independent artifact writes rather than emitting
  them serially.
- Optional subagent fan-out at `deep` effort, with roll-up judgment retained by
  the orchestrating skill.

Deferred:

- Any CLI changes. Where a behavior is later superseded by CLI-written records
  (0012–0016), the spec may note it but does not depend on it.

## Affected specs & docs

Updated during `In-Progress`, before this change reaches `In-Review` (not now):

- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) — operationalize
      effort, require verified evidence and pinned locators, re-check the
      rating-binding findings, batch artifact writes, allow `deep` subagent fan-out.
- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      — sync the durable skill spec with the operationalized effort, evidence rigor,
      rating-binding re-check, batched writes, and `deep` fan-out rules.
- [x] [`specs/skills/quality-skill/examples/`](../../specs/skills/quality-skill/examples/)
      — refresh the reference artifacts so their plan/evidence/report trail
      demonstrates the stricter evaluation behavior.

## Status

`Done`. Implemented and archived after operationalizing effort levels, evidence and pinned-locator rigor, the rating-binding re-check, batched writes, and confined deep fan-out in the skill.
