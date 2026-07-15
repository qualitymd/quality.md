---
type: Functional Specification
title: Intent-faithful evaluator selection — functional spec
description: Requirements for probe-all auto discovery, candidate and selection-reason reporting, authentication basis, and skill-side transport disambiguation.
tags: [evaluation, evaluator, selection, discovery, skill]
timestamp: 2026-07-15T00:00:00Z
---

# Intent-faithful evaluator selection — functional spec

Delta contract for change case
[0206 — Intent-faithful evaluator selection](../0206-intent-faithful-evaluator-selection.md).
Normative sources of truth this spec's changes land in:
[`specs/cli/evaluation-run.md`](../../../specs/cli/evaluation-run.md) (CLI
evaluator selection), [`specs/evaluation/agent-evaluators.md`](../../../specs/evaluation/agent-evaluators.md)
(agent-runtime transport and authentication policy), and
[`specs/skills/quality-skill/evaluation.md`](../../../specs/skills/quality-skill/evaluation.md)
(skill-side selection and transport explanation). The
[evaluator contract](../../../specs/evaluation/evaluator-contract.md) is
informational here: kind and transport semantics are unchanged.

The key words "MUST", "MUST NOT", "SHOULD", and "MAY" are to be interpreted as
described in BCP 14 when, and only when, they appear in all capitals.

## Background / motivation

See the parent case's [motivation](../0206-intent-faithful-evaluator-selection.md#motivation).
In short: first-usable short-circuit discovery hides viable runtimes from
`evaluatorCandidates`, claude authentication is assumed rather than probed, and
the skill's default `harness` precedence silently resolves provider-named
intent ("have Claude evaluate this") to in-session judgment instead of the
independent SDK evaluator the operator may have meant. Evaluator identity is an
evidence-basis input to a rating; selection must be observable and
intent-faithful.

## Requirements

### R1 — probe-all candidate reporting

When the evaluator resolves to `auto`, `qualitymd evaluation run` MUST probe
every built-in agent-runtime candidate (`codex` and `claude`) and report each
one in `evaluatorCandidates` with its readiness evidence — in dry-run and run
receipts — regardless of which candidate is selected.

> Rationale: first-usable short-circuiting left a usable claude runtime absent
> from the candidate list entirely, so the operator had no signal it was
> viable. Probing the remaining candidate costs one executable lookup.
>
> Durable spec: modify `specs/cli/evaluation-run.md` — replace "each considered
> candidate" with an obligation to probe and report every built-in candidate,
> so the short-circuit cannot satisfy the reporting rule vacuously.

### R2 — ordering-decided selection reason

When more than one probed candidate is usable, the `auto` selection reason MUST
name the selected runtime, state that the deterministic discovery order decided
the selection, and name each usable candidate not selected.

> Rationale: with both runtimes ready, "codex was selected" without "claude was
> also usable" reads as "codex was the only option"; the false premise is the
> observability bug.
>
> Durable spec: modify `specs/cli/evaluation-run.md` — extend the
> selection-reason requirement with the usable-but-unselected content.

### R3 — authentication basis is structured

Each reported candidate's readiness data MUST carry its authentication basis —
verified by a non-interactive probe, or assumed because the runtime documents
no such probe — as a structured field in JSON receipts, not only as prose
evidence, and MUST NOT include credential values.

> Rationale: the verified/assumed asymmetry between codex and claude currently
> lives in a prose evidence string; the skill and other consumers cannot
> present it without parsing prose.
>
> Durable spec: modify `specs/cli/evaluation-run.md` — candidate readiness
> reporting carries the authentication basis distinctly from free-form
> evidence.

### R4 — use a documented claude authentication probe

Where the claude agent runtime documents a non-interactive authentication
status probe, `auto` discovery MUST use it to verify authentication before
reporting the candidate usable, exactly as it does for codex.

> Rationale: assumed authentication makes claude's readiness weaker evidence
> than codex's and is part of why codex ranks first; verifying where possible
> narrows the asymmetry honestly. Whether the runtime documents such a probe
> today is an open question for design.
>
> Durable spec: modify `specs/evaluation/agent-evaluators.md` — state the
> documented-probe policy: verify when documented, assume and say so when not.

### R5 — skill disambiguates provider-named intent

When the user's evaluator request names a provider that is ambiguous between
the in-session harness transport and that provider's SDK evaluator (for
example, "have Claude evaluate this" in a Claude harness), the evaluate
workflow MUST resolve the ambiguity with the user before selection, explaining
that `harness` reuses the current session's judgment and authentication while
the SDK evaluator runs a fresh, independent subprocess, and naming the shortest
path to each (an explicit evaluator request now; `evaluation.evaluator` config
for a durable default).

> Rationale: the skill resolved "Claude" to `harness` — the transport that does
> not give an independent evaluation — with no signal that a different
> transport matched the words better. Wanting "Claude Code to evaluate" is
> genuinely ambiguous; only the user can resolve it.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` — add the
> disambiguation obligation to the evaluator-selection section.

### R6 — default-harness explanation names the alternative

When the workflow selects `harness` by default precedence (no explicit request
and no non-`auto` configuration), the transport explanation it gives before the
first evaluation mutation MUST state that judgment runs in-session and name the
independent SDK-evaluator alternative and how to choose it.

> Rationale: the existing explain-the-transport rule was satisfiable by naming
> `harness` alone; an operator cannot correct a default they do not know they
> received.
>
> Durable spec: modify `specs/skills/quality-skill/evaluation.md` — extend the
> transport-explanation requirement; modify
> `specs/skills/quality-skill/workflows/evaluate.md` — align the evaluate
> step's resolve-and-explain wording.

### R7 — deterministic coverage

Deterministic tests MUST cover: both candidates reported when the first is
usable (R1); the selection reason naming a usable-but-unselected candidate
(R2); the structured authentication basis for a verified and an assumed
candidate (R3); and, if a claude probe is adopted, claude selection gated on
that probe (R4).

> Rationale: the removed short-circuit assertions (`candidates` length 1) are
> the regression path back to under-reporting.
>
> Durable spec: none.

## Requirement-set check

Consistent: R1–R4 govern the CLI's discovery and reporting; R5–R6 govern the
skill's intent handling; no requirement reorders discovery or moves precedence,
so none conflicts with the unchanged ordering non-goal. Complete: the
motivation's three gaps (under-reported candidates, assumed authentication,
silent harness default) each map to at least one requirement; the open
probe-availability question is confined to R4's "where documented" condition.
Able to be validated: satisfying R1–R7 makes selection observable (candidates,
reason, authentication basis) and intent-faithful (disambiguation, named
default) without changing who wins — which is the motivation, no more.

## Open questions

- Does the claude agent runtime document a non-interactive authentication
  status probe today? R4 binds only where one is documented; design confirms
  and either wires it or records its absence.

## Durable spec changes

### To add

None

### To modify

- `specs/cli/evaluation-run.md` — probe-all candidate reporting (R1),
  ordering-decided selection reason (R2), structured authentication basis
  (R3).
- `specs/evaluation/agent-evaluators.md` — documented-authentication-probe
  policy (R4).
- `specs/skills/quality-skill/evaluation.md` — provider-named intent
  disambiguation (R5) and default-harness transport explanation (R6).
- `specs/skills/quality-skill/workflows/evaluate.md` — evaluate-step
  resolve-and-explain wording (R6).

### To rename

None

### To delete

None
