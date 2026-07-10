---
type: Design Doc
title: Evaluation runner token efficiency design
description: How the runner and evaluator adapters deliver cache-friendly prompt layout, single-call requirement judgment, and per-area source reuse without changing evaluation output.
tags: [evaluation, runner, evaluator, design]
timestamp: 2026-07-09T00:00:00Z
---

# Evaluation runner token efficiency design

## Context

This design answers the
[Evaluation runner token efficiency functional spec](spec.md). The runner spends
~8.8M input tokens against ~0.3M output on a full-model run because the per-area
source bundle is re-uploaded per requirement past the prompt cache boundary,
every requirement costs two evaluator calls, and the default CLI path reuses no
context. The three levers are independent and land together; the single binding
constraint is observational equivalence — persisted payloads and reports stay
byte-equivalent, so every change is transport- or scheduling-shaped, never a
change to judgment.

The work touches `internal/evaluator/` (prompt rendering, API/CLI adapters,
`WorkRequest`/`Usage` types) and `internal/runner/` (work graph, request
assembly, engine result handling, source packaging, logging).

## Approach

### 1. A prompt with an explicit stable/delta boundary

`evaluator.BuildPrompt` today concatenates one user string with source rendered
_last_, after the mutable per-requirement `Context`. Change it to render
stable-first and to return the boundary:

```
system
──────────── stable prefix ────────────
Work-unit header · Task · Expected schema
Source files (data under evaluation)
Stable context (area-level frame)
──────────── delta ────────────
Per-work-unit context (requirement + requirement frame)
"Return ONLY the JSON object now."
```

`BuildPrompt` returns `(system string, stablePrefix string, delta string)`
instead of `(system, user)`. To classify context, split `WorkRequest.Context`
into two maps: `SharedContext` (stable across the area's units — the
`areaEvaluationFrame`) and `Context` (per-unit — `requirement`,
`requirementEvaluationFrame`). `WorkRequest.Source` is already a separate field
and moves into the stable region unchanged.

CLI evaluators keep today's behavior by concatenating
`system + "\n\n" + stablePrefix + delta`; only the _ordering_ changes, which is
itself a cache win for providers with automatic prefix caching. API evaluators
use the split to place a cache breakpoint (below).

`PromptPrefixHash` widens from `hash(instructions, schema)` to cover the whole
stable prefix (instructions, schema, source hash, shared context), so the log
records what is actually cacheable.

### 2. Provider caching in the API adapters

**Anthropic** (`anthropicEvaluator.Evaluate`): send `system` and the stable
prefix as cacheable content blocks and the delta as a trailing plain block:

```go
"system": []map[string]any{{"type": "text", "text": system,
    "cache_control": map[string]any{"type": "ephemeral"}}},
"messages": []map[string]any{{"role": "user", "content": []map[string]any{
    {"type": "text", "text": stablePrefix,
        "cache_control": map[string]any{"type": "ephemeral"}},
    {"type": "text", "text": delta},
}}},
```

Parse `usage.cache_read_input_tokens` and `cache_creation_input_tokens`.
**OpenAI** benefits from ordering alone (automatic prefix caching); read
`usage.prompt_tokens_details.cached_tokens`. `Usage` gains
`CachedInputTokens *int64`; `logs/evaluator-calls.jsonl` records it.

### 3. One evaluator call per requirement

Replace the `assessRequirement` + `rateRequirement` unit pair with a single
evaluator-backed `assessRateRequirement` unit:

- **Graph** (`addRequirementUnits`): emit one unit per requirement; map
  `rateUnits[req.Ref]` to its ID so factor-analysis dependencies retarget to it
  with no downstream graph edits.
- **Request** (`fillAssessRateRequest`): package source + both frames + combined
  instructions; `ExpectedSchema` is a composite
  `{"assessment": RequirementAssessmentResult, "rating": RequirementRatingResult}`
  built from the two existing typed schemas.
- **Engine**: on accept, split the composite into two payloads, validate each
  against its own `DataKind`, and persist both through the existing
  `acceptUnit(unit, state, []map[string]any{assessment, rating}, hash)` —
  `Results.Merge(unit.ID, payloads)` already stores multiple payloads per unit.
  Downstream reads (`findingIndex`, factor-analysis `directRequirementRatings`)
  resolve the requirement's assessment and rating from the combined unit's
  results by kind rather than by a separate rate-unit ID.
- **Partial result**: if the composite lacks a valid `assessment` or `rating`,
  the unit fails under the existing retry policy and persists nothing, so
  roll-up never sees an unrated requirement.

### 4. Per-area source reuse

- **Packaging** (`source.go` + engine): memoize the packaged bundle by area
  reference for the run (`map[string]*SourceBundle` on the engine). Packaging is
  deterministic, so the memoized bundle is identical to a re-packaged one, same
  `sourcePackageHash`; this removes redundant directory walks and gives adapters
  one stable object to cache or transmit.
- **CLI session reuse (MAY)**: where a CLI evaluator supports continuation
  (codex `exec` thread, claude `--resume <session>`), the runner passes a
  per-area reuse key; the adapter keeps the last session per key and, for a
  second-or-later unit in the same area, resumes it and sends only the delta.
  Off by default and independent of correctness: a dropped session just resends
  the prefix. The `threadId`/`sessionId` already captured in `ContextMeta` is the
  handle.

## Spec response

- **Prompt layout and caching** — §1 gives the stable-before-delta ordering and
  the exposed boundary; §2 applies provider caching and records cached tokens.
  The [caching-not-required-for-correctness requirement](spec.md#prompt-layout-and-caching)
  holds because the split is presentational: a cache miss re-sends identical
  bytes and produces the same result.
- **Single-call requirement assessment and rating** — §3 persists the same two
  payload kinds under the same identities and preserves the rating dependency
  (the combined unit is the factor-analysis dependency). The partial-result rule
  is the retry-failure path.
- **Per-area source reuse** — §4 packages once and reuses; session reuse stays
  reconstructible because every unit's delta fully names its subject, so a run
  can resume with a cold session or freshly packaged bundle.
- **Observability** — §2's cached-token logging plus the merged work-unit count
  surfaced by `--dry-run --json` (the graph emits one evaluator unit per
  requirement, so the count follows automatically).

## Alternatives

- **Cache only, keep two calls.** Rejected: leaves the 2× call structure and its
  ~660k input, and helps only API evaluators — the default codex path still pays
  full price. The merge is transport-independent.
- **A new combined payload kind.** Rejected: would change the
  [payload-kinds](../../../specs/evaluation/records/payload-kinds.md) contract and
  the report inputs. Splitting the composite back into the two existing kinds
  keeps the data model and reports untouched.
- **One session for the whole run.** Rejected: source differs per area, so a
  single thread would accumulate unbounded, mixed context. The area is the
  natural cache/session unit and bounds context growth.
- **Requirement-scoped source excerpting.** Deferred (not rejected): it attacks
  the same dominant term but changes _which_ evidence each call sees, which can
  move a rating — out of scope for a change whose whole premise is unchanged
  output.

## Trade-offs & risks

- **Bigger single response.** The combined call returns assessment + rating in
  one object, a larger and more complex output; a malformed half-response fails
  the whole requirement. Mitigated by the existing retry loop and the explicit
  partial-result failure rule.
- **Provider-specific caching code.** `cache_control` and cached-usage parsing
  are Anthropic-shaped; guarded so a miss or an unsupported provider costs only
  tokens. OpenAI needs no markers.
- **Session reuse is subtle.** Context bleed across an area's requirements is the
  real hazard; kept `MAY` and off by default, with each delta fully specifying
  its subject so a stale or resumed session cannot change judgment. It is the
  last piece to enable and the easiest to defer if a provider misbehaves.
- **Resume cache keys shift.** New unit IDs and the widened `PromptPrefixHash`
  mean an in-flight run started before this change won't reuse prior results on
  resume. Acceptable under early-alpha; a fresh run is unaffected.

## Open questions

None blocking. Sequencing within the case: land the prompt reorder, provider
caching, and the assess+rate merge first (they carry most of the saving and are
transport-independent for the merge); enable CLI per-area session reuse last,
behind its `MAY`, once codex/claude continuation behavior is verified against a
real run.
