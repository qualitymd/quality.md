---
type: Change Case
title: Evaluation runner token efficiency
description: Cut evaluation input-token cost by making the assess prompt's stable prefix cacheable, collapsing requirement assess+rate into one evaluator call, and reusing packaged source across an area's work units.
status: Done
tags: [evaluation, runner, evaluator, performance]
timestamp: 2026-07-09T00:00:00Z
---

# Evaluation runner token efficiency

## Motivation

A full-model run measured on this repo's `QUALITY.md` spends **~8.8M input
tokens against ~0.3M output** — a 29:1 input:output ratio, with input at 96.7%
of all tokens. The cost is not judgment; it is re-shipping the same context.
Three structural causes, all confirmed in `logs/evaluator-calls.jsonl`:

1. **The stable prefix is never cached.** `evaluator.BuildPrompt` emits the
   large per-area source bundle _after_ the mutable per-requirement context
   block, so any provider prefix cache breaks before it reaches the source. The
   evaluator contract already asks for "stable cache-friendly prefixes followed
   by work-unit-specific deltas," but the implementation's layout defeats it and
   no evaluator sets a cache breakpoint. The same area source is re-uploaded for
   every requirement in the area (observed: `cli` ×4, `format-spec` ×5,
   `quality-skill` ×4, byte-identical `sourcePackageHash`, back-to-back).
2. **Every requirement costs two evaluator calls.** `assessRequirement` and
   `rateRequirement` are separate judgment units; the rate call re-ships context
   to score an assessment the same evaluator just produced (~660k projected
   input across the run).
3. **CLI evaluators reuse nothing.** The default `auto` evaluator (codex CLI)
   runs a fresh subprocess per call and discards the `threadId`/`sessionId` it
   captures, so none of the caching design reaches the default path.

These are transport and layout inefficiencies, not evaluation-correctness
issues. Fixing them lowers cost and latency while keeping ratings identical.

## Scope

Covered: prompt layout and provider prompt caching for API evaluators; a single
evaluator call that both assesses and rates a requirement; and reuse of an
area's packaged source across its work units, including per-area CLI session
reuse where the evaluator supports it. Usage logging gains cached-input-token
visibility so the win is measurable.

Deferred / non-goals: changing rating semantics, roll-up, report contents, or
persisted payload kinds; requirement-scoped source excerpting (packaging
heuristics stay as-is beyond reuse); making any provider API or caching a
correctness dependency; parallel execution (owned by 0192's deferred slice).

## Affected artifacts

- **Code:** `internal/evaluator/prompt.go` (cache-boundary segmentation, source
  before mutable context), `internal/evaluator/api.go` (Anthropic/OpenAI cache
  markers, cached-token usage), `internal/evaluator/cli.go` (stable-prefix
  prompt assembly, cached-token usage parsing for the claude and codex CLIs),
  `internal/evaluator/evaluator.go` (`WorkRequest.SharedContext`, `Usage`
  cached-token field); `internal/runner/requests.go` (combined assess+rate
  request, widened prompt-prefix hash), `internal/runner/graph.go` (requirement
  judgment as one evaluator-backed unit), `internal/runner/engine.go` (split the
  combined result into the two persisted payloads; record cached input tokens
  in evaluator-call log entries), `internal/runner/source.go` (per-area bundle
  reuse). Per-area CLI session reuse (the spec's `MAY`) is implemented by
  nothing yet: it stays a permitted strategy in the durable specs, to be
  enabled once codex/claude continuation behavior is verified against a real
  run.
- **Durable specs:** modify `specs/evaluation/evaluator-contract.md`,
  `specs/evaluation/runner.md`, `specs/evaluation/orchestration.md`, and
  `specs/evaluation/records/payload-kinds.md`; touch
  `specs/cli/evaluation-run.md` only for the dry-run work-unit-count description.
  See the spec's [Durable spec changes](0193-evaluation-token-efficiency/spec.md#durable-spec-changes).
- **Bundled skill:** no behavior change. `skills/quality/SKILL.md` and
  `workflows/evaluate.md` describe dry-run work-unit counts generically and do
  not hardcode the assess/rate split, so no runtime skill edit is required.
- **Docs:** none. No README/guide/Mintlify page documents the per-requirement
  call count or prompt layout.
- **Generated artifacts:** unchanged. Persisted payloads and reports are
  byte-equivalent to a pre-change run for the same model and evidence.

## Children

- [Functional spec](0193-evaluation-token-efficiency/spec.md)
- [Design doc](0193-evaluation-token-efficiency/design.md)
