---
type: Change Case
title: Transport-aware evaluator concurrency
description: Resolve and enforce evaluator concurrency from transport capabilities while keeping CLI scheduling authoritative and agent/subagent judgment bounded.
status: In-Progress
tags: [evaluation, evaluator, concurrency, agents, runner]
timestamp: 2026-07-15T00:00:00Z
---

# Transport-aware evaluator concurrency

Status note: **In-Progress**; R1–R8 passed individual requirement-quality and
requirement-set review, the technical design is settled, and implementation is
underway.

## Motivation

The evaluation architecture has the right authority boundary: the CLI owns the
work graph and authoritative artifacts, while evaluators provide bounded
judgment. Its concurrency implementation does not yet express that boundary
honestly across transports.

Automatic concurrency is derived from host CPU count even though evaluator work
is constrained by provider sessions, quotas, cost, and harness worker slots.
The capability booleans conflate simultaneous SDK calls with harness subagent
delegation. Provider runs are created through the harness path and patched to
the selected evaluator afterward, which can stage more pending requests than
the provider's resolved cap. SDK-backed execution also waits for a whole batch
before submitting any result, leaving a freed slot idle behind the slowest call.

This case makes concurrency a transport-aware runner contract. The runner stays
the sole scheduler and persistence owner; direct SDK sessions and harness
subagents remain isolated workers that receive one dependency-ready request at
a time.

## Scope

Covered:

- replace the ambiguous evaluator concurrency booleans with explicit direct-
  call, delegated-request, automatic-cap, and optional maximum-cap capability
  semantics;
- resolve omitted concurrency from the selected evaluator rather than host CPU
  count, while retaining `evaluation.concurrency` as the workspace hard cap;
- create every run with its selected evaluator, capabilities, and resolved
  concurrency before staging work requests;
- keep direct evaluator calls in a completion-driven bounded worker pool that
  validates and persists each accepted result before topping up the freed slot;
- preserve the harness's rolling outstanding-request window while leaving
  actual direct-versus-subagent execution to the invoking agent;
- make the distinction between a concurrency cap, outstanding harness work,
  and observed active calls explicit in receipts, progress wording, and
  run-local diagnostics;
- advance the runner artifact schema cleanly for the evaluator-capability shape;
  and
- add deterministic concurrency, top-up, ordering, cancellation, retry,
  initialization, and harness-delegation coverage.

Deferred:

- adaptive concurrency from live rate limits, latency, token use, or cost;
- provider quota discovery and dynamic capacity negotiation with an invoking
  harness;
- cache-aware warm-up barriers or cohort scheduling, which remain separate from
  the prompt-shaping and telemetry work in 0203;
- a command-line `--concurrency` override;
- a persistent JSON-RPC or JSONL harness transport; and
- recursive subagent trees inside evaluator workers.

Non-goals:

- changing the evaluation work graph, judgment semantics, evidence contract,
  ratings, recommendation logic, or report content;
- making subagents an evaluation orchestration engine;
- sharing, resuming, or forking provider conversations between work units;
- adding direct model-API evaluators; or
- making provider cache state part of correctness or `evaluation.json`.

## Affected artifacts

Derived from repository searches for concurrency configuration and resolution,
evaluator capabilities, provider and harness dispatch, checkpoint windows,
receipts, artifact manifests, progress wording, and tests.

- **Change record:** this parent, `spec.md`, and `design.md` under
  `changes/0204-transport-aware-evaluator-concurrency/`; `changes/index.md` and
  `changes/log.md`; a later review ledger during implementation; all archived
  together when done.
- **Durable evaluation specs:** `specs/evaluation/runner.md`,
  `specs/evaluation/orchestration.md`,
  `specs/evaluation/evaluator-contract.md`,
  `specs/evaluation/agent-evaluators.md`, and
  `specs/evaluation/evaluation-json.md` absorb runner authority, transport-
  aware resolution, completion-driven scheduling, capability semantics,
  persistence, and schema-version rules; `specs/evaluation/log.md` records the
  revision. `specs/evaluation/protocol.md` has no behavior change because work-
  unit kinds and judgment moves are unchanged.
- **Durable CLI spec:** `specs/cli/evaluation-run.md` updates dry-run and receipt
  concurrency semantics. Other command specs are unaffected because no flag or
  status lifecycle changes.
- **Durable skill specs:** `specs/skills/quality-skill/evaluation.md` and
  `specs/skills/quality-skill/workflows/evaluate.md` clarify that an outstanding
  window is a cap, the harness chooses actual worker fan-out, and each delegated
  worker receives only one bounded request. Their local logs record the
  revision; the parent skill contract needs no new workflow or invocation.
- **Domain code:** `src/domain/evaluator/types.ts` owns the capability shape;
  `src/domain/evaluation/run.ts` and a focused evaluation-domain concurrency
  module own pure resolution, run creation, receipts, schema version, and
  deterministic ready-frontier ordering.
- **Application code:** `src/application/evaluation-run.ts`,
  `src/application/evaluation-execute.ts`,
  `src/application/evaluation-provider.ts`, and
  `src/application/evaluation-resume.ts` move to selected-evaluator-first run
  creation, serialized acceptance, and completion-driven top-up;
  `src/application/evaluation-report.ts` advances the current artifact version
  consumed by report rebuild.
- **Adapter and service code:** `src/adapters/evaluator.ts` declares honest
  built-in transport capabilities; `src/services/host-runtime.ts` no longer
  supplies hardware concurrency for evaluation scheduling if no other consumer
  needs it. Workspace parsing keeps the existing positive
  `evaluation.concurrency` setting and evaluator-profile shape.
- **Tests:** pure resolution/capability tests plus
  `test/domain/evaluation/run.test.ts`,
  `test/application/evaluator-selection.test.ts`,
  `test/integration/evaluation-execute.test.ts`,
  `test/integration/evaluation-provider.test.ts`, and
  `test/integration/cli.test.ts` cover initialization, observed parallelism,
  rolling top-up, ordering, parity, partial submission, retry, resume, and
  cancellation. Deterministic artifact snapshots update for schema version 9.
- **Bundled skill runtime:** `skills/quality/SKILL.md` and
  `skills/quality/workflows/evaluate.md` update concurrency-cap and delegation
  wording; their local logs record the revision.
- **Release notes and generated artifacts:** `CHANGELOG.md` records the clean
  capability/artifact change and automatic-concurrency behavior. Regenerate
  `mintlify/cli.mdx` only if its owning CLI-doc task observes source drift; no
  report-gallery, JSON Schema, or specification-Mintlify generation is expected.
- **Durable docs:** no README, install, contributor, or user-guide change is
  planned; concurrency remains runner configuration surfaced primarily through
  the agent-mediated evaluation workflow.
- **Format specification and project model:** no `SPECIFICATION.md`,
  `quality.schema.json`, or `QUALITY.md` change; model and evaluation meaning are
  unchanged.
- **Dependencies, packaging, and scaffold:** no dependency, installer, release-
  pipeline, or scaffold change is planned.

## Children

- [Functional spec](0204-transport-aware-evaluator-concurrency/spec.md) —
  authority, capability, resolution, dispatch, harness, durability, artifact,
  and verification requirements.
- [Design doc](0204-transport-aware-evaluator-concurrency/design.md) — capability
  model, shared run preparation, coordinator-owned provider pool, harness
  boundary, persistence, and verification design.

## Status

`In-Progress`. The functional spec and design are settled; implementation now
advances the runtime, durable current-behavior specs, skill guidance, and tests
together before review.
