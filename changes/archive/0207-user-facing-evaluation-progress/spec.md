---
type: Functional Specification
title: User-facing evaluation progress — functional spec
description: Requirements for keeping evaluator protocol behind the agent interface and presenting evaluation phases, choices, coverage, and recovery in user-facing language.
tags: [evaluation, skill, agent-mediated-ux, progress]
timestamp: 2026-07-15T00:00:00Z
---

# User-facing evaluation progress — functional spec

Delta contract for change case
[0207 — User-facing evaluation progress](../0207-user-facing-evaluation-progress.md).
Normative sources of truth this spec changes:
[`specs/skills/quality-skill/quality-skill.md`](../../../specs/skills/quality-skill/quality-skill.md)
(shared user-interaction contract),
[`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
(evaluation wrapper), and
[`specs/skills/quality-skill/workflows/evaluate.md`](../../../specs/skills/quality-skill/workflows/evaluate.md)
(evaluate workflow). The runner, protocol, and CLI specs are informational:
their mechanics do not change.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

The agent is the `/quality` interface, but current evaluate guidance leaks the
runner's protocol into ordinary progress. A user sees work-unit categories,
outstanding request windows, payload schemas, concurrency caps, worker fan-out,
and resume-loop mechanics where they need only the current phase, meaningful
coverage, whether attention is required, and what comes next. The same guidance
can invite a current-run evaluator preference and continue without waiting, and
its prescribed order writes a feedback log before announcing the "first"
mutation. Protecting the interface boundary makes status truthful and useful
without weakening diagnostics or runner ownership.

## Requirements

### R1 — protect the implementation boundary

The `/quality` skill MUST present ordinary workflow state in terms of the
user's quality task — status, scope or coverage, attention needed, artifact or
result, and next action — and MUST NOT expose evaluator protocol, serialized
payload, worker orchestration, or command-loop mechanics unless a specific
detail is necessary for a user decision or recovery action.

> Rationale: protocol detail can be technically correct and still make the user
> reconstruct the product state from implementation internals. Necessary
> recovery details remain available; routine plumbing does not become the UI.
> — 0207
>
> Durable spec: modify `specs/skills/quality-skill/quality-skill.md` — add the
> shared implementation-boundary and progressive-disclosure rule.

### R2 — evaluate progress uses quality-task phases

During a running evaluation, `/quality evaluate` MUST summarize progress at
meaningful quality-task phase changes, such as preflight, evidence review,
report generation, and completion, and MUST state whether the user needs to act
when that is not otherwise obvious.

> Rationale: phase changes are stable user concepts even when the runner changes
> its graph, batching, evaluator transport, or concurrency implementation.
> — 0207
>
> Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` —
> replace protocol-state progress with phase-based progress; modify
> `specs/skills/quality-skill/evaluation.md` — align the wrapper flow.

### R3 — progress counts describe meaningful coverage

When `/quality evaluate` includes quantitative progress, it MUST use a count the
user can interpret as evaluation coverage, such as modeled areas or requirements
reviewed, and MUST NOT present work-unit counts, outstanding-request windows,
concurrency caps, active-worker counts, or payload counts as ordinary progress.

> Rationale: an exact internal count is not useful when the user cannot tell
> what one unit means or how heterogeneous units contribute to completion.
> — 0207
>
> Durable spec: modify `specs/skills/quality-skill/workflows/evaluate.md` —
> replace the outstanding-cap obligation with meaningful optional coverage;
> modify `specs/skills/quality-skill/evaluation.md` — keep harness request-window
> mechanics internal to the wrapper procedure.

### R4 — evaluator alternatives do not create false choices

When default precedence selects `harness` and the user did not request an
evaluator choice, `/quality evaluate` MUST present the selection as concise
information, state that it uses the current session, and MAY name the explicit
or configured independent-evaluator path for a future invocation. If the
workflow offers changing the evaluator for the current run, it MUST render a
real choice and wait for the answer before any evaluation mutation.

> Rationale: "unless you prefer otherwise" is a choice, not decoration. Either
> the default is settled and informational, or the workflow gives the user a
> gate and waits. — 0207
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` — supersede
> the current-run alternative requirement from 0206 with the
> informational-or-gated distinction; modify
> `specs/skills/quality-skill/workflows/evaluate.md` — align the workflow step.

### R5 — evaluator selection precedes the first write

`/quality evaluate` MUST resolve and explain evaluator selection before it
creates the evaluate feedback log or invokes the runner, and its immediately
preceding progress beat MUST accurately identify that evaluation artifacts and
the local feedback log are about to be written.

> Rationale: the current ordered procedure opens the feedback log before a
> later message calls the runner the first mutation. A user-visible boundary
> must agree with the actual mutation order. — 0207
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` and
> `specs/skills/quality-skill/workflows/evaluate.md` — move selection ahead of
> feedback-log creation and align the pre-mutation beat.

## Requirement-set check

Consistent: R1 is the shared boundary; R2–R3 specialize ordinary evaluation
progress; R4 decides when evaluator-selection prose is informational versus a
gate; R5 makes the mutation claim truthful. None changes runner mechanics,
evaluator precedence, or recovery. Complete: the observed leaks (protocol
jargon, worker narration, meaningless counts), false preference invitation,
and mutation-order contradiction each have one governing requirement. Able to
be validated: source review can trace each requirement from durable contracts
to runtime guidance and examples; the repository and release gates prove the
edited bundles, links, formatting, and package metadata remain coherent.

## Durable spec changes

### To add

None

### To modify

- `specs/skills/quality-skill/quality-skill.md` — shared implementation boundary
  and diagnostic progressive disclosure (R1).
- `specs/skills/quality-skill/evaluation.md` — phase-based wrapper progress,
  internal harness mechanics, evaluator alternative posture, and selection
  before first write (R2–R5).
- `specs/skills/quality-skill/workflows/evaluate.md` — ordinary progress counts,
  evaluator choice behavior, and pre-mutation workflow order (R2–R5).

### To rename

None

### To delete

None
