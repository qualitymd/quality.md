---
type: Change Case
title: Setup intro preview
description: Make /quality setup start with useful orientation and an early setup preview before discovery, while preserving setup's review gate and narrow mutation boundary.
status: Done
tags: [skill, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Setup intro preview

A **Change Case** to improve the opening moments of `/quality setup`. Setup
should teach enough to orient first-time users, then quickly show what the agent
has inferred from the project before asking discovery questions or writing
`QUALITY.md`.

Detail lives in:

- [Functional spec](0096-setup-intro-preview/spec.md) - what setup must present
  up front, how the setup preview works, and how feedback-log timing changes.
- [Design doc](0096-setup-intro-preview/design.md) - how setup adds orientation,
  preview, and delayed feedback-log creation without changing CLI behavior.

## Motivation

The current setup workflow has the right ingredients, but its order makes the
first meaningful user value arrive late. It verifies tooling, emits a process
frame, creates a feedback log, reads authoring guidance, scans context, and
builds the setup brief before the user sees project-specific insight.

Setup should instead make the product concept legible immediately and then
surface the context scan as a preview: likely root, domain, visible evidence,
candidate shape, missing context, and the next calibration step. That preview
turns discovery into review and correction of a useful draft, not a cold
questionnaire after an opaque internal phase.

## Scope

Covered:

- Add an educational setup opening that explains what QUALITY.md gives teams,
  AI assistants, and coding agents.
- Require setup to present a concise project-specific setup preview after the
  bounded context scan and before discovery questions.
- Preserve the existing discovery questions, teaching copy, human context
  checkpoint, review gate, write step, lint, and important-gap closeout.
- Adjust setup feedback-log timing so the log can be created after the first
  meaningful setup preview instead of before context analysis.
- Align the bundled skill runtime guidance and durable skill spec mirrors.

Deferred / non-goals:

- Do not reduce setup's core discovery dimensions or remove their teaching copy.
- Do not add a new public `/quality` command or CLI surface.
- Do not change `QUALITY.md` schema, rating semantics, evaluation records,
  report generation, or Go code.
- Do not make setup create evaluation artifacts, quality-log entries, issues,
  integrations, or automations.

## Affected artifacts

Derived by analysis: searched live setup runtime guidance, durable skill workflow
specs, feedback-log specs, setup-related logs, README setup copy, and the
agent-mediated UX guide for setup framing, run-frame, preview, discovery, and
feedback-log timing. Historical archive entries remain frozen except for this
case once it lands.

### Code

None expected - bundled skill guidance and durable spec mirrors only.

### Format spec (`SPECIFICATION.md`)

None - no QUALITY.md schema, rating-scale, roll-up, or evaluation-semantic
change.

### Durable specs (`specs/`)

The functional spec's
[Durable spec changes](0096-setup-intro-preview/spec.md#durable-spec-changes)
section is authoritative. In summary:

- [x] `specs/skills/quality-skill/workflows/setup.md` - add opening orientation,
      setup preview, and revised feedback-log timing.
- [x] `specs/skills/quality-skill/workflows/setup/feedback-log.md` - allow setup
      feedback-log creation after setup preview while preserving required
      finalization behavior.
- [x] `specs/skills/quality-skill/quality-skill.md` - align the root workflow
      feedback-log contract with setup's delayed feedback-log creation and
      early-stop allowance.
- [x] `specs/skills/quality-skill/workflows/log.md` - record the workflow spec
      update if edited.

### Durable docs

- [x] `README.md` - no edit expected; source copy already supports the setup
      intro language.
- [x] `docs/guides/agent-mediated-ux.md` - no edit expected; existing guidance
      already supports status-first setup preview behavior.
- [x] `CHANGELOG.md` - add a user-facing note when the change is implemented.

### Bundled skill (`skills/quality/`)

- [x] `skills/quality/workflows/setup.md` - primary runtime workflow update.
- [x] `skills/quality/SKILL.md` - align the root feedback-log hard rule with
      setup's delayed feedback-log creation and early-stop allowance.
- [x] `skills/quality/workflows/log.md` - record the bundled workflow update.

### Install / scaffold

None expected - install flow and scaffolded `QUALITY.md` content are unchanged.

## Children

- [Functional spec](0096-setup-intro-preview/spec.md) - what setup must present
  up front, how the setup preview works, and how feedback-log timing changes.
- [Design doc](0096-setup-intro-preview/design.md) - how setup adds orientation,
  preview, and delayed feedback-log creation without changing CLI behavior.

## Status

`Done`. Implemented and archived. Setup now opens with short educational
orientation, performs the first project context pass as an explicit read-only
scan, presents a project-specific setup preview before discovery, and creates the
setup feedback log after that preview when the run continues. Verified with
`mise run fmt-md-check`, `git diff --check`, and `mise run check`. No CLI, Go,
format-schema, rating, roll-up, evaluation-record, or report behavior changed.
