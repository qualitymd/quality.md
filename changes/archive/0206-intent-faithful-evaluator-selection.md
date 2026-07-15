---
type: Change Case
title: Intent-faithful evaluator selection
description: Make evaluator selection observable and intent-faithful — probe and report every auto candidate, distinguish verified from assumed authentication, and disambiguate harness-versus-SDK transport in the skill.
status: Done
tags: [evaluation, evaluator, selection, discovery, skill, transport]
timestamp: 2026-07-15T00:00:00Z
---

# Intent-faithful evaluator selection

Status note: case is **Done** and archived. Implementation, deterministic
coverage, durable specs, runtime skill guidance, release notes, R1–R7 review
evidence, and the full repository gate are complete.

## Motivation

An operator's evaluate feedback log (2026-07-15, an external QUALITY.md
workspace) records an interrupted run whose intent was for the CLI to
orchestrate the Claude Agent SDK — a fresh, independent, headless evaluator —
once a Claude runtime was detected. Two mechanics silently defeated that
intent:

1. **Auto discovery is order-biased and under-reports.** `auto` probes `codex`
   first and short-circuits on the first usable runtime: with codex installed
   and authenticated, `claude` is never probed and is absent from
   `evaluatorCandidates` entirely. The operator gets a different evaluator than
   the installed-and-detected Claude runtime, with no signal that claude was
   viable. The durable CLI spec already requires reporting "each considered
   candidate"; the short-circuit sidesteps it by never considering the rest.
2. **The skill's precedence hides the SDK path, unexplained.** The skill
   resolves `harness` (in-session judgment) above CLI `auto`, so "have Claude
   evaluate this" lands on the transport that does _not_ give an independent
   evaluation, and the transport explanation does not surface the ambiguity or
   the config path to the SDK evaluator.

A third gap compounds both: the claude runtime's authentication is _assumed_
(no non-interactive probe), while codex's is verified — an asymmetry visible
only in prose evidence.

Why this matters: evaluator identity is an evidence-basis input to a rating. A
user who believes the Claude SDK produced an assessment, when the in-session
harness (or codex) actually did, reasons about the result on a false premise.
Selection should be explicit, observable, and intent-faithful — the operator
sees which runtimes were detected, understands harness-versus-SDK, and gets the
transport they asked for without reverse-engineering precedence.

## Scope

Covered:

- probe every built-in agent-runtime candidate during `auto` discovery and
  report all of them, with readiness evidence, in dry-run and run receipts —
  selection stays deterministic;
- make the `auto` selection reason name the ordering decision and any usable
  candidates not selected;
- carry each candidate's authentication basis (verified versus assumed) as
  structured candidate data, and use a real non-interactive claude
  authentication probe where the runtime documents one;
- require the skill to disambiguate provider-named evaluator intent (harness
  in-session versus SDK subprocess) before selection, and to name the
  independent-SDK alternative when `harness` is selected by default precedence.

Deferred:

- capability- or preference-based `auto` ordering, and any workspace setting
  that reorders discovery (`evaluation.evaluator` already pins a choice);
- an interactive CLI evaluator picker.

Non-goals:

- changing the `auto` ordering itself (`codex` then `claude` stays
  deterministic) or the skill's four-tier precedence;
- discovering `harness` from `auto`, or inferring a harness from environment
  variables;
- managing or storing runtime authentication; and
- new evaluator kinds or API-key evaluator methods.

## Affected artifacts

Derived from a repository sweep for evaluator selection, discovery ordering,
candidate reporting, authentication probing, and transport explanation.

- **Change record:** this parent, `spec.md`, `design.md`, and `review.md` under
  `changes/0206-intent-faithful-evaluator-selection/`; `changes/index.md` and
  `changes/log.md`.
- **Application code:** `src/application/evaluation-run.ts` — `selectEvaluator`
  probes all built-in candidates, reports them, and states the ordering-decided
  selection reason and per-candidate authentication basis.
- **Service code:** `src/services/host-runtime.ts` — a `claudeAuthenticated`
  probe alongside `codexAuthenticated`, backed by the documented
  `claude auth status --json` command (confirmed during design).
- **Code deliberately unchanged:** `src/adapters/evaluator.ts` (no selection
  logic) and `src/cli/app.ts` (flag surface unchanged).
- **Tests:** `test/application/evaluator-selection.test.ts` — the
  first-usable short-circuit assertions (`candidates` length 1) become
  probe-all assertions; `test/integration/cli.test.ts` dry-run receipt
  coverage; `test/integration/evaluation-provider.test.ts` and
  `test/integration/evaluation-execute.test.ts` discovery stubs gain the claude
  probe if one is added.
- **Durable CLI spec:** `specs/cli/evaluation-run.md` — evaluator-selection
  section: probe-all candidate reporting, selection-reason content, and the
  verified/assumed authentication basis; `specs/log.md` records the revision.
- **Durable evaluation specs:** `specs/evaluation/agent-evaluators.md` — the
  documented-authentication-probe policy; `specs/evaluation/log.md` records the
  revision.
  `specs/evaluation/evaluator-contract.md` is deliberately unchanged: kind and
  transport semantics do not move.
- **Durable skill specs:** `specs/skills/quality-skill/evaluation.md` — the
  transport-disambiguation and default-harness explanation requirements;
  `specs/skills/quality-skill/workflows/evaluate.md` — the evaluate step's
  selection wording; `specs/log.md` and the workflow-local log record the
  revisions.
- **Bundled skill runtime:** `skills/quality/SKILL.md` (evaluator-selection
  order paragraph) and `skills/quality/workflows/evaluate.md` (resolve-and-
  explain step) absorb the disambiguation and explanation duties; local logs.
- **Release notes:** `CHANGELOG.md` records the receipt-shape and skill-contract
  change when the case lands.
- **Format specification and project model:** no `SPECIFICATION.md`,
  `quality.schema.json`, or `QUALITY.md` change; model and evaluation meaning
  are unchanged.
- **Durable docs, scaffold, dependencies:** no README, guide, install, or
  dependency change is planned.

## Children

- [Functional spec](0206-intent-faithful-evaluator-selection/spec.md) —
  discovery, reporting, authentication-basis, and skill transport-explanation
  requirements.
- [Design doc](0206-intent-faithful-evaluator-selection/design.md) —
  build-then-select discovery, the three-valued claude authentication probe,
  the `authenticationBasis` candidate field, receipt-level candidate
  reporting, and the skill's disambiguation question.
- [Review ledger](0206-intent-faithful-evaluator-selection/review.md) — R1–R7
  implementation, durable-contract, integration, and repository-gate evidence.

## Status

`Done`. Every R1–R7 requirement has direct implementation and verification
evidence in the
[review ledger](0206-intent-faithful-evaluator-selection/review.md). The case
and its child bundle are archived together.
