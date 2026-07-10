---
type: Design Doc
title: Deterministic evaluation runner design
description: Go runner and evaluator-adapter design for replacing skill-orchestrated evaluation.
tags: [evaluation, cli, design]
timestamp: 2026-07-09T00:00:00Z
---

# Deterministic evaluation runner design

## Context

This design answers the
[Deterministic evaluation runner functional spec](spec.md). The goal is to make
`qualitymd evaluation run` the evaluation engine while preserving `/quality
evaluate` as the agent-mediated user interface. The runner owns deterministic
workflow state and artifact writes; evaluators perform bounded judgment work.

## Approach

Implement the runner in Go inside `qualitymd`, close to the existing model,
workspace, evaluation, and report packages.

The runner has six main components:

- **Command adapter:** `qualitymd evaluation run` parses flags, resolves the
  workspace/model/config, chooses the evaluator, and renders human or JSON
  command output.
- **Runner core:** builds the evaluation work graph from the model and requested
  scope, schedules ready work units, applies retry/failure policy, and updates
  run state.
- **Execution planner:** resolves the configured execution strategy, concurrency
  cap, context-reuse policy, and evaluator capability use for the run.
- **Evaluator interface:** dispatches one typed work-unit request at a time to a
  CLI, API, shell, manual, subagent-backed, or context-reusing evaluator
  implementation and returns a typed result envelope.
- **Run store:** owns atomic reads/writes of `evaluation.json` and append-only
  run logs.
- **Report bridge:** renders Markdown reports from `evaluation.json` through the
  existing report renderer, refactored as needed to consume the new aggregate
  run artifact.

The high-level flow:

```text
qualitymd evaluation run
  -> resolve workspace, model, scope, evaluator
  -> create or load evaluation.json
  -> build deterministic work graph
  -> resolve execution strategy and evaluator capabilities
  -> schedule ready work units
  -> call evaluator for judgment work
  -> validate result envelope
  -> merge result into evaluation.json atomically
  -> append execution/evaluator logs
  -> build Markdown reports
  -> emit receipt
```

## Run state

Use one authoritative run-root file:

```text
.quality/evaluations/0007-full-eval/evaluation.json
```

The file is shaped around stable sections:

```json
{
  "schemaVersion": 4,
  "kind": "EvaluationRun",
  "manifest": {},
  "state": {},
  "results": {},
  "outputs": {}
}
```

- `manifest` carries immutable identity and setup: evaluation ID, created time,
  model path, requested scope, planned scope, selected evaluator profile, and
  resolved execution strategy.
- `state` carries execution lifecycle: run status, work-unit statuses,
  dependency hashes, failures, attempts, timestamps, strategy fallback
  decisions, and resume metadata. Provider context identifiers and prompt-cache
  status live only in the evaluator-call log, never here.
- `results` carries evaluation judgment: frames, assessments, ratings, factor
  and area analyses, advice, rankings, and recommendations.
- `outputs` carries generated report refs after report build.

The run store writes `evaluation.json` through a temp-file plus rename
sequence. It merges each accepted work-unit result and persists the file before
the work unit is marked complete, so a crash or cancellation never loses
accepted judgment work. A single store-owned writer serializes merges under
parallel execution, so renames never interleave. Completed runs keep their full
execution state; compacting transient state is deferred until a real size
problem appears.

Resume verifies compatibility first: the artifact's `schemaVersion` must be
supported, the manifest must resolve to the same workspace and model path, and
an explicit `--evaluator` flag must match the manifest's recorded evaluator
profile. Incompatibility fails with `run_state_invalid` — or a refusal naming
the evaluator conflict — rather than guessing. The runner then rebuilds the
graph from the current model and compares it with saved state. Reusable
completed work units stay complete; missing, malformed, incompatible, or
dependency-stale work units rerun.

## Work graph

Keep the current evaluation protocol order, but represent it as an explicit
internal DAG. Work-unit IDs are deterministic strings derived from routine kind
and model reference, for example:

```text
frameEvaluation
frameAreaEvaluation:area:root
assessRequirement:requirement:docs::links-resolve
rateRequirement:requirement:docs::links-resolve
analyzeFactor:factor:docs::usability
analyzeArea:area:docs
rankFindings
recommend
rankRecommendations
buildReports
```

Each work unit declares:

- kind;
- subject reference;
- dependencies;
- input selectors;
- expected result schema;
- whether it is deterministic-only or evaluator-backed;
- whether it is safe to run concurrently with sibling work;
- context affinity, such as run-wide, area-wide, or isolated;
- retry policy;
- merge target in `evaluation.json`.

The first implementation slice ships `auto` and `sequential` only, with `auto`
resolving to sequential execution. The graph still exposes readiness, context
affinity, and concurrency eligibility so bounded parallel execution and
subagent-backed work can land later without contract change, preserving
deterministic model order for outputs and report content.

## Execution strategy

Add an execution planner that turns evaluator capabilities and workspace config
into one resolved strategy for the run:

```text
auto
  -> sequential
  -> bounded parallel
  -> subagent-backed bounded parallel
  -> API context reuse plus bounded parallel where safe
```

The planner never changes the graph or output contract. It only decides which
ready work units can run at the same time, which evaluator capability handles
them, and whether a work unit should receive shared context metadata.

The strategy vocabulary stays small (the first slice implements `auto` and
`sequential`; `parallel` is specified now and lands later):

- `auto`: choose the best safe strategy the evaluator declares.
- `sequential`: run one ready work unit at a time.
- `parallel`: run independent ready work units up to a configured cap.

CLI-backed evaluators such as Codex and Claude can implement `parallel` by
spawning native subagents or worker sessions when the CLI exposes them. API
evaluators can implement it with ordinary concurrent calls, with context reuse
where supported. If the evaluator cannot prove the strategy is available, the
planner falls back to the next simpler strategy and records that decision.

Subagent-backed work remains transport detail. A subagent receives a bounded
work request, returns the same typed result envelope as any other evaluator
call, and never writes run files or expands scope.

For API evaluators, the planner should prepare model requests in cache-friendly
layers:

```text
stable runner/evaluator instructions
stable schema and rating vocabulary
stable model snapshot and scope context
area or source package context
work-unit-specific request
```

The first layers should remain byte-stable across sibling work units when their
inputs are the same. Prompt-cache hits are logged when reported by the provider,
but cache hits are never required for correctness.

Where an API supports reusable conversation, thread, session, or previous
response state, the planner may create a base context for the run or area and
fork or chain work units from it. Provider context IDs are recorded in the
evaluator-call log only, never in `evaluation.json`. Resume must be able to
rebuild equivalent requests from the model snapshot, source package, prompt
plan, and `evaluation.json` if the provider-side context no longer exists.

## Cancellation

SIGINT or SIGTERM cancels the run's root context. In-flight evaluator calls are
cancelled; their work units keep their attempt counts and stay incomplete. The
run store flushes final state, the runner marks the run `cancelled` in run
state and the event log, and the receipt reports the cancelled run with
`--resume` as the next action. Parallel execution needs no special shutdown
path: workers observe the same context, and the store's serialized writer
completes or discards the in-progress merge atomically.

## Evaluator interface

Use one Go interface for evaluator-backed work:

```go
type Evaluator interface {
    Capabilities() Capabilities
    Evaluate(ctx context.Context, req WorkRequest) (WorkResult, error)
}
```

`Capabilities` is the declaration the execution planner reads before
dispatching work: supported execution strategies, subagent support,
reusable-context kinds, and whether usage metadata is reported. The planner
resolves `auto` from this declaration and never schedules a strategy the
evaluator does not declare; when a declared capability fails at run time, it
degrades to the next simpler strategy and logs the fallback. CLI adapters fold
a probe of the installed CLI's version and flag support into their declaration,
so an incompatible CLI fails at selection rather than mid-run.

`WorkRequest` includes:

- run identity;
- work-unit kind and ID;
- subject model reference;
- requirement/factor/area context needed for the judgment;
- rating criteria;
- bounded source bundle;
- expected JSON schema or schema kind;
- safety instructions;
- execution strategy hints;
- context affinity and reusable context metadata when available;
- stable prompt-prefix and source-package hashes;
- correlation ID for logs.

`WorkResult` includes:

- work-unit ID;
- result payload;
- evaluator metadata;
- execution strategy metadata;
- context/cache metadata when available;
- usage metadata when available;
- failure category when the evaluator completed but could not produce ordinary
  judgment output.

Evaluator implementations:

- `codex`: subprocess adapter for Codex CLI.
- `claude`: subprocess adapter for Claude Code CLI.
- `openai`: direct OpenAI API adapter using an API key from the configured env
  var.
- `anthropic`: direct Anthropic API adapter using an API key from the configured
  env var.
- `shell`: deterministic local command/check adapter when a work unit declares a
  shell-compatible routine.
- `manual`: future adapter that emits task payloads and waits for submitted
  results.

CLI adapters should be implemented as subprocess transports around the same
`WorkRequest` JSON and output schema used by API adapters. Each adapter invokes
its CLI in its non-interactive mode with machine-readable output — for Claude
Code, print mode with JSON output; for Codex, its exec mode with JSON output —
sends the work request and expected schema in the prompt payload, and validates
the returned envelope before handing it to the runner, retrying per the
runner's policy. A CLI that cannot honor non-interactive structured invocation
fails selection with `evaluator_incompatible` instead of degrading into
unparseable runs. Adapters do not get special authority to read or write run
files. When a CLI adapter uses native subagents, the adapter still returns one
accepted `WorkResult` per work unit to the runner.

## Evaluator selection

Resolution order:

1. `qualitymd evaluation run --evaluator <name>`
2. `.quality/config.yaml` `evaluation.evaluator`
3. `auto`

`auto` uses deterministic local discovery:

1. an active or installed authenticated Codex CLI when detectable;
2. an active or installed authenticated Claude CLI when detectable;
3. configured API profiles whose required key env var is present;
4. a clear non-interactive failure listing available remedies.

The command does not prompt in non-interactive contexts. `/quality evaluate` may
ask the user to choose when the CLI reports an ambiguous or missing evaluator.

Execution strategy resolution is separate from evaluator selection. The default
strategy is `auto`; `auto` prefers the most capable safe strategy the selected
evaluator advertises, subject to the configured concurrency cap and the work
graph's context-affinity constraints.

## Source packaging

Source packaging is runner-owned. The runner resolves area `source` values,
applies deterministic file walking, excludes unsuitable files, and builds a
bounded source bundle per evaluator work unit.

The initial source packager should:

- resolve all paths model-relative and repository-contained;
- skip binary files and obvious generated/vendor dependency directories unless
  directly selected by the model source;
- cap bundle size with a deterministic truncation or blocked-work result;
- attach stable source refs and content hashes;
- preserve "source as data" instructions in every evaluator request.

Raw prompt bodies are not persisted by default. The evaluator-call log records
input hashes, output hashes, schema kind, evaluator kind, duration, attempts,
and usage metadata when available.

## Logging

Write run-local logs:

```text
logs/events.jsonl
logs/evaluator-calls.jsonl
```

`events.jsonl` records runner lifecycle events such as run creation, work-unit
start/completion/failure, strategy selection, concurrency fallback, retry,
resume decision, report build, and command completion.

`evaluator-calls.jsonl` records evaluator-call metadata only. It deliberately
omits raw prompt bodies, source bundles, raw responses, API keys, auth tokens,
and environment values. It can record prompt-prefix hashes, source-package
hashes, provider context IDs, prompt-cache status, token usage, and cost when
available.

The existing `.quality/logs/*-evaluate-feedback-log.md` remains a separate
workflow-experience artifact owned by the skill when `/quality evaluate` is the
entrypoint.

## Report rendering

Refactor the report builder to consume `evaluation.json` as its source of truth
for new runs. The first implementation can adapt `evaluation.json` into the
existing internal report model rather than rewriting report rendering.

Generated reports keep their current reader-oriented tree where useful:

```text
report.md
findings.md
recommendations.md
areas/...
```

Their source-data sections point to `evaluation.json` for new runs. Historical
multi-file runs may continue through the old report path while supported.

## Spec response

The command adapter satisfies the CLI-surface and config requirements by
centralizing flag parsing, default evaluator resolution, non-interactive errors,
and `--dry-run --json` receipts.

The runner core satisfies orchestration requirements by making the protocol an
explicit DAG with deterministic work-unit IDs, dependency checks, retry policy,
and resume decisions.

The execution planner satisfies parallelism, subagent, prompt-cache, and context
reuse requirements by treating them as strategy choices below the runner-owned
work graph.

The evaluator interface satisfies portability requirements by making Codex,
Claude, OpenAI, Anthropic, shell, and manual paths transports over the same typed
work-unit request/result contract, and satisfies the capability-declaration
requirement through `Capabilities()`, which the planner reads before scheduling.

The run store satisfies artifact and observability requirements by separating
`evaluation.json` from append-only logs and writing the authoritative file
atomically.

The report bridge satisfies report requirements by keeping Markdown output as a
deterministic projection rather than an input or second source of truth.

## Alternatives

### Keep the skill as orchestrator

Rejected. It preserves the current portability problem: every harness has to
reconstruct the same workflow from prompt instructions. It also prevents direct
API evaluators from using the same execution contract as subscription-backed CLI
evaluators.

### Use Mastra or another TypeScript/Python orchestration framework

Rejected for the core runner. Frameworks such as Mastra or LangGraph-like DAG
runtimes are useful for API-first model apps and meta-harness composition, but
this project already has a Go CLI, Go model parser, Go validation, and Go report
generation. Adding a second runtime would make distribution and compatibility
harder without improving the command contract.

Mastra-style SDK agents remain relevant as an adapter option: a future evaluator
could call into a hosted or local agent harness that provides tracing, workflow
composition, or subagent delegation. That adapter still must return the runner's
typed work-unit result envelope and must not become the authoritative run store.

### Let evaluator sessions own context and resume

Rejected. Provider conversations, CLI sessions, prompt caches, and subagent
threads are useful for latency and context management, but they are not a stable
QUALITY.md artifact. Resume and review need to work from `evaluation.json`, the
model snapshot, source packages, and logs even when provider-side state expires
or is unavailable.

### Write one JSON file per routine result

Rejected for new runner-created runs. The multi-file tree was useful when
agents authored discrete payloads and the CLI validated/persisted them one at a
time. A CLI-owned runner can keep routine boundaries internally while presenting
one authoritative run artifact.

### Write only a final `report.json`

Rejected. The structured artifact is not a report; it is the full evaluation
run. Naming it `evaluation.json` keeps the source of truth distinct from the
generated Markdown report tree.

### Make provider APIs the only model path

Rejected. Codex and Claude subscription users should be able to use their local
CLI authentication. API-key profiles are important, but not the only supported
path.

## Trade-offs and risks

- The runner becomes a larger CLI responsibility. That is intentional, but it
  raises the bar for tests and schema discipline.
- `evaluation.json` reduces file count but increases the importance of careful
  atomic writes and merge validation. Per-result persistence also makes write
  frequency proportional to work-unit count, so merges must stay cheap and the
  store's writer must stay serialized.
- CLI-backed evaluator automation may expose less usage/cost metadata than API
  evaluators. The usage schema must tolerate unavailable data.
- Prompt caching and context reuse improve cost and latency only when provider
  behavior cooperates. The runner must treat cache/session state as optional
  metadata, which limits how much performance can be guaranteed.
- Subagent-backed CLI execution improves coverage throughput, but it also
  increases token usage and makes failure surfaces more complex. The runner needs
  strict concurrency caps, structured results, and clear fallback to simpler
  strategies.
- Source packaging quality will strongly affect evaluation quality. The initial
  implementation should prefer explicit blocked/insufficient evidence results
  over oversized or lossy prompts that look complete.
- Historical multi-file runs will coexist with new single-artifact runs for some
  period. Status and report code need a clear dispatch boundary so compatibility
  does not become dual writing.

## Open questions

None. The spec's scope decisions settle the first slice: `manual` stays a
reserved name with no file-drop workflow, completed runs keep full execution
state with compaction deferred, the first slice ships `auto` and `sequential`
only, and provider context identifiers live in `logs/evaluator-calls.jsonl`
only.
