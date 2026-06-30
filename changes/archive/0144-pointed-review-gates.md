---
type: Change Case
title: Pointed Review Gates
description: Make review gates infer the user's purpose and ask for reaction to the consequential assumption instead of a generic adjustment prompt.
status: Done
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Pointed Review Gates

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0144-pointed-review-gates/spec.md) - what the case must do.
- [Design doc](0144-pointed-review-gates/design.md) - how it's built, and why.

## Motivation

The review-gate structure introduced in earlier cases correctly makes feedback
requests real gates, but a gate that ends with only "what should I adjust?" can
still feel procedural. When a user asks for a model or content change, the agent
often can infer the purpose and the main assumption that would change the edit.
Naming that purpose and asking the user to react to the consequential assumption
is more useful than a broad catch-all prompt: it lets the user either recognize
their intent or steer the agent before mutation.

Future Change Cases also need to consult the guides that govern the artifact or
phase they are changing, so the process guidance should make that consultation
explicit.

## Scope

Covered:

- Update the shared agent-mediated UX guide so review gates state the inferred
  purpose and ask for reaction to the most consequential assumption.
- Update `/quality` direct model-authoring runtime guidance to use the same
  pointed review-gate shape.
- Align durable `/quality` skill specs and authoring guide specs.
- Update the Change Case guide to require applicable guide consultation before
  phase work and status advances.
- Update logs and release notes.

Deferred:

- Native review UI implementation.
- New CLI support for direct authoring.
- New public workflow names.
- Changes to setup's existing final review gate.

## Affected artifacts

Derived by sweeping for review gates, direct model authoring, planned-change
wording, value prop, `looks good`, `go`, applicable guidance, and Change Case
phase guidance.

**Code**

- [x] No Go or TypeScript code impact; this is a docs/spec/skill guidance change.

**Durable specs** (substance in the [functional spec](0144-pointed-review-gates/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md` - require direct model
      authoring checkpoints to state inferred purpose and ask for reaction to the
      consequential assumption.
- [x] `specs/skills/quality-skill/guides/authoring.md` - align direct-edit
      authoring guidance with pointed review gates.

**Durable docs / bundled skill runtime**

- [x] `docs/guides/agent-mediated-ux.md` - update review-gate doctrine and the
      Security example.
- [x] `docs/guides/work-with-change-cases.md` - require applicable guidance
      consultation during Change Case work.
- [x] `docs/log.md` - record the guide updates.
- [x] `skills/quality/SKILL.md` - update direct model-authoring checkpoint
      wording.
- [x] `skills/quality/guides/authoring.md` - align direct-edit guidance.
- [x] `skills/quality/log.md` - append runtime-history entry.
- [x] `skills/quality/guides/log.md` - append runtime-guide history entry.
- [x] `specs/log.md` - append durable-spec history entry.
- [x] `specs/skills/quality-skill/guides/log.md` - append durable-guide history
      entry.
- [x] `CHANGELOG.md` - note the `/quality` skill UX refinement.

No planned impact: `SPECIFICATION.md`, CLI command specs, setup workflow
mechanics, evaluation artifacts, generated schemas, or install docs.

## Status

`Done`. Implemented and archived after shared UX guidance, Change Case guidance,
`/quality` runtime skill guidance, durable skill specs, runtime authoring guides,
logs, release notes, and verification were updated. `mise run fmt-md-check` and
`mise run check` pass.
