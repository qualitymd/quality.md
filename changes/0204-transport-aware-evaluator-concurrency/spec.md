---
type: Functional Specification
title: Transport-aware evaluator concurrency — functional spec
description: Requirements for runner-owned, transport-aware evaluator concurrency with bounded direct calls and harness delegation.
tags: [evaluation, evaluator, concurrency, agents, runner]
timestamp: 2026-07-15T00:00:00Z
---

# Transport-aware evaluator concurrency — functional spec

Companion to the
[Transport-aware evaluator concurrency](../0204-transport-aware-evaluator-concurrency.md)
Change Case. This spec defines what must change; the
[design doc](design.md) owns the capability representation, shared run-
preparation path, worker coordination, and test seams.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

Evaluation concurrency spans two different transports. SDK-backed evaluators
can keep several fresh provider sessions active in one CLI process. A harness-
backed evaluation checkpoints several dependency-ready requests so the invoking
agent can judge them directly or place them with native subagents. Those are
different execution mechanisms under one invariant: the runner decides which
work is ready and remains the only authority that validates and persists it.

The current runtime derives automatic concurrency from host CPU count, uses
ambiguous capability booleans for direct calls and subagent delegation, creates
provider runs through a harness-shaped bootstrap, and drains a complete SDK
batch before refilling it. These choices make the recorded cap a weak statement
about actual transport capacity and leave avoidable idle time and state-
initialization inconsistencies.

## Scope

This change covers evaluator capability semantics, automatic and configured
concurrency resolution, run initialization, direct-call scheduling, harness
delegation boundaries, accepted-result durability, observability, and
verification.

Deferred are adaptive rate/cost control, provider quota probing, dynamic harness
capacity negotiation, cache-aware warm-up or cohort scheduling, a per-run CLI
concurrency flag, persistent harness transport, and recursive subagent trees.

The work graph, evaluator selection order, workspace inspection policy, payload
schemas, evidence sealing, rating behavior, recommendations, and report content
are unchanged.

## Terms

- **Configured cap** — the positive `evaluation.concurrency` value when present
  in workspace configuration.
- **Automatic cap** — the positive transport default declared by the selected
  evaluator when the configured cap is absent.
- **Maximum cap** — an optional positive upper bound declared by an evaluator
  transport.
- **Resolved concurrency** — the runner-owned cap after applying configuration
  and the selected evaluator's capabilities.
- **Active direct calls** — SDK-backed evaluator calls currently running.
- **Outstanding harness requests** — checkpointed requests awaiting harness
  results; they are eligible for parallel service but are not proof of active
  subagents.

## Requirements

### R1 — Runner scheduling authority

The evaluation runner **MUST** remain the sole owner of work-graph construction,
dependency readiness, concurrency resolution, work-unit identity, retry,
cancellation, result acceptance, evidence sealing, persistence, output ordering,
and report generation.

> Rationale: letting a provider session or subagent own graph progress recreates
> multiple evaluation engines and makes resume and artifact semantics depend on
> the selected harness. — 0192, 0198, 0204
>
> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/orchestration.md`, and
> `specs/evaluation/evaluator-contract.md` — preserve CLI control-plane
> authority while distinguishing direct and delegated execution transports.

An evaluator worker **MUST** receive one runner-issued dependency-ready work
request and return one correlated result or classified failure. It **MUST NOT**
construct or advance the work graph, choose sibling work, write run artifacts,
or delegate evaluation orchestration recursively.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/agent-evaluators.md` — define evaluator and nested-worker
> boundaries; modify `specs/skills/quality-skill/evaluation.md` — preserve the
> harness wrapper boundary.

### R2 — Honest transport capabilities

Every evaluator capability record **MUST** distinguish support for simultaneous
direct evaluator calls from support for servicing multiple checkpointed
requests through external harness delegation. It **MUST** carry a positive
automatic cap and **MAY** carry a positive maximum cap.

> Rationale: simultaneous SDK calls and parent-agent subagent fan-out both allow
> a cap above one, but they imply different scheduling, progress, cancellation,
> and observability behavior. One generic `subagents` boolean cannot describe
> both honestly. — 0204
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` — replace the
> ambiguous concurrency/subagent capability contract; modify
> `specs/evaluation/evaluation-json.md` — update the persisted capability record.

An SDK-backed evaluator that disables nested agents **MUST NOT** advertise
external delegation merely because its provider or host product supports
subagents elsewhere. A harness evaluator **MAY** advertise external delegation
without advertising simultaneous in-process calls.

> Durable spec: modify `specs/evaluation/agent-evaluators.md` — record the
> difference between fresh direct SDK sessions and invoking-harness delegation.

The built-in capabilities **MUST** initially resolve as follows:

| Evaluator | Direct calls | External delegation | Automatic cap | Maximum cap |
| --------- | ------------ | ------------------- | ------------: | ----------: |
| `harness` | no           | yes                 |             4 | unspecified |
| `codex`   | yes          | no                  |             4 | unspecified |
| `claude`  | no           | no                  |             1 |           1 |

Configured Codex or Claude profiles **MUST** inherit the concurrency capability
of their `kind`; this change **MUST NOT** add profile-level concurrency fields.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/agent-evaluators.md` — define built-in and profile
> capability behavior.

### R3 — Transport-aware concurrency resolution

When `evaluation.concurrency` is present, the runner **MUST** treat it as the
requested hard cap and reject a non-positive or non-integer value. When it is
absent, the runner **MUST** use the selected evaluator's automatic cap and
**MUST NOT** derive evaluation concurrency from host CPU count.

> Rationale: provider calls and harness workers are constrained by remote
> capacity, cost, and agent slots rather than local compute threads. — 0204
>
> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/evaluator-contract.md` — replace the CPU-derived automatic
> rule with selected-transport resolution.

If the selected evaluator declares a maximum cap, the runner **MUST** clamp the
requested or automatic cap to that maximum. If it supports neither simultaneous
direct calls nor external delegation, the runner **MUST** resolve concurrency to
`1` regardless of a higher configured value.

> Durable spec: modify `specs/evaluation/runner.md` — define the complete
> resolution rule and clamping behavior.

The public configuration surface **MUST** remain
`evaluation.concurrency`; this change **MUST NOT** add a `--concurrency` flag or
expose scheduler-strategy names.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/cli/evaluation-run.md` — retain the existing configuration and CLI
> boundary.

### R4 — Selected-evaluator-first run creation

Before the runner stages any evaluator-backed request, it **MUST** select the
evaluator, validate its capabilities, resolve concurrency, and create the run
manifest with that evaluator kind, capability record, and resolved cap.

> Rationale: creating a provider run as a harness run and patching it afterward
> can persist a pending window that was sized against the wrong transport. — 0204
>
> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/orchestration.md`, and
> `specs/evaluation/evaluation-json.md` — require immutable selected-evaluator
> setup before dispatch state exists.

At every persisted state boundary, active direct calls or outstanding harness
requests **MUST NOT** exceed the manifest's resolved concurrency. Evaluator
identity, capabilities, and resolved concurrency **MUST NOT** be patched after
requests have been issued.

> Durable spec: modify `specs/evaluation/evaluation-json.md` — strengthen the
> manifest and pending-state invariant.

### R5 — Completion-driven direct dispatch

For a direct-call evaluator, the runner **MUST** keep up to the resolved
concurrency of dependency-ready evaluator calls active. When one call completes,
the runner **MUST** validate and persist its accepted result, advance dependent
deterministic work, and fill the freed slot with newly ready evaluator work
without waiting for unrelated active calls to finish.

> Rationale: waiting for a complete batch turns the slowest call into a barrier
> and leaves capacity idle even though the work graph exposes a ready frontier.
> — 0195, 0204
>
> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md` — require completion-driven top-up for
> direct transports.

Each direct work unit **MUST** continue to use a fresh evaluator session. Call
completion order **MUST NOT** change persisted payload order, ratings, evidence,
reports, or any other deterministic runner projection.

> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/agent-evaluators.md` — preserve fresh-session isolation and
> observational equivalence under completion-driven scheduling.

### R6 — Bounded harness delegation

For a harness evaluator, resolved concurrency **MUST** remain the maximum number
of outstanding dependency-ready checkpoint requests. The runner **MUST** accept
any correlated subset, persist each accepted result, keep unsubmitted requests
outstanding without retry cost, and top the window up as capacity becomes
available.

> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md` — retain the rolling checkpoint contract
> while using the clarified delegated-request capability.

The invoking skill or agent **MAY** judge an outstanding request itself or give
it to one native subagent. A delegated worker **MUST** receive only that
self-contained request and **MUST NOT** receive authority to choose other work,
write run artifacts, run an independent QC pass, or recursively orchestrate the
evaluation.

> Rationale: the agent knows its live worker availability; the CLI knows the
> dependency graph and artifact state. Keeping those authorities separate gains
> native parallelism without making subagent behavior part of evaluation
> correctness. — 0198, 0204
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — define one-request worker
> delegation and prohibit alternate orchestration.

Harness progress **MUST** describe the receipt as up to `N` outstanding requests
or as a concurrency cap. It **MUST NOT** claim that all `N` requests are actively
running concurrently unless the harness has actually dispatched them.

> Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` — make
> progress wording distinguish eligibility from observed activity.

### R7 — Per-result durability, retry, resume, and cancellation

The runner **MUST** serialize result acceptance and persist each accepted work
unit before treating its slot as free or scheduling a dependent work unit.
Evaluator workers **MUST NOT** mutate shared run state directly.

> Rationale: accepted judgment must survive interruption independently of the
> completion order of sibling calls. — 0192, 0198, 0204
>
> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/evaluation-json.md` — preserve single-store, per-result
> durability under both transports.

A retryable failure **MUST** consume only its work unit's retry budget. A
terminal failure or cancellation **MUST** stop new dispatch, interrupt active
direct calls, keep already accepted results durable, and leave unaccepted work
in a valid resumable state under the existing failure taxonomy.

> Durable spec: modify `specs/evaluation/runner.md` and
> `specs/evaluation/orchestration.md` — apply existing retry, cancellation, and
> resume semantics to the completion-driven coordinator.

### R8 — Artifact, observability, and verification

`qualitymd evaluation run --dry-run --json`, run receipts, and
`evaluation.json` **MUST** continue to report resolved `concurrency` as a cap.
Run-local diagnostics **MUST** record the resolution source (configured or
automatic), declared automatic and maximum caps, transport mode, and observed
peak active direct calls or outstanding harness requests without treating those
diagnostics as evaluation results.

> Durable spec: modify `specs/evaluation/runner.md`,
> `specs/evaluation/evaluation-json.md`, and
> `specs/cli/evaluation-run.md` — clarify cap semantics and place detailed
> resolution/activity data in run-local logs.

The runner artifact schema **MUST** advance from version `8` to version `9` for
the new persisted evaluator-capability shape. The runtime **MUST** refuse to
resume an incompatible version-8 in-flight run rather than add a dual reader or
migration shim; completed historical artifacts remain historical data.

> Durable spec: modify `specs/evaluation/evaluation-json.md` — define schema
> version 9 and the clean resume boundary.

Focused tests **MUST** prove the concurrency-resolution matrix, selected-
evaluator-first pending cap, direct-call peak bound, top-up after one early
completion while a sibling remains active, deterministic ordering under
out-of-order completion, concurrency-1 parity, harness partial submission,
per-result durability, retry isolation, cancellation, resume, progress wording,
and schema-version refusal. The full repository gate **MUST** pass before
review.

> Durable spec: none — this is change acceptance evidence, not an enduring
> runtime contract.

## Requirement-set review

R1 fixes the authority boundary. R2 gives each transport enough honest capacity
information for R3 to resolve a cap without CPU heuristics. R4 ensures that cap
and evaluator identity govern the run before any state is issued. R5 applies the
cap to direct sessions without batch barriers; R6 applies it to harness-eligible
work without pretending the CLI controls agent slots. R7 preserves durability
and recovery as calls finish independently. R8 makes the distinction observable
and proves the concurrency invariants.

Together R1–R8 achieve the motivation: concurrency becomes a runner-owned,
transport-aware bound, while evaluators and subagents remain isolated workers.
The set does not reopen evaluator selection, prompt shaping, caching policy,
judgment semantics, evidence, or reports. Every requirement has a direct
verification path, and no unresolved decision remains inline.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/runner.md` — runner authority, transport-derived automatic
  caps, clamping, selected-evaluator-first creation, completion-driven direct
  dispatch, retry/cancellation application, and concurrency diagnostics
  (R1–R5, R7–R8).
- `specs/evaluation/orchestration.md` — direct completion-driven scheduling,
  harness top-up, deterministic ordering, serialized per-result acceptance,
  and recovery (R1, R4–R7).
- `specs/evaluation/evaluator-contract.md` — direct-call versus external-
  delegation capabilities, built-in defaults, profile inheritance, config
  boundary, and worker authority (R1–R3).
- `specs/evaluation/agent-evaluators.md` — honest built-in capabilities, flat
  fresh SDK sessions, no nested delegation, and scheduling isolation (R1–R2,
  R5).
- `specs/evaluation/evaluation-json.md` — immutable evaluator/cap setup, bounded
  pending state, capability record, serialized persistence, schema version 9,
  and resume refusal (R2, R4, R7–R8).
- `specs/cli/evaluation-run.md` — retained configuration/flag boundary and
  resolved-cap semantics in dry-run and receipts (R3, R8).
- `specs/skills/quality-skill/evaluation.md` — runner/worker authority and one-
  request harness delegation (R1, R6).
- `specs/skills/quality-skill/workflows/evaluate.md` — bounded delegation,
  partial submission, and honest outstanding-versus-active progress wording
  (R6).

### To rename

None.

### To delete

None.
