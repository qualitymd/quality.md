---
type: Functional Specification
title: Evaluation runner token efficiency
description: Requirements for cache-friendly prompt layout, provider prompt caching, single-call requirement assessment-and-rating, and per-area source reuse in the evaluation runner.
tags: [evaluation, runner, evaluator, performance]
timestamp: 2026-07-09T00:00:00Z
---

# Evaluation runner token efficiency

This change makes the evaluation runner spend far fewer input tokens without
changing what it judges. It governs how the runner lays out evaluator prompts,
how many evaluator calls a requirement costs, and how an area's packaged source
is reused across the work that consumes it. It defers to the durable
[Evaluation runner](../../../specs/evaluation/runner.md),
[Evaluator contract](../../../specs/evaluation/evaluator-contract.md), and
[Evaluation orchestration](../../../specs/evaluation/orchestration.md) specs, which
this change amends, and inherits the observational-equivalence and
source-as-data invariants from [Evaluation](../../../specs/evaluation/evaluation.md).

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are to be
interpreted as described in BCP 14 when, and only when, they appear in all
capitals.

## Background / motivation

A full-model run on this repo spends ~8.8M input tokens against ~0.3M output
(96.7% input). The dominant term is the per-area source bundle, re-uploaded once
per requirement in the area, positioned in the prompt where provider prefix
caching cannot reach it; the runner also pays two evaluator calls per
requirement and reuses no context on the default CLI path. The evaluator
contract already prescribes "stable cache-friendly prefixes followed by
work-unit-specific deltas," so the primary gap is conformance, not new policy.
The binding constraint throughout: none of these optimizations may change a
rating, a persisted payload, or a report — caching and call-shape are cost
levers under the same work graph, never alternate orchestrators.

## Scope

Covered:

- prompt-layer ordering and the cache boundary between stable and per-work-unit
  content;
- provider prompt caching for API-backed evaluators, and cached-token usage
  observability;
- assessing and rating a requirement in one evaluator call;
- reuse of an area's packaged source across that area's work units, including
  per-area CLI session reuse where supported.

Deferred:

- requirement-scoped source excerpting or other packaging-heuristic changes
  beyond reuse;
- parallel and subagent execution strategies (owned by the runner spec's
  existing deferred slice).

Non-goals:

- changing rating semantics, roll-up, report contents, or the set of persisted
  payload kinds;
- making any provider API, prompt cache, or retained session a correctness,
  resume, or determinism dependency.

## Assumptions & dependencies

- The runner executes an area's requirement work units consecutively and in a
  deterministic order (per
  [orchestration](../../../specs/evaluation/orchestration.md#work-graph)); per-area
  source reuse and per-area session reuse rely on that adjacency for cache and
  session hits, but not for correctness.
- Providers report cached input tokens in usage
  (`cache_read_input_tokens` / equivalent); when a provider does not, the
  cached-token field is simply absent.

## Requirements

### Prompt layout and caching

`evaluator.BuildPrompt` **MUST** order the rendered prompt so that all content
that is stable across an area's work units — standing instructions, task,
expected schema, packaged source, and area-level frame context — precedes any
per-work-unit-varying context (the requirement and its requirement-level frame).

> Rationale: a provider prefix cache is valid only up to the first byte that
> changes between calls. Today the large source bundle sits after the mutable
> per-requirement context, so it can never be cached even though it repeats
> verbatim across an area's requirements.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` (Prompt shaping
> and reusable context) to require stable-before-delta ordering with source in
> the stable region.

`evaluator.BuildPrompt` **MUST** expose the boundary between the stable prefix
and the per-work-unit delta to evaluator implementations, so an evaluator can
mark or reuse the prefix without re-deriving it.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` (Work-unit
> envelope / Prompt shaping) to state that the rendered prompt carries a
> stable/delta split.

Where an API-backed evaluator supports provider prompt caching, that evaluator
**MUST** apply the provider's caching mechanism to the stable prefix.

> Rationale: the contract's existing "SHOULD shape cache-friendly prefixes" has
> a layout fix (above) but no evaluator actually sets a cache breakpoint, so the
> intended saving is never realized on any path. Strengthening to MUST for
> providers that support it closes the gap the measurement exposed.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` (Built-in
> evaluators / Prompt shaping) to require API evaluators to apply provider
> caching to the stable prefix.

Provider prompt caching **MUST NOT** be required for correctness, resume,
validation, or deterministic report output.

> Rationale: reaffirms the durable invariant against the new caching code; a
> cache miss must only cost tokens, never change output. This restates an
> existing runner-spec invariant intentionally because this change adds the
> code it constrains.
>
> Durable spec: none (already required by
> `specs/evaluation/runner.md#execution-strategy`).

When a provider reports cached input tokens, the evaluator **MUST** record them
in call usage, and the runner **MUST** include the cached-input-token count in
`logs/evaluator-calls.jsonl`.

> Rationale: without a cached-vs-fresh input signal the saving cannot be
> verified or regression-tested; the log is where per-call usage already lives.
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` (Result
> envelope usage) and `specs/evaluation/runner.md` (Logging) to add cached input
> tokens to reported usage and evaluator-call metadata.

### Single-call requirement assessment and rating

The runner **MUST** produce a requirement's `RequirementAssessmentResult` and
`RequirementRatingResult` from a single evaluator call.

> Rationale: the separate rate call re-ships context to score an assessment the
> same evaluator just produced, roughly doubling requirement calls and adding
> ~660k projected input on this repo for no new evidence.
>
> Durable spec: modify `specs/evaluation/orchestration.md` (Work graph) to make
> requirement assessment-and-rating one evaluator-backed unit; modify
> `specs/evaluation/records/payload-kinds.md` to state that one requirement
> judgment call yields both payload kinds.

The single requirement judgment call **MUST** persist the same two payload kinds,
under the same identities and schemas, that the separate assess and rate units
persist today, and **MUST** preserve every dependency the roll-up relies on: a
`RequirementRatingResult` still exists and is valid before any factor node that
depends on the requirement is analyzed.

> Rationale: the merge is a call-shape change, not a data-model change. Reports,
> roll-up, and resume must not be able to tell the difference.
>
> Durable spec: modify `specs/evaluation/orchestration.md` (Dependencies) to
> reflect that the rating dependency is satisfied by the combined unit.

If the combined call returns an assessment but no valid rating (or the reverse),
then the runner **MUST** treat the work unit as failed under the existing retry
and failure policy, rather than persisting a partial requirement result.

> Rationale: decides the divergence the merge introduces — a half-answered
> combined call — so it cannot silently persist an unrated requirement that
> roll-up would later trip over.
>
> Durable spec: modify `specs/evaluation/orchestration.md` (Retry and failure)
> to name partial combined-result output as a retryable unit failure.

### Per-area source reuse

The runner **MUST** package an area's source once per run and reuse that bundle
for every work unit in the area, rather than re-packaging per requirement.

> Rationale: packaging is deterministic, so re-packaging yields an identical
> bundle and identical `sourcePackageHash`; computing it once removes redundant
> filesystem work and gives evaluators a stable object to cache or reuse.
>
> Durable spec: modify `specs/evaluation/runner.md` (Source packaging) to state
> that an area's bundle is packaged once and reused across its work units, with
> unchanged determinism and hashing.

Where a CLI-backed evaluator supports session or thread continuation, the runner
**MAY** reuse one session per area so the area's source is transmitted once and
subsequent requirement calls in that area send only the delta.

> Rationale: the default `auto` path is the codex CLI, which provider prompt
> caching does not reach; per-area session reuse is the CLI analog and the only
> lever that helps the default run. Kept MAY because session-continuation
> semantics are provider-specific and must never be required for correctness.
>
> Durable spec: modify `specs/evaluation/runner.md` (Execution strategy) and
> `specs/evaluation/evaluator-contract.md` (Prompt shaping and reusable context)
> to name per-area session reuse as a permitted CLI context-reuse strategy under
> the existing reusable-context invariant.

Reused source bundles and reused evaluator sessions **MUST** remain
reconstructible execution metadata: a resumed run **MUST NOT** depend on a
retained session or in-memory bundle that cannot be rebuilt from the model
snapshot, packaged source, runner prompts, and `evaluation.json`.

> Rationale: extends the existing reusable-context invariant to the new reuse
> paths so this change cannot weaken resume.
>
> Durable spec: none (already required by
> `specs/evaluation/evaluator-contract.md#prompt-shaping-and-reusable-context`).

### Observability of the saving

`qualitymd evaluation run --dry-run --json` **MUST** report the requirement
judgment work-unit count consistent with one evaluator call per requirement.

> Rationale: the dry-run preview is the pre-run source of truth for work-unit
> counts; it must match the merged call shape so the preview and the run agree.
>
> Durable spec: modify `specs/cli/evaluation-run.md` only where it describes
> dry-run work-unit counts.

## Durable spec changes

### To add

None

### To modify

- `specs/evaluation/evaluator-contract.md` — stable-before-delta prompt ordering
  with source in the stable region; a rendered stable/delta split exposed to
  evaluators; API evaluators MUST apply provider caching to the stable prefix;
  cached input tokens in reported usage; per-area CLI session reuse as a
  permitted reusable-context strategy (per the prompt-layout, cache-boundary,
  API-caching, usage, and session-reuse requirements above).
- `specs/evaluation/runner.md` — an area's source is packaged once and reused
  across its work units; cached input tokens added to evaluator-call logging;
  per-area session reuse named under execution strategy (per the source-reuse,
  cached-token logging, and session-reuse requirements above).
- `specs/evaluation/orchestration.md` — requirement assessment-and-rating is one
  evaluator-backed work unit; the rating dependency is satisfied by the combined
  unit; a partial combined result is a retryable unit failure (per the
  single-call requirements above).
- `specs/evaluation/records/payload-kinds.md` — one requirement judgment call
  yields both `RequirementAssessmentResult` and `RequirementRatingResult`; the
  kinds, identities, and schemas are unchanged (per the single-call payload
  requirement above).
- `specs/cli/evaluation-run.md` — dry-run work-unit-count description reflects
  one evaluator call per requirement (per the observability requirement above).

### To rename

None

### To delete

None

## Open questions

None. Per-area session reuse stays `MAY` because provider session semantics are
not uniform; if a concrete CLI proves it reliably, a later change can strengthen
it.
