---
type: Functional Specification
title: Evaluation runner
description: CLI-owned deterministic evaluation engine, concurrency, run-local logging, and failure taxonomy.
tags: [evaluation, runner, orchestration]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluation runner

The evaluation runner is the CLI-owned engine behind
[`qualitymd evaluation run`](../cli/evaluation-run.md). It executes the
[evaluation protocol](protocol.md) as the deterministic work graph defined by
the [orchestration contract](orchestration.md), invokes evaluators through the
[evaluator contract](evaluator-contract.md), and persists the
[`evaluation.json`](evaluation-json.md) run artifact.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

Evaluation needs the same orchestration and artifact guarantees whether it is
run through Codex, Claude, or an invoking agent harness. When a skill
orchestrated evaluation, workflow integrity depended on the invoking harness:
every harness had to reconstruct the same protocol from prompt instructions.
The runner lets `qualitymd` own repeatable workflow behavior while preserving
agent- and model-mediated judgment as bounded evaluator work units. — 0192

## Scope

Deferred:

- a persistent long-lived harness request/result transport (for example one
  JSON-RPC or JSONL subprocess) as an alternative to discrete resume calls, and
  adaptive window sizing beyond the resolved concurrency cap;
- custom evaluator/plugin protocols beyond configured Codex and Claude
  agent-runtime profiles;
- raw prompt/response tracing beyond a future explicit opt-in;
- compacting completed-run execution state in `evaluation.json`.

Non-goals:

- making provider APIs mandatory or requiring API keys from CLI subscription
  users;
- making any external orchestration framework mandatory for the core runner;
- allowing prompt-cache hits, provider context, or subagent availability to
  affect evaluation correctness or persisted output semantics;
- allowing evaluators to write run artifacts directly.

## Ownership

The runner **MUST** own run creation, resume, work-graph execution, evaluator
invocation, result validation, persistence, report generation, and the final
receipt.

> Rationale: if any of these phases lives outside the runner, the workflow
> keeps two orchestrators and cannot be deterministic across harnesses. — 0192

The runner builds and schedules the work graph per the
[orchestration contract](orchestration.md), which owns the graph shape,
dependencies, scheduling, retry, resume, and cancellation semantics.

## Result validation

Before accepting an evaluator-backed result, the runner **MUST** normalize the
payload envelope — `schemaVersion`, `kind`, and subject identity — and validate
the payload against the evaluation data contracts in
[payload kinds](records/payload-kinds.md) and the run's `model-snapshot.md`.

The runner **MUST** verify finding-coverage completeness before accepting the
recommendation-ranking result, executing the protocol's
`accountForFindingCoverage` move as acceptance validation.

The runner **MUST** assign canonical `qrec_` recommendation IDs to accepted
recommendations before persisting them.

## Concurrency

The runner **MUST** own concurrency resolution for ready evaluator-backed work
units.

> Rationale: parallelism, native subagents, and provider context reuse are
> useful only while they remain scheduling and transport choices under the same
> work graph. If they become alternate orchestration engines, the portability
> problem returns. `concurrency` is the user-facing control; scheduler strategy
> names are an implementation detail. — 0192, 0195

The workspace config **MAY** set `evaluation.concurrency` in
`.quality/config.yaml`. When present, it **MUST** be a positive integer and is
the requested hard cap. When absent, the runner **MUST** use the selected
evaluator's positive `automaticConcurrency`; host CPU count **MUST NOT** enter
evaluation concurrency resolution.

The runner **MUST** select the evaluator, validate its dispatch capability, and
resolve concurrency before it stages evaluator work or writes the initial run
artifact. It **MUST** clamp the requested or automatic value to an evaluator's
optional `maxConcurrency`. An evaluator that supports neither concurrent direct
calls nor delegated requests **MUST** declare automatic and maximum concurrency
of `1`, and the runner **MUST** resolve it to `1`.

> Rationale: provider sessions and harness worker slots are remote transport
> constraints, not local CPU work. Direct calls and external delegation need
> separate capability claims because they have different dispatch and progress
> semantics. — 0198, 0204

The runner **MUST** surface resolved concurrency as a cap in dry-run JSON,
`evaluation.json`, and run receipts. Run-local events **MUST** additionally
record whether resolution was configured or automatic, the requested,
automatic, optional maximum, and resolved values, whether clamping occurred,
the dispatch mode, and the observed peak active direct calls or outstanding
harness requests. The runner **MUST NOT** expose a `--concurrency` flag or
public scheduler-strategy names.

Concurrent execution **MUST** preserve the
[observational equivalence](orchestration.md#scheduling-and-parallelism) the
orchestration contract requires.

Prompt-cache hits and provider-retained context **MUST NOT** be required for
correctness, resume, validation, or deterministic report output. Provider
context identifiers and prompt-cache status are recorded only in run-local
logs, never in `evaluation.json`, per the
[`evaluation.json` contract](evaluation-json.md#state).

Each requirement uses a fresh evaluator session. Provider prefix caching **MAY**
reuse stable request instructions and schema, but a provider conversation,
workspace transcript, or selected evidence **MUST NOT** be shared between
requirements.

For a direct-call evaluator, the runner **MUST** keep up to resolved concurrency
fresh calls active. It **MUST** accept and atomically persist each completed
result before releasing that slot, advance deterministic dependencies, and
immediately top up newly ready evaluator work without waiting for unrelated
active calls. Terminal failure or cancellation stops new dispatch and
interrupts active calls; accepted results remain durable and unaccepted work
remains resumable.

For a harness-backed run, resolved concurrency is the cap on outstanding
checkpointed work requests, per the [harness checkpoints
contract](#harness-checkpoints) below. Outstanding requests are eligible for
parallel service; they do not prove that the harness has active workers.

## Harness checkpoints

For a harness-backed run, the runner **MUST** keep every ownership listed
above across checkpoints: it builds each bounded work request, atomically
persists the awaiting-evaluator checkpoint with every outstanding call's
correlation metadata (request identity, work-unit identity, input hash,
correlation ID, and attempt) before any request leaves stdout, and later
validates each submitted result against its own checkpoint entry.

> Rationale: the CLI command returns control to the invoking agent between
> request and result, so the run must be resumable before the work requests
> are emitted. — 0194

While a harness run has dependency-ready judgment work that is not yet
outstanding, the runner **MUST** keep up to the run's resolved concurrency
work requests outstanding — a rolling window drawn from the dependency-ready
frontier — emitting each newly-ready request as capacity frees, without
waiting for other outstanding requests to be judged first. Each emitted
request **MUST** carry the same complete bounded work request a
single-request checkpoint carries. At resolved concurrency `1`, the runner
**MUST** keep exactly one request outstanding, preserving single-request
behavior.

> Rationale: the transport — not a configuration choice — was what pinned
> harness runs to one judgment at a time. Requests are drawn only from the
> dependency-ready frontier, so parallelism stays bounded by the work-graph
> DAG, and emitting as capacity frees avoids the barrier-latency loss of
> draining a whole wave before starting the next. The window is a transport
> optimization only: the accepted results of a windowed run must match the
> sequential path. — 0198

Pending-call metadata **MUST NOT** persist raw prompt, workspace content, or
result bodies; each pending request is rebuilt deterministically from the model
snapshot, work-graph state, effective source selector, and inspection policy,
and a rebuilt request
whose input hash no longer matches its checkpoint entry **MUST** fail with
`run_state_invalid` rather than accept judgment for changed input.

The runner **MUST** normalize, schema-validate, retry, accept, log, and
persist each harness payload through the same paths used for CLI- and
SDK-backed evaluator payloads, **MUST** accept each valid member of a
submission independently and persist it atomically with its work unit's
completion — so an interruption leaves every accepted member durable and
every still-outstanding member recoverable for resume — **MUST** reject a
mismatched, duplicate, or unsolicited result without discarding or altering
any other member's accepted result, and **MUST** enforce the
[harness identity binding](evaluator-contract.md#harness-evaluator) on
resume. Invalid output is never repaired outside the runner.

> Rationale: per-member binding by the runner's rebuilt-request hash means a
> stale or mismatched result cannot be accepted against evidence it was not
> judged against — the single-request integrity check, preserved per member.
> — 0194, 0198

Evaluator-call logging for harness dispatch **MUST** follow the same boundary
as every other transport: hashes, identities, durations, attempt state, usage
when available, and failure categories — never raw harness requests, source
contents, result bodies, credentials, or tokens.

## Requirement inspection and evidence

The runner **MUST** derive each area's effective source selector using the
format's inheritance rules and record its path, glob, or prose form for display
and validation. It **MUST NOT** walk, package, cap, truncate, or select source
files for judgment. The source identifies the evaluated subject; the
requirement evaluator determines which material is relevant.

Each `assessRateRequirement` request **MUST** carry the requirement and area
identity, effective source selector, requirement and area frames, applied rating
criteria, applicable body guidance, expected result schema, authorized
workspace root, and inspection policy. It **MUST** start a fresh evaluator
session with read-only workspace access, disabled network, no approval
escalation, no workspace writes, and repository instructions and discovered
content treated as untrusted data.

The evaluator **MAY** inspect supporting workspace context outside a concrete
source selector when the requirement needs interpretation or comparison. It
**MUST** classify each proposed observation as `evaluated` or `supporting` and
**MUST NOT** silently widen the area or requirement. Executable verification is
available only when an adapter declares and enforces a mediated, sandboxed path.
When unavailable, the evaluator records the limit and adjusts its assessment or
rating rather than receiving host shell access.

The requirement response **MUST** contain its assessment, rating, and evidence
proposal together. Before accepting it, the runner **MUST**:

- validate the proposal's exact shape, unique IDs, roles, and limits;
- resolve every file path against the authorized workspace and reject absolute
  paths, missing or non-regular files, non-UTF-8 content, lexical escapes, and
  symlink escapes;
- validate that `evaluated` file evidence falls within a path or glob source
  selector while permitting separately classified supporting evidence;
- validate line-range and Markdown-heading locators against the current file;
- compute file byte counts and SHA-256 digests itself;
- verify every assessment `sourceRef` names an accepted observation; and
- compute the canonical manifest hash.

The runner **MUST** persist the sealed evidence manifest atomically with the
paired assessment and rating. It **MUST NOT** persist file bodies, prompts,
hidden reasoning, tool transcripts, credentials, or command output. A malformed
or unverifiable proposal fails acceptance as `evidence_invalid`; honest missing
or inconclusive evidence remains a judgment limit or status, not an
infrastructure failure.

Accepted evidence is immutable for resume. A retry of an unaccepted requirement
opens a fresh session. A new run may discover different evidence; determinism
means the runner preserves and projects what it accepted, not that agentic
inspection returns identical judgments across runs.

## Logging

Run-local logs **MUST** be separate from `evaluation.json`.

> Rationale: execution telemetry changes frequently and is append-only. Keeping
> it outside the authoritative judgment artifact avoids noisy rewrites and
> keeps logs from becoming evaluation data. — 0192

The runner **MUST** write a run-local structured event log at
`logs/events.jsonl` recording lifecycle events: run creation, concurrency
resolution and dispatch mode, observed peak activity, work-unit
start/completion/failure, retry, resume decisions, report build, and run
completion.

The runner **MUST** write evaluator-call metadata to a run-local structured log
at `logs/evaluator-calls.jsonl`. Call metadata **SHOULD** include evaluator
kind, model when known, work-unit kind, attempt number, duration, input hash,
output hash, resolved concurrency, context/cache metadata when available,
declared capabilities, and usage when available. When the evaluator reports
cache-read or cache-creation input tokens, the runner **MUST** preserve each
count separately in the call's usage metadata, including a reported zero. An
unreported count **MUST** remain absent. Cache usage **MUST NOT** enter
`evaluation.json`.

> Rationale: the call log is where per-call usage lives; the cached-vs-fresh
> input split is what makes the prompt-caching saving measurable. Cache creation
> distinguishes a newly written prefix from an unavailable metric. — 0193, 0203

Run-local logs **MUST NOT** record raw prompts, workspace file bodies, raw model
responses, API keys, auth tokens, or environment variable values by default.

> Rationale: evaluation often touches proprietary source and may include
> prompt-injection text from evaluated files. Default observability must be
> useful without turning logs into a data leak. — 0192

The runner **MUST** log request input hashes and accepted evidence-manifest
hashes without storing raw prompt or workspace bodies by default.

## Failure taxonomy

The runner **MUST** use these stable machine-readable failure categories for
run, work-unit, evaluator, source, validation, and report-build failures:

| Category                    | Meaning                                                                                   |
| --------------------------- | ----------------------------------------------------------------------------------------- |
| `missing_evaluator`         | No usable evaluator could be selected.                                                    |
| `evaluator_unauthenticated` | The selected evaluator is present but not authenticated.                                  |
| `evaluator_incompatible`    | An installed CLI cannot honor non-interactive structured invocation.                      |
| `rate_limited`              | The evaluator call was rate limited (retryable).                                          |
| `timeout`                   | The evaluator call timed out (retryable).                                                 |
| `invalid_evaluator_output`  | The evaluator returned unparseable output (retryable).                                    |
| `schema_invalid_output`     | The evaluator output parsed but failed schema validation (retryable).                     |
| `evidence_invalid`          | Proposed evidence failed containment, locator, role, or reference validation (retryable). |
| `unsafe_source_content`     | Source packaging or evaluation refused unsafe source content.                             |
| `workspace_access_denied`   | The evaluator could not honor the authorized workspace boundary.                          |
| `run_state_invalid`         | The run artifact is missing, unsupported, or incompatible for resume.                     |
| `cancelled`                 | The run was interrupted by user cancellation or a termination signal.                     |
| `report_build_failed`       | Report generation failed.                                                                 |
| `internal_error`            | A bug or I/O failure prevented the requested action.                                      |

The retryable categories and attempt budget are defined by the
[orchestration retry policy](orchestration.md#retry-and-failure).

The runner **MUST** surface failure categories in `evaluation.json`,
`logs/events.jsonl`, and `--json` command receipts when a failure affects the
run result or command result.
