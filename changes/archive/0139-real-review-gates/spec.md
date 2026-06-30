---
type: Functional Specification
title: Real Review Gates - functional spec
description: What the UX guide and /quality skill must do so feedback invitations are real gates.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Real Review Gates - functional spec

Companion to the [Real Review Gates](../0139-real-review-gates.md) change case.
This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Agent-mediated UX relies on the user being able to tell whether a sentence is
status, preview, question, or gate. If an agent asks for feedback and then
immediately proceeds, the interaction becomes a false affordance: the user sees a
choice but has no chance to affect the outcome. Direct `QUALITY.md` authoring is
especially sensitive because model edits reshape future evaluation judgment.

## Scope

Covered: shared UX guide doctrine, direct `/quality` model-authoring runtime
guidance, durable skill specs, authoring guide guidance, runtime/spec/docs logs,
release notes, and Markdown verification.

Deferred:

- new command surfaces;
- native UI implementation details; and
- behavior changes to setup's existing final review gate.

## Requirements

### Shared UX guidance

- `docs/guides/agent-mediated-ux.md` **MUST** state that feedback invitations are
  gates: when an agent asks for adjustments, concerns, corrections, or
  confirmation such as `looks good`, it waits for the user before mutating.

  > Rationale: asking for feedback while proceeding removes the agency the
  > sentence appears to offer.
  >
  > Durable spec: none - this is durable docs guidance, not a `specs/` bundle
  > contract.

- `docs/guides/agent-mediated-ux.md` **MUST** distinguish informational previews,
  review gates, and decision gates, including that review and decision gates wait
  for a response before mutation.

  > Durable spec: none - this is durable docs guidance.

- `docs/guides/agent-mediated-ux.md` **MUST** require a quick acknowledgement
  before long non-workflow reads or inference that precede a mutation, including
  direct `QUALITY.md` authoring as an example.

  > Durable spec: none - this is durable docs guidance.

### Direct QUALITY.md authoring

- Direct model authoring in `skills/quality/SKILL.md` **MUST** acknowledge the
  likely `QUALITY.md` model change before long model or guide reads, naming that
  it will inspect the current model and relevant authoring guidance and show the
  intended edit for feedback before changing files.

  > Rationale: direct authoring is not a public workflow with a run frame, so it
  > still needs an immediate visible acknowledgement to avoid a silent runway.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > immediate acknowledgement for direct model authoring before long reads.

- Direct model authoring **MUST** wait for a user response after presenting its
  intent checkpoint. It **MUST NOT** ask what the user wants adjusted and then
  proceed in the same turn.

  > Rationale: the checkpoint is a review gate, not a decorative summary.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > the direct-authoring checkpoint to stop and wait.

- The direct-authoring checkpoint example in `skills/quality/SKILL.md` **SHOULD**
  phrase the feedback ask as "before I make that edit" or equivalent, making the
  wait state concrete.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - align
  > checkpoint wording with the review-gate semantics.

- Direct model authoring **SHOULD** use a review gate for `QUALITY.md` model
  changes that reshape future judgment even when intent seems clear. Clear intent
  may remove follow-up questions, but it **MUST NOT** remove the user's chance to
  review the intended edit before mutation.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - preserve
  > review for judgment-shaping model changes.

### Authoring guide alignment

- `skills/quality/guides/authoring.md` **MUST** say direct edits that invite
  concerns, goals, needs, worries, constraints, or `looks good` wait for that
  response before mutating.

  > Durable spec: modify `specs/skills/quality-skill/guides/authoring.md` - align
  > runtime authoring guidance with real review gates.

## Verification

- Source inspection **MUST** show the UX guide's feedback-gate principle,
  non-workflow acknowledgement pattern, Review Gates section, false-affordance
  example, and checklist entries.
- Source inspection **MUST** show direct model-authoring acknowledgement and wait
  rules in `skills/quality/SKILL.md`.
- Source inspection **MUST** show corresponding durable skill spec requirements.
- Source inspection **MUST** show authoring guide and durable authoring guide
  alignment.
- Markdown formatting checks **SHOULD** pass for touched Markdown files.

## Durable spec changes

Rollup of the per-requirement `Durable spec:` annotations above (those are
authoritative) - the [`specs/`](../../../specs/index.md) bundle and
[`SPECIFICATION.md`](../../../SPECIFICATION.md).

### To add

None.

### To modify

- `specs/skills/quality-skill/quality-skill.md` - require direct model authoring
  to acknowledge before long reads, wait after review checkpoints, avoid
  feedback-while-proceeding wording, and retain review for judgment-shaping model
  edits.
- `specs/skills/quality-skill/guides/authoring.md` - require direct-edit
  guidance to treat feedback invitations as wait-for-response review gates.

### To rename

None.

### To delete

None.
