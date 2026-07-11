---
type: Functional Specification
title: Evaluation runner
description: CLI-owned deterministic evaluation engine, concurrency, run-local logging, and failure taxonomy.
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

## Concurrency

The runner **MUST** own concurrency resolution for ready evaluator-backed work
units.

> Rationale: parallelism, native subagents, and provider context reuse are
> useful only while they remain scheduling and transport choices under the same
> work graph. If they become alternate orchestration engines, the portability
> problem returns. `concurrency` is the user-facing control; scheduler strategy
> names are an implementation detail. — 0192, 0195

The workspace config **MAY** set `evaluation.concurrency` in
`.quality/config.yaml`. When it is present, it **MUST** be a positive integer. A
configured value of `1` selects sequential execution; a configured value greater
than `1` permits concurrent execution of dependency-ready evaluator-backed work
units up to that limit.

When `evaluation.concurrency` is absent, the runner **MUST** request automatic
concurrency equal to `max(2, runtime.NumCPU()*2)`.

The runner **MUST** resolve requested concurrency against the selected
evaluator's declared support for concurrent calls before execution begins. If
the selected evaluator does not support concurrent calls, then the resolved
concurrency **MUST** be `1`.

The runner **MUST** surface the resolved concurrency in dry-run JSON,
`evaluation.json`, run-local logs, and run receipts. The runner **MUST NOT**
expose public execution-strategy or strategy-fallback fields.

Concurrent execution **MUST** preserve the
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

Harness-backed runs **MUST** resolve to `concurrency: 1` until the runner and
harness contracts define multiple pending evaluator checkpoints.

## Harness checkpoints

For a harness-backed run, the runner **MUST** keep every ownership listed
above across checkpoints: it builds the bounded work request, atomically
persists an awaiting-evaluator checkpoint with the pending call's correlation
metadata (request identity, work-unit identity, input hash, correlation ID,
and attempt) before the request leaves stdout, and later validates the
submitted result against that checkpoint.

> Rationale: the CLI command returns control to the invoking agent between
> request and result, so the run must be resumable before the work request is
> emitted. — 0194

Pending-call metadata **MUST NOT** persist raw prompt, source, or result
bodies; the pending request is rebuilt deterministically from the model
snapshot, work-graph state, and current source package, and a rebuilt request
whose input hash no longer matches the checkpoint **MUST** fail with
`run_state_invalid` rather than accept judgment for changed input.

The runner **MUST** normalize, schema-validate, retry, accept, log, and
persist a harness payload through the same paths used for CLI- and API-backed
evaluator payloads, **MUST** reject a mismatched, duplicate, or unsolicited
result without advancing the work graph, and **MUST** enforce the
[harness identity binding](evaluator-contract.md#harness-evaluator) on
resume. Invalid output is never repaired outside the runner.

Evaluator-call logging for harness dispatch **MUST** follow the same boundary
as every other transport: hashes, identities, durations, attempt state, usage
when available, and failure categories — never raw harness requests, source
contents, result bodies, credentials, or tokens.

## Source packaging

Source reaches judgment in three steps: the effective selector is **detected**
to a kind, **resolved** by the resolver serving that kind, and **packaged**
as the bounded, hashed evidence bundle judgment consumes. Path and glob
resolution is the deterministic walk below, unchanged; other kinds resolve
through evaluator dispatch under the same bundle contract.

The runner **MUST** resolve each area's effective source selector as the
format defines it ([`SPECIFICATION.md`](../../SPECIFICATION.md) — Source
resolution): the area's own declared `source`, else the nearest declaring
ancestor's, else the document default — the directory containing the
QUALITY.md file. A source-less root area **MUST** package the model's
directory and its descendants, not an empty bundle, and a source-less child
area **MUST** inherit its nearest ancestor's effective selector.

> Rationale: the runner and `status` once derived source provenance
> independently, and the runner used only an area's own declared value —
> source-less areas were silently judged against empty evidence. Both surfaces
> now consume one shared model-layer resolver so they cannot drift apart
> again. — 0196

### Selector kind detection

The runner **MUST** detect each effective selector's kind from the bare
selector string, in this order: a selector containing glob metacharacters is
a **glob**; otherwise a selector that names an existing filesystem entry is a
**path**; otherwise the selector is **prose**. A selector that is absolute or
escapes the workspace **MUST** remain a filesystem selector under the
workspace-containment rules below and **MUST NOT** fall back to prose.

> Rationale: detection order makes the filesystem interpretation always win,
> so prose is only ever the meaning of a selector that cannot be filesystem
> material — and `source` stays a hand-authorable scalar with an unchanged
> frontmatter shape. The accepted hazard is a mistyped path detecting as
> prose; the pinned per-area record, dry-run surfacing, and the resolution
> instructions' path-like-selector guidance mitigate it. — 0197

A run's detected kinds **MUST** be recorded at run creation in the
[`evaluation.json` sources record](evaluation-json.md#sources) and honored on
resume, so a filesystem change mid-run cannot silently re-dispatch a selector
to a different resolver or change the graph shape. Switching kinds requires a
new run. Dry-run **MUST** surface each in-scope area's detected kind and its
serving resolver before anything runs.

> Rationale: a file created mid-run (often by the evaluating agent itself)
> would flip a prose area to path, change the graph shape under a pending
> checkpoint, and strand the captured bundle. — 0197

### Resolution

Material gathered by any resolver **MUST** be captured into the same bounded,
hashed, persisted source bundle before judgment, and a judgment work unit
**MUST NOT** receive evidence that is not part of a captured bundle.

> Rationale: the runner's guarantees — resumability, input-hash guards,
> evidence-bound re-judgment, source-as-data — are properties of the resolved
> bundle, not of how it was gathered. An agent-gathered "all specs" is
> legitimate only through the bundle contract, never as free repo exploration
> by the judge. — 0197

When an effective selector's kind has no deterministic resolver and the run's
evaluator dispatch can serve source resolution
([evaluator contract](evaluator-contract.md#capability-declaration)), the
runner **MUST** dispatch a `resolveSource` work unit for the area through the
same checkpoint transport as harness judgment, and **MUST** validate and
capture the returned material into the bundle — unique non-empty file paths,
the packaging caps below with truncation marks, SHA-256 per file, and the
shared bundle hash — before any dependent judgment work is dispatched. The
resolution request carries the selector, its kind, and the area frame, and an
empty source bundle: the resolver is fed a description, never pre-gathered
evidence.

> Rationale: a prose or live-system selector needs tools to gather; the
> harness checkpoint transport already carries correlated, validated,
> resumable work both ways. Resolution and judgment stay distinct requests so
> the gatherer and the judge are never the same uncontrolled step. — 0197

Captured prose bundles are frozen: workspace writes **MUST NOT** invalidate a
captured prose bundle, and re-gathering requires a new run. Gathering itself
is not reproducible — two runs over the same model may capture different
evidence for the same prose selector; reproducibility of record comes from the
captured, hashed, persisted bundle.

### Packaging

The runner **MUST** select and package source for evaluator work units through
a deterministic process:

- effective selectors resolve workspace-relative, and resolved paths **MUST
  NOT** escape the repository;
- a selector containing glob metacharacters (`*`, `?`, `[`, including a `**`
  segment spanning directories) **MUST** be expanded as a pattern over the
  workspace rather than looked up as a literal path, its matches feeding the
  same sorted, hashed, capped bundling as a directory walk;
- directory walks are sorted;
- `.git`, `.quality`, `node_modules`, `vendor`, and `dist` directories are
  skipped, including during glob expansion — except that a skipped directory
  named as a literal leading segment of a glob (`vendor/**`) is an explicit
  selection and **MUST** be walked;
- non-regular filesystem entries — symlinked directories, symlinked files,
  sockets, devices — **MUST** be skipped, and the walk **MUST NOT** error on
  them;
- binary files are skipped;
- bundles are capped at 64 KB per file and 512 KB per bundle, with explicit
  truncation marks when a cap applies; and
- every packaged file carries a SHA-256 content hash, and the bundle carries a
  bundle hash.

> Rationale: the format resolves "paths and globs" relative to the containing
> QUALITY.md file, but the runner stat-ed selectors literally, so a glob was
> always "missing". Committed symlinked directories (an agent-skill directory
> symlinked into `.claude/`) crashed the walk with `is a directory` and made
> runs non-resumable; skipping non-regular entries also keeps packaging inside
> the workspace boundary. — 0196

When an area's effective path or glob selector packages zero readable files —
it matches nothing, or matches only binary or unreadable material — the
runner **MUST** fail the affected work with the `source_unavailable`
category, naming the unresolved selector, and **MUST NOT** dispatch judgment
against an empty bundle. The document default always contains at least the
QUALITY.md file itself, so only a selector that resolves to nothing trips
this failure. A harness resolver reporting that the material a prose selector
describes does not exist classifies the same way: the material is missing.

> Rationale: a selector that resolved to nothing was silently packaged as an
> empty bundle and judged against zero evidence — the failure mode behind
> cascades of empty-evidence judgments in real runs. No evidence must be loud,
> not silent. — 0196

If no resolver is available for an effective selector — the selected
evaluator dispatch cannot serve resolution requests for its kind — then the
runner **MUST** fail the run with the `selector_unsupported` category before
any judgment is dispatched, naming the selector, its detected kind, and the
remedy (evaluate through harness dispatch, or change the selector), and
**MUST NOT** fall back to the document default or an empty bundle. The
classification boundary: `selector_unsupported` means _this run cannot
resolve this kind of selector_; `source_unavailable` keeps meaning _the
material is missing_.

> Rationale: a prose selector once failed as `source_unavailable`, which
> misdiagnosed an unsupported kind as missing material and told the author to
> "add the material it names". The remedy differs (change the selector or the
> evaluator, not the material), so the classification must too. — 0197

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
output hash, resolved concurrency, context/cache metadata when available, and
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
| `source_unavailable`        | The material a source selector names is missing.                           |
| `selector_unsupported`      | No available resolver serves the selector's detected kind.                 |
| `run_state_invalid`         | The run artifact is missing, unsupported, or incompatible for resume.      |
| `cancelled`                 | The run was interrupted by user cancellation or a termination signal.      |
| `report_build_failed`       | Report generation failed.                                                  |
| `internal_error`            | A bug or I/O failure prevented the requested action.                       |

The retryable categories and attempt budget are defined by the
[orchestration retry policy](orchestration.md#retry-and-failure).

The runner **MUST** surface failure categories in `evaluation.json`,
`logs/events.jsonl`, and `--json` command receipts when a failure affects the
run result or command result.
