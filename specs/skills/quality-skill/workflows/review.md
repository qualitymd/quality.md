---
type: Functional Specification
title: /quality review
description: Behavioral component spec for the /quality review workflow stub.
tags: [skill, quality, review, workflow]
timestamp: 2026-06-27T00:00:00Z
---

# /quality review

`review` is the `/quality` skill workflow that inspects an Evaluation result, the
`QUALITY.md` model, or a specific quality concern before the user decides
whether to act. It implements the shared contracts in the parent
[/quality skill](../quality-skill.md) spec and owns only the review-specific
behavior below.

The runtime procedure lives at
[`skills/quality/workflows/review.md`](../../../../skills/quality/workflows/review.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../../docs/reference/rfc8174.md) when, and only when, they
appear in all capitals.

## Purpose and routing

`review` is selected when the user asks to review an Evaluation result, review
the `QUALITY.md` model, or review a specific quality concern.

The workflow's purpose is judgment and routing. It helps the user understand
evidence, model fit, and likely next action before mutation.

## Focus

`review` **MUST** support these focus values:

- Evaluation result — latest or selected run, report, rating, finding,
  candidate action, evidence limit, or incomplete input;
- model — `QUALITY.md` usefulness, clarity, coverage, assessability, stale
  assumptions, or authoring quality; and
- concern — a named quality concern that is not clearly a model element or
  Evaluation artifact.

When focus is absent or ambiguous, `review` **MUST** infer likely focus from user
text and lifecycle state before asking. When inference is not strong enough,
`review` **MUST** ask a single-select closed-choice focus question with the
recommended focus first and an explicit shortest answer path.

## Mutation surface and artifacts

`review` is read-only by default. It **MUST NOT** edit evaluated source, edit
`QUALITY.md`, write Evaluation records, write the quality changelog, create external
issues, update tooling, or create workflow feedback logs.

## Required flow

Before tool inspection, `review` **MUST** emit the public `/quality` run frame
required by the parent skill contract. The frame **MUST** include the resolved
or provisional focus and **MUST** name the mutation boundary as read-only.

After focus is resolved, `review` **MUST** confirm the focus before deep review
work. Confirmation may be a concise statement when inference is strong, or a
single-select choice when it is not.

For Evaluation-result focus, `review` **SHOULD** inspect the latest reportable
run when no run is named, summarize available ratings, findings, evidence
limits, incomplete inputs, and likely improve focuses.

For model focus, `review` **SHOULD** route through existing authoring and
top-10-check guidance, inspect `QUALITY.md` for usefulness, clarity, coverage,
assessability, and stale assumptions, and report suggested improve focuses
without changing the model.

For concern focus, `review` **SHOULD** inspect the named concern against the
model and available project context, then recommend one next action: scoped
`evaluate`, model-focused `improve`, work-focused `improve`, or stop.

## Completion criteria

`review` is complete when it reports what was reviewed, the evidence limits, the
recommended next action, concrete alternatives when useful, and the read-only
boundary.
The closeout **MUST** use labeled fields for reviewed subject, signal, evidence
limits, recommended next action, alternatives when useful, and not-changed
boundary.
