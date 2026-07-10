---
type: Functional Specification
title: Deterministic evaluation runner
description: Requirements for replacing skill-orchestrated evaluation with a CLI-owned deterministic runner and pluggable evaluators.
tags: [evaluation, cli, agents, orchestration]
timestamp: 2026-07-09T00:00:00Z
---

# Deterministic evaluation runner

This change specifies a CLI-owned evaluation runner for QUALITY.md. It changes
the evaluation operating model from "skill orchestrates judgment and writes
structured payloads through the CLI" to "CLI owns the deterministic work graph
and invokes pluggable evaluators for bounded judgment work units."

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

Evaluation needs the same results whether it is run through Codex, Claude, a
direct OpenAI or Anthropic API key, a deterministic local checker, or a future
runtime. The current skill-first protocol makes that hard because orchestration
quality depends on the invoking harness. A deterministic runner lets
`qualitymd` own repeatable workflow behavior while preserving agent- and
model-mediated judgment as bounded evaluator work.

## Scope

This spec covers:

- the `qualitymd evaluation run` command;
- evaluator selection and configuration;
- evaluator work-unit boundaries;
- runner-owned execution strategy selection, including sequential execution,
  parallel workers, subagent fanout, and evaluator context reuse;
- the authoritative run artifact;
- run logging and observability;
- `/quality evaluate` integration with the runner;
- compatibility expectations for older evaluation runs.

Deferred:

- exact prompt text for each evaluator work unit;
- concrete source-bundling heuristics beyond safety and determinism
  requirements;
- custom user prompt templates;
- evaluator marketplace or plugin discovery;
- provider-specific conversation, thread, or cache API details beyond the
  runner's portable strategy contract;
- raw prompt/response tracing beyond an explicit future opt-in;
- `shell` and `manual` evaluator implementations — the first implementation
  slice documents them as reserved names only;
- bounded parallel and subagent-backed execution in the first implementation
  slice — the strategy contract below is normative now, and the first slice may
  resolve `auto` to sequential execution only;
- compacting completed-run execution state in `evaluation.json` — completed
  runs keep their full execution state until a real size problem appears;
- positional natural-language scope input, until the explicit flag surface
  settles.

Non-goals:

- making provider APIs mandatory;
- making CLI subscription users use API keys;
- making Mastra, LangGraph, the OpenAI Agents SDK, or any other orchestration
  framework mandatory for the core runner;
- allowing prompt-cache hits, stored conversation state, or subagent availability
  to affect evaluation correctness or persisted output semantics;
- allowing evaluators to write run artifacts directly;
- preserving the existing multi-file routine payload tree for new runs.

## Requirements

### Command surface

`qualitymd evaluation run` **MUST** execute a complete evaluation run for a
resolved `QUALITY.md` model.

> Durable spec: add `specs/cli/evaluation-run.md`; modify `specs/cli.md` to list
> `evaluation run`; modify `specs/evaluation/evaluation.md` to name the runner
> as the default evaluation execution surface.

`qualitymd evaluation run` **MUST** own run creation, resume, work-graph
execution, evaluator invocation, validation, persistence, report generation, and
final receipt output.

> Rationale: If the skill still owns any of these phases, the workflow keeps two
> orchestrators and cannot be made deterministic across harnesses.
>
> Durable spec: modify `specs/evaluation/orchestration.md` and
> `specs/evaluation/protocol.md`.

`qualitymd evaluation run` **MUST** support `--evaluator <name>`, where `<name>`
selects a built-in evaluator or configured evaluator profile.

> Durable spec: add to `specs/cli/evaluation-run.md`.

`qualitymd evaluation run` **MUST** treat an omitted `--evaluator` the same as
`--evaluator auto`.

> Rationale: Default evaluation should work for subscription-based Codex or
> Claude users without requiring a config file on first use.
>
> Durable spec: add to `specs/cli/evaluation-run.md` and modify
> `specs/cli/evaluation-status.md` if status reports evaluator readiness.

`qualitymd evaluation run --dry-run --json` **MUST** emit a deterministic
machine-readable preview of the resolved model, scope, evaluator, execution
strategy, work-unit counts, expected run path behavior, and next actions without
invoking an evaluator or writing evaluation judgment data.

> Durable spec: add to `specs/cli/evaluation-run.md`.

`qualitymd evaluation run` **MUST** support `--resume <run>` for resuming an
existing run whose authoritative run artifact is present and compatible.

> Durable spec: add to `specs/cli/evaluation-run.md`; modify
> `specs/evaluation/orchestration.md`.

A run artifact is resume-compatible when its schema version is supported by the
running `qualitymd` version and its manifest resolves to the same workspace and
model path. If compatibility verification fails, then `--resume` **MUST** fail
with `run_state_invalid` and report whether starting a new run is the remedy.

> Durable spec: add to `specs/cli/evaluation-run.md` and the `evaluation.json`
> artifact contract.

If `--resume` is combined with `--evaluator` naming a different evaluator than
the run manifest records, then `qualitymd evaluation run` **MUST** refuse the
run and report the conflict.

> Rationale: A run's judgments stay attributable to one evaluator profile.
> Re-evaluating with a different evaluator is a new run, not a resume.
>
> Durable spec: add to `specs/cli/evaluation-run.md`.

When a run is interrupted by user cancellation or a termination signal, the
runner **MUST** leave `evaluation.json` valid and resumable, record the
interruption in run state and the event log, and report the run as `cancelled`
rather than failed.

> Rationale: Stopping a long evaluation is an expected user action, not an
> internal error; accepted work must survive for `--resume`.
>
> Durable spec: add to `specs/cli/evaluation-run.md` and modify
> `specs/evaluation/orchestration.md`.

`qualitymd evaluation run` **SHOULD** accept explicit scope flags before adding
positional natural-language scope input.

> Rationale: Explicit flags keep the new deterministic command easier to parse,
> document, and drive from agents. Positional convenience can be added after the
> runner contract is stable.
>
> Durable spec: add to `specs/cli/evaluation-run.md`.

### Evaluator configuration

The workspace config **MUST** support `evaluation.evaluator`, whose value is
`auto` or the name of an evaluator profile.

> Durable spec: modify the workspace configuration contract in
> `specs/evaluation/evaluation.md` or a new configuration spec if one is split
> out.

When `evaluation.evaluator` is absent, `qualitymd evaluation run` **MUST** behave
as though `evaluation.evaluator: auto` were configured.

> Durable spec: modify the workspace configuration contract.

The workspace config **MAY** support an `evaluators` map of named profiles.
Each profile **MUST** declare a `kind`.

> Durable spec: modify the workspace configuration contract.

The workspace config **MAY** support `evaluation.executionStrategy`, whose value
selects `auto`, sequential execution, bounded parallel execution, or a named
configured execution strategy.

> Rationale: scope remains the user-facing way to bound evaluation breadth, while
> execution strategy lets the runner use subagents, workers, or provider context
> optimizations without changing evaluation semantics.
>
> Durable spec: modify the workspace configuration contract and add execution
> strategy requirements under `specs/evaluation/`.

Built-in evaluator names **MUST** be reserved and **MUST NOT** be shadowed by
custom evaluator profiles.

Reserved names are:

- `auto`
- `codex`
- `claude`
- `openai`
- `anthropic`
- `shell`
- `manual`

> Durable spec: modify the workspace configuration contract.

API-key evaluator profiles **MUST** reference secrets by environment-variable
name, not by secret value.

> Rationale: Config files are likely to be committed. A secret field would make
> accidental credential disclosure the path of least resistance.
>
> Durable spec: modify the workspace configuration contract and CLI error
> contract.

### Evaluator contract

An evaluator **MUST** be treated as the runtime used for bounded evaluation
judgment work units.

> Durable spec: add a durable evaluator contract under `specs/evaluation/`.

Evaluators **MUST NOT** own run state, scope expansion, dependency ordering,
artifact paths, report generation, or final rating/report authority outside the
typed work-unit result they return.

> Rationale: The runner is the orchestrator. Evaluators are interchangeable only
> if they cannot become hidden workflow engines.
>
> Durable spec: add a durable evaluator contract under `specs/evaluation/`.

Every evaluator kind **MUST** consume the same work-unit envelope and return the
same schema-valid result envelope for the same work-unit kind.

> Durable spec: add a durable evaluator contract under `specs/evaluation/`.

Every evaluator **MUST** declare its execution capabilities — supported
execution strategies, subagent support, reusable-context support, and usage
reporting — through a capability declaration the runner reads before
dispatching work.

> Rationale: The execution planner can only choose strategies an evaluator
> provably supports; assuming undeclared capabilities reintroduces
> harness-dependent behavior.
>
> Durable spec: add to the evaluator contract under `specs/evaluation/`.

CLI-backed evaluators **MUST** be invoked non-interactively and **MUST** return
machine-readable structured output that the runner validates against the work
unit's expected schema.

> Durable spec: add to the evaluator contract under `specs/evaluation/`.

If an installed CLI cannot honor non-interactive structured invocation, then
the runner **MUST** classify the failure as `evaluator_incompatible` and report
remediation, including any other available evaluators.

> Rationale: The runner's determinism rests on structured evaluator output. A
> CLI version that cannot provide it must fail loudly at selection, not degrade
> into unparseable runs.
>
> Durable spec: add to the evaluator contract and failure taxonomy under
> `specs/evaluation/`.

The built-in evaluator set **SHOULD** include CLI-backed `codex` and `claude`
evaluators, direct API-backed `openai` and `anthropic` evaluators, and
non-provider `shell` or `manual` evaluators when those are implemented.

> Rationale: The runner must work for subscription users and API-key users. CLI
> adapters keep Codex and Claude subscriptions first-class instead of forcing an
> API-only product shape.
>
> Durable spec: add a durable evaluator contract under `specs/evaluation/` and
> modify relevant skill specs.

When an evaluator returns invalid JSON, schema-invalid JSON, or a result that
does not match the requested work unit, the runner **MUST** classify the failure
and retry only according to the runner's retry policy.

> Durable spec: add a durable evaluator contract and modify
> `specs/evaluation/orchestration.md`.

### Execution and context strategy

The runner **MUST** own execution strategy selection for ready work units.

> Rationale: parallelism, native subagents, and provider context reuse are useful
> only when they remain scheduling and transport choices under the same work graph.
> If they become alternate orchestration engines, the portability problem returns.
>
> Durable spec: modify `specs/evaluation/orchestration.md` and add execution
> strategy requirements under `specs/evaluation/`.

Execution strategies **MUST** preserve observational equivalence to deterministic
sequential execution in model order.

> Durable spec: modify `specs/evaluation/orchestration.md`.

The runner **MUST NOT** select an execution strategy the selected evaluator does
not declare. If a declared capability proves unavailable at run time, then the
runner **MUST** fall back to the next simpler strategy and record the fallback
in run state and logs.

> Durable spec: add execution strategy requirements under `specs/evaluation/`
> and modify `specs/evaluation/orchestration.md`.

The runner **MUST** cap concurrency through a deterministic policy and surface
the resolved concurrency in dry-run JSON, run state, and logs.

> Durable spec: add execution strategy and logging requirements under
> `specs/evaluation/` and `specs/cli/evaluation-run.md`.

Where a CLI-backed evaluator exposes native subagents or worker threads, the
runner **MAY** use them for independent evidence collection, verification,
completeness sweeps, or other ready work units.

> Durable spec: add execution strategy requirements under `specs/evaluation/`.

Subagent-backed work **MUST** return the same typed result envelope as any other
evaluator-backed work, and subagents **MUST NOT** write run artifacts, expand
scope, change dependency ordering, or produce final authority outside the
accepted result envelope.

> Durable spec: add to the evaluator contract and modify
> `specs/evaluation/orchestration.md`.

Where an API-backed evaluator supports prompt caching, the runner **SHOULD**
shape model requests with stable cache-friendly prefixes followed by
work-unit-specific deltas.

> Rationale: stable prefixes can reduce cost and latency, but cache availability
> is provider- and time-dependent. Correctness must come from the work graph and
> persisted artifacts, not from a cache hit.
>
> Durable spec: add evaluator prompt-shaping requirements under
> `specs/evaluation/`.

Prompt-cache hits **MUST NOT** be required for correctness, resume, validation,
or deterministic report output.

> Durable spec: add execution strategy requirements under `specs/evaluation/`.

Where an evaluator supports reusable conversation, thread, session, or previous
response state, the runner **MAY** create shared base context and fork or chain
work units from it.

> Durable spec: add evaluator context-state requirements under
> `specs/evaluation/`.

Reusable evaluator context **MUST** be treated as reconstructible execution
metadata, not as authoritative run state. A resumed run **MUST NOT** require a
provider-retained conversation, thread, session, or cache entry that cannot be
recreated from the model snapshot, source package, runner prompts, and
`evaluation.json` state.

> Durable spec: add to the `evaluation.json` artifact contract and execution
> strategy requirements under `specs/evaluation/`.

The runner **MUST** record provider context identifiers and prompt-cache status
only in run-local logs, not in `evaluation.json`.

> Rationale: Provider-retained identifiers expire outside the run's control.
> Keeping them out of the authoritative artifact keeps resume honest about what
> is reconstructible.
>
> Durable spec: add to the `evaluation.json` artifact contract and logging
> requirements under `specs/evaluation/`.

### Run artifacts

New runner-created evaluations **MUST** write one authoritative structured run
artifact named `evaluation.json` at the evaluation run root.

> Durable spec: replace the new-run data layout in
> `specs/evaluation/records/data-layout.md`; modify report source-data contracts
> in `specs/evaluation/reports/report-tree.md`.

`evaluation.json` **MUST** contain the run manifest, execution state, structured
evaluation results, advice results, and generated output refs needed to validate,
resume, review, and render the run.

> Durable spec: add a durable artifact contract for `evaluation.json`.

`evaluation.json` **MUST** be the source of truth for new runner-generated
Markdown reports.

> Durable spec: modify `specs/evaluation/reports/report-tree.md`.

Generated Markdown reports **MUST NOT** be used as input to regenerate,
validate, or resume a run.

> Durable spec: preserve and modify `specs/evaluation/reports/report-tree.md`.

The runner **MUST** write `evaluation.json` atomically when updating persisted
run state.

> Durable spec: add to the `evaluation.json` artifact contract.

The runner **MUST** merge each accepted work-unit result into `evaluation.json`
and persist it atomically before treating the work unit as complete for
scheduling or resume.

> Rationale: An interrupted run must resume without repeating accepted judgment
> work; batched persistence silently discards completed evaluator calls.
>
> Durable spec: add to the `evaluation.json` artifact contract and modify
> `specs/evaluation/orchestration.md`.

While parallel execution is active, the runner **MUST** serialize
`evaluation.json` writes so concurrent merges cannot interleave.

> Durable spec: add to the `evaluation.json` artifact contract.

Run-local logs **MUST** be separate from `evaluation.json`.

> Rationale: Execution telemetry changes frequently and may be append-only.
> Keeping it outside the authoritative judgment artifact avoids noisy rewrites
> and keeps logs from becoming evaluation data.
>
> Durable spec: add logging requirements under `specs/evaluation/`.

### Logging and observability

`qualitymd evaluation run` **MUST** write human progress diagnostics to stderr
and keep stdout reserved for the command payload.

> Durable spec: add to `specs/cli/evaluation-run.md`.

The runner **MUST** write a run-local structured event log at
`logs/events.jsonl`.

> Durable spec: add logging requirements under `specs/evaluation/`.

The runner **MUST** write evaluator-call metadata to a run-local structured log
without recording raw prompts, raw source bundles, raw model responses, API keys,
auth tokens, or environment variable values by default.

> Rationale: Evaluation often touches proprietary source and may include
> prompt-injection text from evaluated files. Default observability must be
> useful without turning logs into a data leak.
>
> Durable spec: add logging requirements under `specs/evaluation/`.

Evaluator-call metadata **SHOULD** include evaluator kind, model when known,
schema/work-unit kind, input hash, output hash, attempt number, duration, result
status, execution strategy, context/cache metadata when available, and token/cost
usage when available.

> Durable spec: add logging requirements under `specs/evaluation/`.

Usage and cost metadata **MUST** be optional and **MUST** distinguish unavailable
usage from zero usage.

> Durable spec: add logging requirements under `specs/evaluation/`.

### Source and safety

The runner **MUST** treat evaluated source content as data, not instructions.

> Durable spec: preserve the invariant in `specs/evaluation/evaluation.md` and
> apply it to runner source packaging and evaluator prompts.

The runner **MUST** select and package source for evaluator work units through a
deterministic process.

> Durable spec: add source-packaging requirements under `specs/evaluation/`.

The runner **MUST** log source and prompt input hashes for evaluator calls
without storing raw prompt bodies by default.

> Durable spec: add logging and source-packaging requirements under
> `specs/evaluation/`.

### Failure taxonomy

The runner **MUST** use stable machine-readable failure categories for run,
work-unit, evaluator, source, validation, and report-build failures.

Initial categories are:

- `missing_evaluator`
- `evaluator_unauthenticated`
- `evaluator_incompatible`
- `missing_api_key`
- `rate_limited`
- `timeout`
- `invalid_evaluator_output`
- `schema_invalid_output`
- `unsafe_source_content`
- `insufficient_evidence`
- `source_unavailable`
- `run_state_invalid`
- `cancelled`
- `report_build_failed`
- `internal_error`

> Durable spec: add failure taxonomy requirements under `specs/evaluation/` and
> `specs/cli/evaluation-run.md`.

The runner **MUST** surface failure categories in `evaluation.json`,
`logs/events.jsonl`, and `--json` command receipts when a failure affects the
run result or command result.

> Durable spec: add failure taxonomy requirements.

### Compatibility

New runner-created runs **MUST NOT** write the existing multi-file routine
payload tree as their authoritative structured data.

> Rationale: The multi-file tree exists primarily to let skill-authored routine
> payloads be validated and persisted incrementally. A CLI-owned runner can keep
> the same structured concepts without exposing each routine as a separate file.
>
> Durable spec: modify `specs/evaluation/records/data-layout.md` and
> `specs/evaluation/records/payload-kinds.md`.

Existing multi-file evaluation runs **MAY** remain readable as historical runs,
but this change **MUST NOT** require migrations, dual writers, or compatibility
payload copies for new runs.

> Durable spec: modify `specs/evaluation/evaluation.md`,
> `specs/evaluation/records/data-layout.md`, and relevant CLI status/report
> specs.

### `/quality evaluate` integration

The `/quality evaluate` workflow **MUST** invoke `qualitymd evaluation run`
rather than orchestrating evaluation directly.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md`.

The `/quality evaluate` workflow **MUST** continue to provide the
agent-mediated user interface: intent parsing, run frame, evaluator/default
selection explanation, CLI invocation, progress summary, result summary, and
next workflow routing.

> Durable spec: modify `specs/skills/quality-skill/evaluation.md`,
> `specs/skills/quality-skill/quality-skill.md`, and runtime skill files.

The `/quality evaluate` workflow **MUST NOT** independently collect evidence,
assign ratings, run a parallel QC loop, or second-guess the runner's authoritative
evaluation result by default.

> Rationale: A wrapper that re-evaluates the source recreates the two-engine
> architecture this change removes.
>
> Durable spec: modify skill specs and runtime skill files.

## Durable spec changes

### To add

- `specs/cli/evaluation-run.md` for `qualitymd evaluation run`.
- A durable evaluator contract under `specs/evaluation/`.
- A durable `evaluation.json` artifact contract under `specs/evaluation/`.
- Durable execution strategy, prompt-cache, and evaluator context-state
  requirements under `specs/evaluation/`.
- Durable logging and failure-taxonomy requirements under `specs/evaluation/`.

### To modify

- `specs/cli.md` to list `evaluation run`.
- Existing per-command CLI specs for status, report, and data commands where
  they interact with new runner-created runs.
- `specs/evaluation/evaluation.md`, `protocol.md`, `orchestration.md`,
  `records/data-layout.md`, `records/payload-kinds.md`, and
  `reports/report-tree.md`.
- `specs/skills/quality-skill/quality-skill.md`,
  `specs/skills/quality-skill/evaluation.md`, and
  `specs/skills/quality-skill/workflows/evaluate.md`.

### To rename

None.

### To delete

None. Historical specs and archived change records stay intact.

## Open questions

None. The first-slice scope questions (reserved `shell`/`manual` names,
sequential-only `auto`, state compaction, positional scope) are resolved as
deferred items under [Scope](#scope), and the context-metadata split is decided
by the execution and context strategy requirements.
