---
type: Change Case
title: Durable spec alignment
description: Align durable specs with the artifact-spec versus behavioral-component spec guidance.
status: Done
tags: [specs, quality-skill]
timestamp: 2026-06-22T00:00:00Z
---

# Durable spec alignment

This case brings the durable specs into alignment with the updated spec
granularity guidance: 1:1 artifact specs keep artifact-oriented names, while
durable behavioral components get capability-oriented specs.

- [Functional spec](0052-durable-spec-alignment/spec.md) - what the alignment
  must do.

## Motivation

The durable specs currently mix three shapes: broad parent specs, 1:1 artifact
specs, and behavioral components embedded as sections inside larger specs. The
new guidance in
[Writing functional specs](../../docs/guides/write-functional-specs.md#conventions)
settles the distinction, but the cumulative `specs/` bundle has not been
reconciled against it.

The immediate pressure is the `/quality` skill: runtime mode files already exist
as standalone behavior components, but their durable contract mostly lives in the
large parent `/quality` spec. That makes the parent harder to review and hides
mode-specific contracts that deserve isolated review.

## Scope

In scope:

- audit the durable `specs/` bundle against the updated granularity guidance;
- add behavioral component specs for `/quality` modes where the mode behavior is
  independently reviewable;
- keep 1:1 artifact specs named after the artifact they govern;
- update indexes, links, and logs needed to keep the OKF bundle navigable.

Deferred:

- changing runtime skill behavior, CLI behavior, or generated artifact formats;
- renaming already-aligned artifact specs such as report and guide artifact
  contracts;
- adding new spec taxonomy beyond the parent/artifact/component distinction.

## Affected artifacts

### Code

None.

### Format spec

None expected.

### Durable specs

- `specs/skills/quality-skill/quality-skill.md` - narrow the parent spec to the
  shared `/quality` contract and link to child specs instead of carrying full
  mode, evaluation, reporting, and quality-log contracts inline.
- `specs/skills/quality-skill/modes/` - add behavioral component specs for the
  `/quality` modes.
- `specs/skills/quality-skill/evaluation.md` - add a behavioral component spec
  for the cross-mode evaluation workflow.
- `specs/skills/quality-skill/reporting.md` - add a component spec for
  evaluation report and run-artifact behavior.
- `specs/skills/quality-skill/quality-log.md` - add a component spec for the
  convention-first quality log.
- `specs/skills/quality-skill/index.md` - list the new mode specs and preserve
  the distinction between the skill parent spec, component specs, guide artifact
  specs, examples, and mode behavior specs.
- `specs/log.md` - record the durable spec alignment.

### Durable docs

None expected; the governing guide changes already landed independently.

### Bundled skill

None expected; this case aligns durable specs, not runtime prompt files.

### Install, scaffold, and packaging

None.

## Status

`Done`: landed and archived. Added component specs for the `/quality` modes,
evaluation workflow, reporting contract, and quality log; narrowed the parent
skill spec to shared contracts and links; updated general spec-splitting
guidance; verified with `mise run fmt-md-check` and `git diff --check`.
