---
type: Change Case
title: Casual Review Gate Wording
description: Make direct QUALITY.md review gates state the planned change and value prop in simple conversational prose.
status: Done
tags: [docs, skill, ux, authoring]
timestamp: 2026-06-27T00:00:00Z
---

# Casual Review Gate Wording

A **Change Case** capturing the _why_ and _status_; the detail lives in its
children:

- [Functional spec](0140-casual-review-gate-wording/spec.md) - what the case
  must do.
- [Design doc](0140-casual-review-gate-wording/design.md) - how it's built, and
  why.

## Motivation

The real review-gate rule is correct, but the example wording is still more
formal than the interaction should feel. A direct `QUALITY.md` authoring
checkpoint should make the planned edit obvious, explain the value in plain
terms, and invite the user to shape the edit before mutation. The checkpoint
should not bury the planned action in prose or sound like a procedural form.

The preferred shape is:

```text
Here's what I'm planning to do:

<simple common-sense prose of the change>, so that <value prop>.

What about that?
```

The skill and shared UX guide should carry that casual structure while preserving
the important rule from 0139: if the agent asks for feedback, it waits.

## Scope

Covered:

- Update the shared agent-mediated UX review-gate guidance to prefer simple
  planned-change prose with an explicit value prop.
- Update `/quality` direct model-authoring runtime guidance to use the same
  casual checkpoint shape.
- Align durable `/quality` skill specs and authoring guide specs.
- Update logs and release notes.

Deferred:

- New CLI support for direct authoring.
- A new public `/quality author` workflow.
- Changes to setup's existing final review gate.
- Harness-specific native review UI behavior beyond existing progressive
  enhancement guidance.

## Affected artifacts

Derived by sweeping for review gates, direct model authoring, `looks good`,
checkpoint wording, and feedback invitations.

**Code**

- [x] No Go or TypeScript code impact; this is a docs/spec/skill guidance change.

**Durable specs** (substance in the [functional spec](0140-casual-review-gate-wording/spec.md))

- [x] `specs/skills/quality-skill/quality-skill.md` - require direct model
      authoring checkpoints to state the planned change and value prop in simple
      prose and keep the welcoming review gate.
- [x] `specs/skills/quality-skill/guides/authoring.md` - align direct-edit
      guidance with planned-change/value-prop review gates.

**Durable docs / bundled skill runtime**

- [x] `docs/guides/agent-mediated-ux.md` - update review-gate doctrine and
      examples.
- [x] `docs/log.md` - record the guide update.
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

`Done`. Implemented and archived after shared UX guidance, `/quality` runtime
skill guidance, durable skill specs, authoring guides, logs, release notes, and
full verification were updated. `mise run fmt-md-check` and `mise run check`
pass.
