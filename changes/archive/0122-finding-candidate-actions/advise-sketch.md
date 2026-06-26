---
type: Sketch
title: Advise phase — design sketch
description: Non-binding, forward-looking notes for the future Evaluation Advise phase (option B) that consumes finding-level candidate actions and produces a final recommendation set.
tags: [evaluation, advise, recommendations, exploration]
timestamp: 2026-06-26T00:00:00Z
---

# Advise phase — design sketch

> 🧭 **Non-binding exploration.** This is a *sketch*, not a spec or a design doc.
> It captures ideas and considerations for the future **Advise** phase (option
> **B**) so the harvest layer this change ships (option **A**, the typed
> [candidate actions](spec.md)) is shaped to feed it. None of this is committed
> work; the Advise phase is **deferred** and Evaluation v0 still
> [forbids recommendation generation](../../../../specs/evaluation/evaluation.md).
> When B is built, it gets its own change case and functional spec.

## Where it sits

The format-level Evaluation contract is **Define → Assess and Rate → Analyze →
Advise → Report** ([SPECIFICATION.md](../../../../SPECIFICATION.md), "Advice"). The
Advise phase is named and specced-as-optional but not built. Its job, per the
spec: identify **key gaps** (shortcomings most responsible for held-down
ratings), enumerate **options**, and select **recommendations** with rationale.

Today the only actionable guidance lives ad hoc in the root report's summary /
rating-driver narrative. The Advise phase makes that a durable, addressable,
closure-trackable artifact.

## Core thesis

Three distinct things, kept distinct:

- **Candidate actions** (option A, shipping here) — finding-local, non-binding
  *inputs*. Born where evidence is richest. Noisy and un-prioritized by design.
- **Rating drivers** (already produced by analysis) — the *spine*. They name the
  shortcomings that actually hold a rating down, so they are the natural
  prioritization anchor.
- **Recommendations** (option B) — the synthesized *output*. A small, ranked set,
  each tracing to the rating it would move.

Candidate actions are necessary but **not sufficient**: a final set requires
cross-finding synthesis (dedup, clustering, ranking) that only exists once the
whole analysis is done. The recurring "0 tests" finding that appeared under four
requirements in the motivating evaluation is the canonical example — four
near-duplicate leads must collapse into one prioritized recommendation.

## Approaches under consideration

Compose, don't pick one:

- **Dedicated Advise phase + a recommendation payload kind.** A new agent-written
  kind (working name `EvaluationAdviceResult`) alongside the existing kinds in
  [payload-kinds.md](../../../../specs/evaluation/records/payload-kinds.md), with a
  routine contract, CLI validation, and report rendering. The heavyweight piece;
  the one that matches the named phase.
- **Gap-anchored synthesis (how the phase should work internally).** Synthesize
  from the **key gaps** the analysis already surfaces (rating drivers, held-down
  factor ratings), not from "all findings." Generate options/recommendations only
  for gaps that hold a rating down. Mirrors the spec's *key gaps → options →
  recommendations* exactly, and naturally dedups and bounds the set.
- **Trace recommendations to ratings, not just findings.** Each recommendation
  cites which rating it would move and the from→to it targets, so it carries a
  verifiable value claim and a later re-evaluation can check closure ("did the
  rating actually move?"). This is what lets
  [recommendation follow-up](../../../../specs/skills/quality-skill/recommendation-follow-up.md)
  act on structured artifacts instead of prose.
- **Model options, not just an action.** Per the spec's three-part advice: a gap
  gets 2–3 options (effort / risk / impact), a selected recommendation, and a
  rationale. Probably a `deep`-rigor enrichment, not the base shape.
- **Subagent fan-out for synthesis at `deep` rigor.** Mirror the existing
  assessment fan-out: a per-area recommendation synthesizer returns candidates;
  the orchestrator does the global dedup/prioritize. Keeps cross-cutting
  prioritization central.

A likely composition: a dedicated Advise phase, gap-anchored internally,
consuming candidate actions as inputs, emitting a small ranked recommendation set
where each entry traces to the rating it moves; options/fan-out as later
enrichments.

## Considerations seeded by the harvest layer (A)

- **Per-action identity.** A ships candidate actions with no stable id (addressed
  by finding id + position). If B needs to reference, supersede, or track an
  individual action to closure, a stable per-action id may be worth adding then —
  deferred until there is a consumer.
- **Action shape evolution.** A keeps the candidate action minimal
  (`description` + optional `rationale`). B may want effort/impact/risk hints. Add
  them on the *recommendation*, not the candidate action, to keep assessment
  diagnostic and synthesis prescriptive.
- **Inputs are leads, not promises.** B must be free to reframe, merge, or drop
  candidate actions; they are raw material, never a recommendation B is obliged to
  honor. The harvest layer's "non-binding, finding-local" framing exists to make
  this safe.

## v0 → v1 boundary

Option A ships under v0 as **inert data**: candidate actions persist in `data/`
but never appear in the report or closeout. Turning on the Advise phase — the
moment recommendations are synthesized and presented — is the v0 → v1 flip. The
relevant v0 prohibitions to lift/retarget at that point:

- [`specs/evaluation/evaluation.md`](../../../../specs/evaluation/evaluation.md) —
  "Recommendation generation MUST NOT be part of the v0 evaluation protocol" and
  the `Deferred: recommendation generation` scope note.
- [`specs/skills/quality-skill/reporting.md`](../../../../specs/skills/quality-skill/reporting.md)
  — "Evaluation v0 MUST NOT present generated recommendations."
- [`skills/quality/workflows/evaluate.md`](../../../../skills/quality/workflows/evaluate.md)
  — "Do not generate recommendations in Evaluation v0."

## Open questions for the B change case

- Payload kind name and shape: one `EvaluationAdviceResult` per run, or per Area?
- Where does Advise run in the orchestration — after root Area analysis, before
  report build?
- How does the report render recommendations distinctly from rating drivers so the
  two are not conflated?
- How does a re-evaluation reconcile prior recommendations (closed / still-open /
  superseded) — and does that reuse the existing recommendation-follow-up surface?
