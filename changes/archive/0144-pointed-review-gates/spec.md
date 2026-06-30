---
type: Functional Specification
title: Pointed Review Gates - functional spec
description: What shared UX guidance and the /quality skill must do so review gates infer purpose and ask for reaction to the consequential assumption.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Pointed Review Gates - functional spec

Companion to the [Pointed Review Gates](../0144-pointed-review-gates.md) change
case. This spec states _what_ the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Review gates should let users correct both the planned edit and the reason the
agent thinks the edit matters. A generic "anything to adjust?" prompt is often
too diffuse to elicit useful steering. When the agent can infer a purpose or a
consequential scope/risk assumption, the gate should make that assumption easy to
accept or reject.

Change Case authors also need explicit process guidance to read the guides that
govern their current artifact and phase, so the case's spec, design,
implementation, and review notes are grounded in the relevant local contracts.

## Scope

Covered: shared UX review-gate guidance, direct `/quality` model-authoring
runtime guidance, durable skill specs, authoring guide guidance, Change Case
process guidance, runtime/spec/docs logs, release notes, and Markdown
verification.

Deferred:

- native review UI implementation details;
- new command or workflow surfaces; and
- behavior changes to setup's existing final review gate.

## Requirements

### Shared UX guidance

- `docs/guides/agent-mediated-ux.md` **MUST** say that content and model review
  gates should state the inferred intent and why the change appears needed before
  describing the planned change.

  > Rationale: the user's fastest correction path is often "yes, that is why" or
  > "no, the purpose is different," not a mechanical edit adjustment.
  >
  > Durable spec: none - this is durable docs guidance, not a `specs/` bundle
  > contract.

- `docs/guides/agent-mediated-ux.md` **MUST** say that review gates should ask
  the user to react to the most consequential inferred scope, risk, naming, or
  boundary assumption when such an assumption is visible.

  > Rationale: a pointed steering axis elicits a better correction than a broad
  > catch-all prompt when one assumption would materially change the edit.
  >
  > Durable spec: none - this is durable docs guidance.

- `docs/guides/agent-mediated-ux.md` **SHOULD NOT** end its model/content review
  gate example with only a generic adjustment prompt when a concrete assumption
  can be named.

  > Durable spec: none - this is durable docs guidance.

### Direct QUALITY.md authoring

- Direct model authoring in `skills/quality/SKILL.md` **MUST** present a
  pre-mutation checkpoint that states the inferred intent, inferred purpose or
  reason the change appears needed, planned change, value prop, important
  boundaries, and quality-log decision when relevant.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > inferred purpose/reason in the direct-authoring checkpoint.

- Direct model authoring in `skills/quality/SKILL.md` **MUST** ask the user to
  react to the most consequential inferred scope, risk, naming, or boundary
  assumption when one is visible, and **SHOULD NOT** use only a generic
  adjustment prompt when a narrower steering axis would better expose the
  assumption most likely to change the edit.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > pointed review-gate steering for direct model authoring.

- After presenting the checkpoint, direct model authoring **MUST** stop and wait
  for the user's response before mutating.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - preserve
  > the real review-gate rule.

### Authoring guide alignment

- `skills/quality/guides/authoring.md` **MUST** tell agents to preserve inferred
  purpose and the consequential assumption for direct edits before inviting
  feedback.

  > Durable spec: modify `specs/skills/quality-skill/guides/authoring.md` - align
  > runtime authoring guidance with pointed review gates.

### Change Case guidance

- `docs/guides/work-with-change-cases.md` **MUST** instruct authors to identify
  and read applicable repository guidance at the start of a Change Case and
  before status advances.

  > Durable spec: none - this is durable process guidance.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` - direct model authoring
  checkpoints state inferred purpose/reason and ask for reaction to the
  consequential assumption (per the Direct QUALITY.md authoring requirements).
- `specs/skills/quality-skill/guides/authoring.md` - direct edit guidance
  preserves inferred purpose and consequential assumptions (per the Authoring
  guide alignment requirement).

### To rename

None

### To delete

None

## Verification

- Source inspection **MUST** show the shared UX guide's purpose-first,
  consequential-assumption review-gate guidance and updated Security example.
- Source inspection **MUST** show matching checkpoint guidance in
  `skills/quality/SKILL.md`.
- Source inspection **MUST** show corresponding durable skill spec and authoring
  guide requirements.
- Source inspection **MUST** show applicable-guidance consultation in
  `docs/guides/work-with-change-cases.md`.
- Markdown formatting checks **SHOULD** pass for touched Markdown files.
