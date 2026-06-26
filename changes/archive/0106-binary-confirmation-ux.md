---
type: Change Case
title: Binary confirmation UX
description: Make /quality mutation gates use y/n for true binary confirmations while preserving numbered choices for multi-option prompts.
status: Done
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Binary confirmation UX

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0106-binary-confirmation-ux/spec.md) - what the change must
  do.
- [Design doc](0106-binary-confirmation-ux/design.md) - how it is implemented,
  and why.

## Motivation

The shared [agent-mediated UX guide](../../docs/guides/agent-mediated-ux.md) now
distinguishes two interaction shapes:

- small non-binary closed-choice prompts use numbered options with `1` as the
  shortest accept path; and
- true binary confirmations, especially mutation gates, use visible `y`/`n`
  responses because the question naturally means yes or no.

The bundled `/quality` skill still carries older wording that treats all small
closed choices as `1`-first, and several mutation decision-brief templates omit
an `Answer` line. That omission produced an update prompt that asked users to
`reply 1 to apply, or skip`, even though the decision was a binary yes/no
confirmation.

This case aligns the durable `/quality` skill specs and runtime guidance with
the clarified UX guide while preserving numbered responses for scope choices,
setup calibration questions, and apply-vs-handoff routing.

## Scope

Covered:

- add the binary-confirmation exception to the durable and runtime shared
  `/quality` interaction contract;
- add visible `y`/`n` answer paths to true binary mutation gates in update,
  recommendation follow-up, and setup's fallback existing-file edit gate;
- preserve numbered responses for multi-option choices and routing prompts; and
- update append-only skill/spec logs as needed for the touched durable specs and
  bundled skill guidance.

Deferred / non-goals:

- no CLI, Go, format-schema, rating, roll-up, evaluation-record, report-rendering,
  or `QUALITY.md` format change;
- no change to setup's discovery dimensions, model authoring semantics, review
  gate fast path, or evaluation judgment semantics;
- no change to issue-tracker mechanics or update installation mechanics; and
- no replacement of the agent-mediated UX guide beyond the already-applied narrow
  clarification.

## Affected artifacts

### Code

- [ ] None - no Go, CLI, or generated report implementation change expected.

### Format spec

- [ ] None - `SPECIFICATION.md` and the QUALITY.md format are unaffected.

### Durable specs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - add the binary-confirmation exception to the shared interaction contract
      and decision-brief requirements.
- [x] [`specs/skills/quality-skill/workflows/update.md`](../../specs/skills/quality-skill/workflows/update.md)
      - require the update-plan confirmation brief to show `y`/`n`.
- [x] [`specs/skills/quality-skill/recommendation-follow-up.md`](../../specs/skills/quality-skill/recommendation-follow-up.md)
      - require `y`/`n` for local apply and external issue creation gates.
- [x] [`specs/skills/quality-skill/workflows/setup.md`](../../specs/skills/quality-skill/workflows/setup.md)
      - require `y`/`n` on the fallback existing-file update brief, while
      preserving the final review's `looks good` fast path.

### Durable docs / bundled skill

- [x] [`docs/guides/agent-mediated-ux.md`](../../docs/guides/agent-mediated-ux.md)
      - already clarified that true binary confirmations use visible `y`/`n`.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md) - mirror the
      binary-confirmation exception and decision-brief answer guidance.
- [x] [`skills/quality/workflows/update.md`](../../skills/quality/workflows/update.md)
      - add `y`/`n` answer guidance to the update-plan confirmation brief.
- [x] [`skills/quality/guides/recommendation-follow-up.md`](../../skills/quality/guides/recommendation-follow-up.md)
      - add `y`/`n` answer guidance to local apply and issue-creation briefs.
- [x] [`skills/quality/workflows/setup.md`](../../skills/quality/workflows/setup.md)
      - add `y`/`n` answer guidance to the fallback existing-file update brief.

### Suggested new durable specs

- None. The existing `/quality` skill specs already own the affected interaction
  contracts.

## Status

`Done`. Implemented across durable skill specs and bundled runtime skill
guidance, verified with `mise run fmt-md-check`, and archived. The
agent-mediated UX guide clarification is included in the same landed change. See
the [status lifecycle](../index.md#status-lifecycle).
