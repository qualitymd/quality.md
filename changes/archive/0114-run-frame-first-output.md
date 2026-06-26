---
type: Change Case
title: Run frame as first output
description: Make the run-frame "first output before any tool call" timing rule a shared contract and bring the evaluate workflow into line, with a provisional-value allowance for fields that need a tool to resolve.
status: Done
tags: [skill, quality, ux, agent-mediated]
timestamp: 2026-06-26T00:00:00Z
---

# Run frame as first output

This parent concept captures the why and status; the detail lives in its
children:

- [Functional spec](0114-run-frame-first-output/spec.md) - what the change must
  do.
- [Design doc](0114-run-frame-first-output/design.md) - how it is implemented,
  and why.

## Motivation

The run frame earns its place by letting the user catch a wrong inference before
the skill spends effort or mutates anything (the 0038 rationale). That value only
exists if the frame reaches the user *first* — before the agent reads files and
runs commands. A field run of `setup` showed the failure mode directly (the 0096
annotation): the agent front-loaded CLI checks, repository scans, and the
feedback-log write before flushing any text, so the frame arrived after one to
two minutes of silence. `setup` (0096) and `update` (its Required flow) were each
given an explicit "before any tool call" / "before tool inspection" timing guard.

Two gaps remain:

- **`evaluate` never got the guard.** Its runtime procedure resolves the
  QUALITY.md workspace (step 1) before emitting the run frame (step 2), and the
  durable `evaluate` spec's Required flow does not mention the run frame at all.
  Workspace resolution invites tool calls, so the frame can slip behind a *silent
  runway* of reads and commands — the same failure 0096 fixed for `setup`,
  reproduced here.
- **The rule lives only as per-workflow prose.** The two shared homes — the
  `SKILL.md` dispatcher instruction and the durable spec's `Run frames` section —
  say the frame is emitted "at the start of a public workflow" without requiring
  it be the *first output before any tool call*. The timing rule is re-derived
  per workflow rather than bound once where every workflow inherits it.

The companion design guidance now states the rule:
[Designing agent-mediated UX](../../docs/guides/agent-mediated-ux.md) gained an
**Opening** section naming the run frame's two jobs (intent reflection and path
preview), the first-output timing rule, the silent-runway anti-pattern, and a
provisional-value allowance for fields that genuinely need a tool to resolve.
This case brings the skill contract and the lagging `evaluate` workflow into line
with that guidance.

## Scope

Covered:

- lift the "run frame is the first output, before any tool call" timing rule to
  the shared homes — the `SKILL.md` dispatcher run-frame instruction and the
  durable spec's `Run frames` section — so all public workflows inherit it;
- allow a field that genuinely needs a tool to resolve (notably `evaluate`'s
  scope across many modeled Areas) to render a provisional / `resolving…` value
  in the first-output frame, confirmed in a later message, rather than blocking
  the frame on a tool result;
- reorder the `evaluate` runtime procedure so the run frame is the first output
  before workspace resolution or any other tool call, with a provisional scope
  value allowed; and
- add a run-frame-first requirement to the durable `evaluate` spec's Required
  flow, matching `setup` and `update`.

Deferred / non-goals:

- no change to the run frame's *fields* (header, model file, scope, rigor,
  mutation, artifacts, next gate are unchanged);
- no behavior change to `setup` or `update`, which already satisfy the timing
  rule; no vocabulary churn there;
- the internal term "run frame" is retained;
- no CLI, Go, format-schema, rating, roll-up, evaluation-record, or report
  change.

## Affected artifacts

### Code

- [ ] None - no Go, CLI, or generated report implementation change.

### Format spec

- [ ] None - `SPECIFICATION.md` and the QUALITY.md format are unaffected.

### Durable specs

- [x] [`specs/skills/quality-skill/quality-skill.md`](../../specs/skills/quality-skill/quality-skill.md)
      - `Run frames` section: require the frame be the first output before any
      tool call, and allow a provisional / `resolving…` value for a field that
      needs a tool to resolve.
- [x] [`specs/skills/quality-skill/workflows/evaluate.md`](../../specs/skills/quality-skill/workflows/evaluate.md)
      - Required flow: require the run frame be emitted as the first output
      before tool inspection, matching `setup` and `update`.

### Durable docs / bundled skill

- [x] [`docs/guides/agent-mediated-ux.md`](../../docs/guides/agent-mediated-ux.md)
      - already synced ahead of this case: the new **Opening** section states the
      rule this case enforces.
- [x] [`skills/quality/SKILL.md`](../../skills/quality/SKILL.md)
      - dispatcher run-frame instruction: add the first-output/pre-tool timing
      rule and the provisional-value allowance.
- [x] [`skills/quality/workflows/evaluate.md`](../../skills/quality/workflows/evaluate.md)
      - reorder so the run frame is the first output before workspace resolution;
      allow a provisional `Scope: resolving…` value.

### Suggested new durable specs

- None. The existing `/quality` skill specs already own the run-frame contract.

## Status

`Done`. The first-output/pre-tool timing rule and the provisional-value
allowance are now stated in the shared homes — the `SKILL.md` dispatcher
instruction and the durable spec's `Run frames` section — and the lagging
`evaluate` workflow is brought into line: its runtime procedure emits the frame
first (using the invocation-derived model path and a provisional `Scope:
resolving…`) before workspace resolution, and its durable spec's Required flow
states the requirement. `setup` and `update` were already compliant and are
unchanged. The companion guide (`docs/guides/agent-mediated-ux.md`) was synced
ahead with the **Opening** section. Verified with `mise run fmt-md-check`.
Archived on landing. See the [status lifecycle](../index.md#status-lifecycle).
