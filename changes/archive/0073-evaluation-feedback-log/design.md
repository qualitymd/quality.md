---
type: Design Doc
title: Evaluation feedback log - design doc
description: Design for replacing evaluation debug logging with an evaluate workflow feedback log.
tags: [skill, evaluation, logging, feedback]
timestamp: 2026-06-23T00:00:00Z
---

# Evaluation feedback log - design doc

Design behind the
[Evaluation feedback log](../0073-evaluation-feedback-log.md) and its
[functional spec](spec.md).

## Context

Evaluation already has a run-local `debug-log.md`, seeded by
`qualitymd evaluation create`, but setup has established a better artifact:
`.quality/logs/<timestamp>-setup-feedback-log.md`. The new design makes evaluate
join that workflow feedback-log family instead of renaming the old run-local
debug file in place.

## Approach

Create a shared durable feedback-log spec at
`specs/skills/quality-skill/workflow-feedback-log.md`. That file owns the common
artifact contract: central `.quality/logs/` location, timestamp/workflow naming,
lifecycle frontmatter, body sections, redaction, and the rule that feedback logs
are local and non-authoritative.

Keep workflow-specific behavior in adopter specs:

- `specs/skills/quality-skill/workflows/setup/feedback-log.md` keeps setup's
  creation/finalization rules and points to the shared contract.
- `specs/skills/quality-skill/workflows/evaluate/feedback-log.md` defines the
  evaluate lifecycle, material evaluation-experience events, relationship to the
  numbered evaluation run, and the evidence/rating boundary.

The runtime `/quality evaluate` workflow becomes responsible for creating
`.quality/logs/<timestamp>-evaluate-feedback-log.md` after the run frame and
before assessment evidence collection. The file starts `in-progress`, records
the evaluation run path after `qualitymd evaluation create` returns it, is
updated for material workflow-experience events, and is finalized on close or
best-effort stop.

`qualitymd evaluation create` stops seeding `debug-log.md`. The CLI still owns
the deterministic evaluation record scaffold: `model.md`, `design.md`, `plan.md`,
and the record directories. The feedback log stays skill-authored, matching
setup; the CLI does not create `.quality/logs/` for evaluate.

Historical `debug-log.md` files remain valid historical artifacts. Live specs
mention them only where compatibility context is useful; no migration is needed.

## Alternatives

**Rename `debug-log.md` to `feedback-log.md` inside the run folder.** Rejected
because it preserves the old placement and creates a second feedback-log family
instead of aligning with setup's central `.quality/logs/` convention.

**Write both a central evaluate feedback log and a run-local feedback log.**
Rejected because two process/feedback artifacts invite drift and force readers to
check both before understanding what happened.

**Keep `debug-log.md` and add only wording guidance.** Rejected because the
current name and minimal seed are the source of the confusion. The new artifact
needs the setup-style lifecycle and sections, not just better prose.

**Move feedback-log creation into the CLI.** Rejected for this change because
setup's feedback log is skill-written and workflow-experience content depends on
agent judgment. The CLI should keep owning deterministic record scaffolding.

## Trade-offs and risks

Removing `debug-log.md` from new run folders is a compatibility change for any
local scripts that expected it. The repo's documented evaluation record history
already permits historical differences, and the artifact was non-authoritative,
so the impact is limited. The changelog and specs should call out that older
runs may still contain `debug-log.md`.

The new evaluate feedback log is not mechanically enforced by the CLI. That is
consistent with setup, but it means the runtime skill instructions must be
precise about creation and finalization.

## Open questions

None for this change.
