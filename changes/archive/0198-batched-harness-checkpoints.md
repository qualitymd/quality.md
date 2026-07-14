---
type: Change Case
title: Batched harness checkpoints
description: Let the harness evaluator dispatch a batch of dependency-ready judgment requests per checkpoint so the invoking agent can fan them out to subagents, removing the concurrency-1 tax on harness-backed runs.
status: Done
tags: [evaluation, evaluator, agents, concurrency]
timestamp: 2026-07-11T00:00:00Z
---

# Batched harness checkpoints

## Motivation

Harness-backed evaluation is pinned to concurrency 1. That is not a
configuration choice: the harness evaluator declares `Concurrent: false`, and
the runner clamps any requested concurrency above 1 down to 1 for it. The cause
is the checkpoint transport itself — the runner stages exactly one pending work
request, checkpoints with `status: awaiting_evaluator`, and exits, so there is
never more than one request outstanding for the invoking harness to judge. A
full model run therefore services every bounded judgment work unit strictly one
at a time, even when many of those units are independent and could be judged in
parallel.

The [evaluation concurrency](0195-evaluation-concurrency.md)
work already built a dependency-aware concurrent scheduler and a coordinator
that owns artifact mutation, but it deliberately kept harness-backed runs
sequential "until multi-checkpoint harness work is designed." The
[harness-native evaluator dispatch](0194-harness-native-evaluator-dispatch.md)
case likewise deferred "parallel harness dispatch, native subagent scheduling"
and reserved a `Subagents` capability field (currently `false` everywhere) for
exactly this. This case designs that deferred work.

The value is throughput on the runs users actually launch. Harness is the
skill's default transport because it keeps judgment in the invoking agent's own
session, model, and authentication with no nested agent or extra credential —
the coherence and auth wins that motivated 0194. Batching the checkpoint keeps
all of that while letting the harness fan independent, self-contained judgment
requests out to subagents. The alternative — telling users to switch to an
`api`/`cli` evaluator for speed — sacrifices the very properties that made
harness the default.

## Scope

Covered:

- extend the harness checkpoint transport to keep a bounded **rolling window**
  of dependency-ready work requests outstanding, capped by the run's resolved
  concurrency, emitting each newly-ready request as an earlier one is accepted;
- accept one or more correlated result envelopes per resume call, binding each
  to its request by `requestId`/`inputHash` and validating, hashing, and
  persisting each on the same runner-owned path a single result takes today;
- reuse the existing dependency-ready scheduler to top up the window so the
  runner remains the sole scheduler and parallelism stays bounded by the
  work-graph DAG;
- advertise the harness's subagent-delegation capability and let
  `resolveConcurrency` carry harness concurrency above 1 as the
  outstanding-window cap;
- update the `/quality` evaluate workflow so the checkpoint loop judges the
  outstanding requests, may delegate independent requests to subagents, and
  streams envelopes back as results become ready;
- preserve resume, partial-submission, retry, and interruption behavior over the
  window (a schema-invalid or failed member re-emits for retry, a not-yet-
  submitted member stays outstanding at no retry cost, and neither corrupts
  other members' accepted results); and
- surface resolved concurrency and window width in dry-run previews and run
  receipts.

Deferred:

- a persistent long-lived request/result transport (e.g. one JSON-RPC or JSONL
  subprocess) as an alternative to discrete resume calls;
- adaptive window sizing or subagent-count heuristics beyond the resolved
  concurrency cap;
- native subagent scheduling _by the runner_ (the runner emits the window; how
  the harness fans it out is the harness's concern); and
- `shell`/`manual` evaluators and provider API changes.

Non-goals:

- changing evaluation judgment semantics, accepted payload schemas, rating
  roll-up, generated reports, or run folder layout — the accepted results of a
  batched run MUST match those of the sequential path;
- granting the harness authority to schedule units, widen source, persist
  results, or alter accepted output — the runner still owns the work graph and
  every mutation; and
- adding a command-line concurrency flag (owned by 0195's deferred list).

Implementation depends on 0194 (harness transport) and 0195 (concurrent
scheduler and `evaluation.concurrency`) having landed; both are `Done`.

## Affected artifacts

Derived by sweeping the harness/checkpoint transport, the concurrent scheduler,
and the awaiting-run inspection surfaces; reconciled with the implementation at
In-Review. Substance of durable spec changes lives in the
[functional spec](0198-batched-harness-checkpoints/spec.md); this is the index.

- **Code:**
  - `internal/runner/harness.go`, `internal/runner/engine.go` — pending-request
    state from singular (`PendingEvaluatorCall`, `awaitingRequest`) to a set;
    harness path routed through a dedicated apply/top-up/emit loop reusing the
    concurrent scheduler's `nextReadyEvaluationStep` frontier; multi-envelope
    intake (`loadHarnessResults`) with per-member correlation and retry
    accounting.
  - `internal/runner/concurrent.go` — no edit needed: `nextReadyEvaluationStep`
    was reusable as-is from the harness loop.
  - `internal/evaluator/harness.go` — advertise the subagent-delegation
    capability; the `HarnessResultEnvelope` shape is unchanged (a submission
    carries one or an array of them; intake lives in `internal/runner`).
  - `internal/evaluator/evaluator.go` — no edit needed: the reserved
    `Subagents` capability field already existed.
  - `internal/runner/runner.go` — `resolveConcurrency` no longer clamps an
    evaluator declaring subagent delegation; concurrency becomes the
    outstanding-window cap; the receipt carries `evaluatorRequests` (plural).
  - `internal/cli/evaluation_run.go` — `--evaluator-result` help and awaiting
    receipt rendering for the outstanding set.
  - `internal/runner/artifact.go` — artifact schema version bump (`6 → 7`) for
    multi-outstanding checkpoint state and resume.
  - `internal/runner/dryrun.go` — no edit needed: the preview already carries
    the resolved concurrency, which is the window cap.
  - `internal/evaluation/runner_support.go`, `internal/evaluation/load.go` —
    awaiting-run inspection reads `pendingEvaluatorCalls` and summarizes every
    outstanding request; `internal/status/` needed no edit (lifecycle-only).
  - Tests alongside each of the above (window cap, rolling top-up, partial and
    mixed submission, duplicate-envelope rejection, concurrency parity of
    accepted results, single-request behavior at concurrency 1).
- **Durable specs (see spec's Durable spec changes for the delta):**
  `specs/evaluation/runner.md`, `specs/evaluation/evaluator-contract.md`,
  `specs/evaluation/orchestration.md`, `specs/evaluation/evaluation-json.md`,
  `specs/cli/evaluation-run.md`, `specs/cli/evaluation-status.md`,
  `specs/skills/quality-skill/quality-skill.md`,
  `specs/skills/quality-skill/evaluation.md`,
  `specs/skills/quality-skill/workflows/evaluate.md`, plus bundle `log.md`
  entries. `specs/evaluation/protocol.md` needed no edit — the checkpoint
  transport lives in the runner and orchestration contracts, and protocol
  moves and judgment semantics are unchanged.
- **Format spec:** `SPECIFICATION.md` — no impact (evaluation semantics and
  payload schemas unchanged).
- **Durable docs (bundled skill runtime):**
  `skills/quality/workflows/evaluate.md`,
  `skills/quality/SKILL.md`,
  `skills/quality/resources/cli-workflow-conventions.md`,
  `skills/quality/log.md`, `skills/quality/workflows/log.md`.
- **Generated docs:** `mintlify/cli.mdx` (regenerated `--evaluator-result`
  flag help).
- **Install/scaffold:** none.
- **Changelog:** `CHANGELOG.md` — entry added under Unreleased.

## Children

- [Functional spec](0198-batched-harness-checkpoints/spec.md) — the batch
  checkpoint request/result contract, correlation and validation per member,
  partial-failure and resume behavior, and the durable spec changes.
- [Design doc](0198-batched-harness-checkpoints/design.md) — pending state as a
  set, reuse of the concurrent scheduler's ready-frontier computation to top up a
  rolling window, streamed result intake, partial-submission and
  not-submitted-vs-failed accounting, the capability/cap change, and the skill
  loop.
