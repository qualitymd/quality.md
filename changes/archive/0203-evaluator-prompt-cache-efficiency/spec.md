---
type: Functional Specification
title: Evaluator prompt cache efficiency — functional spec
description: Requirements for deterministic cacheable evaluator prefixes and complete provider cache-usage telemetry while preserving fresh-session isolation.
tags: [evaluation, evaluator, prompt-caching, tokens, codex, claude]
timestamp: 2026-07-15T00:00:00Z
---

# Evaluator prompt cache efficiency — functional spec

Companion to the
[Evaluator prompt cache efficiency](../0203-evaluator-prompt-cache-efficiency.md)
Change Case. This spec defines what must change; the design doc will own the
prompt-part representation, provider option shape, and test seams.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

Fresh requirement sessions prevent order-dependent judgment and transcript
leakage, but they repeatedly send the same evaluator policy, task instructions,
model guidance, area frame, tools, and output schema. Prompt caching can
amortize provider prefill only when those repeated inputs form an exact prefix.
Putting work-unit identity before shared context leaves the large repeated
material after the first difference and turns an intended caching contract into
an implementation comment rather than observable source behavior.

Caching is an economic and latency optimization, not a smaller logical context
window and not durable evaluation state. The evaluator must therefore remain
correct on every cache miss, while run-local usage preserves enough provider
telemetry to show whether the optimization is working.

## Requirements

### R1 — Layered evaluator prompt

The evaluator prompt renderer **MUST** produce two explicit parts: a cacheable
prefix and a work-unit-specific suffix. The complete prompt **MUST** be exactly
those parts joined by one stable boundary.

> Rationale: explicit parts make prefix equality testable without logging raw
> prompts and keep later edits from accidentally moving a varying field ahead
> of shared context. — 0203
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` — replace the
> advisory prompt-order statement with the layered prefix/suffix contract.

For requests with the same evaluator kind, instructions, body guidance, shared
context, and inspection context, the cacheable prefix **MUST** be byte-identical
even when run ID, work-unit ID, subject, and work-unit context differ. The
suffix **MUST** carry work-unit ID, subject, and work-unit context.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — define the
> stable and varying prompt inputs.

The prefix **MUST** place evaluator safety/output policy first, followed by
work-kind instructions, applicable QUALITY.md body guidance, shared accepted
context, and inspection availability and policy. No run identity, work-unit
identity, subject, or work-unit context may appear before the stable boundary.

> Durable spec: modify `specs/evaluation/evaluator-contract.md` — define prompt
> ordering; modify `specs/evaluation/agent-evaluators.md` — apply it to fresh
> agent inspection sessions.

### R2 — Deterministic structured blocks

Structured values rendered into either prompt part **MUST** use canonical JSON
with recursively sorted object keys and stable escaping. Semantically identical
objects that differ only in insertion order **MUST** render identically.

> Rationale: provider caches match bytes, while JavaScript object insertion
> order is not part of the evaluator contract. Canonicalization prevents an
> internal construction detail from becoming a cache miss. — 0203
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` — require
> deterministic serialization for structured prompt blocks.

### R3 — Cache-usage preservation

Evaluator usage **MUST** represent provider-reported cache reads separately
from ordinary input tokens and provider-reported cache writes separately from
both. An unavailable metric **MUST** remain absent; a reported zero **MUST**
remain zero.

> Rationale: absence and zero answer different operational questions, and cache
> creation cost cannot be evaluated if it is folded into ordinary input. — 0203
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` — extend the
> provider-neutral usage contract; modify `specs/evaluation/runner.md` — extend
> evaluator-call logging.

The Codex adapter **MUST** preserve every cache metric exposed by its pinned SDK.
The Claude adapter **MUST** preserve its reported cache-read and cache-creation
input-token counts. The runner **MUST** carry those optional values unchanged
into the call's `usage` object in `logs/evaluator-calls.jsonl` and **MUST NOT**
persist them in `evaluation.json`.

> Durable spec: modify `specs/evaluation/agent-evaluators.md` — state provider
> adapter reporting; modify `specs/evaluation/runner.md` — state the run-local
> log boundary.

### R4 — Claude cache-stable system prefix

The Claude adapter **MUST** use the pinned SDK's Claude Code preset system prompt
with dynamic runtime sections excluded from the globally cacheable system
prefix and reintroduced through the SDK's supported dynamic path. It **MUST**
retain the existing neutral temporary directory, empty setting sources,
non-persisted session, read/search-only tools, and disallowed mutation, shell,
and nested-agent tools.

> Rationale: every session receives a random neutral directory. Letting that
> path vary inside the system prefix prevents cross-session reuse even though it
> is not evaluator authority. The SDK provides a supported split that preserves
> the context after the cache boundary. — 0203
>
> Durable spec: modify `specs/evaluation/agent-evaluators.md` — record the
> cache-stable preset configuration within the unchanged neutral-session policy.

### R5 — Fresh-session and cache independence

Every requirement **MUST** continue to start a fresh, non-persisted evaluator
session with no sibling transcript, earlier inspection transcript, shared
provider conversation, resumed response chain, or forked session. Cache
availability, cache contents, and cache telemetry **MUST NOT** affect request
identity, acceptance, retry, resume, ratings, evidence, reports, or output
ordering.

> Rationale: provider key/value reuse is safe only because it is invisible to
> judgment semantics; transcript reuse would couple otherwise independent work
> and still occupy the logical context window. — 0203
>
> Durable spec: modify `specs/evaluation/evaluator-contract.md` and
> `specs/evaluation/agent-evaluators.md` — preserve the no-transcript-reuse
> boundary while explaining why prompt caching is permitted.

### R6 — Verification and observability

Focused tests **MUST** prove exact prefix equality across differing work units,
prefix divergence when a declared prefix input changes, suffix variation,
canonical structured rendering, cache field preservation including zero versus
absence, Claude cache-stable system configuration, and run-local logging without
`evaluation.json` leakage.

> Durable spec: none.

Before review, two otherwise identical small scoped live evaluations **MUST** be
run within the provider cache lifetime. The review ledger **MUST** record scope,
evaluator/model, concurrency, per-call input/cache-read/cache-write values when
reported, durations, and limits on attribution. A missing or zero cache read is
valid evidence to revisit prompt shaping; it **MUST NOT** be rewritten as a
successful cache hit.

> Rationale: aggregate coding-agent usage can include both cross-session and
> within-session tool-loop caching, so the evidence must be reported honestly
> rather than claimed as a causal benchmark. — 0203
>
> Durable spec: none — this is change acceptance evidence, not enduring runtime
> behavior.

## Requirement-set review

R1 and R2 make the repeated provider input exact and testable. R3 makes cache
reads and writes visible without promoting provider state into the evaluation
artifact. R4 applies the same principle to the provider-owned Claude system
prompt while preserving every existing safety option. R5 prevents the
optimization from becoming transcript coupling or a correctness dependency.
R6 proves the source contract and records bounded live evidence without
over-claiming attribution. Together they achieve the motivation — cache-friendly
fresh sessions with measurable provider reuse — while leaving scheduling,
judgment, artifacts, and public workflow behavior unchanged.

Every requirement has an observable verification path: pure prompt and usage
tests for R1–R3, option-shape plus live Claude coverage for R4, existing
fresh-session/integration contracts for R5, and the explicit test/live-evidence
ledger for R6. No unresolved design decision remains in the requirement set.

## Durable spec changes

### To add

None.

### To modify

- `specs/evaluation/evaluator-contract.md` — layered deterministic prompt
  prefix, cache usage, and cache-independent fresh-session rules (R1–R3, R5).
- `specs/evaluation/agent-evaluators.md` — agent-session prompt shaping, Claude
  cache-stable preset use, provider usage reporting, and no transcript reuse
  (R1, R3–R5).
- `specs/evaluation/runner.md` — separate optional cache-read and cache-write
  token counts in run-local evaluator-call logs, never `evaluation.json` (R3).

### To rename

None.

### To delete

None.
