---
type: Design Doc
title: Evaluation orchestration overhaul — design
description: How the /quality skill prompt and durable specs deliver one best-quality evaluate workflow with parallel collection and an always-on two-pronged QC phase.
status: Done
tags: [skill, evaluation, orchestration, qc]
timestamp: 2026-06-26T00:00:00Z
---

# Evaluation orchestration overhaul — design

## Context

This answers the [0129 functional spec](spec.md): collapse evaluation rigor into
one best-quality workflow, make exhaustive coverage and a two-pronged QC phase
the mandatory contract, and use subagent fan-out as the default execution
strategy where available. The change is entirely skill-prompt + durable-spec; no
Go code moves. The work lives in `skills/quality/` (the runtime prompt) and
`specs/skills/quality-skill/` (the durable contract it implements).

## Approach

The current evaluate procedure is already a linear pipeline of CLI-mediated
phases (frame → preflight/lint → history → create run → frame payloads →
per-Area/Requirement assessment → Factor/Area roll-up → binding re-check →
batch-write → report build). This design keeps that spine and changes three
things: it deletes the rigor dial, it restructures the per-Requirement assessment
and the binding re-check into a **collect → QC → roll-up** shape, and it makes the
collect and QC phases fan out.

### Phase shape

```
Frame  →  Preflight (lint, history, create run, model-snapshot IDs)
            │
            ▼
   ┌──────────────────────────────────────────────────────────────┐
   │  COLLECT  (exhaustive, fan-out per-Area / per-Requirement)     │
   │    structured findings only; every locator file:line/search   │
   └──────────────────────────────────────────────────────────────┘
            │
            ▼
   ┌──────────────────────────────────────────────────────────────┐
   │  QC  (always on, two prongs, parallel where supported)        │
   │    verify ∥ completeness-sweep                                 │
   │       verify: re-run cmd/search for every roll-up-binding +    │
   │               low-confidence finding                          │
   │       sweep:  coverage ledger · re-examine quiet zones ·       │
   │               escalate thin-evidence                          │
   └──────────────────────────────────────────────────────────────┘
        │ new findings ▲                       │ converged / bound hit
        └──── re-collect (bounded) ────────────┘
            │
            ▼
   ROLL-UP  (orchestrator-only: Factor → Area → headline ratings)
            │
            ▼
   Batch-write payloads  →  evaluation status  →  report build
```

The orchestrator owns the loop, the roll-up, and the rating-binding re-check.
Subagents only ever do collection and verification legwork and hand back
structured findings (the `RequirementAssessmentResult` / finding payloads the CLI
already defines), never authoritative ratings — preserving the existing invariant
verbatim (spec R10).

### Fan-out and fallback

Fan-out is an execution detail the prompt expresses as: *if a subagent capability
is present, dispatch independent per-Area (or per-Requirement) collection and the
two QC prongs as concurrent subagents; otherwise do the identical work serially.*
The skill does not name a specific subagent tool — it tests for the capability,
matching how the agent-mediated UX contract already tells the skill to choose the
richest affordance it has and fall back when absent. Each subagent prompt is
seeded with the resolved scope, the relevant snapshot Requirement IDs, the
secret-handling rule, the source-as-data rule, and "return structured findings
only" (spec R9). Because all authoritative writes go through one batched
`qualitymd evaluation data set` at the end, parallel collectors never contend on
the run's data store.

### The two QC prongs

The verify prong generalizes today's headline-only re-check (it already exists in
the procedure as step 16) to *every* roll-up-binding finding plus low-confidence
findings, and keeps the re-*run* (not re-read) rule and its rationale.

The completeness sweep is the genuinely new machinery. It is three concrete
checks the orchestrator (or a per-Area sweep subagent) performs:

1. **Coverage ledger** — diff the in-scope Requirement ID set (from the model
   snapshot, already queried in the create-run step) against the Requirements
   that reached a terminal evidentiary state. Any ID with neither a rating nor a
   reasoned *not assessed* fails the sweep and re-enters collection.
2. **Quiet-zone re-examination** — any Area/Requirement whose first pass yielded
   only `strength` findings or none gets a fresh adversarial "find the gap/risk"
   look. This is where missed findings hide.
3. **Thin-evidence escalation** — any Requirement rated on a single weak
   observation gets an independent second look before its rating stands.

Anything these surface re-enters collection and is then verified by the verify
prong before it can bind a rating (spec R14).

### Convergence

The loop terminates on the converged state (no new in-scope findings *and* every
in-scope Requirement terminal), bounded by a fixed cap of **two** re-collection
rounds. If the cap is hit first, the orchestrator proceeds to roll-up and lists
every still-unexamined or unresolved zone as an explicit limitation in the report
— never a silent drop (spec R15, R16). Two rounds is the design's choice, not a
spec conformance point; it is large enough to absorb a sweep's discoveries and
re-verify them, small enough to bound cost.

## Spec response

- **Rigor removal (R1–R3)** — delete the Rigor Levels table, the `Rigor:` frame
  field, the `/quality evaluate deep` invocation, and the feedback-log `rigor:`
  field from the runtime files; mirror the deletions in the durable specs. Scope
  resolution is untouched.
- **Exhaustive coverage (R4–R6)** — the collect phase reads the whole in-scope
  source; the coverage ledger (QC check 1) enforces the terminal-state rule; the
  existing "report what was not assessed" obligation moves from the rigor table
  into the coverage contract prose.
- **Execution strategy (R7–R9)** — capability-gated fan-out with an identical
  serial fallback; subagent prompt contract preserved and pointed at snapshot IDs.
- **Orchestrator-owned (R10)** — unchanged invariant, restated against the new
  phase names.
- **QC phase (R11–R14)** — the two-pronged phase above; both prongs always run.
- **Convergence (R15–R16)** — bounded loop with honest residue disclosure.

## Alternatives

- **Keep a cost/budget cap as the surviving knob.** Rejected (and the user
  confirmed "scope is the only dial"): a budget cap re-introduces a quality dial
  by another name and competes with scope, which already bounds cost honestly.
- **Universal verify (re-check every finding).** Rejected: even cheap-in-parallel,
  re-litigating obviously-confirmed, non-binding findings spends agent budget
  where neither failure mode lives. Effort is targeted at decision-binding
  verification and at coverage gaps instead.
- **Verify-only QC (no completeness sweep).** Rejected: it guards false positives
  only and leaves the user's primary worry — missed findings — unaddressed. The
  sweep is the higher-value prong precisely because "found nothing" is the
  likeliest place a real gap was missed.
- **Unbounded loop-until-dry.** Rejected: a large or adversarial source may never
  go dry; a fixed bound plus disclosed residue is safer than both an infinite
  loop and a silent cap.
- **Split QC into a child spec now.** Deferred: QC is entangled with collection
  and roll-up in the shared evaluation contract; a child spec is revisitable if
  `evaluation.md` grows, but is not earned today.

## Trade-offs & risks

- **More tokens/agents per run.** A best-quality default costs more than the old
  `quick`. Mitigated by scope narrowing (the user's lever) and by targeting QC at
  binding/quiet zones rather than universal re-verification.
- **Sweep quality depends on the prompt.** The quiet-zone re-examination is only
  as good as its adversarial framing; a weak prompt regresses to "looks clean."
  The spec pins the *what* (R13); the prompt must carry a genuinely skeptical
  lens, which the design calls out explicitly.
- **Capability detection.** Misjudging whether a subagent tool exists would push a
  run onto the serial path — slower but not wrong, since coverage and QC are
  identical either way. The failure mode is latency, not correctness.

## Open questions

- The re-collection bound is set to two rounds here; if real runs show sweeps
  routinely surfacing fresh findings in round two, revisit whether the cap or the
  collect-phase thoroughness needs tuning. Tracked as a prompt heuristic, not a
  spec change.
