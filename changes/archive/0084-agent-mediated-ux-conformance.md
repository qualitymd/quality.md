---
type: Change Case
title: Agent-mediated UX conformance
description: Bring live agent-mediated workflow guidance and durable skill specs into conformance with the new agent-mediated UX guide, emphasizing primary questions and calls to action, clearer progress, decision gates, and closeouts without changing quality semantics or CLI behavior.
status: Done
tags: [skill, ux, agents, docs, workflows]
timestamp: 2026-06-24T00:00:00Z
---

# Agent-mediated UX conformance

A **Change Case** to bring the repo's live agent-mediated workflow surfaces into
conformance with
[Designing agent-mediated UX](../../docs/guides/agent-mediated-ux.md). The guide
now states the interaction standard: the agent is the interface; user-facing
workflow output should make state, primary question or call to action,
recommendation, evidence, mutation boundary, and next action easy to scan.

Detail lives in:

- [Functional spec](0084-agent-mediated-ux-conformance/spec.md) — what must be
  brought into conformance.
- [Design doc](0084-agent-mediated-ux-conformance/design.md) — how the
  conformance pass is applied across the skill, specs, and docs.

## Motivation

The `/quality` workflow already has strong mechanics: run frames, discovery
questions, confirmation gates, feedback logs, and status-first closeouts. The
weakness is presentation consistency. Setup has accumulated several good
interaction rules, but they live locally in the setup workflow rather than as a
repo-wide standard. Evaluate, update, recommendation follow-up, and the parent
skill contract have their own output rules, but they do not yet consistently
require the primary user action to stand out visually or require the surrounding
status/progress fields to follow the new guide.

That inconsistency matters because the agent is the user's interface. A workflow
can be correct and still feel dense, uncertain, or more effortful than necessary
when the main question is buried below explanation, decision gates are formatted
differently across modes, or closeouts make the user reconstruct what changed and
what to do next. A shared agent-mediated UX standard lets contributors improve
the workflow presentation without changing the underlying quality model,
evaluation semantics, or CLI mechanics.

## Scope

Covered:

- Audit live agent-mediated workflow surfaces across the bundled `/quality` skill,
  durable `/quality` skill specs, contributor docs, and public usage docs.
- Bring user-facing workflow guidance into conformance with the new guide:
  status-first framing, bold primary questions or calls to action, bold scanning
  labels, adjacent recommendation/evidence, clear shortest acceptable responses,
  visible progress for multi-step workflows, semantic emoji only, decision gates,
  and closeouts.
- Align the parent `/quality` interaction contract and each public workflow
  (`setup`, `evaluate`, `update`, and recommendation follow-up) so the same
  presentation rules apply across modes.
- Add the agent-mediated UX guide to contributor/agent routing where design of
  agent-run workflows is discussed.
- Verify whether any CLI code or hardcoded examples are actually agent-mediated
  surfaces; update the affected-artifact list before `In-Review` if the audit
  finds code impact.

Deferred / non-goals:

- No change to the QUALITY.md format, model schema, rating semantics, evaluation
  semantics, or CLI command behavior.
- No redesign of CLI human output as CLI UX; that remains governed by
  [Designing CLI interfaces](../../docs/guides/cli-design.md).
- No mass rewrite of historical Change Cases, archived specs, append-only logs,
  or recorded example evaluation artifacts. Historical material may remain in its
  original style.
- No requirement that every agent message use emoji. Emoji remains semantic and
  sparse.
- No implementation of a shared renderer or template library unless the
  implementation phase discovers repeated code that justifies one.

## Affected artifacts

Derived by sweeping for run frames, decision briefs, setup discovery prompts,
checkpoints, closeouts, confirmation gates, progress language, and user-facing
skill output. Empty or verification-only kinds are deliberate.

### Code

- [x] `cmd/` and `internal/` — verify no hardcoded agent-mediated workflow output
      needs conformance. CLI-only output remains governed by the CLI design guide;
      no code impact found.

### Durable specs

- [x] `specs/skills/quality-skill/quality-skill.md` — make the shared user
      interaction contract reference or absorb the agent-mediated UX rules:
      primary CTA/question emphasis, output hierarchy, progress, decision
      briefs, and closeout shape.
- [x] `specs/skills/quality-skill/workflows/setup.md` — align setup-specific
      discovery, checkpoint, review, and closeout presentation with the guide.
- [x] `specs/skills/quality-skill/workflows/evaluate.md` — align evaluation
      run framing, stop responses, report summary, limitations, and next-action
      closeout with the guide.
- [x] `specs/skills/quality-skill/workflows/update.md` — align update planning,
      confirmation, verification, and restart/reload guidance with the guide.
- [x] `specs/skills/quality-skill/guides/recommendation-follow-up-md.md` and
      `specs/skills/quality-skill/recommendation-follow-up.md` — align
      recommendation apply/handoff confirmations and closeouts with the guide.
- [x] `specs/skills/quality-skill/reporting.md` — verify whether user-facing
      evaluation summaries need presentation rules or a pointer to the shared
      interaction contract.
- [x] `specs/log.md`, `specs/skills/quality-skill/workflows/log.md`, and
      `specs/skills/quality-skill/guides/log.md` — record durable spec and guide
      contract updates.

### Format spec

- [x] None — this is a skill/workflow presentation change, not a QUALITY.md
      format change.

### Durable docs (guides, AGENTS, README, bundled skill)

- [x] `docs/guides/agent-mediated-ux.md` — source guide; no missing guidance found.
- [x] `docs/guides/index.md` and `docs/log.md` — already register the new guide;
      no additional guide change needed during this case.
- [x] `AGENTS.md` — add the guide to the Guides table for designing or reshaping
      agent-run workflows.
- [x] `skills/quality/SKILL.md` — align the shared runtime interaction contract
      with the guide.
- [x] `skills/quality/workflows/setup.md` — align runtime setup prompts,
      checkpoints, review gate, and closeout.
- [x] `skills/quality/workflows/evaluate.md` — align runtime evaluation progress,
      stop responses, summaries, limitations, and next action.
- [x] `skills/quality/workflows/update.md` — align runtime update plan and
      confirmation output.
- [x] `skills/quality/guides/recommendation-follow-up.md` — align recommendation
      follow-up decision and completion output.
- [x] `README.md` and `install.md` — verify public setup/evaluate examples do not
      conflict with the guide; update only if they teach stale workflow shape.
- [x] `CHANGELOG.md` — add a release note if the conformance pass changes
      user-visible skill behavior.

## Status

`Done`. Implemented, verified, and archived.
