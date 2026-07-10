---
type: Functional Specification
title: Evaluation runner
description: CLI-owned deterministic evaluation engine, execution strategy, run-local logging, and failure taxonomy.
tags: [evaluation, runner, orchestration]
timestamp: 2026-07-09T00:00:00Z
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

Evaluation needs the same results whether it is run through Codex, Claude, a
direct OpenAI or Anthropic API key, or a future runtime. When a skill
orchestrated evaluation, workflow quality depended on the invoking harness:
every harness had to reconstruct the same protocol from prompt instructions.
The runner lets `qualitymd` own repeatable workflow behavior while preserving
agent- and model-mediated judgment as bounded evaluator work units. — 0192

## Scope

Deferred:

- bounded parallel and subagent-backed execution — the strategy contract below
  is normative now, and the first implementation slice resolves `auto` to
  sequential execution only;
- `shell` and `manual` evaluator implementations (reserved names);
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

## Execution strategy

The runner **MUST** own execution strategy selection for ready work units.

> Rationale: parallelism, native subagents, and provider context reuse are
> useful only while they remain scheduling and transport choices under the same
> work graph. If they become alternate orchestration engines, the portability
> problem returns. — 0192

The strategy vocabulary is `auto`, `sequential`, and `parallel`. The workspace
config **MAY** set `evaluation.executionStrategy` in `.quality/config.yaml`;
when it is absent, the runner **MUST** behave as though `auto` were configured.
If the configured value is unknown, then the runner **MUST** fall back to
`sequential` and record the fallback.

The first implementation slice resolves `auto` to `sequential` only. While
`parallel` execution is unimplemented, a `parallel` selection **MUST** fall
back to `sequential` with the fallback recorded in run state and logs.

The runner **MUST NOT** select an execution strategy the selected evaluator
does not declare. If a declared capability proves unavailable at run time, then
the runner **MUST** fall back to the next simpler strategy and record the
fallback in run state and logs.

The runner **MUST** cap concurrency through a deterministic policy (currently
`1`) and surface the resolved concurrency in dry-run JSON, run state, and
receipts.

Every strategy **MUST** preserve the
[observational equivalence](orchestration.md#scheduling-and-parallelism) the
orchestration contract requires.

Prompt-cache hits and provider-retained context **MUST NOT** be required for
correctness, resume, validation, or deterministic report output. Provider
context identifiers and prompt-cache status are recorded only in run-local
logs, never in `evaluation.json`, per the
[`evaluation.json` contract](evaluation-json.md#state).

Per-area evaluator session reuse — resuming one CLI session per area so later
work units in the area send only their prompt deltas — is a permitted
context-reuse strategy under the
[evaluator contract](evaluator-contract.md#prompt-shaping-and-reusable-context)
and the invariant above: a dropped or unavailable session only costs
re-transmitted tokens, never output changes.

## Source packaging

The runner **MUST** select and package source for evaluator work units through
a deterministic process:

- area `source` values resolve workspace-relative, and resolved paths **MUST
  NOT** escape the repository;
- directory walks are sorted;
- `.git`, `.quality`, `node_modules`, `vendor`, and `dist` directories are
  skipped;
- binary files are skipped;
- bundles are capped at 64 KB per file and 512 KB per bundle, with explicit
  truncation marks when a cap applies; and
- every packaged file carries a SHA-256 content hash, and the bundle carries a
  bundle hash.

Packaged source **MUST** be presented to evaluators as data accompanied by
standing safety instructions, applying the source-as-data invariant in
[Evaluation](evaluation.md#shared-invariants).

The runner **MUST** package an area's source once per run and reuse that
bundle for every work unit in the area, with determinism, truncation, and
hashing unchanged: the reused bundle is byte-identical to a re-packaged one
and carries the same bundle hash.

> Rationale: packaging is deterministic, so re-packaging per requirement only
> repeats filesystem work; one bundle per area also gives evaluators a stable
> object to cache or transmit once. — 0193

## Logging

Run-local logs **MUST** be separate from `evaluation.json`.

> Rationale: execution telemetry changes frequently and is append-only. Keeping
> it outside the authoritative judgment artifact avoids noisy rewrites and
> keeps logs from becoming evaluation data. — 0192

The runner **MUST** write a run-local structured event log at
`logs/events.jsonl` recording lifecycle events: run creation, work-unit
start/completion/failure, strategy selection and fallback, retry, resume
decisions, report build, and run completion.

The runner **MUST** write evaluator-call metadata to a run-local structured log
at `logs/evaluator-calls.jsonl`. Call metadata **SHOULD** include evaluator
kind, model when known, work-unit kind, attempt number, duration, input hash,
output hash, execution strategy, context/cache metadata when available, and
usage when available. When the evaluator reports cached input tokens, the
runner **MUST** include the cached-input-token count in the call's usage
metadata.

> Rationale: the call log is where per-call usage lives; the cached-vs-fresh
> input split is what makes the prompt-caching saving measurable. — 0193

Run-local logs **MUST NOT** record raw prompts, raw source bundles, raw model
responses, API keys, auth tokens, or environment variable values by default.

> Rationale: evaluation often touches proprietary source and may include
> prompt-injection text from evaluated files. Default observability must be
> useful without turning logs into a data leak. — 0192

The runner **MUST** log prompt and source input hashes for evaluator calls
without storing raw prompt bodies by default.

## Failure taxonomy

The runner **MUST** use these stable machine-readable failure categories for
run, work-unit, evaluator, source, validation, and report-build failures:

| Category                    | Meaning                                                                    |
| --------------------------- | -------------------------------------------------------------------------- |
| `missing_evaluator`         | No usable evaluator could be selected.                                     |
| `evaluator_unauthenticated` | The selected evaluator is present but not authenticated.                   |
| `evaluator_incompatible`    | An installed CLI cannot honor non-interactive structured invocation.       |
| `missing_api_key`           | A selected API profile's key environment variable is unset.                |
| `rate_limited`              | The evaluator call was rate limited (retryable).                           |
| `timeout`                   | The evaluator call timed out (retryable).                                  |
| `invalid_evaluator_output`  | The evaluator returned unparseable output (retryable).                     |
| `schema_invalid_output`     | The evaluator output parsed but failed schema validation (retryable).      |
| `unsafe_source_content`     | Source packaging or evaluation refused unsafe source content.              |
| `insufficient_evidence`     | A work unit could not produce defensible judgment from available evidence. |
| `source_unavailable`        | Modeled source could not be resolved.                                      |
| `run_state_invalid`         | The run artifact is missing, unsupported, or incompatible for resume.      |
| `cancelled`                 | The run was interrupted by user cancellation or a termination signal.      |
| `report_build_failed`       | Report generation failed.                                                  |
| `internal_error`            | A bug or I/O failure prevented the requested action.                       |

The retryable categories and attempt budget are defined by the
[orchestration retry policy](orchestration.md#retry-and-failure).

The runner **MUST** surface failure categories in `evaluation.json`,
`logs/events.jsonl`, and `--json` command receipts when a failure affects the
run result or command result.
