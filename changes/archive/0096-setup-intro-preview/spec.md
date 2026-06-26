---
type: Functional Specification
title: Setup intro preview - functional spec
description: Requirements for making /quality setup open with educational orientation and an early project-specific setup preview before discovery.
tags: [skill, setup, ux]
timestamp: 2026-06-26T00:00:00Z
---

# Setup intro preview - functional spec

Companion to [Setup intro preview](../0096-setup-intro-preview.md). This spec
states what `/quality setup` must present before discovery and how feedback-log
timing changes. The format itself is governed by
[`SPECIFICATION.md`](../../../SPECIFICATION.md). This case changes bundled skill
guidance and durable skill spec mirrors only; it adds no normative format rule
and no CLI behavior.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / Motivation

Setup is educational as well as mechanical. A first-time user should quickly
understand what QUALITY.md is for, why the workflow is reading project context,
and what value the scan produced. When setup waits until discovery questions to
show meaningful content, the workflow can feel like internal process followed by
an interrogation. The first visible interaction should therefore orient the user,
and the first project-specific interaction should be a setup preview that turns
discovery into calibration of an evidence-backed draft.

## Scope

Covered:

- Add a short educational opening to `/quality setup`.
- Present a project-specific setup preview after bounded context analysis and
  before discovery questions.
- Keep existing discovery dimensions and teaching copy.
- Keep the review gate before writing `QUALITY.md`.
- Move setup feedback-log creation from mandatory preflight to after the setup
  preview, while preserving eventual feedback-log creation and finalization for
  setup runs that continue past preflight/context scan.

Deferred / non-goals:

- No removal of the setup discovery questions, human context checkpoint, or
  authored teaching copy.
- No new command, flag, CLI behavior, Go behavior, format-schema rule, rating
  semantics, roll-up behavior, evaluation-record behavior, or report behavior.
- No change to setup's mutation boundary except the already-existing feedback log
  under `.quality/logs/`.
- No marketing-style splash screen, mascot, or decorative branding.

## Requirements

### Open with educational orientation

`/quality setup` **MUST** present a short opening orientation before long-running
context work. The orientation **MUST** explain that QUALITY.md gives AI
assistants, coding agents, and teams a holistic definition of quality tailored to
their project, so they can stay aligned, identify critical risks and issues, and
keep improving.

The opening **MUST** tell the user that the first phase is a read-only context
scan and that they will review before changes are written. It **SHOULD** avoid
marketing slogans, decorative branding, or broad philosophical exposition.

> Rationale: setup should teach enough to make the workflow legible before it
> asks questions. The opening copy should carry the README's value proposition in
> operational language, not introduce a separate brand voice. - 0096

### Present a setup preview before discovery

After bounded repository context analysis and before discovery questions, setup
**MUST** present a concise setup preview. The preview **MUST** include:

- likely root Area or boundary;
- likely domain or quality context;
- visible evidence used for the inference;
- likely candidate model shape, such as initial Areas, Factors, or quality
  concerns;
- missing or not-agent-accessible context;
- the next user action.

The preview **MUST** be framed as draft context for correction, not as confirmed
fact. It **MUST** distinguish high-confidence evidence-backed inferences from
low-confidence or missing context. It **MUST NOT** replace the required discovery
questions, human context checkpoint, or final review gate.

> Rationale: a setup preview gives the user meaningful value before the
> questionnaire. It also makes later questions easier because the user is
> correcting a concrete draft instead of inventing context cold. - 0096

### Preserve discovery and review gates

Setup **MUST** continue to ask or present every required discovery input with its
teaching copy. The setup preview **MAY** reduce repetition by referring back to
the preview's evidence, but it **MUST NOT** silently confirm any discovery answer
or omit the human context checkpoint.

Setup **MUST** continue to recap the final discovery answers and wait for an
explicit review-gate response before writing `QUALITY.md`.

### Adjust feedback-log timing

Setup **MUST** still write a workflow feedback log under `.quality/logs/` for
setup runs that proceed past preflight/context analysis into discovery or
authoring. The feedback log **MAY** be created after the setup preview rather
than during preflight.

When the feedback log is created after the preview, it **SHOULD** record the
creation point and may summarize any material pre-log workflow-experience events,
such as unusually slow context scan, tooling friction, or redaction decisions. It
**MUST NOT** duplicate the setup preview's model content or turn the feedback log
into an evidence store.

If setup stops before the feedback log can be created because CLI support is
missing, preflight fails, or the user stops after the initial read-only context
scan, the absence of a feedback log **MUST NOT** be treated as a failed setup
artifact. Once the feedback log exists, normal update and finalization rules
continue to apply.

> Rationale: the feedback log improves the workflow, but creating and explaining
> it before the user sees value contributes to setup feeling process-heavy. Moving
> creation after the preview preserves the artifact without front-loading it. -
> 0096

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - add the educational opening,
  setup preview stage, and revised feedback-log timing while preserving
  discovery, review, write, lint, and closeout behavior (per all requirements
  above).
- `specs/skills/quality-skill/workflows/setup/feedback-log.md` - allow setup
  feedback-log creation after setup preview and clarify pre-log stop behavior
  (per the feedback-log timing requirement above).
- `specs/skills/quality-skill/quality-skill.md` - align the root workflow
  feedback-log contract with setup's delayed feedback-log creation and early-stop
  allowance (per the feedback-log timing requirement above).
- `specs/skills/quality-skill/workflows/log.md` - record the workflow spec mirror
  change if implementation edits the workflow spec log.

### To rename

None

### To delete

None
