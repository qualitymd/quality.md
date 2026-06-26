---
type: Functional Specification
title: /quality recommendation follow-up
description: Non-mode workflow for applying or handing off /quality evaluation recommendations.
tags: [skill, quality, recommendation]
timestamp: 2026-06-22T00:00:00Z
---

# /quality recommendation follow-up

This spec owns the `/quality` skill's non-mode recommendation follow-up
workflow: how the skill helps users act on evaluation recommendations after
`evaluate` produces them or read-only orientation routes to them. It composes the shared
contracts in the parent [/quality skill](quality-skill.md) spec and the
recommendation artifact contract in [/quality reporting](reporting.md).

This document uses BCP 14 keywords only for testable conformance requirements.
The key words "MUST", "MUST NOT", and "SHOULD" are to be interpreted as
described in [RFC 2119](../../../docs/reference/rfc2119.md) and
[RFC 8174](../../../docs/reference/rfc8174.md) when, and only when, they appear
in all capitals.

## Purpose and routing

Recommendation follow-up is selected when the user asks to apply, act on,
improve from, or hand off an active evaluation recommendation. It is not a
runtime mode and does not replace `evaluate`: recommendations remain outputs of
evaluation.

The workflow's purpose is to turn a selected recommendation into one of two
productive outcomes: a confirmed local apply or an issue-tracker handoff.

## Outcomes

Recommendation follow-up **MUST** offer only two explicit outcomes:

1. apply a confirmed recommendation option now;
2. hand off the recommendation to an issue tracker.

When the user has not already chosen one of those outcomes, the workflow
**MUST** present them as numbered options and include an `Answer` line. When
recommendation evidence supports one path, that path **SHOULD** be option `1`;
otherwise the workflow **MUST NOT** invent a recommendation solely to satisfy
the numbered prompt shape.

The workflow **MUST NOT** present `defer`, `skip`, or `keep open` as formal
follow-up options. If the user does not choose apply or handoff, the workflow
**MUST** stop without mutating evaluated source, `QUALITY.md`, the quality log,
or external systems.

## Apply now

Before applying a recommendation, the skill **MUST** present a decision brief
that names the recommendation, selected option, artifact class being changed,
evidence or reason, risk when relevant, done criterion, and verification path.
The primary apply question or call to action **MUST** be visually emphasized,
and the decision brief **MUST** keep the changed artifacts, evidence/reason,
recommended option, alternatives, and done criterion in a consistent, scannable
shape.

The skill **MUST NOT** edit evaluated source files, edit `QUALITY.md`, or write
the quality log until the user explicitly confirms the recommendation option and
mutation surface. It **MUST NOT** treat an obvious or recommended fix as consent.

After applying a confirmed option, the skill **SHOULD** verify the done criterion
with the narrowest useful evidence. When the done criterion is rating-bound or
depends on the QUALITY.md model, the skill **SHOULD** create a new numbered
evaluation run for the affected scope and compare the new evidence to the
recommendation's done criterion.

The result report **MUST** state the recommendation, outcome, applied option,
changed artifacts, verification performed, rating movement when known, remaining
gaps or limits, and next action. It **SHOULD** include `Not done` when the
mutation or verification boundary matters, such as no evaluation rerun, no issue
creation, no model change, or no quality-log entry. It **MUST** use the shared
agent-mediated UX contract: status first, with scannable labels and a clear next
action when work remains. If verification is incomplete, the result **MUST** be
labeled limited rather than fully confirmed.

When a confirmed apply changes the QUALITY.md model, the skill **MUST** write one
quality-log entry for the coherent model change, cross-linking the source
evaluation run and recommendation when applicable. Evaluated-source fixes that do
not change the model **MUST NOT** write the quality log.

Before applying a model-changing recommendation, the skill **MUST** read the
authoring entry guide, the routed authoring sub-guide for the model element being
changed, and the quality-log authoring sub-guide.

## Issue-tracker handoff

Issue-tracker handoff **MUST** produce issue-ready content that includes the
recommendation ID and title, source evaluation run, affected area/factor or
requirement, current rating when known, target or done criterion, evidence
summary with locators, suggested implementation option, verification path, and
links or paths to the generated report and recommendation artifact.

Creating an external issue **MUST** require explicit user confirmation,
available issue-tracker tooling, and a decision brief that names the external
artifact to create, the local artifacts that will not change, the evidence or
reason, the recommended option, alternatives, and verification. If tooling is
unavailable or the user has not confirmed external creation, the skill **MUST**
stop after producing issue-ready text.

Issue-tracker handoff **MUST NOT** mutate evaluated source, `QUALITY.md`, or the
quality log.
