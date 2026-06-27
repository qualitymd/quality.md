---
type: Functional Specification
title: Casual Review Gate Wording - functional spec
description: What the UX guide and /quality skill must do so review gates use simple planned-change and value-prop wording.
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Casual Review Gate Wording - functional spec

Companion to the
[Casual Review Gate Wording](../0140-casual-review-gate-wording.md) change case.
This spec states *what* the change must do.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / Motivation

Review gates need to feel lightweight enough for routine authoring while still
being real gates. The existing guidance correctly requires the agent to wait, but
the checkpoint example reads more like a formal artifact summary than a
conversation. The user should quickly see what the agent plans to change, why
that helps, and how to redirect it.

## Scope

Covered: shared UX guide review-gate wording, direct `/quality`
model-authoring runtime guidance, durable skill specs, authoring guide guidance,
runtime/spec/docs logs, release notes, and Markdown verification.

Deferred:

- new command surfaces;
- native UI implementation details; and
- behavior changes to setup's existing final review gate.

## Requirements

### Shared UX guidance

- `docs/guides/agent-mediated-ux.md` **MUST** say that review gates for content
  or model edits should state the planned change in simple conversational prose
  and include the value prop, preferably using a `so that` clause when it fits.

  > Rationale: the plan and the benefit are the user's fastest handles for
  > correcting scope, naming, boundaries, or intent.
  >
  > Durable spec: none - this is durable docs guidance, not a `specs/` bundle
  > contract.

- `docs/guides/agent-mediated-ux.md` **SHOULD** show the default review-gate
  shape as:

  ```text
  Here's what I'm planning to do:

  <simple common-sense prose of the change>, so that <value prop>.
  ```

  > Durable spec: none - this is durable docs guidance.

- `docs/guides/agent-mediated-ux.md` **MUST** keep the 0139 rule that a feedback
  invitation is a real gate: after asking what the user wants adjusted, the agent
  stops and waits before mutating.

  > Durable spec: none - this is durable docs guidance.

- `docs/guides/agent-mediated-ux.md` **SHOULD** reserve numbered planned-action
  lists for edits where a prose sentence would obscure multiple independent
  actions.

  > Durable spec: none - this is durable docs guidance.

### Direct QUALITY.md authoring

- Direct model authoring in `skills/quality/SKILL.md` **MUST** present a
  pre-mutation checkpoint that states the inferred intent, then uses a simple
  "Here's what I'm planning to do" block to describe the planned change and its
  value prop.

  > Rationale: direct authoring should be conversational and inspectable without
  > becoming a heavy decision brief for every routine model edit.
  >
  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - require
  > planned-change and value-prop wording in the direct-authoring checkpoint.

- The checkpoint **MUST** keep important boundaries and the quality-log decision
  visible when relevant, and **MUST** invite adjustment with welcoming,
  inquisitive wording that makes `looks good` a valid approval path.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - preserve
  > boundaries, quality-log routing, and explicit approval path.

- After presenting the checkpoint, direct model authoring **MUST** stop and wait
  for the user's response before mutating.

  > Durable spec: modify `specs/skills/quality-skill/quality-skill.md` - preserve
  > the real review-gate rule from 0139.

### Authoring guide alignment

- `skills/quality/guides/authoring.md` **MUST** tell agents to preserve and state
  the planned change and value prop for direct edits before inviting feedback.

  > Durable spec: modify `specs/skills/quality-skill/guides/authoring.md` - align
  > runtime authoring guidance with planned-change/value-prop review gates.

## Verification

- Source inspection **MUST** show the shared UX guide's planned-change/value-prop
  review-gate shape, welcoming feedback ask, real wait rule, and prose-first
  guidance.
- Source inspection **MUST** show the same checkpoint shape in
  `skills/quality/SKILL.md`.
- Source inspection **MUST** show corresponding durable skill spec and authoring
  guide requirements.
- Markdown formatting checks **SHOULD** pass for touched Markdown files.
