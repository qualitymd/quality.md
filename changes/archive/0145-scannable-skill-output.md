---
type: Change Case
title: Scannable Skill Output
description: Adopt labeled, five-second-scan output shapes across the /quality runtime skill and workflow guidance.
status: Done
tags: [docs, skill, ux, workflows]
timestamp: 2026-06-27T00:00:00Z
---

# Scannable Skill Output

A **Change Case** capturing the *why* and *status*; the detail lives in its
children:

- [Functional spec](0145-scannable-skill-output/spec.md) - what the case must do.
- [Design doc](0145-scannable-skill-output/design.md) - how it's built, and why.

## Motivation

The shared agent-mediated UX guide now teaches five-second-scan output: labeled
blocks, short bullets, explicit answer paths, and fewer dense paragraphs. The
runtime `/quality` skill still has several templates and closeout instructions
that list required content without prescribing a scannable shape. That leaves
agents likely to emit long prose even when the content is right.

This case adopts the new guidance across the skill's user-facing runtime
surfaces so setup, evaluate, review, improve, update, recommendation follow-up,
direct model authoring, and model-review guidance all produce output that is
easy to scan.

## Scope

Covered:

- Adopt labeled output blocks in the shared UX guide's new scannability section.
- Update `/quality` direct model-authoring checkpoint guidance.
- Update runtime workflow templates for setup, evaluate, review, improve, and
  update closeouts or gates.
- Update recommendation follow-up, top-10 checks, and getting-started guide
  templates.
- Align durable `/quality` skill, workflow, guide, recommendation-follow-up, and
  reporting specs.
- Update logs and release notes.

Deferred:

- Native interaction UI implementation.
- CLI output formatting.
- Generated Evaluation report rendering.
- New public workflow names or command surfaces.

## Affected artifacts

Derived by sweeping for scannable output, review gates, closeouts, summaries,
status-first templates, `Answer`, `Changed`, `Verification`, `Next`, and
workflow/guidance specs.

**Code**

- [x] No Go or TypeScript code impact expected; this is a docs/spec/skill
      guidance change.

**Durable specs** (substance in the [functional spec](0145-scannable-skill-output/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md`
- [x] `specs/skills/quality-skill/workflows/setup.md`
- [x] `specs/skills/quality-skill/workflows/evaluate.md`
- [x] `specs/skills/quality-skill/workflows/review.md`
- [x] `specs/skills/quality-skill/workflows/improve.md`
- [x] `specs/skills/quality-skill/workflows/update.md`
- [x] `specs/skills/quality-skill/guides/recommendation-follow-up-md.md`
- [x] `specs/skills/quality-skill/guides/top-10-quality-md-checks-md.md`
- [x] `specs/skills/quality-skill/guides/getting-started-md.md`
- [x] `specs/skills/quality-skill/reporting.md`

**Durable docs / bundled skill runtime**

- [x] `docs/guides/agent-mediated-ux.md`
- [x] `docs/log.md`
- [x] `skills/quality/SKILL.md`
- [x] `skills/quality/workflows/setup.md`
- [x] `skills/quality/workflows/evaluate.md`
- [x] `skills/quality/workflows/review.md`
- [x] `skills/quality/workflows/improve.md`
- [x] `skills/quality/workflows/update.md`
- [x] `skills/quality/guides/recommendation-follow-up.md`
- [x] `skills/quality/guides/top-10-quality-md-checks.md`
- [x] `skills/quality/guides/getting-started.md`
- [x] `skills/quality/log.md`
- [x] `skills/quality/workflows/log.md`
- [x] `skills/quality/guides/log.md`
- [x] `specs/log.md`
- [x] `specs/skills/quality-skill/workflows/log.md`
- [x] `specs/skills/quality-skill/guides/log.md`
- [x] `CHANGELOG.md`

No planned impact: `SPECIFICATION.md`, CLI command specs, setup/evaluate feedback
log schemas, generated schemas, install docs, or README.

## Status

`Done`. The shared UX guide now defines five-second-scan output; runtime skill
guidance, durable skill specs, workflow and guide templates, logs, and release
notes are aligned around labeled, scannable output shapes. `mise run
fmt-md-check` and `mise run check` pass.
