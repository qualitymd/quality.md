---
type: Functional Specification
title: Provider-affine SDK evaluator selection — functional spec
description: Requirements for selecting a provider-affine SDK evaluator before harness judgment and explaining the determined transport without a user choice.
tags: [evaluation, evaluator, selection, skill, sdk, harness]
timestamp: 2026-07-15T00:00:00Z
---

# Provider-affine SDK evaluator selection — functional spec

Delta contract for change case
[0208 — Provider-affine SDK evaluator selection](../0208-provider-affine-sdk-selection.md).
Normative sources of truth this spec changes:
[`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
(evaluate summary),
[`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
(evaluation wrapper), and
[`specs/skills/quality-skill/workflows/evaluate.md`](../../../specs/skills/quality-skill/workflows/evaluate.md)
(ordered evaluate workflow). The CLI and evaluator specs are informational:
their automatic discovery and transport contracts do not change.

The key words "MUST", "MUST NOT", and "MAY" are to be interpreted as described
in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The current skill prefers in-session harness judgment over ready SDK-backed
evaluators. That default weakens isolation and makes the common path depend on
the invoking conversation and checkpoint servicing even though the runner can
use fresh provider sessions. Provider-named requests also incur a choice whose
answer can be determined consistently: a provider named as evaluator means its
independent SDK evaluator; the current session remains available through an
explicit `harness` request or as the no-SDK fallback.

## Requirements

### R1 — automatic SDK discovery precedes harness fallback

When neither an explicit evaluator request nor a non-`auto` workspace
`evaluation.evaluator` determines the evaluator, `/quality evaluate` **MUST**
inspect CLI `auto` discovery before selecting `harness`, and **MUST** select
`harness` automatically only when discovery reports no usable built-in SDK
evaluator and the invoking harness can service checkpoints.

> Rationale: harness judgment is a compatibility path; a ready SDK evaluator
> gives the runner fresh, isolated judgment sessions without the invoking
> conversation becoming the evidence basis. — 0208
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — move automatic discovery
> ahead of the harness fallback.

### R2 — automatic selection prefers the invoking provider

When CLI discovery reports that the built-in SDK evaluator matching a known
Codex or Claude invoking harness is usable, `/quality evaluate` **MUST** select
that evaluator explicitly. When no matching provider candidate is usable but
another built-in candidate is usable, the workflow **MUST** use the evaluator
selected by CLI `auto` discovery.

> Rationale: a Claude-mediated evaluation should not silently select Codex only
> because standalone discovery has a fixed Codex-first order, and vice versa;
> provider affinity gives the agent-mediated workflow a predictable default
> without teaching the standalone CLI to guess its parent. — 0208
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — add provider-affine
> selection over the CLI candidate receipt.

### R3 — provider-named requests resolve to the SDK evaluator

When an explicit evaluator request names Codex or Claude without naming a
transport, `/quality evaluate` **MUST** select that provider's built-in SDK
evaluator without asking the user to choose between SDK and harness transport.
The workflow **MUST** select the current session only for an explicit `harness`
request or through the automatic no-SDK fallback in R1.

> Rationale: in an evaluation request, the provider name identifies the
> independent evaluator; asking a second transport question adds friction and
> makes the same words resolve differently by harness. — 0208
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — replace provider-name
> disambiguation with deterministic SDK mapping.

### R4 — selection is visible before mutation

Before creating the evaluate feedback log or invoking the runner,
`/quality evaluate` **MUST** state the selected evaluator and the precedence,
provider-affinity, CLI discovery result, or no-SDK fallback that determined it.
The explanation **MUST** distinguish a fresh independent SDK evaluator from
current-session harness judgment and **MUST NOT** invite a transport choice for
the current run.

> Rationale: evaluator identity changes how a reader interprets the evidence
> basis, but observability does not require blocking every run on a preference
> question. — 0208
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md`,
> `specs/skills/quality-skill/workflows/evaluate.md`, and
> `specs/skills/quality-skill/quality-skill.md` — make determined selection and
> its reason part of the pre-mutation interface.

### R5 — selected runs remain evaluator-pinned

Automatic fallback in R1 and R2 **MUST** finish before run creation. After the
runner records the selected evaluator, `/quality evaluate` **MUST NOT** switch
providers or cross between SDK and harness transport; a different evaluator
requires a new run. An unavailable explicit or non-`auto` configured evaluator
**MUST** surface its concrete remedy rather than entering the automatic
fallback chain.

> Rationale: fallback is safe while determining a new run's evaluator, not
> after judgment has become attributable to one evaluator. Explicit and
> configured intent must not be weakened by an unrelated available runtime.
> — 0208
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — preserve evaluator
> pinning while distinguishing pre-run automatic fallback.

## Requirement-set check

Consistent: R1 establishes SDK-before-harness precedence; R2 decides which
usable SDK wins in an agent-mediated invocation; R3 decides explicit
provider-name semantics; R4 makes the result observable without a choice; and
R5 bounds fallback to pre-run determination. Complete: explicit intent,
configuration, automatic provider affinity, alternate SDK discovery, harness
fallback, explanation, failure, and post-selection pinning each have one home.
Able to be validated: durable/runtime source review can trace the full decision
table, targeted searches can prove the old choice and harness-first rules are
gone from active surfaces, and the existing CLI tests plus repository and
release gates prove that the unchanged discovery receipt remains compatible.
Satisfying R1-R5 achieves SDK-first provider-affine selection without changing
standalone CLI behavior or the QUALITY.md format.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` — selected evaluator visibility
  in the evaluate summary (R4).
- `specs/skills/quality-skill/evaluation.md` — SDK-before-harness precedence,
  provider affinity, provider-name mapping, visible determination, and pre-run
  fallback boundary (R1-R5).
- `specs/skills/quality-skill/workflows/evaluate.md` — ordered selection and
  invocation behavior (R1-R5).

### To rename

None

### To delete

None
