---
type: Functional Specification
title: Setup open-ended review gate - functional spec
description: Requirements for /quality setup's final review prompt before authoring QUALITY.md.
tags: [quality-skill, setup, ux]
timestamp: 2026-06-23T00:00:00Z
---

# Setup open-ended review gate - functional spec

Companion to
[Setup open-ended review gate](../0071-setup-open-ended-review-gate.md). This
spec states what the change must do; the [design doc](design.md) covers how the
runtime skill and durable spec carry it.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Background / Motivation

The setup final recap is the user's last chance to shape the first
`QUALITY.md` before the skill writes it. A prompt that asks only for
`"looks good"` or corrections protects the authoring gate, but it makes
additional useful context feel like exception handling. Setup should invite the
kind of late, human context that improves the model: priorities, worries,
terminology preferences, edge cases, hidden constraints, and facts the repo does
not show.

## Scope

This change governs the wording and handling of `/quality setup`'s final review
gate after discovery and before authoring. It does not change the discovery
question set, the requirement to wait for a user response, or the QUALITY.md
format.

## Requirements

After presenting the consolidated final review recap, `setup` **MUST** ask for a
user response with friendly, colloquial, open-ended wording.

The final review prompt **MUST** preserve the confirmation fast path by telling
the user that saying `"looks good"` or equivalent confirmation advances to
authoring.

The final review prompt **MUST** invite additional context beyond corrections,
including priorities, worries, wording, edge cases, repo-invisible context, or
anything else the user considers important.

`setup` **MUST** treat additional review-gate input as useful setup context, not
only as a correction to an error.

When the user provides additional review-gate input, `setup` **MUST** incorporate
it into the working setup answers before writing or editing `QUALITY.md`.

`setup` **MUST NOT** require the user to provide substantive comments; explicit
confirmation remains sufficient to proceed.

The runtime setup workflow **SHOULD** use this final review prompt or wording
with materially equivalent meaning:

```text
How's this looking? If it feels right, say "looks good" and I'll write QUALITY.md. If anything else is on your mind, send it over too: priorities, worries, wording, edge cases, things the repo doesn't show, or anything that feels important.
```

## Acceptance criteria

- The runtime setup workflow's final review gate includes friendly,
  open-ended prompt wording.
- The runtime prompt explicitly keeps the `"looks good"` fast path.
- The runtime prompt explicitly invites extra context beyond corrections,
  including priorities, worries, wording, edge cases, repo-invisible context, or
  other important considerations.
- The durable setup workflow spec mirrors the runtime behavior.
- The workflow still waits for a user response before writing or editing
  `QUALITY.md`.
- The workflow still incorporates review-gate corrections and cross-cutting
  comments before authoring.
- The workflow does not require the user to provide substantive comments.
- Markdown formatting/checks pass.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/workflows/setup.md` - update the final review
  recap contract to require friendly, open-ended wording that preserves the
  confirmation fast path while inviting broader last-call setup context.

### To rename

None

### To delete

None
