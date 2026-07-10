---
type: Functional Specification
title: Evaluator contract
description: Capability, work-unit envelope, result envelope, and configuration contract for evaluation runner evaluators.
tags: [evaluation, evaluator, agents]
timestamp: 2026-07-09T00:00:00Z
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

Every evaluator **MUST** declare its execution capabilities — supported
execution strategies, subagent support, reusable-context kinds, and usage
reporting — and the runner **MUST** read the declaration before dispatching
work.

> Rationale: the runner can only choose strategies an evaluator provably
> supports; assuming undeclared capabilities reintroduces harness-dependent
> behavior. — 0192

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
- the execution strategy used;
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

The built-in evaluator kinds are `codex` (Codex CLI subprocess), `claude`
(Claude Code CLI subprocess), `openai` (direct OpenAI API), and `anthropic`
(direct Anthropic API). `shell` and `manual` are reserved names with no
implementations.

CLI-backed evaluators **MUST** be invoked non-interactively with
machine-readable structured output — Claude Code in print mode with JSON
output, Codex in exec mode with JSON output — which the runner validates
against the work unit's expected schema.

If an installed CLI cannot honor non-interactive structured invocation, then
evaluator selection **MUST** fail with `evaluator_incompatible` and report
remediation, including any other available evaluators.

> Rationale: the runner's determinism rests on structured evaluator output. A
> CLI version that cannot provide it must fail loudly at selection, not degrade
> into unparseable runs. — 0192

API-backed evaluators **MUST** read their API key from the profile's configured
environment variable. If that variable is unset when an API-backed evaluator is
selected, then selection **MUST** fail with `missing_api_key`.

## Subagent-backed work

Where a CLI-backed evaluator exposes native subagents or worker threads,
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

Where an evaluator supports reusable conversation, thread, session, or previous
response state, the runner **MAY** create shared base context and fork or chain
work units from it. In particular, where a CLI-backed evaluator supports
session or thread continuation, the runner **MAY** reuse one session per area
so the area's stable prefix is transmitted once and later work units in that
area send only their deltas.

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
- `evaluators` — an optional map of named profiles. Each profile **MUST**
  declare a `kind` — one of `codex`, `claude`, `openai`, or `anthropic`
  (`shell` and `manual` stay reserved) — and **MAY** declare `model`,
  `apiKeyEnv`, `baseUrl`, and `command`.

```yaml
evaluation:
  evaluator: team-openai
  executionStrategy: auto # see the evaluation runner spec
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
