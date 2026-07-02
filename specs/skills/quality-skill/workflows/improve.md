---
type: Functional Specification
title: /quality improve
description: Behavioral component spec for the /quality improve workflow stub.
tags: [skill, quality, improve, workflow]
timestamp: 2026-06-27T00:00:00Z
---

# /quality improve

`improve` is the `/quality` skill workflow that acts on quality judgment after
focus and mutation surface are confirmed. It implements the shared contracts in
the parent [/quality skill](../quality-skill.md) spec and owns only the
improve-specific behavior below.

The runtime procedure lives at
[`skills/quality/workflows/improve.md`](../../../../skills/quality/workflows/improve.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`improve` is selected when the user asks to improve from an evaluation result,
improve the model, improve a specific quality concern, apply or hand off a
compatible recommendation artifact, or otherwise act on quality judgment.

The workflow's purpose is safe action routing. It starts read-only, confirms the
focus and mutation surface, then delegates to an existing safe route or stops at
an explicit stub boundary.

## Focus

`improve` **MUST** support these focus values:

- Evaluation result — latest or selected run, finding, candidate action, evidence
  limit, or model gap;
- model — `QUALITY.md` body, structure, areas, factors, requirements, rating
  scale, coverage, clarity, or assessability;
- concern — a named quality concern that is not clearly an evaluation artifact,
  recommendation artifact, or model element; and
- recommendation — a compatible existing recommendation artifact or explicit
  recommendation ID.

When focus is absent or ambiguous, `improve` **MUST** infer likely focus from
user text and lifecycle state before asking. When inference is not strong enough,
`improve` **MUST** ask a single-select closed-choice focus question with the
recommended focus first and an explicit shortest answer path.

## Mutation surface and artifacts

`improve` starts read-only. It **MUST** confirm both focus and mutation surface
before editing evaluated source, editing `QUALITY.md`, writing the quality changelog,
creating an external issue, or updating tooling.

`improve` **MUST NOT** create numbered evaluation records or reports itself. If
verification requires a fresh rating, it routes to `evaluate` for the affected
scope.

## Required flow

Before tool inspection, `improve` **MUST** emit the public `/quality` run frame
required by the parent skill contract. The frame **MUST** include the resolved
or provisional focus and **MUST** name the mutation surface as `read-only until
confirmed` when unresolved.

After focus is resolved, `improve` **MUST** identify the likely mutation surface
before acting:

- model focus maps to `QUALITY.md` and, for meaningful model changes, the quality
  log;
- recommendation focus maps through recommendation follow-up to evaluated source,
  `QUALITY.md`, quality changelog, or external issue;
- Evaluation-result focus remains read-only until a finding, candidate action,
  model gap, or work target is selected; and
- concern focus maps to evaluated source, `QUALITY.md`, external issue, or no
  mutation yet.

For model focus, `improve` **MUST** delegate to the parent spec's direct
model-authoring route after focus and mutation surface are confirmed.

For recommendation focus, `improve` **MUST** delegate to
[recommendation follow-up](../recommendation-follow-up.md) when a compatible
recommendation artifact exists. It **MUST NOT** synthesize a recommendation when
no compatible artifact exists.

For evaluation-result focus without a compatible recommendation artifact,
`improve` **SHOULD** help select a finding, candidate action, model gap, or work
target, then either route to model-focused/direct improvement, recommend a scoped
`evaluate`, or stop with the deferred deeper-workflow boundary.

For concern focus, `improve` **SHOULD** confirm whether the user wants to improve
the work, improve `QUALITY.md`, hand off an issue, or first run/review an
evaluation.

## Completion criteria

`improve` is complete when it reports the confirmed focus, mutation surface,
changed artifacts when any, verification performed, remaining limits, and what
was not changed. If it stops at the stub boundary, it must name which deeper
improve behavior is deferred and offer the nearest runnable next workflow.
The closeout **MUST** use labeled fields for focus, changed artifacts,
verification, remaining limits, not-changed boundary, and next action.
