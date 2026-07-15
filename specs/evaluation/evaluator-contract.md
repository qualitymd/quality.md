---
type: Functional Specification
title: Evaluator contract
description: Capability, work-unit envelope, result envelope, and configuration contract for evaluation runner evaluators.
tags: [evaluation, evaluator, agents]
timestamp: 2026-07-11T00:00:00Z
---

# Evaluator contract

An evaluator is the runtime the [evaluation runner](runner.md) uses for bounded
evaluation judgment work units. This spec defines the contract every evaluator
kind honors and the workspace configuration that names evaluators.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

Evaluation judgment can come from a Codex or Claude subscription CLI, a direct
OpenAI or Anthropic API key, or a future runtime. Keeping every path
first-class — without forcing an API-only product shape on subscription users —
requires one typed request/result contract that every kind is a transport over.
— 0192

## Boundaries

The runner **MUST** treat an evaluator as the runtime for bounded evaluation
judgment work units, and nothing more.

Evaluators **MUST NOT** own run state, scope expansion, dependency ordering,
artifact paths, report generation, or final rating/report authority outside the
typed work-unit result they return.

> Rationale: the runner is the orchestrator. Evaluators are interchangeable
> only if they cannot become hidden workflow engines. — 0192

## Capability declaration

Every evaluator **MUST** declare structured output, source resolution, tool use,
concurrent calls, nested subagents, fresh-session isolation, cancellation,
usage reporting, turn limits, token or cost limits, context-window visibility,
compaction control, sandbox control, and executable override. Unsupported
controls **MUST** be represented as unsupported rather than silently ignored,
and the runner **MUST** read the declaration before dispatching work.

> Rationale: the runner can only choose strategies an evaluator provably
> supports; assuming undeclared capabilities reintroduces harness-dependent
> behavior. — 0192

Source-resolution support — serving `resolveSource` work requests that gather
the material a non-deterministic source selector describes — is a dedicated
capability, distinct from subagent support. An evaluator **MUST NOT** be
dispatched resolution work it does not declare; the runner fails such runs
with `selector_unsupported` per the
[runner source contract](runner.md#source-packaging).

> Rationale: subagent support describes an evaluator's internal parallelism
> for judgment work; serving a distinct request kind is a different promise,
> and conflating them would make every subagent-capable evaluator implicitly
> claim resolution. — 0197

When a requested policy depends on an unsupported capability, planning **MUST**
choose a documented safe fallback that preserves the policy or fail before
evaluator work with an actionable classified error.

## Work-unit envelope

Every evaluator kind **MUST** consume the same work-unit envelope for the same
work-unit kind. The envelope carries:

- the run ID;
- the work-unit ID and kind;
- the subject model reference;
- the work-unit instructions;
- context payloads — the upstream accepted results the judgment needs — split
  into shared context (stable across the subject area's work units, such as
  the area evaluation frame) and per-work-unit context;
- the bounded, hashed source bundle;
- the expected JSON schema for the result payload;
- stable prompt-prefix and source-package hashes; and
- a correlation ID for logs.

The envelope's rendered prompt **MUST** carry an explicit split between the
stable prefix and the per-work-unit delta, so an evaluator can mark or reuse
the prefix without re-deriving it. Prompt ordering is defined under
[prompt shaping](#prompt-shaping-and-reusable-context).

## Result envelope

Every evaluator kind **MUST** return the same schema-valid result envelope for
the same work-unit kind. The envelope carries:

- the result payload;
- the evaluator kind and, when known, the model used;
- context metadata, when available;
- usage metadata, when available; and
- a failure category and detail when the evaluator completed but could not
  produce ordinary judgment output.

Usage metadata **MUST** be optional and **MUST** distinguish unavailable usage
from zero usage.

When a provider reports cached input tokens, the evaluator **MUST** record
them in the result's usage metadata, distinct from total input tokens.

> Rationale: without a cached-vs-fresh input signal, the prompt-caching saving
> can be neither verified nor regression-tested. — 0193

When an evaluator returns invalid JSON, schema-invalid JSON, or a result that
does not match the requested work unit, the runner **MUST** classify the
failure and retry only per the
[orchestration retry policy](orchestration.md#retry-and-failure).

## Built-in evaluators

The built-in evaluator kinds are `harness` (checkpointed dispatch to the
invoking agent harness), `codex` (Codex SDK and authenticated local runtime),
`claude` (Claude Agent SDK and authenticated local runtime), `openai` (direct
OpenAI API), and `anthropic` (direct Anthropic API). `shell` and `manual` are
reserved names with no implementations.

SDK-backed agent evaluators **MUST** be invoked non-interactively with
machine-readable structured output and a fresh isolated judgment session. A
provider-managed child runtime remains inside the evaluator boundary; it is not
a project-owned sidecar and **MUST NOT** be required when another evaluator or a
non-evaluation command is used. The runner validates every result against the
work unit's expected schema.

Where a supported SDK advertises JSON Schema output enforcement or ephemeral /
no-session-persistence controls for its non-interactive mode, its built-in
evaluator adapter **MUST** use those controls for bounded work requests,
detecting the capability from the installed CLI before the first judgment
call. The runner still **MUST** validate the returned payload independently.

> Rationale: prompt-only JSON and retained one-off sessions add avoidable
> output failures and local state without replacing runner validation. — 0194

If an installed provider runtime cannot honor non-interactive structured invocation, then
evaluator selection **MUST** fail with `evaluator_incompatible` and report
remediation, including any other available evaluators.

> Rationale: the runner's determinism rests on structured evaluator output. A
> CLI version that cannot provide it must fail loudly at selection, not degrade
> into unparseable runs. — 0192

API-backed evaluators **MUST** read their API key from the profile's configured
environment variable. If that variable is unset when an API-backed evaluator is
selected, then selection **MUST** fail with `missing_api_key`.

Evaluator credentials are not interchangeable across kinds: local
subscription/access-token authentication belongs to SDK-backed agent evaluators, direct
provider API keys belong to the API evaluators, and the harness evaluator uses
the invoking agent's own authentication and **MUST NOT** require a provider
API key. Guidance and configuration **MUST** reference secrets by
environment-variable name only.

## Harness evaluator

`harness` is the reserved built-in evaluator for judgment supplied by the
invoking agent harness. It **MUST** be selected explicitly — by the quality
skill or the caller — never by `auto` discovery, and it **MUST** delegate only
bounded evaluator work requests: the runner checkpoints the run with the
complete typed work request instead of calling a subprocess, and the harness
submits a typed result envelope through resume.

The harness evaluator is the built-in evaluator that declares
source-resolution support: alongside judgment requests, its checkpoint
transport carries `resolveSource` work requests — the selector, its detected
kind, the area frame, and an empty source bundle — whose returned
workspace-relative file paths the runner validates, rereads, and captures per the
[runner source contract](runner.md#source-packaging). A resolution result
envelope is the ordinary result envelope; when the material the selector
describes does not exist, the harness returns a classified
`source_unavailable` failure instead of improvised evidence.

The harness evaluator **MUST** declare subagent delegation: its checkpointed
requests are self-contained bounded evidence boundaries the invoking harness
may judge itself or fan out to native subagents. The runner **MUST NOT**
reduce a harness run's resolved concurrency to `1` on capability grounds; the
resolved concurrency bounds the outstanding checkpoint window per the
[runner harness contract](runner.md#harness-checkpoints).

> Rationale: the harness never takes simultaneous in-process calls, so
> concurrent-call support is the wrong capability to gate on; what makes the
> window serviceable is that each request is complete and independent, so the
> harness can delegate members without sharing judgment state. — 0198

A resume submission **MUST** be accepted as one result envelope or several —
any subset of the outstanding requests, one envelope per request — with each
member validated and correlated independently per the
[orchestration checkpoint contract](orchestration.md#harness-checkpoints).

The harness evaluator **MUST NOT** own run creation, scope expansion,
work-graph ordering, retry policy, result validation, persistence, report
generation, or final authority outside the result envelope it submits.

> Rationale: making the outer agent the workflow orchestrator would recreate
> the harness-dependent evaluation behavior 0192 removed. — 0194

A submitted harness result envelope **MUST** identify the pending request and
its input hash, **MUST** identify the harness runtime supplying the judgment,
and **MUST** carry either a result payload or a classified evaluator failure.
Optional usage and model metadata follow the ordinary result envelope.

The first accepted harness result **MUST** bind the run to the envelope's
harness runtime identity; later results from a different runtime **MUST** be
refused with `run_state_invalid`. The runner **SHOULD** record the model per
call when the harness reports it, without making model metadata a correctness,
resume, or validation dependency.

> Rationale: a resumable run must not silently mix Claude, Codex, or another
> harness after accepted judgments exist, while model metadata is not
> uniformly available across surfaces. — 0194

## Subagent-backed work

Where an evaluator exposes native subagents or worker threads,
subagent-backed work **MUST** return the same typed result envelope as any
other evaluator-backed work, and subagents **MUST NOT** write run artifacts,
expand scope, change dependency ordering, or produce final authority outside
the accepted result envelope.

## Prompt shaping and reusable context

The runner **MUST** render every work-request prompt with all content that is
stable across the subject area's work units — standing instructions, task,
expected schema, packaged source, and shared area-level context — preceding
any per-work-unit-varying content, and **MUST** expose the boundary between
the stable prefix and the per-work-unit delta to evaluator implementations.

> Rationale: a provider prefix cache is valid only up to the first byte that
> changes between calls. With source rendered after mutable per-requirement
> context, the largest repeated content could never be cached even though it
> repeats verbatim across an area's requirements. — 0193

Where an API-backed evaluator supports provider prompt caching, that evaluator
**MUST** apply the provider's caching mechanism to the stable prefix.

> Rationale: the earlier "SHOULD shape cache-friendly prefixes" left the layout
> fixed but no evaluator setting a cache breakpoint, so the intended saving was
> never realized on any path. Cache availability stays provider- and
> time-dependent: correctness comes from the work graph and persisted
> artifacts, not from a cache hit. — 0192, 0193

Where an evaluator supports reusable prompt or prefix state, the runner **MAY**
reuse the immutable area prefix. Every requirement still **MUST** run in a fresh
session or thread; one requirement's transcript **MUST NOT** become another's
context. Correctness **MUST NOT** depend on a provider cache hit.

Reusable evaluator context **MUST** be treated as reconstructible execution
metadata, not as authoritative run state. A resumed run **MUST NOT** require a
provider-retained conversation, thread, session, or cache entry that cannot be
recreated from the model snapshot, source package, runner prompts, and
`evaluation.json` state.

## Configuration

The workspace config file `.quality/config.yaml` names evaluators:

- `evaluation.evaluator` — `auto` or the name of a built-in evaluator or
  configured profile. When absent, the runner **MUST** behave as though
  `evaluation.evaluator: auto` were configured.
- `evaluation.concurrency` — optional positive integer maximum for active
  evaluator calls. When absent, the runner uses the
  [runner concurrency default](runner.md#concurrency).
- `evaluators` — an optional map of named profiles. Each profile **MUST**
  declare a `kind` — one of `codex`, `claude`, `openai`, or `anthropic`
  (`shell` and `manual` stay reserved) — and **MAY** declare `model`,
  `apiKeyEnv`, `baseUrl`, and `command`.

```yaml
evaluation:
  evaluator: team-openai
  concurrency: 8
evaluators:
  team-openai:
    kind: openai
    model: gpt-5
    apiKeyEnv: OPENAI_API_KEY
```

Built-in evaluator names **MUST** be reserved and **MUST NOT** be shadowed by
custom evaluator profiles. The reserved names are `auto`, `codex`, `claude`,
`openai`, `anthropic`, `shell`, and `manual`.

API-key evaluator profiles **MUST** reference secrets by environment-variable
name (`apiKeyEnv`), never by secret value.

> Rationale: config files are likely to be committed. A secret field would make
> accidental credential disclosure the path of least resistance. — 0192

Evaluator resolution order and `auto` discovery are the
[`qualitymd evaluation run`](../cli/evaluation-run.md#evaluator-selection)
command's contract.

Evaluators **MUST** declare whether they support concurrent calls. The runner
**MUST NOT** resolve concurrency above `1` for an evaluator that declares
neither concurrent-call support nor subagent delegation; for a checkpointed
evaluator, subagent delegation is the declaration that makes an outstanding
window above `1` serviceable.

## Runner authority and cancellation

The runner is the sole owner of the work graph, dependency readiness,
top-level concurrency, timeout, cancellation, retry budget, accepted-result
validation, deterministic persistence, and final report assembly. SDK agent
loops are workers inside one work unit.

Requirement judgment **MUST** disable nested provider subagents by default. A
source-resolution profile may opt in only when the adapter declares support and
the profile supplies explicit depth and concurrency caps.

Adapters **MUST** propagate runner cancellation to SDK streams and provider
child processes. A late result after cancellation or request supersession
**MUST NOT** be accepted. Provider session identifiers are diagnostic only;
resume remains based on runner input hashes and accepted artifacts.

## Safety, environment, and observability

Evaluators **MUST** treat captured source and tool output as untrusted data, not
instructions, and use the least tools, workspace access, network access,
sandbox authority, and approval authority needed for the work kind.

Provider child environments **MUST** be allowlisted. They may inherit common
process variables plus the selected evaluator's documented authentication and
configuration variables, but **MUST NOT** inherit unrelated credentials.
Configuration continues to store secret locators rather than secret values.

Evaluator-call logs and run metadata **MUST** identify evaluator kind/profile,
provider model when reported, duration, attempt, classified failure, declared
capabilities, and provider-reported usage when available. They **MUST NOT**
record raw prompts, source bodies, tool transcripts, result bodies, secrets, or
environment values.

`qualitymd` adds no project telemetry and does not enable provider telemetry on
the user's behalf. Provider data behavior inherent to a selected SDK or local
runtime belongs to that evaluator's documented boundary.

The complete agent session and context contract lives in
[Agent evaluators](agent-evaluators.md).
